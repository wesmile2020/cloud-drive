package routers

import (
	"cloud-drive/internal/controllers"
	"cloud-drive/internal/services"
	"cloud-drive/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(engine *gin.Engine, db *gorm.DB, pathUtil *utils.PathUtil) {
	group := engine.Group("/api")

	// 注册路由
	userService := services.NewUserService(db)
	userHandler := controllers.NewUserController(userService)
	userHandler.RegisterRoutes(group)

	// 注册路由
	fileService := services.NewFileService(db, pathUtil)
	fileHandler := controllers.NewFileController(fileService)
	fileHandler.RegisterRoutes(group)

}
