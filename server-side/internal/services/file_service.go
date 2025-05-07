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
		parentDirectory := &models.DBDirectory{}
		if err := service.DB.Where("id = ?", directory.ParentID).First(parentDirectory).Error; err == nil {
			parentPublic = parentDirectory.Public
		}
	}

	dbDirectory := directory.ToDBDirectory()
	dbDirectory.ParentPublic = parentPublic
	return service.DB.Create(dbDirectory).Error
}

func (service *FileService) GetFiles(parentId uint, userID uint) []*models.APIFile {
	// var dbDirectories []*models.DBDirectory
	var files []*models.APIFile

	// 处理共享的文件夹
	if parentId == 0 {
		rows, err := service.DB.Rows()

		if err == nil {
			columns, err := rows.Columns()
			if err == nil {
				// 处理错误
				log.Printf("Error querying root public directories: %v", columns)
			}
			for rows.Next() {
			}
		}
		// if err := service.DB.Where("public = ? AND parent_public = ?", true, false).Find(&dbDirectories).Error; err == nil {
		// 	for _, dbDirectory := range dbDirectories {
		// 		files = append(files, dbDirectory.ToAPIFile())
		// 	}
		// } else {
		// 	log.Printf("Error querying root public directories: %v", err)
		// }
	} else {
		// 处理当前文件夹下的共享文件夹
		// if err := service.DB.Where("parent_id = ? AND public = ? AND user_id != ?", directoryID, true, userID).Find(&dbDirectories).Error; err == nil {
		// 	for _, dbDirectory := range dbDirectories {
		// 		files = append(files, dbDirectory.ToAPIFile())
		// 	}
		// } else {
		// 	log.Printf("Error querying public directories: %v", err)
		// }
	}

	// 查询指定目录下的所有文件夹
	// if err := service.DB.Where("parent_id = ? AND user_id = ?", directoryID, userID).Find(&dbDirectories).Error; err == nil {
	// 	for _, dbDirectory := range dbDirectories {
	// 		files = append(files, dbDirectory.ToAPIFile())
	// 	}
	// } else {
	// 	log.Printf("Error querying directories: %v", err)
	// }

	return files
}
