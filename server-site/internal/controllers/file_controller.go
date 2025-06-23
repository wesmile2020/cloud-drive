package controllers

import (
	"cloud-drive/internal/models"
	"cloud-drive/internal/services"
	"cloud-drive/middlewares"
	"fmt"
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

func (controller *FileController) UpdateDirectory(ctx *gin.Context) {
	// Implement the logic to update a directory
	var request models.UpdateDirectoryRequest
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
	}
	directoryID := ctx.Param("id")
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
		Permission: *request.Permission,
	}

	if err := controller.service.UpdateDirectory(uint(dirId), directory); err != nil {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	})
}

func (controller *FileController) DeleteDirectory(ctx *gin.Context) {
	// Implement the logic to delete a directory
	directoryID := ctx.Param("id")
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
	if !exists {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user ID",
			Data:    nil,
		})
		return
	}

	if err := controller.service.DeleteDirectory(uint(dirId), userID.(uint)); err != nil {
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
		claims, err := middlewares.ParseJWTToken(tokenString, controller.service.DB)
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

func (controller *FileController) UploadFile(ctx *gin.Context) {
	// Implement the logic to upload a file
	userId, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user ID",
			Data:    nil,
		})
		return
	}
	uid := userId.(uint)

	var request models.UploadFileRequest
	if err := ctx.ShouldBind(&request); err != nil {
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

	fileId, err := controller.service.UploadFile(&request, uid)
	if err != nil {
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
		Data:    gin.H{"fileId": fileId},
	})
}

func (controller *FileController) DeleteFile(ctx *gin.Context) {
	// Implement the logic to download a file
	userId, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user ID",
			Data:    nil,
		})
	}
	uid := userId.(uint)

	fileId := ctx.Param("id")
	if fileId == "" {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "File ID is required",
			Data:    nil,
		})
	}

	id, err := strconv.ParseUint(fileId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid file ID",
			Data:    nil,
		})
		return
	}

	if err := controller.service.DeleteFile(uint(id), uid); err != nil {
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

func (controller *FileController) UpdateFile(ctx *gin.Context) {
	userId, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user ID",
			Data:    nil,
		})
	}
	uid := userId.(uint)

	fileId := ctx.Param("id")
	if fileId == "" {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "File ID is required",
			Data:    nil,
		})
	}

	id, err := strconv.ParseUint(fileId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid file ID",
			Data:    nil,
		})
	}

	var request models.UpdateFileRequest
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
	}

	if err := controller.service.UpdateFile(uint(id), uid, &request); err != nil {
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

func (controller *FileController) DownloadFile(ctx *gin.Context) {
	// Implement the logic to download a file
	fileId := ctx.Param("id")
	if fileId == "" {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("file ID is required"))
		return
	}

	id, err := strconv.ParseUint(fileId, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var uid uint = 0
	userId, exists := ctx.Get("userID")
	if exists {
		uid = userId.(uint)
	}
	fileUrl, err := controller.service.DownloadFile(uint(id), uid)
	if err != nil {
		ctx.AbortWithError(http.StatusOK, err)
		return
	}
	ctx.File(fileUrl)
}

func (controller *FileController) RegisterRoutes(router *gin.RouterGroup) {
	fileGroup := router.Group("/file")
	{
		fileGroup.GET("/:directoryID", controller.GetFiles)
		fileGroup.GET("/download/:id", controller.DownloadFile)
	}

	authGroup := router.Group("/file", middlewares.AuthMiddleware(controller.service.DB))
	{
		authGroup.POST("/directory", controller.CreateDirectory)
		authGroup.PUT("/directory/:id", controller.UpdateDirectory)
		authGroup.DELETE("/directory/:id", controller.DeleteDirectory)
		authGroup.POST("/upload", controller.UploadFile)
		authGroup.DELETE("/:id", controller.DeleteFile)
		authGroup.PUT("/:id", controller.UpdateFile)
	}
}
