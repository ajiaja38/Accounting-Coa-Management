package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserName  string         `gorm:"type:varchar(100);uniqueIndex;not null"          json:"userName"`
	Email     string         `gorm:"type:varchar(150);uniqueIndex;not null"          json:"email"`
	Password  string         `gorm:"type:varchar(255);not null"                     json:"-"`
	Role      string         `gorm:"type:varchar(20);not null;default:'user'"       json:"role"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                          json:"-"`
}
