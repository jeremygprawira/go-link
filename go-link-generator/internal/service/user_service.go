package service

import (
	"context"

	"github.com/jeremygprawira/go-link-generator/internal/models"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/errorc"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/formatter"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/generator"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/logger"
)

type UserService interface {
	Create(ctx context.Context, request *models.CreateUserRequest) (*models.User, error)
	GetByAccountNumber(ctx context.Context, accountNumber string) (*models.User, error)
}

type userService struct {
	d *Dependencies
}

func NewUserService(d *Dependencies) UserService {
	return &userService{d: d}
}

func (us *userService) Create(ctx context.Context, request *models.CreateUserRequest) (*models.User, error) {
	// Enrich wide event with business context
	logger.Add(ctx, "user", map[string]any{
		"operation": "create",
	})

	// Process phone number formatting if provided
	phoneNumber := request.PhoneNumber.Number
	phoneCountryCode := request.PhoneNumber.CountryCode

	// Format phone number to international format if provided
	if phoneNumber != "" {
		if phoneCountryCode == "" {
			phoneCountryCode = "ID"
		}
		formattedPhoneNumber, err := formatter.PhoneNumber(models.PhoneNumber{
			Number:      phoneNumber,
			CountryCode: phoneCountryCode,
		})
		if err != nil {
			logger.Add(ctx, "phone_format_error", err.Error())
			return nil, errorc.Error(errorc.ErrorInvalidInput, "Invalid phone number format")
		}
		phoneNumber = *formattedPhoneNumber
	}

	// Check if user already exists
	isUserExist, err := us.d.Repository.Postgre.User.CheckByEmailOrPhoneNumber(ctx, request.Email, phoneNumber)
	if err != nil {
		logger.AddError(ctx, &logger.ErrorContext{
			Type:      "DatabaseError",
			Code:      "USER_CHECK_FAILED",
			Message:   err.Error(),
			Retriable: true, // Database queries can be retried
		})
		return nil, errorc.Error(errorc.ErrorDatabase)
	}

	if isUserExist {
		logger.AddMap(ctx, map[string]any{
			"conflict_email": request.Email,
			"conflict_phone": phoneNumber,
		})
		return nil, errorc.Error(errorc.ErrorAlreadyExist, "User with the same email or phone number already exists")
	}

	// Generate unique account number
	accountNumber, err := generator.AccountNumber()
	if err != nil {
		logger.AddError(ctx, &logger.ErrorContext{
			Type:      "GenerationError",
			Code:      "ACCOUNT_NUMBER_GENERATION_FAILED",
			Message:   err.Error(),
			Retriable: false,
		})
		return nil, errorc.Error(errorc.ErrorInternalServer, "Failed to generate account number")
	}

	// Hash password
	hashedPassword, err := generator.Hash(request.Password)
	if err != nil {
		logger.AddError(ctx, &logger.ErrorContext{
			Type:      "HashError",
			Code:      "PASSWORD_HASH_FAILED",
			Message:   err.Error(),
			Retriable: false,
		})
		return nil, errorc.Error(errorc.ErrorInternalServer, "Failed to hash password")
	}

	var emailPtr *string
	if request.Email != "" {
		emailPtr = &request.Email
	}

	var phonePtr *string
	if phoneNumber != "" {
		phonePtr = &phoneNumber
	}

	// Create user model
	user := &models.User{
		AccountNumber:    accountNumber,
		Name:             request.Name,
		Email:            emailPtr,
		Password:         hashedPassword, // Use local variable, don't mutate input
		PhoneNumber:      phonePtr,
		PhoneCountryCode: phoneCountryCode,
		// CreatedAt and UpdatedAt will be set by database defaults
	}

	// Persist user to database
	err = us.d.Repository.Postgre.User.Create(ctx, user)
	if err != nil {
		logger.AddError(ctx, &logger.ErrorContext{
			Type:      "DatabaseError",
			Code:      "USER_CREATE_FAILED",
			Message:   err.Error(),
			Retriable: false, // Creation failures usually indicate constraint violations
		})
		return nil, errorc.Error(errorc.ErrorDatabase, "Failed to create user")
	}

	// Enrich wide event with success metrics
	logger.AddMap(ctx, map[string]any{
		"user_account_number": accountNumber,
		"user_has_email":      request.Email != "",
		"user_has_phone":      phoneNumber != "",
	})

	return user, nil
}

func (us *userService) GetByAccountNumber(ctx context.Context, accountNumber string) (*models.User, error) {
	user, err := us.d.Repository.Postgre.User.GetOneByAccountNumber(ctx, accountNumber)
	if err != nil {
		logger.AddError(ctx, &logger.ErrorContext{
			Type:      "DatabaseError",
			Code:      "USER_GET_FAILED",
			Message:   err.Error(),
			Retriable: false,
		})
		return nil, errorc.Error(errorc.ErrorDatabase, "Failed to get user")
	}

	if user == nil {
		logger.AddMap(ctx, map[string]any{
			"account_number": accountNumber,
		})
		return nil, errorc.Error(errorc.ErrorDataNotFound, "User not found")
	}

	return user, nil
}
