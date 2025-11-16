package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ynsssss/pr-manager/internal/domain"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) TeamExists(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`, name).
		Scan(&exists)
	return exists, err
}

func (r *TeamRepository) CreateTeam(
	ctx context.Context,
	team *domain.Team,
) (*domain.Team, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO teams (team_name) VALUES ($1)`, team.Name)
	return team, err
}

func (r *TeamRepository) GetTeamByName(
	ctx context.Context,
	teamName string,
) (*domain.Team, error) {
	row := r.db.QueryRowContext(ctx, `SELECT team_name FROM teams WHERE team_name = $1`, teamName)

	var name string
	if err := row.Scan(&name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1
	`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]domain.TeamMember, 0)
	for rows.Next() {
		var m domain.TeamMember
		if err := rows.Scan(&m.UserID, &m.Username, &m.IsActive); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	return &domain.Team{
		Name:    teamName,
		Members: members,
	}, nil
}

// TODO: move to userRepo and rename to GetUserTeam
func (r *TeamRepository) GetTeamWithUser(
	ctx context.Context,
	userID string,
) (*domain.Team, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT team_name FROM users WHERE user_id = $1`,
		userID,
	)
	var teamName string
	if err := row.Scan(&teamName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return r.GetTeamByName(ctx, teamName)
}
