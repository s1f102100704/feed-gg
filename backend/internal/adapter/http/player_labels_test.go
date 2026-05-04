package httpadapter

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"feed-gg/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type fakePlayerLabelsUsecase struct {
	labelsResult  *usecase.PlayerLabelsResult
	voteResult    *usecase.PlayerLabelVoteResult
	statusCode    int
	err           error
	lastListInput usecase.PlayerLabelsInput
	lastVoteInput usecase.PlayerLabelVoteInput
}

func (f *fakePlayerLabelsUsecase) List(
	ctx context.Context,
	input usecase.PlayerLabelsInput,
) (*usecase.PlayerLabelsResult, int, error) {
	f.lastListInput = input
	return f.labelsResult, f.statusCode, f.err
}

func (f *fakePlayerLabelsUsecase) Vote(
	ctx context.Context,
	input usecase.PlayerLabelVoteInput,
) (*usecase.PlayerLabelVoteResult, int, error) {
	f.lastVoteInput = input
	return f.voteResult, f.statusCode, f.err
}

func TestPlayerLabelsHandler_List(t *testing.T) {
	t.Parallel()

	fake := &fakePlayerLabelsUsecase{
		statusCode: http.StatusOK,
		labelsResult: &usecase.PlayerLabelsResult{
			Labels:     []usecase.PlayerLabelSummary{{ID: 1, Name: "寄り遅め", VoteCount: 2}},
			TotalVotes: 2,
		},
	}
	handler := NewPlayerLabelsHandler(fake, "test-salt", false)
	router := chi.NewRouter()
	router.Get("/api/players/{puuid}/labels", handler.List)
	req := httptest.NewRequest(http.MethodGet, "/api/players/test-puuid/labels", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusOK)
	}
	if fake.lastListInput.PUUID != "test-puuid" {
		t.Fatalf("lastListInput = %+v, want puuid", fake.lastListInput)
	}

	wantBody := "{\"labels\":[{\"id\":1,\"name\":\"寄り遅め\",\"voteCount\":2}],\"totalVotes\":2}\n"
	if rec.Body.String() != wantBody {
		t.Fatalf("body = %q, want %q", rec.Body.String(), wantBody)
	}
}

func TestPlayerLabelsHandler_Vote(t *testing.T) {
	t.Parallel()

	fake := &fakePlayerLabelsUsecase{
		statusCode: http.StatusOK,
		voteResult: &usecase.PlayerLabelVoteResult{
			SelectedLabel: usecase.PlayerLabelSummary{ID: 1, Name: "寄り遅め", VoteCount: 2},
			Labels:        []usecase.PlayerLabelSummary{{ID: 1, Name: "寄り遅め", VoteCount: 2}},
			TotalVotes:    2,
		},
	}
	handler := NewPlayerLabelsHandler(fake, "test-salt", false)
	router := chi.NewRouter()
	router.Post("/api/players/{puuid}/labels", handler.Vote)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/players/test-puuid/labels",
		strings.NewReader(`{"labelId":1}`),
	)
	req.RemoteAddr = "203.0.113.10:12345"
	req.Header.Set("User-Agent", "test-agent")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusOK)
	}
	if fake.lastVoteInput.PUUID != "test-puuid" || fake.lastVoteInput.LabelID != 1 {
		t.Fatalf("lastVoteInput = %+v, want decoded input", fake.lastVoteInput)
	}
	if fake.lastVoteInput.VoterKey == "" || strings.Contains(fake.lastVoteInput.VoterKey, "203.0.113.10") {
		t.Fatalf("VoterKey = %q, want hashed voter key", fake.lastVoteInput.VoterKey)
	}
}

func TestPlayerLabelsHandler_VoteReturnsUsecaseError(t *testing.T) {
	t.Parallel()

	fake := &fakePlayerLabelsUsecase{
		statusCode: http.StatusBadRequest,
		err:        errors.New("invalid"),
	}
	handler := NewPlayerLabelsHandler(fake, "test-salt", false)
	req := httptest.NewRequest(http.MethodPost, "/api/players/test-puuid/labels", strings.NewReader(`{"labelId":0}`))
	rec := httptest.NewRecorder()

	handler.Vote(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if rec.Body.String() != "{\"error\":\"invalid\"}\n" {
		t.Fatalf("body = %q, want error response", rec.Body.String())
	}
}

func TestClientIPIgnoresForwardedHeadersByDefault(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/api/players/test-puuid/labels", nil)
	req.RemoteAddr = "203.0.113.10:12345"
	req.Header.Set("X-Forwarded-For", "198.51.100.1")

	if got := clientIP(req, false); got != "203.0.113.10" {
		t.Fatalf("clientIP = %q, want RemoteAddr host", got)
	}
}

func TestClientIPUsesValidForwardedHeaderWhenTrusted(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/api/players/test-puuid/labels", nil)
	req.RemoteAddr = "203.0.113.10:12345"
	req.Header.Set("X-Forwarded-For", "not-an-ip, 198.51.100.1")

	if got := clientIP(req, true); got != "198.51.100.1" {
		t.Fatalf("clientIP = %q, want first valid forwarded IP", got)
	}
}

func TestPlayerLabelVoterKeyIsScopedByPlayer(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/api/players/test-puuid/labels", nil)
	req.RemoteAddr = "203.0.113.10:12345"
	req.Header.Set("User-Agent", "test-agent")

	first := playerLabelVoterKey(req, "first-puuid", "test-salt", false)
	second := playerLabelVoterKey(req, "second-puuid", "test-salt", false)
	if first == second {
		t.Fatal("voter keys match across players, want player-scoped keys")
	}
}
