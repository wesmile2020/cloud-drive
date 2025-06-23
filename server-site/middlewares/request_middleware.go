package middlewares

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 在别的中间件处理完成之后在执行这个中间件
		ctx.Next()
		// 判断是否处于abort状态，如果是则同步读取body，防止和前端的链接直接断开，导致前端出现ERR_CONNECTION_RESET错误
		if ctx.IsAborted() {
			_, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				logrus.Errorf("Failed to read body: %v", err)
			}
		}
	}
}
