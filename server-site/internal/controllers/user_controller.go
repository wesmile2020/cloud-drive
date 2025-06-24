package controllers

import (
	"cloud-drive/internal/models"
	"cloud-drive/internal/services"
	"cloud-drive/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	userService *services.UserService
	mailService *services.MailService
}

func NewUserController(userService *services.UserService, mailService *services.MailService) *UserController {

	return &UserController{
		userService: userService,
		mailService: mailService,
	}
}

// RegisterUser 处理用户注册请求
func (controller *UserController) RegisterUser(ctx *gin.Context) {
	var request models.RegisterUserRequest
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

	apiUser := &models.APIUser{
		Name:  request.Name,
		Email: request.Email,
		Phone: request.Phone,
	}

	if err := controller.userService.RegisterUser(apiUser, request.Password); err != nil {
		response := models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := models.Response{
		Code:    http.StatusOK,
		Message: "注册成功",
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *UserController) LoginUser(ctx *gin.Context) {
	var request models.LoginUserRequest
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

	apiUser, err := controller.userService.LoginUser(request.Account, request.Password)
	if err != nil {
		response := models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}
	tokenString, err := middlewares.GenerateJWTToken(apiUser.ID)
	if err != nil {
		response := models.Response{
			Code:    http.StatusInternalServerError,
			Message: "generate token failed",
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}

	// 将token存入数据库中
	if err := controller.userService.SaveToken(tokenString, apiUser.ID, time.Now().Add(middlewares.Duration)); err != nil {
		response := models.Response{
			Code:    http.StatusInternalServerError,
			Message: "save token failed",
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
	}

	// token 有效期是一周，设置cookie
	ctx.SetCookie("token", tokenString, int(middlewares.Duration*3600), "/", "", false, true)

	response := models.Response{
		Code:    http.StatusOK,
		Message: "登录成功",
		Data:    apiUser,
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *UserController) Logout(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		response := models.Response{
			Code:    http.StatusInternalServerError,
			Message: "get token failed",
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}

	if err := controller.userService.Logout(tokenString); err != nil {
		response := models.Response{
			Code:    http.StatusInternalServerError,
			Message: "logout failed",
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.SetCookie("token", "", -1, "/", "", false, true)

	response := models.Response{
		Code:    http.StatusOK,
		Message: "退出登录成功",
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *UserController) GetUserInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		response := models.Response{
			Code:    http.StatusUnauthorized,
			Message: "未登录",
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}

	user, err := controller.userService.GetUserInfo(userID.(uint))
	if err != nil {
		response := models.Response{
			Code:    http.StatusInternalServerError,
			Message: "用户不存在",
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    user,
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *UserController) EditUserInfo(ctx *gin.Context) {
	var request models.EditUserInfoRequest
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

	apiUser := &models.APIUser{
		Name:  request.Name,
		Email: request.Email,
		Phone: request.Phone,
	}
	userID, exists := ctx.Get("userID")
	if !exists {
		response := models.Response{
			Code:    http.StatusUnauthorized,
			Message: "未登录",
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}
	if err := controller.userService.EditUserInfo(userID.(uint), apiUser); err != nil {
		response := models.Response{
			Code:    http.StatusInternalServerError,
			Message: "编辑用户信息失败",
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := models.Response{
		Code:    http.StatusOK,
		Message: "编辑用户信息成功",
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *UserController) UpdatePassword(ctx *gin.Context) {
	var request models.UpdatePasswordRequest
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
		response := models.Response{
			Code:    http.StatusUnauthorized,
			Message: "未登录",
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}
	if err := controller.userService.EditPassword(userID.(uint), request.OldPassword, request.NewPassword); err != nil {
		response := models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := models.Response{
		Code:    http.StatusOK,
		Message: "修改密码成功",
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, response)
}

// 获取验证码
func (controller *UserController) GetVerifyCode(ctx *gin.Context) {
	var request models.GetVerifyCodeRequest
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
		response := models.Response{
			Code:    http.StatusUnauthorized,
			Message: "未登录",
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}
	// 发送验证码
	if err := controller.mailService.SendVerifyCode(userID.(uint), request.Email); err != nil {
		response := models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := models.Response{
		Code:    http.StatusOK,
		Message: "验证码发送成功",
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *UserController) RetrievePassword(ctx *gin.Context) {
	var request models.RetrievePasswordRequest
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
		response := models.Response{
			Code:    http.StatusUnauthorized,
			Message: "未登录",
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}
	if err := controller.userService.RetrievePassword(userID.(uint), request.Code, request.Password); err != nil {
		response := models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := models.Response{
		Code:    http.StatusOK,
		Message: "重置密码成功",
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *UserController) RegisterRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/user")
	{
		userGroup.POST("/register", controller.RegisterUser)
		userGroup.POST("/login", controller.LoginUser)
		userGroup.POST("/logout", controller.Logout)
	}
	// 验证token的中间件
	authGroup := router.Group("/user", middlewares.AuthMiddleware(controller.userService.DB))
	{
		authGroup.GET("/info", controller.GetUserInfo)
		authGroup.PUT("/info", controller.EditUserInfo)
		authGroup.PUT("/password", controller.UpdatePassword)
		authGroup.POST("/verify_code", controller.GetVerifyCode)
		authGroup.POST("/retrieve_password", controller.RetrievePassword)
	}
}
