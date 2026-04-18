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
	db      *sql.DB
	queries *db.Queries
}

func NewRepository(sqlDB *sql.DB, queries *db.Queries) *Repository {
	return &Repository{
		db:      sqlDB,
		queries: queries,
	}
}

func (r *Repository) FindSavedPlayer(
	ctx context.Context,
	input usecase.PlayerSearchInput,
) (*usecase.PlayerSearchResult, error) {
	return nil, nil
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
