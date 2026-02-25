package domain

import (
	"time"

	"gorm.io/gorm"
)

type AccountType string

const (
	AccountTypeAsset     AccountType = "asset"
	AccountTypeLiability AccountType = "liability"
	AccountTypeEquity    AccountType = "equity"
	AccountTypeRevenue   AccountType = "revenue"
	AccountTypeExpense   AccountType = "expense"
)

type ChartOfAccount struct {
	Code       string         `gorm:"type:varchar(20);primaryKey"   json:"code"`
	Name       string         `gorm:"type:varchar(200);not null"    json:"name"`
	Type       AccountType    `gorm:"type:varchar(20);not null"     json:"type"`
	ParentCode *string        `gorm:"type:varchar(20);index"        json:"parentCode,omitempty"`
	IsActive   bool           `gorm:"not null;default:true"          json:"isActive"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index"                         json:"-"`
}
