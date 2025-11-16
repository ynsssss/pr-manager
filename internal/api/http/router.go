package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ynsssss/pr-manager/internal/api/http/handlers"
	"github.com/ynsssss/pr-manager/internal/service"
)

func NewRouter(
	userService *service.UserService,
	teamService *service.TeamService,
	prService *service.PullRequestService,
) *mux.Router {
	router := mux.NewRouter()

	// Users
	userHandler := handlers.NewUserHandler(userService, prService)
	router.HandleFunc("/users/setIsActive", userHandler.SetIsActive).Methods(http.MethodPost)
	router.HandleFunc("/users/getReview", userHandler.GetReview).Methods(http.MethodGet)

	// Teams
	teamHandler := handlers.NewTeamHandler(teamService)
	router.HandleFunc("/team/add", teamHandler.Add).Methods(http.MethodPost)
	router.HandleFunc("/team/get", teamHandler.GetByName).Methods(http.MethodGet)

	// Pull Requests
	prHandler := handlers.NewPRHandler(prService)
	router.HandleFunc("/pullRequest/create", prHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/pullRequest/merge", prHandler.Merge).Methods(http.MethodPost)
	router.HandleFunc("/pullRequest/reassign", prHandler.Reassign).Methods(http.MethodPost)

	return router
}
