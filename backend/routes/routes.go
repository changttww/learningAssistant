package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 初始化RAG服务（确保在任务路由使用前初始化）
	initRAGServices()

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
		registerTaskStatRoutes(tasks)

		// 团队相关路由
		teams := v1.Group("/teams")
		registerTeamRoutes(teams)

		// 通知相关路由
		notifications := v1.Group("/notifications")
		registerNotificationRoutes(notifications)

		// 学习室相关路由
		study := v1.Group("/study")
		registerStudyRoutes(study)
		registerStudyWebsocketRoutes(study)
		registerStudySessionRoutes(study)
		registerStudyNotesRoutes(study)
		{
			rooms := study.Group("/rooms")
			rooms.GET("/:roomId/chat/history", handleGetRoomChatHistory)
			rooms.POST("/:roomId/chat", handlePostRoomChat)
		}

		analysis := v1.Group("/analysis")
		registerAnalysisRoutes(analysis)

		// AI相关路由
		ai := v1.Group("/ai")
		registerAIRoutes(ai)

		// 智能笔记增强路由
		notes := v1.Group("/notes")
		registerNoteEnhanceRoutes(notes)

		// 知识库相关路由
		knowledge := v1.Group("")
		registerKnowledgeBaseRoutes(knowledge)
		registerKnowledgeSyncRoutes(knowledge)
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
		registerTaskStatRoutes(tasksLegacy)

		teamsLegacy := legacy.Group("/teams")
		registerTeamRoutes(teamsLegacy)

		notificationsLegacy := legacy.Group("/notifications")
		registerNotificationRoutes(notificationsLegacy)

		studyLegacy := legacy.Group("/study")
		registerStudyRoutes(studyLegacy)
		registerStudyWebsocketRoutes(studyLegacy)
		registerStudySessionRoutes(studyLegacy)
		registerStudyNotesRoutes(studyLegacy)
		{
			roomsLegacy := studyLegacy.Group("/rooms")
			roomsLegacy.GET("/:roomId/chat/history", handleGetRoomChatHistory)
			roomsLegacy.POST("/:roomId/chat", handlePostRoomChat)
		}

		analysisLegacy := legacy.Group("/analysis")
		registerAnalysisRoutes(analysisLegacy)

		// AI相关路由
		aiLegacy := legacy.Group("/ai")
		registerAIRoutes(aiLegacy)

		// 智能笔记增强路由
		notesLegacy := legacy.Group("/notes")
		registerNoteEnhanceRoutes(notesLegacy)

		// 知识库相关路由
		knowledgeLegacy := legacy.Group("")
		registerKnowledgeBaseRoutes(knowledgeLegacy)
		registerKnowledgeSyncRoutes(knowledgeLegacy)
	}
}
