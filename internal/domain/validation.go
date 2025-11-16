package domain

import "errors"

type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

func NewValidationError(msg string) error {
	return ValidationError(msg)
}

func IsValidationError(err error) bool {
	var ve ValidationError
	return errors.As(err, &ve)
}
