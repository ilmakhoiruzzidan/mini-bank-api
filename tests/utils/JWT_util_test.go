package utils

import (
	"github.com/stretchr/testify/assert"
	"mini-bank-api/models"
	"mini-bank-api/utils"
	"os"
	"testing"
)

func TestGetSecretKey(t *testing.T) {

	// Set environment variable for testing
	err := os.Setenv("JWT_SECRET_KEY", "testsecretkey")
	if err != nil {
		t.Fatalf("failed to set JWT_SECRET_KEY: %v", err)
	}

	// Call GetSecretKey function
	secretKey := utils.GetSecretKey()

	// Check if the secret key is correct
	assert.Equal(t, "testsecretkey", secretKey)

	// Clean up: remove the environment variable after the test
	defer os.Unsetenv("JWT_SECRET_KEY")
}

func TestGenerateJWT(t *testing.T) {
	// Prepare mock user
	customer := &models.Customer{
		ID:       "123",
		Username: "testuser",
		Password: "password123",
	}

	// Set environment variable for testing
	os.Setenv("JWT_SECRET_KEY", "testsecretkey")

	// Call GenerateJWT function
	token, err := utils.GenerateJWT(customer)

	// Check if token is generated successfully and no error occurs
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Try to parse the generated token to verify its validity
	claims, err := utils.ParseToken(token)
	assert.NoError(t, err)

	// Check if claims contain expected user data
	assert.Equal(t, "123", claims["id"])
	assert.Equal(t, "testuser", claims["username"])
	assert.Contains(t, claims, "exp")
}

func TestParseToken_InvalidToken(t *testing.T) {
	// Set environment variable for testing
	os.Setenv("JWT_SECRET_KEY", "testsecretkey")

	// Invalid token (fake token for testing)
	invalidToken := "invalid.token.string"

	// Call ParseToken with an invalid token
	claims, err := utils.ParseToken(invalidToken)

	// Check if an error is returned
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestParseToken_ValidToken(t *testing.T) {
	// Prepare mock user
	customer := &models.Customer{
		ID:       "123",
		Username: "testuser",
		Password: "password123",
	}

	// Set environment variable for testing
	os.Setenv("JWT_SECRET_KEY", "testsecretkey")

	// Generate a valid token
	token, err := utils.GenerateJWT(customer)
	assert.NoError(t, err)

	// Parse the generated token
	claims, err := utils.ParseToken(token)
	assert.NoError(t, err)

	// Check if claims contain expected user data
	assert.Equal(t, "123", claims["id"])
	assert.Equal(t, "testuser", claims["username"])
	assert.Contains(t, claims, "exp")
}
