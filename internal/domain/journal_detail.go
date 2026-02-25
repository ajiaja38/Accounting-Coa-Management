package domain

import (
	"gorm.io/gorm"
)

type JournalEntryDetail struct {
	ID             string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	JournalEntryID string         `gorm:"type:uuid;not null;index"                       json:"journalEntryId"`
	CoaCode        string         `gorm:"type:varchar(20);not null;index"                json:"coaCode"`
	Debit          float64        `gorm:"type:numeric(20,2);not null;default:0"          json:"debit"`
	Credit         float64        `gorm:"type:numeric(20,2);not null;default:0"          json:"credit"`
	Description    string         `gorm:"type:text"                                      json:"description"`
	DeletedAt      gorm.DeletedAt `gorm:"index"                                          json:"-"`
}
