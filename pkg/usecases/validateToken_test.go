package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	defaultToken string = "test"
)

func TestInvalidToken(t *testing.T) {
	interactor := ValidateToken{
		SecretToken: defaultToken,
	}

	err := interactor.CleanAndMatchToken("")
	assert.Error(t, err)
}

func TestValidToken(t *testing.T) {
	interactor := ValidateToken{
		SecretToken: defaultToken,
	}

	err := interactor.CleanAndMatchToken("test")
	assert.NoError(t, err)
}
func TestNoCheckToken(t *testing.T) {
	interactor := ValidateToken{
		SecretToken: "",
	}

	err := interactor.CleanAndMatchToken("test")
	assert.NoError(t, err)
}
