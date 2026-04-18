package httpadapter

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"feed-gg/backend/internal/usecase"
)

type PlayerSearchHandler struct {
	usecase PlayerSearchUsecase
}

type PlayerSearchUsecase interface {
	Execute(
		ctx context.Context,
		input usecase.PlayerSearchInput,
	) (*usecase.PlayerSearchResult, int, error)
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

func NewPlayerSearchHandler(usecase PlayerSearchUsecase) *PlayerSearchHandler {
	return &PlayerSearchHandler{
		usecase: usecase,
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

	player, statusCode, err := h.usecase.Execute(r.Context(), usecase.PlayerSearchInput{
		Region:   req.Region,
		GameName: req.GameName,
		TagLine:  req.TagLine,
	})
	if err != nil {
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
