package repository

import (
	"encoding/json"
	"io/ioutil"
	"mini-bank-api/models"
	"mini-bank-api/repository"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFileWithContent(t *testing.T, content []models.Customer) *os.File {
	tempFile, err := ioutil.TempFile("", "test-customers-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("failed to marshal content: %v", err)
	}

	if _, err := tempFile.Write(data); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	return tempFile
}

func TestJSONCustomerRepository_LoadAll(t *testing.T) {
	customers := []models.Customer{
		{ID: "1", Username: "user1", Password: "User One"},
		{ID: "2", Username: "user2", Password: "User Two"},
	}
	tempFile := createTempFileWithContent(t, customers)
	defer os.Remove(tempFile.Name())

	repo := &repository.JSONCustomerRepository{FilePath: tempFile.Name()}
	result, err := repo.LoadAll()

	assert.NoError(t, err)
	assert.Equal(t, customers, result)
}

func TestJSONCustomerRepository_LoadByUsername(t *testing.T) {
	customers := []models.Customer{
		{ID: "1", Username: "user1", Password: "User One"},
		{ID: "2", Username: "user2", Password: "User Two"},
	}
	tempFile := createTempFileWithContent(t, customers)
	defer os.Remove(tempFile.Name())

	repo := &repository.JSONCustomerRepository{FilePath: tempFile.Name()}

	t.Run("username exists", func(t *testing.T) {
		result, err := repo.LoadByUsername("user1")

		assert.NoError(t, err)
		assert.Equal(t, &customers[0], result)
	})

	t.Run("username does not exist", func(t *testing.T) {
		result, err := repo.LoadByUsername("unknown")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestJSONCustomerRepository_FindCustomerByID(t *testing.T) {
	customers := []models.Customer{
		{ID: "1", Username: "user1", Password: "User One"},
		{ID: "2", Username: "user2", Password: "User Two"},
	}
	tempFile := createTempFileWithContent(t, customers)
	defer os.Remove(tempFile.Name())

	repo := &repository.JSONCustomerRepository{FilePath: tempFile.Name()}

	t.Run("ID exists", func(t *testing.T) {
		result, err := repo.FindCustomerByID("1")

		assert.NoError(t, err)
		assert.Equal(t, &customers[0], result)
	})

	t.Run("ID does not exist", func(t *testing.T) {
		result, err := repo.FindCustomerByID("unknown")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestJSONCustomerRepository_SaveCustomers(t *testing.T) {
	customers := []models.Customer{
		{ID: "1", Username: "user1", Password: "User One"},
		{ID: "2", Username: "user2", Password: "User Two"},
	}

	tempFile, err := ioutil.TempFile("", "test-customers-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	repo := &repository.JSONCustomerRepository{FilePath: tempFile.Name()}
	err = repo.SaveCustomers(customers)
	assert.NoError(t, err)

	// Verify the saved data
	savedData, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	var savedCustomers []models.Customer
	err = json.Unmarshal(savedData, &savedCustomers)
	assert.NoError(t, err)

	assert.Equal(t, customers, savedCustomers)
}
