package handler

import (
	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	"gym-pro-2026-ptit/internal/domain/user"
	useruc "gym-pro-2026-ptit/internal/usecase/user"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	registerOTPUC   *useruc.RegisterRequestOTPUseCase
	verifyOTPUC     *useruc.VerifyOTPUseCase
	loginUC         *useruc.LoginUseCase
	getProfileUC    *useruc.GetProfileUseCase
	updateProfileUC *useruc.UpdateProfileUseCase
	refreshTokenUC  *useruc.RefreshTokenUseCase
	// TODO: Add OAuth use cases when implemented
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	registerOTPUC *useruc.RegisterRequestOTPUseCase,
	verifyOTPUC *useruc.VerifyOTPUseCase,
	loginUC *useruc.LoginUseCase,
	getProfileUC *useruc.GetProfileUseCase,
	updateProfileUC *useruc.UpdateProfileUseCase,
	refreshTokenUC *useruc.RefreshTokenUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerOTPUC:   registerOTPUC,
		verifyOTPUC:     verifyOTPUC,
		loginUC:         loginUC,
		getProfileUC:    getProfileUC,
		updateProfileUC: updateProfileUC,
		refreshTokenUC:  refreshTokenUC,
	}
}

// RegisterRequestOTP godoc
// @Summary Request OTP for registration
// @Description Request an OTP code to be sent via email for registration
// @Tags auth
// @Accept json
// @Produce json
// @Param request body useruc.RegisterRequestOTPInput true "Registration OTP request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /auth/register/request [post]
func (h *AuthHandler) RegisterRequestOTP(c *gin.Context) {
	var input useruc.RegisterRequestOTPInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	err := h.registerOTPUC.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "OTP sent to your email. Please verify within 5 minutes.",
	})
}

// VerifyOTP godoc
// @Summary Verify OTP and complete registration
// @Description Verify OTP code and create user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body useruc.VerifyOTPInput true "OTP verification"
// @Success 201 {object} response.Response{data=useruc.TokenPair}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /auth/register/verify [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var input useruc.VerifyOTPInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	result, err := h.verifyOTPUC.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, result)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.LoginInput true "Login credentials"
// @Success 200 {object} response.Response{data=useruc.TokenPair}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input user.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	result, err := h.loginUC.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response{data=RefreshTokenResponse}
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input useruc.RefreshTokenRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}
	result, err := h.refreshTokenUC.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=user.User}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	u, err := h.getProfileUC.Execute(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, u)
}

// UpdateMe godoc
// @Summary Update current user profile
// @Description Update authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.UpdateProfileInput true "Profile update data"
// @Success 200 {object} response.Response{data=user.User}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /users/me [put]
func (h *AuthHandler) UpdateMe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input user.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.BadRequest("invalid request body"))
		return
	}

	u, err := h.updateProfileUC.Execute(c.Request.Context(), userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, u)
}

// GoogleOAuth godoc
// @Summary Google OAuth login
// @Description Redirect to Google OAuth consent screen
// @Tags auth
// @Accept json
// @Produce json
// @Router /auth/oauth/google [get]
func (h *AuthHandler) GoogleOAuth(c *gin.Context) {
	// TODO: Implement Google OAuth redirect
	response.Error(c, errors.InternalServer("not implemented", nil))
}

// GoogleOAuthCallback godoc
// @Summary Google OAuth callback
// @Description Handle Google OAuth callback
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code"
// @Param state query string true "State parameter"
// @Success 200 {object} response.Response{data=useruc.TokenPair}
// @Failure 401 {object} response.Response
// @Router /auth/oauth/google/callback [get]
func (h *AuthHandler) GoogleOAuthCallback(c *gin.Context) {
	// TODO: Implement Google OAuth callback
	response.Error(c, errors.InternalServer("not implemented", nil))
}

// FacebookOAuth godoc
// @Summary Facebook OAuth login
// @Description Redirect to Facebook OAuth consent screen
// @Tags auth
// @Accept json
// @Produce json
// @Router /auth/oauth/facebook [get]
func (h *AuthHandler) FacebookOAuth(c *gin.Context) {
	// TODO: Implement Facebook OAuth redirect
	response.Error(c, errors.InternalServer("not implemented", nil))
}

// FacebookOAuthCallback godoc
// @Summary Facebook OAuth callback
// @Description Handle Facebook OAuth callback
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code"
// @Param state query string true "State parameter"
// @Success 200 {object} response.Response{data=useruc.TokenPair}
// @Failure 401 {object} response.Response
// @Router /auth/oauth/facebook/callback [get]
func (h *AuthHandler) FacebookOAuthCallback(c *gin.Context) {
	// TODO: Implement Facebook OAuth callback
	response.Error(c, errors.InternalServer("not implemented", nil))
}
