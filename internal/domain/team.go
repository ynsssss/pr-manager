package domain

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (t *TeamMember) Validate() error {
	if t.UserID == "" {
		return ErrEmptyTeamMemberID
	}
	if t.Username == "" {
		return ErrEmptyTeamMemberName
	}

	return nil
}

type Team struct {
	Name    string       `json:"team_name"`
	Members []TeamMember `json:"members"`
}

func (t *Team) Validate() error {
	if t.Name == "" {
		return ErrEmptyTeamName
	}

	for _, member := range t.Members {
		if err := member.Validate(); err != nil {
			return err
		}
	}

	return nil
}
