package models

import (
	"cloud-drive/permissions"

	"gorm.io/gorm"
)

type APIFileTree struct {
	ID         uint         `json:"id"`         // 文件ID
	UserID     uint         `json:"userId"`     // 用户ID
	Name       string       `json:"name"`       // 文件名称
	Public     bool         `json:"public"`     // 是否公开
	Permission uint         `json:"permission"` // 权限 0:私有 1:继承父目录的权限 2:公开
	Parent     *APIFileTree `json:"parent"`     // 父文件
}

type APIFile struct {
	User        APIUser `json:"user"`        // 用户信息
	ID          uint    `json:"id"`          // 文件ID
	Name        string  `json:"name"`        // 文件名称
	Size        int64   `json:"size"`        // 文件大小，单位为字节
	FileID      uint    `json:"fileId"`      // 文件ID
	ParentID    uint    `json:"parentId"`    // 父文件夹ID
	Public      bool    `json:"public"`      // 是否公开
	Permission  uint    `json:"permission"`  // 权限 0:私有 1:继承父目录的权限 2:公开
	IsDirectory bool    `json:"isDirectory"` // 是否是目录
	Timestamp   int64   `json:"timestamp"`   // 时间戳，单位为秒
}

type DBFile struct {
	gorm.Model
	Name         string `gorm:"not null"` // 文件名称
	Size         int64  `gorm:"not null"` // 文件大小，单位为字节
	FileID       uint   `gorm:"not null"` // 文件ID
	UserID       uint   `gorm:"not null"` // 用户ID
	User         DBUser `gorm:"foreignKey:UserID"`
	ParentID     uint   `gorm:"not null"` // 父文件夹ID
	Public       bool   `gorm:"not null"` // 是否公开
	ParentPublic bool   `gorm:"not null"` // 父文件夹是否公开
	Permission   uint   `gorm:"not null"` // 权限 0:私有 1:继承父目录的权限 2:公开
}

func (db *DBFile) TableName() string {
	return "file"
}

type APIDirectory struct {
	Name       string `json:"name"`       // 文件夹名称
	UserID     uint   `json:"userId"`     // 用户ID
	ParentID   uint   `json:"parentId"`   // 父文件夹ID
	Permission uint   `json:"permission"` // 权限 0:私有 1:继承父目录的权限 2:公开
}

func (directory *APIDirectory) ToDBDirectory(parentPublic bool) *DBDirectory {
	return &DBDirectory{
		Name:         directory.Name,
		UserID:       directory.UserID,
		ParentID:     directory.ParentID,
		Permission:   directory.Permission,
		Public:       permissions.CalculatePublic(parentPublic, directory.Permission),
		ParentPublic: parentPublic,
	}
}

type DBDirectory struct {
	gorm.Model
	Name         string `gorm:"not null"` // 文件夹名称
	UserID       uint   `gorm:"not null"` // 用户ID
	User         DBUser `gorm:"foreignKey:UserID;references:ID"`
	ParentID     uint   `gorm:"not null"` // 父文件夹ID
	Public       bool   `gorm:"not null"` // 是否公开
	ParentPublic bool   `gorm:"not null"` // 父文件夹是否公开
	Permission   uint   `gorm:"not null"` // 权限 0:私有 1:继承父目录的权限 2:公开
}

func (directory *DBDirectory) TableName() string {
	return "directory"
}

func (directory *DBDirectory) ToAPIFile() *APIFile {
	return &APIFile{
		User:        directory.User.ToAPIUser(),
		ID:          directory.ID,
		Name:        directory.Name,
		ParentID:    directory.ParentID,
		Public:      directory.Public,
		Permission:  directory.Permission,
		IsDirectory: true,
		Timestamp:   directory.UpdatedAt.Unix(),
	}
}

type CreateDirectoryRequest struct {
	Name       string `json:"name" binding:"required"`       // 文件夹名称
	ParentID   *uint  `json:"parentId" binding:"required"`   // 父文件夹ID
	Permission *uint  `json:"permission" binding:"required"` // 权限 0:私有 1:继承父目录的权限 2:公开
}

type UpdateDirectoryRequest struct {
	Name       string `json:"name" binding:"required"`       // 文件夹名称
	Permission *uint  `json:"permission" binding:"required"` // 权限 0:私有 1:继承父目录的权限 2:公开
}
