package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/FelliniFeed/AuthApp.git/models"
	"github.com/FelliniFeed/AuthApp.git/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt = "gn2u39ighn39tn10328th"
	signingKey = "q!@51fewg4nG@!^&51gWQE"
	tockenTTL = 12 * time.Hour
)
	
type tokenClaims struct {
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

func (s *AuthService) GenerateToken(username, password string) (string, string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))

	if err != nil {
		return "","", err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tockenTTL).Unix(),
			IssuedAt: time.Now().Unix(),
		},
		user.ID,
	})

	tockenStr, err := accessToken.SignedString([]byte(signingKey))

	if err != nil {
		return "","", err
	}

	refreshToken := generateRefreshToken()

	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)

	if err != nil {
		return "","", err
	}

	err = s.repo.CreateRefreshToken(string(hash), user.ID)

	if err != nil {
		return "","", err
	}

	return tockenStr, refreshToken, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func generateRefreshToken() string {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	token := base64.StdEncoding.EncodeToString(bytes)
	return token
}