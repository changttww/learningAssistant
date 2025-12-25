package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"learningAssistant-backend/database"
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/models"
)

// registerNoteEnhanceRoutes 注册智能笔记增强路由
func registerNoteEnhanceRoutes(router *gin.RouterGroup) {
	router.POST("/enhance", middleware.AuthMiddleware(), handleEnhanceNote)
	router.POST("/generate-summary", middleware.AuthMiddleware(), handleGenerateSummary)
	router.POST("/extract-keywords", middleware.AuthMiddleware(), handleExtractKeywords)
	router.POST("/generate-mindmap", middleware.AuthMiddleware(), handleGenerateMindmap)
	router.POST("/generate-questions", middleware.AuthMiddleware(), handleGenerateQuestions)
	router.POST("/polish", middleware.AuthMiddleware(), handlePolishNote)
}

// NoteEnhanceRequest 笔记增强请求
type NoteEnhanceRequest struct {
	NoteID  uint64 `json:"note_id"`  // 笔记ID（可选，如果提供则从数据库读取）
	Content string `json:"content"`  // 笔记内容
	Title   string `json:"title"`    // 笔记标题
	Type    string `json:"type"`     // 增强类型：all/summary/keywords/mindmap/questions/polish
}

// NoteEnhanceResponse 笔记增强响应
type NoteEnhanceResponse struct {
	Summary    *NoteSummary    `json:"summary,omitempty"`
	Keywords   *NoteKeywords   `json:"keywords,omitempty"`
	Mindmap    *NoteMindmap    `json:"mindmap,omitempty"`
	Questions  *NoteQuestions  `json:"questions,omitempty"`
	Polish     *NotePolish     `json:"polish,omitempty"`
}

// NoteSummary 笔记摘要
type NoteSummary struct {
	Brief       string   `json:"brief"`        // 一句话总结
	KeyPoints   []string `json:"key_points"`   // 核心要点（3-5条）
	Conclusion  string   `json:"conclusion"`   // 结论/总结
	ReadingTime int      `json:"reading_time"` // 预计阅读时间（分钟）
}

// NoteKeywords 关键词提取
type NoteKeywords struct {
	MainKeywords   []KeywordItem `json:"main_keywords"`   // 主要关键词
	RelatedConcepts []string     `json:"related_concepts"` // 相关概念
	Category       string        `json:"category"`         // 自动分类
	Tags           []string      `json:"tags"`             // 推荐标签
}

// KeywordItem 关键词项
type KeywordItem struct {
	Word       string  `json:"word"`
	Weight     float64 `json:"weight"`     // 重要性权重 0-1
	Definition string  `json:"definition"` // 简要定义
}

// NoteMindmap 思维导图
type NoteMindmap struct {
	Root     MindmapNode `json:"root"`
	Markdown string      `json:"markdown"` // Markdown 格式的思维导图
}

// MindmapNode 思维导图节点
type MindmapNode struct {
	Text     string        `json:"text"`
	Children []MindmapNode `json:"children,omitempty"`
	Level    int           `json:"level"`
}

// NoteQuestions 知识点问题
type NoteQuestions struct {
	ReviewQuestions []ReviewQuestion `json:"review_questions"` // 复习问题
	ThinkingQuestions []string       `json:"thinking_questions"` // 思考题
}

// ReviewQuestion 复习问题
type ReviewQuestion struct {
	Question   string   `json:"question"`
	Answer     string   `json:"answer"`
	Difficulty string   `json:"difficulty"` // easy/medium/hard
	Type       string   `json:"type"`       // concept/application/analysis
}

// NotePolish 笔记润色
type NotePolish struct {
	PolishedContent string     `json:"polished_content"` // 润色后的内容
	Improvements    []string   `json:"improvements"`     // 改进说明
	Structure       []Section  `json:"structure"`        // 优化后的结构
	Suggestions     []string   `json:"suggestions"`      // 进一步建议
}

// Section 章节
type Section struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Order   int    `json:"order"`
}

// handleEnhanceNote 全面增强笔记
func handleEnhanceNote(c *gin.Context) {
	var req NoteEnhanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 如果提供了笔记ID，从数据库获取内容
	if req.NoteID > 0 {
		var note models.StudyNote
		if err := database.GetDB().First(&note, req.NoteID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "笔记不存在"})
			return
		}
		req.Content = note.Content
		req.Title = note.Title
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "笔记内容不能为空"})
		return
	}

	enhanceType := req.Type
	if enhanceType == "" {
		enhanceType = "all"
	}

	response := &NoteEnhanceResponse{}

	apiKey := getQwenAPIKey()
	if apiKey == "" {
		// 无 AI，使用基础处理
		response = generateBasicEnhancement(req.Title, req.Content, enhanceType)
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": response})
		return
	}

	// 根据类型调用不同的增强功能
	switch enhanceType {
	case "all":
		response.Summary, _ = generateAISummary(apiKey, req.Title, req.Content)
		response.Keywords, _ = generateAIKeywords(apiKey, req.Title, req.Content)
		response.Mindmap, _ = generateAIMindmap(apiKey, req.Title, req.Content)
		response.Questions, _ = generateAIQuestions(apiKey, req.Title, req.Content)
		response.Polish, _ = generateAIPolish(apiKey, req.Title, req.Content)
	case "summary":
		response.Summary, _ = generateAISummary(apiKey, req.Title, req.Content)
	case "keywords":
		response.Keywords, _ = generateAIKeywords(apiKey, req.Title, req.Content)
	case "mindmap":
		response.Mindmap, _ = generateAIMindmap(apiKey, req.Title, req.Content)
	case "questions":
		response.Questions, _ = generateAIQuestions(apiKey, req.Title, req.Content)
	case "polish":
		response.Polish, _ = generateAIPolish(apiKey, req.Title, req.Content)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}

// handleGenerateSummary 生成摘要
func handleGenerateSummary(c *gin.Context) {
	var req NoteEnhanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "内容不能为空"})
		return
	}

	apiKey := getQwenAPIKey()
	if apiKey == "" {
		summary := generateBasicSummary(req.Title, req.Content)
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": summary})
		return
	}

	summary, err := generateAISummary(apiKey, req.Title, req.Content)
	if err != nil {
		summary = generateBasicSummary(req.Title, req.Content)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": summary})
}

// handleExtractKeywords 提取关键词
func handleExtractKeywords(c *gin.Context) {
	var req NoteEnhanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "内容不能为空"})
		return
	}

	apiKey := getQwenAPIKey()
	if apiKey == "" {
		keywords := generateBasicKeywords(req.Content)
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": keywords})
		return
	}

	keywords, err := generateAIKeywords(apiKey, req.Title, req.Content)
	if err != nil {
		keywords = generateBasicKeywords(req.Content)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": keywords})
}

// handleGenerateMindmap 生成思维导图
func handleGenerateMindmap(c *gin.Context) {
	var req NoteEnhanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "内容不能为空"})
		return
	}

	apiKey := getQwenAPIKey()
	if apiKey == "" {
		mindmap := generateBasicMindmap(req.Title, req.Content)
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": mindmap})
		return
	}

	mindmap, err := generateAIMindmap(apiKey, req.Title, req.Content)
	if err != nil {
		mindmap = generateBasicMindmap(req.Title, req.Content)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": mindmap})
}

// handleGenerateQuestions 生成复习问题
func handleGenerateQuestions(c *gin.Context) {
	var req NoteEnhanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "内容不能为空"})
		return
	}

	apiKey := getQwenAPIKey()
	if apiKey == "" {
		questions := generateBasicQuestions(req.Content)
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": questions})
		return
	}

	questions, err := generateAIQuestions(apiKey, req.Title, req.Content)
	if err != nil {
		questions = generateBasicQuestions(req.Content)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": questions})
}

// handlePolishNote 润色笔记
func handlePolishNote(c *gin.Context) {
	var req NoteEnhanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "内容不能为空"})
		return
	}

	apiKey := getQwenAPIKey()
	if apiKey == "" {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": &NotePolish{
			PolishedContent: req.Content,
			Improvements:    []string{"需要配置 AI 服务才能使用润色功能"},
		}})
		return
	}

	polish, err := generateAIPolish(apiKey, req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "润色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": polish})
}

// ============ AI 生成函数 ============

// generateAISummary 使用 AI 生成摘要
func generateAISummary(apiKey, title, content string) (*NoteSummary, error) {
	prompt := fmt.Sprintf(`请为以下学习笔记生成摘要。

标题: %s

内容:
%s

请严格按照以下 JSON 格式返回（只返回纯 JSON，不要 markdown 代码块）：
{
  "brief": "一句话总结（不超过50字）",
  "key_points": ["核心要点1", "核心要点2", "核心要点3"],
  "conclusion": "结论或总结（100字以内）",
  "reading_time": 5
}

要求：
1. brief: 用一句话概括笔记的核心内容
2. key_points: 提取3-5个最重要的知识点
3. conclusion: 总结学习这些内容的意义和价值
4. reading_time: 预计阅读原文需要的分钟数`, title, content)

	result, err := callQwenForNoteEnhance(apiKey, prompt)
	if err != nil {
		return nil, err
	}

	var summary NoteSummary
	if err := json.Unmarshal([]byte(result), &summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

// generateAIKeywords 使用 AI 提取关键词
func generateAIKeywords(apiKey, title, content string) (*NoteKeywords, error) {
	prompt := fmt.Sprintf(`请为以下学习笔记提取关键词和概念。

标题: %s

内容:
%s

请严格按照以下 JSON 格式返回（只返回纯 JSON，不要 markdown 代码块）：
{
  "main_keywords": [
    {"word": "关键词", "weight": 0.9, "definition": "简要定义"}
  ],
  "related_concepts": ["相关概念1", "相关概念2"],
  "category": "学科分类",
  "tags": ["标签1", "标签2", "标签3"]
}

要求：
1. main_keywords: 提取5-8个核心关键词，weight表示重要性(0-1)
2. related_concepts: 与本笔记相关的延伸概念
3. category: 判断笔记属于哪个学科领域
4. tags: 推荐3-5个适合的标签`, title, content)

	result, err := callQwenForNoteEnhance(apiKey, prompt)
	if err != nil {
		return nil, err
	}

	var keywords NoteKeywords
	if err := json.Unmarshal([]byte(result), &keywords); err != nil {
		return nil, err
	}

	return &keywords, nil
}

// generateAIMindmap 使用 AI 生成思维导图
func generateAIMindmap(apiKey, title, content string) (*NoteMindmap, error) {
	prompt := fmt.Sprintf(`请为以下学习笔记生成思维导图结构。

标题: %s

内容:
%s

请严格按照以下 JSON 格式返回（只返回纯 JSON，不要 markdown 代码块）：
{
  "root": {
    "text": "中心主题",
    "level": 0,
    "children": [
      {
        "text": "一级分支1",
        "level": 1,
        "children": [
          {"text": "二级分支1", "level": 2, "children": []},
          {"text": "二级分支2", "level": 2, "children": []}
        ]
      }
    ]
  },
  "markdown": "# 中心主题\n## 一级分支1\n- 二级分支1\n- 二级分支2"
}

要求：
1. 提取笔记的层次结构，构建清晰的思维导图
2. 中心主题应该是笔记的核心概念
3. 一级分支是主要的知识模块
4. 二级分支是具体的知识点
5. markdown 字段提供 Markdown 格式的大纲`, title, content)

	result, err := callQwenForNoteEnhance(apiKey, prompt)
	if err != nil {
		return nil, err
	}

	var mindmap NoteMindmap
	if err := json.Unmarshal([]byte(result), &mindmap); err != nil {
		return nil, err
	}

	return &mindmap, nil
}

// generateAIQuestions 使用 AI 生成复习问题
func generateAIQuestions(apiKey, title, content string) (*NoteQuestions, error) {
	prompt := fmt.Sprintf(`请根据以下学习笔记生成复习问题。

标题: %s

内容:
%s

请严格按照以下 JSON 格式返回（只返回纯 JSON，不要 markdown 代码块）：
{
  "review_questions": [
    {
      "question": "问题描述",
      "answer": "参考答案",
      "difficulty": "easy/medium/hard",
      "type": "concept/application/analysis"
    }
  ],
  "thinking_questions": ["开放性思考题1", "开放性思考题2"]
}

要求：
1. review_questions: 生成5-8道复习题，覆盖主要知识点
   - concept: 概念理解题
   - application: 应用题
   - analysis: 分析题
2. thinking_questions: 2-3道开放性思考题，帮助深入理解`, title, content)

	result, err := callQwenForNoteEnhance(apiKey, prompt)
	if err != nil {
		return nil, err
	}

	var questions NoteQuestions
	if err := json.Unmarshal([]byte(result), &questions); err != nil {
		return nil, err
	}

	return &questions, nil
}

// generateAIPolish 使用 AI 润色笔记
func generateAIPolish(apiKey, title, content string) (*NotePolish, error) {
	prompt := fmt.Sprintf(`请润色并优化以下学习笔记。

标题: %s

原始内容:
%s

请严格按照以下 JSON 格式返回（只返回纯 JSON，不要 markdown 代码块）：
{
  "polished_content": "润色后的完整内容（使用Markdown格式）",
  "improvements": ["改进点1", "改进点2", "改进点3"],
  "structure": [
    {"title": "章节标题", "content": "章节内容", "order": 1}
  ],
  "suggestions": ["进一步建议1", "进一步建议2"]
}

要求：
1. polished_content: 润色后的笔记，保持原意的同时优化表达和结构
2. improvements: 说明做了哪些改进
3. structure: 建议的章节结构
4. suggestions: 进一步完善笔记的建议`, title, content)

	result, err := callQwenForNoteEnhance(apiKey, prompt)
	if err != nil {
		return nil, err
	}

	var polish NotePolish
	if err := json.Unmarshal([]byte(result), &polish); err != nil {
		return nil, err
	}

	return &polish, nil
}

// callQwenForNoteEnhance 调用千问 API
func callQwenForNoteEnhance(apiKey, prompt string) (string, error) {
	reqBody := QwenRequest{
		Model: "qwen-plus",
		Messages: []QwenMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	apiURL := qwenChatURL()

	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("[note-enhance] AI response length: %d", len(body))

	var qwenResp QwenResponse
	if err := json.Unmarshal(body, &qwenResp); err != nil {
		return "", fmt.Errorf("解析响应失败")
	}

	if qwenResp.Error != nil {
		return "", fmt.Errorf("AI 错误: %s", qwenResp.Error.Message)
	}

	if len(qwenResp.Choices) == 0 {
		return "", fmt.Errorf("AI 返回为空")
	}

	content := qwenResp.Choices[0].Message.Content
	content = extractJSON(content)

	return content, nil
}

// ============ 基础处理函数（无 AI 降级方案）============

func generateBasicEnhancement(title, content, enhanceType string) *NoteEnhanceResponse {
	response := &NoteEnhanceResponse{}

	switch enhanceType {
	case "all":
		response.Summary = generateBasicSummary(title, content)
		response.Keywords = generateBasicKeywords(content)
		response.Mindmap = generateBasicMindmap(title, content)
		response.Questions = generateBasicQuestions(content)
		response.Polish = generateBasicPolish(title, content)
	case "summary":
		response.Summary = generateBasicSummary(title, content)
	case "keywords":
		response.Keywords = generateBasicKeywords(content)
	case "mindmap":
		response.Mindmap = generateBasicMindmap(title, content)
	case "questions":
		response.Questions = generateBasicQuestions(content)
	case "polish":
		response.Polish = generateBasicPolish(title, content)
	}

	return response
}

func generateBasicSummary(title, content string) *NoteSummary {
	// 简单截取作为摘要
	brief := content
	if len(brief) > 100 {
		brief = brief[:100] + "..."
	}

	// 按句子分割提取要点
	sentences := regexp.MustCompile(`[。！？.!?]`).Split(content, -1)
	var keyPoints []string
	for i, s := range sentences {
		s = strings.TrimSpace(s)
		if len(s) > 10 && i < 5 {
			keyPoints = append(keyPoints, s)
		}
	}

	readingTime := len(content) / 500 // 假设每分钟500字
	if readingTime < 1 {
		readingTime = 1
	}

	return &NoteSummary{
		Brief:       brief,
		KeyPoints:   keyPoints,
		Conclusion:  "建议复习上述要点以巩固理解。",
		ReadingTime: readingTime,
	}
}

func generateBasicKeywords(content string) *NoteKeywords {
	// 简单分词提取
	words := strings.Fields(content)
	wordCount := make(map[string]int)

	punctuation := "，。！？、；：\"\"''（）,.!?;:()[]{}'"
	for _, w := range words {
		w = strings.Trim(w, punctuation)
		if len(w) >= 2 {
			wordCount[w]++
		}
	}

	var keywords []KeywordItem
	for word, count := range wordCount {
		if count >= 2 && len(keywords) < 8 {
			keywords = append(keywords, KeywordItem{
				Word:       word,
				Weight:     float64(count) / float64(len(words)),
				Definition: "请参考笔记内容",
			})
		}
	}

	return &NoteKeywords{
		MainKeywords:    keywords,
		RelatedConcepts: []string{},
		Category:        "其他",
		Tags:            []string{"学习笔记"},
	}
}

func generateBasicMindmap(title, content string) *NoteMindmap {
	root := MindmapNode{
		Text:  title,
		Level: 0,
	}

	// 按段落分割
	paragraphs := strings.Split(content, "\n")
	for i, p := range paragraphs {
		p = strings.TrimSpace(p)
		if len(p) > 5 && i < 5 {
			child := MindmapNode{
				Text:  truncateText(p, 30),
				Level: 1,
			}
			root.Children = append(root.Children, child)
		}
	}

	// 生成 Markdown
	var md strings.Builder
	md.WriteString("# " + title + "\n")
	for _, child := range root.Children {
		md.WriteString("## " + child.Text + "\n")
	}

	return &NoteMindmap{
		Root:     root,
		Markdown: md.String(),
	}
}

func generateBasicQuestions(content string) *NoteQuestions {
	return &NoteQuestions{
		ReviewQuestions: []ReviewQuestion{
			{
				Question:   "请总结本笔记的主要内容？",
				Answer:     "请根据笔记内容进行总结。",
				Difficulty: "medium",
				Type:       "concept",
			},
			{
				Question:   "本笔记中最重要的知识点是什么？",
				Answer:     "请从笔记中提取关键知识点。",
				Difficulty: "easy",
				Type:       "concept",
			},
		},
		ThinkingQuestions: []string{
			"如何将这些知识应用到实际问题中？",
			"这些概念与其他知识有什么联系？",
		},
	}
}

func truncateText(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

// generateBasicPolish 基础润色（无AI时的降级方案）
func generateBasicPolish(title, content string) *NotePolish {
	// 简单的格式化处理
	polished := content
	
	// 分段处理
	paragraphs := strings.Split(content, "\n")
	var structure []Section
	order := 1
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if len(p) > 10 {
			structure = append(structure, Section{
				Title:   fmt.Sprintf("段落 %d", order),
				Content: p,
				Order:   order,
			})
			order++
		}
	}

	return &NotePolish{
		PolishedContent: polished,
		Improvements:    []string{"已保持原文格式", "建议使用AI润色获得更好效果"},
		Structure:       structure,
		Suggestions:     []string{"可以添加标题和小节", "建议补充具体示例"},
	}
}
