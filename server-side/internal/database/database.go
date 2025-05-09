package database

import (
	"cloud-drive/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.DBUser{})
	db.AutoMigrate(&models.DBPassword{})
	db.AutoMigrate(&models.DBFile{})
	db.AutoMigrate(&models.DBDirectory{})
	db.AutoMigrate(&models.DBToken{})

	return db, nil
}
