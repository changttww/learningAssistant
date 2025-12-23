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

// RAGService RAG服务接口
type RAGService interface {
	// 添加文档到知识库
	AddDocument(userID uint64, sourceType int8, sourceID uint64, title, content string) (*models.KnowledgeBaseEntry, error)
	// 删除文档
	RemoveDocument(entryID uint64) error
	// 搜索知识库
	SearchKnowledge(userID uint64, query string, limit int) ([]models.KnowledgeBaseEntry, error)
	// 获取用户知识库统计
	GetUserKnowledgeStats(userID uint64) (map[string]interface{}, error)
	// 更新知识点掌握等级
	UpdateKnowledgeLevel(entryID uint64, level int8) error
	// 获取知识点关系
	GetKnowledgeRelations(entryID uint64) ([]models.KnowledgeRelation, error)
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
		}
		if err := db.Model(&existingEntry).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("更新知识库条目失败: %w", err)
		}

		// 更新向量缓存
		if vector, err := r.embeddingService.GenerateEmbedding(cleanTitle + " " + summary); err == nil {
			contentHash := md5Hash(cleanContent)
			db.Where("entry_id = ?", existingEntry.ID).Delete(&models.KnowledgeVectorCache{})
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
		if similarity > 0.3 { // 阈值
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
		"数学": {"数学", "math", "代数", "几何", "微积分", "函数", "方程", "公式", "计算", "概率", "统计"},
		"物理": {"物理", "physics", "力学", "电学", "磁学", "光学", "热学", "能量", "牛顿", "运动"},
		"化学": {"化学", "chemistry", "元素", "分子", "原子", "反应", "酸碱", "有机"},
		"生物": {"生物", "biology", "细胞", "遗传", "基因", "生命", "生态", "动物", "植物"},
		// 文科
		"语文": {"语文", "chinese", "文言文", "古诗", "阅读理解", "写作", "作文", "文学", "诗词"},
		"英语": {"英语", "english", "单词", "语法", "词汇", "口语", "听力", "翻译"},
		"历史": {"历史", "history", "朝代", "战争", "革命", "古代", "近代", "历史事件"},
		"地理": {"地理", "geography", "气候", "地形", "区域", "城市", "自然", "环境"},
		"政治": {"政治", "politics", "政策", "制度", "法律", "哲学", "思想"},
		// 技能
		"编程":  {"编程", "programming", "代码", "python", "java", "javascript", "程序", "开发", "算法"},
		"计算机": {"计算机", "computer", "软件", "硬件", "网络", "系统", "数据库"},
		"艺术":  {"艺术", "art", "绘画", "美术", "设计", "色彩", "创作"},
		"音乐":  {"音乐", "music", "歌曲", "乐器", "旋律", "节奏", "音符"},
		"体育":  {"体育", "sports", "运动", "锻炼", "健身", "比赛", "训练"},
		// 学习通识
		"学习方法": {"学习方法", "学习技巧", "记忆", "复习", "笔记", "思维导图", "效率"},
		"考试技巧": {"考试", "exam", "测验", "答题", "解题", "技巧"},
		"阅读":   {"阅读", "reading", "书籍", "文章", "理解"},
		"思维训练": {"思维", "逻辑", "推理", "思考", "分析"},
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
