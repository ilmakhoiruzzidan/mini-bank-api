package repository

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"mini-bank-api/models"
	"mini-bank-api/repository"
	"os"
	"testing"
	"time"
)

func TestLogAction(t *testing.T) {
	// Buat file sementara untuk pengujian
	tempFile, err := os.CreateTemp("", "histories_test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Hapus file setelah pengujian selesai

	// Tuliskan data awal ke file sementara
	initialHistories := []models.History{
		{
			ID:         "test-id-1",
			CustomerID: "customer-1",
			Action:     "initial-action",
			Timestamp:  time.Now(),
		},
	}
	initialData, err := json.MarshalIndent(initialHistories, "", "  ")
	assert.NoError(t, err)
	err = os.WriteFile(tempFile.Name(), initialData, os.ModePerm)
	assert.NoError(t, err)

	// Inisialisasi repository dengan file sementara
	repo := &repository.JSONHistoryRepository{FilePath: tempFile.Name()}

	// Log aksi baru
	err = repo.LogAction("customer-2", "new-action")
	assert.NoError(t, err)

	// Baca file untuk memverifikasi hasil
	data, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	var updatedHistories []models.History
	err = json.Unmarshal(data, &updatedHistories)
	assert.NoError(t, err)

	// Verifikasi hasil
	assert.Equal(t, 2, len(updatedHistories))
	assert.Equal(t, "customer-1", updatedHistories[0].CustomerID)
	assert.Equal(t, "customer-2", updatedHistories[1].CustomerID)
	assert.Equal(t, "new-action", updatedHistories[1].Action)
	assert.NotEmpty(t, updatedHistories[1].ID)
}
