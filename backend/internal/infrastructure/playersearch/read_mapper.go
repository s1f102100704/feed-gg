package playersearch

import (
	"context"
	"database/sql"

	db "feed-gg/backend/internal/infrastructure/db/sqlc"
	"feed-gg/backend/internal/usecase"
)

const savedPlayerRecentMatchLimit int32 = 3

type ProfileIconURLResolver interface {
	ProfileIconURL(ctx context.Context, profileIconID int) (string, error)
}

func mapSavedPlayerSearchResult(
	ctx context.Context,
	profileIconResolver ProfileIconURLResolver,
	player db.GetSavedPlayerByPuuidRow,
	ranks []db.ListPlayerCurrentRanksByPlayerIDRow,
	matches []db.ListRecentMatchHistoriesByPlayerIDRow,
	matchParticipants map[int64][]db.ListMatchParticipantsByMatchHistoryIDRow,
) *usecase.PlayerSearchResult {
	return &usecase.PlayerSearchResult{
		Region:         player.RegionName,
		PUUID:          player.Puuid,
		GameName:       player.GameName.String,
		TagLine:        player.TagLine.String,
		SummonerLevel:  nullInt64Value(player.SummonerLevel),
		ProfileIconID:  nullInt32AsInt(player.ProfileIconID),
		ProfileIconURL: savedProfileIconURL(ctx, profileIconResolver, player.ProfileIconID),
		RevisionDate:   nullTimeAsUnixMillis(player.RevisionDate),
		SoloRank:       savedRankByQueue(ranks, "RANKED_SOLO_5x5"),
		FlexRank:       savedRankByQueue(ranks, "RANKED_FLEX_SR"),
		Matches:        mapSavedMatchSummaries(matches, matchParticipants),
	}
}

func savedProfileIconURL(
	ctx context.Context,
	profileIconResolver ProfileIconURLResolver,
	profileIconID sql.NullInt32,
) string {
	if profileIconResolver == nil || !profileIconID.Valid {
		return ""
	}

	iconURL, err := profileIconResolver.ProfileIconURL(ctx, int(profileIconID.Int32))
	if err != nil {
		return ""
	}

	return iconURL
}

func savedRankByQueue(
	ranks []db.ListPlayerCurrentRanksByPlayerIDRow,
	queueType string,
) *usecase.RankedQueue {
	for _, rank := range ranks {
		if rank.QueueType != queueType {
			continue
		}

		return &usecase.RankedQueue{
			Tier:         rank.Tier,
			Rank:         rank.Division,
			LeaguePoints: nullInt32AsInt(rank.LeaguePoints),
			Wins:         nullInt32AsInt(rank.Wins),
			Losses:       nullInt32AsInt(rank.Losses),
		}
	}

	return nil
}

func mapSavedMatchSummaries(
	matches []db.ListRecentMatchHistoriesByPlayerIDRow,
	matchParticipants map[int64][]db.ListMatchParticipantsByMatchHistoryIDRow,
) []usecase.MatchSummary {
	summaries := make([]usecase.MatchSummary, 0, len(matches))
	for _, match := range matches {
		summaries = append(summaries, usecase.MatchSummary{
			MatchID:          match.MatchID,
			PlayedAt:         match.PlayedAt.UTC().UnixMilli(),
			GameVersion:      match.GameVersion,
			GameMode:         match.GameMode,
			QueueID:          int(match.QueueID),
			TeamID:           int(match.TeamID),
			ChampionName:     match.ChampionName,
			Role:             savedRole(match.TeamPosition, match.Role),
			Win:              match.Win,
			Kills:            int(match.Kills),
			Deaths:           int(match.Deaths),
			Assists:          int(match.Assists),
			SummonerSpell1ID: int(match.SummonerSpell1ID),
			SummonerSpell2ID: int(match.SummonerSpell2ID),
			DurationSeconds:  int64(match.DurationSeconds),
			Participants:     mapSavedMatchParticipants(matchParticipants[match.MatchHistoryID]),
		})
	}
	return summaries
}

func mapSavedMatchParticipants(
	participants []db.ListMatchParticipantsByMatchHistoryIDRow,
) []usecase.MatchParticipant {
	summaries := make([]usecase.MatchParticipant, 0, len(participants))
	for _, participant := range participants {
		summaries = append(summaries, usecase.MatchParticipant{
			PUUID:            participant.Puuid,
			GameName:         participant.GameNameSnapshot.String,
			TagLine:          participant.TagLineSnapshot.String,
			TeamID:           int(participant.TeamID),
			ChampionName:     participant.ChampionName,
			Role:             savedRole(participant.TeamPosition, participant.Role),
			Win:              participant.Win,
			Kills:            int(participant.Kills),
			Deaths:           int(participant.Deaths),
			Assists:          int(participant.Assists),
			SummonerSpell1ID: int(participant.SummonerSpell1ID),
			SummonerSpell2ID: int(participant.SummonerSpell2ID),
		})
	}
	return summaries
}

func savedRole(teamPosition sql.NullString, role sql.NullString) string {
	if teamPosition.Valid && teamPosition.String != "" {
		return teamPosition.String
	}
	if role.Valid && role.String != "" {
		return role.String
	}
	return "UNKNOWN"
}

func nullInt32AsInt(value sql.NullInt32) int {
	if !value.Valid {
		return 0
	}
	return int(value.Int32)
}

func nullInt64Value(value sql.NullInt64) int64 {
	if !value.Valid {
		return 0
	}
	return value.Int64
}

func nullTimeAsUnixMillis(value sql.NullTime) int64 {
	if !value.Valid {
		return 0
	}
	return value.Time.UTC().UnixMilli()
}
