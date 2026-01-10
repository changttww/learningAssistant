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

	if ragService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "RAG服务未初始化"})
		return
	}

	db := database.GetDB()

	// 获取用户所有任务（不仅仅是已完成的）
	var tasks []models.Task
	if err := db.Where("created_by = ? OR owner_user_id = ?", userID.(uint64), userID.(uint64)).
		Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务失败"})
		return
	}

	syncedCount := 0
	skippedCount := 0

	for _, task := range tasks {
		// 跳过标题和描述都为空的任务
		if task.Title == "" && task.Description == "" {
			skippedCount++
			continue
		}

		// 使用标题作为内容（如果描述为空）
		content := task.Description
		if content == "" {
			content = task.Title
		}

		_, err := ragService.AddDocument(userID.(uint64), 1, task.ID, task.Title, content)
		if err != nil {
			skippedCount++
		} else {
			syncedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"synced_count":  syncedCount,
			"skipped_count": skippedCount,
			"total_tasks":   len(tasks),
		},
		"msg": "任务同步完成",
	})
}

// syncNotesToKnowledge 将用户笔记同步到知识库
func syncNotesToKnowledge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	if ragService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "RAG服务未初始化"})
		return
	}

	db := database.GetDB()

	// 获取用户所有笔记
	var notes []models.StudyNote
	if err := db.Where("user_id = ?", userID.(uint64)).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记失败"})
		return
	}

	syncedCount := 0
	skippedCount := 0

	for _, note := range notes {
		// 跳过标题和内容都为空的笔记
		if note.Title == "" && note.Content == "" {
			skippedCount++
			continue
		}

		// 使用标题作为内容（如果内容为空）
		content := note.Content
		if content == "" {
			content = note.Title
		}

		_, err := ragService.AddDocument(userID.(uint64), 2, note.ID, note.Title, content)
		if err != nil {
			skippedCount++
		} else {
			syncedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"synced_count":  syncedCount,
			"skipped_count": skippedCount,
			"total_notes":   len(notes),
		},
		"msg": "笔记同步完成",
	})
}

// syncAllToKnowledge 同步所有任务和笔记到知识库
func syncAllToKnowledge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	if ragService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "RAG服务未初始化"})
		return
	}

	db := database.GetDB()

	// 同步任务 - 同步所有任务（不仅仅是已完成的）
	var tasks []models.Task
	db.Where("created_by = ? OR owner_user_id = ?", userID.(uint64), userID.(uint64)).Find(&tasks)

	taskSyncedCount := 0
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
		if _, err := ragService.AddDocument(userID.(uint64), 1, task.ID, task.Title, content); err == nil {
			taskSyncedCount++
		}
	}

	// 同步笔记
	var notes []models.StudyNote
	db.Where("user_id = ?", userID.(uint64)).Find(&notes)

	noteSyncedCount := 0
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
		if _, err := ragService.AddDocument(userID.(uint64), 2, note.ID, note.Title, content); err == nil {
			noteSyncedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"tasks_synced": taskSyncedCount,
			"notes_synced": noteSyncedCount,
			"total_synced": taskSyncedCount + noteSyncedCount,
		},
		"msg": "全部内容同步完成",
	})
}

// syncTeamKnowledge 将指定团队的任务同步到知识库（写入 team_id）
func syncTeamKnowledge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

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
		Where("id = ? AND owner_user_id = ?", req.TeamID, userID.(uint64)).
		Or("id = ? AND id IN (SELECT team_id FROM team_members WHERE team_id = ? AND user_id = ?)", req.TeamID, req.TeamID, userID.(uint64)).
		Count(&allowed)
	if allowed == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权同步该团队"})
		return
	}

	// 拉取团队任务
	var tasks []models.Task
	if err := db.Where("owner_team_id = ?", req.TeamID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取团队任务失败"})
		return
	}

	syncedCount := 0
	skippedCount := 0

	for _, task := range tasks {
		if task.Title == "" && task.Description == "" {
			skippedCount++
			continue
		}

		content := task.Description
		if content == "" {
			content = task.Title
		}

		entry, err := ragService.AddDocument(userID.(uint64), 1, task.ID, task.Title, content)
		if err != nil {
			skippedCount++
			continue
		}

		// 绑定 team_id，确保团队列表可见
		if err := db.Model(entry).Update("team_id", req.TeamID).Error; err != nil {
			// 更新失败不阻塞其他任务
			continue
		}

		syncedCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"tasks_synced":  syncedCount,
			"skipped_count": skippedCount,
			"total_tasks":   len(tasks),
		},
		"msg": "团队任务同步完成",
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
