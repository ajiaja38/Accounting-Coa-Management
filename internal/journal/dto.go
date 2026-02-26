package journal

import (
	"time"
)

type JournalDetailRequest struct {
	CoaCode     string  `json:"coaCode"     validate:"required"  example:"5-1001"`
	Debit       float64 `json:"debit"       validate:"min=0"     example:"5000000"`
	Credit      float64 `json:"credit"      validate:"min=0"     example:"0"`
	Description string  `json:"description" validate:"omitempty" example:"Pembayaran gaji bulan Februari"`
}

type CreateJournalRequest struct {
	Description string                 `json:"description" validate:"omitempty" example:"Pembayaran gaji bulan Februari"`
	Details     []JournalDetailRequest `json:"details"     validate:"required,min=2,dive"`
}

type JournalDetailResponse struct {
	ID          string  `json:"id"`
	CoaCode     string  `json:"coaCode"`
	CoaName     string  `json:"coaName"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Description string  `json:"description"`
}

type JournalListResponse struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	Reference   string    `json:"reference"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"createdBy"`
	TotalDebit  float64   `json:"totalDebit"`
	TotalCredit float64   `json:"totalCredit"`
}

type JournalDetailedResponse struct {
	ID          string                  `json:"id"`
	Date        time.Time               `json:"date"`
	Reference   string                  `json:"reference"`
	Description string                  `json:"description"`
	Status      string                  `json:"status"`
	CreatedBy   string                  `json:"createdBy"`
	Details     []JournalDetailResponse `json:"details"`
}

// Swagger Responses

type SwaggerJournalResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    JournalDetailedResponse `json:"data"`
}
