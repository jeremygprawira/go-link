package service

import (
	"github.com/jeremygprawira/go-link-generator/internal/config"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/broker"
	"github.com/jeremygprawira/go-link-generator/internal/repository"
)

// Dependencies groups all service dependencies.
type Dependencies struct {
	Repository repository.Repository
	Config     *config.Configuration
	Producer   *broker.Producer

	Health HealthService
	Url    UrlService
}

// Service groups all business-logic services.
type Service struct {
	Health HealthService
	Url    UrlService
}

func New(d Dependencies) *Service {
	health := NewHealthService(&d)
	url := NewUrlService(&d)

	d.Health = health
	d.Url = url

	return &Service{
		Health: health,
		Url:    url,
	}
}
