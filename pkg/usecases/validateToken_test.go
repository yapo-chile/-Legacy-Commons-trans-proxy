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

	statusCode, err := interactor.CleanAndMatchToken("")
	assert.Error(t, err)
	assert.Equal(t, StatusUnauthorized, statusCode)
}

func TestValidToken(t *testing.T) {
	interactor := ValidateToken{
		SecretToken: defaultToken,
	}

	statusCode, err := interactor.CleanAndMatchToken("test")
	assert.NoError(t, err)
	assert.Equal(t, StatusOK, statusCode)
}
func TestNoCheckToken(t *testing.T) {
	interactor := ValidateToken{
		SecretToken: "",
	}

	statusCode, err := interactor.CleanAndMatchToken("test")
	assert.NoError(t, err)
	assert.Equal(t, StatusOK, statusCode)
}
