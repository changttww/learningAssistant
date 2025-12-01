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

		// 任务相关路由
		tasks := v1.Group("/tasks")
		registerTaskRoutes(tasks)

		// 团队相关路由
		teams := v1.Group("/teams")
		{
			teams.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Teams endpoint"})
			})
		}

		// 学习室相关路由
		study := v1.Group("/study")
		registerStudyRoutes(study)
		registerStudyWebsocketRoutes(study)
		registerStudyNotesRoutes(study)
		{
			rooms := study.Group("/rooms")
			rooms.GET("/:roomId/chat/history", handleGetRoomChatHistory)
			rooms.POST("/:roomId/chat", handlePostRoomChat)
		}
	}

	// 兼容旧版未带版本号的前缀 /api/**
	legacy := r.Group("/api")
	{
		authLegacy := legacy.Group("/auth")
		registerAuthRoutes(authLegacy)

		usersLegacy := legacy.Group("/users")
		registerUserRoutes(usersLegacy)

		tasksLegacy := legacy.Group("/tasks")
		registerTaskRoutes(tasksLegacy)

		studyLegacy := legacy.Group("/study")
		registerStudyRoutes(studyLegacy)
		registerStudyWebsocketRoutes(studyLegacy)
		registerStudyNotesRoutes(studyLegacy)
		{
			roomsLegacy := studyLegacy.Group("/rooms")
			roomsLegacy.GET("/:roomId/chat/history", handleGetRoomChatHistory)
			roomsLegacy.POST("/:roomId/chat", handlePostRoomChat)
		}
	}
}
