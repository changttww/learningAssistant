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

		// 获取用户所属的所有团队ID
		var teamIDs []uint64
		// 1. 作为成员
		db.Model(&models.TeamMember{}).Where("user_id = ?", targetUserID).Pluck("team_id", &teamIDs)
		// 2. 作为Owner
		var ownedTeamIDs []uint64
		db.Model(&models.Team{}).Where("owner_user_id = ?", targetUserID).Pluck("id", &ownedTeamIDs)
		teamIDs = append(teamIDs, ownedTeamIDs...)

		// 构建查询：
		// 1. 个人任务：(owner_team_id IS NULL) AND (created_by = me OR owner_user_id = me)
		// 2. 团队任务：(owner_team_id IN my_teams) AND (parent_id IS NULL OR owner_user_id = me)
		//    即：读取团队的主任务(parent_id IS NULL) 和 分配给自己的子任务
		var tasks []models.Task

		query := db.Where("owner_team_id IS NULL").
			Where("created_by = ? OR owner_user_id = ?", targetUserID, targetUserID)

		if len(teamIDs) > 0 {
			query = query.Or(
				db.Where("owner_team_id IN ?", teamIDs).
					Where("parent_id IS NULL OR owner_user_id = ?", targetUserID),
			)
		}

		if err := query.Order("id desc").Find(&tasks).Error; err != nil {
			// 异步任务中的错误只能打印日志
			return
		}

		for _, task := range tasks {
			// 跳过标题和描述都为空的任务
			if task.Title == "" && task.Description == "" {
				continue
			}

			// 使用任务聚合方式同步知识点（任务+笔记）
			_, _ = ragService.AddTaskKnowledge(targetUserID, task.ID)
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
		if err := db.Where("user_id = ?", targetUserID).Order("id desc").Find(&notes).Error; err != nil {
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

		// 获取用户所属的所有团队ID
		var teamIDs []uint64
		db.Model(&models.TeamMember{}).Where("user_id = ?", targetUserID).Pluck("team_id", &teamIDs)
		var ownedTeamIDs []uint64
		db.Model(&models.Team{}).Where("owner_user_id = ?", targetUserID).Pluck("id", &ownedTeamIDs)
		teamIDs = append(teamIDs, ownedTeamIDs...)

		// 同步任务 - 使用任务聚合方式（任务+笔记聚合为一个知识点）
		var tasks []models.Task

		query := db.Where("owner_team_id IS NULL").
			Where("created_by = ? OR owner_user_id = ?", targetUserID, targetUserID)

		if len(teamIDs) > 0 {
			query = query.Or(
				db.Where("owner_team_id IN ?", teamIDs).
					Where("parent_id IS NULL OR owner_user_id = ?", targetUserID),
			)
		}

		query.Order("id desc").Find(&tasks)

		for _, task := range tasks {
			// 如果标题和描述都为空，跳过
			if task.Title == "" && task.Description == "" {
				continue
			}
			_, _ = ragService.AddTaskKnowledge(targetUserID, task.ID)
		}

		// 同步独立笔记（不关联任务的笔记）
		var notes []models.StudyNote
		db.Where("user_id = ? AND task_id IS NULL", targetUserID).Order("id desc").Find(&notes)

		for _, note := range notes {
			// 如果标题和内容都为空，跳过
			if note.Title == "" && note.Content == "" {
				continue
			}
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
		if err := db.Where("owner_team_id = ?", targetTeamID).Order("id desc").Find(&tasks).Error; err != nil {
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
