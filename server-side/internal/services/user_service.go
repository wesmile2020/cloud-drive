package services

import (
	"cloud-drive/internal/models"
	"errors"

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
func (service *UserService) RegisterUser(serviceUser *models.ServiceUser, password string) error {
	// 检查phone是否已存在
	var existingUser models.DBUser
	if err := service.DB.Where("phone = ?", serviceUser.Phone).First(&existingUser).Error; err == nil {
		return errors.New("手机号已被注册")
	}

	// 检查email是否已存在
	if err := service.DB.Where("email = ?", serviceUser.Email).First(&existingUser).Error; err == nil {
		return errors.New("邮箱已被注册")
	}

	// 插入用户表
	dbUser := serviceUser.ToDBUser()
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
func (service *UserService) LoginUser(account, password string) (*models.ServiceUser, error) {
	// 根据 account 查询用户，account 可能是 phone 或 email
	var existingUser models.DBUser
	if err := service.DB.Where("phone = ? OR email = ?", account, account).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("账号或者密码错误")
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
		return nil, errors.New("账号或者密码错误")
	}

	serviceUser := existingUser.ToServiceUser()

	return &serviceUser, nil
}

// 获取用户登录信息
func (service *UserService) GetUserInfo(userID uint) (*models.ServiceUser, error) {
	var dbUser models.DBUser
	if err := service.DB.First(&dbUser, userID).Error; err != nil {
		return nil, err
	}

	serviceUser := dbUser.ToServiceUser()

	return &serviceUser, nil
}
