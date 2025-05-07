package models

import (
	"gorm.io/gorm"
)

type APIFile struct {
	User        APIUser // 用户信息
	ID          uint    `json:"id"`          // 文件ID
	Name        string  `json:"name"`        // 文件名称
	Size        int64   `json:"size"`        // 文件大小，单位为字节
	ParentID    uint    `json:"parentId"`    // 父文件夹ID
	Public      bool    `json:"public"`      // 是否公开
	IsDirectory bool    `json:"isDirectory"` // 是否是目录
}

type DBFile struct {
	gorm.Model
	Name         string `gorm:"not null"` // 文件名称
	Size         int64  `gorm:"not null"` // 文件大小，单位为字节
	FileID       uint   `gorm:"not null"` // 文件ID
	UserID       uint   `gorm:"not null"` // 用户ID
	ParentID     uint   `gorm:"not null"` // 父文件夹ID
	Public       bool   `gorm:"not null"` // 是否公开
	ParentPublic bool   `gorm:"not null"` // 父文件夹是否公开
}

func (db *DBFile) TableName() string {
	return "file"
}

type APIDirectory struct {
	Name     string `json:"name"`     // 文件夹名称
	UserID   uint   `json:"userId"`   // 用户ID
	ParentID uint   `json:"parentId"` // 父文件夹ID
	Public   bool   `json:"public"`   // 是否公开
}

func (directory *APIDirectory) ToDBDirectory() *DBDirectory {
	return &DBDirectory{
		Name:     directory.Name,
		UserID:   directory.UserID,
		ParentID: directory.ParentID,
		Public:   directory.Public,
	}
}

type DBDirectory struct {
	gorm.Model
	Name         string `gorm:"not null"` // 文件夹名称
	UserID       uint   `gorm:"not null"` // 用户ID
	ParentID     uint   `gorm:"not null"` // 父文件夹ID
	Public       bool   `gorm:"not null"` // 是否公开
	ParentPublic bool   `gorm:"not null"` // 父文件夹是否公开
}

func (db *DBDirectory) TableName() string {
	return "directory"
}

type CreateDirectoryRequest struct {
	Name     string `json:"name" binding:"required"`     // 文件夹名称
	ParentID *uint  `json:"parentId" binding:"required"` // 父文件夹ID
	Public   *bool  `json:"public" binding:"required"`   // 是否公开
}
