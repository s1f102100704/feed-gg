package httpadapter

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"feed-gg/backend/internal/infrastructure/riot"
)

type PlayerSearchHandler struct {
	riotClient  PlayerSearcher
	regionStore RegionChecker
}

type PlayerSearcher interface {
	SearchPlayerByRiotID(
		ctx context.Context,
		platformRegion string,
		gameName string,
		tagLine string,
	) (*riot.PlayerProfile, int, error)
}

type RegionChecker interface {
	RegionExists(ctx context.Context, name string) (bool, error)
}

type playerSearchRequest struct {
	Region   string `json:"region"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type errorResponse struct {
	Error string `json:"error"`
}

const playerSearchRequestBodyLimit = 1 << 10

func NewPlayerSearchHandler(riotClient PlayerSearcher, regionStore RegionChecker) *PlayerSearchHandler {
	return &PlayerSearchHandler{
		riotClient:  riotClient,
		regionStore: regionStore,
	}
}

func (h *PlayerSearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req playerSearchRequest
	r.Body = http.MaxBytesReader(w, r.Body, playerSearchRequestBodyLimit)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	req.Region = strings.TrimSpace(req.Region)
	req.GameName = strings.TrimSpace(req.GameName)
	req.TagLine = strings.TrimSpace(req.TagLine)

	if req.Region == "" || req.GameName == "" || req.TagLine == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "region, gameName, and tagLine are required"})
		return
	}

	exists, err := h.regionStore.RegionExists(r.Context(), req.Region)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to load region master"})
		return
	}
	if !exists {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "unsupported region"})
		return
	}

	player, statusCode, err := h.riotClient.SearchPlayerByRiotID(r.Context(), req.Region, req.GameName, req.TagLine)
	if err != nil {
		if errors.Is(err, riot.ErrInvalidRegion) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "unsupported region"})
			return
		}
		writeJSON(w, statusCode, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, player)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
