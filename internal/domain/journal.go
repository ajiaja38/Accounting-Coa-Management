package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JournalStatus string

const (
	JournalStatusDraft  JournalStatus = "draft"
	JournalStatusPosted JournalStatus = "posted"
)

type JournalEntry struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Date        time.Time      `gorm:"type:date;not null"                             json:"date"`
	Reference   string         `gorm:"type:varchar(100);uniqueIndex;not null"          json:"reference"`
	Description string         `gorm:"type:text"                                      json:"description"`
	Status      JournalStatus  `gorm:"type:varchar(20);not null;default:'draft'"      json:"status"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null"                             json:"createdBy"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                                          json:"-"`
}
