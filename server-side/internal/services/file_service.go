package services

import (
	"cloud-drive/internal/models"
	"cloud-drive/permissions"
	"cloud-drive/utils"
	"fmt"
	"log"
	"path/filepath"

	"gorm.io/gorm"
)

type FileService struct {
	DB      *gorm.DB
	rootDir string
	fileDir string
	tempDir string
}

func NewFileService(db *gorm.DB, rootDir string) *FileService {
	fileDir := filepath.Join(rootDir, ".cloud_drive_files")
	tempDir := filepath.Join(fileDir, "temp")
	if err := utils.CreateDir(fileDir); err != nil {
		log.Fatalf("Failed to create file directory: %v", err)
	}
	if err := utils.CreateDir(tempDir); err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}

	return &FileService{
		DB:      db,
		rootDir: rootDir,
		fileDir: fileDir,
		tempDir: tempDir,
	}
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

func (service *FileService) UpdateDirectory(directoryID uint, directory *models.APIDirectory) error {
	var dbDirectory models.DBDirectory
	if err := service.DB.Where("id = ?", directoryID).First(&dbDirectory).Error; err != nil {
		return fmt.Errorf("文件夹不存在")
	}
	if dbDirectory.UserID != directory.UserID {
		return fmt.Errorf("没有权限更新文件夹")
	}

	parentPublic := true
	if dbDirectory.ParentID != 0 {
		var parentDirectory models.DBDirectory
		if err := service.DB.Where("id = ?", dbDirectory.ParentID).First(&parentDirectory).Error; err == nil {
			parentPublic = parentDirectory.Public
		}
	}
	dbDirectory.Permission = directory.Permission
	dbDirectory.Name = directory.Name
	dbDirectory.Public = permissions.CalculatePublic(parentPublic, directory.Permission)
	if err := service.DB.Save(&dbDirectory).Error; err != nil {
		return err
	}

	// 更新所有子文件夹的权限
	var updateChildError error = nil
	var childDirectories []models.DBDirectory
	if err := service.DB.Where("parent_id = ?", directoryID).Find(&childDirectories).Error; err == nil {
		for _, childDirectory := range childDirectories {
			childDirectory.Public = permissions.CalculatePublic(dbDirectory.Public, childDirectory.Permission)
			childDirectory.ParentPublic = dbDirectory.Public
			if err := service.DB.Save(&childDirectory).Error; err != nil {
				updateChildError = err
			}
		}
	}

	return updateChildError
}

func (service *FileService) DeleteDirectory(directoryID uint, userID uint) error {
	var dbDirectory models.DBDirectory
	if err := service.DB.Where("id = ?", directoryID).First(&dbDirectory).Error; err != nil {
		return fmt.Errorf("文件夹不存在")
	}
	if dbDirectory.UserID != userID {
		return fmt.Errorf("没有权限删除文件夹")
	}

	// 删除所有子文件
	var deleteFileError error = nil
	var dbFiles []models.DBFile
	if err := service.DB.Where("parent_id =?", directoryID).Find(&dbFiles).Error; err == nil {
		for _, dbFile := range dbFiles {
			// 删除文件
			filePath := filepath.Join(service.fileDir, dbFile.FileID)
			if err := utils.RemoveFile(filePath); err != nil {
				deleteFileError = err
			}
			if err := service.DB.Delete(&dbFile).Error; err != nil {
				deleteFileError = err
			}
		}
	}
	if deleteFileError != nil {
		return deleteFileError
	}

	// 删除所有子文件夹
	var deleteDirectoryError error = nil
	var childDirectories []models.DBDirectory
	if err := service.DB.Where("parent_id = ?", directoryID).Find(&childDirectories).Error; err == nil {
		for _, childDirectory := range childDirectories {
			if err := service.DeleteDirectory(childDirectory.ID, userID); err != nil {
				deleteDirectoryError = err
			}
		}
	}
	if deleteDirectoryError != nil {
		return deleteDirectoryError
	}

	return nil
}

func (service *FileService) GetFileTree(directoryID uint, userID uint) *models.APIFileTree {
	var dbDirectory models.DBDirectory
	var tree *models.APIFileTree = nil
	if err := service.DB.Where("id = ? and (user_id = ? or public = ?)", directoryID, userID, true).First(&dbDirectory).Error; err == nil {
		tree = &models.APIFileTree{}
		tree.ID = dbDirectory.ID
		tree.UserID = dbDirectory.UserID
		tree.Name = dbDirectory.Name
		tree.Public = dbDirectory.Public
		tree.Permission = dbDirectory.Permission
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
