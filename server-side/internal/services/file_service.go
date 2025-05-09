package services

import (
	"cloud-drive/internal/models"
	"fmt"
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
	parentPublic := true

	if directory.ParentID != 0 {
		var parentDirectory models.DBDirectory
		if err := service.DB.Where("id = ?", directory.ParentID).First(&parentDirectory).Error; err == nil {
			// 判断是否有权限创建文件夹
			if parentDirectory.UserID != directory.UserID {
				return fmt.Errorf("没有权限创建文件夹")
			}

			parentPublic = parentDirectory.Public
		} else {
			return fmt.Errorf("父文件夹不存在")
		}
	}

	dbDirectory := directory.ToDBDirectory(parentPublic)
	dbDirectory.ParentPublic = parentPublic
	return service.DB.Create(dbDirectory).Error
}

func (service *FileService) DeleteDirectory(directoryID uint) error {
	return nil
}

func (service *FileService) GetFileTree(directoryID uint, userID uint) *models.APIFileTree {
	var dbDirectory models.DBDirectory
	var tree *models.APIFileTree = nil
	if err := service.DB.Where("id = ? and (user_id = ? or public = ?)", directoryID, userID, true).First(&dbDirectory).Error; err == nil {
		tree = &models.APIFileTree{}
		tree.ID = dbDirectory.ID
		tree.Name = dbDirectory.Name
		if dbDirectory.ParentID != 0 {
			tree.Parent = service.GetFileTree(dbDirectory.ParentID, userID)
		}
	}

	return tree
}

func (service *FileService) GetFiles(directoryID uint, userID uint) []*models.APIFile {
	var dbDirectories []*models.DBDirectory
	var files []*models.APIFile = []*models.APIFile{}

	// 处理共享的文件夹
	if err := service.DB.Preload("User").Where("parent_id = ? AND public = ? AND user_id != ?", directoryID, true, userID).Find(&dbDirectories).Error; err == nil {
		for _, dbDirectory := range dbDirectories {
			files = append(files, dbDirectory.ToAPIFile())
		}
	} else {
		log.Printf("Error querying public directories: %v", err)
	}
	if directoryID == 0 {
		// 处理父文件夹是私有但是子文件夹是公开的情况
		if err := service.DB.Preload("User").Where("public = ? and parent_public = ? and user_id != ?", true, false, userID).Find(&dbDirectories).Error; err == nil {
			for _, dbDirectory := range dbDirectories {
				files = append(files, dbDirectory.ToAPIFile())
			}
		} else {
			log.Printf("Error querying root public directories: %v", err)
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
