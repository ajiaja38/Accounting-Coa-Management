package journal

import (
	"fmt"
	"math"
	"strings"
	"time"

	"fiber.com/session-api/internal/domain"
	"fiber.com/session-api/pkg/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetAll(req *model.PaginationRequest) ([]JournalListResponse, *model.MetaPagination, error)
	GetByID(id uuid.UUID) (*JournalDetailedResponse, error)
	Create(req *CreateJournalRequest, createdBy uuid.UUID, tx *gorm.DB) (*JournalDetailedResponse, error)
	PostJournal(id uuid.UUID) error
	Delete(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(req *model.PaginationRequest) ([]JournalListResponse, *model.MetaPagination, error) {
	entries, total, err := s.repo.FindAll(req)
	if err != nil {
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))
	meta := &model.MetaPagination{
		Page:      req.Page,
		Limit:     req.Limit,
		TotalPage: totalPage,
		TotalData: int(total),
	}

	return entries, meta, nil
}

func (s *service) GetByID(id uuid.UUID) (*JournalDetailedResponse, error) {
	entry, details, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if entry == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Journal entry not found")
	}

	detailResponses := make([]JournalDetailResponse, len(details))
	for i, d := range details {
		detailResponses[i] = JournalDetailResponse{
			ID:          d.ID,
			CoaCode:     d.CoaCode,
			CoaName:     d.CoaName,
			Debit:       d.Debit,
			Credit:      d.Credit,
			Description: d.Description,
		}
	}

	return &JournalDetailedResponse{
		ID:          entry.ID.String(),
		Date:        entry.Date,
		Reference:   entry.Reference,
		Description: entry.Description,
		Status:      string(entry.Status),
		CreatedBy:   entry.CreatedBy.String(),
		Details:     detailResponses,
	}, nil
}

func (s *service) Create(req *CreateJournalRequest, createdBy uuid.UUID, tx *gorm.DB) (*JournalDetailedResponse, error) {
	txRepo := NewRepository(tx)

	entryID := uuid.New()

	now := time.Now()
	refSuffix := strings.ToUpper(uuid.New().String()[0:4])
	reference := fmt.Sprintf("JRN-%s-%s", now.Format("20060102"), refSuffix)

	entry := &domain.JournalEntry{
		ID:          entryID,
		Date:        now,
		Reference:   reference,
		Description: req.Description,
		Status:      domain.JournalStatusDraft,
		CreatedBy:   createdBy,
	}

	details := make([]domain.JournalEntryDetail, len(req.Details))
	for i, d := range req.Details {
		details[i] = domain.JournalEntryDetail{
			JournalEntryID: entryID.String(),
			CoaCode:        d.CoaCode,
			Debit:          d.Debit,
			Credit:         d.Credit,
			Description:    d.Description,
		}
	}

	if err := txRepo.Create(entry, details); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	entryResult, detailsResult, err := txRepo.FindByID(entryID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if entryResult == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Journal entry not found (transaction issue)")
	}

	detailResponses := make([]JournalDetailResponse, len(detailsResult))
	for i, d := range detailsResult {
		detailResponses[i] = JournalDetailResponse{
			ID:          d.ID,
			CoaCode:     d.CoaCode,
			CoaName:     d.CoaName,
			Debit:       d.Debit,
			Credit:      d.Credit,
			Description: d.Description,
		}
	}

	result := &JournalDetailedResponse{
		ID:          entryResult.ID.String(),
		Date:        entryResult.Date,
		Reference:   entryResult.Reference,
		Description: entryResult.Description,
		Status:      string(entryResult.Status),
		CreatedBy:   entryResult.CreatedBy.String(),
		Details:     detailResponses,
	}

	return result, nil
}

func (s *service) PostJournal(id uuid.UUID) error {
	entry, _, err := s.repo.FindByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if entry == nil {
		return fiber.NewError(fiber.StatusNotFound, "Journal entry not found")
	}

	return s.repo.PostJournal(id)
}

func (s *service) Delete(id uuid.UUID) error {
	entry, _, err := s.repo.FindByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if entry == nil {
		return fiber.NewError(fiber.StatusNotFound, "Journal entry not found")
	}

	return s.repo.Delete(id)
}
