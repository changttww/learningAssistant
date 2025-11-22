package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

// registerTaskRoutes 注册任务相关路由
func registerTaskRoutes(rg *gin.RouterGroup) {
	rg.GET("/team", handleGetTeamTasks)
	rg.GET("/:taskId", handleGetTaskDetail)
	rg.POST("/", handleCreateTask)
	rg.PUT("/:taskId/progress", handleUpdateTaskProgress)
}

// handleGetTeamTasks 返回团队任务列表
func handleGetTeamTasks(c *gin.Context) {
	db := database.GetDB()

	var tasks []models.Task
	if err := db.Where("task_type = ?", 2).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to query tasks"})
		return
	}

	// 构造响应列表，计算每个任务的进度（基于 assignees 或 status）
	var out []gin.H
	for _, t := range tasks {
		// 查询分配者进度
		var assignees []models.TaskAssignee
		_ = db.Where("task_id = ?", t.ID).Find(&assignees).Error

		progress := 0
		if len(assignees) > 0 {
			sum := 0
			for _, a := range assignees {
				sum += int(a.Progress)
			}
			progress = sum / len(assignees)
		} else {
			// fallback: use status
			if t.Status == 2 || t.CompletedAt != nil {
				progress = 100
			} else if t.Status == 1 {
				progress = 50
			} else {
				progress = 0
			}
		}

		ownerName := ""
		if t.OwnerUserID != nil {
			var u models.User
			if err := db.First(&u, *t.OwnerUserID).Error; err == nil {
				ownerName = u.DisplayName
			}
		}

		due := ""
		if t.DueAt != nil {
			due = t.DueAt.Format("2006-01-02")
		}

		out = append(out, gin.H{
			"id":           t.ID,
			"title":        t.Title,
			"description":  t.Description,
			"owner_name":   ownerName,
			"due_at":       due,
			"status":       t.Status,
			"status_label": "",
			"progress":     progress,
			"created_at":   t.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": out})
}

// handleGetTaskDetail 返回单个任务详情
func handleGetTaskDetail(c *gin.Context) {
	idStr := c.Param("taskId")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	db := database.GetDB()

	var t models.Task
	if err := db.First(&t, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": t})
}

// handleCreateTask 创建任务
func handleCreateTask(c *gin.Context) {
	var body struct {
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description"`
		TaskType    int8    `json:"task_type"`
		DueAt       *string `json:"due_at"`
		OwnerName   string  `json:"owner_name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid payload"})
		return
	}

	db := database.GetDB()
	t := models.Task{
		Title:       body.Title,
		Description: body.Description,
		TaskType:    body.TaskType,
		Status:      0,
	}

	if body.DueAt != nil && *body.DueAt != "" {
		if parsed, err := time.Parse("2006-01-02", *body.DueAt); err == nil {
			t.DueAt = &parsed
		}
	}

	if err := db.Create(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to create task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": t})
}

// handleUpdateTaskProgress 更新任务进度（将进度写入所有 assignees 并在 100% 时标记完成）
func handleUpdateTaskProgress(c *gin.Context) {
	idStr := c.Param("taskId")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid task id"})
		return
	}

	var body struct {
		Progress int `json:"progress" binding:"required"
        "`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid payload"})
		return
	}

	if body.Progress < 0 {
		body.Progress = 0
	}
	if body.Progress > 100 {
		body.Progress = 100
	}

	db := database.GetDB()

	var task models.Task
	if err := db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "task not found"})
		return
	}

	// 更新 task status based on progress
	updates := map[string]interface{}{}
	if body.Progress >= 100 {
		updates["status"] = int8(2)
		now := time.Now()
		updates["completed_at"] = &now
	} else if body.Progress > 0 {
		updates["status"] = int8(1)
	}

	if len(updates) > 0 {
		if err := db.Model(&task).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to update task status"})
			return
		}
	}

	// 更新所有 assignees 的进度（简单策略：把每个 assignee 的 progress 设为传入值）
	if err := db.Model(&models.TaskAssignee{}).Where("task_id = ?", id).Update("progress", int8(body.Progress)).Error; err != nil {
		// 仅记录错误，不影响主流程
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"id": id, "progress": body.Progress}})
}
