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

type fakePlayerSearchCache struct {
	cachedValue *PlayerSearchResult
	cachedHit   bool
	gotKey      PlayerSearchKey
	setKey      PlayerSearchKey
	setValue    *PlayerSearchResult
}

type fakePlayerSearchRepository struct {
	savedPlayer    *PlayerSearchResult
	findErr        error
	saveErr        error
	gotFindInput   PlayerSearchInput
	gotSaveFetched *riot.PlayerProfile
}

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

func (f *fakePlayerSearchCache) Get(key PlayerSearchKey) (*PlayerSearchResult, bool) {
	f.gotKey = key
	return f.cachedValue, f.cachedHit
}

func (f *fakePlayerSearchCache) Set(key PlayerSearchKey, value *PlayerSearchResult) {
	f.setKey = key
	f.setValue = value
}

func (f *fakePlayerSearchRepository) FindSavedPlayer(
	ctx context.Context,
	input PlayerSearchInput,
) (*PlayerSearchResult, error) {
	f.gotFindInput = input
	return f.savedPlayer, f.findErr
}

func (f *fakePlayerSearchRepository) SaveFetchedPlayer(ctx context.Context, fetched *riot.PlayerProfile) error {
	f.gotSaveFetched = fetched
	return f.saveErr
}

func TestPlayerSearchExecuteReturnsUnsupportedRegionWhenMissingInMaster(t *testing.T) {
	t.Parallel()

	searcher := &fakePlayerSearcher{}
	usecase := NewPlayerSearch(
		&fakePlayerSearchCache{},
		&fakePlayerSearchRepository{},
		searcher,
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

func TestPlayerSearchExecuteNormalizesInputBeforeLookup(t *testing.T) {
	t.Parallel()

	cache := &fakePlayerSearchCache{}
	searcher := &fakePlayerSearcher{
		player:     &riot.PlayerProfile{Region: "JP1", GameName: "hide on bush", TagLine: "KR1"},
		statusCode: http.StatusOK,
	}
	repository := &fakePlayerSearchRepository{}
	usecase := NewPlayerSearch(
		cache,
		repository,
		searcher,
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
	if repository.gotFindInput.Region != "JP1" ||
		repository.gotFindInput.GameName != "hide on bush" ||
		repository.gotFindInput.TagLine != "KR1" {
		t.Fatalf("repository got %+v, want normalized input", repository.gotFindInput)
	}
	if cache.gotKey != PlayerSearchKey("JP1:hide on bush:KR1") {
		t.Fatalf("cache got key %q, want normalized key", cache.gotKey)
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

func TestPlayerSearchExecuteReturnsCachedPlayerOnCacheHit(t *testing.T) {
	t.Parallel()

	cache := &fakePlayerSearchCache{
		cachedHit: true,
		cachedValue: &PlayerSearchResult{
			Region:   "JP1",
			PUUID:    "cached-puuid",
			GameName: "hide on bush",
			TagLine:  "KR1",
		},
	}
	repository := &fakePlayerSearchRepository{}
	searcher := &fakePlayerSearcher{}
	usecase := NewPlayerSearch(
		cache,
		repository,
		searcher,
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
	if result == nil || result.PUUID != "cached-puuid" {
		t.Fatalf("result = %+v, want cached player", result)
	}
	if repository.gotFindInput != (PlayerSearchInput{}) {
		t.Fatalf("repository got %+v, want not called", repository.gotFindInput)
	}
	if searcher.gotRegion != "" {
		t.Fatalf("riot client called with region %q, want not called", searcher.gotRegion)
	}
}

func TestPlayerSearchExecuteReturnsInvalidInputAfterNormalization(t *testing.T) {
	t.Parallel()

	searcher := &fakePlayerSearcher{}
	usecase := NewPlayerSearch(
		&fakePlayerSearchCache{},
		&fakePlayerSearchRepository{},
		searcher,
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
	if searcher.gotRegion != "" {
		t.Fatalf("riot client got region %q, want not called", searcher.gotRegion)
	}
}

func TestPlayerSearchExecuteReturnsSavedPlayerOnDBHit(t *testing.T) {
	t.Parallel()

	cache := &fakePlayerSearchCache{}
	repository := &fakePlayerSearchRepository{
		savedPlayer: &PlayerSearchResult{
			Region:   "JP1",
			PUUID:    "saved-puuid",
			GameName: "hide on bush",
			TagLine:  "KR1",
		},
	}
	searcher := &fakePlayerSearcher{}
	usecase := NewPlayerSearch(
		cache,
		repository,
		searcher,
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
	if result == nil || result.PUUID != "saved-puuid" {
		t.Fatalf("result = %+v, want saved player", result)
	}
	if searcher.gotRegion != "" {
		t.Fatalf("riot client called with region %q, want not called", searcher.gotRegion)
	}
	if repository.gotSaveFetched != nil {
		t.Fatalf("repository saved %+v, want not called", repository.gotSaveFetched)
	}
	if cache.setValue != repository.savedPlayer {
		t.Fatalf("cache set value %+v, want saved player %+v", cache.setValue, repository.savedPlayer)
	}
}

func TestPlayerSearchExecuteReturnsMappedResultAfterDBMiss(t *testing.T) {
	t.Parallel()

	cache := &fakePlayerSearchCache{}
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
	repository := &fakePlayerSearchRepository{}
	usecase := NewPlayerSearch(
		cache,
		repository,
		searcher,
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
	if repository.gotSaveFetched != searcher.player {
		t.Fatalf("repository saved %+v, want riot player %+v", repository.gotSaveFetched, searcher.player)
	}
	if cache.setValue == nil || cache.setValue.PUUID != "test-puuid" {
		t.Fatalf("cache set value %+v, want mapped riot result", cache.setValue)
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

func TestPlayerSearchExecuteReturnsRepositoryLookupFailure(t *testing.T) {
	t.Parallel()

	usecase := NewPlayerSearch(
		&fakePlayerSearchCache{},
		&fakePlayerSearchRepository{findErr: errors.New("db down")},
		&fakePlayerSearcher{},
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
	if !errors.Is(err, ErrSavedPlayerLookupFailed) {
		t.Fatalf("err = %v, want ErrSavedPlayerLookupFailed", err)
	}
}

func TestPlayerSearchExecuteReturnsSaveFailureAfterRiotFetch(t *testing.T) {
	t.Parallel()

	searcher := &fakePlayerSearcher{
		player:     &riot.PlayerProfile{Region: "JP1", GameName: "hide on bush", TagLine: "KR1"},
		statusCode: http.StatusOK,
	}
	repository := &fakePlayerSearchRepository{saveErr: errors.New("write failed")}
	usecase := NewPlayerSearch(
		&fakePlayerSearchCache{},
		repository,
		searcher,
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
	if !errors.Is(err, ErrFetchedPlayerSaveFailed) {
		t.Fatalf("err = %v, want ErrFetchedPlayerSaveFailed", err)
	}
	if repository.gotSaveFetched != searcher.player {
		t.Fatalf("repository saved %+v, want riot player %+v", repository.gotSaveFetched, searcher.player)
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
