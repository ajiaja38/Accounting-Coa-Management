package coa

import (
	"math"

	"fiber.com/session-api/internal/domain"
	"fiber.com/session-api/pkg/model"

	"github.com/gofiber/fiber/v2"
)

// Service defines the business logic layer for COA.
type Service interface {
	GetAll(req *model.PaginationRequest) ([]COAResponse, *model.MetaPagination, error)
	GetByCode(code string) (*COAResponse, error)
	Create(req *CreateCOARequest) (*COAResponse, error)
	Update(code string, req *UpdateCOARequest) (*COAResponse, error)
	Delete(code string) error
}

type service struct {
	repo Repository
}

// NewService creates a new COA service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func toResponse(c *domain.ChartOfAccount) *COAResponse {
	return &COAResponse{
		Code:       c.Code,
		Name:       c.Name,
		Type:       c.Type,
		ParentCode: c.ParentCode,
		IsActive:   c.IsActive,
	}
}

func (s *service) GetAll(req *model.PaginationRequest) ([]COAResponse, *model.MetaPagination, error) {
	accounts, total, err := s.repo.FindAll(req)
	if err != nil {
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	responses := make([]COAResponse, len(accounts))
	for i, a := range accounts {
		responses[i] = *toResponse(&a)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))
	meta := &model.MetaPagination{
		Page:      req.Page,
		Limit:     req.Limit,
		TotalPage: totalPage,
		TotalData: int(total),
	}

	return responses, meta, nil
}

func (s *service) GetByCode(code string) (*COAResponse, error) {
	coa, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if coa == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "COA not found")
	}
	return toResponse(coa), nil
}

func (s *service) Create(req *CreateCOARequest) (*COAResponse, error) {
	existing, err := s.repo.FindByCode(req.Code)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if existing != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "COA code already exists")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	coa := &domain.ChartOfAccount{
		Code:       req.Code,
		Name:       req.Name,
		Type:       req.Type,
		ParentCode: req.ParentCode,
		IsActive:   isActive,
	}

	if err := s.repo.Create(coa); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	created, err := s.repo.FindByCode(req.Code)
	if err != nil || created == nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch created COA")
	}

	return toResponse(created), nil
}

func (s *service) Update(code string, req *UpdateCOARequest) (*COAResponse, error) {
	existing, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if existing == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "COA not found")
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Type != "" {
		existing.Type = req.Type
	}
	if req.ParentCode != nil {
		existing.ParentCode = req.ParentCode
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return toResponse(existing), nil
}

func (s *service) Delete(code string) error {
	return s.repo.Delete(code)
}
