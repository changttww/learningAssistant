package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	defaultPermissions = []string{"profile:view", "task:manage", "team:view"}
	defaultRoles       = []string{"student"}
	tokenTTLSeconds    = 7200
)

func registerAuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", handleRegister)
	router.POST("/login", handleLogin)
	router.POST("/logout", handleLogout)
}

type registerRequestBody struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

func handleRegister(c *gin.Context) {
	var req registerRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数格式不正确",
		})
		return
	}

	summary, err := registerUserAccount(req.Username, req.Email, req.Password, req.DisplayName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data": gin.H{
			"user": summary,
		},
	})
}

type loginRequestBody struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	Remember   bool   `json:"remember"`
}

func handleLogin(c *gin.Context) {
	var req loginRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数格式不正确",
		})
		return
	}

	req.Identifier = strings.TrimSpace(req.Identifier)
	req.Password = strings.TrimSpace(req.Password)

	if req.Identifier == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请输入用户名/邮箱和密码",
		})
		return
	}

	summary, err := authenticateUser(req.Identifier, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	token, refreshToken := generateTokens(summary.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token":           token,
			"refresh_token":   refreshToken,
			"expires_in":      tokenTTLSeconds,
			"user":            summary,
			"permissions":     defaultPermissions,
			"roles":           defaultRoles,
			"token_type":      "Bearer",
			"remember":        req.Remember,
			"issued_at":       time.Now().Unix(),
			"refresh_expires": time.Now().Add(24 * time.Hour).Unix(),
		},
	})
}

func handleLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "退出成功",
	})
}

type refreshRequestBody struct {
	RefreshToken string `json:"refresh_token"`
}

func generateTokens(userID uint64) (string, string) {
	now := time.Now().UnixNano()
	token := fmt.Sprintf("mock-token-%d-%d", userID, now)
	refreshToken := fmt.Sprintf("mock-refresh-%d-%d", userID, now)
	return token, refreshToken
}

func extractUserIDFromToken(token, prefix string) (uint64, error) {
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
