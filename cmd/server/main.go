package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	httpserver "github.com/ynsssss/pr-manager/internal/api/http"
	sqlrepo "github.com/ynsssss/pr-manager/internal/repository/sql"
	"github.com/ynsssss/pr-manager/internal/service"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed to open DB connection: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	userRepo := sqlrepo.NewUserRepository(db)
	teamRepo := sqlrepo.NewTeamRepository(db)
	prRepo := sqlrepo.NewPullRequestRepository(db)

	userService := service.NewUserService(userRepo)
	teamService := service.NewTeamService(teamRepo, userRepo)
	prService := service.NewPullRequestService(prRepo, userRepo, teamRepo)

	router := httpserver.NewRouter(userService, teamService, prService)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting PR Manager service on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
