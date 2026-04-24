package httpadapter

import (
	"context"
	"net/http"

	"feed-gg/backend/internal/usecase"
)

type LabelsUsecase interface {
	Execute(ctx context.Context, input usecase.LabelsInput) ([]usecase.Label, error)
}

type LabelsHandler struct {
	usecase LabelsUsecase
}

type labelsResponse struct {
	Labels []usecase.Label `json:"labels"`
}

func NewLabelsHandler(usecase LabelsUsecase) *LabelsHandler {
	return &LabelsHandler{usecase: usecase}
}

func (h *LabelsHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	labels, err := h.usecase.Execute(r.Context(), usecase.LabelsInput{
		ForceRefresh: shouldRefresh(r),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, labelsResponse{Labels: labels})
}

func shouldRefresh(r *http.Request) bool {
	value := r.URL.Query().Get("refresh")
	switch value {
	case "1", "true", "TRUE", "yes", "YES":
		return true
	default:
		return false
	}
}
