package services

import (
	"cloud-drive/internal/models"
	"cloud-drive/permissions"
	"cloud-drive/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FileService struct {
	DB       *gorm.DB
	PathUtil *utils.PathUtil
}

func NewFileService(db *gorm.DB, pathUtil *utils.PathUtil) *FileService {
	return &FileService{
		DB:       db,
		PathUtil: pathUtil,
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

	// 更新所有子文件的权限
	var childFiles []models.DBFile
	if err := service.DB.Where("parent_id =?", directoryID).Find(&childFiles).Error; err == nil {
		for _, childFile := range childFiles {
			childFile.Public = permissions.CalculatePublic(dbDirectory.Public, childFile.Permission)
			childFile.ParentPublic = dbDirectory.Public
			if err := service.DB.Save(&childFile).Error; err != nil {
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
			filePath := filepath.Join(service.PathUtil.GetFileDir(), dbFile.FileID)
			if err := utils.RemoveFile(filePath); err != nil {
				deleteFileError = err
			}
			if err := service.DB.Unscoped().Delete(&dbFile).Error; err != nil {
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

	// 删除文件夹
	if err := service.DB.Unscoped().Delete(&dbDirectory).Error; err != nil {
		return err
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
	var files []*models.APIFile = []*models.APIFile{}

	// 处理共享的文件夹
	var dbDirectories []*models.DBDirectory
	if err := service.DB.Preload("User").Where("parent_id = ? AND public = ? AND user_id != ?", directoryID, true, userID).Find(&dbDirectories).Error; err == nil {
		for _, dbDirectory := range dbDirectories {
			files = append(files, dbDirectory.ToAPIFile())
		}
	} else {
		logrus.Errorf("Error querying public directories: %v", err)
	}
	if directoryID == 0 {
		// 处理父文件夹是私有但是子文件夹是公开的情况
		if err := service.DB.Preload("User").Where("public = ? and parent_public = ? and user_id != ?", true, false, userID).Find(&dbDirectories).Error; err == nil {
			for _, dbDirectory := range dbDirectories {
				files = append(files, dbDirectory.ToAPIFile())
			}
		} else {
			logrus.Errorf("Error querying root public directories: %v", err)
		}
	}

	// 查询指定目录下的所有文件夹
	if err := service.DB.Preload("User").Where("parent_id = ? AND user_id = ?", directoryID, userID).Find(&dbDirectories).Error; err == nil {
		for _, dbDirectory := range dbDirectories {
			files = append(files, dbDirectory.ToAPIFile())
		}
	} else {
		logrus.Errorf("Error querying directories: %v", err)
	}

	// 处理共享的文件
	var dbFiles []*models.DBFile
	if err := service.DB.Preload("User").Where("parent_id = ? AND public = ? AND user_id != ?", directoryID, true, userID).Find(&dbFiles).Error; err == nil {
		for _, dbFile := range dbFiles {
			files = append(files, dbFile.ToAPIFile())
		}
	} else {
		logrus.Errorf("Error querying public files: %v", err)
	}
	if directoryID == 0 {
		// 处理父文件夹是私有但是子文件是公开的情况
		if err := service.DB.Preload("User").Where("public = ? and parent_public = ? and user_id != ?", true, false, userID).Find(&dbFiles).Error; err == nil {
			for _, dbFile := range dbFiles {
				files = append(files, dbFile.ToAPIFile())
			}
		} else {
			logrus.Errorf("Error querying root public files: %v", err)
		}
	}
	// 查询指定目录下的所有文件
	if err := service.DB.Preload("User").Where("parent_id = ? AND user_id = ?", directoryID, userID).Find(&dbFiles).Error; err == nil {
		for _, dbFile := range dbFiles {
			files = append(files, dbFile.ToAPIFile())
		}
	} else {
		logrus.Errorf("Error querying files: %v", err)
	}

	return files
}

func (service *FileService) UploadFile(request *models.UploadFileRequest, userID uint) (string, error) {
	var parentPublic bool = true

	if *request.ParentID != 0 {
		var parentDirectory models.DBDirectory
		if err := service.DB.Where("id = ?", *request.ParentID).First(&parentDirectory).Error; err != nil {
			return "", fmt.Errorf("文件夹不存在")
		}
		if parentDirectory.UserID != userID {
			return "", fmt.Errorf("没有权限上传文件")
		}
		parentPublic = parentDirectory.Public
	}

	fileId := request.FileID
	if fileId == "" {
		fileId = uuid.New().String()
		var countfileId int64 = 0
		service.DB.Model(&models.DBFileChunk{}).Where("file_id = ?", fileId).Count(&countfileId)
		if countfileId != 0 {
			return "", fmt.Errorf("重复生成文件ID")
		}
	}
	var currentSize uint = uint(request.File.Size)

	var fileChunk models.DBFileChunk
	if err := service.DB.Where("file_id = ?", fileId).First(&fileChunk).Error; err != nil {
		// 写入临时表
		fileChunk = models.DBFileChunk{
			FileID:      fileId,
			TotalSize:   *request.Total,
			CurrentSize: currentSize,
		}
		if err := service.DB.Create(&fileChunk).Error; err != nil {
			return "", err
		}
	} else {
		currentSize += fileChunk.CurrentSize
		// 更新临时表
		fileChunk.TotalSize = *request.Total
		fileChunk.CurrentSize = currentSize
		if err := service.DB.Save(&fileChunk).Error; err != nil {
			return "", err
		}
	}

	// 写入临时文件
	tempFilePath := filepath.Join(service.PathUtil.GetTempDir(), fileId)
	file, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 写入临时文件
	source, err := request.File.Open()
	if err != nil {
		return "", err
	}
	defer source.Close()

	bytes, err := io.ReadAll(source)
	if err != nil {
		return "", err
	}
	if _, err := file.WriteAt(bytes, int64(*request.Index)); err != nil {
		return "", err
	}

	// 判断是否上传完成
	if currentSize == *request.Total {
		// 移动到正式目录
		file.Close()
		filePath := filepath.Join(service.PathUtil.GetFileDir(), fileId)
		if err := os.Rename(tempFilePath, filePath); err != nil {
			return "", err
		}
		// 写入正式表
		dbFile := &models.DBFile{
			Name:         request.Name,
			Size:         int64(currentSize),
			FileID:       fileId,
			UserID:       userID,
			ParentID:     *request.ParentID,
			Public:       permissions.CalculatePublic(parentPublic, *request.Permission),
			ParentPublic: parentPublic,
			Permission:   *request.Permission,
		}
		if err := service.DB.Create(dbFile).Error; err != nil {
			return "", err
		}
		// 从临时表里面移除
		if err := service.DB.Unscoped().Delete(&fileChunk).Error; err != nil {
			return "", err
		}
	}

	return fileId, nil
}

func (service *FileService) DeleteFile(id uint, userID uint) error {
	var dbFile models.DBFile
	if err := service.DB.Where("id = ?", id).First(&dbFile).Error; err != nil {
		return fmt.Errorf("文件不存在")
	}
	if dbFile.UserID != userID {
		return fmt.Errorf("没有权限删除文件")
	}

	// 删除文件
	filePath := filepath.Join(service.PathUtil.GetFileDir(), dbFile.FileID)
	if err := utils.RemoveFile(filePath); err != nil {
		return err
	}

	// 删除文件
	if err := service.DB.Unscoped().Delete(&dbFile).Error; err != nil {
		return err
	}

	return nil
}

func (service *FileService) UpdateFile(id uint, userID uint, request *models.UpdateFileRequest) error {
	var dbFile models.DBFile
	if err := service.DB.Where("id =?", id).First(&dbFile).Error; err != nil {
		return fmt.Errorf("文件不存在")
	}
	if dbFile.UserID != userID {
		return fmt.Errorf("没有权限更新文件")
	}
	var parentPublic bool = true
	if dbFile.ParentID != 0 {
		var parentDirectory models.DBDirectory
		if err := service.DB.Where("id =?", dbFile.ParentID).First(&parentDirectory).Error; err != nil {
			return fmt.Errorf("文件夹不存在")
		}
		parentPublic = parentDirectory.Public
	}
	dbFile.Name = request.Name
	dbFile.Public = permissions.CalculatePublic(parentPublic, *request.Permission)
	dbFile.Permission = *request.Permission
	if err := service.DB.Save(&dbFile).Error; err != nil {
		return err
	}

	return nil
}

func (service *FileService) DownloadFile(id uint, userID uint) (string, error) {
	var dbFile models.DBFile
	if err := service.DB.Where("id = ?", id).First(&dbFile).Error; err != nil {
		return "", fmt.Errorf("文件不存在")
	}
	if dbFile.UserID != userID && !dbFile.Public {
		return "", fmt.Errorf("没有权限下载文件")
	}

	fileUrl := filepath.Join(service.PathUtil.GetFileDir(), dbFile.FileID)
	if _, err := os.Stat(fileUrl); os.IsNotExist(err) {
		return "", fmt.Errorf("文件不存在")
	}

	return fileUrl, nil
}
