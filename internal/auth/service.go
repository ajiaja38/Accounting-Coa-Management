package auth

import (
	"fiber.com/session-api/internal/domain"
	"fiber.com/session-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req *RegisterRequest) (*domain.User, error)
	Login(req *LoginRequest) (string, *AuthResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(req *RegisterRequest) (*domain.User, error) {
	emailExists, err := s.repo.EmailExists(req.Email)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if emailExists {
		return nil, fiber.NewError(fiber.StatusConflict, "Email already registered")
	}

	userNameExists, err := s.repo.UserNameExists(req.UserName)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if userNameExists {
		return nil, fiber.NewError(fiber.StatusConflict, "Username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user := &domain.User{
		ID:       uuid.New(),
		UserName: req.UserName,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return user, nil
}

func (s *service) Login(req *LoginRequest) (string, *AuthResponse, error) {
	user, err := s.repo.FindUserByEmail(req.Email)

	if err != nil {
		return "", nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return "", nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	tokenStr, err := utils.GenerateToken(user.ID.String(), user.UserName, user.Email, user.Role)

	if err != nil {
		return "", nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	authResp := &AuthResponse{
		UserID:   user.ID.String(),
		UserName: user.UserName,
		Email:    user.Email,
		Role:     user.Role,
	}

	return tokenStr, authResp, nil
}
