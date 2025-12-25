package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/models"
	"learningAssistant-backend/services/rag"
)

// 初始化RAG相关的全局变量
var (
	ragService          rag.RAGService
	aiAnalysisService   *rag.AIAnalysisService
	hybridSearchService *rag.HybridSearchService
)

// initRAGServices 初始化RAG服务
func initRAGServices() {
	// 优先使用 Qwen Embedding API（真正的语义向量化）
	// 需要设置环境变量 QWEN_API_KEY 或 DASHSCOPE_API_KEY
	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("DASHSCOPE_API_KEY")
	}

	var embeddingService rag.EmbeddingService
	if apiKey != "" {
		// 使用 Qwen Embedding API - 真正的语义理解
		embeddingService = rag.NewQwenEmbeddingService(apiKey)
		fmt.Println("[RAG] 使用 Qwen Embedding API (text-embedding-v3)")
	} else {
		// 降级到本地简单 Embedding（基于字符特征，效果有限）
		embeddingService = rag.NewLocalEmbeddingService()
		fmt.Println("[RAG] 警告: 未配置 QWEN_API_KEY，使用本地 Embedding（效果有限）")
		fmt.Println("[RAG] 建议: 设置环境变量 QWEN_API_KEY 以获得更好的语义理解能力")
	}

	ragService = rag.NewRAGService(embeddingService)
	aiAnalysisService = rag.NewAIAnalysisService(apiKey)
	hybridSearchService = rag.NewHybridSearchService(embeddingService)
}

// registerKnowledgeBaseRoutes 注册知识库路由
func registerKnowledgeBaseRoutes(router *gin.RouterGroup) {
	// 初始化RAG服务
	if ragService == nil {
		initRAGServices()
	}

	kb := router.Group("/knowledge-base")
	kb.Use(middleware.AuthMiddleware())

	// 知识库管理
	kb.POST("/add", addKnowledgeEntry)
	kb.POST("/add-from-task", addKnowledgeFromTask)
	kb.POST("/add-from-note", addKnowledgeFromNote)
	kb.GET("/search", searchKnowledge)
	kb.GET("/entry/:id", getKnowledgeEntry)
	kb.PUT("/entry/:id/level", updateKnowledgeLevel)
	kb.DELETE("/entry/:id", deleteKnowledgeEntry)
	kb.GET("/stats", getUserKnowledgeStats)
	kb.GET("/list", listUserKnowledge)

	// AI分析
	kb.GET("/analysis", analyzeUserKnowledge)
	kb.GET("/distribution", getKnowledgeDistribution)
	kb.GET("/skill-radar", getSkillRadarData)
	kb.GET("/trends", getLearningTrends)

	// 知识关系
	kb.GET("/relations/:id", getKnowledgeRelations)
	kb.POST("/relations", createKnowledgeRelation)

	// 知识图谱
	kb.GET("/graph", getKnowledgeGraph)

	// RAG 问答（带引用溯源）
	kb.POST("/chat", ragChat)
}

// addKnowledgeEntry 添加知识库条目
// @Summary 添加知识库条目
// @Description 创建新的知识库条目
// @Tags Knowledge Base
// @Accept json
// @Produce json
// @Param request body object true "请求体"
// @Router /knowledge-base/add [post]
func addKnowledgeEntry(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req struct {
		Title    string   `json:"title" binding:"required"`
		Content  string   `json:"content" binding:"required"`
		Category string   `json:"category"`
		Tags     []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := ragService.AddDocument(userID.(uint64), 3, 0, req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": entry,
		"msg":  "知识库条目创建成功",
	})
}

// addKnowledgeFromTask 从任务创建知识库条目
func addKnowledgeFromTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req struct {
		TaskID  uint64 `json:"task_id" binding:"required"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取任务信息
	db := database.GetDB()
	var task models.Task
	if err := db.First(&task, req.TaskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	// 如果没有提供标题和内容，使用任务的信息
	title := req.Title
	if title == "" {
		title = task.Title
	}
	content := req.Content
	if content == "" {
		content = task.Description
	}

	entry, err := ragService.AddDocument(userID.(uint64), 1, req.TaskID, title, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": entry,
		"msg":  "从任务创建知识库条目成功",
	})
}

// addKnowledgeFromNote 从笔记创建知识库条目
func addKnowledgeFromNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req struct {
		NoteID uint64 `json:"note_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取笔记信息
	db := database.GetDB()
	var note models.StudyNote
	if err := db.First(&note, req.NoteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	entry, err := ragService.AddDocument(userID.(uint64), 2, req.NoteID, note.Title, note.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": entry,
		"msg":  "从笔记创建知识库条目成功",
	})
}

// searchKnowledge 搜索知识库（使用混合检索）
func searchKnowledge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	// 使用混合检索（向量 + BM25）
	hybridResults, err := hybridSearchService.Search(userID.(uint64), query, limit, 0.6)
	if err != nil {
		// 降级到原有检索
		entries, fallbackErr := ragService.SearchKnowledge(userID.(uint64), query, limit)
		if fallbackErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fallbackErr.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"results": entries,
				"total":   len(entries),
			},
			"msg": "搜索成功",
		})
		return
	}

	// 构建返回结果（包含匹配信息）
	type EnhancedResult struct {
		models.KnowledgeBaseEntry
		Score        float32  `json:"score"`
		MatchedTerms []string `json:"matched_terms,omitempty"`
		Highlight    string   `json:"highlight,omitempty"`
	}

	results := make([]EnhancedResult, 0, len(hybridResults))
	for _, hr := range hybridResults {
		results = append(results, EnhancedResult{
			KnowledgeBaseEntry: hr.Entry,
			Score:              hr.FinalScore,
			MatchedTerms:       hr.MatchedTerms,
			Highlight:          hr.MatchHighlight,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"results": results,
			"total":   len(results),
		},
		"msg": "搜索成功",
	})
}

// getKnowledgeEntry 获取单个知识库条目
func getKnowledgeEntry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	db := database.GetDB()
	var entry models.KnowledgeBaseEntry
	if err := db.First(&entry, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "条目不存在"})
		return
	}

	// 更新浏览次数
	db.Model(&entry).Update("view_count", entry.ViewCount+1)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": entry,
	})
}

// updateKnowledgeLevel 更新知识点掌握等级
func updateKnowledgeLevel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req struct {
		Level *int8 `json:"level"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Level == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "等级不能为空"})
		return
	}

	if *req.Level < 0 || *req.Level > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "等级必须在0-4之间"})
		return
	}

	if err := ragService.UpdateKnowledgeLevel(id, *req.Level); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "等级更新成功",
	})
}

// deleteKnowledgeEntry 删除知识库条目
func deleteKnowledgeEntry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := ragService.RemoveDocument(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}

// getUserKnowledgeStats 获取用户知识库统计
func getUserKnowledgeStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	stats, err := ragService.GetUserKnowledgeStats(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}

// listUserKnowledge 列表用户知识库
func listUserKnowledge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	db := database.GetDB()
	var entries []models.KnowledgeBaseEntry

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	category := c.Query("category")
	level := c.Query("level")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := db.Where("user_id = ? AND status = 1", userID.(uint64))

	// 构建筛选条件的查询（用于计算总数）
	countQuery := db.Model(&models.KnowledgeBaseEntry{}).Where("user_id = ? AND status = 1", userID.(uint64))

	if category != "" {
		query = query.Where("category = ?", category)
		countQuery = countQuery.Where("category = ?", category)
	}
	if level != "" {
		levelInt, err := strconv.Atoi(level)
		if err == nil {
			query = query.Where("level = ?", int8(levelInt))
			countQuery = countQuery.Where("level = ?", int8(levelInt))
		}
	}

	offset := (page - 1) * pageSize
	if err := query.Order("level DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	var total int64
	countQuery.Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":     entries,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// analyzeUserKnowledge AI分析用户知识库
func analyzeUserKnowledge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	report, err := aiAnalysisService.AnalyzeUserKnowledge(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": report,
	})
}

// getKnowledgeDistribution 获取知识点分布
func getKnowledgeDistribution(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	report, err := aiAnalysisService.AnalyzeUserKnowledge(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": report.KnowledgeDistribution,
	})
}

// getSkillRadarData 获取技能雷达数据
func getSkillRadarData(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	report, err := aiAnalysisService.AnalyzeUserKnowledge(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": report.SkillRadar,
	})
}

// getLearningTrends 获取学习趋势
func getLearningTrends(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// range: 30 | 90 | year
	rangeKey := strings.TrimSpace(c.DefaultQuery("range", "30"))
	now := time.Now()
	from := now.AddDate(0, 0, -29)
	granularity := rag.TrendGranularityDay
	if rangeKey == "90" {
		from = now.AddDate(0, 0, -89)
		granularity = rag.TrendGranularityWeek
	} else if rangeKey == "year" {
		y := now.Year()
		from = time.Date(y, 1, 1, 0, 0, 0, 0, now.Location())
		// 到本月月底不需要，直接到今天；展示时仍是 12 个月桶（后端会补零到当前月）
		granularity = rag.TrendGranularityMonth
	}

	// 这里不走 AnalyzeUserKnowledge，避免为趋势多做一次全量分析
	db := database.GetDB()
	var entries []models.KnowledgeBaseEntry
	if err := db.Where("user_id = ? AND status = 1 AND created_at >= ? AND created_at <= ?", userID.(uint64), from, now).
		Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	trends := aiAnalysisService.AnalyzeUserLearningTrendsRange(userID.(uint64), entries, from, now, granularity)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": trends,
	})
}

// getKnowledgeRelations 获取知识关系
func getKnowledgeRelations(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	relations, err := ragService.GetKnowledgeRelations(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": relations,
	})
}

// createKnowledgeRelation 创建知识关系
func createKnowledgeRelation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req struct {
		SourceEntryID uint64  `json:"source_entry_id" binding:"required"`
		TargetEntryID uint64  `json:"target_entry_id" binding:"required"`
		RelationType  int8    `json:"relation_type" binding:"required,min=1,max=4"`
		Strength      float32 `json:"strength"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	relation := &models.KnowledgeRelation{
		UserID:        userID.(uint64),
		SourceEntryID: req.SourceEntryID,
		TargetEntryID: req.TargetEntryID,
		RelationType:  req.RelationType,
		Strength:      req.Strength,
	}

	db := database.GetDB()
	if err := db.Create(relation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建关系失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": relation,
		"msg":  "关系创建成功",
	})
}

// getKnowledgeGraph 获取知识图谱数据
func getKnowledgeGraph(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 添加调试信息：查询不限 status 的条目数
	db := database.GetDB()
	var totalCount, publishedCount int64
	db.Model(&models.KnowledgeBaseEntry{}).Where("user_id = ?", userID.(uint64)).Count(&totalCount)
	db.Model(&models.KnowledgeBaseEntry{}).Where("user_id = ? AND status = 1", userID.(uint64)).Count(&publishedCount)

	fmt.Printf("[DEBUG] GetKnowledgeGraph - userID: %d, totalEntries: %d, publishedEntries(status=1): %d\n",
		userID.(uint64), totalCount, publishedCount)

	graphData, err := ragService.GetKnowledgeGraph(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("[DEBUG] GetKnowledgeGraph - nodes: %d, links: %d\n",
		len(graphData.Nodes), len(graphData.Links))

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": graphData,
		"debug": gin.H{
			"user_id":           userID.(uint64),
			"total_entries":     totalCount,
			"published_entries": publishedCount,
		},
	})
}

// ragChat RAG问答（带引用溯源）
// 业界标准流程：Query理解 → 向量检索 → 上下文组装 → LLM生成 → 返回结果+引用
func ragChat(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req struct {
		Query string `json:"query" binding:"required"`
		Limit int    `json:"limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入问题"})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 5
	}

	db := database.GetDB()
	uid := userID.(uint64)

	// 步骤1: 获取用户所有知识点（用于上下文）
	var allEntries []models.KnowledgeBaseEntry
	db.Where("user_id = ? AND status = 1", uid).
		Order("level DESC, view_count DESC").
		Limit(20). // 最多取20条作为候选
		Find(&allEntries)

	// 步骤2: 使用混合检索（向量 + BM25）- 业界标准方案
	hybridResults, err := hybridSearchService.Search(uid, req.Query, req.Limit, 0.6) // 向量权重60%

	// 将混合检索结果转换为 SearchResult 格式
	// 只保留相似度 >= 35% 的结果，过滤掉不相关内容
	var searchResults []rag.SearchResult
	if err == nil && len(hybridResults) > 0 {
		for _, hr := range hybridResults {
			// 只添加相关性足够高的结果
			if hr.FinalScore >= 0.35 {
				searchResults = append(searchResults, rag.SearchResult{
					Entry:      hr.Entry,
					Similarity: hr.FinalScore,
				})
			}
		}
	}

	// 降级: 如果混合检索无结果，使用智能关键词匹配
	// 但不强行返回低相关性结果
	if len(searchResults) == 0 {
		keywords := extractQueryKeywords(req.Query)
		relevantEntries := smartKeywordSearch(db, uid, keywords, req.Limit)
		for _, entry := range relevantEntries {
			// 关键词匹配给较低的相似度，但仍显示
			searchResults = append(searchResults, rag.SearchResult{
				Entry:      entry,
				Similarity: 0.4, // 关键词匹配默认 40%
			})
		}
	}

	// 步骤3: 构建引用信息
	citations := make([]rag.Citation, 0, len(searchResults))
	contextParts := make([]string, 0, len(searchResults))
	for i, result := range searchResults {
		citations = append(citations, rag.Citation{
			ID:         result.Entry.ID,
			Title:      result.Entry.Title,
			Category:   result.Entry.Category,
			Summary:    result.Entry.Summary,
			Similarity: result.Similarity,
		})
		// 构建上下文片段（编号便于引用）
		content := result.Entry.Summary
		if content == "" {
			content = truncateContent(result.Entry.Content, 300)
		}
		contextParts = append(contextParts,
			"["+strconv.Itoa(i+1)+"] 《"+result.Entry.Title+"》("+result.Entry.Category+")\n"+content)
	}

	// 步骤4: 构建知识库概览（即使没有精确匹配，也让AI知道用户学了什么）
	knowledgeOverview := buildKnowledgeOverview(allEntries)

	// 步骤5: 调用AI生成回答（真正的RAG）
	answer, err := generateEnhancedRAGAnswer(req.Query, contextParts, knowledgeOverview, len(allEntries))
	if err != nil || answer == "" {
		// 降级：基于知识库生成结构化回答
		answer = generateSmartFallbackAnswer(req.Query, searchResults, allEntries)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rag.RAGQueryResult{
			Answer:    answer,
			Citations: citations,
			Query:     req.Query,
		},
	})
}

// extractQueryKeywords 从问题中提取关键词
func extractQueryKeywords(query string) []string {
	// 移除常见的疑问词和停用词
	stopWords := []string{
		"帮我", "请", "一下", "什么", "怎么", "如何", "为什么", "哪些", "哪个",
		"总结", "介绍", "说明", "告诉", "给我", "我想", "我要", "知道",
		"的", "了", "吗", "呢", "啊", "是", "在", "有", "和", "与",
	}

	result := query
	for _, word := range stopWords {
		result = strings.ReplaceAll(result, word, " ")
	}

	// 分词
	words := strings.Fields(result)
	keywords := make([]string, 0)
	for _, w := range words {
		w = strings.TrimSpace(w)
		if len(w) >= 2 { // 至少2个字符
			keywords = append(keywords, w)
		}
	}

	return keywords
}

// smartKeywordSearch 智能关键词搜索
func smartKeywordSearch(db *gorm.DB, userID uint64, keywords []string, limit int) []models.KnowledgeBaseEntry {
	var entries []models.KnowledgeBaseEntry

	if len(keywords) == 0 {
		// 没有关键词，返回最近的知识点
		db.Where("user_id = ? AND status = 1", userID).
			Order("created_at DESC").
			Limit(limit).
			Find(&entries)
		return entries
	}

	// 构建 OR 查询条件
	query := db.Where("user_id = ? AND status = 1", userID)

	// 对每个关键词进行模糊匹配
	var conditions []string
	var args []interface{}
	for _, kw := range keywords {
		pattern := "%" + kw + "%"
		conditions = append(conditions, "(title LIKE ? OR content LIKE ? OR category LIKE ? OR keywords LIKE ?)")
		args = append(args, pattern, pattern, pattern, pattern)
	}

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	query.Order("level DESC, view_count DESC").
		Limit(limit).
		Find(&entries)

	return entries
}

// buildKnowledgeOverview 构建知识库概览
func buildKnowledgeOverview(entries []models.KnowledgeBaseEntry) string {
	if len(entries) == 0 {
		return "用户知识库暂无内容。"
	}

	// 按分类统计
	categoryCount := make(map[string]int)
	for _, e := range entries {
		categoryCount[e.Category]++
	}

	var parts []string
	for cat, count := range categoryCount {
		parts = append(parts, cat+"("+strconv.Itoa(count)+"个)")
	}

	return "用户知识库共有" + strconv.Itoa(len(entries)) + "个知识点，涵盖：" + strings.Join(parts, "、")
}

// truncateContent 截断内容
func truncateContent(content string, maxLen int) string {
	runes := []rune(content)
	if len(runes) <= maxLen {
		return content
	}
	return string(runes[:maxLen]) + "..."
}

// generateEnhancedRAGAnswer 增强版RAG回答生成
func generateEnhancedRAGAnswer(query string, contextParts []string, knowledgeOverview string, totalKnowledge int) (string, error) {
	apiKey := getQwenAPIKey()
	if apiKey == "" {
		return "", nil // 没有API Key，使用降级策略
	}

	context := ""
	hasRelevantKnowledge := len(contextParts) > 0
	if hasRelevantKnowledge {
		context = "【相关知识点】\n" + strings.Join(contextParts, "\n\n")
	} else {
		context = "【重要提示】在用户的知识库中没有找到与问题「" + query + "」直接相关的内容。请根据下面的要求诚实回答。"
	}

	prompt := `你是一个智能学习助手，负责基于用户的个人知识库回答问题。

【用户知识库概况】
` + knowledgeOverview + `

` + context + `

【用户问题】
` + query + `

【回答要求】
1. 如果找到相关知识点（上面有【相关知识点】部分），请基于知识点内容回答，并在回答中引用编号如"根据[1]..."
2. 如果没有找到直接相关的知识点：
   - 首先明确告知用户"在您的知识库中暂未找到与此问题直接相关的内容"
   - 然后提供简要的通用知识帮助（如果你知道的话）
   - 最后建议用户补充相关知识到知识库
3. 对于"总结学过的内容"这类问题，请基于知识库概况给出分析
4. 回答要有条理，使用中文，适当使用 Markdown 格式
5. 不要编造用户知识库中不存在的内容

请回答：`

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

	client := &http.Client{Timeout: 60 * time.Second} // 增加超时时间
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var qwenResp QwenResponse
	if err := json.Unmarshal(body, &qwenResp); err != nil {
		return "", err
	}

	if qwenResp.Error != nil {
		return "", fmt.Errorf("AI API error: %s", qwenResp.Error.Message)
	}

	if len(qwenResp.Choices) == 0 {
		return "", nil
	}

	return qwenResp.Choices[0].Message.Content, nil
}

// generateSmartFallbackAnswer 智能降级回答
func generateSmartFallbackAnswer(query string, results []rag.SearchResult, allEntries []models.KnowledgeBaseEntry) string {
	var sb strings.Builder

	// 判断问题类型
	isSummaryQuery := strings.Contains(query, "总结") ||
		strings.Contains(query, "学过") ||
		strings.Contains(query, "概览") ||
		strings.Contains(query, "有哪些")

	if len(allEntries) == 0 {
		sb.WriteString("📚 您的知识库暂时是空的。\n\n")
		sb.WriteString("建议您：\n")
		sb.WriteString("1. 在完成学习任务时，将重要内容添加到知识库\n")
		sb.WriteString("2. 使用「同步知识库」功能导入已有的笔记和任务\n")
		sb.WriteString("3. 手动添加学习心得和知识点\n")
		return sb.String()
	}

	if isSummaryQuery {
		// 总结类问题：展示知识库概览
		sb.WriteString("📊 **您的知识库概览**\n\n")

		// 按分类统计
		categoryEntries := make(map[string][]models.KnowledgeBaseEntry)
		for _, e := range allEntries {
			categoryEntries[e.Category] = append(categoryEntries[e.Category], e)
		}

		sb.WriteString("您共积累了 **" + strconv.Itoa(len(allEntries)) + "** 个知识点，分布如下：\n\n")

		for cat, entries := range categoryEntries {
			sb.WriteString("### " + cat + " (" + strconv.Itoa(len(entries)) + "个)\n")
			for i, e := range entries {
				if i >= 3 { // 每个分类最多显示3个
					sb.WriteString("- ...还有" + strconv.Itoa(len(entries)-3) + "个\n")
					break
				}
				sb.WriteString("- " + e.Title + "\n")
			}
			sb.WriteString("\n")
		}

		return sb.String()
	}

	if len(results) > 0 {
		sb.WriteString("📖 根据您的知识库，找到以下相关内容：\n\n")
		for i, result := range results {
			sb.WriteString(strconv.Itoa(i+1) + ". **" + result.Entry.Title + "**")
			if result.Entry.Category != "" {
				sb.WriteString(" [" + result.Entry.Category + "]")
			}
			sb.WriteString("\n")
			if result.Entry.Summary != "" {
				sb.WriteString("   " + result.Entry.Summary + "\n")
			}
			sb.WriteString("\n")
		}
		sb.WriteString("💡 点击上方的引用来源可查看详细内容。")
	} else {
		sb.WriteString("🔍 在您的知识库中暂未找到与「" + query + "」直接相关的内容。\n\n")
		sb.WriteString("您可以：\n")
		sb.WriteString("1. 尝试使用不同的关键词\n")
		sb.WriteString("2. 将相关知识添加到知识库中\n")
	}

	return sb.String()
}
