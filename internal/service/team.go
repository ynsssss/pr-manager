package service

import (
	"context"

	"github.com/ynsssss/pr-manager/internal/domain"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	TeamExists(ctx context.Context, name string) (bool, error)

	GetTeamByName(ctx context.Context, teamName string) (*domain.Team, error)
	GetTeamWithUser(ctx context.Context, userID string) (*domain.Team, error)
}

type TeamService struct {
	teamRepo TeamRepository
	userRepo UserRepository
}

func NewTeamService(teamRepo TeamRepository, userRepo UserRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// NOTE: operations here could be used in a single db transaction
// but given the RPS and the use case, it doesn't seem right but
// I could make them a transaction in repository implementation
// if that would be needed
func (s *TeamService) AddTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	if err := team.Validate(); err != nil {
		return nil, err
	}
	// Or GetTeam and then check if it's nil
	exists, err := s.teamRepo.TeamExists(ctx, team.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrTeamExists
	}

	newTeam, err := s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(team.Members))
	for _, member := range team.Members {
		users = append(users, domain.User{
			ID:       member.UserID,
			Username: member.Username,
			TeamName: team.Name,
			IsActive: member.IsActive,
		})
	}

	err = s.userRepo.UpsertUsers(ctx, users)
	if err != nil {
		return nil, err
	}

	return newTeam, nil
}

func (s *TeamService) GetByName(ctx context.Context, teamName string) (*domain.Team, error) {
	team, err := s.teamRepo.GetTeamByName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	return team, nil
}
