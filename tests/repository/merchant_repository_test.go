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

func createTempFileWithMerchants(t *testing.T, content []models.Merchant) *os.File {
	tempFile, err := ioutil.TempFile("", "test-merchants-*.json")
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

func TestJSONMerchantRepository_LoadAll(t *testing.T) {
	merchants := []models.Merchant{
		{ID: "1", Name: "Merchant One", Balance: 1000},
		{ID: "2", Name: "Merchant Two", Balance: 2000},
	}
	tempFile := createTempFileWithMerchants(t, merchants)
	defer os.Remove(tempFile.Name())

	repo := &repository.JSONMerchantRepository{FilePath: tempFile.Name()}
	result, err := repo.LoadAll()

	assert.NoError(t, err)
	assert.Equal(t, merchants, result)
}

func TestJSONMerchantRepository_FindMerchantByID(t *testing.T) {
	merchants := []models.Merchant{
		{ID: "1", Name: "Merchant One", Balance: 1000},
		{ID: "2", Name: "Merchant Two", Balance: 2000},
	}
	tempFile := createTempFileWithMerchants(t, merchants)
	defer os.Remove(tempFile.Name())

	repo := &repository.JSONMerchantRepository{FilePath: tempFile.Name()}

	t.Run("ID exists", func(t *testing.T) {
		result, err := repo.FindMerchantByID("1")
		assert.NoError(t, err)
		assert.Equal(t, &merchants[0], result)
	})

	t.Run("ID does not exist", func(t *testing.T) {
		result, err := repo.FindMerchantByID("3")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestJSONMerchantRepository_SaveMerchants(t *testing.T) {
	merchants := []models.Merchant{
		{ID: "1", Name: "Merchant One", Balance: 1000},
		{ID: "2", Name: "Merchant Two", Balance: 2000},
	}

	tempFile, err := ioutil.TempFile("", "test-merchants-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	repo := &repository.JSONMerchantRepository{FilePath: tempFile.Name()}
	err = repo.SaveMerchants(merchants)
	assert.NoError(t, err)

	// Verify the saved data
	savedData, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	var savedMerchants []models.Merchant
	err = json.Unmarshal(savedData, &savedMerchants)
	assert.NoError(t, err)

	assert.Equal(t, merchants, savedMerchants)
}

func TestJSONMerchantRepository_UpdateMerchantBalance(t *testing.T) {
	merchants := []models.Merchant{
		{ID: "1", Name: "Merchant One", Balance: 1000},
		{ID: "2", Name: "Merchant Two", Balance: 2000},
	}
	tempFile := createTempFileWithMerchants(t, merchants)
	defer os.Remove(tempFile.Name())

	repo := &repository.JSONMerchantRepository{FilePath: tempFile.Name()}

	t.Run("update balance successfully", func(t *testing.T) {
		err := repo.UpdateMerchantBalance(merchants, "1", 500)

		assert.NoError(t, err)

		// Verify that the balance has been updated
		updatedMerchants, err := repo.LoadAll()
		assert.NoError(t, err)

		assert.Equal(t, float64(1500), updatedMerchants[0].Balance)
	})

	t.Run("merchant not found", func(t *testing.T) {
		err := repo.UpdateMerchantBalance(merchants, "3", 500)
		if err != nil {
			assert.Error(t, err)
		}
		assert.Error(t, err)
	})
}
