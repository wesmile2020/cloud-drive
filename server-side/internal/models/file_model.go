package models

import "gorm.io/gorm"

type DBFile struct {
	gorm.Model
	Name     string `gorm:"not null"` // 文件名称
	Size     int64  `gorm:"not null"` // 文件大小，单位为字节
	FileID   uint   `gorm:"not null"` // 文件ID
	UserID   uint   `gorm:"not null"` // 用户ID
	ParentID uint   `gorm:"not null"` // 父文件夹ID
	Public   bool   `gorm:"not null"` // 是否公开
}

func (db *DBFile) TableName() string {
	return "file"
}

type ServiceDirectory struct {
	Name     string `json:"name"`     // 文件夹名称
	UserID   uint   `json:"user_id"`  // 用户ID
	ParentID uint   `json:"parentId"` // 父文件夹ID
	Public   bool   `json:"public"`   // 是否公开
}

func (directory *ServiceDirectory) ToDBDirectory() *DBDirectory {
	return &DBDirectory{
		Name:     directory.Name,
		UserID:   directory.UserID,
		ParentID: directory.ParentID,
		Public:   directory.Public,
	}
}

type DBDirectory struct {
	gorm.Model
	Name     string `gorm:"not null"` // 文件夹名称
	UserID   uint   `gorm:"not null"` // 用户ID
	ParentID uint   `gorm:"not null"` // 父文件夹ID
	Public   bool   `gorm:"not null"` // 是否公开
}

func (db *DBDirectory) TableName() string {
	return "directory"
}

type CreateDirectoryRequest struct {
	Name     string `json:"name" binding:"required"`     // 文件夹名称
	ParentID *uint  `json:"parentId" binding:"required"` // 父文件夹ID
	Public   *bool  `json:"public" binding:"required"`   // 是否公开
}
