package services

import (
	"cloud-drive/internal/models"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// UserService 用户服务结构体
type UserService struct {
	DB *gorm.DB
}

// NewUserService 创建一个新的 UserService 实例
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// RegisterUser 注册新用户
func (service *UserService) RegisterUser(apiUser *models.APIUser, password string) error {
	var existCount int64 = 0
	// 检查phone是否已存在
	service.DB.Table("user").Where("phone = ?", apiUser.Phone).Count(&existCount)
	if existCount > 0 {
		return fmt.Errorf("手机号已被注册")
	}

	// 检查email是否已存在
	service.DB.Table("user").Where("email = ?", apiUser.Email).Count(&existCount)
	if existCount > 0 {
		return fmt.Errorf("邮箱已被注册")
	}

	// 插入用户表
	dbUser := apiUser.ToDBUser()
	if err := service.DB.Create(&dbUser).Error; err != nil {
		return err
	}

	// 对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 插入密码表
	dbPassword := models.DBPassword{UserID: dbUser.ID, Password: string(hashedPassword)}
	if err := service.DB.Create(&dbPassword).Error; err != nil {
		return err
	}

	return nil
}

// LoginUser 用户登录
func (service *UserService) LoginUser(account, password string) (*models.APIUser, error) {
	// 根据 account 查询用户，account 可能是 phone 或 email
	var existingUser models.DBUser
	if err := service.DB.Where("phone = ? OR email = ?", account, account).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("账号或者密码错误")
		}
		return nil, err
	}

	// 查询用户密码
	var dbPassword models.DBPassword
	if err := service.DB.Where("user_id = ?", existingUser.ID).First(&dbPassword).Error; err != nil {
		return nil, err
	}

	// 对比密码
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("账号或者密码错误")
	}

	apiUser := existingUser.ToAPIUser()

	return &apiUser, nil
}

func (service *UserService) SaveToken(tokenString string, userID uint, expiredAt time.Time) error {
	dbToken := models.DBToken{
		Token:     tokenString,
		UserID:    userID,
		ExpiredAt: expiredAt,
	}
	return service.DB.Create(&dbToken).Error
}

func (service *UserService) Logout(tokenString string) error {
	return service.DB.Unscoped().Delete(&models.DBToken{}, "token = ?", tokenString).Error
}

// 获取用户登录信息
func (service *UserService) GetUserInfo(userID uint) (*models.APIUser, error) {
	var dbUser models.DBUser
	if err := service.DB.First(&dbUser, userID).Error; err != nil {
		return nil, err
	}

	apiUser := dbUser.ToAPIUser()

	return &apiUser, nil
}

func (service *UserService) EditUserInfo(userID uint, apiUser *models.APIUser) error {
	var dbUser models.DBUser
	if err := service.DB.First(&dbUser, userID).Error; err != nil {
		return err
	}
	dbUser.Name = apiUser.Name
	dbUser.Email = apiUser.Email
	dbUser.Phone = apiUser.Phone
	return service.DB.Save(&dbUser).Error
}

func (service *UserService) EditPassword(userID uint, oldPassword, newPassword string) error {
	var dbPassword models.DBPassword
	if err := service.DB.Where("user_id = ?", userID).First(&dbPassword).Error; err != nil {
		return err
	}

	// 对比旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("旧密码错误")
	}

	// 对新密码进行加密
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	dbPassword.Password = string(hashedNewPassword)
	if err := service.DB.Save(&dbPassword).Error; err != nil {
		return err
	}

	// 移除旧的token
	if err := service.DB.Unscoped().Delete(&models.DBToken{}, "user_id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}

func (service *UserService) RetrievePassword(userID uint, code string, password string) error {
	var dbVerifyCode models.DBVerifyCode
	if err := service.DB.Where("user_id = ? AND code = ? AND expired_at > ?", userID, code, time.Now()).First(&dbVerifyCode).Error; err != nil {
		return fmt.Errorf("验证码错误")
	}

	if err := service.DB.Unscoped().Delete(&models.DBVerifyCode{}, "id = ?", dbVerifyCode.ID).Error; err != nil {
		return err
	}

	// 对新密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	dbPassword := models.DBPassword{UserID: userID, Password: string(hashedPassword)}
	if err := service.DB.Save(&dbPassword).Error; err != nil {
		return err
	}

	// 移除旧的token
	if err := service.DB.Unscoped().Delete(&models.DBToken{}, "user_id = ?", userID).Error; err != nil {
		return err
	}

	return nil
}
