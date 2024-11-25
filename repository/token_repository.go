package repository

import (
	"errors"
	"github.com/goccy/go-json"
	"log"
	"mini-bank-api/models"
	"os"
	"time"
)

type TokenRepositoryInterface interface {
	LoadRevokedTokens() ([]models.RevokedToken, error)
	IsTokenRevoked(token string) bool
	AddToRevocationList(token string) error
}

type JSONTokenRepository struct {
	FilePath string
}

func NewJSONTokenRepository() TokenRepositoryInterface {
	return &JSONTokenRepository{
		FilePath: "data/revoked_tokens.json",
	}
}

func (repo *JSONTokenRepository) AddToRevocationList(token string) error {
	revokedTokens, err := repo.LoadRevokedTokens()
	if err != nil {
		return err
	}

	newRevokedToken := models.RevokedToken{
		Token:      token,
		LogoutTime: time.Now(),
	}

	revokedTokens = append(revokedTokens, newRevokedToken)

	file, err := os.Create(repo.FilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	err = json.NewEncoder(file).Encode(revokedTokens)
	if err != nil {
		return err
	}

	return nil
}

func (repo *JSONTokenRepository) IsTokenRevoked(token string) bool {
	revokedTokens, err := repo.LoadRevokedTokens()
	if err != nil {
		log.Println("Error loading revoked tokens:", err)
		return true
	}

	for _, revoked := range revokedTokens {
		if revoked.Token == token {
			return true
		}
	}
	return false

}

func (repo *JSONTokenRepository) LoadRevokedTokens() ([]models.RevokedToken, error) {
	err := InitializeRevokedTokensFile(repo.FilePath)
	if err != nil {
		log.Fatalf("Failed to initialize revoked tokens file: %v", err)
	}

	file, err := os.Open(repo.FilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []models.RevokedToken{}, nil
		}
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() == 0 {
		return []models.RevokedToken{}, nil
	}

	var tokens []models.RevokedToken
	err = json.NewDecoder(file).Decode(&tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func InitializeRevokedTokensFile(filePath string) error {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.Write([]byte("[]")) // Tulis array kosong
		if err != nil {
			return err
		}
	}
	return nil
}
