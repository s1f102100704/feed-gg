package httpadapter

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net"
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
	voterKeySalt      string
	trustProxyHeaders bool
}

type playerLabelVoteRequest struct {
	LabelID int16 `json:"labelId"`
}

func NewPlayerLabelsHandler(
	usecase PlayerLabelsUsecase,
	voterKeySalt string,
	trustProxyHeaders bool,
) *PlayerLabelsHandler {
	return &PlayerLabelsHandler{
		usecase:           usecase,
		voterKeySalt:      voterKeySalt,
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

	puuid := chi.URLParam(r, "puuid")
	result, statusCode, err := h.usecase.Vote(r.Context(), usecase.PlayerLabelVoteInput{
		PUUID:    puuid,
		LabelID:  request.LabelID,
		VoterKey: playerLabelVoterKey(r, puuid, h.voterKeySalt, h.trustProxyHeaders),
	})
	if err != nil {
		writeJSON(w, statusCode, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func playerLabelVoterKey(
	r *http.Request,
	puuid string,
	salt string,
	trustProxyHeaders bool,
) string {
	source := puuid + "|" + clientIP(r, trustProxyHeaders) + "|" + r.UserAgent()
	mac := hmac.New(sha256.New, []byte(salt))
	_, _ = mac.Write([]byte(source))
	return hex.EncodeToString(mac.Sum(nil))
}

func clientIP(r *http.Request, trustProxyHeaders bool) string {
	if trustProxyHeaders {
		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			if ip := firstValidForwardedIP(forwardedFor); ip != "" {
				return ip
			}
		}

		if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); net.ParseIP(realIP) != nil {
			return realIP
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}

func firstValidForwardedIP(forwardedFor string) string {
	parts := strings.Split(forwardedFor, ",")
	for _, part := range parts {
		ip := strings.TrimSpace(part)
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	return ""
}
