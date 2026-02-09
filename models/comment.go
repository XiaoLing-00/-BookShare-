package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	BookID    uint           `json:"book_id" gorm:"not null"`
	Book      Book           `json:"book"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user"`
	Content   string         `json:"content" gorm:"not null;type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}