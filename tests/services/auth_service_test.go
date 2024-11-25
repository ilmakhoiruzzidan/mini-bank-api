package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mini-bank-api/models"
	"mini-bank-api/services"
	"mini-bank-api/utils"
	"testing"
)

// Mock Customer Repository
type MockCustomerRepo struct {
	mock.Mock
}

func (m *MockCustomerRepo) LoadAll() ([]models.Customer, error) {
	var customers []models.Customer
	args := m.Called(customers)
	return args.Get(0).([]models.Customer), args.Error(0)
}

func (m *MockCustomerRepo) SaveCustomers([]models.Customer) error {
	var customers []models.Customer
	args := m.Called(customers)
	return args.Error(0)
}

func (m *MockCustomerRepo) LoadByUsername(username string) (*models.Customer, error) {
	args := m.Called(username)
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepo) FindCustomerByID(id string) (*models.Customer, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Customer), args.Error(1)
}

// Mock History Repository
type MockHistoryRepo struct {
	mock.Mock
}

func (m *MockHistoryRepo) LogAction(customerID, action string) error {
	args := m.Called(customerID, action)
	return args.Error(0)
}

// Mock Token Repository
type MockTokenRepo struct {
	mock.Mock
}

func (m *MockTokenRepo) LoadRevokedTokens() ([]models.RevokedToken, error) {
	args := m.Called()
	return args.Get(0).([]models.RevokedToken), args.Error(1)
}

func (m *MockTokenRepo) IsTokenRevoked(token string) bool {
	args := m.Called(token)
	return args.Bool(0)
}

func (m *MockTokenRepo) AddToRevocationList(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func TestAuthService_Login(t *testing.T) {
	mockCustomerRepo := new(MockCustomerRepo)
	mockHistoryRepo := new(MockHistoryRepo)
	mockTokenRepo := new(MockTokenRepo)

	authService := services.NewAuthService(mockCustomerRepo, mockHistoryRepo, mockTokenRepo)

	// Setup mock for successful login
	customer := &models.Customer{
		ID:       "123",
		Username: "testuser",
		Password: "password123",
	}
	mockCustomerRepo.On("LoadByUsername", "testuser").Return(customer, nil)
	mockHistoryRepo.On("LogAction", "123", "login").Return(nil)

	// Mock the GenerateJWT call from utils package
	utils.GenerateJWT(customer)

	// Test successful login
	token, err := authService.Login("testuser", "password123")
	assert.NoError(t, err)
	assert.Equal(t, "access_token", token)

	// Test failed login due to incorrect password
	token, err = authService.Login("testuser", "wrongpassword")
	assert.Error(t, err)
	assert.Empty(t, token)

	// Verify the expectations
	mockCustomerRepo.AssertExpectations(t)
	mockHistoryRepo.AssertExpectations(t)
}

func TestAuthService_Logout(t *testing.T) {
	mockCustomerRepo := new(MockCustomerRepo)
	mockHistoryRepo := new(MockHistoryRepo)
	mockTokenRepo := new(MockTokenRepo)

	authService := services.NewAuthService(mockCustomerRepo, mockHistoryRepo, mockTokenRepo)

	// Simulate claims for valid token

	accessToken := "access_token"

	// Mock token parsing to return valid claims
	_, err := utils.ParseToken(accessToken)

	// Set up mock for IsTokenRevoked and AddToRevocationList
	mockTokenRepo.On("IsTokenRevoked", "valid_token").Return(false)
	mockTokenRepo.On("AddToRevocationList", "valid_token").Return(nil)
	mockHistoryRepo.On("LogAction", "123", "logout").Return(nil)

	// Test successful logout
	err = authService.Logout("valid_token")
	assert.NoError(t, err)

	// Test trying to logout with an already revoked token
	mockTokenRepo.On("IsTokenRevoked", "revoked_token").Return(true)
	err = authService.Logout("revoked_token")
	assert.EqualError(t, err, "token already revoked")

	// Verify the expectations
	mockTokenRepo.AssertExpectations(t)
	mockHistoryRepo.AssertExpectations(t)
}

func TestAuthService_GetCurrentUserInfo(t *testing.T) {
	mockCustomerRepo := new(MockCustomerRepo)
	mockHistoryRepo := new(MockHistoryRepo)
	mockTokenRepo := new(MockTokenRepo)

	authService := services.NewAuthService(mockCustomerRepo, mockHistoryRepo, mockTokenRepo)

	// Simulate claims for valid token

	accessToken := "access_token"

	// Mock token parsing to return valid claims
	_, err := utils.ParseToken(accessToken)

	// Set up mock for finding customer by ID
	customer := &models.Customer{
		ID:       "123",
		Username: "testuser",
		Password: "password123",
	}
	mockCustomerRepo.On("FindCustomerByID", "123").Return(customer, nil)
	mockHistoryRepo.On("LogAction", "123", "check profile").Return(nil)

	// Test valid GetCurrentUserInfo
	result, err := authService.GetCurrentUserInfo("valid_token")
	assert.NoError(t, err)
	assert.Equal(t, customer, result)

	accessToken = "access_token"
	// Test invalid token (invalid claims)
	_, err = utils.ParseToken(accessToken)

	_, err = authService.GetCurrentUserInfo("invalid_token")
	assert.Error(t, err)

	// Verify the expectations
	mockCustomerRepo.AssertExpectations(t)
	mockHistoryRepo.AssertExpectations(t)
}
