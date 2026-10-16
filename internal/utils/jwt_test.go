package utils

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_Success(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	token, err := GenerateToken(userID, "test@test.com", "user", "secret", time.Hour)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken_Success(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	token, _ := GenerateToken(userID, "test@test.com", "user", "secret", time.Hour)

	claims, err := ValidateToken(token, "secret")

	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "test@test.com", claims.Email)
	assert.Equal(t, "user", claims.Role)
}

func TestValidateToken_InvalidSecret(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	token, _ := GenerateToken(userID, "test@test.com", "user", "secret", time.Hour)

	_, err := ValidateToken(token, "wrong-secret")

	assert.Error(t, err)
}

func TestValidateToken_Expired(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	token, _ := GenerateToken(userID, "test@test.com", "user", "secret", -time.Hour)

	_, err := ValidateToken(token, "secret")

	assert.Error(t, err)
}

func TestGenerateToken_UniqueForIdenticalInput(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())

	first, err := GenerateToken(userID, "test@test.com", "user", "secret", time.Hour)
	assert.NoError(t, err)
	second, err := GenerateToken(userID, "test@test.com", "user", "secret", time.Hour)
	assert.NoError(t, err)

	assert.NotEqual(t, first, second)
}

func TestValidateToken_HasUUIDTokenID(t *testing.T) {
	token, _ := GenerateToken(uuid.Must(uuid.NewV7()), "test@test.com", "user", "secret", time.Hour)

	claims, err := ValidateToken(token, "secret")

	assert.NoError(t, err)
	_, parseErr := uuid.Parse(claims.ID)
	assert.NoError(t, parseErr)
}
