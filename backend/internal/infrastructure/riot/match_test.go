package riot

import "testing"

func TestMapMatchSummary(t *testing.T) {
	t.Parallel()

	match := &MatchDTO{
		Metadata: MatchMetadataDTO{
			MatchID: "JP1_123456789",
		},
		Info: MatchInfoDTO{
			GameCreation:       1710000000000,
			GameStartTimestamp: 1710000100000,
			GameEndTimestamp:   1710002200000,
			GameDuration:       2100,
			GameVersion:        "15.7.123.4567",
			GameMode:           "CLASSIC",
			QueueID:            420,
			Participants: []MatchParticipantDTO{
				{
					RiotIDGameName: "target",
					RiotIDTagline:  "JP1",
					PUUID:          "target-puuid",
					ChampionName:   "Ahri",
					TeamPosition:   "MIDDLE",
					Kills:          10,
					Deaths:         2,
					Assists:        8,
					Win:            true,
					Summoner1ID:    4,
					Summoner2ID:    14,
				},
				{
					RiotIDGameName: "teammate",
					RiotIDTagline:  "KR1",
					PUUID:          "teammate-puuid",
					ChampionName:   "LeeSin",
					TeamPosition:   "JUNGLE",
					Kills:          4,
					Deaths:         5,
					Assists:        12,
					Win:            true,
					Summoner1ID:    11,
					Summoner2ID:    4,
				},
			},
		},
	}

	summary, err := mapMatchSummary(match, "target-puuid")
	if err != nil {
		t.Fatalf("mapMatchSummary returned error: %v", err)
	}

	if summary.MatchID != "JP1_123456789" {
		t.Fatalf("MatchID = %q, want JP1_123456789", summary.MatchID)
	}
	if summary.PlayedAt != 1710002200000 {
		t.Fatalf("PlayedAt = %d, want 1710002200000", summary.PlayedAt)
	}
	if summary.QueueID != 420 {
		t.Fatalf("QueueID = %d, want 420", summary.QueueID)
	}
	if summary.ChampionName != "Ahri" {
		t.Fatalf("ChampionName = %q, want Ahri", summary.ChampionName)
	}
	if summary.Role != "MIDDLE" {
		t.Fatalf("Role = %q, want MIDDLE", summary.Role)
	}
	if !summary.Win {
		t.Fatal("Win = false, want true")
	}
	if summary.Kills != 10 || summary.Deaths != 2 || summary.Assists != 8 {
		t.Fatalf("KDA = %d/%d/%d, want 10/2/8", summary.Kills, summary.Deaths, summary.Assists)
	}
	if summary.SummonerSpell1ID != 4 || summary.SummonerSpell2ID != 14 {
		t.Fatalf(
			"Summoner spells = %d/%d, want 4/14",
			summary.SummonerSpell1ID,
			summary.SummonerSpell2ID,
		)
	}
	if summary.DurationSeconds != 2100 {
		t.Fatalf("DurationSeconds = %d, want 2100", summary.DurationSeconds)
	}
	if len(summary.Participants) != 2 {
		t.Fatalf("len(Participants) = %d, want 2", len(summary.Participants))
	}
	if summary.Participants[0].GameName != "target" || summary.Participants[0].TagLine != "JP1" {
		t.Fatalf(
			"Participants[0] name = %q#%q, want target#JP1",
			summary.Participants[0].GameName,
			summary.Participants[0].TagLine,
		)
	}
	if summary.Participants[1].Role != "JUNGLE" {
		t.Fatalf("Participants[1].Role = %q, want JUNGLE", summary.Participants[1].Role)
	}
}

func TestMapMatchSummary_UsesFallbackRole(t *testing.T) {
	t.Parallel()

	match := &MatchDTO{
		Metadata: MatchMetadataDTO{MatchID: "JP1_1"},
		Info: MatchInfoDTO{
			Participants: []MatchParticipantDTO{
				{
					PUUID:              "target-puuid",
					ChampionName:       "Jinx",
					IndividualPosition: "BOTTOM",
				},
			},
		},
	}

	summary, err := mapMatchSummary(match, "target-puuid")
	if err != nil {
		t.Fatalf("mapMatchSummary returned error: %v", err)
	}

	if summary.Role != "BOTTOM" {
		t.Fatalf("Role = %q, want BOTTOM", summary.Role)
	}
}

func TestMapMatchSummary_ReturnsErrorWhenParticipantMissing(t *testing.T) {
	t.Parallel()

	match := &MatchDTO{
		Metadata: MatchMetadataDTO{MatchID: "JP1_1"},
		Info: MatchInfoDTO{
			Participants: []MatchParticipantDTO{
				{PUUID: "someone-else"},
			},
		},
	}

	if _, err := mapMatchSummary(match, "target-puuid"); err == nil {
		t.Fatal("mapMatchSummary error = nil, want non-nil")
	}
}

func TestMatchDurationSeconds_UsesMillisecondsWhenNeeded(t *testing.T) {
	t.Parallel()

	duration := matchDurationSeconds(MatchInfoDTO{
		GameDuration: 1650000,
	})

	if duration != 1650 {
		t.Fatalf("DurationSeconds = %d, want 1650", duration)
	}
}
