package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Learning Assistant API is running",
		})
	})

	// API 版本分组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := v1.Group("/auth")
		registerAuthRoutes(auth)

		// 用户相关路由
		users := v1.Group("/users")
		registerUserRoutes(users)

		// 任务相关路由（团队任务）
		tasks := v1.Group("/tasks")
		{
			// 注册团队任务路由处理器（在 routes/teamtasks.go 中实现）
			registerTeamTaskRoutes(tasks)
		}

		// 团队相关路由
		teams := v1.Group("/teams")
		{
			teams.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Teams endpoint"})
			})
		}

		// 学习室相关路由
		studyRooms := v1.Group("/study-rooms")
		{
			studyRooms.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Study rooms endpoint"})
			})
		}
	}

	// 兼容旧版未带版本号的前缀 /api/**
	legacy := r.Group("/api")
	{
		authLegacy := legacy.Group("/auth")
		registerAuthRoutes(authLegacy)

		usersLegacy := legacy.Group("/users")
		registerUserRoutes(usersLegacy)
	}
}
