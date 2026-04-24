package httpadapter

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"feed-gg/backend/internal/usecase"
)

type fakeLabelsUsecase struct {
	labels    []usecase.Label
	err       error
	lastInput usecase.LabelsInput
}

func (f *fakeLabelsUsecase) Execute(ctx context.Context, input usecase.LabelsInput) ([]usecase.Label, error) {
	f.lastInput = input
	return f.labels, f.err
}

func TestLabelsHandler_List(t *testing.T) {
	t.Parallel()

	t.Run("returns labels", func(t *testing.T) {
		t.Parallel()

		fake := &fakeLabelsUsecase{
			labels: []usecase.Label{
				{ID: 1, Name: "寄り遅め"},
				{ID: 2, Name: "単独行動"},
			},
		}

		handler := NewLabelsHandler(fake)
		req := httptest.NewRequest(http.MethodGet, "/api/labels", nil)
		rec := httptest.NewRecorder()

		handler.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status code = %d, want %d", rec.Code, http.StatusOK)
		}

		wantBody := "{\"labels\":[{\"id\":1,\"name\":\"寄り遅め\"},{\"id\":2,\"name\":\"単独行動\"}]}\n"
		if rec.Body.String() != wantBody {
			t.Fatalf("body = %q, want %q", rec.Body.String(), wantBody)
		}

		if fake.lastInput.ForceRefresh {
			t.Fatal("ForceRefresh = true, want false")
		}
	})

	t.Run("supports refresh query", func(t *testing.T) {
		t.Parallel()

		fake := &fakeLabelsUsecase{}
		handler := NewLabelsHandler(fake)
		req := httptest.NewRequest(http.MethodGet, "/api/labels?refresh=1", nil)
		rec := httptest.NewRecorder()

		handler.List(rec, req)

		if !fake.lastInput.ForceRefresh {
			t.Fatal("ForceRefresh = false, want true")
		}
	})

	t.Run("returns internal server error", func(t *testing.T) {
		t.Parallel()

		fake := &fakeLabelsUsecase{err: errors.New("boom")}
		handler := NewLabelsHandler(fake)
		req := httptest.NewRequest(http.MethodGet, "/api/labels", nil)
		rec := httptest.NewRecorder()

		handler.List(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status code = %d, want %d", rec.Code, http.StatusInternalServerError)
		}

		wantBody := "{\"error\":\"boom\"}\n"
		if rec.Body.String() != wantBody {
			t.Fatalf("body = %q, want %q", rec.Body.String(), wantBody)
		}
	})
}
