package user

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/infrastructure/email"
	"gym-pro-2026-ptit/internal/infrastructure/otp"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"
)

type RegisterRequestOTPInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequestOTPUseCase struct {
	userRepo     user.Repository
	otpService   otp.Service
	emailService email.Service
	validator    *validator.Validator
}

func NewRegisterRequestOTPUseCase(
	userRepo user.Repository,
	otpService otp.Service,
	emailService email.Service,
	validator *validator.Validator,
) *RegisterRequestOTPUseCase {
	return &RegisterRequestOTPUseCase{
		userRepo:     userRepo,
		otpService:   otpService,
		emailService: emailService,
		validator:    validator,
	}
}

func (uc *RegisterRequestOTPUseCase) Execute(ctx context.Context, input RegisterRequestOTPInput) error {
	if err := uc.validator.Validate(input); err != nil {
		return errors.Validation(err.Error())
	}

	exists, err := uc.userRepo.Exists(ctx, input.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.Conflict("email already registered")
	}

	otpCode, err := uc.otpService.Generate(ctx, input.Email)
	if err != nil {
		return err
	}

	if err := uc.emailService.SendOTP(input.Email, otpCode); err != nil {
		return errors.InternalServer("failed to send OTP email", err)
	}

	return nil
}
