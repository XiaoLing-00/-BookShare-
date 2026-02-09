package models

import (
	"time"

	"gorm.io/gorm"
)

type UserBookRelation struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null"`
	User         User           `json:"user"`
	BookID       uint           `json:"book_id" gorm:"not null"`
	Book         Book           `json:"book"`
	RelationType string         `json:"relation_type" gorm:"not null;type:varchar(20)"` // collected, read
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}