package playersearch

import (
	"errors"
	"testing"
	"time"

	db "feed-gg/backend/internal/infrastructure/db/sqlc"
	"feed-gg/backend/internal/infrastructure/riot"
)

func TestRepositorySaveFetchedPlayerRequiresFetchedPlayer(t *testing.T) {
	t.Parallel()

	repo := &Repository{}

	err := repo.SaveFetchedPlayer(t.Context(), nil)
	if !errors.Is(err, ErrFetchedPlayerRequired) {
		t.Fatalf("err = %v, want ErrFetchedPlayerRequired", err)
	}
}

func TestRepositorySaveFetchedPlayerRequiresConfiguration(t *testing.T) {
	t.Parallel()

	repo := &Repository{}

	err := repo.SaveFetchedPlayer(t.Context(), &riot.PlayerProfile{})
	if !errors.Is(err, ErrRepositoryNotConfigured) {
		t.Fatalf("err = %v, want ErrRepositoryNotConfigured", err)
	}
}

func TestCurrentRankEntries(t *testing.T) {
	t.Parallel()

	recordedAt := time.Unix(1710000000, 0).UTC()
	fetched := &riot.PlayerProfile{
		SoloRank: &riot.RankedQueue{Tier: "GOLD", Rank: "II"},
		FlexRank: &riot.RankedQueue{Tier: "SILVER", Rank: "I"},
	}

	entries := currentRankEntries(fetched, recordedAt)
	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2", len(entries))
	}
	if entries[0].queueType != "RANKED_SOLO_5x5" {
		t.Fatalf("entries[0].queueType = %q, want RANKED_SOLO_5x5", entries[0].queueType)
	}
	if entries[1].queueType != "RANKED_FLEX_SR" {
		t.Fatalf("entries[1].queueType = %q, want RANKED_FLEX_SR", entries[1].queueType)
	}
	if !entries[0].recordedAt.Equal(recordedAt) || !entries[1].recordedAt.Equal(recordedAt) {
		t.Fatalf("recordedAt = %+v, want %v", entries, recordedAt)
	}
}

func TestNewPlayerProfileParams(t *testing.T) {
	t.Parallel()

	params, err := newPlayerProfileParams(1, &riot.PlayerProfile{
		PUUID:         "test-puuid",
		GameName:      "hide on bush",
		TagLine:       "KR1",
		SummonerLevel: 999,
		ProfileIconID: 1234,
		RevisionDate:  1710000000000,
	})
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	if params.Puuid != "test-puuid" {
		t.Fatalf("params.Puuid = %q, want test-puuid", params.Puuid)
	}
	if !params.GameName.Valid || params.GameName.String != "hide on bush" {
		t.Fatalf("params.GameName = %+v, want valid hide on bush", params.GameName)
	}
	if !params.TagLine.Valid || params.TagLine.String != "KR1" {
		t.Fatalf("params.TagLine = %+v, want valid KR1", params.TagLine)
	}
	if !params.ProfileIconID.Valid || params.ProfileIconID.Int32 != 1234 {
		t.Fatalf("params.ProfileIconID = %+v, want valid 1234", params.ProfileIconID)
	}
	if !params.RevisionDate.Valid || !params.RevisionDate.Time.Equal(time.UnixMilli(1710000000000).UTC()) {
		t.Fatalf("params.RevisionDate = %+v, want valid unix millis", params.RevisionDate)
	}
}

func TestNewMatchHistoryParams(t *testing.T) {
	t.Parallel()

	params, err := newMatchHistoryParams(1, riot.MatchSummary{
		MatchID:         "JP1_1",
		PlayedAt:        1710000000000,
		GameVersion:     "14.1.123.4567",
		GameMode:        "CLASSIC",
		QueueID:         420,
		DurationSeconds: 1800,
	})
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	want := db.UpsertMatchHistoryParams{
		MatchID:         "JP1_1",
		RegionID:        1,
		QueueID:         420,
		GameMode:        "CLASSIC",
		GameVersion:     "14.1.123.4567",
		DurationSeconds: 1800,
		PlayedAt:        time.UnixMilli(1710000000000).UTC(),
	}
	if params != want {
		t.Fatalf("params = %+v, want %+v", params, want)
	}
}

func TestNewParticipantPlayerParams(t *testing.T) {
	t.Parallel()

	params := newParticipantPlayerParams(1, riot.MatchParticipantSummary{
		PUUID:    "participant-puuid",
		GameName: "ally",
		TagLine:  "JP1",
	})

	if params.Puuid != "participant-puuid" {
		t.Fatalf("params.Puuid = %q, want participant-puuid", params.Puuid)
	}
	if !params.GameName.Valid || params.GameName.String != "ally" {
		t.Fatalf("params.GameName = %+v, want valid ally", params.GameName)
	}
	if !params.TagLine.Valid || params.TagLine.String != "JP1" {
		t.Fatalf("params.TagLine = %+v, want valid JP1", params.TagLine)
	}
}

func TestNewMatchParticipantParams(t *testing.T) {
	t.Parallel()

	params := newMatchParticipantParams(10, 20, riot.MatchParticipantSummary{
		GameName:         "ally",
		TagLine:          "JP1",
		TeamID:           100,
		ChampionName:     "Ahri",
		Role:             "MIDDLE",
		Win:              true,
		Kills:            10,
		Deaths:           2,
		Assists:          5,
		SummonerSpell1ID: 4,
		SummonerSpell2ID: 14,
	})

	if params.MatchHistoryID != 10 || params.PlayerID != 20 {
		t.Fatalf("ids = %+v, want matchHistoryID=10 playerID=20", params)
	}
	if !params.GameNameSnapshot.Valid || params.GameNameSnapshot.String != "ally" {
		t.Fatalf("params.GameNameSnapshot = %+v, want valid ally", params.GameNameSnapshot)
	}
	if !params.Role.Valid || params.Role.String != "MIDDLE" {
		t.Fatalf("params.Role = %+v, want valid MIDDLE", params.Role)
	}
	if params.TeamPosition.Valid {
		t.Fatalf("params.TeamPosition = %+v, want invalid", params.TeamPosition)
	}
	if params.TeamID != 100 || params.Kills != 10 || params.SummonerSpell1ID != 4 {
		t.Fatalf("params = %+v, want mapped combat fields", params)
	}
}
