package services

import (
	"cloud-drive/configs"
	"crypto/tls"
	"net/smtp"
)

type MailService struct {
	config *configs.EmailConfig
}

func NewMailService(config *configs.EmailConfig) *MailService {
	return &MailService{
		config: config,
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

	message := []byte("From: cloud_drive <" + service.config.Username + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	_, err = wc.Write(message)

	if err != nil {
		return err
	}
	return nil
}
