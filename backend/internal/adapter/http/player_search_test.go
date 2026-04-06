package httpadapter

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"feed-gg/backend/internal/infrastructure/riot"
)

type fakePlayerSearcher struct {
	player     *riot.PlayerProfile
	statusCode int
	err        error
}

func (f *fakePlayerSearcher) SearchPlayerByRiotID(
	ctx context.Context,
	platformRegion string,
	gameName string,
	tagLine string,
) (*riot.PlayerProfile, int, error) {
	return f.player, f.statusCode, f.err
}

func TestPlayerSearchHandler_Search(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		body           string
		searcher       PlayerSearcher
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "returns player profile",
			body: `{"region":"JP1","gameName":"hide on bush","tagLine":"KR1"}`,
			searcher: &fakePlayerSearcher{
				player: &riot.PlayerProfile{
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
			name:           "returns bad request for missing required fields",
			body:           `{"region":"JP1","gameName":"","tagLine":"KR1"}`,
			searcher:       &fakePlayerSearcher{},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"error\":\"region, gameName, and tagLine are required\"}\n",
		},
		{
			name:           "returns bad request for invalid json",
			body:           `{`,
			searcher:       &fakePlayerSearcher{},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"error\":\"invalid request body\"}\n",
		},
		{
			name: "returns bad request for unsupported region",
			body: `{"region":"XXX","gameName":"hide on bush","tagLine":"KR1"}`,
			searcher: &fakePlayerSearcher{
				statusCode: http.StatusBadRequest,
				err:        riot.ErrInvalidRegion,
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"error\":\"unsupported region\"}\n",
		},
		{
			name: "returns upstream status code",
			body: `{"region":"JP1","gameName":"hide on bush","tagLine":"KR1"}`,
			searcher: &fakePlayerSearcher{
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

			handler := NewPlayerSearchHandler(tt.searcher)
			req := httptest.NewRequest(http.MethodPost, "/api/players/search", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()

			handler.Search(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Fatalf("status code = %d, want %d", rec.Code, tt.wantStatusCode)
			}
			if rec.Body.String() != tt.wantBody {
				t.Fatalf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestPlayerSearchHandler_Search_TooLargeBody(t *testing.T) {
	t.Parallel()

	handler := NewPlayerSearchHandler(&fakePlayerSearcher{})
	body := bytes.Repeat([]byte("a"), playerSearchRequestBodyLimit+1)
	req := httptest.NewRequest(http.MethodPost, "/api/players/search", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Search(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if rec.Body.String() != "{\"error\":\"invalid request body\"}\n" {
		t.Fatalf("body = %q, want invalid request body", rec.Body.String())
	}
}
