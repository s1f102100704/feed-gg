package playersearch

import (
	"database/sql"
	"fmt"
	"time"

	db "feed-gg/backend/internal/infrastructure/db/sqlc"
	"feed-gg/backend/internal/infrastructure/riot"
)

type currentRankEntry struct {
	queueType  string
	rank       *riot.RankedQueue
	recordedAt time.Time
}

func currentRankEntries(fetched *riot.PlayerProfile, recordedAt time.Time) []currentRankEntry {
	entries := make([]currentRankEntry, 0, 2)
	if fetched == nil {
		return entries
	}
	if fetched.SoloRank != nil {
		entries = append(entries, currentRankEntry{
			queueType:  "RANKED_SOLO_5x5",
			rank:       fetched.SoloRank,
			recordedAt: recordedAt,
		})
	}
	if fetched.FlexRank != nil {
		entries = append(entries, currentRankEntry{
			queueType:  "RANKED_FLEX_SR",
			rank:       fetched.FlexRank,
			recordedAt: recordedAt,
		})
	}
	return entries
}

func newPlayerProfileParams(regionID int16, fetched *riot.PlayerProfile) (db.UpsertPlayerProfileParams, error) {
	if fetched == nil {
		return db.UpsertPlayerProfileParams{}, ErrFetchedPlayerRequired
	}

	profileIconID, err := toNullInt32(fetched.ProfileIconID)
	if err != nil {
		return db.UpsertPlayerProfileParams{}, fmt.Errorf("convert profile icon id: %w", err)
	}

	return db.UpsertPlayerProfileParams{
		Puuid:         fetched.PUUID,
		GameName:      nullString(fetched.GameName),
		TagLine:       nullString(fetched.TagLine),
		RegionID:      regionID,
		ProfileIconID: profileIconID,
		SummonerLevel: sql.NullInt64{Int64: fetched.SummonerLevel, Valid: true},
		RevisionDate:  nullTimeFromMillis(fetched.RevisionDate),
	}, nil
}

func newMatchHistoryParams(regionID int16, match riot.MatchSummary) (db.UpsertMatchHistoryParams, error) {
	queueID, err := toInt32(match.QueueID)
	if err != nil {
		return db.UpsertMatchHistoryParams{}, fmt.Errorf("convert queue id for match %q: %w", match.MatchID, err)
	}
	durationSeconds, err := toInt32(match.DurationSeconds)
	if err != nil {
		return db.UpsertMatchHistoryParams{}, fmt.Errorf(
			"convert duration seconds for match %q: %w",
			match.MatchID,
			err,
		)
	}

	return db.UpsertMatchHistoryParams{
		MatchID:         match.MatchID,
		RegionID:        regionID,
		QueueID:         queueID,
		GameMode:        match.GameMode,
		GameVersion:     match.GameVersion,
		DurationSeconds: durationSeconds,
		PlayedAt:        timeFromMillis(match.PlayedAt),
	}, nil
}

func newParticipantPlayerParams(regionID int16, participant riot.MatchParticipantSummary) db.UpsertParticipantPlayerParams {
	return db.UpsertParticipantPlayerParams{
		Puuid:    participant.PUUID,
		GameName: nullString(participant.GameName),
		TagLine:  nullString(participant.TagLine),
		RegionID: regionID,
	}
}

func newMatchParticipantParams(
	matchHistoryID int64,
	playerID int64,
	participant riot.MatchParticipantSummary,
) db.UpsertMatchParticipantParams {
	return db.UpsertMatchParticipantParams{
		MatchHistoryID:   matchHistoryID,
		PlayerID:         playerID,
		GameNameSnapshot: nullString(participant.GameName),
		TagLineSnapshot:  nullString(participant.TagLine),
		TeamID:           int16(participant.TeamID),
		ChampionName:     participant.ChampionName,
		TeamPosition:     sql.NullString{},
		Role:             nullString(participant.Role),
		Win:              participant.Win,
		Kills:            int32(participant.Kills),
		Deaths:           int32(participant.Deaths),
		Assists:          int32(participant.Assists),
		SummonerSpell1ID: int32(participant.SummonerSpell1ID),
		SummonerSpell2ID: int32(participant.SummonerSpell2ID),
	}
}
