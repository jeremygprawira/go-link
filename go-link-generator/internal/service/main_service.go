package service

import (
	"github.com/jeremygprawira/go-link-generator/internal/config"
	"github.com/jeremygprawira/go-link-generator/internal/repository"
)

// Dependencies groups all service dependencies.
type Dependencies struct {
	Repository repository.Repository
	Config     *config.Configuration
}

// Service groups all business-logic services.
type Service struct {
	Health HealthService
	User   UserService
}

func New(d Dependencies) *Service {
	return &Service{
		Health: NewHealthService(&d),
		User:   NewUserService(&d),
	}
}
