package domain

import "errors"


// TODO: make custom errors
var (
	ErrTeamExists  = errors.New("team already exists")
	ErrNotFound    = errors.New("resource not found")
	ErrPRExists    = errors.New("PR already exists")
	ErrPRMerged    = errors.New("PR is already merged")
	ErrNotAssigned = errors.New("reviewer not assigned to this PR")
	ErrNoCandidate = errors.New("no active replacement candidate in team")
)

// Team specific domain errors
var (
	ErrEmptyTeamName       = NewValidationError("team name is empty")
	ErrEmptyTeamMemberID   = NewValidationError("team member id is empty")
	ErrEmptyTeamMemberName = NewValidationError("team member name is empty")
)

// User specific domain errors
var (
	ErrEmptyUserID   = NewValidationError("user ID is empty")
	ErrEmptyUsername = NewValidationError("username is empty")
)

// Pull request specific domain errors
var (
	ErrEmptyID          = NewValidationError("pull request ID is empty")
	ErrEmptyName        = NewValidationError("pull request name is empty")
	ErrEmptyAuthorID    = NewValidationError("author ID is empty")
	ErrInvalidStatus    = NewValidationError("pull request status is invalid")
	ErrTooManyReviewers = NewValidationError("too many assigned reviewers")
)
