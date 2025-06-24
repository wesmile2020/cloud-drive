package services

import (
	"cloud-drive/configs"
	"cloud-drive/internal/models"
	"crypto/tls"
	"math/rand"
	"net/smtp"
	"time"

	"gorm.io/gorm"
)

type MailService struct {
	config *configs.EmailConfig
	DB     *gorm.DB
}

func NewMailService(config *configs.EmailConfig, db *gorm.DB) *MailService {
	return &MailService{
		config: config,
		DB:     db,
	}
}

func (service *MailService) SendEmail(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", service.config.Username, service.config.Password, service.config.Host)
	// 发送邮件
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         service.config.Host,
	}
	connection, err := tls.Dial("tcp", service.config.Host+":"+service.config.Port, tlsConfig)
	if err != nil {
		return err
	}
	defer connection.Close()
	client, err := smtp.NewClient(connection, service.config.Host)
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Auth(auth)
	if err != nil {
		return err
	}
	err = client.Mail(service.config.Username)
	if err != nil {
		return err
	}
	err = client.Rcpt(to)
	if err != nil {
		return err
	}
	wc, err := client.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	message := []byte(
		"From: cloud_drive <" + service.config.Username + ">\n" +
			"To: " + to + "\n" +
			"Subject: " + subject + "\n" +
			"\n" +
			body + "\n",
	)
	_, err = wc.Write(message)

	if err != nil {
		return err
	}
	return nil
}

func (service *MailService) SendVerifyCode(userID uint, to string) error {
	// 生成6位验证码
	codeLength := 6
	code := make([]byte, codeLength)

	for i := range code {
		code[i] = byte(rand.Intn(10)) + '0'
	}

	strCode := string(code)
	// 先删除旧的验证码
	if err := service.DB.Unscoped().Where("user_id = ?", userID).Delete(&models.DBVerifyCode{}).Error; err != nil {
		return err
	}
	// 保存验证码到数据库
	verifyCode := &models.DBVerifyCode{
		UserID:    userID,
		Code:      strCode,
		ExpiredAt: time.Now().Add(time.Minute * 5),
	}
	if err := service.DB.Create(verifyCode).Error; err != nil {
		return err
	}
	subject := "验证码"
	body := "请不要将验证码告诉他人！\n您的验证码为：" + strCode + "，有效期为5分钟。"

	return service.SendEmail(to, subject, body)
}
