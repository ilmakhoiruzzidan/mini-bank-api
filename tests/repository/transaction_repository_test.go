package repository

import (
	"github.com/stretchr/testify/assert"
	"mini-bank-api/repository"
	"os"
	"path/filepath"
	"testing"
)

func TestJSONTransactionRepository_SaveTransaction(t *testing.T) {
	// Menyiapkan file sementara
	tempFile := filepath.Join(t.TempDir(), "transactions.json")

	// Membuat repository baru dengan path file sementara
	repo := &repository.JSONTransactionRepository{FilePath: tempFile}

	// Menyimpan transaksi
	senderID := "sender-123"
	merchantID := "merchant-123"
	amount := 100.0
	transactionID, err := repo.SaveTransaction(senderID, merchantID, amount)

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan transaksi ID yang baru saja dibuat ada
	assert.NotEmpty(t, transactionID)

	// Memastikan file transaksi ada
	_, err = os.Stat(tempFile)
	assert.NoError(t, err)
}

func TestJSONTransactionRepository_LoadTransactions(t *testing.T) {
	// Menyiapkan file sementara
	tempFile := filepath.Join(t.TempDir(), "transactions.json")

	// Membuat repository baru dengan path file sementara
	repo := &repository.JSONTransactionRepository{FilePath: tempFile}

	// Membuat transaksi pertama
	senderID := "sender-123"
	merchantID := "merchant-123"
	amount := 100.0
	_, err := repo.SaveTransaction(senderID, merchantID, amount)
	assert.NoError(t, err)

	// Memuat transaksi dari file
	transactions, err := repo.LoadTransactions()

	// Memastikan tidak ada error saat load
	assert.NoError(t, err)

	// Memastikan ada transaksi yang dimuat
	assert.Len(t, transactions, 1)

	// Memastikan transaksi yang dimuat sesuai
	assert.Equal(t, senderID, transactions[0].SenderID)
	assert.Equal(t, merchantID, transactions[0].MerchantID)
	assert.Equal(t, amount, transactions[0].Amount)
}

func TestJSONTransactionRepository_LoadTransactions_NoFile(t *testing.T) {
	// Menyiapkan file sementara dengan path yang tidak ada
	tempFile := filepath.Join(t.TempDir(), "nonexistent_transactions.json")

	// Membuat repository baru dengan path file yang tidak ada
	repo := &repository.JSONTransactionRepository{FilePath: tempFile}

	// Memuat transaksi dari file yang tidak ada
	transactions, err := repo.LoadTransactions()

	// Memastikan tidak ada error saat file tidak ada (harus mengembalikan array kosong)
	assert.NoError(t, err)

	// Memastikan tidak ada transaksi yang dimuat
	assert.Len(t, transactions, 0)
}

func TestJSONTransactionRepository_SaveTransaction_Error(t *testing.T) {
	// Menyebabkan error dengan memberikan path file yang tidak valid
	repo := &repository.JSONTransactionRepository{FilePath: "/invalid/path/transactions.json"}

	// Mencoba menyimpan transaksi
	_, err := repo.SaveTransaction("sender-123", "merchant-123", 100.0)

	// Memastikan error terjadi
	assert.Error(t, err)
}
