package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const TYPE_URL = "url"

type Url struct {
	ID            uuid.UUID        `json:"id"`
	Code          string           `json:"code"`
	Name          string           `json:"name"`
	Url           string           `json:"url"`
	AccountNumber string           `json:"accountNumber"`
	ClickCount    int              `json:"clickCount"`
	State         string           `json:"state"`
	Metadata      *json.RawMessage `json:"metadata" swaggertype:"object"` // Nullable
	ExpiredAt     *time.Time       `json:"expiredAt"`                     // Nullable
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
	DeletedAt     *time.Time       `json:"deletedAt"`
}

type (
	CreateUrlRequest struct {
		Name          string           `json:"name" validate:"omitempty,noStartEndSpaces" example:"John Doe"`
		Url           string           `json:"url" validate:"required,url" example:"john.doe@example.com"`
		Code          string           `json:"code" validate:"omitempty,noStartEndSpaces,alphanumericsWithDelimiter" example:"john-doe"`
		AccountNumber string           `json:"accountNumber" validate:"required" example:"93090290290"`
		State         string           `json:"state" validate:"required"`
		Metadata      *json.RawMessage `json:"metadata" swaggertype:"object"`
		ExpiredAt     *time.Time       `json:"expiredAt"`
	}

	CreateUrlResponse struct {
		Type          string           `json:"type" example:"url"`
		Code          string           `json:"code"`
		Name          string           `json:"name" example:"example.com - shortened"`
		Url           string           `json:"url" example:"http://example.com"`
		AccountNumber string           `json:"accountNumber" example:"8989489484"`
		State         string           `json:"state" example:"active"`
		Metadata      *json.RawMessage `json:"metadata" swaggertype:"object"`
		ExpiredAt     *time.Time       `json:"expiredAt" example:"2026-01-24T15:57:37+07:00"`
		CreatedAt     time.Time        `json:"createdAt" example:"2026-01-24T15:57:37+07:00"`
		UpdatedAt     time.Time        `json:"updatedAt" example:"2026-01-24T15:57:37+07:00"`
	}
)

func (u *Url) Response() *CreateUrlResponse {
	return &CreateUrlResponse{
		Type:          TYPE_URL,
		Code:          u.Code,
		Name:          u.Name,
		Url:           u.Url,
		AccountNumber: u.AccountNumber,
		State:         u.State,
		Metadata:      u.Metadata,
		ExpiredAt:     u.ExpiredAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
