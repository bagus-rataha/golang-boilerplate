package handlers

import (
	"fiber-api-boilerplate/internal/dto"
	"fiber-api-boilerplate/internal/services"
	"fiber-api-boilerplate/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterInput true "Register request"
// @Success 201 {object} utils.Response{data=dto.TokenResponse}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input dto.RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errors := utils.ValidateStruct(&input); len(errors) > 0 {
		return utils.ValidationErrorResponse(c, errors)
	}

	result, err := h.authService.Register(input)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "User registered successfully", result)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Param request body dto.LoginInput true "Login request"
// @Success 200 {object} utils.Response{data=dto.TokenResponse}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input dto.LoginInput

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errors := utils.ValidateStruct(&input); len(errors) > 0 {
		return utils.ValidationErrorResponse(c, errors)
	}

	result, err := h.authService.Login(input)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", result)
}

// RefreshToken godoc
// @Summary Refresh token
// @Tags auth
// @Param request body dto.RefreshTokenInput true "Refresh token"
// @Success 200 {object} utils.Response{data=dto.TokenResponse}
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var input dto.RefreshTokenInput

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errors := utils.ValidateStruct(&input); len(errors) > 0 {
		return utils.ValidationErrorResponse(c, errors)
	}

	result, err := h.authService.RefreshToken(input.RefreshToken)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Token refreshed successfully", result)
}
