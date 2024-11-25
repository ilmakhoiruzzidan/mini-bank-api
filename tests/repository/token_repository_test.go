package repository

import (
	"github.com/stretchr/testify/assert"
	"mini-bank-api/repository"
	"os"
	"path/filepath"
	"testing"
)

func TestJSONTokenRepository_AddToRevocationList(t *testing.T) {
	// Menyiapkan file sementara
	tempFile := filepath.Join(t.TempDir(), "revoked_tokens.json")

	// Membuat repository baru dengan path file sementara
	repo := &repository.JSONTokenRepository{FilePath: tempFile}

	// Menambahkan token ke dalam daftar revokasi
	token := "sample-token"
	err := repo.AddToRevocationList(token)

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan file revoked_tokens.json ada setelah penambahan
	_, err = os.Stat(tempFile)
	assert.NoError(t, err)

	// Memastikan token ada dalam file
	revokedTokens, err := repo.LoadRevokedTokens()
	assert.NoError(t, err)
	assert.Len(t, revokedTokens, 1)
	assert.Equal(t, token, revokedTokens[0].Token)
}

func TestJSONTokenRepository_IsTokenRevoked(t *testing.T) {
	// Menyiapkan file sementara
	tempFile := filepath.Join(t.TempDir(), "revoked_tokens.json")

	// Membuat repository baru dengan path file sementara
	repo := &repository.JSONTokenRepository{FilePath: tempFile}

	// Menambahkan token ke dalam daftar revokasi
	token := "sample-token"
	err := repo.AddToRevocationList(token)
	assert.NoError(t, err)

	// Memeriksa apakah token telah di-revokasi
	isRevoked := repo.IsTokenRevoked(token)
	assert.True(t, isRevoked)

	// Memeriksa token yang belum di-revokasi
	nonRevokedToken := "non-revoked-token"
	isRevoked = repo.IsTokenRevoked(nonRevokedToken)
	assert.False(t, isRevoked)
}

func TestJSONTokenRepository_LoadRevokedTokens_EmptyFile(t *testing.T) {
	// Menyiapkan file sementara dengan file kosong
	tempFile := filepath.Join(t.TempDir(), "revoked_tokens.json")

	// Membuat repository baru dengan path file sementara
	repo := &repository.JSONTokenRepository{FilePath: tempFile}

	// Memuat revoked tokens dari file yang kosong
	revokedTokens, err := repo.LoadRevokedTokens()

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan tidak ada token yang dimuat
	assert.Len(t, revokedTokens, 0)
}

func TestJSONTokenRepository_LoadRevokedTokens(t *testing.T) {
	// Menyiapkan file sementara
	tempFile := filepath.Join(t.TempDir(), "revoked_tokens.json")

	// Membuat repository baru dengan path file sementara
	repo := &repository.JSONTokenRepository{FilePath: tempFile}

	// Menambahkan token ke dalam daftar revokasi
	token := "sample-token"
	err := repo.AddToRevocationList(token)
	assert.NoError(t, err)

	// Memuat revoked tokens dari file
	revokedTokens, err := repo.LoadRevokedTokens()

	// Memastikan tidak ada error
	assert.NoError(t, err)

	// Memastikan token yang ditambahkan ada dalam list revoked tokens
	assert.Len(t, revokedTokens, 1)
	assert.Equal(t, token, revokedTokens[0].Token)
}

func TestJSONTokenRepository_AddToRevocationList_Error(t *testing.T) {
	// Menggunakan path sementara yang valid
	repo := &repository.JSONTokenRepository{FilePath: filepath.Join(t.TempDir())}

	// Mencoba menambahkan token ke dalam revokasi
	err := repo.AddToRevocationList("sample-token")

	// Memastikan error terjadi
	assert.Error(t, err)
}
