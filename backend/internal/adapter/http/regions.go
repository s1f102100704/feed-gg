package httpadapter

import (
	"context"
	"net/http"
)

type RegionsHandler struct {
	regionStore RegionLister
}

type RegionLister interface {
	ListRegionNames(ctx context.Context) ([]string, error)
}

type regionsResponse struct {
	Regions []string `json:"regions"`
}

func NewRegionsHandler(regionStore RegionLister) *RegionsHandler {
	return &RegionsHandler{regionStore: regionStore}
}

func (h *RegionsHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	regions, err := h.regionStore.ListRegionNames(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to load region master"})
		return
	}

	writeJSON(w, http.StatusOK, regionsResponse{Regions: regions})
}
