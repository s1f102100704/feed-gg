package riot

import "testing"

func TestMapRanksByQueue(t *testing.T) {
	t.Parallel()

	entries := []LeagueEntry{
		{
			QueueType:    "RANKED_SOLO_5x5",
			Tier:         "GOLD",
			Rank:         "II",
			LeaguePoints: 75,
			Wins:         40,
			Losses:       32,
		},
		{
			QueueType:    "RANKED_FLEX_SR",
			Tier:         "PLATINUM",
			Rank:         "IV",
			LeaguePoints: 12,
			Wins:         10,
			Losses:       8,
		},
	}

	soloRank, flexRank := mapRanksByQueue(entries)

	if soloRank == nil {
		t.Fatal("soloRank = nil, want non-nil")
	}
	if soloRank.Tier != "GOLD" || soloRank.Rank != "II" || soloRank.LeaguePoints != 75 {
		t.Fatalf("soloRank = %+v, want GOLD II 75LP", soloRank)
	}

	if flexRank == nil {
		t.Fatal("flexRank = nil, want non-nil")
	}
	if flexRank.Tier != "PLATINUM" || flexRank.Rank != "IV" || flexRank.LeaguePoints != 12 {
		t.Fatalf("flexRank = %+v, want PLATINUM IV 12LP", flexRank)
	}
}

func TestMapRanksByQueue_IgnoresUnknownQueues(t *testing.T) {
	t.Parallel()

	soloRank, flexRank := mapRanksByQueue([]LeagueEntry{
		{QueueType: "CHERRY", Tier: "DIAMOND", Rank: "I", LeaguePoints: 99},
	})

	if soloRank != nil {
		t.Fatalf("soloRank = %+v, want nil", soloRank)
	}
	if flexRank != nil {
		t.Fatalf("flexRank = %+v, want nil", flexRank)
	}
}
