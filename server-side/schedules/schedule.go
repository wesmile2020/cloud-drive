package schedules

import (
	"cloud-drive/internal/models"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Task(db *gorm.DB) {
	logrus.Infof("Running scheduled task")
	// 执行清理逻辑，例如删除过期的文件
	if err := db.Unscoped().Where("expired_at < ?", time.Now().Unix()).Delete(&models.DBToken{}).Error; err != nil {
		logrus.Errorf("Failed to delete expired files: %v", err)
	}
}

func Run(db *gorm.DB) {
	Task(db)

	// 实现定时任务的逻辑
	cron := cron.New()
	// 每7天凌晨1点
	cron.AddFunc("0 1 * * 0", func() {
		Task(db)
	})

	cron.Start()
}
