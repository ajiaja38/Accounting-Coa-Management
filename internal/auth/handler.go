package auth

import (
	"time"

	"fiber.com/session-api/config"
	"fiber.com/session-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body RegisterRequest true "Register payload"
// @Success      201  {object}  model.SwaggerAuthResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      409  {object}  model.SwaggerErrorResponse
// @Failure      500  {object}  model.SwaggerErrorResponse
// @Router       /auth/register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := h.service.Register(&req)
	if err != nil {
		return err
	}

	resp := &AuthResponse{
		UserID:   user.ID.String(),
		UserName: user.UserName,
		Email:    user.Email,
		Role:     user.Role,
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "User registered successfully", resp)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates a user and sets an HttpOnly JWT cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body LoginRequest true "Login payload"
// @Success      200  {object}  model.SwaggerAuthResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Failure      500  {object}  model.SwaggerErrorResponse
// @Router       /auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	tokenStr, authResp, err := h.service.Login(&req)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    tokenStr,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Expires:  time.Now().Add(time.Duration(config.AppConfig.JWTExpiresHour) * time.Hour),
	})

	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", authResp)
}

// Logout godoc
// @Summary      Logout user
// @Description  Clears the JWT auth cookie
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  model.SwaggerEmptyResponse
// @Security     CookieAuth
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		HTTPOnly: true,
		Expires:  time.Now().Add(-time.Hour),
	})

	return utils.SuccessResponse[any](c, fiber.StatusOK, "Logout successful", nil)
}

// Me godoc
// @Summary      Get current authenticated user
// @Description  Returns information of the currently authenticated user
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  model.SwaggerAuthResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /auth/me [get]
func (h *Handler) Me(c *fiber.Ctx) error {
	resp := &AuthResponse{
		UserID:   c.Locals("userId").(string),
		UserName: c.Locals("userName").(string),
		Email:    c.Locals("email").(string),
		Role:     c.Locals("role").(string),
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success get current user", resp)
}
