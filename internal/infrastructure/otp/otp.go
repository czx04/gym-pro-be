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
}

func NewOTPService(cache *cache.Cache) Service {
	return &otpService{cache: cache}
}

func (s *otpService) Generate(ctx context.Context, email string) (string, error) {
	otp, err := generateOTP(OTPLength)
	if err != nil {
		return "", errors.InternalServer("generate OTP", err)
	}

	// Store in Redis with expiration
	key := otpKey(email)
	if err := s.cache.Set(ctx, key, otp, OTPExpiration); err != nil {
		logger.Error("Failed to store OTP in cache", "email", email, "err", err)
		return "", errors.InternalServer("store OTP", err)
	}

	logger.Info("OTP generated", "email", email, "expires_in", OTPExpiration.String())

	return otp, nil
}

func (s *otpService) Verify(ctx context.Context, email, otp string) error {
	key := otpKey(email)

	// Get stored OTP
	storedOTP, err := s.cache.Get(ctx, key)
	if err != nil {
		logger.Warn("OTP verification failed - not found or expired", "email", email)
		return errors.BadRequest("Mã OTP đã hết hạn hoặc không hợp lệ")
	}

	// Compare OTP
	if storedOTP != otp {
		logger.Warn("OTP verification failed - mismatch", "email", email)
		return errors.BadRequest("Mã OTP không đúng")
	}

	logger.Info("OTP verified successfully", "email", email)

	if err := s.Delete(ctx, email); err != nil {
		logger.Warn("Failed to delete OTP after verification", "email", email, "err", err)
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
