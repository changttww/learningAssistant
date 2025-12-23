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
type LearningTrend struct {
	Date          string  `json:"date"`
	NewKnowledge  int     `json:"new_knowledge"`
	ReviewCount   int     `json:"review_count"`
	MasteredCount int     `json:"mastered_count"`
	StudyHours    float32 `json:"study_hours"`
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

// generateSkillRadar 生成技能雷达数据
func (a *AIAnalysisService) generateSkillRadar(entries []models.KnowledgeBaseEntry) []SkillRadarData {
	skillMap := make(map[string]*SkillRadarData)

	// 从知识库条目中提取技能
	for _, entry := range entries {
		skills := a.extractSkills(entry)
		for _, skill := range skills {
			if _, exists := skillMap[skill]; !exists {
				skillMap[skill] = &SkillRadarData{
					Skill:    skill,
					Value:    0,
					MaxValue: 100,
					Category: entry.Category,
				}
			}

			skillMap[skill].Value += int(entry.Level) * 20
		}
	}

	// 转换为切片
	radar := make([]SkillRadarData, 0, len(skillMap))
	for _, skill := range skillMap {
		if skill.Value > 100 {
			skill.Value = 100
		}
		skill.Progress = skill.Value
		skill.Level = a.getLevelLabel(skill.Value)
		radar = append(radar, *skill)
	}

	// 按技能值排序并取前10个
	sort.Slice(radar, func(i, j int) bool {
		return radar[i].Value > radar[j].Value
	})

	if len(radar) > 10 {
		radar = radar[:10]
	}

	return radar
}

// analyzeLearningTrends 分析学习趋势（最近30天）
func (a *AIAnalysisService) analyzeLearningTrends(userID uint64, entries []models.KnowledgeBaseEntry) []LearningTrend {
	db := database.GetDB()

	// 构建日期到条目的映射
	trendMap := make(map[string]*LearningTrend)
	now := time.Now()

	for i := 29; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		trendMap[dateStr] = &LearningTrend{
			Date:          dateStr,
			NewKnowledge:  0,
			ReviewCount:   0,
			MasteredCount: 0,
			StudyHours:    0,
		}
	}

	// 统计每日创建的条目
	for _, entry := range entries {
		dateStr := entry.CreatedAt.Format("2006-01-02")
		if trend, exists := trendMap[dateStr]; exists {
			trend.NewKnowledge++
		}
	}

	// 统计学习时长
	var stats []models.DailyStudyStat
	fromDate := now.AddDate(0, 0, -30)
	db.Where("user_id = ? AND date >= ?", userID, fromDate).
		Find(&stats)

	for _, stat := range stats {
		dateStr := stat.Date.Format("2006-01-02")
		if trend, exists := trendMap[dateStr]; exists {
			trend.StudyHours = float32(stat.Minutes) / 60.0
		}
	}

	// 转换为切片
	trends := make([]LearningTrend, 0, len(trendMap))
	for _, trend := range trendMap {
		trends = append(trends, *trend)
	}

	// 按日期排序
	sort.Slice(trends, func(i, j int) bool {
		return trends[i].Date < trends[j].Date
	})

	return trends
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
