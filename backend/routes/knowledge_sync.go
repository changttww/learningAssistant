package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"learningAssistant-backend/database"
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/models"
)

// registerKnowledgeSyncRoutes 注册知识同步路由
func registerKnowledgeSyncRoutes(router *gin.RouterGroup) {
	sync := router.Group("/knowledge-base")
	sync.Use(middleware.AuthMiddleware())

	// 同步用户的所有已完成任务到知识库
	sync.POST("/sync-tasks", syncTasksToKnowledge)
	// 同步用户的所有笔记到知识库
	sync.POST("/sync-notes", syncNotesToKnowledge)
	// 一键同步所有内容
	sync.POST("/sync-all", syncAllToKnowledge)
	// 按团队同步团队任务到知识库
	sync.POST("/team/sync", syncTeamKnowledge)
	// 批量发布所有草稿状态的知识条目
	sync.POST("/publish-all", publishAllDraftEntries)
}

// syncTasksToKnowledge 将用户任务同步到知识库
func syncTasksToKnowledge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	uid := userID.(uint64)

	if ragService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "RAG服务未初始化"})
		return
	}

	// 异步处理同步任务
	go func(targetUserID uint64) {
		db := database.GetDB()

		// 获取用户所有任务（不仅仅是已完成的）
		var tasks []models.Task
		if err := db.Where("created_by = ? OR owner_user_id = ?", targetUserID, targetUserID).
			Find(&tasks).Error; err != nil {
			// 异步任务中的错误只能打印日志
			return
		}

		for _, task := range tasks {
			// 跳过标题和描述都为空的任务
			if task.Title == "" && task.Description == "" {
				continue
			}

			// 使用标题作为内容（如果描述为空）
			content := task.Description
			if content == "" {
				content = task.Title
			}

			// 忽略错误，继续处理下一个
			_, _ = ragService.AddDocument(targetUserID, 1, task.ID, task.Title, content)
		}
	}(uid)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"status": "processing",
		},
		"msg": "任务同步已在后台开始",
	})
}

// syncNotesToKnowledge 将用户笔记同步到知识库
func syncNotesToKnowledge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	uid := userID.(uint64)

	if ragService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "RAG服务未初始化"})
		return
	}

	// 异步处理同步任务
	go func(targetUserID uint64) {
		db := database.GetDB()

		// 获取用户所有笔记
		var notes []models.StudyNote
		if err := db.Where("user_id = ?", targetUserID).Find(&notes).Error; err != nil {
			return
		}

		for _, note := range notes {
			// 跳过标题和内容都为空的笔记
			if note.Title == "" && note.Content == "" {
				continue
			}

			// 使用标题作为内容（如果内容为空）
			content := note.Content
			if content == "" {
				content = note.Title
			}

			_, _ = ragService.AddDocument(targetUserID, 2, note.ID, note.Title, content)
		}
	}(uid)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"status": "processing",
		},
		"msg": "笔记同步已在后台开始",
	})
}

// syncAllToKnowledge 同步所有任务和笔记到知识库
func syncAllToKnowledge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	uid := userID.(uint64)

	if ragService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "RAG服务未初始化"})
		return
	}

	// 异步处理同步任务
	go func(targetUserID uint64) {
		db := database.GetDB()

		// 同步任务 - 同步所有任务（不仅仅是已完成的）
		var tasks []models.Task
		db.Where("created_by = ? OR owner_user_id = ?", targetUserID, targetUserID).Find(&tasks)

		for _, task := range tasks {
			// 如果标题和描述都为空，跳过
			if task.Title == "" && task.Description == "" {
				continue
			}
			// 使用标题作为内容（如果描述为空）
			content := task.Description
			if content == "" {
				content = task.Title
			}
			_, _ = ragService.AddDocument(targetUserID, 1, task.ID, task.Title, content)
		}

		// 同步笔记
		var notes []models.StudyNote
		db.Where("user_id = ?", targetUserID).Find(&notes)

		for _, note := range notes {
			// 如果标题和内容都为空，跳过
			if note.Title == "" && note.Content == "" {
				continue
			}
			// 使用标题作为内容（如果内容为空）
			content := note.Content
			if content == "" {
				content = note.Title
			}
			_, _ = ragService.AddDocument(targetUserID, 2, note.ID, note.Title, content)
		}
	}(uid)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"status": "processing",
		},
		"msg": "全部内容同步已在后台开始",
	})
}

// syncTeamKnowledge 将指定团队的任务同步到知识库（写入 team_id）
func syncTeamKnowledge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	uid := userID.(uint64)

	if ragService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "RAG服务未初始化"})
		return
	}

	var req struct {
		TeamID uint64 `json:"team_id"`
	}

	// 允许 body 或 query，body 解析失败不直接返回，继续尝试 query
	_ = c.ShouldBindJSON(&req)
	if req.TeamID == 0 {
		if teamIDStr := c.Query("team_id"); teamIDStr != "" {
			if parsed, errParse := strconv.ParseUint(teamIDStr, 10, 64); errParse == nil {
				req.TeamID = parsed
			}
		}
	}

	if req.TeamID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供有效的 team_id"})
		return
	}

	db := database.GetDB()

	// 校验团队存在且用户属于该团队（owner 或成员）
	var allowed int64
	db.Model(&models.Team{}).
		Where("id = ? AND owner_user_id = ?", req.TeamID, uid).
		Or("id = ? AND id IN (SELECT team_id FROM team_members WHERE team_id = ? AND user_id = ?)", req.TeamID, req.TeamID, uid).
		Count(&allowed)
	if allowed == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权同步该团队"})
		return
	}

	// 异步处理同步任务
	go func(targetUserID uint64, targetTeamID uint64) {
		db := database.GetDB()

		// 拉取团队任务
		var tasks []models.Task
		if err := db.Where("owner_team_id = ?", targetTeamID).Find(&tasks).Error; err != nil {
			return
		}

		for _, task := range tasks {
			if task.Title == "" && task.Description == "" {
				continue
			}

			content := task.Description
			if content == "" {
				content = task.Title
			}

			entry, err := ragService.AddDocument(targetUserID, 1, task.ID, task.Title, content)
			if err != nil {
				continue
			}

			// 绑定 team_id，确保团队列表可见
			// 忽略错误
			_ = db.Model(entry).Update("team_id", targetTeamID).Error
		}
	}(uid, req.TeamID)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"status": "processing",
		},
		"msg": "团队任务同步已在后台开始",
	})
}

// publishAllDraftEntries 批量发布所有草稿状态的知识条目
func publishAllDraftEntries(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	db := database.GetDB()

	// 统计草稿数量
	var draftCount int64
	db.Model(&models.KnowledgeBaseEntry{}).
		Where("user_id = ? AND status = 0", userID.(uint64)).
		Count(&draftCount)

	// 批量更新 status 从 0 改为 1
	result := db.Model(&models.KnowledgeBaseEntry{}).
		Where("user_id = ? AND status = 0", userID.(uint64)).
		Update("status", 1)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"published_count": result.RowsAffected,
			"original_drafts": draftCount,
		},
		"msg": "批量发布成功",
	})
}
