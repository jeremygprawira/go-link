package service

import (
	"context"
	"fmt"

	"github.com/jeremygprawira/go-link-generator/internal/models"
)

// HealthService defines the health-check interface.
type HealthService interface {
	Check(ctx context.Context) (*models.HealthResponse, error)
}

type healthService struct {
	d *Dependencies
}

func NewHealthService(d *Dependencies) HealthService {
	return &healthService{d: d}
}

func (hs *healthService) Check(ctx context.Context) (*models.HealthResponse, error) {
	var healthDetail []models.HealthDetailResponse

	var postgreHealth models.HealthDetailResponse
	postgreErr := hs.d.Repository.Postgre.Health.Check(ctx)
	kafkaProducerErr := hs.d.Producer.Check(ctx)
	if postgreErr != nil {
		postgreHealth = models.HealthDetailResponse{
			Type:        models.TYPE_HEALTH,
			Component:   "PostgreSQL",
			Status:      "ERROR",
			Description: fmt.Sprintf("PostgreSQL is not healthy, due to %v", postgreErr),
		}
		healthDetail = append(healthDetail, postgreHealth)
	} else {
		postgreHealth = models.HealthDetailResponse{
			Type:      models.TYPE_HEALTH,
			Component: "PostgreSQL",
			Status:    "OK",
		}
		healthDetail = append(healthDetail, postgreHealth)
	}

	var kafkaHealth models.HealthDetailResponse
	if kafkaProducerErr != nil {
		kafkaHealth = models.HealthDetailResponse{
			Type:        models.TYPE_HEALTH,
			Component:   "Kafka Producer",
			Status:      "ERROR",
			Description: fmt.Sprintf("Kafka Producer is not healthy, due to: %v", kafkaProducerErr),
		}
		healthDetail = append(healthDetail, kafkaHealth)
	} else {
		kafkaHealth = models.HealthDetailResponse{
			Type:      models.TYPE_HEALTH,
			Component: "Kafka Producer",
			Status:    "OK",
		}
		healthDetail = append(healthDetail, kafkaHealth)
	}

	return &models.HealthResponse{
		Description:  "Service is healthy",
		Dependencies: healthDetail,
	}, nil
}
