package rag

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

// AIAnalysisService AI分析服务 - 基于知识库生成分析数据
type AIAnalysisService struct {
	apiKey     string
	httpClient *http.Client
}

// NewAIAnalysisService 创建AI分析服务
func NewAIAnalysisService(apiKey string) *AIAnalysisService {
	if apiKey == "" {
		apiKey = os.Getenv("QWEN_API_KEY")
	}

	return &AIAnalysisService{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// KnowledgeDistribution 知识点分布分析
type KnowledgeDistribution struct {
	Category      string  `json:"category"`
	Count         int     `json:"count"`
	Percentage    float32 `json:"percentage"`
	MasteredCount int     `json:"mastered_count"`
	LearningCount int     `json:"learning_count"`
	Color         string  `json:"color"`
	Gradient      string  `json:"gradient"`   // 渐变色
	Icon          string  `json:"icon"`       // 图标名
	LightBg       string  `json:"light_bg"`   // 浅色背景
	TextColor     string  `json:"text_color"` // 文字颜色
}

// SkillRadarData 技能雷达数据
type SkillRadarData struct {
	Skill    string `json:"skill"`
	Value    int    `json:"value"` // 0-100
	MaxValue int    `json:"max_value"`
	Level    string `json:"level"`    // beginner, intermediate, advanced, expert
	Progress int    `json:"progress"` // 0-100
	Category string `json:"category"`
}

// LearningTrend 学习趋势
// 说明：学习时长口径容易产生歧义（如跨天、挂机等）。
// 这里改为使用“当日完成任务数 / 当日创建笔记数 / 当日新增知识点数”作为趋势指标。
type LearningTrend struct {
	Date         string `json:"date"`
	NewKnowledge int    `json:"new_knowledge"`
	NewNotes     int    `json:"new_notes"`
	DoneTasks    int    `json:"done_tasks"`
}

// AnalysisReport AI分析报告
type AnalysisReport struct {
	UserID                  uint64                  `json:"user_id"`
	GeneratedAt             time.Time               `json:"generated_at"`
	KnowledgeDistribution   []KnowledgeDistribution `json:"knowledge_distribution"`
	SkillRadar              []SkillRadarData        `json:"skill_radar"`
	LearningTrends          []LearningTrend         `json:"learning_trends"`
	LearningInsights        []string                `json:"learning_insights"`
	RecommendedTopics       []string                `json:"recommended_topics"`
	MasteredSkillsCount     int                     `json:"mastered_skills_count"`
	LearningSkillsCount     int                     `json:"learning_skills_count"`
	TotalKnowledgePoints    int                     `json:"total_knowledge_points"`
	EstimatedCompletionDays int                     `json:"estimated_completion_days"`
}

// AnalyzeUserKnowledge 分析用户知识库
func (a *AIAnalysisService) AnalyzeUserKnowledge(userID uint64) (*AnalysisReport, error) {
	db := database.GetDB()

	// 获取用户的所有知识库条目
	var entries []models.KnowledgeBaseEntry
	if err := db.Where("user_id = ? AND status = 1", userID).
		Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("获取知识库条目失败: %w", err)
	}

	// 生成知识点分布
	distribution := a.analyzeKnowledgeDistribution(entries)

	// 生成技能雷达数据
	skillRadar := a.generateSkillRadar(entries)

	// 生成学习趋势
	trends := a.analyzeLearningTrends(userID, entries)

	// 生成学习洞察
	insights := a.generateInsights(entries, distribution)

	// 推荐主题
	recommended := a.recommendTopics(entries)

	// 统计数据
	masteredCount := 0
	learningCount := 0
	for _, entry := range entries {
		if entry.Level == 4 {
			masteredCount++
		} else if entry.Level > 0 {
			learningCount++
		}
	}

	report := &AnalysisReport{
		UserID:                  userID,
		GeneratedAt:             time.Now(),
		KnowledgeDistribution:   distribution,
		SkillRadar:              skillRadar,
		LearningTrends:          trends,
		LearningInsights:        insights,
		RecommendedTopics:       recommended,
		MasteredSkillsCount:     masteredCount,
		LearningSkillsCount:     learningCount,
		TotalKnowledgePoints:    len(entries),
		EstimatedCompletionDays: a.estimateCompletionDays(entries),
	}

	return report, nil
}

const dateLayoutYMD = "2006-01-02"

// analyzeKnowledgeDistribution 分析知识点分布
func (a *AIAnalysisService) analyzeKnowledgeDistribution(entries []models.KnowledgeBaseEntry) []KnowledgeDistribution {
	categoryMap := make(map[string]*KnowledgeDistribution)

	// 分类统计
	for _, entry := range entries {
		category := entry.Category
		if category == "" {
			category = "其他"
		}

		if _, exists := categoryMap[category]; !exists {
			categoryMap[category] = &KnowledgeDistribution{
				Category:      category,
				Count:         0,
				MasteredCount: 0,
				LearningCount: 0,
			}
		}

		categoryMap[category].Count++
		if entry.Level == 4 {
			categoryMap[category].MasteredCount++
		} else if entry.Level > 0 {
			categoryMap[category].LearningCount++
		}
	}

	// 计算百分比并转换为切片，填充显示配置
	totalCount := len(entries)
	if totalCount == 0 {
		totalCount = 1 // 避免除零
	}
	distribution := make([]KnowledgeDistribution, 0, len(categoryMap))
	for _, dist := range categoryMap {
		config := GetDisplayConfigForCategory(dist.Category)
		dist.Percentage = float32(dist.Count) / float32(totalCount) * 100
		dist.Color = config.Color
		dist.Gradient = config.Gradient
		dist.Icon = config.Icon
		dist.LightBg = config.LightBg
		dist.TextColor = config.TextColor
		distribution = append(distribution, *dist)
	}

	// 按数量排序
	sort.Slice(distribution, func(i, j int) bool {
		return distribution[i].Count > distribution[j].Count
	})

	return distribution
}

// generateSkillRadar 生成技能雷达数据 - 基于分类统计
func (a *AIAnalysisService) generateSkillRadar(entries []models.KnowledgeBaseEntry) []SkillRadarData {
	dims := a.getSkillRadarDims()
	bucket := a.buildSkillRadarBucket(dims, entries)
	return a.buildSkillRadarResult(dims, bucket)
}

type skillRadarDim struct {
	Key        string
	Label      string
	Categories []string
}

type skillRadarStat struct {
	totalLevel int
	count      int
	mastered   int
	learning   int
}

func (a *AIAnalysisService) getSkillRadarDims() []skillRadarDim {
	return []skillRadarDim{
		{Key: "chinese", Label: "语文能力", Categories: []string{"语文", "文学", "阅读", "写作"}},
		{Key: "math", Label: "数学能力", Categories: []string{"数学", "高等数学", "线性代数", "概率论"}},
		{Key: "physics", Label: "物理能力", Categories: []string{"物理"}},
		{Key: "programming", Label: "编程能力", Categories: []string{"编程语言", "前端开发", "后端开发", "算法", "数据库", "DevOps", "人工智能", "数据科学"}},
		{Key: "exam", Label: "考试技巧能力", Categories: []string{"考试", "应试", "刷题", "面试"}},
		{Key: "other", Label: "其他能力", Categories: []string{"其他", "通用", "软技能"}},
	}
}

func (a *AIAnalysisService) buildSkillRadarBucket(dims []skillRadarDim, entries []models.KnowledgeBaseEntry) map[string]*skillRadarStat {
	bucket := make(map[string]*skillRadarStat, len(dims))
	for _, d := range dims {
		bucket[d.Key] = &skillRadarStat{}
	}

	categoryToDim := a.buildSkillRadarCategoryToDimMap(dims)

	for _, entry := range entries {
		dimKey := a.resolveSkillRadarDimKey(categoryToDim, &entry)
		st := bucket[dimKey]
		st.totalLevel += int(entry.Level)
		st.count++
		if entry.Level == 4 {
			st.mastered++
		} else if entry.Level > 0 {
			st.learning++
		}
	}

	return bucket
}

func (a *AIAnalysisService) buildSkillRadarCategoryToDimMap(dims []skillRadarDim) map[string]string {
	categoryToDim := make(map[string]string)
	for _, d := range dims {
		for _, c := range d.Categories {
			categoryToDim[a.normalizeSkillCategoryKey(c)] = d.Key
		}
	}
	return categoryToDim
}

func (a *AIAnalysisService) resolveSkillRadarDimKey(categoryToDim map[string]string, entry *models.KnowledgeBaseEntry) string {
	if entry == nil {
		return "other"
	}

	category := strings.TrimSpace(entry.Category)
	if category == "" {
		category = "其他"
	}

	normalizedCategory := a.normalizeSkillCategoryKey(category)
	if dimKey, ok := categoryToDim[normalizedCategory]; ok {
		return dimKey
	}
	if aliasKey, aliasOk := a.skillRadarAliasDimKey(normalizedCategory); aliasOk {
		return aliasKey
	}
	if guessed := a.guessSkillRadarDimKeyByContent(entry); guessed != "" {
		return guessed
	}
	return "other"
}

func (a *AIAnalysisService) normalizeSkillCategoryKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func (a *AIAnalysisService) skillRadarAliasDimKey(categoryLower string) (string, bool) {
	// 这里用“完全等值”的别名，避免误伤（模糊匹配交给 guessSkillRadarDimKeyByContent）
	alias := map[string]string{
		// programming
		"编程":      "programming",
		"程序设计":    "programming",
		"软件开发":    "programming",
		"前端":      "programming",
		"后端":      "programming",
		"算法与数据结构": "programming",
		"数据结构":    "programming",
		"计算机":     "programming",
		"计算机科学":   "programming",
		"计算机系统":   "programming",
		"操作系统":    "programming",
		"网络":      "programming",
		"计算机网络":   "programming",
		"数据库":     "programming",
		"人工智能":    "programming",
		"机器学习":    "programming",
		"深度学习":    "programming",
		"数据科学":    "programming",

		// exam
		"考试技巧": "exam",
		"应试技巧": "exam",
		"刷题":   "exam",
		"面试":   "exam",
		"测验":   "exam",
	}
	v, ok := alias[categoryLower]
	return v, ok
}

func (a *AIAnalysisService) guessSkillRadarDimKeyByContent(entry *models.KnowledgeBaseEntry) string {
	if entry == nil {
		return "other"
	}

	// 1) category 子串关键词
	cat := strings.ToLower(strings.TrimSpace(entry.Category))
	title := strings.ToLower(strings.TrimSpace(entry.Title))

	// 2) tags/keywords（JSON）
	joined := cat + " " + title + " " + a.joinJSONStrings(entry.Tags) + " " + a.joinJSONStrings(entry.Keywords)

	// programming keywords
	programmingKeywords := []string{
		"编程", "程序", "代码", "开发", "前端", "后端", "vue", "react", "node", "go", "golang", "java", "python", "c++", "c#",
		"算法", "数据结构", "数据库", "mysql", "redis", "http", "rest", "api", "socket", "操作系统", "计算机网络", "linux",
		"机器学习", "深度学习", "ai", "llm",
	}
	for _, kw := range programmingKeywords {
		if strings.Contains(joined, strings.ToLower(kw)) {
			return "programming"
		}
	}

	// exam keywords
	examKeywords := []string{"考试", "应试", "刷题", "面试", "题", "真题", "技巧", "解题"}
	for _, kw := range examKeywords {
		if strings.Contains(joined, strings.ToLower(kw)) {
			return "exam"
		}
	}

	return "other"
}

func (a *AIAnalysisService) joinJSONStrings(raw interface{}) string {
	if raw == nil {
		return ""
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return ""
	}
	var arr []string
	if err := json.Unmarshal(bytes, &arr); err != nil {
		return ""
	}
	return strings.Join(arr, " ")
}

func (a *AIAnalysisService) scoreSkillRadar(st *skillRadarStat) int {
	if st == nil || st.count == 0 {
		return 0
	}

	avgLevel := float64(st.totalLevel) / float64(st.count)
	value := int(avgLevel * 25)

	masteredRatio := float64(st.mastered) / float64(st.count)
	value += int(masteredRatio * 20)

	if st.count >= 3 {
		value += 5
	}
	if st.count >= 6 {
		value += 5
	}
	if st.count >= 10 {
		value += 5
	}

	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func (a *AIAnalysisService) buildSkillRadarResult(dims []skillRadarDim, bucket map[string]*skillRadarStat) []SkillRadarData {
	radar := make([]SkillRadarData, 0, len(dims))
	for _, d := range dims {
		value := a.scoreSkillRadar(bucket[d.Key])
		radar = append(radar, SkillRadarData{
			Skill:    d.Label,
			Value:    value,
			MaxValue: 100,
			Category: d.Label,
			Progress: value,
			Level:    a.getLevelLabel(value),
		})
	}
	return radar
}

// TrendGranularity 趋势聚合粒度
// - day: 按天聚合（适合 30 天）
// - week: 按周聚合（适合 90 天，避免点太密）
// - month: 按月聚合（适合 本年度 12 个月）
type TrendGranularity string

const (
	TrendGranularityDay   TrendGranularity = "day"
	TrendGranularityWeek  TrendGranularity = "week"
	TrendGranularityMonth TrendGranularity = "month"
)

// AnalyzeUserLearningTrendsRange 对外暴露：按范围与粒度生成学习趋势
func (a *AIAnalysisService) AnalyzeUserLearningTrendsRange(userID uint64, entries []models.KnowledgeBaseEntry, from, to time.Time, granularity TrendGranularity) []LearningTrend {
	return a.analyzeLearningTrendsRange(userID, entries, from, to, granularity)
}

// analyzeLearningTrendsRange 分析学习趋势（可配置范围与粒度）
// label 规则：
// - day:  YYYY-MM-DD
// - week: YYYY-Www（ISO 周，周一为起始）
// - month: YYYY-MM
func (a *AIAnalysisService) analyzeLearningTrendsRange(userID uint64, entries []models.KnowledgeBaseEntry, from, to time.Time, granularity TrendGranularity) []LearningTrend {
	bucket := make(map[string]*LearningTrend)

	add := func(label string, apply func(t *LearningTrend)) {
		tr, ok := bucket[label]
		if !ok {
			tr = &LearningTrend{Date: label}
			bucket[label] = tr
		}
		apply(tr)
	}

	a.applyNewKnowledgeToBucket(add, entries, from, to, granularity)
	a.applyDoneTasksToBucket(add, userID, from, to, granularity)
	a.applyNewNotesToBucket(add, userID, from, to, granularity)

	// 生成完整桶（补零），保证 90 天/年度不会缺点
	labels := a.trendBucketLabels(from, to, granularity)
	trends := make([]LearningTrend, 0, len(labels))
	for _, label := range labels {
		if t, ok := bucket[label]; ok {
			trends = append(trends, *t)
		} else {
			trends = append(trends, LearningTrend{Date: label})
		}
	}
	return trends
}

func (a *AIAnalysisService) applyNewKnowledgeToBucket(add func(string, func(*LearningTrend)), entries []models.KnowledgeBaseEntry, from, to time.Time, granularity TrendGranularity) {
	for _, entry := range entries {
		if entry.CreatedAt.Before(from) || entry.CreatedAt.After(to) {
			continue
		}
		label := a.trendBucketLabel(entry.CreatedAt, granularity)
		add(label, func(t *LearningTrend) { t.NewKnowledge++ })
	}
}

func (a *AIAnalysisService) applyDoneTasksToBucket(add func(string, func(*LearningTrend)), userID uint64, from, to time.Time, granularity TrendGranularity) {
	db := database.GetDB()
	var tasks []models.Task
	if err := db.Where("(owner_user_id = ? OR created_by = ?) AND deleted_at IS NULL AND status = ? AND completed_at IS NOT NULL AND completed_at >= ? AND completed_at <= ?", userID, userID, int8(2), from, to).
		Select("id, completed_at").Find(&tasks).Error; err != nil {
		return
	}
	for _, task := range tasks {
		if task.CompletedAt == nil {
			continue
		}
		label := a.trendBucketLabel(*task.CompletedAt, granularity)
		add(label, func(t *LearningTrend) { t.DoneTasks++ })
	}
}

func (a *AIAnalysisService) applyNewNotesToBucket(add func(string, func(*LearningTrend)), userID uint64, from, to time.Time, granularity TrendGranularity) {
	db := database.GetDB()
	var notes []models.StudyNote
	if err := db.Where("user_id = ? AND deleted_at IS NULL AND created_at >= ? AND created_at <= ?", userID, from, to).
		Select("id, created_at").Find(&notes).Error; err != nil {
		return
	}
	for _, note := range notes {
		label := a.trendBucketLabel(note.CreatedAt, granularity)
		add(label, func(t *LearningTrend) { t.NewNotes++ })
	}
}

// analyzeLearningTrends 分析学习趋势（兼容旧调用：默认最近30天按天）
func (a *AIAnalysisService) analyzeLearningTrends(userID uint64, entries []models.KnowledgeBaseEntry) []LearningTrend {
	now := time.Now()
	from := now.AddDate(0, 0, -29)
	return a.analyzeLearningTrendsRange(userID, entries, from, now, TrendGranularityDay)
}

// generateInsights 生成学习洞察
func (a *AIAnalysisService) generateInsights(entries []models.KnowledgeBaseEntry, distribution []KnowledgeDistribution) []string {
	insights := make([]string, 0)

	if len(entries) == 0 {
		return append(insights, "你还没有添加任何知识点。从完成任务或创建笔记开始吧！")
	}

	// 统计掌握情况
	masteredCount := 0
	learningCount := 0
	for _, entry := range entries {
		if entry.Level == 4 {
			masteredCount++
		} else if entry.Level > 0 {
			learningCount++
		}
	}

	// 掌握情况洞察
	if masteredCount > 0 {
		percentage := (masteredCount * 100) / len(entries)
		insights = append(insights, fmt.Sprintf("你已掌握 %d%% 的知识点（%d/%d），继续加油！",
			percentage, masteredCount, len(entries)))
	}

	// 最强和最弱的分类
	if len(distribution) > 0 {
		strongest := distribution[0]
		insights = append(insights, fmt.Sprintf("你在「%s」领域的知识点最多，已有 %d 个。",
			strongest.Category, strongest.Count))

		if len(distribution) > 1 {
			weakest := distribution[len(distribution)-1]
			insights = append(insights, fmt.Sprintf("你需要加强「%s」领域的学习，目前只有 %d 个知识点。",
				weakest.Category, weakest.Count))
		}
	}

	// 学习多样性
	if len(distribution) >= 5 {
		insights = append(insights, "你的学习覆盖面很广，涉及多个领域。保持这种多元学习很棒！")
	}

	// 复习建议
	if learningCount > masteredCount {
		insights = append(insights, "你有很多知识点还在学习中，建议定期复习以加深理解。")
	}

	return insights
}

// recommendTopics 推荐学习主题
func (a *AIAnalysisService) recommendTopics(entries []models.KnowledgeBaseEntry) []string {
	recommended := make([]string, 0)

	// 统计各分类的未掌握知识
	categoryMap := make(map[string]int)
	for _, entry := range entries {
		if entry.Level < 4 {
			categoryMap[entry.Category]++
		}
	}

	// 推荐未掌握最多的分类
	type categoryCount struct {
		category string
		count    int
	}
	var counts []categoryCount
	for cat, count := range categoryMap {
		counts = append(counts, categoryCount{cat, count})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].count > counts[j].count
	})

	for i := 0; i < 5 && i < len(counts); i++ {
		recommended = append(recommended, fmt.Sprintf("继续学习「%s」", counts[i].category))
	}

	// 如果推荐少于5个，添加通用建议
	if len(recommended) < 5 {
		suggestions := []string{
			"掌握核心基础概念",
			"实践项目应用",
			"深入学习算法",
			"提高代码质量",
			"探索新技术领域",
		}
		for _, suggestion := range suggestions {
			if len(recommended) >= 5 {
				break
			}
			if !contains(recommended, suggestion) {
				recommended = append(recommended, suggestion)
			}
		}
	}

	return recommended
}

// extractSkills 从知识库条目中提取技能
func (a *AIAnalysisService) extractSkills(entry models.KnowledgeBaseEntry) []string {
	skills := make([]string, 0)

	// 从关键词中提取
	var keywords []string
	if err := json.Unmarshal(entry.Keywords, &keywords); err == nil {
		for _, keyword := range keywords {
			if len(keyword) > 2 {
				skills = append(skills, keyword)
			}
		}
	}

	// 如果关键词不足，从标题和摘要中提取
	if len(skills) == 0 {
		words := strings.Fields(entry.Title + " " + entry.Summary)
		for _, word := range words {
			if len(word) > 3 && !isStopword(word) {
				skills = append(skills, word)
			}
		}
	}

	// 去重并限制数量
	return uniqueStrings(skills)[:min(len(uniqueStrings(skills)), 3)]
}

// 辅助函数

// SubjectDisplayConfig 学科显示配置
type SubjectDisplayConfig struct {
	Color     string `json:"color"`      // 主色调
	Gradient  string `json:"gradient"`   // 渐变色(CSS格式)
	Icon      string `json:"icon"`       // Iconify图标名
	LightBg   string `json:"light_bg"`   // 浅色背景
	TextColor string `json:"text_color"` // 文字颜色
}

// GetSubjectDisplayConfig 获取学科显示配置 - 面向学习助手的完整配置
func GetSubjectDisplayConfig() map[string]SubjectDisplayConfig {
	return map[string]SubjectDisplayConfig{
		// 理科类
		"数学": {
			Color: "#3b82f6", Gradient: "linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)",
			Icon: "mdi:calculator-variant", LightBg: "#eff6ff", TextColor: "#1e40af",
		},
		"物理": {
			Color: "#8b5cf6", Gradient: "linear-gradient(135deg, #8b5cf6 0%, #6d28d9 100%)",
			Icon: "mdi:atom", LightBg: "#f5f3ff", TextColor: "#5b21b6",
		},
		"化学": {
			Color: "#06b6d4", Gradient: "linear-gradient(135deg, #06b6d4 0%, #0891b2 100%)",
			Icon: "mdi:flask-outline", LightBg: "#ecfeff", TextColor: "#0e7490",
		},
		"生物": {
			Color: "#10b981", Gradient: "linear-gradient(135deg, #10b981 0%, #059669 100%)",
			Icon: "mdi:dna", LightBg: "#ecfdf5", TextColor: "#047857",
		},
		// 文科类
		"语文": {
			Color: "#f59e0b", Gradient: "linear-gradient(135deg, #f59e0b 0%, #d97706 100%)",
			Icon: "mdi:book-open-page-variant", LightBg: "#fffbeb", TextColor: "#b45309",
		},
		"英语": {
			Color: "#ec4899", Gradient: "linear-gradient(135deg, #ec4899 0%, #db2777 100%)",
			Icon: "mdi:alphabetical", LightBg: "#fdf2f8", TextColor: "#be185d",
		},
		"历史": {
			Color: "#78350f", Gradient: "linear-gradient(135deg, #92400e 0%, #78350f 100%)",
			Icon: "mdi:castle", LightBg: "#fef3c7", TextColor: "#78350f",
		},
		"地理": {
			Color: "#16a34a", Gradient: "linear-gradient(135deg, #22c55e 0%, #16a34a 100%)",
			Icon: "mdi:earth", LightBg: "#f0fdf4", TextColor: "#15803d",
		},
		"政治": {
			Color: "#dc2626", Gradient: "linear-gradient(135deg, #ef4444 0%, #dc2626 100%)",
			Icon: "mdi:bank", LightBg: "#fef2f2", TextColor: "#b91c1c",
		},
		// 技能类
		"编程": {
			Color: "#0ea5e9", Gradient: "linear-gradient(135deg, #0ea5e9 0%, #0284c7 100%)",
			Icon: "mdi:code-braces", LightBg: "#f0f9ff", TextColor: "#0369a1",
		},
		"计算机": {
			Color: "#6366f1", Gradient: "linear-gradient(135deg, #6366f1 0%, #4f46e5 100%)",
			Icon: "mdi:laptop", LightBg: "#eef2ff", TextColor: "#4338ca",
		},
		"艺术": {
			Color: "#f472b6", Gradient: "linear-gradient(135deg, #f472b6 0%, #ec4899 100%)",
			Icon: "mdi:palette", LightBg: "#fdf2f8", TextColor: "#db2777",
		},
		"音乐": {
			Color: "#a855f7", Gradient: "linear-gradient(135deg, #a855f7 0%, #9333ea 100%)",
			Icon: "mdi:music", LightBg: "#faf5ff", TextColor: "#7e22ce",
		},
		"体育": {
			Color: "#f97316", Gradient: "linear-gradient(135deg, #f97316 0%, #ea580c 100%)",
			Icon: "mdi:basketball", LightBg: "#fff7ed", TextColor: "#c2410c",
		},
		// 通识类
		"学习方法": {
			Color: "#14b8a6", Gradient: "linear-gradient(135deg, #14b8a6 0%, #0d9488 100%)",
			Icon: "mdi:lightbulb-on", LightBg: "#f0fdfa", TextColor: "#0f766e",
		},
		"考试技巧": {
			Color: "#eab308", Gradient: "linear-gradient(135deg, #eab308 0%, #ca8a04 100%)",
			Icon: "mdi:school", LightBg: "#fefce8", TextColor: "#a16207",
		},
		"阅读": {
			Color: "#84cc16", Gradient: "linear-gradient(135deg, #84cc16 0%, #65a30d 100%)",
			Icon: "mdi:book-open-variant", LightBg: "#f7fee7", TextColor: "#4d7c0f",
		},
		"写作": {
			Color: "#0891b2", Gradient: "linear-gradient(135deg, #22d3ee 0%, #0891b2 100%)",
			Icon: "mdi:pencil", LightBg: "#ecfeff", TextColor: "#0e7490",
		},
		"思维训练": {
			Color: "#7c3aed", Gradient: "linear-gradient(135deg, #a78bfa 0%, #7c3aed 100%)",
			Icon: "mdi:brain", LightBg: "#f5f3ff", TextColor: "#6d28d9",
		},
		"项目": {
			Color: "#2563eb", Gradient: "linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)",
			Icon: "mdi:folder-star", LightBg: "#eff6ff", TextColor: "#1d4ed8",
		},
		"笔记": {
			Color: "#ea580c", Gradient: "linear-gradient(135deg, #fb923c 0%, #ea580c 100%)",
			Icon: "mdi:notebook-outline", LightBg: "#fff7ed", TextColor: "#c2410c",
		},
		// 默认
		"其他": {
			Color: "#64748b", Gradient: "linear-gradient(135deg, #94a3b8 0%, #64748b 100%)",
			Icon: "mdi:bookshelf", LightBg: "#f1f5f9", TextColor: "#475569",
		},
		"默认": {
			Color: "#64748b", Gradient: "linear-gradient(135deg, #94a3b8 0%, #64748b 100%)",
			Icon: "mdi:book", LightBg: "#f1f5f9", TextColor: "#475569",
		},
	}
}

// GetDisplayConfigForCategory 根据分类获取显示配置
func GetDisplayConfigForCategory(category string) SubjectDisplayConfig {
	configs := GetSubjectDisplayConfig()

	// 直接匹配
	if config, exists := configs[category]; exists {
		return config
	}

	// 模糊匹配
	categoryLower := strings.ToLower(category)
	keywordMapping := map[string]string{
		"math": "数学", "数学": "数学", "代数": "数学", "几何": "数学", "微积分": "数学",
		"physics": "物理", "物理": "物理", "力学": "物理", "电学": "物理",
		"chemistry": "化学", "化学": "化学",
		"biology": "生物", "生物": "生物", "生命": "生物",
		"chinese": "语文", "语文": "语文", "文学": "语文", "作文": "语文",
		"english": "英语", "英语": "英语", "外语": "英语",
		"history": "历史", "历史": "历史",
		"geography": "地理", "地理": "地理",
		"politics": "政治", "政治": "政治", "思想": "政治",
		"programming": "编程", "编程": "编程", "代码": "编程", "开发": "编程",
		"computer": "计算机", "计算机": "计算机", "电脑": "计算机",
		"art": "艺术", "艺术": "艺术", "美术": "艺术", "绘画": "艺术",
		"music": "音乐", "音乐": "音乐",
		"sports": "体育", "体育": "体育", "运动": "体育",
		"study": "学习方法", "学习": "学习方法", "方法": "学习方法",
		"exam": "考试技巧", "考试": "考试技巧", "测验": "考试技巧",
		"reading": "阅读", "阅读": "阅读",
		"writing": "写作", "写作": "写作",
		"thinking": "思维训练", "思维": "思维训练", "逻辑": "思维训练",
		"project": "项目", "项目": "项目",
		"note": "笔记", "笔记": "笔记",
	}

	for keyword, subject := range keywordMapping {
		if strings.Contains(categoryLower, keyword) {
			return configs[subject]
		}
	}

	// 默认配置
	return configs["默认"]
}

func (a *AIAnalysisService) getCategoryColor(category string) string {
	config := GetDisplayConfigForCategory(category)
	return config.Color
}

func (a *AIAnalysisService) getCategoryIcon(category string) string {
	config := GetDisplayConfigForCategory(category)
	return config.Icon
}

func (a *AIAnalysisService) getLevelLabel(value int) string {
	if value >= 80 {
		return "expert"
	} else if value >= 60 {
		return "advanced"
	} else if value >= 40 {
		return "intermediate"
	}
	return "beginner"
}

func (a *AIAnalysisService) estimateCompletionDays(entries []models.KnowledgeBaseEntry) int {
	toLearnCount := 0
	for _, entry := range entries {
		if entry.Level < 4 {
			toLearnCount++
		}
	}

	// 粗略估计：每天学习1-2个知识点
	if toLearnCount > 20 {
		return toLearnCount * 2
	}
	return toLearnCount
}

func isStopword(word string) bool {
	stopwords := map[string]bool{
		"是": true, "的": true, "了": true, "和": true,
		"in": true, "is": true, "the": true, "a": true,
		"to": true, "of": true, "for": true, "on": true,
	}
	return stopwords[strings.ToLower(word)]
}

func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, s := range strs {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (a *AIAnalysisService) trendBucketLabel(t time.Time, granularity TrendGranularity) string {
	loc := t.Location()
	lt := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	switch granularity {
	case TrendGranularityMonth:
		return fmt.Sprintf("%04d-%02d", lt.Year(), int(lt.Month()))
	case TrendGranularityWeek:
		y, w := lt.ISOWeek()
		return fmt.Sprintf("%04d-W%02d", y, w)
	case TrendGranularityDay:
		fallthrough
	default:
		return lt.Format(dateLayoutYMD)
	}
}

func (a *AIAnalysisService) trendBucketLabels(from, to time.Time, granularity TrendGranularity) []string {
	loc := to.Location()
	// 统一到当天零点，避免边界问题
	start := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, loc)
	end := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, loc)

	switch granularity {
	case TrendGranularityMonth:
		// 从 start 月到 end 月
		cur := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, loc)
		endM := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, loc)
		labels := make([]string, 0, 12)
		for !cur.After(endM) {
			labels = append(labels, fmt.Sprintf("%04d-%02d", cur.Year(), int(cur.Month())))
			cur = cur.AddDate(0, 1, 0)
		}
		return labels

	case TrendGranularityWeek:
		// ISO 周 label 序列：按周一递增
		cur := start
		// 对齐到周一
		wd := int(cur.Weekday())
		if wd == 0 {
			wd = 7 // Sunday -> 7
		}
		cur = cur.AddDate(0, 0, -(wd - 1))
		labels := make([]string, 0, 16)
		for !cur.After(end) {
			y, w := cur.ISOWeek()
			labels = append(labels, fmt.Sprintf("%04d-W%02d", y, w))
			cur = cur.AddDate(0, 0, 7)
		}
		// 去重（跨年边界可能重复）
		seen := make(map[string]struct{}, len(labels))
		out := make([]string, 0, len(labels))
		for _, l := range labels {
			if _, ok := seen[l]; ok {
				continue
			}
			seen[l] = struct{}{}
			out = append(out, l)
		}
		return out

	case TrendGranularityDay:
		fallthrough
	default:
		labels := make([]string, 0, 90)
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			labels = append(labels, d.Format(dateLayoutYMD))
		}
		return labels
	}
}
