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

	// 团队知识库
	kb.GET("/team/list", listTeamKnowledge)
	kb.GET("/team/stats", getTeamKnowledgeStats)

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

	// AI关系挖掘
	kb.POST("/mine-relations", mineKnowledgeRelations)        // 为单个知识点挖掘关系
	kb.POST("/mine-all-relations", mineAllKnowledgeRelations) // 批量挖掘所有关系

	// RAG 问答（带引用溯源）
	kb.POST("/chat", ragChat)

	// 数据清洗：批量重分类知识点
	kb.POST("/reclassify", reclassifyKnowledgeEntries)
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

	db := database.GetDB()

	// 【防环检测】对于定向关系（前置/扩展），检查是否存在反向关系
	if req.RelationType == 1 || req.RelationType == 3 {
		var reverseCount int64
		db.Model(&models.KnowledgeRelation{}).
			Where("user_id = ? AND source_entry_id = ? AND target_entry_id = ? AND relation_type IN (1, 3)",
				userID.(uint64), req.TargetEntryID, req.SourceEntryID).
			Count(&reverseCount)

		if reverseCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "无法创建关系：已存在反向的定向关系，会形成互指环",
				"code":  400,
			})
			return
		}
	}

	// 检查是否已存在相同关系
	var existingCount int64
	db.Model(&models.KnowledgeRelation{}).
		Where("user_id = ? AND source_entry_id = ? AND target_entry_id = ?",
			userID.(uint64), req.SourceEntryID, req.TargetEntryID).
		Count(&existingCount)

	if existingCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "关系已存在",
			"code":  400,
		})
		return
	}

	relation := &models.KnowledgeRelation{
		UserID:        userID.(uint64),
		SourceEntryID: req.SourceEntryID,
		TargetEntryID: req.TargetEntryID,
		RelationType:  req.RelationType,
		Strength:      req.Strength,
	}

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

	// 获取 TeamID 参数
	var teamID *uint64
	teamIDStr := c.Query("team_id")
	if teamIDStr != "" {
		tid, err := strconv.ParseUint(teamIDStr, 10, 64)
		if err == nil {
			teamID = &tid
		}
	}

	// 添加调试信息：查询不限 status 的条目数
	db := database.GetDB()
	var totalCount, publishedCount int64

	query := db.Model(&models.KnowledgeBaseEntry{}).Where("user_id = ?", userID.(uint64))
	if teamID != nil {
		query = query.Where("team_id = ?", *teamID)
	}

	query.Count(&totalCount)
	query.Where("status = 1").Count(&publishedCount)

	fmt.Printf("[DEBUG] GetKnowledgeGraph - userID: %d, teamID: %v, totalEntries: %d, publishedEntries(status=1): %d\n",
		userID.(uint64), teamID, totalCount, publishedCount)

	graphData, err := ragService.GetKnowledgeGraph(userID.(uint64), teamID)
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
			"team_id":           teamID,
			"total_entries":     totalCount,
			"published_entries": publishedCount,
		},
	})
}

// mineKnowledgeRelations 为单个知识点挖掘关系
// @Summary AI自动挖掘知识点关系
// @Description 使用向量检索+AI推理，为指定知识点挖掘与其他知识点的逻辑关系
// @Tags Knowledge Base
// @Accept json
// @Produce json
// @Router /knowledge-base/mine-relations [post]
func mineKnowledgeRelations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req struct {
		EntryID uint64 `json:"entry_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供知识点ID"})
		return
	}

	// 初始化关系挖掘服务
	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("DASHSCOPE_API_KEY")
	}
	embeddingService := rag.NewQwenEmbeddingService(apiKey)
	miningService := rag.NewRelationMiningService(embeddingService)

	// 执行关系挖掘
	result, err := miningService.MineRelationsForEntry(userID.(uint64), req.EntryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  -1,
			"error": "挖掘关系失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
		"msg":  fmt.Sprintf("成功挖掘 %d 个关系", result.RelationsFound),
	})
}

// mineAllKnowledgeRelations 批量挖掘所有知识点的关系
// @Summary 批量挖掘所有知识点关系
// @Description 为用户的所有知识点执行AI关系挖掘，建立知识图谱的逻辑连线
// @Tags Knowledge Base
// @Accept json
// @Produce json
// @Router /knowledge-base/mine-all-relations [post]
func mineAllKnowledgeRelations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 初始化关系挖掘服务
	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("DASHSCOPE_API_KEY")
	}
	embeddingService := rag.NewQwenEmbeddingService(apiKey)
	miningService := rag.NewRelationMiningService(embeddingService)

	// 异步执行批量挖掘
	go func(uid uint64) {
		totalRelations, err := miningService.MineAllRelations(uid)
		if err != nil {
			fmt.Printf("[RelationMining] 批量挖掘失败: %v\n", err)
			return
		}
		fmt.Printf("[RelationMining] 用户 %d 批量挖掘完成，共发现 %d 个关系\n", uid, totalRelations)
	}(userID.(uint64))

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"status": "processing",
		},
		"msg": "批量关系挖掘已在后台开始，请稍后刷新知识图谱查看结果",
	})
}

// reclassifyKnowledgeEntries 批量重分类知识点（数据清洗）
// @Summary 批量重分类知识点
// @Description 使用AI重新分类用户的所有知识点，修复分类错误
// @Tags Knowledge Base
// @Accept json
// @Produce json
// @Router /knowledge-base/reclassify [post]
func reclassifyKnowledgeEntries(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 调用RAG服务的批量重分类方法
	successCount, err := ragService.ReclassifyAllEntries(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  -1,
			"error": "重分类失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"success_count": successCount,
		},
		"msg": fmt.Sprintf("成功重分类 %d 个知识点", successCount),
	})
}

// ragChat RAG问答（带引用溯源）
// 优化版流程：
// 第一步：Query理解 - AI判断问题所属领域，缩小检索范围
// 第二步：漏斗式检索 - 分类过滤 → 向量检索 → 关键词补漏
// 第三步：防幻觉生成 - 严格基于参考资料回答
// 第四步：引用溯源 - 返回带编号的引用信息
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

	// ========== 第一步：Query Understanding（问题理解）==========
	// 用AI判断用户问题属于哪个知识领域，缩小检索范围
	queryCategory := classifyUserQuery(req.Query)
	fmt.Printf("[RAG] 问题分类结果: %s, 原问题: %s\n", queryCategory, req.Query)

	// ========== 第二步：Funnel Retrieval（漏斗式检索）==========
	var searchResults []rag.SearchResult

	// 第一层：硬过滤 - 基于分类缩小范围
	var candidateEntries []models.KnowledgeBaseEntry
	if queryCategory != "" && queryCategory != "其他" && queryCategory != "未分类" {
		// 有明确分类，先按分类过滤
		db.Where("user_id = ? AND status = 1 AND category = ?", uid, queryCategory).
			Order("level DESC, view_count DESC").
			Find(&candidateEntries)
		fmt.Printf("[RAG] 分类过滤后候选: %d 条\n", len(candidateEntries))
	}

	// 如果分类过滤结果太少，扩大到全库搜索
	if len(candidateEntries) < 3 {
		db.Where("user_id = ? AND status = 1", uid).
			Order("level DESC, view_count DESC").
			Find(&candidateEntries)
		fmt.Printf("[RAG] 全库搜索候选: %d 条\n", len(candidateEntries))
	}

	// 第二层：软检索 - 向量相似度搜索（使用高阈值 0.60）
	highThreshold := float32(0.60)
	// minThreshold := float32(0.40) // 备用

	hybridResults, err := hybridSearchService.Search(uid, req.Query, req.Limit*2, 0.6)
	if err == nil && len(hybridResults) > 0 {
		for _, hr := range hybridResults {
			// 使用高阈值过滤，只保留真正相关的内容
			if hr.FinalScore >= highThreshold {
				// 如果有分类过滤，进一步检查分类匹配
				if queryCategory != "" && queryCategory != "其他" && queryCategory != "未分类" {
					if hr.Entry.Category == queryCategory {
						searchResults = append(searchResults, rag.SearchResult{
							Entry:      hr.Entry,
							Similarity: hr.FinalScore,
						})
					}
				} else {
					searchResults = append(searchResults, rag.SearchResult{
						Entry:      hr.Entry,
						Similarity: hr.FinalScore,
					})
				}
			}
		}
		fmt.Printf("[RAG] 高阈值(%.2f)向量搜索结果: %d 条\n", highThreshold, len(searchResults))
	}

	// 第三层：关键词补漏 - 如果向量搜索结果不足，用关键词搜索补充
	if len(searchResults) < 3 {
		keywords := extractQueryKeywords(req.Query)
		keywordEntries := smartKeywordSearch(db, uid, keywords, req.Limit)

		// 去重添加
		existingIDs := make(map[uint64]bool)
		for _, sr := range searchResults {
			existingIDs[sr.Entry.ID] = true
		}

		for _, entry := range keywordEntries {
			if !existingIDs[entry.ID] {
				// 如果有分类过滤，优先匹配同分类的
				if queryCategory != "" && queryCategory != "其他" && entry.Category == queryCategory {
					searchResults = append(searchResults, rag.SearchResult{
						Entry:      entry,
						Similarity: 0.45, // 关键词+分类匹配
					})
				} else if len(searchResults) < req.Limit {
					searchResults = append(searchResults, rag.SearchResult{
						Entry:      entry,
						Similarity: 0.35, // 纯关键词匹配，较低相似度
					})
				}
				existingIDs[entry.ID] = true
			}
		}
		fmt.Printf("[RAG] 关键词补漏后结果: %d 条\n", len(searchResults))
	}

	// 限制结果数量
	if len(searchResults) > req.Limit {
		searchResults = searchResults[:req.Limit]
	}

	// ========== 第三步：构建引用信息 ==========
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
		content := result.Entry.Content
		if len([]rune(content)) > 500 {
			content = string([]rune(content)[:500]) + "..."
		}
		contextParts = append(contextParts,
			fmt.Sprintf("[%d] 标题：%s\n分类：%s\n内容：%s",
				i+1, result.Entry.Title, result.Entry.Category, content))
	}

	// ========== 第四步：防幻觉生成（Grounded Generation）==========
	answer, err := generateGroundedRAGAnswer(req.Query, contextParts, len(searchResults) > 0)
	if err != nil || answer == "" {
		// 降级：基于知识库生成结构化回答
		var allEntries []models.KnowledgeBaseEntry
		db.Where("user_id = ? AND status = 1", uid).Limit(20).Find(&allEntries)
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

// classifyUserQuery 使用AI对用户问题进行分类（Query Understanding）
func classifyUserQuery(query string) string {
	apiKey := getQwenAPIKey()
	if apiKey == "" {
		// 无API Key时使用关键词规则降级
		return fallbackQueryClassify(query)
	}

	// 构造分类Prompt
	categories := []string{"计算机", "人文社科", "数理逻辑", "自然科学", "经济管理", "艺术体育", "其他"}
	prompt := fmt.Sprintf(`请判断用户问题属于以下哪个知识库大类。

可选分类：%s

用户问题：%s

【要求】
1. 只返回一个分类名称，不要有任何其他文字
2. 如果问题是通用性的（如"总结我学过的内容"），返回"其他"
3. 根据问题的核心主题判断，不要被个别词汇误导

分类结果：`, strings.Join(categories, "、"), query)

	reqBody := QwenRequest{
		Model: qwenFastModel(), // 用轻量模型，快速分类
		Messages: []QwenMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	apiURL := qwenChatURL()

	httpReq, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("[RAG] Query分类请求失败: %v\n", err)
		return fallbackQueryClassify(query)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var qwenResp QwenResponse
	if err := json.Unmarshal(body, &qwenResp); err != nil {
		return fallbackQueryClassify(query)
	}

	if len(qwenResp.Choices) == 0 {
		return fallbackQueryClassify(query)
	}

	result := strings.TrimSpace(qwenResp.Choices[0].Message.Content)
	// 验证结果是否在枚举中
	for _, cat := range categories {
		if result == cat {
			return result
		}
	}
	return "其他"
}

// fallbackQueryClassify 降级的问题分类（基于关键词）
func fallbackQueryClassify(query string) string {
	q := strings.ToLower(query)

	// 通用问题检测
	if strings.Contains(q, "总结") || strings.Contains(q, "学过") ||
		strings.Contains(q, "概览") || strings.Contains(q, "有哪些") {
		return "其他"
	}

	// 计算机类
	if strings.Contains(q, "编程") || strings.Contains(q, "代码") ||
		strings.Contains(q, "算法") || strings.Contains(q, "数据库") ||
		strings.Contains(q, "python") || strings.Contains(q, "java") ||
		strings.Contains(q, "前端") || strings.Contains(q, "后端") {
		return "计算机"
	}

	// 人文社科
	if strings.Contains(q, "文学") || strings.Contains(q, "历史") ||
		strings.Contains(q, "哲学") || strings.Contains(q, "政治") ||
		strings.Contains(q, "语文") || strings.Contains(q, "英语") {
		return "人文社科"
	}

	// 数理逻辑
	if strings.Contains(q, "数学") || strings.Contains(q, "物理") ||
		strings.Contains(q, "公式") || strings.Contains(q, "定理") {
		return "数理逻辑"
	}

	// 自然科学
	if strings.Contains(q, "化学") || strings.Contains(q, "生物") ||
		strings.Contains(q, "地理") {
		return "自然科学"
	}

	// 经济管理
	if strings.Contains(q, "经济") || strings.Contains(q, "金融") ||
		strings.Contains(q, "管理") || strings.Contains(q, "会计") {
		return "经济管理"
	}

	// 艺术体育
	if strings.Contains(q, "艺术") || strings.Contains(q, "音乐") ||
		strings.Contains(q, "体育") || strings.Contains(q, "设计") {
		return "艺术体育"
	}

	return "其他"
}

// generateGroundedRAGAnswer 防幻觉的RAG回答生成
// 严格要求AI只能基于参考资料回答，无相关资料时诚实拒绝
func generateGroundedRAGAnswer(query string, contextParts []string, hasRelevantKnowledge bool) (string, error) {
	apiKey := getQwenAPIKey()
	if apiKey == "" {
		return "", nil
	}

	var prompt string
	if hasRelevantKnowledge {
		// 有相关知识点时的Prompt
		prompt = fmt.Sprintf(`你是"智学空间"的AI助教。请严格基于以下【参考资料】回答用户问题。

【参考资料】
%s

【用户问题】
%s

【回答要求】
1. 引用标注：回答时，必须在相关内容后标注信息来源，如"...根据[1]的内容..."或"...[1]"
2. 拒绝编造：如果【参考资料】中没有包含问题答案的相关信息，请直接回答："抱歉，我的知识库中暂时没有收录相关内容。建议您补充相关知识点到知识库。"
3. 严禁使用你自己的训练数据编造答案，只能使用上面提供的参考资料
4. 保持客观：回答要简洁有条理，使用中文，可适当使用Markdown格式
5. 如果参考资料只是部分相关，请明确说明哪些内容来自知识库，哪些是你的补充建议

请回答：`, strings.Join(contextParts, "\n\n"), query)
	} else {
		// 无相关知识点时的Prompt
		prompt = fmt.Sprintf(`你是"智学空间"的AI助教。用户问了一个问题，但在他的个人知识库中没有找到相关内容。

【用户问题】
%s

【回答要求】
1. 首先明确告知：「抱歉，在您的知识库中暂时没有收录与"[问题关键词]"相关的内容。」
2. 然后可以简要提供一些通用的学习建议（2-3句话即可）
3. 最后建议用户将相关知识添加到知识库，以便下次查询
4. 不要长篇大论，保持简洁

请回答：`, query)
	}

	reqBody := QwenRequest{
		Model: qwenChatModel(),
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
		Model: qwenChatModel(),
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

// listTeamKnowledge 获取团队知识库列�?
func listTeamKnowledge(c *gin.Context) {
	teamIDStr := c.Query("team_id")
	if teamIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供team_id"})
		return
	}
	teamID, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的team_id"})
		return
	}

	db := database.GetDB()
	var entries []models.KnowledgeBaseEntry

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	category := c.Query("category")
	level := c.Query("level")
	search := c.Query("search")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := db.Where("team_id = ? AND status = 1", teamID)
	countQuery := db.Model(&models.KnowledgeBaseEntry{}).Where("team_id = ? AND status = 1", teamID)

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
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("title LIKE ? OR content LIKE ? OR summary LIKE ?", searchPattern, searchPattern, searchPattern)
		countQuery = countQuery.Where("title LIKE ? OR content LIKE ? OR summary LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("level DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取团队知识列表失败"})
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

// getTeamKnowledgeStats 获取团队知识库统�?
func getTeamKnowledgeStats(c *gin.Context) {
	teamIDStr := c.Query("team_id")
	if teamIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供team_id"})
		return
	}
	teamID, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的team_id"})
		return
	}

	db := database.GetDB()

	type LevelCount struct {
		Level int8
		Count int64
	}
	var levelCounts []LevelCount

	if err := db.Model(&models.KnowledgeBaseEntry{}).
		Select("level, count(*) as count").
		Where("team_id = ? AND status = 1", teamID).
		Group("level").
		Find(&levelCounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取团队统计失败"})
		return
	}

	stats := map[string]interface{}{
		"level_0_count": int64(0),
		"level_1_count": int64(0),
		"level_2_count": int64(0),
		"level_3_count": int64(0),
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

	// 查询需要复习的知识点（超过7天未复习且等级低�?�?
	var reviewCount int64
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	db.Model(&models.KnowledgeBaseEntry{}).
		Where("team_id = ? AND status = 1 AND level < 3 AND (last_review_at IS NULL OR last_review_at < ?)", teamID, sevenDaysAgo).
		Count(&reviewCount)
	stats["review_needed"] = reviewCount

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}
