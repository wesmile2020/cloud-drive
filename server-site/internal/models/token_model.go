package models

import (
	"time"

	"gorm.io/gorm"
)

type DBToken struct {
	gorm.Model
	Token     string    `gorm:"unique;not null"`
	UserID    uint      `gorm:"not null"`
	ExpiredAt time.Time `gorm:"not null"`
}

func (DBToken) TableName() string {
	return "token"
}
