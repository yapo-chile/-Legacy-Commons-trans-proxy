package usecases

import (
	"fmt"
	"strings"
)

const (
	StatusUnauthorized string = "no_auth"
	StatusOK           string = "ok"
	InvalidToken       string = "invalid token"
)

// ValidateTokenInteractor defines the methods that a Validate Token
// should have
type ValidateTokenInteractor interface {
	CleanAndMatchToken(token string) (string, error)
}

// ValidateToken defines the interactor
type ValidateToken struct {
	SecretToken string
}

// CleanAndMatchToken validates de input token according a default one
func (interactor *ValidateToken) CleanAndMatchToken(token string) (string, error) {
	if interactor.SecretToken != "" {
		// clean token
		token = strings.ReplaceAll(token, "Bearer ", "")
		token = strings.ReplaceAll(token, " ", "")

		// auth token validation
		if token == "" || token != interactor.SecretToken {
			return StatusUnauthorized, fmt.Errorf(InvalidToken)
		}
	}
	return StatusOK, nil
}
