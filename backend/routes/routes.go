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
		// 用户相关路由
		users := v1.Group("/users")
		{
			users.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Users endpoint"})
			})
		}

		// 任务相关路由
		tasks := v1.Group("/tasks")
		{
			tasks.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Tasks endpoint"})
			})
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
}