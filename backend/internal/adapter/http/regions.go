package httpadapter

import (
	"net/http"

	"feed-gg/backend/internal/domain/regions"
)

type RegionsHandler struct{}

type regionsResponse struct {
	Regions []string `json:"regions"`
}

func NewRegionsHandler() *RegionsHandler {
	return &RegionsHandler{}
}

func (h *RegionsHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	writeJSON(w, http.StatusOK, regionsResponse{Regions: regions.Names()})
}
