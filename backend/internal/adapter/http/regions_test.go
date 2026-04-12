package httpadapter

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeRegionLister struct {
	regions []string
	err     error
}

func (f *fakeRegionLister) ListRegionNames(ctx context.Context) ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.regions, nil
}

func TestRegionsHandler_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		regionStore    RegionLister
		method         string
		wantStatusCode int
		wantBody       string
	}{
		{
			name:           "returns regions from master",
			method:         http.MethodGet,
			regionStore:    &fakeRegionLister{regions: []string{"KR", "JP1", "NA1"}},
			wantStatusCode: http.StatusOK,
			wantBody:       "{\"regions\":[\"KR\",\"JP1\",\"NA1\"]}\n",
		},
		{
			name:           "returns no content for options",
			method:         http.MethodOptions,
			regionStore:    &fakeRegionLister{},
			wantStatusCode: http.StatusNoContent,
			wantBody:       "",
		},
		{
			name:           "returns server error on master failure",
			method:         http.MethodGet,
			regionStore:    &fakeRegionLister{err: errors.New("db down")},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       "{\"error\":\"failed to load region master\"}\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := NewRegionsHandler(tt.regionStore)
			req := httptest.NewRequest(tt.method, "/api/regions", nil)
			rec := httptest.NewRecorder()

			handler.List(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Fatalf("status code = %d, want %d", rec.Code, tt.wantStatusCode)
			}
			if rec.Body.String() != tt.wantBody {
				t.Fatalf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}
