package controllers

import (
	"cloud-drive/internal/models"
	"cloud-drive/internal/services"
	"cloud-drive/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FileController struct {
	service *services.FileService
}

func NewFileController(service *services.FileService) *FileController {
	return &FileController{service: service}
}

func (controller *FileController) CreateDirectory(ctx *gin.Context) {
	// Implement the logic to create a directory
	var request models.CreateDirectoryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		var errorMessages []string
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validationError := range validationErrors {
				errorMessages = append(errorMessages, validationError.Error())
			}
		}
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Data:    errorMessages,
		})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user ID",
			Data:    nil,
		})
	}
	directory := &models.APIDirectory{
		UserID:     userID.(uint),
		Name:       request.Name,
		ParentID:   *request.ParentID,
		Permission: *request.Permission,
	}

	if err := controller.service.CreateDirectory(directory); err != nil {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	})
}

func (controller *FileController) GetFiles(ctx *gin.Context) {
	directoryID := ctx.Param("directoryID")
	if directoryID == "" {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "Directory ID is required",
			Data:    nil,
		})
		return
	}

	dirId, err := strconv.ParseUint(directoryID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid directory ID",
			Data:    nil,
		})
		return
	}

	var uid uint = 0
	tokenString, err := ctx.Cookie("token")
	if err == nil {
		claims, err := middleware.ParseJWTToken(tokenString, controller.service.DB)
		if err == nil {
			uid = claims.UserID
		}
	}

	files := controller.service.GetFiles(uint(dirId), uid)

	tree := controller.service.GetFileTree(uint(dirId), uid)

	ctx.JSON(http.StatusOK, &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    gin.H{"tree": tree, "files": files},
	})
}

func (controller *FileController) RegisterRoutes(router *gin.RouterGroup) {
	fileGroup := router.Group("/file")
	{
		fileGroup.GET("/:directoryID", controller.GetFiles)
	}

	authGroup := router.Group("/file", middleware.AuthMiddleware(controller.service.DB))
	{
		authGroup.POST("/directory", controller.CreateDirectory)
	}
}
