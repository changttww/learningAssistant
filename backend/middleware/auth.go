package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "缺少授权信息",
			})
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "授权信息格式不正确",
			})
			c.Abort()
			return
		}

		userID, err := extractUserIDFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "访问令牌已失效",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

// extractUserIDFromToken 从token中提取用户ID
func extractUserIDFromToken(token string) (uint64, error) {
	prefix := "mock-token-"
	if !strings.HasPrefix(token, prefix) {
		return 0, fmt.Errorf("invalid token prefix")
	}

	trimmed := strings.TrimPrefix(token, prefix)
	parts := strings.Split(trimmed, "-")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid token format")
	}

	return strconv.ParseUint(parts[0], 10, 64)
}
