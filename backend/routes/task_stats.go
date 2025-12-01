package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
	taskservice "learningAssistant-backend/services/task"
)

func registerTaskStatRoutes(router *gin.RouterGroup) {
	router.GET("/stats/bar", handleGetBarStats)
	router.GET("/users/:userId/today", handleGetUserTodayTasks)
}

func handleGetBarStats(c *gin.Context) {
	rangeKey := normalizeRangeKeyWithDefault(c.Query("range"), "week")
	handleGetBarStatsWithRange(c, rangeKey)
}

func handleGetBarStatsWithRange(c *gin.Context, rangeKey string) {
	userID, ok := currentUserIDFromHeader(c)
	if !ok {
		return
	}

	stats, err := taskservice.GetBarStats(normalizeRangeKeyWithDefault(rangeKey, "week"), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取任务统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// currentUserIDFromHeader 从 Authorization Bearer 解析当前用户ID
func currentUserIDFromHeader(c *gin.Context) (uint64, bool) {
	authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权，请先登录"})
		return 0, false
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "授权信息无效"})
		return 0, false
	}

	userID, err := extractUserIDFromToken(token, "mock-token-")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "授权信息无效"})
		return 0, false
	}

	return userID, true
}

func normalizeRangeKey(rangeKey string) string {
	switch strings.ToLower(strings.TrimSpace(rangeKey)) {
	case "day", "daily":
		return "day"
	case "week", "weekly":
		return "week"
	case "month", "monthly":
		return "month"
	case "quarter", "quarterly":
		return "quarter"
	default:
		return ""
	}
}

func normalizeRangeKeyWithDefault(rangeKey, defaultKey string) string {
	normalized := normalizeRangeKey(rangeKey)
	if normalized == "" {
		return defaultKey
	}
	return normalized
}

type createTaskRequest struct {
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	TaskType        int8    `json:"task_type"`
	CategoryID      *uint64 `json:"category_id"`
	CreatedBy       uint64  `json:"created_by"`
	OwnerUserID     *uint64 `json:"owner_user_id"`
	OwnerTeamID     *uint64 `json:"owner_team_id"`
	Status          *int8   `json:"status"`
	Priority        *int8   `json:"priority"`
	StartAt         string  `json:"start_at"`
	DueAt           string  `json:"due_at"`
	EstimateMinutes *int    `json:"estimate_minutes"`
}

type taskCategoryBrief struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type todayTaskDetail struct {
	models.Task
	Progress        int                        `json:"progress"`
	Category        *taskCategoryBrief         `json:"category,omitempty"`
	Assignees       []models.TaskAssignee      `json:"assignees"`
	StatusHistory   []models.TaskStatusHistory `json:"status_history"`
	LearningRecords []models.LearningRecord    `json:"learning_records"`
}

func handleCreateTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数格式不正确"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "任务标题不能为空"})
		return
	}

	if req.TaskType != 1 && req.TaskType != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "任务类型不正确"})
		return
	}

	if req.CreatedBy == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "创建人不能为空"})
		return
	}

	startAt, err := parseISOTimePtr(req.StartAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "开始时间格式应为RFC3339"})
		return
	}
	dueAt, err := parseISOTimePtr(req.DueAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "截止时间格式应为RFC3339"})
		return
	}

	status := int8(0)
	if req.Status != nil {
		status = *req.Status
	}
	priority := int8(0)
	if req.Priority != nil {
		priority = *req.Priority
	}

	task := models.Task{
		Title:           title,
		Description:     strings.TrimSpace(req.Description),
		TaskType:        req.TaskType,
		CategoryID:      req.CategoryID,
		CreatedBy:       req.CreatedBy,
		OwnerUserID:     req.OwnerUserID,
		OwnerTeamID:     req.OwnerTeamID,
		Status:          status,
		Priority:        priority,
		StartAt:         startAt,
		DueAt:           dueAt,
		EstimateMinutes: req.EstimateMinutes,
	}

	db := database.GetDB()
	if err := db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "任务创建成功",
		"data": gin.H{
			"task": task,
		},
	})
}

func parseISOTimePtr(value string) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func handleGetUserTodayTasks(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	now := time.Now().In(time.Local)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	db := database.GetDB()
	base := db.Model(&models.Task{}).
		Where("(owner_user_id = ? OR created_by = ? OR id IN (SELECT task_id FROM task_assignees WHERE user_id = ?))", userID, userID, userID).
		Where("(start_at IS NULL OR start_at <= ?) AND (due_at IS NULL OR due_at >= ?)", endOfDay, startOfDay)

	var allTasks []models.Task
	if err := base.Order("priority DESC, due_at ASC, id DESC").Find(&allTasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取今日任务失败"})
		return
	}

	detailed, err := enrichTodayTasks(db, allTasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "整理任务关联数据失败"})
		return
	}

	completed, inProgress, notStarted := splitTodayTasks(detailed)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"date":        startOfDay.Format("2006-01-02"),
			"completed":   completed,
			"in_progress": inProgress,
			"not_started": notStarted,
		},
	})
}

func enrichTodayTasks(db *gorm.DB, tasks []models.Task) ([]todayTaskDetail, error) {
	if len(tasks) == 0 {
		return []todayTaskDetail{}, nil
	}

	taskIDs := make([]uint64, 0, len(tasks))
	categoryIDs := make([]uint64, 0, len(tasks))
	categorySeen := make(map[uint64]struct{})

	for _, t := range tasks {
		taskIDs = append(taskIDs, t.ID)
		if t.CategoryID != nil {
			if _, ok := categorySeen[*t.CategoryID]; !ok {
				categoryIDs = append(categoryIDs, *t.CategoryID)
				categorySeen[*t.CategoryID] = struct{}{}
			}
		}
	}

	assigneeMap := make(map[uint64][]models.TaskAssignee, len(taskIDs))
	var assignees []models.TaskAssignee
	if err := db.Where("task_id IN ?", taskIDs).Order("is_owner DESC, id ASC").Find(&assignees).Error; err != nil {
		return nil, err
	}
	for _, a := range assignees {
		assigneeMap[a.TaskID] = append(assigneeMap[a.TaskID], a)
	}

	statusHistoryMap := make(map[uint64][]models.TaskStatusHistory, len(taskIDs))
	var histories []models.TaskStatusHistory
	if err := db.Where("task_id IN ?", taskIDs).Order("created_at ASC, id ASC").Find(&histories).Error; err != nil {
		return nil, err
	}
	for _, h := range histories {
		statusHistoryMap[h.TaskID] = append(statusHistoryMap[h.TaskID], h)
	}

	learningRecordMap := make(map[uint64][]models.LearningRecord, len(taskIDs))
	var records []models.LearningRecord
	if err := db.Where("task_id IN ?", taskIDs).Order("session_start ASC, id ASC").Find(&records).Error; err != nil {
		return nil, err
	}
	for _, r := range records {
		learningRecordMap[r.TaskID] = append(learningRecordMap[r.TaskID], r)
	}

	categoryMap := make(map[uint64]taskCategoryBrief, len(categoryIDs))
	if len(categoryIDs) > 0 {
		var categories []models.TaskCategory
		if err := db.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
			return nil, err
		}
		for _, cat := range categories {
			categoryMap[cat.ID] = taskCategoryBrief{
				ID:    cat.ID,
				Name:  cat.Name,
				Color: cat.Color,
			}
		}
	}

	detailed := make([]todayTaskDetail, 0, len(tasks))
	for _, t := range tasks {
		assigneeList := assigneeMap[t.ID]
		statusHistory := statusHistoryMap[t.ID]
		learningRecords := learningRecordMap[t.ID]

		if assigneeList == nil {
			assigneeList = []models.TaskAssignee{}
		}
		if statusHistory == nil {
			statusHistory = []models.TaskStatusHistory{}
		}
		if learningRecords == nil {
			learningRecords = []models.LearningRecord{}
		}

		var category *taskCategoryBrief
		if t.CategoryID != nil {
			if brief, ok := categoryMap[*t.CategoryID]; ok {
				category = &brief
			}
		}

		detailed = append(detailed, todayTaskDetail{
			Task:            t,
			Progress:        deriveTaskProgress(t.Status, assigneeList),
			Category:        category,
			Assignees:       assigneeList,
			StatusHistory:   statusHistory,
			LearningRecords: learningRecords,
		})
	}

	return detailed, nil
}

func splitTodayTasks(tasks []todayTaskDetail) (completed, inProgress, notStarted []todayTaskDetail) {
	for _, task := range tasks {
		switch task.Status {
		case 2:
			completed = append(completed, task)
		case 1:
			inProgress = append(inProgress, task)
		default:
			notStarted = append(notStarted, task)
		}
	}
	return
}

func deriveTaskProgress(status int8, assignees []models.TaskAssignee) int {
	progress := 0
	for _, a := range assignees {
		if int(a.Progress) > progress {
			progress = int(a.Progress)
		}
	}

	if progress > 0 {
		if progress > 100 {
			return 100
		}
		return progress
	}

	switch status {
	case 2:
		return 100
	case 1:
		return 60
	default:
		return 0
	}
}
