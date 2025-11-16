package domain

type User struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// Validate checks the invariants of the User entity
func (u *User) Validate() error {
	if u.ID == "" {
		return ErrEmptyUserID
	}
	if u.Username == "" {
		return ErrEmptyUsername
	}
	return nil
}
