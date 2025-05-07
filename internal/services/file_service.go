package services

import (
	"cloud-drive/internal/models"

	"gorm.io/gorm"
)

type FileService struct {
	DB      *gorm.DB
	rootDir string
}

func NewFileService(db *gorm.DB, rootDir string) *FileService {
	return &FileService{DB: db, rootDir: rootDir}
}

func (service *FileService) CreateDirectory(directory *models.ServiceDirectory) error {
	dbDirectory := directory.ToDBDirectory()
	return service.DB.Create(dbDirectory).Error
}
