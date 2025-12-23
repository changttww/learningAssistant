package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"learningAssistant-backend/database"
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/models"
	"learningAssistant-backend/services/rag"
)

// 初始化RAG相关的全局变量
var (
	ragService        rag.RAGService
	aiAnalysisService *rag.AIAnalysisService
)

// initRAGServices 初始化RAG服务
func initRAGServices() {
	// 使用本地embedding服务或Qwen API
	embeddingService := rag.NewLocalEmbeddingService()
	// embeddingService := rag.NewQwenEmbeddingService("") // 若有API key可用

	ragService = rag.NewRAGService(embeddingService)
	aiAnalysisService = rag.NewAIAnalysisService("")
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

// searchKnowledge 搜索知识库
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

	entries, err := ragService.SearchKnowledge(userID.(uint64), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	report, err := aiAnalysisService.AnalyzeUserKnowledge(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": report.LearningTrends,
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
