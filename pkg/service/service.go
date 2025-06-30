package service

import (
	"github.com/FelliniFeed/AuthApp.git/models"
	"github.com/FelliniFeed/AuthApp.git/pkg/repository"
	"github.com/google/uuid"
)

type Authorization interface {
	CreateUser(user models.User) (uuid.UUID, error)
	GenerateTocken(username, password string) (string, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
