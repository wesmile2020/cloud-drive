package main

import (
	"cloud-drive/configs"
	"cloud-drive/formats"
	"cloud-drive/internal/database"
	"cloud-drive/internal/routers"
	"cloud-drive/middlewares"
	"cloud-drive/schedules"
	"cloud-drive/utils"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 加载配置并启动 Gin 服务器
func main() {
	// 获取executable的路径
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
		return
	}
	// 获取executable的目录
	rootDir := filepath.Dir(executablePath)
	pathUtil := utils.NewPathUtil(rootDir)

	logger := &lumberjack.Logger{
		Filename:   filepath.Join(pathUtil.GetRootDir(), "logs/cloud-drive.log"), // 日志文件路径
		MaxSize:    10,                                                           // 每个日志文件的最大大小（MB）
		MaxBackups: 3,                                                            // 最多保留的旧日志文件数量
	}

	logrus.SetOutput(io.MultiWriter(
		os.Stdout,
		logger,
	))
	logrus.SetFormatter(&formats.LogFormatter{})

	// 加载配置
	cfg, err := configs.LoadConfig(pathUtil)
	if err != nil {
		logrus.Errorf("Failed to load config: %v", err)
		return
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	dsn := filepath.Join(rootDir, cfg.Database.DSN)
	db, err := database.InitDB(dsn)
	if err != nil {
		logrus.Errorf("Failed to initialize database: %v", err)
		return
	}

	schedules.Run(db)

	// 初始化路由
	engine := gin.New()
	engine.Use(middlewares.LogMiddleware())
	engine.Use(gin.Recovery())
	engine.Use(middlewares.RequestMiddleWare())

	// 托管静态文件
	engine.Static("/static", filepath.Join(rootDir, "static"))

	// 注册路由
	routers.SetupRouter(engine, db, pathUtil)

	// 启动服务器
	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		logrus.Errorf("Failed to start server: %v", err)
	}
}
