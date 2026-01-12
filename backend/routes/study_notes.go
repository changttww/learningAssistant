package routes

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/models"
)

type createNoteRequest struct {
	Title   string  `json:"title"`
	TaskID  *uint64 `json:"task_id"`
	Content string  `json:"content"`
}

type updateNoteRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func registerStudyNotesRoutes(router *gin.RouterGroup) {
	notes := router.Group("/notes")
	notes.Use(middleware.AuthMiddleware())
	notes.GET("", handleListNotes)
	notes.POST("", handleCreateNote)
	notes.PUT(":id", handleUpdateNote)
	notes.DELETE(":id", handleDeleteNote)
}

func handleListNotes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var notes []models.StudyNote
	db := database.GetDB().Model(&models.StudyNote{})
	db = db.Where("user_id = ?", userID.(uint64))

	if taskIDStr := c.Query("task_id"); taskIDStr != "" {
		if taskID, err := strconv.ParseUint(taskIDStr, 10, 64); err == nil {
			db = db.Where("task_id = ?", taskID)
		}
	}

	if err := db.Order("created_at DESC").Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": notes,
		"msg":  "获取成功",
	})
}

func handleCreateNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req createNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "message": "标题不能为空"})
		return
	}

	// 同一用户不允许创建同标题笔记
	var existing models.StudyNote
	if err := database.GetDB().Model(&models.StudyNote{}).
		Select("id").
		Where("user_id = ? AND title = ?", userID.(uint64), title).
		Limit(1).
		First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询笔记失败"})
			return
		}
	} else {
		c.JSON(http.StatusConflict, gin.H{
			"error":            "创建笔记失败",
			"message":          "已存在同标题的笔记，不允许重复创建",
			"existing_note_id": existing.ID,
		})
		return
	}

	note := models.StudyNote{
		UserID:  userID.(uint64),
		TaskID:  req.TaskID,
		Title:   title,
		Content: req.Content,
	}

	if err := database.GetDB().Create(&note).Error; err != nil {
		// 如果数据库层面存在唯一约束，兜底把重复键映射成 409（避免误导成“标题重复”）
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "创建笔记失败",
				"message": "数据冲突，创建失败（可能已存在相同标题的笔记）",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建笔记失败"})
		return
	}

	// 自动更新任务知识点（如果笔记关联了任务）
	if ragService != nil && req.TaskID != nil && *req.TaskID > 0 {
		go func() {
			_, _ = ragService.AddTaskKnowledge(userID.(uint64), *req.TaskID)
		}()
	} else if ragService != nil && req.Content != "" {
		// 独立笔记（不关联任务）仍使用原方式添加
		go func() {
			_, _ = ragService.AddDocument(userID.(uint64), 2, note.ID, req.Title, req.Content)
		}()
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": note,
		"msg":  "创建成功",
	})
}

func handleUpdateNote(c *gin.Context) {
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var note models.StudyNote
	if err := database.GetDB().Where("id = ? AND user_id = ?", noteID, userID.(uint64)).First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询笔记失败"})
		}
		return
	}

	var req updateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		newTitle := strings.TrimSpace(*req.Title)
		if newTitle == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "message": "标题不能为空"})
			return
		}

		// 同一用户不允许把标题改成已存在的标题（排除自己）
		var existing models.StudyNote
		if err := database.GetDB().Model(&models.StudyNote{}).
			Select("id").
			Where("user_id = ? AND title = ? AND id <> ?", userID.(uint64), newTitle, noteID).
			Limit(1).
			First(&existing).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "查询笔记失败"})
				return
			}
		} else {
			c.JSON(http.StatusConflict, gin.H{
				"error":            "更新笔记失败",
				"message":          "已存在同标题的笔记，不允许使用该标题",
				"existing_note_id": existing.ID,
			})
			return
		}

		updates["title"] = newTitle
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}

	if err := database.GetDB().Model(&note).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	if err := database.GetDB().First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取更新后的笔记失败"})
		return
	}

	// 更新知识库中的对应条目（如果内容有变化）
	if ragService != nil && req.Content != nil && *req.Content != "" {
		go func() {
			// 如果笔记关联了任务，更新任务知识点
			if note.TaskID != nil && *note.TaskID > 0 {
				_, _ = ragService.AddTaskKnowledge(userID.(uint64), *note.TaskID)
			} else {
				// 独立笔记使用原方式更新
				_, _ = ragService.AddDocument(userID.(uint64), 2, note.ID, note.Title, note.Content)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": note,
		"msg":  "更新成功",
	})
}

func handleDeleteNote(c *gin.Context) {
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 原子化删除：同时删除笔记和关联的知识库条目
	if err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		// 1. 先删除基于该笔记的知识库条目的向量缓存
		tx.Where("entry_id IN (?)",
			tx.Model(&models.KnowledgeBaseEntry{}).Select("id").
				Where("user_id = ? AND source_type = 2 AND source_id = ?", userID.(uint64), noteID)).
			Delete(&models.KnowledgeVectorCache{})

		// 2. 删除知识关系
		tx.Where("user_id = ? AND (source_entry_id IN (?) OR target_entry_id IN (?))",
			userID.(uint64),
			tx.Model(&models.KnowledgeBaseEntry{}).Select("id").
				Where("user_id = ? AND source_type = 2 AND source_id = ?", userID.(uint64), noteID),
			tx.Model(&models.KnowledgeBaseEntry{}).Select("id").
				Where("user_id = ? AND source_type = 2 AND source_id = ?", userID.(uint64), noteID)).
			Delete(&models.KnowledgeRelation{})

		// 3. 删除知识库条目（source_type=2 表示学习笔记）
		tx.Where("user_id = ? AND source_type = 2 AND source_id = ?", userID.(uint64), noteID).
			Delete(&models.KnowledgeBaseEntry{})

		// 4. 删除笔记
		if err := tx.Where("id = ? AND user_id = ?", noteID, userID.(uint64)).Delete(&models.StudyNote{}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}
