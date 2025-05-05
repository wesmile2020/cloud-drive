package routers

import (
	"cloud-drive/internal/controllers"
	"cloud-drive/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, rootDir string) *gin.Engine {
	engine := gin.Default()

	group := engine.Group("/api")

	// 注册路由
	userService := services.NewUserService(db)
	userHandler := controllers.NewUserController(userService)
	userHandler.RegisterRoutes(group)

	// 注册路由
	fileService := services.NewFileService(db, rootDir)
	fileHandler := controllers.NewFileController(fileService)
	fileHandler.RegisterRoutes(group)

	return engine
}
