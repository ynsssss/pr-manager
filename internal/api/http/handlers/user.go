package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ynsssss/pr-manager/internal/domain"
	"github.com/ynsssss/pr-manager/internal/service"
)

// TODO: add custom error returner and marshaller

type UserHandler struct {
	userService *service.UserService
	prService   *service.PullRequestService
}

func NewUserHandler(
	userService *service.UserService,
	prService *service.PullRequestService,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		prService:   prService,
	}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req setActiveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, domain.ErrNotFound)
		return
	}

	user, err := h.userService.SetIsActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		sendError(w, err)
		return
	}

	writeJSON(w, 200, user)
	return
}

type setActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	userID := queryParams.Get("user_id")

	prs, err := h.prService.GetPullRequestsForUser(r.Context(), userID)
	if err != nil {
		sendError(w, err)
		return
	}

	response := GetReviewResponse{
		UserID: userID,
		Prs:    prs,
	}

	writeJSON(w, 200, response)
	return
}

// TODO: use pull request short
type GetReviewResponse struct {
	UserID string               `json:"user_id"`
	Prs    []domain.PullRequest `json:"pull_requests"`
}
