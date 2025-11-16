package service

import (
	"context"
	"errors"
	"math/rand"
	"slices"
	"time"

	"github.com/ynsssss/pr-manager/internal/domain"
)

type PullRequestRepository interface {
	UpdateWithFn(
		ctx context.Context,
		id string,
		updateFn func(pr *domain.PullRequest) (*domain.PullRequest, error),
	) (*domain.PullRequest, error)

	Create(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
	GetByID(ctx context.Context, prID string) (*domain.PullRequest, error)

	GetPullRequestsForUser(ctx context.Context, userID string) ([]domain.PullRequest, error)
}

type PullRequestService struct {
	prRepo   PullRequestRepository
	userRepo UserRepository
	teamRepo TeamRepository
}

func NewPullRequestService(
	prRepo PullRequestRepository,
	userRepo UserRepository,
	teamRepo TeamRepository,
) *PullRequestService {
	return &PullRequestService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *PullRequestService) Create(
	ctx context.Context,
	id, title, authorID string,
) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
	}
	if pr != nil {
		return nil, domain.ErrPRExists
	}

	_, err = s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	reviewers, err := s.pickReviewers(ctx, authorID)
	if err != nil {
		return nil, err
	}

	newPrRequest := domain.PullRequest{
		ID:                id,
		Name:              title,
		AuthorID:          authorID,
		Status:            domain.StatusOpen,
		AssignedReviewers: reviewers,
	}

	newPr, err := s.prRepo.Create(ctx, &newPrRequest)
	if err != nil {
		return nil, err
	}

	return newPr, nil
}

func (s *PullRequestService) pickReviewers(
	ctx context.Context,
	authorID string,
) ([]string, error) {
	team, err := s.teamRepo.GetTeamWithUser(ctx, authorID)
	if err != nil {
		return nil, err
	}
	activeMembers := make([]domain.TeamMember, 0, len(team.Members))
	for _, member := range team.Members {
		if member.IsActive && member.UserID != authorID {
			activeMembers = append(activeMembers, member)
		}
	}
	var chosenMembersIDs []string
	for range min(len(activeMembers), 2) {
		id := activeMembers[rand.Intn(len(activeMembers))].UserID
		chosenMembersIDs = append(chosenMembersIDs, id)
	}

	return chosenMembersIDs, nil
}

// NOTE: it is not a transactional operation yet,
// consider using transaction manager
func (s *PullRequestService) ReassignReviewer(
	ctx context.Context,
	prID, oldReviewer string,
) (*domain.PullRequest, string, error) {
	reviewers, err := s.pickReviewers(ctx, oldReviewer)
	if err != nil {
		return nil, "", err
	}
	if len(reviewers) == 0 {
		return nil, "", domain.ErrNoCandidate
	}
	newAssignee := reviewers[0]
	pr, err := s.prRepo.UpdateWithFn(
		ctx,
		prID,
		func(pr *domain.PullRequest) (*domain.PullRequest, error) {
			if pr.Status == domain.StatusMerged {
				return pr, domain.ErrPRMerged
			}

			if !slices.Contains(pr.AssignedReviewers, oldReviewer) {
				return pr, domain.ErrNotAssigned
			}

			for i, id := range pr.AssignedReviewers {
				if id == oldReviewer {
					pr.AssignedReviewers[i] = newAssignee
				}
			}

			return pr, nil
		},
	)
	return pr, newAssignee, err
}

func (s *PullRequestService) Merge(ctx context.Context, prID string) (*domain.PullRequest, error) {
	return s.prRepo.UpdateWithFn(
		ctx,
		prID,
		func(pr *domain.PullRequest) (*domain.PullRequest, error) {
			if pr.Status == domain.StatusMerged {
				return pr, nil
			}

			now := time.Now()
			pr.Status = domain.StatusMerged
			pr.MergedAt = &now
			return pr, nil
		},
	)
}

func (s *PullRequestService) GetPullRequestsForUser(ctx context.Context, userId string) (
	[]domain.PullRequest,
	error,
) {
	return s.prRepo.GetPullRequestsForUser(ctx, userId)
}
