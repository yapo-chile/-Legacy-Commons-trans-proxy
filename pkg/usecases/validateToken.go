package usecases

import (
	"fmt"
	"strings"
)

const (
	InvalidToken string = "invalid token"
)

// ValidateTokenInteractor defines the methods that a Validate Token
// should have
type ValidateTokenInteractor interface {
	CleanAndMatchToken(token string) error
}

// ValidateToken defines the interactor
type ValidateToken struct {
	SecretToken string
}

// CleanAndMatchToken validates de input token according a default one
func (interactor *ValidateToken) CleanAndMatchToken(token string) error {
	if interactor.SecretToken != "" {
		// clean token
		token = strings.ReplaceAll(token, "Bearer ", "")
		token = strings.ReplaceAll(token, " ", "")

		// auth token validation
		if token == "" || token != interactor.SecretToken {
			return fmt.Errorf(InvalidToken)
		}
	}
	return nil
}
