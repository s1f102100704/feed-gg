package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"feed-gg/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type PlayerLabelsUsecase interface {
	List(ctx context.Context, input usecase.PlayerLabelsInput) (*usecase.PlayerLabelsResult, int, error)
	Vote(ctx context.Context, input usecase.PlayerLabelVoteInput) (*usecase.PlayerLabelVoteResult, int, error)
}

type PlayerLabelsHandler struct {
	usecase           PlayerLabelsUsecase
	voterKeyGenerator PlayerLabelVoterKeyGenerator
	trustProxyHeaders bool
}

type playerLabelVoteRequest struct {
	LabelID int16 `json:"labelId"`
}

func NewPlayerLabelsHandler(
	usecase PlayerLabelsUsecase,
	voterKeyGenerator PlayerLabelVoterKeyGenerator,
	trustProxyHeaders bool,
) *PlayerLabelsHandler {
	return &PlayerLabelsHandler{
		usecase:           usecase,
		voterKeyGenerator: voterKeyGenerator,
		trustProxyHeaders: trustProxyHeaders,
	}
}

func (h *PlayerLabelsHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result, statusCode, err := h.usecase.List(r.Context(), usecase.PlayerLabelsInput{
		PUUID: chi.URLParam(r, "puuid"),
	})
	if err != nil {
		writeJSON(w, statusCode, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *PlayerLabelsHandler) Vote(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var request playerLabelVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON request body"})
		return
	}

	puuid := strings.TrimSpace(chi.URLParam(r, "puuid"))
	result, statusCode, err := h.usecase.Vote(r.Context(), usecase.PlayerLabelVoteInput{
		PUUID:   puuid,
		LabelID: request.LabelID,
		VoterKey: h.voterKeyGenerator.Generate(
			puuid,
			newClientIdentity(r, h.trustProxyHeaders),
		),
	})
	if err != nil {
		writeJSON(w, statusCode, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}
