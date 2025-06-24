package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		latency := time.Since(start)

		statusCode := ctx.Writer.Status()

		fields := logrus.Fields{
			"status":  statusCode,
			"method":  ctx.Request.Method,
			"path":    ctx.Request.URL.Path,
			"ip":      ctx.ClientIP(),
			"latency": latency,
		}
		if statusCode >= http.StatusInternalServerError {
			logrus.WithFields(fields).Error(ctx.Errors.String())
		} else if statusCode >= http.StatusBadRequest {
			logrus.WithFields(fields).Warn(ctx.Errors.String())
		} else {
			logrus.WithFields(fields).Debug()
		}
	}
}
