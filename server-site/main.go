package main

import (
	"cloud-drive/configs"
	"cloud-drive/formats"
	"cloud-drive/internal/database"
	"cloud-drive/internal/routers"
	"cloud-drive/middlewares"
	"cloud-drive/schedules"
	"cloud-drive/utils"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

//go:embed static/*
var staticFS embed.FS

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

	// 加载配置
	cfg, err := configs.LoadConfig(pathUtil)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

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

	if cfg.Log.Level == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else if cfg.Log.Level == "info" {
		logrus.SetLevel(logrus.InfoLevel)
	} else if cfg.Log.Level == "warn" {
		logrus.SetLevel(logrus.WarnLevel)
	} else if cfg.Log.Level == "error" {
		logrus.SetLevel(logrus.ErrorLevel)
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

	schedules.Run(db, pathUtil)

	// 初始化路由
	engine := gin.New()
	engine.Use(middlewares.LogMiddleware())
	engine.Use(gin.Recovery())
	engine.Use(middlewares.RequestMiddleWare())

	// 注册路由
	routers.SetupRouter(engine, db, pathUtil, cfg)

	// 托管静态文件
	staticFP, _ := fs.Sub(staticFS, "static")
	httpFS := http.FS(staticFP)
	// engine.StaticFS("/static", http.FS(staticFP))

	engine.NoRoute(gin.WrapH(http.FileServer(httpFS)))

	// 启动服务器
	logrus.Infof("Server is running on http://localhost:%s", cfg.Server.Port)
	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		logrus.Errorf("Failed to start server: %v", err)
	}
}
