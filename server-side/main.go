package main

import (
	"cloud-drive/configs"
	"cloud-drive/internal/database"
	"cloud-drive/internal/routers"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// 加载配置并启动 Gin 服务器
func main() {
	// 加载配置
	cfg, err := configs.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 获取executable的路径
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
		return
	}
	// 获取executable的目录
	rootDir := filepath.Dir(executablePath)

	// 初始化数据库
	dsn := filepath.Join(rootDir, cfg.Database.DSN)
	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		return
	}

	// 初始化路由
	engine := routers.SetupRouter(db, rootDir)

	// 启动服务器
	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
