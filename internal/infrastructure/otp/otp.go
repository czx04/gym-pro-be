package otp

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"gym-pro-2026-ptit/internal/infrastructure/cache"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/errors"

	"go.uber.org/zap"
)

const (
	OTPLength     = 6
	OTPExpiration = 5 * time.Minute
	OTPKeyPrefix  = "otp:"
)

type Service interface {
	Generate(ctx context.Context, email string) (string, error)
	Verify(ctx context.Context, email, otp string) error
	Delete(ctx context.Context, email string) error
}

type otpService struct {
	cache *cache.Cache
	log   logger.Logger
}

func NewOTPService(cache *cache.Cache, log logger.Logger) Service {
	return &otpService{
		cache: cache,
		log:   log,
	}
}

func (s *otpService) Generate(ctx context.Context, email string) (string, error) {
	otp, err := generateOTP(OTPLength)
	if err != nil {
		return "", errors.InternalServer("generate OTP", err)
	}

	// Store in Redis with expiration
	key := otpKey(email)
	if err := s.cache.Set(ctx, key, otp, OTPExpiration); err != nil {
		s.log.Error("Failed to store OTP in cache",
			zap.String("email", email),
			zap.Error(err),
		)
		return "", errors.InternalServer("store OTP", err)
	}

	s.log.Info("OTP generated",
		zap.String("email", email),
		zap.String("expires_in", OTPExpiration.String()),
	)

	return otp, nil
}

func (s *otpService) Verify(ctx context.Context, email, otp string) error {
	key := otpKey(email)

	// Get stored OTP
	storedOTP, err := s.cache.Get(ctx, key)
	if err != nil {
		s.log.Warn("OTP verification failed - not found or expired",
			zap.String("email", email),
		)
		return errors.BadRequest("OTP expired or invalid")
	}

	// Compare OTP
	if storedOTP != otp {
		s.log.Warn("OTP verification failed - mismatch",
			zap.String("email", email),
		)
		return errors.BadRequest("OTP is incorrect")
	}

	s.log.Info("OTP verified successfully",
		zap.String("email", email),
	)

	if err := s.Delete(ctx, email); err != nil {
		s.log.Warn("Failed to delete OTP after verification",
			zap.String("email", email),
			zap.Error(err),
		)
	}

	return nil
}

func (s *otpService) Delete(ctx context.Context, email string) error {
	key := otpKey(email)
	if err := s.cache.Del(ctx, key).Err(); err != nil {
		return errors.InternalServer("delete OTP", err)
	}
	return nil
}

func generateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}

func otpKey(email string) string {
	return fmt.Sprintf("%s%s", OTPKeyPrefix, email)
}
