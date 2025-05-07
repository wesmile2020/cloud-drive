package models

import "gorm.io/gorm"

// ServiceUser 用于和前端交互的用户模型
type ServiceUser struct {
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

// ToDBUser 将 ServiceUser 转换为 DBUser
func (user ServiceUser) ToDBUser() DBUser {
	return DBUser{
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}
}

// ToServiceUser 将 DBUser 转换为 ServiceUser
func (user DBUser) ToServiceUser() ServiceUser {
	return ServiceUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}
}

// RegisterUserRequest 用于接收前端注册用户的参数
type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginUserRequest 用于接收前端登录用户的参数
type LoginUserRequest struct {
	Account  string `json:"account" validate:"required"` // 可以是 phone 或 email
	Password string `json:"password" validate:"required"`
}
