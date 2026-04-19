package httpadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegionsHandler_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		wantStatusCode int
		wantBody       string
	}{
		{
			name:           "returns regions from master",
			method:         http.MethodGet,
			wantStatusCode: http.StatusOK,
			wantBody:       "{\"regions\":[\"BR1\",\"EUN1\",\"EUW1\",\"JP1\",\"KR\",\"LA1\",\"LA2\",\"ME1\",\"NA1\",\"OC1\",\"PH2\",\"RU\",\"SG2\",\"TH2\",\"TR1\",\"TW2\",\"VN2\"]}\n",
		},
		{
			name:           "returns no content for options",
			method:         http.MethodOptions,
			wantStatusCode: http.StatusNoContent,
			wantBody:       "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := NewRegionsHandler()
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
