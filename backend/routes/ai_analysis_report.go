package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"learningAssistant-backend/database"
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/models"
)

// registerAIAnalysisReportRoutes 注册 AI 智能分析报告路由
func registerAIAnalysisReportRoutes(router *gin.RouterGroup) {
	router.POST("/ai-report", middleware.AuthMiddleware(), handleGenerateAIReport)
	router.GET("/ai-report/history", middleware.AuthMiddleware(), handleGetReportHistory)
}

// AI 学习报告请求结构
type AIReportRequest struct {
	Days       int    `json:"days"`        // 分析天数，默认7天
	ReportType string `json:"report_type"` // 报告类型：weekly/monthly/custom
}

// AI 学习报告结构
type AILearningReport struct {
	// 基础信息
	GeneratedAt string `json:"generated_at"`
	ReportType  string `json:"report_type"`
	Period      string `json:"period"`

	// 学习概览
	Overview LearningOverview `json:"overview"`

	// 能力雷达图数据
	AbilityRadar AbilityRadarData `json:"ability_radar"`

	// 学习行为分析
	BehaviorAnalysis BehaviorAnalysis `json:"behavior_analysis"`

	// 知识掌握分析
	KnowledgeAnalysis KnowledgeAnalysis `json:"knowledge_analysis"`

	// AI 个性化建议
	AIAdvice AIAdviceSection `json:"ai_advice"`

	// 下周学习计划
	WeeklyPlan WeeklyPlanSection `json:"weekly_plan"`

	// 激励语句
	Motivation string `json:"motivation"`
}

// LearningOverview 学习概览
type LearningOverview struct {
	TotalStudyHours     float64 `json:"total_study_hours"`
	TotalStudyDays      int     `json:"total_study_days"`
	TasksCompleted      int     `json:"tasks_completed"`
	TaskCompletionRate  float64 `json:"task_completion_rate"`
	KnowledgePoints     int     `json:"knowledge_points"`
	NotesCreated        int     `json:"notes_created"`
	StreakDays          int     `json:"streak_days"`
	ComparedToLastWeek  string  `json:"compared_to_last_week"` // 与上周对比
	EfficiencyScore     int     `json:"efficiency_score"`      // 学习效率评分 0-100
	EfficiencyLevel     string  `json:"efficiency_level"`      // 效率等级描述
}

// AbilityRadarData 能力雷达图数据
type AbilityRadarData struct {
	Dimensions []RadarDimension `json:"dimensions"`
}

// RadarDimension 雷达图维度
type RadarDimension struct {
	Name  string `json:"name"`
	Value int    `json:"value"` // 0-100
	Max   int    `json:"max"`
}

// BehaviorAnalysis 学习行为分析
type BehaviorAnalysis struct {
	PeakStudyTime       string           `json:"peak_study_time"`       // 学习高峰时段
	AverageSessionTime  int              `json:"average_session_time"`  // 平均单次学习时长（分钟）
	MostProductiveDay   string           `json:"most_productive_day"`   // 最高效的日子
	StudyHabits         []string         `json:"study_habits"`          // 学习习惯总结
	FocusDistribution   []FocusItem      `json:"focus_distribution"`    // 专注度分布
	SubjectDistribution []SubjectItem    `json:"subject_distribution"`  // 学科分布
	DailyTrend          []DailyTrendItem `json:"daily_trend"`           // 每日趋势
}

// FocusItem 专注度项
type FocusItem struct {
	Level      string  `json:"level"` // 高度专注/中等专注/低专注
	Percentage float64 `json:"percentage"`
}

// SubjectItem 学科分布项
type SubjectItem struct {
	Subject    string  `json:"subject"`
	Hours      float64 `json:"hours"`
	Percentage float64 `json:"percentage"`
	Color      string  `json:"color"`
}

// DailyTrendItem 每日趋势项
type DailyTrendItem struct {
	Date        string  `json:"date"`
	StudyHours  float64 `json:"study_hours"`
	TasksDone   int     `json:"tasks_done"`
	FocusScore  int     `json:"focus_score"`
}

// KnowledgeAnalysis 知识掌握分析
type KnowledgeAnalysis struct {
	TotalKnowledge   int                     `json:"total_knowledge"`
	MasteredCount    int                     `json:"mastered_count"`
	LearningCount    int                     `json:"learning_count"`
	ToLearnCount     int                     `json:"to_learn_count"`
	MasteryRate      float64                 `json:"mastery_rate"`
	TopCategories    []CategoryMastery       `json:"top_categories"`
	RecentProgress   []KnowledgeProgressItem `json:"recent_progress"`
	WeakPoints       []string                `json:"weak_points"` // AI 分析的薄弱点
}

// CategoryMastery 分类掌握度
type CategoryMastery struct {
	Category     string  `json:"category"`
	MasteryRate  float64 `json:"mastery_rate"`
	TotalCount   int     `json:"total_count"`
	MasteredCount int    `json:"mastered_count"`
}

// KnowledgeProgressItem 知识进度项
type KnowledgeProgressItem struct {
	Date       string `json:"date"`
	NewLearned int    `json:"new_learned"`
	Reviewed   int    `json:"reviewed"`
}

// AIAdviceSection AI 建议部分
type AIAdviceSection struct {
	StrengthAnalysis   []string        `json:"strength_analysis"`   // 优势分析
	ImprovementAreas   []string        `json:"improvement_areas"`   // 待提升领域
	PersonalizedTips   []PersonalTip   `json:"personalized_tips"`   // 个性化建议
	RecommendedActions []ActionItem    `json:"recommended_actions"` // 推荐行动
}

// PersonalTip 个性化建议
type PersonalTip struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"` // high/medium/low
	Icon        string `json:"icon"`
}

// ActionItem 行动项
type ActionItem struct {
	Action     string `json:"action"`
	Reason     string `json:"reason"`
	Impact     string `json:"impact"`
	Difficulty string `json:"difficulty"` // easy/medium/hard
}

// WeeklyPlanSection 周计划部分
type WeeklyPlanSection struct {
	Goals     []WeeklyGoal `json:"goals"`
	Schedule  []DayPlan    `json:"schedule"`
	KeyTasks  []string     `json:"key_tasks"`
}

// WeeklyGoal 周目标
type WeeklyGoal struct {
	Goal     string `json:"goal"`
	Priority int    `json:"priority"`
	Metric   string `json:"metric"`
}

// DayPlan 日计划
type DayPlan struct {
	Day           string   `json:"day"`
	FocusSubjects []string `json:"focus_subjects"`
	SuggestedTime int      `json:"suggested_time"` // 建议学习时长（分钟）
}

// handleGenerateAIReport 生成 AI 学习报告
func handleGenerateAIReport(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req AIReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Days = 7
		req.ReportType = "weekly"
	}
	if req.Days <= 0 {
		req.Days = 7
	}
	if req.Days > 30 {
		req.Days = 30
	}
	if req.ReportType == "" {
		req.ReportType = "weekly"
	}

	// 收集用户数据
	reportData, err := collectReportData(userID.(uint64), req.Days)
	if err != nil {
		log.Printf("[ai-report] collect data error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "收集数据失败"})
		return
	}

	// 调用 AI 生成分析
	report, err := generateAIReportWithLLM(reportData, req)
	if err != nil {
		log.Printf("[ai-report] AI generation error: %v", err)
		// 降级到基础报告
		report = generateBasicReport(reportData, req)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    report,
	})
}

// ReportRawData 收集的原始数据
type ReportRawData struct {
	User           models.User
	Profile        models.UserProfile
	DailyStats     []models.DailyStudyStat
	Tasks          []models.Task
	Notes          []models.StudyNote
	KnowledgeStats map[string]interface{}
	Sessions       []models.StudySession
}

// collectReportData 收集报告数据
func collectReportData(userID uint64, days int) (*ReportRawData, error) {
	db := database.GetDB()
	data := &ReportRawData{}

	// 获取用户信息
	if err := db.First(&data.User, userID).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 获取用户资料
	db.Where("user_id = ?", userID).First(&data.Profile)

	// 获取每日学习统计
	startDate := time.Now().AddDate(0, 0, -days)
	db.Where("user_id = ? AND date >= ?", userID, startDate).
		Order("date ASC").
		Find(&data.DailyStats)

	// 获取任务
	db.Where("user_id = ? AND created_at >= ?", userID, startDate).
		Find(&data.Tasks)

	// 获取笔记
	db.Where("user_id = ? AND created_at >= ?", userID, startDate).
		Find(&data.Notes)

	// 获取学习会话
	db.Where("user_id = ? AND start_at >= ?", userID, startDate).
		Find(&data.Sessions)

	// 获取知识库统计
	if ragService != nil {
		data.KnowledgeStats, _ = ragService.GetUserKnowledgeStats(userID)
	}

	return data, nil
}

// generateAIReportWithLLM 使用大模型生成报告
func generateAIReportWithLLM(data *ReportRawData, req AIReportRequest) (*AILearningReport, error) {
	apiKey := getQwenAPIKey()
	if apiKey == "" {
		return nil, fmt.Errorf("AI 服务未配置")
	}

	// 构建上下文
	context := buildReportContext(data, req)

	prompt := fmt.Sprintf(`你是一位专业的学习分析师，请根据以下学习数据为用户生成一份详细的学习分析报告。

=== 用户学习数据 ===
%s

请严格按照以下 JSON 格式返回分析报告（只返回 JSON，不要 markdown 代码块，不要额外解释）：
{
  "strength_analysis": ["优势1", "优势2", "优势3"],
  "improvement_areas": ["待提升领域1", "待提升领域2"],
  "personalized_tips": [
    {"title": "建议标题", "description": "详细描述", "priority": "high/medium/low", "icon": "💡"}
  ],
  "recommended_actions": [
    {"action": "具体行动", "reason": "原因", "impact": "预期效果", "difficulty": "easy/medium/hard"}
  ],
  "weekly_goals": [
    {"goal": "目标描述", "priority": 1, "metric": "衡量指标"}
  ],
  "weak_points": ["薄弱点1", "薄弱点2"],
  "motivation": "一句激励性的话语",
  "peak_study_time": "学习高峰时段描述",
  "study_habits": ["学习习惯1", "学习习惯2"]
}

请基于数据进行专业、具体、可操作的分析，每条建议要有针对性。`, context)

	reqBody := QwenRequest{
		Model: "qwen-plus",
		Messages: []QwenMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	apiURL := qwenChatURL()

	httpReq, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求 AI 失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("[ai-report] AI response: %s", string(body))

	var qwenResp QwenResponse
	if err := json.Unmarshal(body, &qwenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败")
	}

	if qwenResp.Error != nil {
		return nil, fmt.Errorf("AI 错误: %s", qwenResp.Error.Message)
	}

	if len(qwenResp.Choices) == 0 {
		return nil, fmt.Errorf("AI 返回为空")
	}

	content := qwenResp.Choices[0].Message.Content
	content = extractJSON(content)

	// 解析 AI 返回
	var aiResult struct {
		StrengthAnalysis   []string        `json:"strength_analysis"`
		ImprovementAreas   []string        `json:"improvement_areas"`
		PersonalizedTips   []PersonalTip   `json:"personalized_tips"`
		RecommendedActions []ActionItem    `json:"recommended_actions"`
		WeeklyGoals        []WeeklyGoal    `json:"weekly_goals"`
		WeakPoints         []string        `json:"weak_points"`
		Motivation         string          `json:"motivation"`
		PeakStudyTime      string          `json:"peak_study_time"`
		StudyHabits        []string        `json:"study_habits"`
	}

	if err := json.Unmarshal([]byte(content), &aiResult); err != nil {
		log.Printf("[ai-report] parse AI result error: %v, content: %s", err, content)
		return nil, fmt.Errorf("解析 AI 结果失败")
	}

	// 合并基础数据和 AI 分析
	report := generateBasicReport(data, req)
	report.AIAdvice.StrengthAnalysis = aiResult.StrengthAnalysis
	report.AIAdvice.ImprovementAreas = aiResult.ImprovementAreas
	report.AIAdvice.PersonalizedTips = aiResult.PersonalizedTips
	report.AIAdvice.RecommendedActions = aiResult.RecommendedActions
	report.WeeklyPlan.Goals = aiResult.WeeklyGoals
	report.KnowledgeAnalysis.WeakPoints = aiResult.WeakPoints
	report.Motivation = aiResult.Motivation
	report.BehaviorAnalysis.PeakStudyTime = aiResult.PeakStudyTime
	report.BehaviorAnalysis.StudyHabits = aiResult.StudyHabits

	return report, nil
}

// buildReportContext 构建报告上下文
func buildReportContext(data *ReportRawData, req AIReportRequest) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("用户: %s\n", data.User.DisplayName))
	sb.WriteString(fmt.Sprintf("分析周期: 最近 %d 天\n\n", req.Days))

	// 学习时长统计
	totalMinutes := 0
	studyDays := 0
	for _, stat := range data.DailyStats {
		totalMinutes += stat.Minutes
		if stat.Minutes > 0 {
			studyDays++
		}
	}
	sb.WriteString(fmt.Sprintf("总学习时长: %.1f 小时\n", float64(totalMinutes)/60))
	sb.WriteString(fmt.Sprintf("有效学习天数: %d 天\n", studyDays))

	// 任务统计
	totalTasks := len(data.Tasks)
	completedTasks := 0
	for _, task := range data.Tasks {
		if task.Status == 2 {
			completedTasks++
		}
	}
	sb.WriteString(fmt.Sprintf("任务总数: %d, 已完成: %d\n", totalTasks, completedTasks))
	if totalTasks > 0 {
		sb.WriteString(fmt.Sprintf("任务完成率: %.1f%%\n", float64(completedTasks)/float64(totalTasks)*100))
	}

	// 笔记统计
	sb.WriteString(fmt.Sprintf("创建笔记数: %d\n", len(data.Notes)))

	// 知识库统计
	if data.KnowledgeStats != nil {
		sb.WriteString(fmt.Sprintf("知识点总数: %v\n", data.KnowledgeStats["total_count"]))
		sb.WriteString(fmt.Sprintf("已掌握: %v\n", data.KnowledgeStats["level_3_count"]))
		sb.WriteString(fmt.Sprintf("学习中: %v\n", data.KnowledgeStats["level_1_count"]))
	}

	// 连续学习天数
	sb.WriteString(fmt.Sprintf("连续学习天数: %d\n", data.Profile.StreakDays))

	// 任务类别分布
	categoryCount := make(map[string]int)
	for _, task := range data.Tasks {
		catName := "未分类"
		if task.Category != nil {
			catName = task.Category.Name
		}
		categoryCount[catName]++
	}
	sb.WriteString("\n任务类别分布:\n")
	for cat, count := range categoryCount {
		sb.WriteString(fmt.Sprintf("  - %s: %d 个\n", cat, count))
	}

	return sb.String()
}

// generateBasicReport 生成基础报告（不依赖 AI）
func generateBasicReport(data *ReportRawData, req AIReportRequest) *AILearningReport {
	report := &AILearningReport{
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
		ReportType:  req.ReportType,
		Period:      fmt.Sprintf("最近 %d 天", req.Days),
	}

	// 计算概览数据
	totalMinutes := 0
	studyDays := 0
	for _, stat := range data.DailyStats {
		totalMinutes += stat.Minutes
		if stat.Minutes > 0 {
			studyDays++
		}
	}

	completedTasks := 0
	for _, task := range data.Tasks {
		if task.Status == 2 {
			completedTasks++
		}
	}

	totalTasks := len(data.Tasks)
	completionRate := 0.0
	if totalTasks > 0 {
		completionRate = float64(completedTasks) / float64(totalTasks) * 100
	}

	// 计算效率评分
	efficiencyScore := calculateEfficiencyScore(data, req.Days)

	report.Overview = LearningOverview{
		TotalStudyHours:    float64(totalMinutes) / 60,
		TotalStudyDays:     studyDays,
		TasksCompleted:     completedTasks,
		TaskCompletionRate: completionRate,
		KnowledgePoints:    getIntFromStats(data.KnowledgeStats, "total_count"),
		NotesCreated:       len(data.Notes),
		StreakDays:         data.Profile.StreakDays,
		EfficiencyScore:    efficiencyScore,
		EfficiencyLevel:    getEfficiencyLevel(efficiencyScore),
	}

	// 能力雷达图
	report.AbilityRadar = AbilityRadarData{
		Dimensions: []RadarDimension{
			{Name: "学习时长", Value: min(int(report.Overview.TotalStudyHours*10), 100), Max: 100},
			{Name: "任务完成", Value: int(completionRate), Max: 100},
			{Name: "知识掌握", Value: calculateMasteryScore(data), Max: 100},
			{Name: "学习连续性", Value: min(data.Profile.StreakDays*10, 100), Max: 100},
			{Name: "笔记输出", Value: min(len(data.Notes)*20, 100), Max: 100},
		},
	}

	// 行为分析
	report.BehaviorAnalysis = BehaviorAnalysis{
		PeakStudyTime:      "19:00-21:00",
		AverageSessionTime: calculateAverageSessionTime(data.Sessions),
		MostProductiveDay:  findMostProductiveDay(data.DailyStats),
		StudyHabits:        []string{"保持规律学习", "善用碎片时间"},
		SubjectDistribution: calculateSubjectDistribution(data.Tasks),
		DailyTrend:         calculateDailyTrend(data.DailyStats),
	}

	// 知识分析
	report.KnowledgeAnalysis = KnowledgeAnalysis{
		TotalKnowledge: getIntFromStats(data.KnowledgeStats, "total_count"),
		MasteredCount:  getIntFromStats(data.KnowledgeStats, "level_3_count"),
		LearningCount:  getIntFromStats(data.KnowledgeStats, "level_1_count") + getIntFromStats(data.KnowledgeStats, "level_2_count"),
		ToLearnCount:   getIntFromStats(data.KnowledgeStats, "level_0_count"),
		TopCategories:  calculateTopCategories(data.Tasks),
	}

	if report.KnowledgeAnalysis.TotalKnowledge > 0 {
		report.KnowledgeAnalysis.MasteryRate = float64(report.KnowledgeAnalysis.MasteredCount) / float64(report.KnowledgeAnalysis.TotalKnowledge) * 100
	}

	// 默认建议
	report.AIAdvice = AIAdviceSection{
		StrengthAnalysis: []string{"保持了良好的学习习惯", "任务完成情况稳定"},
		ImprovementAreas: []string{"可以增加复习频率", "尝试更多实践练习"},
		PersonalizedTips: []PersonalTip{
			{Title: "制定每日小目标", Description: "将大任务拆分成小目标，更容易坚持", Priority: "high", Icon: "🎯"},
			{Title: "定期复习", Description: "使用艾宾浩斯遗忘曲线安排复习计划", Priority: "medium", Icon: "📚"},
		},
		RecommendedActions: []ActionItem{
			{Action: "每天固定时间学习30分钟", Reason: "形成习惯", Impact: "提升学习效率20%", Difficulty: "easy"},
		},
	}

	// 周计划
	report.WeeklyPlan = WeeklyPlanSection{
		Goals: []WeeklyGoal{
			{Goal: "完成本周所有任务", Priority: 1, Metric: "任务完成率100%"},
			{Goal: "每天至少学习1小时", Priority: 2, Metric: "累计7小时"},
		},
		KeyTasks: getKeyTasks(data.Tasks),
	}

	report.Motivation = "每一次努力都是在为更好的自己铺路！💪"

	return report
}

// 辅助函数
func calculateEfficiencyScore(data *ReportRawData, days int) int {
	score := 0

	// 学习时长分 (30分)
	totalMinutes := 0
	for _, stat := range data.DailyStats {
		totalMinutes += stat.Minutes
	}
	avgMinutes := float64(totalMinutes) / float64(days)
	score += min(int(avgMinutes/60*10), 30)

	// 任务完成率分 (30分)
	completedTasks := 0
	for _, task := range data.Tasks {
		if task.Status == 2 {
			completedTasks++
		}
	}
	if len(data.Tasks) > 0 {
		score += int(float64(completedTasks) / float64(len(data.Tasks)) * 30)
	}

	// 连续学习分 (20分)
	score += min(data.Profile.StreakDays*2, 20)

	// 笔记输出分 (20分)
	score += min(len(data.Notes)*4, 20)

	return min(score, 100)
}

func getEfficiencyLevel(score int) string {
	if score >= 90 {
		return "卓越 🌟"
	} else if score >= 75 {
		return "优秀 👍"
	} else if score >= 60 {
		return "良好 ✅"
	} else if score >= 40 {
		return "一般 📈"
	}
	return "待提升 💪"
}

func calculateMasteryScore(data *ReportRawData) int {
	total := getIntFromStats(data.KnowledgeStats, "total_count")
	mastered := getIntFromStats(data.KnowledgeStats, "level_3_count")
	if total == 0 {
		return 0
	}
	return int(float64(mastered) / float64(total) * 100)
}

func getIntFromStats(stats map[string]interface{}, key string) int {
	if stats == nil {
		return 0
	}
	if val, ok := stats[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return 0
}

func calculateAverageSessionTime(sessions []models.StudySession) int {
	if len(sessions) == 0 {
		return 30
	}
	total := 0
	count := 0
	for _, s := range sessions {
		if s.DurationMinutes > 0 {
			total += s.DurationMinutes
			count++
		}
	}
	if count == 0 {
		return 30
	}
	return total / count
}

func findMostProductiveDay(stats []models.DailyStudyStat) string {
	if len(stats) == 0 {
		return "暂无数据"
	}
	maxMinutes := 0
	bestDay := ""
	for _, s := range stats {
		if s.Minutes > maxMinutes {
			maxMinutes = s.Minutes
			bestDay = s.Date.Format("01-02")
		}
	}
	return bestDay
}

func calculateSubjectDistribution(tasks []models.Task) []SubjectItem {
	categoryHours := make(map[string]int)
	total := 0
	for _, t := range tasks {
		catName := "未分类"
		if t.Category != nil {
			catName = t.Category.Name
		}
		categoryHours[catName]++
		total++
	}

	colors := map[string]string{
		"数学": "#3b82f6", "物理": "#8b5cf6", "化学": "#10b981",
		"英语": "#ef4444", "语文": "#f59e0b", "计算机": "#06b6d4",
	}

	var result []SubjectItem
	for cat, count := range categoryHours {
		color := colors[cat]
		if color == "" {
			color = "#9ca3af"
		}
		result = append(result, SubjectItem{
			Subject:    cat,
			Hours:      float64(count),
			Percentage: float64(count) / float64(total) * 100,
			Color:      color,
		})
	}
	return result
}

func calculateDailyTrend(stats []models.DailyStudyStat) []DailyTrendItem {
	var result []DailyTrendItem
	for _, s := range stats {
		result = append(result, DailyTrendItem{
			Date:       s.Date.Format("01-02"),
			StudyHours: float64(s.Minutes) / 60,
			TasksDone:  s.SessionCount,
			FocusScore: min(50+s.Minutes/10, 100),
		})
	}
	return result
}

func calculateTopCategories(tasks []models.Task) []CategoryMastery {
	categoryStats := make(map[string]struct {
		total     int
		completed int
	})

	for _, t := range tasks {
		catName := "未分类"
		if t.Category != nil {
			catName = t.Category.Name
		}
		stats := categoryStats[catName]
		stats.total++
		if t.Status == 2 {
			stats.completed++
		}
		categoryStats[catName] = stats
	}

	var result []CategoryMastery
	for cat, stats := range categoryStats {
		rate := 0.0
		if stats.total > 0 {
			rate = float64(stats.completed) / float64(stats.total) * 100
		}
		result = append(result, CategoryMastery{
			Category:      cat,
			TotalCount:    stats.total,
			MasteredCount: stats.completed,
			MasteryRate:   rate,
		})
	}
	return result
}

func getKeyTasks(tasks []models.Task) []string {
	var result []string
	for _, t := range tasks {
		if t.Status != 2 && t.Priority >= 3 {
			result = append(result, t.Title)
			if len(result) >= 5 {
				break
			}
		}
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// handleGetReportHistory 获取报告历史（可选实现）
func handleGetReportHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    []interface{}{},
	})
}
