package models

import "gorm.io/gorm"

type DBToken struct {
	gorm.Model
	Token     string `gorm:"unique;not null"`
	UserID    uint   `gorm:"not null"`
	ExpiredAt int64  `gorm:"not null"`
}

func (DBToken) TableName() string {
	return "token"
}
