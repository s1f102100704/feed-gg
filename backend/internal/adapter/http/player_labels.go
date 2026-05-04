package httpadapter

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"

	"feed-gg/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type PlayerLabelsUsecase interface {
	List(ctx context.Context, input usecase.PlayerLabelsInput) (*usecase.PlayerLabelsResult, int, error)
	Vote(ctx context.Context, input usecase.PlayerLabelVoteInput) (*usecase.PlayerLabelVoteResult, int, error)
}

type PlayerLabelsHandler struct {
	usecase PlayerLabelsUsecase
}

type playerLabelVoteRequest struct {
	LabelID int16 `json:"labelId"`
}

func NewPlayerLabelsHandler(usecase PlayerLabelsUsecase) *PlayerLabelsHandler {
	return &PlayerLabelsHandler{usecase: usecase}
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

	result, statusCode, err := h.usecase.Vote(r.Context(), usecase.PlayerLabelVoteInput{
		PUUID:    chi.URLParam(r, "puuid"),
		LabelID:  request.LabelID,
		VoterKey: playerLabelVoterKey(r),
	})
	if err != nil {
		writeJSON(w, statusCode, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func playerLabelVoterKey(r *http.Request) string {
	source := clientIP(r) + "|" + r.UserAgent()
	salt := os.Getenv("PLAYER_LABEL_VOTER_KEY_SALT")
	hash := sha256.Sum256([]byte(salt + "|" + source))
	return hex.EncodeToString(hash[:])
}

func clientIP(r *http.Request) string {
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}

	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
