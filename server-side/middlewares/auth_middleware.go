package middlewares

import (
	"cloud-drive/internal/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 定义 JWT 的声明结构体
type Claims struct {
	UserID uint `json:"userID"`
	jwt.RegisteredClaims
}

const SECRET_KEY = "27c539cdba924c2e67a46c7c72847158fe24f5f8fd27c34cb789bbbe5a1168e1"

const Duration = 7 * 24 * time.Hour

// GenerateJWTToken 生成 JWT token
func GenerateJWTToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Duration)), // 设置 7 天过期
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY))
}

func ParseJWTToken(tokenString string, db *gorm.DB) (*Claims, error) {
	// 判断 token 是否在token表中
	var count int64 = 0
	db.Table("token").Where("token = ?", tokenString).Count(&count)
	if count == 0 {
		return nil, fmt.Errorf("token is out of work")
	}

	var claims Claims
	parsedToken, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return &claims, nil
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从 cookie 中获取 token
		tokenString, err := ctx.Cookie("token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, &models.Response{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Data:    nil,
			})
			return
		}

		// 解析 token
		claims, err := ParseJWTToken(tokenString, db)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, &models.Response{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		// 将用户信息存入上下文
		ctx.Set("userID", claims.UserID)
	}
}
