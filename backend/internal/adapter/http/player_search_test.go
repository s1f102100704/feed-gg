package httpadapter

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"feed-gg/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type fakePlayerSearchUsecase struct {
	player     *usecase.PlayerSearchResult
	statusCode int
	err        error
}

func (f *fakePlayerSearchUsecase) Execute(
	ctx context.Context,
	input usecase.PlayerSearchInput,
) (*usecase.PlayerSearchResult, int, error) {
	return f.player, f.statusCode, f.err
}

func TestPlayerSearchHandler_Search(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		searchUsecase  PlayerSearchUsecase
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "returns player profile",
			path: "/api/players/JP1/hide%20on%20bush/KR1",
			searchUsecase: &fakePlayerSearchUsecase{
				player: &usecase.PlayerSearchResult{
					Region:   "JP1",
					PUUID:    "test-puuid",
					GameName: "hide on bush",
					TagLine:  "KR1",
				},
				statusCode: http.StatusOK,
			},
			wantStatusCode: http.StatusOK,
			wantBody:       "{\"region\":\"JP1\",\"puuid\":\"test-puuid\",\"gameName\":\"hide on bush\",\"tagLine\":\"KR1\",\"summonerLevel\":0,\"profileIconId\":0,\"profileIconUrl\":\"\",\"revisionDate\":0}\n",
		},
		{
			name: "returns bad request for missing required fields",
			path: "/api/players/JP1//KR1",
			searchUsecase: &fakePlayerSearchUsecase{
				statusCode: http.StatusBadRequest,
				err:        usecase.ErrInvalidPlayerSearchInput,
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"error\":\"region, gameName, and tagLine are required\"}\n",
		},
		{
			name: "returns bad request for unsupported region",
			path: "/api/players/XXX/hide%20on%20bush/KR1",
			searchUsecase: &fakePlayerSearchUsecase{
				statusCode: http.StatusBadRequest,
				err:        usecase.ErrUnsupportedRegion,
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"error\":\"unsupported region\"}\n",
		},
		{
			name: "returns server error when region master lookup fails",
			path: "/api/players/JP1/hide%20on%20bush/KR1",
			searchUsecase: &fakePlayerSearchUsecase{
				statusCode: http.StatusInternalServerError,
				err:        usecase.ErrRegionMasterUnavailable,
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       "{\"error\":\"failed to load region master\"}\n",
		},
		{
			name: "returns bad request when riot client rejects region",
			path: "/api/players/JP1/hide%20on%20bush/KR1",
			searchUsecase: &fakePlayerSearchUsecase{
				statusCode: http.StatusBadRequest,
				err:        usecase.ErrUnsupportedRegion,
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"error\":\"unsupported region\"}\n",
		},
		{
			name: "returns upstream status code",
			path: "/api/players/JP1/hide%20on%20bush/KR1",
			searchUsecase: &fakePlayerSearchUsecase{
				statusCode: http.StatusNotFound,
				err:        errors.New("not found"),
			},
			wantStatusCode: http.StatusNotFound,
			wantBody:       "{\"error\":\"not found\"}\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := NewPlayerSearchHandler(tt.searchUsecase)
			router := chi.NewRouter()
			router.Get("/api/players/{region}/{gameName}/{tagLine}", handler.Search)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Fatalf("status code = %d, want %d", rec.Code, tt.wantStatusCode)
			}
			if rec.Body.String() != tt.wantBody {
				t.Fatalf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}
