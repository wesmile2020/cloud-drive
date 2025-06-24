package routers

import (
	"cloud-drive/configs"
	"cloud-drive/internal/controllers"
	"cloud-drive/internal/services"
	"cloud-drive/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(engine *gin.Engine, db *gorm.DB, pathUtil *utils.PathUtil, config *configs.Config) {

	group := engine.Group("/api")

	mailService := services.NewMailService(&config.Email, db)
	userService := services.NewUserService(db)
	userHandler := controllers.NewUserController(userService, mailService)
	userHandler.RegisterRoutes(group)

	fileService := services.NewFileService(db, pathUtil)
	fileHandler := controllers.NewFileController(fileService)
	fileHandler.RegisterRoutes(group)

}
