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

func Task(db *gorm.DB, pathUtil *utils.PathUtil) {

	logrus.Infof("Running scheduled task")

	// 删除过期的Token
	if err := db.Unscoped().Where("expired_at < ?", time.Now().Unix()).Delete(&models.DBToken{}).Error; err != nil {
		logrus.Errorf("Failed to delete expired files: %v", err)
	}

	// 删除7天前的临时上传文件
	var tempFiles []models.DBFileChunk
	if err := db.Where("updated_at < datetime('now', '-7 days')").Find(&tempFiles).Error; err != nil {
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

func Run(db *gorm.DB, pathUtil *utils.PathUtil) {
	Task(db, pathUtil)

	// 实现定时任务的逻辑
	cron := cron.New()
	// 每7天凌晨1点
	cron.AddFunc("0 1 * * 0", func() {
		Task(db, pathUtil)
	})

	cron.Start()
}
