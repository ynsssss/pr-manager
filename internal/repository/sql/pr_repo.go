package sql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/ynsssss/pr-manager/internal/domain"
)

type PullRequestRepository struct {
	db *sql.DB
}

func NewPullRequestRepository(db *sql.DB) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (r *PullRequestRepository) Create(
	ctx context.Context,
	newPR *domain.PullRequest,
) (*domain.PullRequest, error) {
	query := `
INSERT INTO pull_requests
    (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
`

	newPR.CreatedAt = time.Now()

	var pr domain.PullRequest
	var assigned pq.StringArray
	err := r.db.QueryRowContext(
		ctx,
		query,
		newPR.ID,
		newPR.Name,
		newPR.AuthorID,
		newPR.Status,
		pq.Array(newPR.AssignedReviewers),
		newPR.CreatedAt,
	).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&assigned,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	pr.AssignedReviewers = []string(assigned)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}

func (r *PullRequestRepository) GetByID(
	ctx context.Context,
	prID string,
) (*domain.PullRequest, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT pull_request_id, pull_request_name, author_id, status,
		        assigned_reviewers, created_at, merged_at
		   FROM pull_requests WHERE pull_request_id = $1`,
		prID,
	)

	var pr domain.PullRequest
	var assigned pq.StringArray
	var mergedAt sql.NullTime

	err := row.Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&assigned,
		&pr.CreatedAt,
		&mergedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	pr.AssignedReviewers = []string(assigned)

	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	return &pr, nil
}

func (r *PullRequestRepository) UpdateWithFn(
	ctx context.Context,
	id string,
	updateFn func(pr *domain.PullRequest) (*domain.PullRequest, error),
) (*domain.PullRequest, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	pr, err := r.getByIDTx(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	pr, err = updateFn(pr)
	if err != nil {
		return nil, err
	}

	assigned := pq.StringArray(pr.AssignedReviewers)

	_, err = tx.ExecContext(
		ctx,
		`UPDATE pull_requests
		 SET pull_request_name = $1, status = $2, assigned_reviewers = $3, merged_at = $4
		 WHERE pull_request_id = $5`,
		pr.Name, pr.Status, assigned, pr.MergedAt, pr.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return pr, nil
}

func (r *PullRequestRepository) getByIDTx(
	ctx context.Context,
	tx *sql.Tx,
	prID string,
) (*domain.PullRequest, error) {
	row := tx.QueryRowContext(
		ctx,
		`SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
		 FROM pull_requests
		 WHERE pull_request_id = $1`,
		prID,
	)

	var pr domain.PullRequest
	var assigned pq.StringArray
	var mergedAt sql.NullTime

	if err := row.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &assigned, &pr.CreatedAt, &mergedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	pr.AssignedReviewers = []string(assigned)

	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	return &pr, nil
}

func (r *PullRequestRepository) GetPullRequestsForUser(
	ctx context.Context,
	userID string,
) ([]domain.PullRequest, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
         FROM pull_requests
         WHERE $1 = ANY(assigned_reviewers)`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []domain.PullRequest

	for rows.Next() {
		var pr domain.PullRequest
		var assignedReviewers []string
		var mergedAt sql.NullTime

		if err := rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.AuthorID,
			&pr.Status,
			pq.Array(&assignedReviewers),
			&pr.CreatedAt,
			&mergedAt,
		); err != nil {
			return nil, err
		}

		if mergedAt.Valid {
			pr.MergedAt = &mergedAt.Time
		}

		pr.AssignedReviewers = assignedReviewers

		prs = append(prs, pr)
	}

	return prs, nil
}
