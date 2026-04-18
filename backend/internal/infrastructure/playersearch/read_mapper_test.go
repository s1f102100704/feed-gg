package playersearch

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	db "feed-gg/backend/internal/infrastructure/db/sqlc"
)

type fakeProfileIconURLResolver struct {
	url string
	err error
	got int
}

func (f *fakeProfileIconURLResolver) ProfileIconURL(ctx context.Context, profileIconID int) (string, error) {
	f.got = profileIconID
	return f.url, f.err
}

func TestMapSavedPlayerSearchResult(t *testing.T) {
	t.Parallel()

	resolver := &fakeProfileIconURLResolver{url: "https://example.com/icon.png"}
	result := mapSavedPlayerSearchResult(
		context.Background(),
		resolver,
		db.GetSavedPlayerByPuuidRow{
			PlayerID:      1,
			Puuid:         "test-puuid",
			GameName:      nullString("hide on bush"),
			TagLine:       nullString("KR1"),
			RegionName:    "JP1",
			ProfileIconID: nullInt32(1234),
			SummonerLevel: nullInt64(999),
			RevisionDate:  nullTime(time.UnixMilli(1710000000000).UTC()),
		},
		[]db.ListPlayerCurrentRanksByPlayerIDRow{
			{
				QueueType:    "RANKED_SOLO_5x5",
				Tier:         "GOLD",
				Division:     "II",
				LeaguePoints: nullInt32(75),
				Wins:         nullInt32(10),
				Losses:       nullInt32(5),
			},
			{
				QueueType:    "RANKED_FLEX_SR",
				Tier:         "SILVER",
				Division:     "I",
				LeaguePoints: nullInt32(50),
				Wins:         nullInt32(8),
				Losses:       nullInt32(7),
			},
		},
		[]db.ListRecentMatchHistoriesByPlayerIDRow{
			{
				MatchHistoryID:   10,
				MatchID:          "JP1_1",
				GameVersion:      "14.1.123.4567",
				GameMode:         "CLASSIC",
				QueueID:          420,
				PlayedAt:         time.UnixMilli(1710000000000).UTC(),
				DurationSeconds:  1800,
				TeamID:           100,
				ChampionName:     "Ahri",
				Role:             nullString("MIDDLE"),
				Win:              true,
				Kills:            10,
				Deaths:           2,
				Assists:          5,
				SummonerSpell1ID: 4,
				SummonerSpell2ID: 14,
			},
		},
		map[int64][]db.ListMatchParticipantsByMatchHistoryIDRow{
			10: {
				{
					Puuid:            "test-puuid",
					GameNameSnapshot: nullString("hide on bush"),
					TagLineSnapshot:  nullString("KR1"),
					TeamID:           100,
					ChampionName:     "Ahri",
					Role:             nullString("MIDDLE"),
					Win:              true,
					Kills:            10,
					Deaths:           2,
					Assists:          5,
					SummonerSpell1ID: 4,
					SummonerSpell2ID: 14,
				},
			},
		},
	)

	if result == nil {
		t.Fatal("result = nil, want non-nil")
	}
	if result.Region != "JP1" || result.PUUID != "test-puuid" {
		t.Fatalf("result = %+v, want mapped player identity", result)
	}
	if resolver.got != 1234 {
		t.Fatalf("resolver.got = %d, want 1234", resolver.got)
	}
	if result.ProfileIconURL != "https://example.com/icon.png" {
		t.Fatalf("result.ProfileIconURL = %q, want resolved URL", result.ProfileIconURL)
	}
	if result.SoloRank == nil || result.SoloRank.Tier != "GOLD" {
		t.Fatalf("result.SoloRank = %+v, want GOLD", result.SoloRank)
	}
	if result.FlexRank == nil || result.FlexRank.Tier != "SILVER" {
		t.Fatalf("result.FlexRank = %+v, want SILVER", result.FlexRank)
	}
	if len(result.Matches) != 1 || len(result.Matches[0].Participants) != 1 {
		t.Fatalf("result.Matches = %+v, want one match with one participant", result.Matches)
	}
	if result.Matches[0].Role != "MIDDLE" || result.Matches[0].Participants[0].Role != "MIDDLE" {
		t.Fatalf("roles = %+v, want MIDDLE", result.Matches[0])
	}
}

func TestMapSavedPlayerSearchResultFallsBackWhenIconResolverFails(t *testing.T) {
	t.Parallel()

	result := mapSavedPlayerSearchResult(
		context.Background(),
		&fakeProfileIconURLResolver{err: errors.New("ddragon down")},
		db.GetSavedPlayerByPuuidRow{
			ProfileIconID: nullInt32(1234),
		},
		nil,
		nil,
		nil,
	)

	if result.ProfileIconURL != "" {
		t.Fatalf("result.ProfileIconURL = %q, want empty fallback", result.ProfileIconURL)
	}
}

func TestSavedRole(t *testing.T) {
	t.Parallel()

	if got := savedRole(nullString("TOP"), nullString("SOLO")); got != "TOP" {
		t.Fatalf("savedRole(teamPosition, role) = %q, want TOP", got)
	}
	if got := savedRole(sqlNullStringInvalid(), nullString("MIDDLE")); got != "MIDDLE" {
		t.Fatalf("savedRole(_, role) = %q, want MIDDLE", got)
	}
	if got := savedRole(sqlNullStringInvalid(), sqlNullStringInvalid()); got != "UNKNOWN" {
		t.Fatalf("savedRole(_, _) = %q, want UNKNOWN", got)
	}
}

func nullInt32(value int32) sql.NullInt32 {
	return sql.NullInt32{Int32: value, Valid: true}
}

func nullInt64(value int64) sql.NullInt64 {
	return sql.NullInt64{Int64: value, Valid: true}
}

func nullTime(value time.Time) sql.NullTime {
	return sql.NullTime{Time: value, Valid: true}
}

func sqlNullStringInvalid() sql.NullString {
	return sql.NullString{}
}
