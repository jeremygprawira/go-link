package models

import "time"

const TYPE_USER = "user"

type User struct {
	ID               int        `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountNumber    string     `json:"accountNumber" gorm:"unique;not null"`
	Name             string     `json:"name" gorm:"not null"`
	Email            *string    `json:"email" gorm:"unique"`
	PhoneNumber      *string    `json:"phoneNumber"`
	PhoneCountryCode string     `json:"phoneCountryCode"`
	Password         string     `json:"-" gorm:"not null"`
	IsActive         bool       `json:"isActive" gorm:"default:true"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}


type (
	CreateUserRequest struct {
		Name        string      `json:"name" validate:"required" example:"John Doe"`
		Email       string      `json:"email" validate:"required_without=PhoneNumber,omitempty,emailFormat" example:"john.doe@example.com"`
		PhoneNumber PhoneNumber `json:"phoneNumber" validate:"required_without=Email,omitempty"`
		Password    string      `json:"password" validate:"required" example:"password123"`
	}

	CreateUserResponse struct {
		Type          string      `json:"type" example:"user"`
		AccountNumber string      `json:"accountNumber" example:"1234567890"`
		Name          string      `json:"name" example:"John Doe"`
		Email         string      `json:"email" example:"john.doe@example.com"`
		PhoneNumber   PhoneNumber `json:"phoneNumber"`
		CreatedAt     time.Time   `json:"createdAt" example:"2026-01-24T15:57:37+07:00"`
		UpdatedAt     time.Time   `json:"updatedAt" example:"2026-01-24T15:57:37+07:00"`
	}
)

func (u *User) CreateUserResponse() *CreateUserResponse {
	email := ""
	if u.Email != nil {
		email = *u.Email
	}

	phoneNumber := ""
	if u.PhoneNumber != nil {
		phoneNumber = *u.PhoneNumber
	}

	return &CreateUserResponse{
		Type:          TYPE_USER,
		AccountNumber: u.AccountNumber,
		Name:          u.Name,
		Email:         email,
		PhoneNumber:   PhoneNumber{Number: phoneNumber, CountryCode: u.PhoneCountryCode},
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

type (
	GetUserTokenRequest struct {
		Email       string      `json:"email" validate:"required_without=PhoneNumber,omitempty,emailFormat" example:"john.doe@example.com"`
		PhoneNumber PhoneNumber `json:"phoneNumber" validate:"required_without=Email,omitempty"`
		Password    string      `json:"password" validate:"required" example:"password123"`
	}

	GetUserTokenResponse struct {
		Type          string      `json:"type" example:"user"`
		AccountNumber string      `json:"accountNumber" example:"1234567890"`
		Name          string      `json:"name" example:"John Doe"`
		Email         string      `json:"email" example:"john.doe@example.com"`
		PhoneNumber   PhoneNumber `json:"phoneNumber"`
		Tokens        []Token     `json:"tokens"`
	}
)

type (
	GetUserByAccountNumberResponse struct {
		Type          string      `json:"type" example:"user"`
		AccountNumber string      `json:"accountNumber" example:"1234567890"`
		Name          string      `json:"name" example:"John Doe"`
		Email         string      `json:"email" example:"john.doe@example.com"`
		PhoneNumber   PhoneNumber `json:"phoneNumber"`
		CreatedAt     time.Time   `json:"createdAt" example:"2026-01-24T15:57:37+07:00"`
		UpdatedAt     time.Time   `json:"updatedAt" example:"2026-01-24T15:57:37+07:00"`
	}
)

func (u *User) GetUserByAccountNumberResponse() *GetUserByAccountNumberResponse {
	email := ""
	if u.Email != nil {
		email = *u.Email
	}

	phoneNumber := ""
	if u.PhoneNumber != nil {
		phoneNumber = *u.PhoneNumber
	}

	return &GetUserByAccountNumberResponse{
		Type:          TYPE_USER,
		AccountNumber: u.AccountNumber,
		Name:          u.Name,
		Email:         email,
		PhoneNumber: PhoneNumber{
			Number:      phoneNumber,
			CountryCode: u.PhoneCountryCode,
		},
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
