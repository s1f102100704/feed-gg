package usecase

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"feed-gg/backend/internal/infrastructure/riot"
)

type fakePlayerSearcher struct {
	player      *riot.PlayerProfile
	statusCode  int
	err         error
	gotRegion   string
	gotGameName string
	gotTagLine  string
}

type fakeRegionChecker struct {
	exists  bool
	err     error
	gotName string
}

type fakePlayerSearchCache struct{}

type fakePlayerSearchRepository struct{}

func (f *fakePlayerSearcher) SearchPlayerByRiotID(
	ctx context.Context,
	platformRegion string,
	gameName string,
	tagLine string,
) (*riot.PlayerProfile, int, error) {
	f.gotRegion = platformRegion
	f.gotGameName = gameName
	f.gotTagLine = tagLine
	return f.player, f.statusCode, f.err
}

func (f *fakeRegionChecker) RegionExists(ctx context.Context, name string) (bool, error) {
	f.gotName = name
	return f.exists, f.err
}

func (f *fakePlayerSearchCache) Get(key PlayerSearchKey) (*PlayerSearchResult, bool) {
	return nil, false
}

func (f *fakePlayerSearchCache) Set(key PlayerSearchKey, value *PlayerSearchResult) {}

func (f *fakePlayerSearchRepository) FindSavedPlayer(
	ctx context.Context,
	input PlayerSearchInput,
) (*PlayerSearchResult, error) {
	return nil, nil
}

func (f *fakePlayerSearchRepository) SaveFetchedPlayer(ctx context.Context, fetched *riot.PlayerProfile) error {
	return nil
}

func TestPlayerSearchExecuteReturnsUnsupportedRegionWhenMissingInMaster(t *testing.T) {
	t.Parallel()

	searcher := &fakePlayerSearcher{}
	usecase := NewPlayerSearch(
		&fakePlayerSearchCache{},
		&fakePlayerSearchRepository{},
		searcher,
		&fakeRegionChecker{exists: false},
	)

	result, statusCode, err := usecase.Execute(context.Background(), PlayerSearchInput{
		Region:   "XXX",
		GameName: "hide on bush",
		TagLine:  "KR1",
	})

	if result != nil {
		t.Fatalf("result = %+v, want nil", result)
	}
	if statusCode != http.StatusBadRequest {
		t.Fatalf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if !errors.Is(err, ErrUnsupportedRegion) {
		t.Fatalf("err = %v, want ErrUnsupportedRegion", err)
	}
	if searcher.gotRegion != "" {
		t.Fatalf("riot client called with region %q, want not called", searcher.gotRegion)
	}
}

func TestPlayerSearchExecuteReturnsMasterLookupFailure(t *testing.T) {
	t.Parallel()

	usecase := NewPlayerSearch(
		&fakePlayerSearchCache{},
		&fakePlayerSearchRepository{},
		&fakePlayerSearcher{},
		&fakeRegionChecker{err: errors.New("db down")},
	)

	result, statusCode, err := usecase.Execute(context.Background(), PlayerSearchInput{
		Region:   "JP1",
		GameName: "hide on bush",
		TagLine:  "KR1",
	})

	if result != nil {
		t.Fatalf("result = %+v, want nil", result)
	}
	if statusCode != http.StatusInternalServerError {
		t.Fatalf("statusCode = %d, want %d", statusCode, http.StatusInternalServerError)
	}
	if !errors.Is(err, ErrRegionMasterUnavailable) {
		t.Fatalf("err = %v, want ErrRegionMasterUnavailable", err)
	}
}

func TestPlayerSearchExecuteNormalizesInputBeforeLookup(t *testing.T) {
	t.Parallel()

	searcher := &fakePlayerSearcher{
		player:     &riot.PlayerProfile{Region: "JP1", GameName: "hide on bush", TagLine: "KR1"},
		statusCode: http.StatusOK,
	}
	regionChecker := &fakeRegionChecker{exists: true}
	usecase := NewPlayerSearch(
		&fakePlayerSearchCache{},
		&fakePlayerSearchRepository{},
		searcher,
		regionChecker,
	)

	result, statusCode, err := usecase.Execute(context.Background(), PlayerSearchInput{
		Region:   " jp1 ",
		GameName: " hide on bush ",
		TagLine:  " KR1 ",
	})

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("statusCode = %d, want %d", statusCode, http.StatusOK)
	}
	if result == nil {
		t.Fatal("result = nil, want non-nil")
	}
	if regionChecker.gotName != "JP1" {
		t.Fatalf("region checker got %q, want JP1", regionChecker.gotName)
	}
	if searcher.gotRegion != "JP1" || searcher.gotGameName != "hide on bush" || searcher.gotTagLine != "KR1" {
		t.Fatalf(
			"riot client called with %q %q %q, want JP1 hide on bush KR1",
			searcher.gotRegion,
			searcher.gotGameName,
			searcher.gotTagLine,
		)
	}
}

func TestPlayerSearchExecuteReturnsInvalidInputAfterNormalization(t *testing.T) {
	t.Parallel()

	searcher := &fakePlayerSearcher{}
	regionChecker := &fakeRegionChecker{exists: true}
	usecase := NewPlayerSearch(
		&fakePlayerSearchCache{},
		&fakePlayerSearchRepository{},
		searcher,
		regionChecker,
	)

	result, statusCode, err := usecase.Execute(context.Background(), PlayerSearchInput{
		Region:   "   ",
		GameName: " hide on bush ",
		TagLine:  " KR1 ",
	})

	if result != nil {
		t.Fatalf("result = %+v, want nil", result)
	}
	if statusCode != http.StatusBadRequest {
		t.Fatalf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if !errors.Is(err, ErrInvalidPlayerSearchInput) {
		t.Fatalf("err = %v, want ErrInvalidPlayerSearchInput", err)
	}
	if regionChecker.gotName != "" {
		t.Fatalf("region checker got %q, want not called", regionChecker.gotName)
	}
	if searcher.gotRegion != "" {
		t.Fatalf("riot client got region %q, want not called", searcher.gotRegion)
	}
}

func TestPlayerSearchExecuteReturnsMappedResult(t *testing.T) {
	t.Parallel()

	searcher := &fakePlayerSearcher{
		player: &riot.PlayerProfile{
			Region:         "JP1",
			PUUID:          "test-puuid",
			GameName:       "hide on bush",
			TagLine:        "KR1",
			SummonerLevel:  999,
			ProfileIconID:  1234,
			ProfileIconURL: "https://example.com/icon.png",
			RevisionDate:   1710000000000,
			SoloRank: &riot.RankedQueue{
				Tier:         "GOLD",
				Rank:         "II",
				LeaguePoints: 75,
				Wins:         10,
				Losses:       5,
			},
			Matches: []riot.MatchSummary{
				{
					MatchID:      "JP1_1",
					TeamID:       100,
					ChampionName: "Ahri",
					Participants: []riot.MatchParticipantSummary{
						{PUUID: "test-puuid", GameName: "hide on bush", TagLine: "KR1", TeamID: 100},
					},
				},
			},
		},
		statusCode: http.StatusOK,
	}
	usecase := NewPlayerSearch(
		&fakePlayerSearchCache{},
		&fakePlayerSearchRepository{},
		searcher,
		&fakeRegionChecker{exists: true},
	)

	result, statusCode, err := usecase.Execute(context.Background(), PlayerSearchInput{
		Region:   "JP1",
		GameName: "hide on bush",
		TagLine:  "KR1",
	})

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("statusCode = %d, want %d", statusCode, http.StatusOK)
	}
	if result == nil {
		t.Fatal("result = nil, want non-nil")
	}
	if result.GameName != "hide on bush" || result.TagLine != "KR1" {
		t.Fatalf("result name = %q#%q, want hide on bush#KR1", result.GameName, result.TagLine)
	}
	if result.SoloRank == nil || result.SoloRank.Tier != "GOLD" {
		t.Fatalf("result.SoloRank = %+v, want GOLD", result.SoloRank)
	}
	if len(result.Matches) != 1 || len(result.Matches[0].Participants) != 1 {
		t.Fatalf("result.Matches = %+v, want one mapped match", result.Matches)
	}
	if result.Matches[0].TeamID != 100 || result.Matches[0].Participants[0].TeamID != 100 {
		t.Fatalf("team ids = %+v, want 100", result.Matches[0])
	}
	if searcher.gotRegion != "JP1" || searcher.gotGameName != "hide on bush" || searcher.gotTagLine != "KR1" {
		t.Fatalf(
			"riot client called with %q %q %q, want JP1 hide on bush KR1",
			searcher.gotRegion,
			searcher.gotGameName,
			searcher.gotTagLine,
		)
	}
}

func TestPlayerSearchExecuteMapsRiotInvalidRegion(t *testing.T) {
	t.Parallel()

	usecase := NewPlayerSearch(
		&fakePlayerSearchCache{},
		&fakePlayerSearchRepository{},
		&fakePlayerSearcher{
			statusCode: http.StatusBadRequest,
			err:        riot.ErrInvalidRegion,
		},
		&fakeRegionChecker{exists: true},
	)

	result, statusCode, err := usecase.Execute(context.Background(), PlayerSearchInput{
		Region:   "JP1",
		GameName: "hide on bush",
		TagLine:  "KR1",
	})

	if result != nil {
		t.Fatalf("result = %+v, want nil", result)
	}
	if statusCode != http.StatusBadRequest {
		t.Fatalf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if !errors.Is(err, ErrUnsupportedRegion) {
		t.Fatalf("err = %v, want ErrUnsupportedRegion", err)
	}
}

func TestNewPlayerSearchKeyNormalizesInput(t *testing.T) {
	t.Parallel()

	key := NewPlayerSearchKey(PlayerSearchInput{
		Region:   " jp1 ",
		GameName: " hide on bush ",
		TagLine:  " KR1 ",
	})

	if key != PlayerSearchKey("JP1:hide on bush:KR1") {
		t.Fatalf("key = %q, want %q", key, PlayerSearchKey("JP1:hide on bush:KR1"))
	}
}
