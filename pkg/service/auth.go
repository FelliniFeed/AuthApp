package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/FelliniFeed/AuthApp.git/models"
	"github.com/FelliniFeed/AuthApp.git/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	salt = "gn2u39ighn39tn10328th"
	signingKey = "q!@51fewg4nG@!^&51gWQE"
	tockenTTL = 12 * time.Hour
)
	
type tockenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"userID"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) (uuid.UUID, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateTocken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tockenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tockenTTL).Unix(),
			IssuedAt: time.Now().Unix(),
		},
		user.ID,
	})
	tockenStr, err := token.SignedString([]byte(signingKey))

	if err != nil {
		return "", err
	}

	return tockenStr, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}