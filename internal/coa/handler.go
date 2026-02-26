package coa

import (
	"fmt"

	"fiber.com/session-api/pkg/model"
	"fiber.com/session-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetAll godoc
// @Summary      List all Chart of Accounts
// @Description  Returns a paginated list of COAs with optional search by code or name
// @Tags         COA
// @Produce      json
// @Param        page   query  int     true  "Page number"    minimum(1)
// @Param        limit  query  int     true  "Items per page" minimum(1) maximum(100)
// @Param        search query  string  false "Search by name or code"
// @Success      200  {object}  model.SwaggerCOAListResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Failure      500  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /coa [get]
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

	accounts, meta, err := h.service.GetAll(&req)
	if err != nil {
		return err
	}

	return utils.SuccessResponsePaginate(c, fiber.StatusOK, "Success get all COA", accounts, meta)
}

// GetAllWithChildren godoc
// @Summary      List all Chart of Accounts with children
// @Description  Returns a paginated list of COAs with optional search by code or name or type
// @Tags         COA
// @Produce      json
// @Param        page   query  int     true  "Page number"    minimum(1)
// @Param        limit  query  int     true  "Items per page" minimum(1) maximum(100)
// @Param        search query  string  false "Search by name or code"
// @Success      200  {object}  model.SwaggerCOAListResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Failure      500  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /coa/with-children [get]
func (h *Handler) GetAllWithChildren(c *fiber.Ctx) error {
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

	accounts, meta, err := h.service.GetAllWithChildren(&req)
	if err != nil {
		return err
	}

	return utils.SuccessResponsePaginate(c, fiber.StatusOK, "Success get all COA", accounts, meta)
}

// GetByCode godoc
// @Summary      Get COA by code
// @Description  Returns a single Chart of Account by its code (e.g. "1-1001")
// @Tags         COA
// @Produce      json
// @Param        code  path  string  true  "COA Code (e.g. 1-1001)"
// @Success      200  {object}  model.SwaggerCOAResponse
// @Failure      404  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /coa/{code} [get]
func (h *Handler) GetByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "COA code is required")
	}

	coa, err := h.service.GetByCode(code)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fmt.Sprintf("Success get COA %s", coa.Code), coa)
}

// Create godoc
// @Summary      Create a new COA
// @Description  Creates a new Chart of Account. The code becomes the unique identifier (primary key).
// @Tags         COA
// @Accept       json
// @Produce      json
// @Param        body body CreateCOARequest true "COA payload"
// @Success      201  {object}  model.SwaggerCOAResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Failure      409  {object}  model.SwaggerErrorResponse
// @Failure      500  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /coa [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateCOARequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	coa, err := h.service.Create(&req)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "COA created successfully", coa)
}

// Update godoc
// @Summary      Update a COA
// @Description  Updates name, type, parentCode, or isActive of an existing COA by code
// @Tags         COA
// @Accept       json
// @Produce      json
// @Param        code  path  string          true  "COA Code (e.g. 1-1001)"
// @Param        body  body  UpdateCOARequest true "COA update payload"
// @Success      200  {object}  model.SwaggerCOAResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Failure      404  {object}  model.SwaggerErrorResponse
// @Failure      500  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /coa/{code} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "COA code is required")
	}

	var req UpdateCOARequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	coa, err := h.service.Update(code, &req)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fmt.Sprintf("COA %s updated successfully", coa.Code), coa)
}

// Delete godoc
// @Summary      Delete a COA
// @Description  Soft-deletes a Chart of Account by code
// @Tags         COA
// @Produce      json
// @Param        code  path  string  true  "COA Code (e.g. 1-1001)"
// @Success      200  {object}  model.SwaggerEmptyResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Failure      404  {object}  model.SwaggerErrorResponse
// @Failure      500  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /coa/{code} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "COA code is required")
	}

	if err := h.service.Delete(code); err != nil {
		return err
	}

	return utils.SuccessResponse[any](c, fiber.StatusOK, "COA deleted successfully", nil)
}
