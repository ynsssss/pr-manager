package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ynsssss/pr-manager/internal/domain"
	"github.com/ynsssss/pr-manager/internal/service"
)

type PRHandler struct {
	svc *service.PullRequestService
}

func NewPRHandler(svc *service.PullRequestService) *PRHandler {
	return &PRHandler{svc: svc}
}

func (h *PRHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err)
		return
	}

	pr, err := h.svc.Create(r.Context(), req.ID, req.Name, req.AuthorID)
	if err != nil {
		sendError(w, err)
		return
	}
	writeJSON(w, 201, pr)
}

type CreateRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

func (h *PRHandler) Merge(w http.ResponseWriter, r *http.Request) {
	var req mergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err)
		return
	}

	pr, err := h.svc.Merge(r.Context(), req.PrId)
	if err != nil {
		sendError(w, err)
		return
	}

	writeJSON(w, 200, pr)
}

type mergeRequest struct {
	PrId string `json:"pull_request_id"`
}

func (h *PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	var req reassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err)
		return
	}

	pr, replacedById, err := h.svc.ReassignReviewer(r.Context(), req.PrId, req.OldReviewerId)
	if err != nil {
		sendError(w, err)
		return
	}
	response := reassignResponse{
		Pr:           *pr,
		ReplacedById: replacedById,
	}

	writeJSON(w, 200, response)
}

type reassignRequest struct {
	PrId          string `json:"pull_request_id"`
	OldReviewerId string `json:"old_reviewer_id"`
}

type reassignResponse struct {
	Pr           domain.PullRequest `json:"pr"`
	ReplacedById string             `json:"replaced_by"`
}
