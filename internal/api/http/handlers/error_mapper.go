package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ynsssss/pr-manager/internal/domain"
)

// TODO: rename and restructure
// TODO: write error mapper for every handler

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}


// TODO: create custom error
// TODO: fix status codes
func sendError(w http.ResponseWriter, err error) {
	resp := ErrorResponse{}
	switch {
	case domain.IsValidationError(err):
		// TODO: add code
		resp.Error.Message = err.Error()
		writeJSON(w, http.StatusBadRequest, resp)

	case errors.Is(err, domain.ErrTeamExists):
		resp.Error.Code = "TEAM_EXISTS"
		resp.Error.Message = "team_name already exists"
		writeJSON(w, http.StatusBadRequest, resp)

	case errors.Is(err, domain.ErrPRExists):
		resp.Error.Code = "PR_EXISTS"
		resp.Error.Message = "PR id already exists"
		writeJSON(w, http.StatusConflict, resp)

	case errors.Is(err, domain.ErrPRMerged):
		resp.Error.Code = "PR_MERGED"
		resp.Error.Message = "cannot reassign on merged PR"
		writeJSON(w, http.StatusConflict, resp)

	case errors.Is(err, domain.ErrNotAssigned):
		resp.Error.Code = "NOT_ASSIGNED"
		resp.Error.Message = "reviewer is not assigned to this PR"
		writeJSON(w, http.StatusConflict, resp)

	case errors.Is(err, domain.ErrNoCandidate):
		resp.Error.Code = "NO_CANDIDATE"
		resp.Error.Message = "no active replacement candidate in team"
		writeJSON(w, http.StatusConflict, resp)

	case errors.Is(err, domain.ErrNotFound):
		resp.Error.Code = "NOT_FOUND"
		resp.Error.Message = "resource not found"
		writeJSON(w, http.StatusNotFound, resp)

	default:
		resp.Error.Code = "INTERNAL"
		resp.Error.Message = err.Error()
		writeJSON(w, http.StatusInternalServerError, resp)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
