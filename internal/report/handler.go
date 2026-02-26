package report

import (
	"fiber.com/session-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetLedger godoc
// @Summary      Get General Ledger
// @Description  Get General Ledger transactions for a specific COA
// @Tags         Report
// @Produce      json
// @Param        coaCode   query     string  true  "COA Code"
// @Param        startDate query     string  false "Start Date (YYYY-MM-DD)"
// @Param        endDate   query     string  false "End Date (YYYY-MM-DD)"
// @Success      200  {object}  SwaggerLedgerResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /report/ledger [get]
func (h *Handler) GetLedger(c *fiber.Ctx) error {
	req := new(LedgerQuery)
	if err := c.QueryParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	if req.CoaCode == "" {
		return fiber.NewError(fiber.StatusBadRequest, "coaCode is required")
	}

	res, err := h.service.GetLedger(req)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "General ledger fetched successfully", res)
}

// GetTrialBalance godoc
// @Summary      Get Trial Balance
// @Description  Get Trial Balance report for a specific period
// @Tags         Report
// @Produce      json
// @Param        startDate query     string  false "Start Date (YYYY-MM-DD)"
// @Param        endDate   query     string  false "End Date (YYYY-MM-DD)"
// @Success      200  {object}  SwaggerTrialBalanceResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /report/trial-balance [get]
func (h *Handler) GetTrialBalance(c *fiber.Ctx) error {
	req := new(PeriodQuery)
	if err := c.QueryParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	res, err := h.service.GetTrialBalance(req)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Trial balance fetched successfully", res)
}

// GetProfitLoss godoc
// @Summary      Get Profit & Loss
// @Description  Get Profit & Loss report for a specific period (Income Statement)
// @Tags         Report
// @Produce      json
// @Param        startDate query     string  false "Start Date (YYYY-MM-DD)"
// @Param        endDate   query     string  false "End Date (YYYY-MM-DD)"
// @Success      200  {object}  SwaggerProfitLossResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /report/profit-loss [get]
func (h *Handler) GetProfitLoss(c *fiber.Ctx) error {
	req := new(PeriodQuery)
	if err := c.QueryParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	res, err := h.service.GetProfitLoss(req)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profit & Loss fetched successfully", res)
}

// GetBalanceSheet godoc
// @Summary      Get Balance Sheet
// @Description  Get Balance Sheet report up to a specific date (Financial Position)
// @Tags         Report
// @Produce      json
// @Param        endDate   query     string  false "End Date (YYYY-MM-DD)"
// @Success      200  {object}  SwaggerBalanceSheetResponse
// @Failure      400  {object}  model.SwaggerErrorResponse
// @Failure      401  {object}  model.SwaggerErrorResponse
// @Security     CookieAuth
// @Router       /report/balance-sheet [get]
func (h *Handler) GetBalanceSheet(c *fiber.Ctx) error {
	req := new(PeriodQuery)
	if err := c.QueryParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	res, err := h.service.GetBalanceSheet(req)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Balance Sheet fetched successfully", res)
}
