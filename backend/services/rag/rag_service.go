package rag

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	// 添加/更新任务知识点（聚合任务及其所有笔记为一个知识点）
	AddTaskKnowledge(userID uint64, taskID uint64) (*models.KnowledgeBaseEntry, error)
	// 删除文档
	RemoveDocument(entryID uint64) error
	// 按任务ID删除知识点
	RemoveTaskKnowledge(userID uint64, taskID uint64) error
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
	GetKnowledgeGraph(userID uint64, teamID *uint64) (*KnowledgeGraphData, error)
	// 批量重分类现有知识点（数据清洗）
	ReclassifyAllEntries(userID uint64) (int, error)
}

// KnowledgeGraphNode 知识图谱节点
type KnowledgeGraphNode struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	Level      int8   `json:"level"`
	Value      int    `json:"value"` // 节点大小，基于ViewCount
	Color      string `json:"color"`
	IsVirtual  bool   `json:"is_virtual"`  // 是否为虚拟中心节点
	SymbolSize int    `json:"symbol_size"` // 节点显示大小
}

// KnowledgeGraphLink 知识图谱边
type KnowledgeGraphLink struct {
	Source       uint64  `json:"source"`
	Target       uint64  `json:"target"`
	RelationType int8    `json:"relation_type"` // 1=prerequisite, 2=related, 3=extends, 4=conflict, 5=same_category, 6=归属中心
	Strength     float32 `json:"strength"`
	Label        string  `json:"label"`
}

// KnowledgeGraphData 知识图谱数据
type KnowledgeGraphData struct {
	Nodes      []KnowledgeGraphNode `json:"nodes"`
	Links      []KnowledgeGraphLink `json:"links"`
	Categories []GraphCategory      `json:"categories"` // 分类信息（用于图例）
}

// GraphCategory 图谱分类信息
type GraphCategory struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// AIClassificationResult AI分类结果
type AIClassificationResult struct {
	Category    string `json:"category"`     // 学科大类（环形图内圈）
	SubCategory string `json:"sub_category"` // 细分领域（环形图外圈）
	Subject     string `json:"subject"`      // 能力维度（雷达图）
}

// 预定义的分类枚举（限制AI只能从这些值中选择）
var (
	// CategoryEnum 学科大类（环形图内圈，5-6个）
	CategoryEnum = []string{
		"计算机",  // 编程、数据结构、算法、数据库、网络等
		"人文社科", // 文学、历史、哲学、政治、语言等
		"数理逻辑", // 数学、物理、逻辑推理等
		"自然科学", // 化学、生物、地理、环境等
		"经济管理", // 经济学、管理学、金融、会计等
		"艺术体育", // 音乐、美术、体育、设计等
	}

	// SubjectEnum 能力维度（雷达图，5个）
	SubjectEnum = []string{
		"理论素养", // 底层的知识积累（原: 记忆理解）
		"逻辑思维", // 理性的思考能力
		"实操应用", // 动手解决问题的能力（原: 工程实践）
		"创新思维", // 不仅是设计，强调创新（原: 创意设计）
		"沟通表达", // 输出和传播能力（原: 语言表达）
	}
)

// DefaultRAGService 默认RAG服务实现
type DefaultRAGService struct {
	embeddingService       EmbeddingService
	relationMiningService  *RelationMiningService
	enableAutoRelationMine bool // 是否开启自动关系挖掘
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
		embeddingService:       embeddingService,
		relationMiningService:  NewRelationMiningService(embeddingService),
		enableAutoRelationMine: true, // 默认开启自动关系挖掘
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
	category, subCategory, subject := classifyContent(cleanTitle, cleanContent)

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
			"subject":       subject, // 使用AI返回的能力维度
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
		Subject:      subject, // 使用AI返回的能力维度
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

		// 自动挖掘关系（异步执行，不阻塞主流程）
		if r.enableAutoRelationMine && r.relationMiningService != nil {
			go func(uid, eid uint64) {
				result, err := r.relationMiningService.MineRelationsForEntry(uid, eid)
				if err != nil {
					fmt.Printf("[AutoRelationMine] 为知识点 %d 挖掘关系失败: %v\n", eid, err)
				} else {
					fmt.Printf("[AutoRelationMine] 为知识点 %d 挖掘到 %d 个关系\n", eid, result.RelationsFound)
				}
			}(userID, entry.ID)
		}
	}

	// 更新用户统计
	r.updateUserStats(userID)

	return entry, nil
}

// AddTaskKnowledge 添加/更新任务知识点
// 简化版：只保存摘要和关联的任务/笔记ID，不再拼接完整内容
// 详细内容通过跳转到任务或笔记页面查看
func (r *DefaultRAGService) AddTaskKnowledge(userID uint64, taskID uint64) (*models.KnowledgeBaseEntry, error) {
	db := database.GetDB()

	// 1. 获取任务信息
	var task models.Task
	if err := db.First(&task, taskID).Error; err != nil {
		return nil, fmt.Errorf("获取任务失败: %w", err)
	}

	// 验证任务归属
	if task.CreatedBy != userID && (task.OwnerUserID == nil || *task.OwnerUserID != userID) {
		return nil, fmt.Errorf("无权访问该任务")
	}

	// 2. 获取任务关联的所有笔记ID
	var notes []models.StudyNote
	db.Where("task_id = ? AND user_id = ?", taskID, userID).Select("id", "title").Find(&notes)

	// 收集笔记ID列表
	noteIDs := make([]uint64, 0, len(notes))
	noteTitles := make([]string, 0, len(notes))
	for _, note := range notes {
		noteIDs = append(noteIDs, note.ID)
		if note.Title != "" {
			noteTitles = append(noteTitles, stripHTMLTags(note.Title))
		}
	}
	noteIDsJSON, _ := json.Marshal(noteIDs)

	// 3. 构建用于AI分析的简短内容（只用于生成摘要和分类，不存储完整内容）
	cleanTitle := stripHTMLTags(task.Title)
	cleanDesc := stripHTMLTags(task.Description)

	// 简短内容：任务标题 + 描述摘要 + 笔记标题
	var briefContent strings.Builder
	briefContent.WriteString(cleanTitle)
	if cleanDesc != "" {
		// 只取描述的前200字用于分析
		descRunes := []rune(cleanDesc)
		if len(descRunes) > 200 {
			briefContent.WriteString("\n")
			briefContent.WriteString(string(descRunes[:200]))
		} else {
			briefContent.WriteString("\n")
			briefContent.WriteString(cleanDesc)
		}
	}
	if len(noteTitles) > 0 {
		briefContent.WriteString("\n关联笔记: ")
		briefContent.WriteString(strings.Join(noteTitles, ", "))
	}

	contentForAnalysis := briefContent.String()

	// 4. 生成知识点数据
	summary := generateSummary(contentForAnalysis)
	keywords := extractKeywords(contentForAnalysis)
	category, subCategory, subject := classifyContent(cleanTitle, contentForAnalysis)
	displayConfig := GetDisplayConfigForCategory(category)

	// 5. 简化的存储内容（只存必要信息用于搜索）
	// 不再拼接完整的笔记内容，详情通过跳转查看
	storedContent := cleanTitle
	if cleanDesc != "" {
		descRunes := []rune(cleanDesc)
		if len(descRunes) > 300 {
			storedContent += "\n" + string(descRunes[:300]) + "..."
		} else {
			storedContent += "\n" + cleanDesc
		}
	}

	// 6. 检查是否已存在该任务的知识条目
	var existingEntry models.KnowledgeBaseEntry
	result := db.Where("user_id = ? AND task_id = ?", userID, taskID).First(&existingEntry)
	if result.Error != nil {
		result = db.Where("user_id = ? AND source_type = 1 AND source_id = ?", userID, taskID).First(&existingEntry)
	}

	if result.Error == nil {
		// 已存在，更新内容
		updates := map[string]interface{}{
			"title":         cleanTitle,
			"content":       storedContent,
			"summary":       summary,
			"keywords":      keywords,
			"category":      category,
			"sub_category":  subCategory,
			"status":        1,
			"display_color": displayConfig.Color,
			"display_icon":  displayConfig.Icon,
			"subject":       subject,
			"task_id":       taskID,
			"note_ids":      noteIDsJSON,
		}
		if err := db.Model(&existingEntry).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("更新知识库条目失败: %w", err)
		}

		// 更新向量缓存
		if vector, err := r.embeddingService.GenerateEmbedding(cleanTitle + " " + summary); err == nil {
			contentHash := md5Hash(storedContent)
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
	taskIDPtr := taskID
	entry := &models.KnowledgeBaseEntry{
		UserID:       userID,
		SourceType:   1, // 任务知识点
		SourceID:     taskID,
		TaskID:       &taskIDPtr,
		NoteIDs:      noteIDsJSON,
		Title:        cleanTitle,
		Content:      storedContent,
		Summary:      summary,
		Keywords:     keywords,
		Category:     category,
		SubCategory:  subCategory,
		Level:        0,
		Status:       1,
		DisplayColor: displayConfig.Color,
		DisplayIcon:  displayConfig.Icon,
		Subject:      subject,
	}

	if err := db.Create(entry).Error; err != nil {
		return nil, fmt.Errorf("创建知识库条目失败: %w", err)
	}

	// 生成向量并缓存
	if vector, err := r.embeddingService.GenerateEmbedding(cleanTitle + " " + summary); err == nil {
		contentHash := md5Hash(storedContent)
		cache := &models.KnowledgeVectorCache{
			EntryID:     entry.ID,
			ContentHash: contentHash,
			Vector:      vector,
			VectorDim:   len(vector),
			VectorModel: "qwen-embedding",
		}
		db.Create(cache)

		// 自动挖掘关系（异步执行，不阻塞主流程）
		if r.enableAutoRelationMine && r.relationMiningService != nil {
			go func(uid, eid uint64) {
				result, err := r.relationMiningService.MineRelationsForEntry(uid, eid)
				if err != nil {
					fmt.Printf("[AutoRelationMine] 为知识点 %d 挖掘关系失败: %v\n", eid, err)
				} else {
					fmt.Printf("[AutoRelationMine] 为知识点 %d 挖掘到 %d 个关系\n", eid, result.RelationsFound)
				}
			}(userID, entry.ID)
		}
	}

	// 更新用户统计
	r.updateUserStats(userID)

	return entry, nil
}

// RemoveTaskKnowledge 按任务ID删除知识点
func (r *DefaultRAGService) RemoveTaskKnowledge(userID uint64, taskID uint64) error {
	db := database.GetDB()

	// 查找任务对应的知识条目
	var entry models.KnowledgeBaseEntry
	result := db.Where("user_id = ? AND task_id = ?", userID, taskID).First(&entry)
	if result.Error != nil {
		// 兼容旧数据
		result = db.Where("user_id = ? AND source_type = 1 AND source_id = ?", userID, taskID).First(&entry)
	}

	if result.Error != nil {
		return nil // 不存在则不需要删除
	}

	return r.RemoveDocument(entry.ID)
}

// RemoveDocument 删除文档（硬删除）
func (r *DefaultRAGService) RemoveDocument(entryID uint64) error {
	db := database.GetDB()

	// 使用 Unscoped 进行硬删除，确保数据从数据库中彻底删除

	// 删除向量缓存
	db.Unscoped().Where("entry_id = ?", entryID).Delete(&models.KnowledgeVectorCache{})

	// 删除关系
	db.Unscoped().Where("source_entry_id = ? OR target_entry_id = ?", entryID, entryID).
		Delete(&models.KnowledgeRelation{})

	// 删除条目
	if err := db.Unscoped().Delete(&models.KnowledgeBaseEntry{}, entryID).Error; err != nil {
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

// generateSummary 生成内容摘要，长度不超过 100 字
func generateSummary(content string) string {
	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		// 降级策略：截取前200字符
		if len(content) <= 200 {
			return content
		}
		runeContent := []rune(content)
		if len(runeContent) <= 200 {
			return content
		}
		return string(runeContent[:200]) + "..."
	}

	// 构造 AI 请求
	qwenURL := os.Getenv("QWEN_CHAT_URL")
	if qwenURL == "" {
		qwenURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	}

	prompt := fmt.Sprintf("请为以下内容生成一个简短的摘要，不超过100字，抓住核心要点：\n\n%s", content)

	reqBody := map[string]interface{}{
		"model": "qwen-plus",
		"input": map[string]interface{}{
			"messages": []map[string]string{
				{"role": "user", "content": prompt},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", qwenURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("调用 AI 生成摘要失败: %v\n", err)
		return fallbackSummary(content)
	}
	defer resp.Body.Close()

	var result struct {
		Output struct {
			Text string `json:"text"`
		} `json:"output"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("解析 AI 摘要响应失败: %v\n", err)
		return fallbackSummary(content)
	}

	if result.Output.Text == "" {
		return fallbackSummary(content)
	}

	return strings.TrimSpace(result.Output.Text)
}

func fallbackSummary(content string) string {
	if len(content) <= 200 {
		return content
	}
	runeContent := []rune(content)
	if len(runeContent) <= 200 {
		return content
	}
	return string(runeContent[:200]) + "..."
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

// classifyContent 使用AI语义分类内容
// 返回：category(学科大类), subCategory(细分领域), subject(能力维度)
func classifyContent(title, content string) (category, subCategory, subject string) {
	// 调用AI进行语义分类
	result := aiClassifyContent(title, content)
	return result.Category, result.SubCategory, result.Subject
}

// aiClassifyContent 调用AI进行智能分类
func aiClassifyContent(title, content string) AIClassificationResult {
	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		// 无API Key时使用简单规则降级
		return fallbackClassify(title, content)
	}

	// 构造AI请求
	qwenURL := os.Getenv("QWEN_CHAT_URL")
	if qwenURL == "" {
		qwenURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	}

	// 截取内容避免过长
	truncatedContent := content
	if len([]rune(content)) > 500 {
		truncatedContent = string([]rune(content)[:500]) + "..."
	}

	// 构造Prompt - 强制AI从枚举值中选择
	prompt := fmt.Sprintf(`你是一个专业的学习内容分类管理员。请根据以下标题和内容，进行智能分类。

【重要规则】
1. 如果内容是文学、语言学、比较文学等人文类，即使出现"索引"、"数据"等词，也必须归类为"人文社科"，严禁归类为"计算机"
2. 只有当内容明确讨论编程、数据库技术、软件开发时，才归类为"计算机"
3. 请根据内容的核心主题判断，不要被个别词汇误导

【分类要求】
请严格从以下选项中选择，不要自己编造：

category（学科大类，必须从以下选一个）：
%s

subject（能力维度，必须从以下选一个）：
%s

sub_category（细分领域）：
根据内容自动生成，限制在2-6个字，如"现当代文学"、"数据库"、"微积分"等

【待分类内容】
标题：%s
内容：%s

【输出格式】
请直接输出JSON，不要有任何其他文字：
{"category": "xxx", "sub_category": "xxx", "subject": "xxx"}`,
		strings.Join(CategoryEnum, "、"),
		strings.Join(SubjectEnum, "、"),
		title,
		truncatedContent,
	)

	reqBody := map[string]interface{}{
		"model": "qwen-plus",
		"input": map[string]interface{}{
			"messages": []map[string]string{
				{"role": "user", "content": prompt},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", qwenURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("创建AI分类请求失败: %v\n", err)
		return fallbackClassify(title, content)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("调用AI分类失败: %v\n", err)
		return fallbackClassify(title, content)
	}
	defer resp.Body.Close()

	var apiResult struct {
		Output struct {
			Text string `json:"text"`
		} `json:"output"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		fmt.Printf("解析AI分类响应失败: %v\n", err)
		return fallbackClassify(title, content)
	}

	if apiResult.Output.Text == "" {
		return fallbackClassify(title, content)
	}

	// 解析AI返回的JSON
	var classResult AIClassificationResult
	// 清理可能的markdown代码块标记
	jsonStr := strings.TrimSpace(apiResult.Output.Text)
	jsonStr = strings.TrimPrefix(jsonStr, "```json")
	jsonStr = strings.TrimPrefix(jsonStr, "```")
	jsonStr = strings.TrimSuffix(jsonStr, "```")
	jsonStr = strings.TrimSpace(jsonStr)

	if err := json.Unmarshal([]byte(jsonStr), &classResult); err != nil {
		fmt.Printf("解析AI分类JSON失败: %v, 原始内容: %s\n", err, jsonStr)
		return fallbackClassify(title, content)
	}

	// 验证并修正分类结果
	classResult = validateAndFixClassification(classResult)

	return classResult
}

// validateAndFixClassification 验证并修正分类结果，确保在枚举范围内
func validateAndFixClassification(result AIClassificationResult) AIClassificationResult {
	// 验证Category
	validCategory := false
	for _, c := range CategoryEnum {
		if result.Category == c {
			validCategory = true
			break
		}
	}
	if !validCategory {
		result.Category = "人文社科" // 默认
	}

	// 验证Subject
	validSubject := false
	for _, s := range SubjectEnum {
		if result.Subject == s {
			validSubject = true
			break
		}
	}
	if !validSubject {
		result.Subject = "理论素养" // 默认
	}

	// SubCategory如果为空，设置默认值
	if result.SubCategory == "" {
		result.SubCategory = "综合"
	}
	// 限制SubCategory长度
	if len([]rune(result.SubCategory)) > 10 {
		result.SubCategory = string([]rune(result.SubCategory)[:10])
	}

	return result
}

// fallbackClassify 降级分类（当AI不可用时使用简单规则）
func fallbackClassify(title, content string) AIClassificationResult {
	fullText := strings.ToLower(title + " " + content)

	// 简化的规则匹配，按优先级
	type simpleRule struct {
		category string
		subject  string
		keywords []string
	}

	rules := []simpleRule{
		// 人文社科优先（避免误分类）
		{"人文社科", "沟通表达", []string{"文学", "文言文", "诗词", "小说", "散文", "戏剧", "比较文学", "语文", "阅读", "写作"}},
		{"人文社科", "沟通表达", []string{"英语", "英文", "翻译", "语法", "单词", "口语", "听力"}},
		{"人文社科", "理论素养", []string{"历史", "朝代", "战争", "革命", "哲学", "政治", "思想"}},
		// 计算机类
		{"计算机", "实操应用", []string{"编程", "代码", "程序", "开发", "python", "java", "javascript", "golang"}},
		{"计算机", "实操应用", []string{"数据库", "sql", "mysql", "postgresql", "mongodb"}},
		{"计算机", "逻辑思维", []string{"数据结构", "算法", "链表", "二叉树", "排序", "动态规划"}},
		{"计算机", "实操应用", []string{"操作系统", "进程", "线程", "内存管理", "文件系统"}},
		{"计算机", "实操应用", []string{"计算机网络", "tcp", "http", "网络协议", "socket"}},
		// 数理逻辑
		{"数理逻辑", "逻辑思维", []string{"数学", "微积分", "线性代数", "概率论", "矩阵", "方程", "几何"}},
		{"数理逻辑", "逻辑思维", []string{"物理", "力学", "电磁", "量子", "热力学", "光学"}},
		// 自然科学
		{"自然科学", "理论素养", []string{"化学", "元素", "分子", "反应", "有机", "无机"}},
		{"自然科学", "理论素养", []string{"生物", "细胞", "基因", "遗传", "生态"}},
		{"自然科学", "理论素养", []string{"地理", "气候", "地形", "环境"}},
		// 经济管理
		{"经济管理", "逻辑思维", []string{"经济", "金融", "会计", "管理", "市场", "投资", "财务"}},
		// 艺术体育
		{"艺术体育", "创新思维", []string{"艺术", "绘画", "美术", "设计", "音乐", "舞蹈"}},
		{"艺术体育", "实操应用", []string{"体育", "运动", "健身", "训练"}},
	}

	for _, rule := range rules {
		for _, keyword := range rule.keywords {
			if strings.Contains(fullText, keyword) {
				return AIClassificationResult{
					Category:    rule.category,
					SubCategory: keyword,
					Subject:     rule.subject,
				}
			}
		}
	}

	// 默认分类
	return AIClassificationResult{
		Category:    "人文社科",
		SubCategory: "综合",
		Subject:     "理论素养",
	}
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

// GetKnowledgeGraph 获取用户知识图谱数据（聚类拓扑结构）
// GetKnowledgeGraph 获取知识图谱数据
func (r *DefaultRAGService) GetKnowledgeGraph(userID uint64, teamID *uint64) (*KnowledgeGraphData, error) {
	db := database.GetDB()

	// 获取知识条目查询
	query := db.Where("user_id = ?", userID)
	if teamID != nil {
		query = query.Where("team_id = ?", *teamID)
	} else {
		// 个人图谱：如果不传TeamID，是否应该只看个人的？或者全部？
		// 根据需求，个人任务和团队任务分开。
		// 如果 teamID == nil, 假设是个人图谱，通常包含所有个人条目，或者排除团队条目？
		// 现有逻辑是 WHERE user_id = ?，包含所有。为了区分，这里可以加上 team_id IS NULL 吗？
		// 暂时保持 user_id = ? 包含所有（或者根据产品逻辑，个人图谱是否包含团队任务的分配？通常是包含的）
		// 但用户现在的需求是 "团队任务的知识图谱应该独立于个人任务"。
		// 所以，如果 teamID != nil，只查该 Team。
		// 如果 teamID == nil，为了不混淆，我们可能需要排除 team_id != nil 的吗？
		// 用户的原话："团队任务的知识图谱和个人任务知识图谱是一样的，但是实际上团队任务的知识图谱应该是独立于个人任务的"
		// 这意味着 Team Graph 应该只显示 Team Data。Personal Graph 可能是显示 Personal Data (OR all data).
		// 现在的修改重点是让 Team Graph 独立。

		// 暂时保持 teamID == nil 时查所有 belongs to user (Current behavior),
		// 或者根据需求改为 team_id IS NULL。
		// 既然用户抱怨一样，那说明 Team Graph 里混入了 Personal Data (or vice versa).
		// 其实之前是因为根本没有 filter TeamID。
		// 只要 Team Graph 加上 filter，它就会只显示 Team 的。

		// 但 Personal Graph 是否要排除 Team 的？
		// 通常 "Use Personal Tasks" 可能意味着 exclude team stuff.
		// 但是 models.Task 有 TaskType (1=personal, 2=team).
		// KnowledgeBaseEntry 也有 TeamID.

		// 稍微改一下 query 逻辑：
		// 如果 teamID != nil, explicitly filter by team_id
		// 如果 teamID == nil， 暂时保持查询所有 user_id 的（因为个人可能也会把自己在团队里的贡献视为己有），
		// 或者仅查询 team_id IS NULL。
		// 考虑到用户抱怨 "不能用个人的那些任务创建"，说明 Team Graph 应该纯粹一点。
		// Personal Graph 如果包含 Team 的也许还好？
		// 让我们先只处理 TeamID != nil 的情况。
	}

	var entries []models.KnowledgeBaseEntry
	if err := query.Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("获取知识条目失败: %w", err)
	}

	categoryColorMap := getCategoryColorMap()
	nodes := make([]KnowledgeGraphNode, 0)
	links := make([]KnowledgeGraphLink, 0)
	linkSet := make(map[string]bool)

	// ===== 1. 创建6个虚拟中心节点（分类聚类中心）=====
	virtualNodes := createVirtualCenterNodes(categoryColorMap)
	nodes = append(nodes, virtualNodes...)

	// 构建虚拟节点ID映射 (category -> virtualID)
	virtualIDMap := make(map[string]uint64)
	for _, vn := range virtualNodes {
		virtualIDMap[vn.Category] = vn.ID
	}

	// ===== 2. 添加真实知识点节点 =====
	entryMap := make(map[uint64]models.KnowledgeBaseEntry)
	for _, entry := range entries {
		entryMap[entry.ID] = entry
		color := categoryColorMap[entry.Category]
		if color == "" {
			color = "#9ca3af"
		}

		// 计算节点大小：基础20 + 浏览次数*2，最大60
		symbolSize := 20 + entry.ViewCount*2
		if symbolSize > 60 {
			symbolSize = 60
		}

		nodes = append(nodes, KnowledgeGraphNode{
			ID:         entry.ID,
			Name:       truncateString(entry.Title, 20),
			Category:   entry.Category,
			Level:      entry.Level,
			Value:      entry.ViewCount + 10,
			Color:      color,
			IsVirtual:  false,
			SymbolSize: symbolSize,
		})

		// ===== 3. 创建归属连接（真实节点 -> 虚拟中心）=====
		if virtualID, ok := virtualIDMap[entry.Category]; ok {
			key := fmt.Sprintf("%d-%d", entry.ID, virtualID)
			if !linkSet[key] {
				linkSet[key] = true
				links = append(links, KnowledgeGraphLink{
					Source:       entry.ID,
					Target:       virtualID,
					RelationType: 6, // 归属中心
					Strength:     0.1,
					Label:        "",
				})
			}
		}
	}

	// ===== 4. 添加数据库中的显式关系（AI挖掘的逻辑关系）=====
	var relations []models.KnowledgeRelation
	db.Where("user_id = ?", userID).Find(&relations)

	relationLabels := map[int8]string{
		1: "前置",
		2: "相关",
		3: "扩展",
		4: "冲突",
	}

	for _, rel := range relations {
		// 确保两个节点都存在
		if _, ok := entryMap[rel.SourceEntryID]; !ok {
			continue
		}
		if _, ok := entryMap[rel.TargetEntryID]; !ok {
			continue
		}

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

	// ===== 5. 构建分类信息（用于前端图例）=====
	categories := make([]GraphCategory, 0)
	for name, color := range categoryColorMap {
		if name != "未分类" && name != "其他" {
			categories = append(categories, GraphCategory{
				Name:  name,
				Color: color,
			})
		}
	}

	return &KnowledgeGraphData{
		Nodes:      nodes,
		Links:      links,
		Categories: categories,
	}, nil
}

// createVirtualCenterNodes 创建虚拟中心节点（6大学科分类）
func createVirtualCenterNodes(colorMap map[string]string) []KnowledgeGraphNode {
	// 虚拟节点使用特殊的ID范围（避免与真实ID冲突）
	// 使用 9000000001 - 9000000006 作为虚拟节点ID
	virtualCategories := []struct {
		ID       uint64
		Name     string
		Category string
	}{
		{9000000001, "💻 计算机", "计算机"},
		{9000000002, "📚 人文社科", "人文社科"},
		{9000000003, "🔢 数理逻辑", "数理逻辑"},
		{9000000004, "🔬 自然科学", "自然科学"},
		{9000000005, "💰 经济管理", "经济管理"},
		{9000000006, "🎨 艺术体育", "艺术体育"},
	}

	nodes := make([]KnowledgeGraphNode, 0, len(virtualCategories))
	for _, vc := range virtualCategories {
		color := colorMap[vc.Category]
		if color == "" {
			color = "#9ca3af"
		}
		nodes = append(nodes, KnowledgeGraphNode{
			ID:         vc.ID,
			Name:       vc.Name,
			Category:   vc.Category,
			Level:      0,
			Value:      1000, // 超大节点，位于引力中心
			Color:      color,
			IsVirtual:  true,
			SymbolSize: 100, // 中心节点显示更大，确保能覆盖文字
		})
	}
	return nodes
}

// getCategoryColorMap 获取分类颜色映射（适配新的学科大类）
func getCategoryColorMap() map[string]string {
	return map[string]string{
		// 新的学科大类（环形图内圈）
		"计算机":  "#3b82f6", // 蓝色
		"人文社科": "#f59e0b", // 橙色
		"数理逻辑": "#8b5cf6", // 紫色
		"自然科学": "#10b981", // 绿色
		"经济管理": "#ef4444", // 红色
		"艺术体育": "#ec4899", // 粉色
		// 其他/未分类
		"未分类": "#9ca3af", // 灰色
		"其他":  "#9ca3af", // 灰色
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

// ReclassifyAllEntries 批量重分类用户的所有知识点（数据清洗）
// 返回成功处理的条目数量
func (r *DefaultRAGService) ReclassifyAllEntries(userID uint64) (int, error) {
	db := database.GetDB()

	// 获取用户所有知识条目
	var entries []models.KnowledgeBaseEntry
	if err := db.Where("user_id = ? AND status = 1", userID).Find(&entries).Error; err != nil {
		return 0, fmt.Errorf("获取知识条目失败: %w", err)
	}

	successCount := 0
	for _, entry := range entries {
		// 重新分类
		category, subCategory, subject := classifyContent(entry.Title, entry.Content)

		// 获取新的显示配置
		displayConfig := GetDisplayConfigForCategory(category)

		// 更新数据库
		updates := map[string]interface{}{
			"category":      category,
			"sub_category":  subCategory,
			"subject":       subject,
			"display_color": displayConfig.Color,
			"display_icon":  displayConfig.Icon,
		}

		if err := db.Model(&entry).Updates(updates).Error; err != nil {
			fmt.Printf("更新知识点 %d 分类失败: %v\n", entry.ID, err)
			continue
		}

		successCount++
		fmt.Printf("已重分类知识点 %d: %s -> 大类:%s, 细分:%s, 能力:%s\n",
			entry.ID, entry.Title, category, subCategory, subject)

		// 添加短暂延迟，避免API限流
		time.Sleep(100 * time.Millisecond)
	}

	// 更新用户统计
	r.updateUserStats(userID)

	return successCount, nil
}
