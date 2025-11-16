package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ynsssss/pr-manager/internal/domain"
	"github.com/ynsssss/pr-manager/internal/service"
)

type TeamHandler struct {
	service *service.TeamService
}

func NewTeamHandler(service *service.TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

// POST /team/add
func (h *TeamHandler) Add(w http.ResponseWriter, r *http.Request) {
	var teamRequest domain.Team

	if err := json.NewDecoder(r.Body).Decode(&teamRequest); err != nil {
		sendError(w, err)
		return
	}

	team, err := h.service.AddTeam(r.Context(), &teamRequest)
	if err != nil {
		sendError(w, err)
		return
	}

	writeJSON(w, 201, team)
	return
}

// GET /team/get
func (h *TeamHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	teamName := queryParams.Get("team_name")
	team, err := h.service.GetByName(r.Context(), teamName)
	if err != nil {
		sendError(w, err)
		return
	}

	writeJSON(w, 200, team)
	return
}
