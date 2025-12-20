package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

// registerAnalysisRoutes 注册学习效率分析路由
func registerAnalysisRoutes(router *gin.RouterGroup) {
	// 采用 POST 便于后续扩展过滤参数
	router.POST("/efficiency", handleEfficiencyAnalysis)
	router.POST("/seed-month-study", seedUserMonthStudyData)
	router.GET("/weekly-hours", handleGetWeeklyStudyHours)
}

type weeklyHoursDaily struct {
	Date    string  `json:"date"`
	Minutes int     `json:"minutes"`
	Hours   float64 `json:"hours"`
}

type weeklyHoursResponse struct {
	UserID       uint64             `json:"user_id"`
	From         string             `json:"from"`
	To           string             `json:"to"`
	TotalMinutes int                `json:"total_minutes"`
	TotalHours   float64            `json:"total_hours"`
	Daily        []weeklyHoursDaily `json:"daily"`
}

// handleGetWeeklyStudyHours 返回最近 7 天（含当天）学习时长，单位小时
func handleGetWeeklyStudyHours(c *gin.Context) {
	userIDStr := strings.TrimSpace(c.Query("user_id"))
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 user_id"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "user_id 不正确"})
		return
	}

	db := database.GetDB()
	today := startOfDay(time.Now())
	from := today.AddDate(0, 0, -6)

	var stats []models.DailyStudyStat
	if err := db.Where("user_id = ? AND date >= ? AND date <= ?", userID, from, today).
		Order("date ASC").
		Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取学习时长失败"})
		return
	}

	byDate := make(map[string]models.DailyStudyStat, len(stats))
	for _, s := range stats {
		byDate[s.Date.Format("2006-01-02")] = s
	}

	daily := make([]weeklyHoursDaily, 0, 7)
	totalMinutes := 0
	for i := 0; i < 7; i++ {
		day := from.AddDate(0, 0, i)
		key := day.Format("2006-01-02")
		minutes := 0
		if stat, ok := byDate[key]; ok {
			minutes = stat.Minutes
		}
		totalMinutes += minutes
		daily = append(daily, weeklyHoursDaily{
			Date:    key,
			Minutes: minutes,
			Hours:   math.Round((float64(minutes)/60.0)*100) / 100,
		})
	}

	resp := weeklyHoursResponse{
		UserID:       userID,
		From:         from.Format("2006-01-02"),
		To:           today.Format("2006-01-02"),
		TotalMinutes: totalMinutes,
		TotalHours:   math.Round((float64(totalMinutes)/60.0)*100) / 100,
		Daily:        daily,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

func startOfDay(t time.Time) time.Time {
	loc := t.Location()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

type efficiencyAnalysisRequest struct {
	UserID       uint64 `json:"user_id"`
	Days         int    `json:"days"`
	Model        string `json:"model"`
	ForceRefresh bool   `json:"force_refresh"`
}

type aiError struct {
	Status  int
	Message string
}

func (e aiError) Error() string {
	return e.Message
}

type analysisSummaryView struct {
	WeeklyStudyHours   float64  `json:"weekly_study_hours"`
	FocusScore         int      `json:"focus_score"`
	TaskCompletionRate int      `json:"task_completion_rate"`
	StreakDays         int      `json:"streak_days"`
	StrengthHighlights []string `json:"strength_highlights,omitempty"`
	RiskHighlights     []string `json:"risk_highlights,omitempty"`
}

type trendPointView struct {
	Date       string  `json:"date"`
	StudyHours float64 `json:"study_hours"`
	FocusScore int     `json:"focus_score"`
}

type recommendationView struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Impact string `json:"impact"`
}

type taskInsightView struct {
	TaskTitle string `json:"task_title"`
	Status    string `json:"status"`
	Advice    string `json:"advice"`
	Risk      string `json:"risk,omitempty"`
}

type reviewItemView struct {
	Subject  string `json:"subject"`
	Priority string `json:"priority"`
	Progress int    `json:"progress"`
	DueDate  string `json:"due_date"`
}

type reminderView struct {
	Content string `json:"content"`
	Time    string `json:"time"`
}

type knowledgeMapView struct {
	Mastered int `json:"mastered"`
	Learning int `json:"learning"`
	ToLearn  int `json:"to_learn"`
}

type reviewPlanView struct {
	Summary      string           `json:"summary,omitempty"`
	ReviewItems  []reviewItemView `json:"review_items,omitempty"`
	Reminders    []reminderView   `json:"reminders,omitempty"`
	KnowledgeMap knowledgeMapView `json:"knowledge_map"`
}

type efficiencyAnalysisView struct {
	Source          string               `json:"source"`
	Model           string               `json:"model"`
	GeneratedAt     string               `json:"generated_at"`
	Summary         analysisSummaryView  `json:"summary"`
	StudyTrend      []trendPointView     `json:"study_trend"`
	Recommendations []recommendationView `json:"recommendations"`
	TaskInsights    []taskInsightView    `json:"task_insights"`
	ReviewPlan      reviewPlanView       `json:"review_plan,omitempty"`
	Prompt          string               `json:"prompt,omitempty"`
	Notes           string               `json:"notes,omitempty"`
}

type analysisContext struct {
	User        userBriefContext    `json:"user"`
	Profile     profileBriefContext `json:"profile"`
	TimeRange   timeRangeContext    `json:"time_range"`
	Daily       []dailyStudyStat    `json:"daily_study"`
	Sessions    []sessionBrief      `json:"recent_sessions"`
	TaskStats   taskStatSnapshot    `json:"task_stats"`
	TaskSamples []taskSample        `json:"task_samples"`
}

type userBriefContext struct {
	ID          uint64 `json:"id"`
	DisplayName string `json:"display_name"`
}

type profileBriefContext struct {
	StreakDays     int     `json:"streak_days"`
	TaskDone       int     `json:"tasks_completed"`
	TaskInProgress int     `json:"tasks_in_progress"`
	CompletionRate float32 `json:"completion_rate"`
	TotalStudyMins int     `json:"total_study_minutes"`
}

type timeRangeContext struct {
	Days int    `json:"days"`
	From string `json:"from"`
	To   string `json:"to"`
}

type dailyStudyStat struct {
	Date         string   `json:"date"`
	TotalMinutes int      `json:"total_minutes"`
	SessionCount int      `json:"session_count"`
	FirstSession string   `json:"first_session,omitempty"`
	LastSession  string   `json:"last_session,omitempty"`
	Notes        []string `json:"notes,omitempty"`
}

type sessionBrief struct {
	TaskID          uint64 `json:"task_id"`
	DurationMinutes int    `json:"duration_minutes"`
	StartAt         string `json:"start_at"`
	EndAt           string `json:"end_at"`
	Note            string `json:"note,omitempty"`
}

type taskStatSnapshot struct {
	Total          int `json:"total"`
	Completed      int `json:"completed"`
	InProgress     int `json:"in_progress"`
	Pending        int `json:"pending"`
	CompletionRate int `json:"completion_rate"`
}

type taskSample struct {
	ID       uint64  `json:"id"`
	Title    string  `json:"title"`
	Status   string  `json:"status"`
	Priority int8    `json:"priority"`
	DueDate  *string `json:"due_date,omitempty"`
}

type llmPayload struct {
	Model        string        `json:"model"`
	Temperature  float64       `json:"temperature"`
	SystemPrompt string        `json:"system_prompt"`
	UserPrompt   string        `json:"user_prompt"`
	Schema       llmSchemaSpec `json:"schema"`
}

type llmSchemaSpec struct {
	Name   string                 `json:"name"`
	Schema map[string]interface{} `json:"schema"`
}

type llmAnalysis struct {
	Summary struct {
		WeeklyStudyHours   float64  `json:"weekly_study_hours"`
		FocusScore         int      `json:"focus_score"`
		TaskCompletionRate int      `json:"task_completion_rate"`
		StreakDays         int      `json:"streak_days"`
		KeyStrengths       []string `json:"key_strengths"`
		Risks              []string `json:"risks"`
	} `json:"summary"`
	StudyTrend []struct {
		Date       string  `json:"date"`
		StudyHours float64 `json:"study_hours"`
		FocusScore int     `json:"focus_score"`
	} `json:"study_trend"`
	Recommendations []struct {
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Impact string `json:"impact"`
	} `json:"recommendations"`
	TaskInsights []struct {
		TaskTitle string `json:"task_title"`
		Status    string `json:"status"`
		Advice    string `json:"advice"`
		Risk      string `json:"risk"`
	} `json:"task_insights"`
	ReviewPlan struct {
		Summary      string `json:"summary"`
		KnowledgeMap struct {
			Mastered int `json:"mastered"`
			Learning int `json:"learning"`
			ToLearn  int `json:"to_learn"`
		} `json:"knowledge_map"`
		ReviewItems []struct {
			Subject  string `json:"subject"`
			Priority string `json:"priority"`
			Progress int    `json:"progress"`
			DueDate  string `json:"due_date"`
		} `json:"review_items"`
		Reminders []struct {
			Content string `json:"content"`
			Time    string `json:"time"`
		} `json:"reminders"`
	} `json:"review_plan"`
}

// 宽松解析结构，用于兼容模型未按 schema 返回的字段
type looseAnalysis struct {
	LearningEfficiencySummary struct {
		StreakDays        int     `json:"streak_days"`
		TotalStudyMinutes int     `json:"total_study_minutes"`
		CompletionRate    int     `json:"completion_rate"`
		TasksCompleted    int     `json:"tasks_completed"`
		InProgressRatio   float64 `json:"in_progress_ratio"`
		PendingTasks      int     `json:"pending_tasks"`
		Insight           string  `json:"insight"`
	} `json:"learning_efficiency_summary"`
	KeyTaskInsights []struct {
		Issue    string      `json:"issue"`
		Evidence interface{} `json:"evidence"` // 可能是字符串或字符串数组
		Impact   string      `json:"impact"`
	} `json:"key_task_insights"`
	ReviewReminderList []struct {
		Action     string `json:"action"`
		TaskID     *int   `json:"task_id"`
		Title      string `json:"title"`
		Priority   int    `json:"priority"`
		DueDate    string `json:"due_date"`
		Suggestion string `json:"suggestion"`
	} `json:"review_reminder_list"`
	ReviewAndReminders []struct {
		Type     string `json:"type"`
		Priority string `json:"priority"`
		Title    string `json:"title"`
		Detail   string `json:"detail"`
		Deadline string `json:"deadline"`
	} `json:"review_and_reminders"`
}

type llmMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type llmResponseFormat struct {
	Type       string `json:"type"`
	JsonSchema struct {
		Name   string                 `json:"name"`
		Schema map[string]interface{} `json:"schema"`
		Strict bool                   `json:"strict"`
	} `json:"json_schema"`
}

type llmChatInput struct {
	Messages []llmMessage `json:"messages"`
}

type llmChatParameters struct {
	Temperature  float64 `json:"temperature,omitempty"`
	ResultFormat string  `json:"result_format,omitempty"`
}

type llmChatRequest struct {
	Model          string             `json:"model"`
	Input          llmChatInput       `json:"input"`
	Parameters     llmChatParameters  `json:"parameters"`
	ResponseFormat *llmResponseFormat `json:"response_format,omitempty"`
}

type llmChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	FinishReason string `json:"finish_reason,omitempty"`
}

type llmChatResponse struct {
	Model      string `json:"model,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
	TaskStatus string `json:"task_status,omitempty"`
	Output     struct {
		Choices []llmChoice `json:"choices"`
		Text    string      `json:"text,omitempty"`
	} `json:"output"`
	Usage struct {
		InputTokens  int `json:"input_tokens,omitempty"`
		OutputTokens int `json:"output_tokens,omitempty"`
		TotalTokens  int `json:"total_tokens,omitempty"`
	} `json:"usage,omitempty"`
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// handleEfficiencyAnalysis 生成学习效率分析报告
func handleEfficiencyAnalysis(c *gin.Context) {
	var req efficiencyAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	if req.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "user_id 必填"})
		return
	}
	log.Printf("[analysis] incoming request user=%d days=%d model=%s force=%v", req.UserID, req.Days, req.Model, req.ForceRefresh)
	days := req.Days
	if days <= 0 {
		days = 14
	}
	if days > 60 {
		days = 60
	}

	ctxData, err := buildAnalysisContext(req.UserID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	if isContextEmpty(ctxData) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "暂无学习记录，无法生成分析"})
		return
	}
	log.Printf(
		"[analysis] context ready user=%d days=%d daily=%d sessions=%d tasks=%d completion=%d%%",
		req.UserID,
		days,
		len(ctxData.Daily),
		len(ctxData.Sessions),
		len(ctxData.TaskSamples),
		ctxData.TaskStats.CompletionRate,
	)

	payload, prompt := buildLLMPayload(ctxData, req.Model)
	log.Printf(
		"[analysis] built LLM payload model=%s schema=%s",
		payload.Model,
		payload.Schema.Name,
	)
	analysisView, llmErr := runLLMAnalysis(payload, prompt)

	if llmErr != nil {
		log.Printf("[analysis] LLM failed: %v", llmErr)
		var aErr aiError
		if errors.As(llmErr, &aErr) && aErr.Status != 0 {
			c.JSON(aErr.Status, gin.H{"code": aErr.Status, "message": aErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "AI分析失败，请稍后重试", "error": llmErr.Error()})
		return
	}

	log.Printf(
		"[analysis] LLM success model=%s recs=%d tasks=%d trend=%d",
		analysisView.Model,
		len(analysisView.Recommendations),
		len(analysisView.TaskInsights),
		len(analysisView.StudyTrend),
	)

	log.Printf("[analysis] response ready source=%s model=%s", analysisView.Source, analysisView.Model)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"analysis": analysisView,
		},
	})
}

func buildAnalysisContext(userID uint64, days int) (analysisContext, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if errorsIsNotFound(err) {
			return analysisContext{}, fmt.Errorf("用户不存在")
		}
		return analysisContext{}, fmt.Errorf("加载用户失败")
	}

	var profile models.UserProfile
	_ = db.Where("user_id = ?", userID).First(&profile).Error

	start := time.Now().AddDate(0, 0, -days)

	// 先加载每日聚合，保证 AI 用到的每日数据是真实统计
	var dailyStats []models.DailyStudyStat
	_ = db.Where("user_id = ? AND date >= ? AND date <= ?", userID, startOfDay(start), startOfDay(time.Now())).
		Order("date ASC").
		Find(&dailyStats).Error

	dailyMap := map[string]*dailyStudyStat{}
	for _, stat := range dailyStats {
		dayKey := stat.Date.Format("2006-01-02")
		dailyMap[dayKey] = &dailyStudyStat{
			Date:         dayKey,
			TotalMinutes: stat.Minutes,
			SessionCount: stat.SessionCount,
		}
	}

	// 加载真实学习会话（StudySession），补充 session 样本与首/尾时间
	var sessions []models.StudySession
	if err := db.Where("user_id = ? AND end_time IS NOT NULL AND start_time >= ?", userID, start).
		Order("start_time ASC").
		Find(&sessions).Error; err != nil {
		return analysisContext{}, fmt.Errorf("加载学习会话失败")
	}

	sessionSamples := make([]sessionBrief, 0, len(sessions))
	for _, sess := range sessions {
		if sess.EndTime == nil {
			continue
		}
		dayKey := sess.StartTime.Format("2006-01-02")
		stat, ok := dailyMap[dayKey]
		if !ok {
			stat = &dailyStudyStat{Date: dayKey}
			dailyMap[dayKey] = stat
		}
		minutes := sess.DurationMinutes
		if minutes == 0 {
			minutes = int(sess.EndTime.Sub(sess.StartTime).Minutes())
			if minutes < 1 {
				minutes = 1
			}
		}
		stat.TotalMinutes += minutes
		stat.SessionCount++
		if stat.FirstSession == "" {
			stat.FirstSession = sess.StartTime.Format(time.RFC3339)
		}
		stat.LastSession = sess.EndTime.Format(time.RFC3339)
		note := strings.TrimSpace(sess.Note)
		if note == "" {
			note = fmt.Sprintf("source:%s", strings.TrimSpace(sess.Source))
		}
		if note != "" && len(stat.Notes) < 3 {
			stat.Notes = append(stat.Notes, note)
		}

		sessionSamples = append(sessionSamples, sessionBrief{
			TaskID:          0,
			DurationMinutes: minutes,
			StartAt:         sess.StartTime.Format(time.RFC3339),
			EndAt:           sess.EndTime.Format(time.RFC3339),
			Note:            note,
		})
	}

	daily := make([]dailyStudyStat, 0, len(dailyMap))
	for _, item := range dailyMap {
		daily = append(daily, *item)
	}
	sort.Slice(daily, func(i, j int) bool {
		return daily[i].Date < daily[j].Date
	})

	tasks, stats := loadUserTasks(db, userID)

	taskSamples := make([]taskSample, 0, len(tasks))
	for _, task := range tasks {
		var due *string
		if task.DueAt != nil {
			formatted := task.DueAt.Format("2006-01-02")
			due = &formatted
		}
		taskSamples = append(taskSamples, taskSample{
			ID:       task.ID,
			Title:    strings.TrimSpace(task.Title),
			Status:   taskStatusLabel(task.Status, task.CompletedAt),
			Priority: task.Priority,
			DueDate:  due,
		})
	}

	return analysisContext{
		User: userBriefContext{
			ID:          user.ID,
			DisplayName: user.DisplayName,
		},
		Profile: profileBriefContext{
			StreakDays:     profile.StreakDays,
			TaskDone:       profile.TasksCompleted,
			TaskInProgress: profile.TasksInProgress,
			CompletionRate: profile.TaskCompletionRate,
			TotalStudyMins: profile.TotalStudyMins,
		},
		TimeRange: timeRangeContext{
			Days: days,
			From: start.Format("2006-01-02"),
			To:   time.Now().Format("2006-01-02"),
		},
		Daily:       daily,
		Sessions:    sessionSamples,
		TaskStats:   stats,
		TaskSamples: taskSamples,
	}, nil
}

func loadUserTasks(db *gorm.DB, userID uint64) ([]models.Task, taskStatSnapshot) {
	result := []models.Task{}
	taskIDs := []uint64{}

	var owned []models.Task
	_ = db.Where("owner_user_id = ? OR created_by = ?", userID, userID).Find(&owned).Error
	for _, t := range owned {
		taskIDs = append(taskIDs, t.ID)
		result = append(result, t)
	}

	var assignees []models.TaskAssignee
	_ = db.Where("user_id = ?", userID).Find(&assignees).Error
	for _, a := range assignees {
		taskIDs = append(taskIDs, a.TaskID)
	}

	if len(assignees) > 0 {
		var assigned []models.Task
		_ = db.Where("id IN ?", taskIDs).Find(&assigned).Error
		result = mergeTasks(result, assigned)
	}

	stats := taskStatSnapshot{}
	seen := map[uint64]struct{}{}
	for _, task := range result {
		if _, ok := seen[task.ID]; ok {
			continue
		}
		seen[task.ID] = struct{}{}

		stats.Total++
		statusLabel := strings.ToLower(taskStatusLabel(task.Status, task.CompletedAt))
		switch statusLabel {
		case "done":
			stats.Completed++
		case "in_progress":
			stats.InProgress++
		default:
			stats.Pending++
		}
	}
	stats.CompletionRate = safeRate(stats.Completed, stats.Total)
	return result, stats
}

func mergeTasks(existing, incoming []models.Task) []models.Task {
	index := map[uint64]models.Task{}
	for _, t := range existing {
		index[t.ID] = t
	}
	for _, t := range incoming {
		if _, ok := index[t.ID]; !ok {
			existing = append(existing, t)
			index[t.ID] = t
		}
	}
	return existing
}

func taskStatusLabel(status int8, completedAt *time.Time) string {
	if completedAt != nil {
		return "done"
	}
	switch status {
	case 1:
		return "in_progress"
	case 2:
		return "done"
	default:
		return "pending"
	}
}

func safeRate(done, total int) int {
	if total == 0 {
		return 0
	}
	return int(math.Round(float64(done) / float64(total) * 100))
}

func buildLLMPayload(ctx analysisContext, model string) (llmPayload, string) {
	if model == "" {
		model = "qwen-plus"
	}
	systemPrompt := "你是一名学习效率教练，擅长用数据给出简洁可执行的反馈。用中文回答，避免虚构。"
	contextJSON, _ := json.MarshalIndent(ctx, "", "  ")
	userPrompt := fmt.Sprintf(
		"请基于以下学习记录和任务数据生成结构化分析，完全按照 JSON schema 输出。"+
			" 在输出中同时给出学习效率总结、关键任务洞察以及复习/提醒清单，优先返回高影响、短周期可执行的建议，避免冗长描述。\n数据：\n%s",
		string(contextJSON),
	)

	schema := map[string]interface{}{
		"type":     "object",
		"required": []string{"summary", "study_trend", "recommendations", "task_insights", "review_plan"},
		"properties": map[string]interface{}{
			"summary": map[string]interface{}{
				"type":     "object",
				"required": []string{"weekly_study_hours", "focus_score", "task_completion_rate", "streak_days"},
				"properties": map[string]interface{}{
					"weekly_study_hours":   map[string]interface{}{"type": "number"},
					"focus_score":          map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 100},
					"task_completion_rate": map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 100},
					"streak_days":          map[string]interface{}{"type": "integer", "minimum": 0},
					"key_strengths":        map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "maxItems": 3},
					"risks":                map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "maxItems": 3},
				},
			},
			"review_plan": map[string]interface{}{
				"type":     "object",
				"required": []string{"knowledge_map"},
				"properties": map[string]interface{}{
					"summary": map[string]interface{}{"type": "string"},
					"knowledge_map": map[string]interface{}{
						"type":     "object",
						"required": []string{"mastered", "learning", "to_learn"},
						"properties": map[string]interface{}{
							"mastered": map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 100},
							"learning": map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 100},
							"to_learn": map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 100},
						},
					},
					"review_items": map[string]interface{}{
						"type":     "array",
						"maxItems": 8,
						"items": map[string]interface{}{
							"type":     "object",
							"required": []string{"subject", "priority"},
							"properties": map[string]interface{}{
								"subject":  map[string]interface{}{"type": "string"},
								"priority": map[string]interface{}{"type": "string", "enum": []string{"high", "medium", "low"}},
								"progress": map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 100},
								"due_date": map[string]interface{}{"type": "string"},
							},
						},
					},
					"reminders": map[string]interface{}{
						"type":     "array",
						"maxItems": 5,
						"items": map[string]interface{}{
							"type":     "object",
							"required": []string{"content"},
							"properties": map[string]interface{}{
								"content": map[string]interface{}{"type": "string"},
								"time":    map[string]interface{}{"type": "string"},
							},
						},
					},
				},
			},
			"study_trend": map[string]interface{}{
				"type":     "array",
				"maxItems": 14,
				"items": map[string]interface{}{
					"type":     "object",
					"required": []string{"date", "study_hours", "focus_score"},
					"properties": map[string]interface{}{
						"date":        map[string]interface{}{"type": "string"},
						"study_hours": map[string]interface{}{"type": "number"},
						"focus_score": map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 100},
					},
				},
			},
			"recommendations": map[string]interface{}{
				"type":     "array",
				"maxItems": 5,
				"items": map[string]interface{}{
					"type":     "object",
					"required": []string{"title", "detail", "impact"},
					"properties": map[string]interface{}{
						"title":  map[string]interface{}{"type": "string"},
						"detail": map[string]interface{}{"type": "string"},
						"impact": map[string]interface{}{"type": "string", "enum": []string{"high", "medium", "low"}},
					},
				},
			},
			"task_insights": map[string]interface{}{
				"type":     "array",
				"maxItems": 5,
				"items": map[string]interface{}{
					"type":     "object",
					"required": []string{"task_title", "status", "advice"},
					"properties": map[string]interface{}{
						"task_title": map[string]interface{}{"type": "string"},
						"status":     map[string]interface{}{"type": "string"},
						"advice":     map[string]interface{}{"type": "string"},
						"risk":       map[string]interface{}{"type": "string"},
					},
				},
			},
		},
	}

	return llmPayload{
		Model:        model,
		Temperature:  0.2,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Schema: llmSchemaSpec{
			Name:   "learning_efficiency_report",
			Schema: schema,
		},
	}, userPrompt
}

func runLLMAnalysis(payload llmPayload, prompt string) (efficiencyAnalysisView, error) {
	apiKey := getQwenAPIKey()
	if apiKey == "" {
		return efficiencyAnalysisView{}, fmt.Errorf("未配置 QWEN_API_KEY，当前为空")
	}
	chatURL := qwenChatURL()

	reqBody := struct {
		Model          string       `json:"model"`
		Messages       []llmMessage `json:"messages"`
		Temperature    float64      `json:"temperature,omitempty"`
		ResponseFormat *struct {
			Type       string `json:"type"`
			JsonSchema struct {
				Name   string                 `json:"name"`
				Schema map[string]interface{} `json:"schema"`
				Strict bool                   `json:"strict"`
			} `json:"json_schema,omitempty"`
		} `json:"response_format,omitempty"`
	}{
		Model:       payload.Model,
		Messages:    []llmMessage{{Role: "system", Content: payload.SystemPrompt}, {Role: "user", Content: payload.UserPrompt}},
		Temperature: payload.Temperature,
		ResponseFormat: &struct {
			Type       string `json:"type"`
			JsonSchema struct {
				Name   string                 `json:"name"`
				Schema map[string]interface{} `json:"schema"`
				Strict bool                   `json:"strict"`
			} `json:"json_schema,omitempty"`
		}{
			Type: "json_schema",
			JsonSchema: struct {
				Name   string                 `json:"name"`
				Schema map[string]interface{} `json:"schema"`
				Strict bool                   `json:"strict"`
			}{
				Name:   payload.Schema.Name,
				Schema: payload.Schema.Schema,
				Strict: true,
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, chatURL, bytes.NewReader(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 50 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return efficiencyAnalysisView{}, fmt.Errorf("调用AI接口失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("[analysis] raw AI response: %s", string(body))
	if resp.StatusCode >= 400 {
		msg := fmt.Sprintf("AI接口返回错误状态 %d", resp.StatusCode)
		if strings.Contains(string(body), "insufficient_quota") || resp.StatusCode == http.StatusTooManyRequests {
			return efficiencyAnalysisView{}, aiError{Status: http.StatusTooManyRequests, Message: "AI 配额不足，请更换或充值 QWEN API Key"}
		}
		return efficiencyAnalysisView{}, fmt.Errorf("%s: %s", msg, string(body))
	}

	var chatResp QwenResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return efficiencyAnalysisView{}, fmt.Errorf("解析AI响应失败: %v", err)
	}
	if chatResp.Error != nil {
		if strings.Contains(chatResp.Error.Message, "insufficient_quota") {
			return efficiencyAnalysisView{}, aiError{Status: http.StatusTooManyRequests, Message: "AI 配额不足，请更换或充值 QWEN API Key"}
		}
		return efficiencyAnalysisView{}, aiError{Status: http.StatusBadGateway, Message: "AI接口错误: " + chatResp.Error.Message}
	}
	if len(chatResp.Choices) == 0 {
		return efficiencyAnalysisView{}, fmt.Errorf("AI返回内容为空")
	}

	content := chatResp.Choices[0].Message.Content
	content = extractJSON(content)

	var result llmAnalysis
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		log.Printf("[analysis] unmarshal LLM JSON failed: %v content=%s", err, content)
		// 尝试宽松解析（兼容未按 schema 返回的键名）
		fallback, fbErr := parseLooseAnalysis(content)
		if fbErr != nil {
			return efficiencyAnalysisView{}, fmt.Errorf("解析AI JSON失败: %v；fallback 解析失败: %v", err, fbErr)
		}
		result = fallback
	} else if isEmptyAnalysis(result) {
		// 严格解析虽然成功但内容为空，尝试宽松解析
		if fallback, fbErr := parseLooseAnalysis(content); fbErr == nil && !isEmptyAnalysis(fallback) {
			result = fallback
		}
	}

	return mapLLMToView(result, payload.Model, prompt), nil
}

func mapLLMToView(analysis llmAnalysis, modelUsed, prompt string) efficiencyAnalysisView {
	trend := make([]trendPointView, 0, len(analysis.StudyTrend))
	for _, item := range analysis.StudyTrend {
		trend = append(trend, trendPointView{
			Date:       item.Date,
			StudyHours: round1(item.StudyHours),
			FocusScore: clampInt(item.FocusScore, 0, 100),
		})
	}

	recs := make([]recommendationView, 0, len(analysis.Recommendations))
	for _, rec := range analysis.Recommendations {
		recs = append(recs, recommendationView{
			Title:  rec.Title,
			Detail: rec.Detail,
			Impact: rec.Impact,
		})
	}

	taskInsights := make([]taskInsightView, 0, len(analysis.TaskInsights))
	for _, task := range analysis.TaskInsights {
		taskInsights = append(taskInsights, taskInsightView{
			TaskTitle: task.TaskTitle,
			Status:    task.Status,
			Advice:    task.Advice,
			Risk:      task.Risk,
		})
	}

	reviewPlan := reviewPlanView{
		Summary: analysis.ReviewPlan.Summary,
		KnowledgeMap: knowledgeMapView{
			Mastered: clampInt(analysis.ReviewPlan.KnowledgeMap.Mastered, 0, 100),
			Learning: clampInt(analysis.ReviewPlan.KnowledgeMap.Learning, 0, 100),
			ToLearn:  clampInt(analysis.ReviewPlan.KnowledgeMap.ToLearn, 0, 100),
		},
	}

	for _, item := range analysis.ReviewPlan.ReviewItems {
		reviewPlan.ReviewItems = append(reviewPlan.ReviewItems, reviewItemView{
			Subject:  item.Subject,
			Priority: item.Priority,
			Progress: clampInt(item.Progress, 0, 100),
			DueDate:  item.DueDate,
		})
	}
	for _, reminder := range analysis.ReviewPlan.Reminders {
		reviewPlan.Reminders = append(reviewPlan.Reminders, reminderView{
			Content: reminder.Content,
			Time:    reminder.Time,
		})
	}

	return efficiencyAnalysisView{
		Source:      "llm",
		Model:       modelUsed,
		GeneratedAt: time.Now().Format(time.RFC3339),
		Prompt:      prompt,
		Summary: analysisSummaryView{
			WeeklyStudyHours:   round1(analysis.Summary.WeeklyStudyHours),
			FocusScore:         clampInt(analysis.Summary.FocusScore, 0, 100),
			TaskCompletionRate: clampInt(analysis.Summary.TaskCompletionRate, 0, 100),
			StreakDays:         maxInt(analysis.Summary.StreakDays, 0),
			StrengthHighlights: analysis.Summary.KeyStrengths,
			RiskHighlights:     analysis.Summary.Risks,
		},
		StudyTrend:      trend,
		Recommendations: recs,
		TaskInsights:    taskInsights,
		ReviewPlan:      reviewPlan,
	}
}

func round1(val float64) float64 {
	return math.Round(val*10) / 10
}

func clampInt(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isEmptyAnalysis(a llmAnalysis) bool {
	return a.Summary.WeeklyStudyHours == 0 &&
		a.Summary.FocusScore == 0 &&
		a.Summary.TaskCompletionRate == 0 &&
		a.Summary.StreakDays == 0 &&
		len(a.Recommendations) == 0 &&
		len(a.TaskInsights) == 0 &&
		len(a.StudyTrend) == 0 &&
		len(a.ReviewPlan.ReviewItems) == 0 &&
		len(a.ReviewPlan.Reminders) == 0
}

func isContextEmpty(ctx analysisContext) bool {
	return len(ctx.Daily) == 0 && len(ctx.Sessions) == 0 && len(ctx.TaskSamples) == 0
}

// 兼容宽松 JSON 结构的解析
func parseLooseAnalysis(content string) (llmAnalysis, error) {
	var loose looseAnalysis
	if err := json.Unmarshal([]byte(content), &loose); err != nil {
		return llmAnalysis{}, err
	}

	// 粗略判断是否包含有效信息
	if loose.LearningEfficiencySummary.CompletionRate == 0 &&
		loose.LearningEfficiencySummary.TotalStudyMinutes == 0 &&
		len(loose.KeyTaskInsights) == 0 &&
		len(loose.ReviewReminderList) == 0 &&
		len(loose.ReviewAndReminders) == 0 {
		return llmAnalysis{}, fmt.Errorf("缺少预期字段")
	}

	sum := loose.LearningEfficiencySummary
	result := llmAnalysis{}

	weeklyHours := math.Round((float64(sum.TotalStudyMinutes)/60.0)*10) / 10
	result.Summary.WeeklyStudyHours = weeklyHours
	result.Summary.FocusScore = 80 // 模型未返回时给默认值
	result.Summary.TaskCompletionRate = clampInt(sum.CompletionRate, 0, 100)
	result.Summary.StreakDays = sum.StreakDays
	if sum.Insight != "" {
		result.Summary.KeyStrengths = append(result.Summary.KeyStrengths, sum.Insight)
	}

	for _, item := range loose.KeyTaskInsights {
		evidenceText := ""
		switch v := item.Evidence.(type) {
		case string:
			evidenceText = v
		case []interface{}:
			parts := []string{}
			for _, e := range v {
				if s, ok := e.(string); ok {
					parts = append(parts, s)
				}
			}
			evidenceText = strings.Join(parts, "；")
		}
		result.TaskInsights = append(result.TaskInsights, struct {
			TaskTitle string `json:"task_title"`
			Status    string `json:"status"`
			Advice    string `json:"advice"`
			Risk      string `json:"risk"`
		}{
			TaskTitle: item.Issue,
			Status:    "in_progress",
			Advice:    evidenceText,
			Risk:      item.Impact,
		})
	}

	// 利用复习/提醒列表构造建议与复习计划
	for _, reminder := range loose.ReviewReminderList {
		impact := "medium"
		switch reminder.Priority {
		case 1:
			impact = "high"
		case 2:
			impact = "medium"
		case 3:
			impact = "low"
		}
		result.Recommendations = append(result.Recommendations, struct {
			Title  string `json:"title"`
			Detail string `json:"detail"`
			Impact string `json:"impact"`
		}{
			Title:  reminder.Action,
			Detail: reminder.Suggestion,
			Impact: impact,
		})

		priorityLabel := map[int]string{1: "high", 2: "medium", 3: "low"}[reminder.Priority]
		if priorityLabel == "" {
			priorityLabel = "medium"
		}
		result.ReviewPlan.ReviewItems = append(result.ReviewPlan.ReviewItems, struct {
			Subject  string `json:"subject"`
			Priority string `json:"priority"`
			Progress int    `json:"progress"`
			DueDate  string `json:"due_date"`
		}{
			Subject:  reminder.Title,
			Priority: priorityLabel,
			Progress: 0,
			DueDate:  reminder.DueDate,
		})
		result.ReviewPlan.Reminders = append(result.ReviewPlan.Reminders, struct {
			Content string `json:"content"`
			Time    string `json:"time"`
		}{
			Content: reminder.Action,
			Time:    reminder.DueDate,
		})
	}

	for _, item := range loose.ReviewAndReminders {
		priority := strings.ToLower(item.Priority)
		if priority == "" {
			priority = "medium"
		}
		impact := map[string]string{
			"high":   "high",
			"medium": "medium",
			"low":    "low",
		}[priority]
		if impact == "" {
			impact = "medium"
		}
		result.Recommendations = append(result.Recommendations, struct {
			Title  string `json:"title"`
			Detail string `json:"detail"`
			Impact string `json:"impact"`
		}{
			Title:  item.Title,
			Detail: item.Detail,
			Impact: impact,
		})
		result.ReviewPlan.Reminders = append(result.ReviewPlan.Reminders, struct {
			Content string `json:"content"`
			Time    string `json:"time"`
		}{
			Content: item.Title,
			Time:    item.Deadline,
		})
		result.ReviewPlan.ReviewItems = append(result.ReviewPlan.ReviewItems, struct {
			Subject  string `json:"subject"`
			Priority string `json:"priority"`
			Progress int    `json:"progress"`
			DueDate  string `json:"due_date"`
		}{
			Subject:  item.Title,
			Priority: priority,
			Progress: 0,
			DueDate:  item.Deadline,
		})
	}

	// 粗略构造知识图谱百分比
	mastered := clampInt(sum.CompletionRate, 0, 100)
	learning := clampInt(int(sum.InProgressRatio*100), 0, 100)
	if mastered+learning > 100 {
		learning = maxInt(0, 100-mastered)
	}
	toLearn := maxInt(0, 100-mastered-learning)
	result.ReviewPlan.KnowledgeMap.Mastered = mastered
	result.ReviewPlan.KnowledgeMap.Learning = learning
	result.ReviewPlan.KnowledgeMap.ToLearn = toLearn

	return result, nil
}

// seedUserMonthStudyData 为指定用户生成最近一段时间的学习记录（仅供联调/测试）
func seedUserMonthStudyData(c *gin.Context) {
	type seedRequest struct {
		UserID      uint64 `json:"user_id"`
		Days        int    `json:"days"`
		AvgMinutes  int    `json:"avg_minutes"`  // 每日平均学习分钟
		VarianceMin int    `json:"variance_min"` // 振幅
	}
	var req seedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.UserID == 0 {
		req.UserID = 1
	}
	if req.Days <= 0 {
		req.Days = 30
	}
	if req.AvgMinutes <= 0 {
		req.AvgMinutes = 90
	}
	if req.VarianceMin < 0 {
		req.VarianceMin = 30
	}

	db := database.GetDB()
	var user models.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户不存在"})
		return
	}

	now := time.Now()
	records := make([]models.LearningRecord, 0, req.Days*2)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < req.Days; i++ {
		day := now.AddDate(0, 0, -i)

		// 生成 1-2 个学习 session，合计时长在 avg ± variance
		base := req.AvgMinutes
		offset := rand.Intn(req.VarianceMin*2+1) - req.VarianceMin
		totalMinutes := maxInt(30, base+offset)
		sessionCount := 1
		if totalMinutes >= 90 {
			sessionCount = 2
		}
		remain := totalMinutes
		for s := 0; s < sessionCount; s++ {
			slotMinutes := remain
			if s == 0 && sessionCount == 2 {
				// 第一段占 40%-60%
				slotMinutes = maxInt(25, int(float64(totalMinutes)*0.4)+rand.Intn(20))
				if slotMinutes >= totalMinutes {
					slotMinutes = totalMinutes / 2
				}
			}
			if slotMinutes <= 0 {
				continue
			}
			start := time.Date(day.Year(), day.Month(), day.Day(), 9+rand.Intn(8), rand.Intn(60), 0, 0, time.Local)
			start = start.Add(time.Duration(s) * time.Hour * 3) // sessions 间隔
			end := start.Add(time.Duration(slotMinutes) * time.Minute)
			focusScore := clampInt(55+slotMinutes/3+rand.Intn(15), 50, 95)
			records = append(records, models.LearningRecord{
				TaskID:          0,
				UserID:          req.UserID,
				SessionStart:    start,
				SessionEnd:      end,
				DurationMinutes: slotMinutes,
				Note:            fmt.Sprintf("自动生成学习记录，专注度评分:%d/100", focusScore),
			})
			remain -= slotMinutes
		}
	}

	if len(records) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "未生成任何记录"})
		return
	}

	if err := db.Create(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "写入学习记录失败", "error": err.Error()})
		return
	}

	// 更新用户档案统计
	totalAdd := 0
	for _, r := range records {
		totalAdd += r.DurationMinutes
	}
	db.Model(&models.UserProfile{}).
		Where("user_id = ?", req.UserID).
		UpdateColumn("total_study_mins", gorm.Expr("COALESCE(total_study_mins,0)+?", totalAdd))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data": gin.H{
			"user_id":       req.UserID,
			"records_added": len(records),
			"minutes_added": totalAdd,
			"days":          req.Days,
		},
	})
}
