package services

import (
	"cloud-drive/internal/models"
	"log"

	"gorm.io/gorm"
)

type FileService struct {
	DB      *gorm.DB
	rootDir string
}

func NewFileService(db *gorm.DB, rootDir string) *FileService {
	return &FileService{DB: db, rootDir: rootDir}
}

func (service *FileService) CreateDirectory(directory *models.APIDirectory) error {
	parentPublic := false

	if directory.ParentID != 0 {
		var parentDirectory models.DBDirectory
		if err := service.DB.Where("id = ?", directory.ParentID).First(&parentDirectory).Error; err == nil {
			parentPublic = parentDirectory.Public
		}
	}

	dbDirectory := directory.ToDBDirectory()
	dbDirectory.ParentPublic = parentPublic
	return service.DB.Create(dbDirectory).Error
}

func (service *FileService) DeleteDirectory(directoryID uint) error {
	return nil
}

func (service *FileService) GetFileTree(directoryID uint) *models.APIFileTree {
	var dbDirectory models.DBDirectory
	var tree *models.APIFileTree = nil
	if err := service.DB.Where("id =?", directoryID).First(&dbDirectory).Error; err == nil {
		tree = &models.APIFileTree{}
		tree.ID = dbDirectory.ID
		tree.Name = dbDirectory.Name
		if dbDirectory.ParentID != 0 {
			tree.Parent = service.GetFileTree(dbDirectory.ParentID)
		}
	}

	return tree
}

func (service *FileService) GetFiles(directoryID uint, userID uint) []*models.APIFile {
	var dbDirectories []*models.DBDirectory
	var files []*models.APIFile = []*models.APIFile{}

	// 处理共享的文件夹
	if directoryID == 0 {
		if err := service.DB.Preload("User").Where("public = ? and parent_public = ? and user_id != ?", true, false, userID).Find(&dbDirectories).Error; err == nil {
			for _, dbDirectory := range dbDirectories {
				files = append(files, dbDirectory.ToAPIFile())
			}
		} else {
			log.Printf("Error querying root public directories: %v", err)
		}
	} else {
		// 处理当前文件夹下的共享文件夹
		if err := service.DB.Preload("User").Where("parent_id = ? AND public = ? AND user_id != ?", directoryID, true, userID).Find(&dbDirectories).Error; err == nil {
			for _, dbDirectory := range dbDirectories {
				files = append(files, dbDirectory.ToAPIFile())
			}
		} else {
			log.Printf("Error querying public directories: %v", err)
		}
	}

	// 查询指定目录下的所有文件夹
	if err := service.DB.Preload("User").Where("parent_id = ? AND user_id = ?", directoryID, userID).Find(&dbDirectories).Error; err == nil {
		for _, dbDirectory := range dbDirectories {
			files = append(files, dbDirectory.ToAPIFile())
		}
	} else {
		log.Printf("Error querying directories: %v", err)
	}

	return files
}
