package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ynsssss/pr-manager/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1
	`, userID)

	var u domain.User
	if err := row.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}

	return u, nil
}

func (r *UserRepository) SetIsActive(
	ctx context.Context,
	userID string,
	isActive bool,
) (domain.User, error) {
	var u domain.User
	query := `
		UPDATE users
		SET is_active = $1
		WHERE user_id = $2
		RETURNING user_id, username, team_name, is_active
	`
	err := r.db.QueryRowContext(ctx, query, isActive, userID).
		Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserRepository) UpsertUsers(ctx context.Context, users []domain.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, u := range users {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) DO UPDATE SET
				username = EXCLUDED.username,
				team_name = EXCLUDED.team_name,
				is_active = EXCLUDED.is_active
		`, u.ID, u.Username, u.TeamName, u.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
