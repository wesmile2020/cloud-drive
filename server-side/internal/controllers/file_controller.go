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
		UserID:   userID.(uint),
		Name:     request.Name,
		ParentID: *request.ParentID,
		Public:   *request.Public,
	}

	if err := controller.service.CreateDirectory(directory); err != nil {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create directory",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Code:    http.StatusOK,
		Message: "Directory created successfully",
		Data:    nil,
	})
}

func (controller *FileController) GetFiles(ctx *gin.Context) {
	var uid uint = 0
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

	userID, exists := ctx.Get("userID")
	if exists {
		uid = userID.(uint)
	}
	files := controller.service.GetFiles(uint(dirId), uid)

	tree := controller.service.GetFileTree(uint(dirId))

	ctx.JSON(http.StatusOK, &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    gin.H{"tree": tree, "files": files},
	})
}

func (controller *FileController) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/file")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("/directory", controller.CreateDirectory)
		authGroup.GET("/files/:directoryID", controller.GetFiles)
	}
}
