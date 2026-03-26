package service

import (
	"context"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/jeremygprawira/go-link-generator/internal/models"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/broker"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/errorc"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/generator"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/logger"
)

type UrlService interface {
	Create(ctx context.Context, request *models.CreateUrlRequest) (*models.Url, error)
}

type urlService struct {
	d *Dependencies
}

func NewUrlService(d *Dependencies) UrlService {
	return &urlService{d: d}
}

func (us *urlService) Create(ctx context.Context, request *models.CreateUrlRequest) (*models.Url, error) {
	// Enrich wide event with business context
	logger.Add(ctx, "url", map[string]any{
		"operation": "create",
	})

	id, err := uuid.NewV7()
	if err != nil {
		return nil, errorc.Error(err)
	}
	logger.AddToKey(ctx, "url", "id", id)

	name := request.Name
	if name == "" {
		u, err := url.Parse(request.Url)
		if err != nil {
			return nil, errorc.Error(errorc.ErrorInvalidInput, err)
		}

		name = u.Hostname()

		logger.AddToKey(ctx, "url", "is_name_empty", true)
	}
	logger.AddToKey(ctx, "url", "name", name)

	var code string
	for attempt := 0; attempt < us.d.Config.Url.CodeGenerationRetries; attempt++ {
		if request.Code == "" {
			code, err = generator.SnowflakeID(us.d.Config.Url.SnowflakeMachineID).Base62().AddHMAC(us.d.Config.Url.SecureLength, us.d.Config.Url.Secret)
			if err != nil {
				return nil, err
			}

			logger.AddToKey(ctx, "url", "is_code_provided", false)
		} else {
			code = request.Code
			logger.AddToKey(ctx, "url", "is_code_provided", true)
		}

		exists, err := us.d.Repository.Postgre.Url.CheckByCode(ctx, code)
		if err != nil {
			return nil, err
		}

		if !exists {
			logger.AddToKey(ctx, "url", map[string]any{
				"code":                    code,
				"code_generation_attempt": attempt + 1,
				"code_collision":          false,
			})
			break
		}

		logger.AddToKey(ctx, "url", map[string]any{
			"code_collision":          true,
			"code_generation_attempt": attempt + 1,
		})

		if attempt == us.d.Config.Url.CodeGenerationRetries-1 {
			return nil, errorc.Error(errorc.ErrorAlreadyExist, err)
		} else {
			time.Sleep(time.Duration(us.d.Config.Url.CodeGenerationBackoff) * time.Millisecond * time.Duration(attempt+1)) // Exponential backoff
		}
	}

	url := &models.Url{
		ID:            id,
		Code:          code,
		Name:          name,
		Url:           request.Url,
		AccountNumber: request.AccountNumber,
		ClickCount:    0,
		State:         request.State,
		Metadata:      request.Metadata,
		ExpiredAt:     request.ExpiredAt,
	}

	if err := us.d.Repository.Postgre.Url.Create(ctx, url); err != nil {
		return nil, errorc.Error(err)
	}

	logger.AddProcess(ctx, "kafka", "url_create")
	if err := us.d.Producer.Produce(ctx, us.d.Config.Kafka.Topics.Link, url.ID.String(), url, broker.WithHeader("event_type", "url_created")); err != nil {
		return nil, errorc.Error(err)
	}

	return url, nil
}
