package models

import (
	"time"

	"gorm.io/gorm"
)

// APIUser 用于和前端交互的用户模型
type APIUser struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// DBUser 用于和数据库交互的用户模型
type DBUser struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Email string `gorm:"not null;unique"`
	Phone string `gorm:"not null;unique"`
}

// TableName 自定义 DBUser 表名
func (user DBUser) TableName() string {
	return "user"
}

// DBPassword 用于和数据库交互的密码表模型
type DBPassword struct {
	gorm.Model
	UserID   uint   `gorm:"not null;unique"` // 关联用户 ID
	Password string `gorm:"not null"`        // 用户密码
}

// TableName 自定义 DBPassword 表名
func (password DBPassword) TableName() string {
	return "password"
}

// ToDBUser 将 APIUser 转换为 DBUser
func (user APIUser) ToDBUser() DBUser {
	return DBUser{
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}
}

// ToAPIUser 将 DBUser 转换为 APIUser
func (user DBUser) ToAPIUser() APIUser {
	return APIUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}
}

type DBVerifyCode struct {
	gorm.Model
	Email     string    `gorm:"not null;unique"` // 关联用户邮箱
	Code      string    `gorm:"not null"`        // 验证码
	ExpiredAt time.Time `gorm:"not null"`        // 过期时间
}

// TableName 自定义 DBVerifyCode 表名
func (code DBVerifyCode) TableName() string {
	return "verify_code"
}

// RegisterUserRequest 用于接收前端注册用户的参数
type RegisterUserRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginUserRequest 用于接收前端登录用户的参数
type LoginUserRequest struct {
	Account  string `json:"account" binding:"required"` // 可以是 phone 或 email
	Password string `json:"password" binding:"required"`
}

type EditUserInfoRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=50"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type GetVerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RetrievePasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
