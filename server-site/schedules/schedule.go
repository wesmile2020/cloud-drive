package schedules

import (
	"cloud-drive/internal/models"
	"cloud-drive/utils"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func DayTask(db *gorm.DB, pathUtil *utils.PathUtil) {

	logrus.Infof("Running DayTask task")

	// 删除过期的Token
	if err := db.Unscoped().Where("expired_at < ?", time.Now()).Delete(&models.DBToken{}).Error; err != nil {
		logrus.Errorf("Failed to delete expired files: %v", err)
	}

	// 删除7天前的临时上传文件
	var tempFiles []models.DBFileChunk
	if err := db.Where("expired_at < ?", time.Now()).Find(&tempFiles).Error; err != nil {
		logrus.Errorf("Failed to delete expired files: %v", err)
	}
	for _, file := range tempFiles {
		fileUrl := filepath.Join(pathUtil.GetTempDir(), file.FileID)
		if err := utils.RemoveFile(fileUrl); err != nil {
			logrus.Errorf("Failed to delete expired files: %v", err)
		}
		if err := db.Unscoped().Where("file_id = ?", file.FileID).Delete(&models.DBFileChunk{}).Error; err != nil {
			logrus.Errorf("Failed to delete expired files: %v", err)
		}
	}
}

func MinuteTask(db *gorm.DB) {
	logrus.Infof("Running MinuteTask task")

	// 删除过期验证码
	if err := db.Unscoped().Where("expired_at < ?", time.Now()).Delete(&models.DBVerifyCode{}).Error; err != nil {
		logrus.Errorf("Failed to delete expired verify code: %v", err)
	}
}

func Run(db *gorm.DB, pathUtil *utils.PathUtil) {
	DayTask(db, pathUtil)
	MinuteTask(db)

	// 实现定时任务的逻辑
	cron := cron.New()
	// 每天凌晨1点
	cron.AddFunc("0 0 1 * *", func() {
		DayTask(db, pathUtil)
	})

	// 每5分钟执行一次
	cron.AddFunc("*/5 * * * *", func() {
		MinuteTask(db)
	})

	cron.Start()
}
