package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/models"
)

// CreateTaskRequest 创建任务请求结构
type CreateTaskRequest struct {
	Title           string     `json:"title" binding:"required"`
	Description     string     `json:"description"`
	TaskType        int8       `json:"task_type" binding:"required"`
	CategoryID      *uint64    `json:"category_id"`
	Priority        int8       `json:"priority"`
	StartAt         *time.Time `json:"start_at"`
	DueAt           *time.Time `json:"due_at"`
	EstimateMinutes *int       `json:"estimate_minutes"`
	EffortPoints    int        `json:"effort_points"`
}

// UpdateTaskRequest 更新任务请求结构
type UpdateTaskRequest struct {
	Title           *string    `json:"title"`
	Description     *string    `json:"description"`
	CategoryID      *uint64    `json:"category_id"`
	Priority        *int8      `json:"priority"`
	Status          *int8      `json:"status"`
	StartAt         *time.Time `json:"start_at"`
	DueAt           *time.Time `json:"due_at"`
	EstimateMinutes *int       `json:"estimate_minutes"`
	EffortPoints    *int       `json:"effort_points"`
}

// TaskResponse 任务响应结构
type TaskResponse struct {
	ID              uint64                `json:"id"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	TaskType        int8                  `json:"task_type"`
	CategoryID      *uint64               `json:"category_id"`
	Category        *TaskCategoryResponse `json:"category,omitempty"`
	CreatedBy       uint64                `json:"created_by"`
	OwnerUserID     *uint64               `json:"owner_user_id"`
	OwnerTeamID     *uint64               `json:"owner_team_id"`
	Status          int8                  `json:"status"`
	Priority        int8                  `json:"priority"`
	StartAt         *time.Time            `json:"start_at"`
	DueAt           *time.Time            `json:"due_at"`
	CompletedAt     *time.Time            `json:"completed_at"`
	EstimateMinutes *int                  `json:"estimate_minutes"`
	EffortPoints    int                   `json:"effort_points"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

// TaskCategoryResponse 任务分类响应结构
type TaskCategoryResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// registerTaskRoutes 注册任务路由
func registerTaskRoutes(r *gin.RouterGroup) {
	// 应用认证中间件到所有任务路由
	r.Use(middleware.AuthMiddleware())

	r.POST("", createTask)
	r.POST("/", createTask) // 兼容两种路径
	r.GET("", getTaskList)
	r.GET("/", getTaskList) // 兼容两种路径
	r.GET("/personal", getPersonalTasks)
	r.GET("/team", getTeamTasks)
	r.GET("/:id", getTaskDetail)
	r.PUT("/:id", updateTask)
	r.DELETE("/:id", deleteTask)
	r.POST("/:id/complete", completeTask)
	r.POST("/:id/uncomplete", uncompleteTask)
	r.GET("/categories", getTaskCategories)
	r.GET("/statistics", getTaskStatistics)
}

// createTask 创建任务
func createTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT中获取用户ID (这里假设你已经有JWT中间件)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	task := models.Task{
		Title:           req.Title,
		Description:     req.Description,
		TaskType:        req.TaskType,
		CategoryID:      req.CategoryID,
		CreatedBy:       userID.(uint64),
		Priority:        req.Priority,
		StartAt:         req.StartAt,
		DueAt:           req.DueAt,
		EstimateMinutes: req.EstimateMinutes,
		EffortPoints:    req.EffortPoints,
		Status:          0, // 默认状态为待处理
	}

	// 如果是个人任务，设置所有者为当前用户
	if req.TaskType == 1 {
		userIDValue := userID.(uint64)
		task.OwnerUserID = &userIDValue
	}

	if err := database.GetDB().Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建任务失败"})
		return
	}

	// 加载分类信息
	var taskWithCategory models.Task
	database.GetDB().Preload("Category").First(&taskWithCategory, task.ID)

	response := convertTaskToResponse(taskWithCategory)
	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": response,
		"msg":  "任务创建成功",
	})
}

// getTaskList 获取任务列表
func getTaskList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var tasks []models.Task
	db := database.GetDB().Preload("Category")

	// 获取用户相关的任务
	db = db.Where("created_by = ? OR owner_user_id = ?", userID.(uint64), userID.(uint64))

	if err := db.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务列表失败"})
		return
	}

	var responses []TaskResponse
	for _, task := range tasks {
		responses = append(responses, convertTaskToResponse(task))
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": responses,
		"msg":  "获取成功",
	})
}

// getPersonalTasks 获取个人任务列表
func getPersonalTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var tasks []models.Task
	db := database.GetDB().Preload("Category")

	// 只获取个人任务 (task_type = 1)
	db = db.Where("task_type = ? AND (created_by = ? OR owner_user_id = ?)", 1, userID.(uint64), userID.(uint64))

	// 支持状态过滤
	if status := c.Query("status"); status != "" {
		if statusInt, err := strconv.Atoi(status); err == nil {
			db = db.Where("status = ?", statusInt)
		}
	}

	// 支持优先级过滤
	if priority := c.Query("priority"); priority != "" {
		if priorityInt, err := strconv.Atoi(priority); err == nil {
			db = db.Where("priority = ?", priorityInt)
		}
	}

	// 支持分类过滤
	if categoryID := c.Query("category_id"); categoryID != "" {
		if categoryIDInt, err := strconv.ParseUint(categoryID, 10, 64); err == nil {
			db = db.Where("category_id = ?", categoryIDInt)
		}
	}

	// 排序
	orderBy := c.DefaultQuery("order_by", "created_at")
	orderDir := c.DefaultQuery("order_dir", "desc")
	db = db.Order(orderBy + " " + orderDir)

	if err := db.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取个人任务失败"})
		return
	}

	var responses []TaskResponse
	for _, task := range tasks {
		responses = append(responses, convertTaskToResponse(task))
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": responses,
		"msg":  "获取成功",
	})
}

// getTeamTasks 获取团队任务列表
func getTeamTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var tasks []models.Task
	db := database.GetDB().Preload("Category")

	// 获取团队任务 (task_type = 2) - 这里需要根据用户的团队关系进行查询
	// 暂时简单实现为创建者或所有者是当前用户的团队任务
	db = db.Where("task_type = ? AND (created_by = ? OR owner_team_id IN (SELECT team_id FROM user_teams WHERE user_id = ?))",
		2, userID.(uint64), userID.(uint64))

	if err := db.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取团队任务失败"})
		return
	}

	var responses []TaskResponse
	for _, task := range tasks {
		responses = append(responses, convertTaskToResponse(task))
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": responses,
		"msg":  "获取成功",
	})
}

// getTaskDetail 获取任务详情
func getTaskDetail(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var task models.Task
	db := database.GetDB().Preload("Category")

	// 确保用户只能访问自己相关的任务
	if err := db.Where("id = ? AND (created_by = ? OR owner_user_id = ?)", taskID, userID.(uint64), userID.(uint64)).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或无权限访问"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务详情失败"})
		}
		return
	}

	response := convertTaskToResponse(task)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": response,
		"msg":  "获取成功",
	})
}

// updateTask 更新任务
func updateTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查任务是否存在且用户有权限
	var task models.Task
	if err := database.GetDB().Where("id = ? AND (created_by = ? OR owner_user_id = ?)", taskID, userID.(uint64), userID.(uint64)).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或无权限访问"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询任务失败"})
		}
		return
	}

	// 更新字段
	updateData := make(map[string]interface{})
	if req.Title != nil {
		updateData["title"] = *req.Title
	}
	if req.Description != nil {
		updateData["description"] = *req.Description
	}
	if req.CategoryID != nil {
		updateData["category_id"] = *req.CategoryID
	}
	if req.Priority != nil {
		updateData["priority"] = *req.Priority
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
		// 如果状态改为已完成，记录完成时间
		if *req.Status == 2 {
			now := time.Now()
			updateData["completed_at"] = now
		} else {
			updateData["completed_at"] = nil
		}
	}
	if req.StartAt != nil {
		updateData["start_at"] = *req.StartAt
	}
	if req.DueAt != nil {
		updateData["due_at"] = *req.DueAt
	}
	if req.EstimateMinutes != nil {
		updateData["estimate_minutes"] = *req.EstimateMinutes
	}
	if req.EffortPoints != nil {
		updateData["effort_points"] = *req.EffortPoints
	}

	if err := database.GetDB().Model(&task).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新任务失败"})
		return
	}

	// 重新查询更新后的任务
	database.GetDB().Preload("Category").First(&task, taskID)
	response := convertTaskToResponse(task)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": response,
		"msg":  "更新成功",
	})
}

// deleteTask 删除任务
func deleteTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 检查任务是否存在且用户有权限
	var task models.Task
	if err := database.GetDB().Where("id = ? AND (created_by = ? OR owner_user_id = ?)", taskID, userID.(uint64), userID.(uint64)).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或无权限访问"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询任务失败"})
		}
		return
	}

	if err := database.GetDB().Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}

// completeTask 完成任务
func completeTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 检查任务是否存在且用户有权限
	var task models.Task
	if err := database.GetDB().Where("id = ? AND (created_by = ? OR owner_user_id = ?)", taskID, userID.(uint64), userID.(uint64)).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或无权限访问"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询任务失败"})
		}
		return
	}

	now := time.Now()
	updateData := map[string]interface{}{
		"status":       2, // 已完成
		"completed_at": now,
	}

	if err := database.GetDB().Model(&task).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "完成任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "任务已完成",
	})
}

// uncompleteTask 取消完成任务
func uncompleteTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 检查任务是否存在且用户有权限
	var task models.Task
	if err := database.GetDB().Where("id = ? AND (created_by = ? OR owner_user_id = ?)", taskID, userID.(uint64), userID.(uint64)).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或无权限访问"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询任务失败"})
		}
		return
	}

	updateData := map[string]interface{}{
		"status":       1, // 进行中
		"completed_at": nil,
	}

	if err := database.GetDB().Model(&task).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消完成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "已取消完成状态",
	})
}

// getTaskCategories 获取任务分类
func getTaskCategories(c *gin.Context) {
	var categories []models.TaskCategory
	if err := database.GetDB().Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类失败"})
		return
	}

	var responses []TaskCategoryResponse
	for _, category := range categories {
		responses = append(responses, TaskCategoryResponse{
			ID:    category.ID,
			Name:  category.Name,
			Color: category.Color,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": responses,
		"msg":  "获取成功",
	})
}

// getTaskStatistics 获取任务统计
func getTaskStatistics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	type TaskStats struct {
		Total      int64 `json:"total"`
		Completed  int64 `json:"completed"`
		InProgress int64 `json:"in_progress"`
		Pending    int64 `json:"pending"`
		Overdue    int64 `json:"overdue"`
	}

	var stats TaskStats
	db := database.GetDB().Model(&models.Task{}).Where("task_type = ? AND (created_by = ? OR owner_user_id = ?)", 1, userID.(uint64), userID.(uint64))

	// 总任务数
	db.Count(&stats.Total)

	// 已完成任务数
	db.Where("status = ?", 2).Count(&stats.Completed)

	// 进行中任务数
	db.Where("status = ?", 1).Count(&stats.InProgress)

	// 待处理任务数
	db.Where("status = ?", 0).Count(&stats.Pending)

	// 逾期任务数
	now := time.Now()
	db.Where("status != ? AND due_at < ?", 2, now).Count(&stats.Overdue)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
		"msg":  "获取成功",
	})
}

// convertTaskToResponse 将Task模型转换为响应结构
func convertTaskToResponse(task models.Task) TaskResponse {
	response := TaskResponse{
		ID:              task.ID,
		Title:           task.Title,
		Description:     task.Description,
		TaskType:        task.TaskType,
		CategoryID:      task.CategoryID,
		CreatedBy:       task.CreatedBy,
		OwnerUserID:     task.OwnerUserID,
		OwnerTeamID:     task.OwnerTeamID,
		Status:          task.Status,
		Priority:        task.Priority,
		StartAt:         task.StartAt,
		DueAt:           task.DueAt,
		CompletedAt:     task.CompletedAt,
		EstimateMinutes: task.EstimateMinutes,
		EffortPoints:    task.EffortPoints,
		CreatedAt:       task.CreatedAt,
		UpdatedAt:       task.UpdatedAt,
	}

	// 如果有分类信息，添加到响应中
	if task.Category != nil {
		response.Category = &TaskCategoryResponse{
			ID:    task.Category.ID,
			Name:  task.Category.Name,
			Color: task.Category.Color,
		}
	}

	return response
}
