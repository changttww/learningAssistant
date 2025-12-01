package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	note := models.StudyNote{
		UserID:  userID.(uint64),
		TaskID:  req.TaskID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := database.GetDB().Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建笔记失败"})
		return
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
		updates["title"] = *req.Title
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

	if err := database.GetDB().Where("id = ? AND user_id = ?", noteID, userID.(uint64)).Delete(&models.StudyNote{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}
