package rag

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

// SearchResult 带相似度的搜索结果
type SearchResult struct {
	Entry      models.KnowledgeBaseEntry `json:"entry"`
	Similarity float32                   `json:"similarity"`
}

// Citation 引用信息
type Citation struct {
	ID         uint64  `json:"id"`
	Title      string  `json:"title"`
	Category   string  `json:"category"`
	Summary    string  `json:"summary"`
	Similarity float32 `json:"similarity"`
}

// RAGQueryResult RAG问答结果（带引用溯源）
type RAGQueryResult struct {
	Answer    string     `json:"answer"`
	Citations []Citation `json:"citations"`
	Query     string     `json:"query"`
}

// RAGService RAG服务接口
type RAGService interface {
	// 添加文档到知识库
	AddDocument(userID uint64, sourceType int8, sourceID uint64, title, content string) (*models.KnowledgeBaseEntry, error)
	// 删除文档
	RemoveDocument(entryID uint64) error
	// 搜索知识库
	SearchKnowledge(userID uint64, query string, limit int) ([]models.KnowledgeBaseEntry, error)
	// 搜索知识库（带相似度）
	SearchKnowledgeWithScore(userID uint64, query string, limit int) ([]SearchResult, error)
	// 获取用户知识库统计
	GetUserKnowledgeStats(userID uint64) (map[string]interface{}, error)
	// 更新知识点掌握等级
	UpdateKnowledgeLevel(entryID uint64, level int8) error
	// 获取知识点关系
	GetKnowledgeRelations(entryID uint64) ([]models.KnowledgeRelation, error)
	// 获取用户知识图谱数据
	GetKnowledgeGraph(userID uint64) (*KnowledgeGraphData, error)
}

// KnowledgeGraphNode 知识图谱节点
type KnowledgeGraphNode struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Level    int8   `json:"level"`
	Value    int    `json:"value"` // 节点大小，基于ViewCount
	Color    string `json:"color"`
}

// KnowledgeGraphLink 知识图谱边
type KnowledgeGraphLink struct {
	Source       uint64  `json:"source"`
	Target       uint64  `json:"target"`
	RelationType int8    `json:"relation_type"` // 1=prerequisite, 2=related, 3=extends, 4=conflict, 5=same_category, 6=same_tag
	Strength     float32 `json:"strength"`
	Label        string  `json:"label"`
}

// KnowledgeGraphData 知识图谱数据
type KnowledgeGraphData struct {
	Nodes []KnowledgeGraphNode `json:"nodes"`
	Links []KnowledgeGraphLink `json:"links"`
}

// DefaultRAGService 默认RAG服务实现
type DefaultRAGService struct {
	embeddingService EmbeddingService
}

// EmbeddingService 向量化服务接口
type EmbeddingService interface {
	// 生成文本向量
	GenerateEmbedding(text string) (models.Vector, error)
	// 计算向量相似度
	CosineSimilarity(vec1, vec2 models.Vector) float32
	// 批量生成向量
	GenerateEmbeddings(texts []string) ([]models.Vector, error)
}

// 预编译正则表达式以提高性能
var (
	htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
	spaceRegex   = regexp.MustCompile(`\s+`)
)

// stripHTMLTags 去除 HTML 标签，只保留纯文本
func stripHTMLTags(content string) string {
	if content == "" {
		return ""
	}

	// 移除 HTML 标签
	text := htmlTagRegex.ReplaceAllString(content, " ")

	// 处理 HTML 实体
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")

	// 移除多余的空白
	text = spaceRegex.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

// NewRAGService 创建RAG服务
func NewRAGService(embeddingService EmbeddingService) RAGService {
	return &DefaultRAGService{
		embeddingService: embeddingService,
	}
}

// AddDocument 添加文档到知识库（如果已存在则更新）
func (r *DefaultRAGService) AddDocument(userID uint64, sourceType int8, sourceID uint64, title, content string) (*models.KnowledgeBaseEntry, error) {
	db := database.GetDB()

	// 关联任务/笔记ID，便于后续按 task_id/note_id 追溯
	var taskID *uint64
	var noteID *uint64
	if sourceID > 0 {
		switch sourceType {
		case 1: // 任务
			taskID = &sourceID
		case 2: // 笔记
			noteID = &sourceID
		}
	}

	// 清理 HTML 标签，保存纯文本
	cleanTitle := stripHTMLTags(title)
	cleanContent := stripHTMLTags(content)

	// 生成摘要和关键词
	summary := generateSummary(cleanContent)
	keywords := extractKeywords(cleanContent)
	category, subCategory := classifyContent(cleanTitle, cleanContent)

	// 获取显示配置
	displayConfig := GetDisplayConfigForCategory(category)

	// 检查是否已存在同来源的知识条目
	var existingEntry models.KnowledgeBaseEntry
	result := db.Where("user_id = ? AND source_type = ? AND source_id = ?", userID, sourceType, sourceID).First(&existingEntry)

	if result.Error == nil {
		// 已存在，更新内容
		updates := map[string]interface{}{
			"title":         cleanTitle,
			"content":       cleanContent,
			"summary":       summary,
			"keywords":      keywords,
			"category":      category,
			"sub_category":  subCategory,
			"status":        1,
			"display_color": displayConfig.Color,
			"display_icon":  displayConfig.Icon,
			"subject":       category,
			"task_id":       taskID,
			"note_id":       noteID,
		}
		if err := db.Model(&existingEntry).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("更新知识库条目失败: %w", err)
		}

		// 更新向量缓存
		if vector, err := r.embeddingService.GenerateEmbedding(cleanTitle + " " + summary); err == nil {
			contentHash := md5Hash(cleanContent)
			// 使用 Unscoped 硬删除旧的向量缓存，避免唯一索引冲突
			db.Unscoped().Where("entry_id = ?", existingEntry.ID).Delete(&models.KnowledgeVectorCache{})
			cache := &models.KnowledgeVectorCache{
				EntryID:     existingEntry.ID,
				ContentHash: contentHash,
				Vector:      vector,
				VectorDim:   len(vector),
				VectorModel: "qwen-embedding",
			}
			db.Create(cache)
		}

		return &existingEntry, nil
	}

	// 不存在，创建新条目
	entry := &models.KnowledgeBaseEntry{
		UserID:       userID,
		SourceType:   sourceType,
		SourceID:     sourceID,
		TaskID:       taskID,
		NoteID:       noteID,
		Title:        cleanTitle,
		Content:      cleanContent,
		Summary:      summary,
		Keywords:     keywords,
		Category:     category,
		SubCategory:  subCategory,
		Level:        0, // 初始等级为未学习
		Status:       1, // 默认发布
		DisplayColor: displayConfig.Color,
		DisplayIcon:  displayConfig.Icon,
		Subject:      category,
	}

	if err := db.Create(entry).Error; err != nil {
		return nil, fmt.Errorf("创建知识库条目失败: %w", err)
	}

	// 生成向量并缓存
	if vector, err := r.embeddingService.GenerateEmbedding(cleanTitle + " " + summary); err == nil {
		contentHash := md5Hash(cleanContent)
		cache := &models.KnowledgeVectorCache{
			EntryID:     entry.ID,
			ContentHash: contentHash,
			Vector:      vector,
			VectorDim:   len(vector),
			VectorModel: "qwen-embedding",
		}
		db.Create(cache)
	}

	// 更新用户统计
	r.updateUserStats(userID)

	return entry, nil
}

// RemoveDocument 删除文档
func (r *DefaultRAGService) RemoveDocument(entryID uint64) error {
	db := database.GetDB()

	// 删除向量缓存
	db.Where("entry_id = ?", entryID).Delete(&models.KnowledgeVectorCache{})

	// 删除关系
	db.Where("source_entry_id = ? OR target_entry_id = ?", entryID, entryID).
		Delete(&models.KnowledgeRelation{})

	// 删除条目
	if err := db.Delete(&models.KnowledgeBaseEntry{}, entryID).Error; err != nil {
		return fmt.Errorf("删除知识库条目失败: %w", err)
	}

	return nil
}

// SearchKnowledge 搜索知识库
func (r *DefaultRAGService) SearchKnowledge(userID uint64, query string, limit int) ([]models.KnowledgeBaseEntry, error) {
	if limit <= 0 {
		limit = 10
	}

	db := database.GetDB()
	var entries []models.KnowledgeBaseEntry

	// 首先尝试向量相似度搜索
	queryVector, err := r.embeddingService.GenerateEmbedding(query)
	if err == nil && len(queryVector) > 0 {
		vectorResults, vectorErr := r.vectorSearch(userID, queryVector, limit)
		// 只有向量搜索成功且有结果时才返回
		if vectorErr == nil && len(vectorResults) > 0 {
			return vectorResults, nil
		}
		// 否则降级到关键词搜索
	}

	// 降级到关键词搜索
	searchPattern := "%" + query + "%"
	if err := db.Where("user_id = ? AND status = 1", userID).
		Where("title LIKE ? OR content LIKE ? OR keywords LIKE ? OR category LIKE ?", searchPattern, searchPattern, searchPattern, searchPattern).
		Order("level DESC, view_count DESC").
		Limit(limit).
		Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("搜索知识库失败: %w", err)
	}

	return entries, nil
}

// vectorSearch 向量相似度搜索
func (r *DefaultRAGService) vectorSearch(userID uint64, queryVector models.Vector, limit int) ([]models.KnowledgeBaseEntry, error) {
	db := database.GetDB()
	var entries []models.KnowledgeBaseEntry

	// 获取用户的所有向量缓存
	var caches []models.KnowledgeVectorCache
	if err := db.Where("entry_id IN (?)",
		db.Table("knowledge_base_entries").
			Select("id").
			Where("user_id = ? AND status = 1", userID)).
		Find(&caches).Error; err != nil {
		return nil, err
	}

	// 计算相似度并排序
	type scoredEntry struct {
		entry   models.KnowledgeBaseEntry
		score   float32
		cacheID uint64
	}

	var scored []scoredEntry
	for _, cache := range caches {
		similarity := r.embeddingService.CosineSimilarity(queryVector, cache.Vector)
		// 提高阈值到 0.35，避免返回不相关内容
		if similarity >= 0.35 {
			scored = append(scored, scoredEntry{
				score:   similarity,
				cacheID: cache.ID,
			})
		}
	}

	// 按相似度排序
	// 这里简化处理，实际可用更完善的排序
	if len(scored) == 0 {
		return []models.KnowledgeBaseEntry{}, nil
	}

	// 获取top-k的条目
	for i := 0; i < len(scored) && i < limit; i++ {
		var cache models.KnowledgeVectorCache
		db.First(&cache, scored[i].cacheID)
		var entry models.KnowledgeBaseEntry
		db.First(&entry, cache.EntryID)
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetUserKnowledgeStats 获取用户知识库统计
func (r *DefaultRAGService) GetUserKnowledgeStats(userID uint64) (map[string]interface{}, error) {
	db := database.GetDB()

	// 直接从知识库条目统计各等级数量
	type LevelCount struct {
		Level int8
		Count int64
	}
	var levelCounts []LevelCount

	if err := db.Model(&models.KnowledgeBaseEntry{}).
		Select("level, count(*) as count").
		Where("user_id = ? AND status = 1", userID).
		Group("level").
		Find(&levelCounts).Error; err != nil {
		return nil, fmt.Errorf("获取知识库统计失败: %w", err)
	}

	// 构建统计结果
	stats := map[string]interface{}{
		"level_0_count": int64(0), // 待学习
		"level_1_count": int64(0), // 了解
		"level_2_count": int64(0), // 熟悉
		"level_3_count": int64(0), // 已掌握
		"total_count":   int64(0),
		"review_needed": int64(0),
	}

	for _, lc := range levelCounts {
		switch lc.Level {
		case 0:
			stats["level_0_count"] = lc.Count
		case 1:
			stats["level_1_count"] = lc.Count
		case 2:
			stats["level_2_count"] = lc.Count
		case 3:
			stats["level_3_count"] = lc.Count
		}
		stats["total_count"] = stats["total_count"].(int64) + lc.Count
	}

	// 查询需要复习的知识点（超过7天未复习且等级低于3）
	var reviewCount int64
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	db.Model(&models.KnowledgeBaseEntry{}).
		Where("user_id = ? AND status = 1 AND level < 3 AND (last_review_at IS NULL OR last_review_at < ?)", userID, sevenDaysAgo).
		Count(&reviewCount)
	stats["review_needed"] = reviewCount

	return stats, nil
}

// UpdateKnowledgeLevel 更新知识点掌握等级
func (r *DefaultRAGService) UpdateKnowledgeLevel(entryID uint64, level int8) error {
	db := database.GetDB()

	if level < 0 || level > 4 {
		return fmt.Errorf("无效的等级: %d", level)
	}

	return db.Model(&models.KnowledgeBaseEntry{}).
		Where("id = ?", entryID).
		Update("level", level).Error
}

// GetKnowledgeRelations 获取知识点关系
func (r *DefaultRAGService) GetKnowledgeRelations(entryID uint64) ([]models.KnowledgeRelation, error) {
	db := database.GetDB()
	var relations []models.KnowledgeRelation

	if err := db.Where("source_entry_id = ? OR target_entry_id = ?", entryID, entryID).
		Find(&relations).Error; err != nil {
		return nil, fmt.Errorf("获取知识关系失败: %w", err)
	}

	return relations, nil
}

// 私有方法

// md5Hash 计算字符串MD5哈希
func md5Hash(text string) string {
	h := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", h)
}

// generateSummary 生成内容摘要
func generateSummary(content string) string {
	// 简单实现：取前200字符
	if len(content) <= 200 {
		return content
	}
	return content[:200] + "..."
}

// extractKeywords 提取关键词
func extractKeywords(content string) []byte {
	// 简单实现：分词
	words := strings.Fields(content)
	keywords := make([]string, 0)

	// 过滤短单词和常用词
	stopwords := map[string]bool{
		"是": true, "的": true, "了": true, "和": true,
		"in": true, "is": true, "the": true, "a": true,
	}

	for _, word := range words {
		if len(word) > 2 && !stopwords[strings.ToLower(word)] {
			keywords = append(keywords, word)
			if len(keywords) >= 10 {
				break
			}
		}
	}

	// 转换为JSON
	data, _ := json.Marshal(keywords)
	return data
}

// classifyContent 分类内容 - 面向学习助手场景
func classifyContent(title, content string) (category, subCategory string) {
	// 将标题和内容转为小写进行匹配
	fullText := strings.ToLower(title + " " + content)

	// 学习场景的分类体系
	categories := map[string][]string{
		// 理科
		"数学": {"数学", "math", "代数", "几何", "微积分", "函数", "方程", "公式", "计算", "概率", "统计", "线性代数", "Linear Algebra", "高等数学", "Advanced Mathematics"},
		"物理": {"物理", "physics", "力学", "电学", "磁学", "光学", "热学", "能量", "牛顿", "运动", "大学物理", "College Physics"},
		// （化学、生物分类保持不变）
		"化学": {"化学", "chemistry", "元素", "分子", "原子", "反应", "酸碱", "有机"},
		"生物": {"生物", "biology", "细胞", "遗传", "基因", "生命", "生态", "动物", "植物"},
		// 文科（保持不变）
		"语文": {"语文", "chinese", "文言文", "古诗", "阅读理解", "写作", "作文", "文学", "诗词"},
		"英语": {"英语", "english", "单词", "语法", "词汇", "口语", "听力", "翻译"},
		"历史": {"历史", "history", "朝代", "战争", "革命", "古代", "近代", "历史事件"},
		"地理": {"地理", "geography", "气候", "地形", "区域", "城市", "自然", "环境"},
		"政治": {"政治", "politics", "政策", "制度", "法律", "哲学", "思想"},
		// 技能
		"编程":  {"编程", "programming", "代码", "python", "java", "javascript", "程序", "开发", "算法", "程序设计", "Computer Programming", "程序分析", "Program Analysis"},
		"计算机": {"计算机", "computer", "软件", "硬件", "网络", "系统", "数据库", "计算机网络", "Computer Networks", "操作系统", "Operating Systems", "编译原理", "Principles of Compilers", "数字电路与逻辑设计", "Digital Circuits and Logical Design"},
		// （艺术、音乐、体育分类保持不变）
		"艺术": {"艺术", "art", "绘画", "美术", "设计", "色彩", "创作"},
		"音乐": {"音乐", "music", "歌曲", "乐器", "旋律", "节奏", "音符"},
		"体育": {"体育", "sports", "运动", "锻炼", "健身", "比赛", "训练"},
		// 学习通识（保持不变）
		"学习方法": {"学习方法", "学习技巧", "记忆", "复习", "笔记", "思维导图", "效率"},
		"考试技巧": {"考试", "exam", "测验", "答题", "解题", "技巧"},
		"阅读":   {"阅读", "reading", "书籍", "文章", "理解"},
		"思维训练": {"思维", "逻辑", "推理", "思考", "分析"},
		// 软件工程专业核心分类（新增图中课程对应的关键词）
		"软件工程":     {"软件工程", "software engineering", "需求分析", "软件设计", "项目管理", "软件生命周期", "工程化", "架构设计", "模块化", "软件工程实训", "Software Engineering Training"},
		"软件测试":     {"软件测试", "software testing", "黑盒测试", "白盒测试", "自动化测试", "单元测试", "集成测试", "测试用例", "bug", "质量保证", "软件质量保证", "Software Quality Assurance", "软件测试实验", "Software Testing Experiment"},
		"数据库工程":    {"数据库工程", "database engineering", "sql", "mysql", "oracle", "postgres", "数据建模", "索引", "事务", "分库分表", "数据备份", "数据库系统原理", "Principles of Database Systems", "数据库系统实验", "Database Systems Laboratory"},
		"软件开发框架":   {"框架", "framework", "spring", "django", "flask", "vue", "react", "angular", "mybatis", "微服务"},
		"版本控制":     {"版本控制", "version control", "git", "github", "gitlab", "commit", "branch", "merge", "revert", "reset"},
		"操作系统与内核":  {"操作系统", "os", "linux", "windows", "内核", "进程", "线程", "内存管理", "文件系统", "驱动", "操作系统原理", "Principles of Operating Systems", "操作系统实验", "Operating Systems Laboratory"},
		"计算机网络与接口": {"网络编程", "network", "http", "tcp/ip", "restful", "api", "接口", "socket", "网关", "负载均衡", "计算机网络实验", "Computer Networks Laboratory"},
		"软件架构":     {"软件架构", "architecture", "单体架构", "微服务架构", "分布式架构", "云原生", "docker", "k8s", "服务网格"},
		"编译原理与实现":  {"编译原理", "Principles of Compilers", "编译器构造实验", "Compilers Construction Laboratory"},
		"软件需求分析":   {"软件需求分析与设计", "Software Analysis and Design", "软件需求分析与设计实验", "Software Analysis and Design Laboratory"},
	}

	// 匹配分类
	for cat, keywords := range categories {
		for _, keyword := range keywords {
			if strings.Contains(fullText, keyword) {
				category = cat
				subCategory = keyword
				return
			}
		}
	}

	category = "其他"
	subCategory = ""
	return
}

// updateUserStats 更新用户统计
func (r *DefaultRAGService) updateUserStats(userID uint64) error {
	db := database.GetDB()

	var stats models.UserKnowledgeStats
	db.FirstOrCreate(&stats, models.UserKnowledgeStats{UserID: userID})

	// 统计各等级知识点数
	var masteredCount, learningCount, toLearnCount int64
	db.Model(&models.KnowledgeBaseEntry{}).
		Where("user_id = ? AND status = 1 AND level = 4", userID).
		Count(&masteredCount)
	db.Model(&models.KnowledgeBaseEntry{}).
		Where("user_id = ? AND status = 1 AND level IN (1, 2, 3)", userID).
		Count(&learningCount)
	db.Model(&models.KnowledgeBaseEntry{}).
		Where("user_id = ? AND status = 1 AND level = 0", userID).
		Count(&toLearnCount)

	var totalCount int64
	db.Model(&models.KnowledgeBaseEntry{}).
		Where("user_id = ? AND status = 1", userID).
		Count(&totalCount)

	now := time.Now()
	return db.Model(&stats).Updates(map[string]interface{}{
		"total_entries":  totalCount,
		"mastered_count": masteredCount,
		"learning_count": learningCount,
		"to_learn_count": toLearnCount,
		"last_update_at": now,
	}).Error
}

// SearchKnowledgeWithScore 搜索知识库并返回相似度分数
func (r *DefaultRAGService) SearchKnowledgeWithScore(userID uint64, query string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 10
	}

	db := database.GetDB()
	var results []SearchResult

	// 首先尝试向量相似度搜索
	queryVector, err := r.embeddingService.GenerateEmbedding(query)
	if err == nil && len(queryVector) > 0 {
		vectorResults, vectorErr := r.vectorSearchWithScore(userID, queryVector, limit)
		if vectorErr == nil && len(vectorResults) > 0 {
			return vectorResults, nil
		}
	}

	// 降级到关键词搜索
	var entries []models.KnowledgeBaseEntry
	searchPattern := "%" + query + "%"
	if err := db.Where("user_id = ? AND status = 1", userID).
		Where("title LIKE ? OR content LIKE ? OR keywords LIKE ? OR category LIKE ?", searchPattern, searchPattern, searchPattern, searchPattern).
		Order("level DESC, view_count DESC").
		Limit(limit).
		Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("搜索知识库失败: %w", err)
	}

	// 关键词匹配给一个默认相似度
	for _, entry := range entries {
		results = append(results, SearchResult{
			Entry:      entry,
			Similarity: 0.5, // 关键词匹配的默认相似度
		})
	}

	return results, nil
}

// vectorSearchWithScore 向量相似度搜索（带分数）
func (r *DefaultRAGService) vectorSearchWithScore(userID uint64, queryVector models.Vector, limit int) ([]SearchResult, error) {
	db := database.GetDB()

	// 获取用户的所有知识条目
	var entries []models.KnowledgeBaseEntry
	if err := db.Where("user_id = ? AND status = 1", userID).Find(&entries).Error; err != nil {
		return nil, err
	}

	// 获取这些条目的向量缓存
	entryIDs := make([]uint64, len(entries))
	entryMap := make(map[uint64]models.KnowledgeBaseEntry)
	for i, entry := range entries {
		entryIDs[i] = entry.ID
		entryMap[entry.ID] = entry
	}

	var caches []models.KnowledgeVectorCache
	if err := db.Where("entry_id IN ?", entryIDs).Find(&caches).Error; err != nil {
		return nil, err
	}

	// 计算相似度
	type scoredResult struct {
		entryID    uint64
		similarity float32
	}
	var scored []scoredResult
	for _, cache := range caches {
		similarity := r.embeddingService.CosineSimilarity(queryVector, cache.Vector)
		// 提高阈值到 0.35，避免返回不相关内容
		if similarity >= 0.35 {
			scored = append(scored, scoredResult{
				entryID:    cache.EntryID,
				similarity: similarity,
			})
		}
	}

	// 按相似度排序（降序）
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].similarity > scored[i].similarity {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// 获取top-k的条目
	var results []SearchResult
	for i := 0; i < len(scored) && i < limit; i++ {
		if entry, ok := entryMap[scored[i].entryID]; ok {
			results = append(results, SearchResult{
				Entry:      entry,
				Similarity: scored[i].similarity,
			})
		}
	}

	return results, nil
}

// GetKnowledgeGraph 获取用户知识图谱数据
func (r *DefaultRAGService) GetKnowledgeGraph(userID uint64) (*KnowledgeGraphData, error) {
	db := database.GetDB()

	// 获取用户的所有知识条目（不限制 status，显示所有条目）
	var entries []models.KnowledgeBaseEntry
	if err := db.Where("user_id = ?", userID).Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("获取知识条目失败: %w", err)
	}

	// 构建节点
	nodes := make([]KnowledgeGraphNode, 0, len(entries))
	categoryColorMap := getCategoryColorMap()
	entryMap := make(map[uint64]models.KnowledgeBaseEntry)

	for _, entry := range entries {
		entryMap[entry.ID] = entry
		color := categoryColorMap[entry.Category]
		if color == "" {
			color = "#9ca3af" // 默认灰色
		}
		nodes = append(nodes, KnowledgeGraphNode{
			ID:       entry.ID,
			Name:     truncateString(entry.Title, 20),
			Category: entry.Category,
			Level:    entry.Level,
			Value:    entry.ViewCount + 10, // 基础大小 + 浏览次数
			Color:    color,
		})
	}

	// 获取显式关系
	var relations []models.KnowledgeRelation
	db.Where("user_id = ?", userID).Find(&relations)

	links := make([]KnowledgeGraphLink, 0)
	linkSet := make(map[string]bool) // 用于去重

	// 添加显式关系
	relationLabels := map[int8]string{
		1: "前置",
		2: "相关",
		3: "扩展",
		4: "冲突",
	}
	for _, rel := range relations {
		key := fmt.Sprintf("%d-%d", rel.SourceEntryID, rel.TargetEntryID)
		if !linkSet[key] {
			linkSet[key] = true
			links = append(links, KnowledgeGraphLink{
				Source:       rel.SourceEntryID,
				Target:       rel.TargetEntryID,
				RelationType: rel.RelationType,
				Strength:     rel.Strength,
				Label:        relationLabels[rel.RelationType],
			})
		}
	}

	// 添加隐式关系：同分类
	categoryEntries := make(map[string][]uint64)
	for _, entry := range entries {
		categoryEntries[entry.Category] = append(categoryEntries[entry.Category], entry.ID)
	}

	for _, entryIDs := range categoryEntries {
		if len(entryIDs) < 2 {
			continue
		}
		// 为同分类的条目创建关系（限制数量避免太多连线）
		maxLinks := 3
		for i := 0; i < len(entryIDs) && i < maxLinks; i++ {
			for j := i + 1; j < len(entryIDs) && j < maxLinks+1; j++ {
				key1 := fmt.Sprintf("%d-%d", entryIDs[i], entryIDs[j])
				key2 := fmt.Sprintf("%d-%d", entryIDs[j], entryIDs[i])
				if !linkSet[key1] && !linkSet[key2] {
					linkSet[key1] = true
					links = append(links, KnowledgeGraphLink{
						Source:       entryIDs[i],
						Target:       entryIDs[j],
						RelationType: 5, // 同分类
						Strength:     0.3,
						Label:        "同分类",
					})
				}
			}
		}
	}

	return &KnowledgeGraphData{
		Nodes: nodes,
		Links: links,
	}, nil
}

// getCategoryColorMap 获取分类颜色映射
func getCategoryColorMap() map[string]string {
	return map[string]string{
		"数学":   "#3b82f6",
		"物理":   "#8b5cf6",
		"化学":   "#10b981",
		"生物":   "#22c55e",
		"语文":   "#f59e0b",
		"英语":   "#ef4444",
		"历史":   "#d97706",
		"地理":   "#14b8a6",
		"政治":   "#6366f1",
		"编程":   "#06b6d4",
		"计算机":  "#0ea5e9",
		"艺术":   "#ec4899",
		"音乐":   "#f472b6",
		"体育":   "#84cc16",
		"学习方法": "#a855f7",
		"考试技巧": "#f97316",
		"阅读":   "#eab308",
		"思维训练": "#8b5cf6",
		"其他":   "#9ca3af",
	}
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
