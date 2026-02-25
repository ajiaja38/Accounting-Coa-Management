package journal

import (
	"fmt"

	"fiber.com/session-api/pkg/model"
	"fiber.com/session-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetAll godoc
// @Summary      List all journal entries
// @Description  Returns a paginated list of journal entries
// @Tags         Journal
// @Produce      json
// @Param        page   query  int     true  "Page number"    minimum(1)
// @Param        limit  query  int     true  "Items per page" minimum(1) maximum(100)
// @Param        search query  string  false "Search by reference or description"
// @Success      200  {object}  model.SwaggerJournalListResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Failure      500  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /journal [get]
func (h *Handler) GetAll(c *fiber.Ctx) error {
	var req model.PaginationRequest
	if err := c.QueryParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	entries, meta, err := h.service.GetAll(&req)
	if err != nil {
		return err
	}

	return utils.SuccessResponsePaginate(c, fiber.StatusOK, "Success get all journal entries", entries, meta)
}

// GetByID godoc
// @Summary      Get journal entry by ID
// @Description  Returns a single journal entry with all its detail lines
// @Tags         Journal
// @Produce      json
// @Param        id   path  string  true  "Journal Entry ID (UUID)"
// @Success      200  {object}  model.SwaggerJournalResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Failure      404  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /journal/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid journal entry ID")
	}

	entry, err := h.service.GetByID(id)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fmt.Sprintf("Success get journal %s", entry.Reference), entry)
}

// Create godoc
// @Summary      Create a new journal entry
// @Description  Creates a journal entry with detail lines. This endpoint uses a DB transaction.
// @Tags         Journal
// @Accept       json
// @Produce      json
// @Param        body body CreateJournalRequest true "Journal entry payload"
// @Success      201  {object}  model.SwaggerJournalResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Failure      500  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /journal [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateJournalRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if len(req.Details) < 2 {
		return fiber.NewError(fiber.StatusBadRequest, "Journal entry must have at least 2 detail lines")
	}

	createdByStr := c.Locals("userId").(string)
	createdBy, err := uuid.Parse(createdByStr)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID in token")
	}

	tx, ok := c.Locals("tx").(*gorm.DB)
	if !ok || tx == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database transaction not available")
	}

	entry, err := h.service.Create(&req, createdBy, tx)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Journal entry created successfully", entry)
}

// PostJournal godoc
// @Summary      Post a draft journal entry
// @Description  Changes journal status from 'draft' to 'posted'
// @Tags         Journal
// @Produce      json
// @Param        id   path  string  true  "Journal Entry ID (UUID)"
// @Success      200  {object}  model.SwaggerEmptyResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Failure      404  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /journal/{id}/post [put]
func (h *Handler) PostJournal(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid journal entry ID")
	}

	if err := h.service.PostJournal(id); err != nil {
		return err
	}

	return utils.SuccessResponse[any](c, fiber.StatusOK, "Journal posted successfully", nil)
}

// Delete godoc
// @Summary      Delete a draft journal entry
// @Description  Soft-deletes a journal entry (only draft entries can be deleted)
// @Tags         Journal
// @Produce      json
// @Param        id   path  string  true  "Journal Entry ID (UUID)"
// @Success      200  {object}  model.SwaggerEmptyResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Failure      404  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /journal/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid journal entry ID")
	}

	if err := h.service.Delete(id); err != nil {
		return err
	}

	return utils.SuccessResponse[any](c, fiber.StatusOK, "Journal entry deleted successfully", nil)
}
