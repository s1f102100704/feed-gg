package playersearch

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	db "feed-gg/backend/internal/infrastructure/db/sqlc"
	"feed-gg/backend/internal/infrastructure/riot"
	"feed-gg/backend/internal/usecase"
)

var (
	ErrRepositoryNotConfigured = errors.New("player search repository is not configured")
	ErrFetchedPlayerRequired   = errors.New("fetched player is required")
)

type Repository struct {
	db                  *sql.DB
	queries             *db.Queries
	profileIconResolver ProfileIconURLResolver
}

func NewRepository(
	sqlDB *sql.DB,
	queries *db.Queries,
	profileIconResolver ProfileIconURLResolver,
) *Repository {
	return &Repository{
		db:                  sqlDB,
		queries:             queries,
		profileIconResolver: profileIconResolver,
	}
}

func (r *Repository) FindSavedPlayer(
	ctx context.Context,
	input usecase.PlayerSearchInput,
) (*usecase.PlayerSearchResult, error) {
	if r.queries == nil {
		return nil, ErrRepositoryNotConfigured
	}

	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}

	region, err := r.queries.GetRegionByName(ctx, input.Region)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("resolve region %q: %w", input.Region, err)
	}

	playerKey, err := r.queries.GetSavedPlayerKeyByRiotID(ctx, db.GetSavedPlayerKeyByRiotIDParams{
		RegionID: region.ID,
		GameName: nullString(input.GameName),
		TagLine:  nullString(input.TagLine),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("lookup saved player key: %w", err)
	}

	playerRow, err := r.queries.GetSavedPlayerByPuuid(ctx, playerKey.Puuid)
	if err != nil {
		return nil, fmt.Errorf("load saved player %q: %w", playerKey.Puuid, err)
	}

	rankRows, err := r.queries.ListPlayerCurrentRanksByPlayerID(ctx, playerRow.PlayerID)
	if err != nil {
		return nil, fmt.Errorf("load player ranks for player_id=%d: %w", playerRow.PlayerID, err)
	}

	matchRows, err := r.queries.ListRecentMatchHistoriesByPlayerID(ctx, db.ListRecentMatchHistoriesByPlayerIDParams{
		PlayerID:   playerRow.PlayerID,
		LimitCount: savedPlayerRecentMatchLimit,
	})
	if err != nil {
		return nil, fmt.Errorf("load recent matches for player_id=%d: %w", playerRow.PlayerID, err)
	}

	matchParticipants := make(map[int64][]db.ListMatchParticipantsByMatchHistoryIDRow, len(matchRows))
	for _, matchRow := range matchRows {
		participants, err := r.queries.ListMatchParticipantsByMatchHistoryID(ctx, matchRow.MatchHistoryID)
		if err != nil {
			return nil, fmt.Errorf(
				"load participants for match_history_id=%d: %w",
				matchRow.MatchHistoryID,
				err,
			)
		}
		matchParticipants[matchRow.MatchHistoryID] = participants
	}

	return mapSavedPlayerSearchResult(
		ctx,
		r.profileIconResolver,
		playerRow,
		rankRows,
		matchRows,
		matchParticipants,
	), nil
}

func (r *Repository) SaveFetchedPlayer(ctx context.Context, fetched *riot.PlayerProfile) (err error) {
	if fetched == nil {
		return ErrFetchedPlayerRequired
	}
	if r.db == nil || r.queries == nil {
		return ErrRepositoryNotConfigured
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin player search save transaction: %w", err)
	}
	defer func() {
		if err == nil {
			return
		}
		_ = tx.Rollback()
	}()

	q := r.queries.WithTx(tx)
	savedAt := nowUTC()

	player, region, err := r.saveRootPlayer(ctx, q, fetched)
	if err != nil {
		return err
	}

	if err := r.replaceCurrentRanks(ctx, q, player.ID, fetched, savedAt); err != nil {
		return err
	}

	if err := r.saveMatches(ctx, q, region.ID, fetched.Matches); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit player search save transaction: %w", err)
	}

	return nil
}

func (r *Repository) saveRootPlayer(
	ctx context.Context,
	q *db.Queries,
	fetched *riot.PlayerProfile,
) (db.Player, db.Region, error) {
	region, err := q.GetRegionByName(ctx, fetched.Region)
	if err != nil {
		return db.Player{}, db.Region{}, fmt.Errorf("resolve region %q: %w", fetched.Region, err)
	}

	playerParams, err := newPlayerProfileParams(region.ID, fetched)
	if err != nil {
		return db.Player{}, db.Region{}, err
	}
	player, err := q.UpsertPlayerProfile(ctx, playerParams)
	if err != nil {
		return db.Player{}, db.Region{}, fmt.Errorf("upsert player profile: %w", err)
	}

	return player, region, nil
}

func (r *Repository) replaceCurrentRanks(
	ctx context.Context,
	q *db.Queries,
	playerID int64,
	fetched *riot.PlayerProfile,
	savedAt time.Time,
) error {
	if err := q.DeletePlayerCurrentRanksByPlayerID(ctx, playerID); err != nil {
		return fmt.Errorf("delete current ranks: %w", err)
	}

	for _, rankEntry := range currentRankEntries(fetched, savedAt) {
		playerRank, err := q.GetPlayerRankByTierDivision(ctx, db.GetPlayerRankByTierDivisionParams{
			Tier:     rankEntry.rank.Tier,
			Division: rankEntry.rank.Rank,
		})
		if err != nil {
			return fmt.Errorf(
				"resolve player rank %s %s for queue %s: %w",
				rankEntry.rank.Tier,
				rankEntry.rank.Rank,
				rankEntry.queueType,
				err,
			)
		}

		if err := q.UpsertPlayerCurrentRank(ctx, db.UpsertPlayerCurrentRankParams{
			PlayerID:      playerID,
			QueueType:     rankEntry.queueType,
			PlayerRanksID: sql.NullInt16{Int16: playerRank.ID, Valid: true},
			LeaguePoints:  sql.NullInt32{Int32: int32(rankEntry.rank.LeaguePoints), Valid: true},
			Wins:          sql.NullInt32{Int32: int32(rankEntry.rank.Wins), Valid: true},
			Losses:        sql.NullInt32{Int32: int32(rankEntry.rank.Losses), Valid: true},
			RecordedAt:    rankEntry.recordedAt,
		}); err != nil {
			return fmt.Errorf("upsert current rank for queue %s: %w", rankEntry.queueType, err)
		}
	}

	return nil
}

func (r *Repository) saveMatches(
	ctx context.Context,
	q *db.Queries,
	regionID int16,
	matches []riot.MatchSummary,
) error {
	for _, match := range matches {
		if err := r.saveMatch(ctx, q, regionID, match); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) saveMatch(
	ctx context.Context,
	q *db.Queries,
	regionID int16,
	match riot.MatchSummary,
) error {
	matchHistoryParams, err := newMatchHistoryParams(regionID, match)
	if err != nil {
		return err
	}
	matchHistory, err := q.UpsertMatchHistory(ctx, matchHistoryParams)
	if err != nil {
		return fmt.Errorf("upsert match history %q: %w", match.MatchID, err)
	}

	for _, participant := range match.Participants {
		if err := r.saveMatchParticipant(ctx, q, regionID, match.MatchID, matchHistory.ID, participant); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) saveMatchParticipant(
	ctx context.Context,
	q *db.Queries,
	regionID int16,
	matchID string,
	matchHistoryID int64,
	participant riot.MatchParticipantSummary,
) error {
	participantPlayer, err := q.UpsertParticipantPlayer(ctx, newParticipantPlayerParams(regionID, participant))
	if err != nil {
		return fmt.Errorf("upsert participant player %q: %w", participant.PUUID, err)
	}

	if err := q.UpsertMatchParticipant(
		ctx,
		newMatchParticipantParams(matchHistoryID, participantPlayer.ID, participant),
	); err != nil {
		return fmt.Errorf(
			"upsert match participant match=%q puuid=%q: %w",
			matchID,
			participant.PUUID,
			err,
		)
	}

	return nil
}

var _ usecase.PlayerSearchRepository = (*Repository)(nil)
