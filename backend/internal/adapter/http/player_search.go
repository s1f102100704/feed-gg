package httpadapter

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"feed-gg/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
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

type errorResponse struct {
	Error string `json:"error"`
}

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

	player, statusCode, err := h.usecase.Execute(r.Context(), usecase.PlayerSearchInput{
		Region:   chi.URLParam(r, "region"),
		GameName: chi.URLParam(r, "gameName"),
		TagLine:  chi.URLParam(r, "tagLine"),
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
