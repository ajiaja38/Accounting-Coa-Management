package coa

import (
	"fiber.com/session-api/internal/domain"
)

type CreateCOARequest struct {
	Code       string             `json:"code"       validate:"required"                                          example:"1-1001"`
	Name       string             `json:"name"       validate:"required"                                          example:"Kas dan Setara Kas"`
	Type       domain.AccountType `json:"type"       validate:"required,oneof=asset liability equity revenue expense" example:"asset"`
	ParentCode *string            `json:"parentCode" validate:"omitempty"                                         example:"1-1000"`
	IsActive   *bool              `json:"isActive"                                                                example:"true"`
}

type UpdateCOARequest struct {
	Name       string             `json:"name"       validate:"omitempty" example:"Kas dan Setara Kas"`
	Type       domain.AccountType `json:"type"       validate:"omitempty,oneof=asset liability equity revenue expense" example:"asset"`
	ParentCode *string            `json:"parentCode" validate:"omitempty" example:"1-1000"`
	IsActive   *bool              `json:"isActive"                        example:"true"`
}

type COAResponse struct {
	Code       string             `json:"code"`
	Name       string             `json:"name"`
	Type       domain.AccountType `json:"type"`
	ParentCode *string            `json:"parentCode,omitempty"`
	IsActive   bool               `json:"isActive"`
}
