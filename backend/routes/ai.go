package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AI 解析请求结构
type ParseTaskRequest struct {
	Input string `json:"input" binding:"required"`
}

// AI 解析响应结构
type ParseTaskResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	StartTime   string `json:"startTime"`
	EndDate     string `json:"endDate"`
	EndTime     string `json:"endTime"`
	Category    string `json:"category"`
}

// 任务指导请求结构
type TaskGuidanceRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// 资源链接结构
type ResourceLink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"` // video/article/course/tool
}

// 任务指导响应结构
type TaskGuidanceResponse struct {
	Steps      []string       `json:"steps"`      // 分步骤指导
	Tips       []string       `json:"tips"`       // 学习技巧
	Resources  []ResourceLink `json:"resources"`  // 相关资源
	TimeAdvice string         `json:"timeAdvice"` // 时间建议
}

// 测验生成请求结构
type QuizGenerateRequest struct {
	Topic        string `json:"topic" binding:"required"` // 主题
	Content      string `json:"content"`                  // 学习内容/笔记（可选）
	Difficulty   string `json:"difficulty"`               // 难度：easy/medium/hard
	QuizCount    int    `json:"quizCount"`                // 题目数量
	IncludeEssay bool   `json:"includeEssay"`             // 是否包含简答题
}

// 选择题结构
type MultipleChoiceQuestion struct {
	Question      string            `json:"question"`      // 题目
	Options       map[string]string `json:"options"`       // 选项 {"A": "选项1", "B": "选项2"}
	CorrectAnswer string            `json:"correctAnswer"` // 正确答案(A/B/C/D)
	Explanation   string            `json:"explanation"`   // 答案解析
	Difficulty    string            `json:"difficulty"`    // 难度
}

// 简答题结构
type EssayQuestion struct {
	Question        string `json:"question"`        // 题目
	StudySuggestion string `json:"studySuggestion"` // 学习建议
}

// 测验响应结构
type QuizResponse struct {
	Questions     []MultipleChoiceQuestion `json:"questions"`               // 选择题列表
	EssayQuestion *EssayQuestion           `json:"essayQuestion,omitempty"` // 简答题(可选)
}

// Qwen API 请求结构
type QwenRequest struct {
	Model        string        `json:"model"`
	Messages     []QwenMessage `json:"messages"`
	EnableSearch bool          `json:"enable_search,omitempty"`
}

type QwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type QwenResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

// ParseTaskWithAI 使用 AI 解析自然语言任务
func ParseTaskWithAI(c *gin.Context) {
	var req ParseTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入任务描述"})
		return
	}

	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		fmt.Println("未配置 QWEN_API_KEY，使用本地解析")
		result := mockParseTask(req.Input)
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
		return
	}

	fmt.Printf("正在调用通义千问 API，输入: %s\n", req.Input)

	result, err := callQwenAPI(apiKey, req.Input)
	if err != nil {
		fmt.Printf("通义千问解析失败: %v, 降级到本地解析\n", err)
		result = mockParseTask(req.Input)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}

// GenerateQuiz 生成智能测验
func GenerateQuiz(c *gin.Context) {
	var req QuizGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入测验主题"})
		return
	}

	// 设置默认值
	if req.Difficulty == "" {
		req.Difficulty = "medium"
	}
	if req.QuizCount == 0 {
		req.QuizCount = 3
	}

	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		fmt.Println("未配置 QWEN_API_KEY，使用默认测验")
		result := mockQuiz(req.Topic, req.QuizCount, req.IncludeEssay)
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
		return
	}

	fmt.Printf("正在生成测验，主题: %s，难度: %s\n", req.Topic, req.Difficulty)

	result, err := callQwenForQuiz(apiKey, req)
	if err != nil {
		fmt.Printf("生成测验失败: %v, 使用默认测验\n", err)
		result = mockQuiz(req.Topic, req.QuizCount, req.IncludeEssay)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}

// GetTaskGuidance 获取任务指导和相关资源
func GetTaskGuidance(c *gin.Context) {
	var req TaskGuidanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入任务信息"})
		return
	}

	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		fmt.Println("未配置 QWEN_API_KEY，使用默认指导")
		result := mockTaskGuidance(req.Title, req.Category)
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
		return
	}

	fmt.Printf("正在获取任务指导，任务: %s\n", req.Title)

	result, err := callQwenForGuidance(apiKey, req.Title, req.Description, req.Category)
	if err != nil {
		fmt.Printf("获取指导失败: %v, 使用默认指导\n", err)
		result = mockTaskGuidance(req.Title, req.Category)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}

// callQwenForGuidance 调用通义千问获取任务指导
func callQwenForGuidance(apiKey, title, description, category string) (*TaskGuidanceResponse, error) {
	// --- 修改开始 ---
	// 核心改动：
	// 1. 明确要求 "先进行联网搜索"
	// 2. 强调 URL 必须是 "验证过的"
	// 3. 警告 "严禁编造"
	prompt := fmt.Sprintf(`你是一个专业的学习指导助手。请务必先【进行联网搜索】，查找与以下任务最匹配的最新学习资源，然后制定指导方案。

任务标题: %s
任务描述: %s
任务类型: %s

请严格按照以下 JSON 格式返回（只返回纯 JSON，不要 markdown，不要解释）：
{
  "steps": ["步骤1", "步骤2", "步骤3"],
  "tips": ["技巧1", "技巧2"],
  "resources": [
    {"title": "资源标题", "url": "真实URL", "type": "video/article/course/tool"}
  ],
  "timeAdvice": "时间安排建议"
}

重要要求：
1. Resources (资源):
   - 必须利用搜索功能找到 3-5 个【真实存在】且【目前可访问】的链接。
   - 优先选择 B站(bilibili)、GitHub、知乎高质量回答、官方文档、MOOC 平台的真实链接。
   - **严禁编造链接**，如果找不到特定资源就不要返回网址，给出资源名称即可。
   - 必须确保 URL 是完整的（以 http/https 开头）。

2. steps: 提供3-5个具体的执行步骤，每个步骤要详细可操作，记住一定要详细
3. tips: 提供2-3个实用的学习技巧或注意事项，一定要详细
4. 关于 resources (资源) 的特殊要求：
   1. 不要尝试生成具体的视频 ID (如 BV号) 或文章具体路径，因为这些极易失效。
   2. 请生成【搜索聚合页】的链接，确保用户一定能打开。
   
   格式示例：
   - 推荐B站教程，URL请填: https://search.bilibili.com/all?keyword=具体的搜索词
   - 推荐书籍/文档，URL请填: https://www.bing.com/search?q=具体的搜索词
   - 推荐GitHub项目，URL请填: https://github.com/search?q=具体的搜索词

5. timeAdvice: 给出合理的时间分配建议`, title, description, category)

	reqBody := QwenRequest{
		Model: "qwen-plus",
		Messages: []QwenMessage{
			{Role: "user", Content: prompt},
		},
		EnableSearch: true,
	}

	jsonData, _ := json.Marshal(reqBody)
	apiURL := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("任务指导响应: %s\n", string(body))

	var qwenResp QwenResponse
	if err := json.Unmarshal(body, &qwenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if qwenResp.Error != nil {
		return nil, fmt.Errorf("API 错误: %s", qwenResp.Error.Message)
	}

	if len(qwenResp.Choices) == 0 {
		return nil, fmt.Errorf("AI 返回内容为空")
	}

	content := qwenResp.Choices[0].Message.Content
	content = extractJSON(content)

	var result TaskGuidanceResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v, 内容: %s", err, content)
	}

	return &result, nil
}

// callQwenForQuiz 调用通义千问生成测验
func callQwenForQuiz(apiKey string, req QuizGenerateRequest) (*QuizResponse, error) {
	essayInstruction := ""
	if req.IncludeEssay {
		essayInstruction = `
  "essay": {
    "question": "简答题题目",
    "keyPoints": ["得分点1", "得分点2", "得分点3"],
    "referenceAnswer": "详细的参考答案",
    "difficulty": "medium"
  },`
	}

	contentContext := ""
	if req.Content != "" {
		contentContext = fmt.Sprintf("\n学习内容/笔记:\n%s\n", req.Content)
	}

	difficultyMap := map[string]string{
		"easy":   "简单（适合初学者）",
		"medium": "中等（需要理解和应用）",
		"hard":   "困难（需要深入分析和综合运用）",
	}
	difficultyDesc := difficultyMap[req.Difficulty]
	if difficultyDesc == "" {
		difficultyDesc = "中等"
	}

	prompt := fmt.Sprintf(`你是一位专业的教育测评专家。请为以下主题生成高质量的测验题目。

主题: %s
难度要求: %s
题目数量: %d 道选择题%s
%s
请严格按照以下 JSON 格式返回（只返回纯 JSON，不要 markdown，不要解释）：
{
  "topic": "%s",
  "multipleChoice": [
    {
      "question": "题目内容",
      "options": ["A. 选项1", "B. 选项2", "C. 选项3", "D. 选项4"],
      "correctAnswer": "A",
      "explanation": "详细的答案解析，说明为什么选这个答案",
      "difficulty": "easy/medium/hard"
    }
  ],%s
  "totalQuestions": %d,
  "estimatedTime": "预计用时",
  "studySuggestions": ["学习建议1", "学习建议2"]
}

重要要求：
1. 选择题要求：
   - 题目要精准，考查核心知识点
   - 4个选项要有迷惑性，避免明显错误
   - 正确答案用字母表示（A/B/C/D）
   - 解析要详细，帮助理解知识点
   - 难度要符合要求，由易到难排列

2. 简答题要求（如果包含）：
   - 题目要开放，考查理解和应用能力
   - 给出3-5个关键得分点
   - 提供完整的参考答案（150-200字）

3. 学习建议：
   - 根据题目内容，给出2-3条针对性的复习建议
   - 指出常见易错点

4. 预计用时：
   - 选择题每题约2-3分钟
   - 简答题约10-15分钟`, req.Topic,
		difficultyDesc,
		req.QuizCount,
		func() string {
			if req.IncludeEssay {
				return " + 1 道简答题"
			}
			return ""
		}(),
		contentContext,
		req.Topic,
		essayInstruction,
		func() int {
			if req.IncludeEssay {
				return req.QuizCount + 1
			}
			return req.QuizCount
		}())

	reqBody := QwenRequest{
		Model: "qwen-plus",
		Messages: []QwenMessage{
			{Role: "user", Content: prompt},
		},
		EnableSearch: false, // 测验生成不需要联网
	}

	jsonData, _ := json.Marshal(reqBody)
	apiURL := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

	httpReq, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("测验生成响应: %s\n", string(body))

	var qwenResp QwenResponse
	if err := json.Unmarshal(body, &qwenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if qwenResp.Error != nil {
		return nil, fmt.Errorf("API 错误: %s", qwenResp.Error.Message)
	}

	if len(qwenResp.Choices) == 0 {
		return nil, fmt.Errorf("AI 返回内容为空")
	}

	content := qwenResp.Choices[0].Message.Content
	content = extractJSON(content)

	// 先解析为临时结构
	var rawResult struct {
		Topic          string `json:"topic"`
		MultipleChoice []struct {
			Question      string   `json:"question"`
			Options       []string `json:"options"` // AI返回的是数组格式
			CorrectAnswer string   `json:"correctAnswer"`
			Explanation   string   `json:"explanation"`
			Difficulty    string   `json:"difficulty"`
		} `json:"multipleChoice"`
		Essay struct {
			Question        string   `json:"question"`
			KeyPoints       []string `json:"keyPoints"`
			ReferenceAnswer string   `json:"referenceAnswer"`
			Difficulty      string   `json:"difficulty"`
		} `json:"essay"`
		TotalQuestions   int      `json:"totalQuestions"`
		EstimatedTime    string   `json:"estimatedTime"`
		StudySuggestions []string `json:"studySuggestions"`
	}

	if err := json.Unmarshal([]byte(content), &rawResult); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v, 内容: %s", err, content)
	}

	// 转换为前端期望的格式
	result := &QuizResponse{
		Questions: make([]MultipleChoiceQuestion, len(rawResult.MultipleChoice)),
	}

	// 转换选择题
	for i, q := range rawResult.MultipleChoice {
		// 将数组格式的options转为map格式
		optionsMap := make(map[string]string)
		for _, opt := range q.Options {
			// 解析 "A. 选项内容" 格式
			if len(opt) >= 3 && opt[1] == '.' {
				key := string(opt[0])
				value := strings.TrimSpace(opt[2:])
				optionsMap[key] = value
			}
		}

		result.Questions[i] = MultipleChoiceQuestion{
			Question:      q.Question,
			Options:       optionsMap,
			CorrectAnswer: q.CorrectAnswer,
			Explanation:   q.Explanation,
			Difficulty:    q.Difficulty,
		}
	}

	// 转换简答题
	if rawResult.Essay.Question != "" {
		studySuggestion := rawResult.Essay.ReferenceAnswer
		if len(rawResult.Essay.KeyPoints) > 0 {
			studySuggestion = "关键点: " + strings.Join(rawResult.Essay.KeyPoints, "; ") + "\n\n" + rawResult.Essay.ReferenceAnswer
		}
		result.EssayQuestion = &EssayQuestion{
			Question:        rawResult.Essay.Question,
			StudySuggestion: studySuggestion,
		}
	}

	return result, nil
}

// mockTaskGuidance 本地生成默认指导
func mockTaskGuidance(title, category string) *TaskGuidanceResponse {
	result := &TaskGuidanceResponse{
		Steps: []string{
			"1. 明确任务目标和要求",
			"2. 收集相关学习资料",
			"3. 制定学习计划和时间表",
			"4. 开始执行并做好笔记",
			"5. 复习总结，检验学习效果",
		},
		Tips: []string{
			"保持专注，使用番茄工作法提高效率",
			"遇到问题及时记录，寻求帮助",
			"定期回顾，加深记忆",
		},
		Resources: []ResourceLink{
			{Title: "B站学习资源", URL: "https://www.bilibili.com", Type: "video"},
			{Title: "中国大学MOOC", URL: "https://www.icourse163.org", Type: "course"},
			{Title: "知乎学习专栏", URL: "https://www.zhihu.com", Type: "article"},
		},
		TimeAdvice: "建议每天安排1-2小时专注学习，分阶段完成任务",
	}

	// 根据分类提供更具体的资源
	switch category {
	case "study":
		result.Resources = []ResourceLink{
			{Title: "B站学习视频", URL: "https://www.bilibili.com/v/knowledge", Type: "video"},
			{Title: "中国大学MOOC", URL: "https://www.icourse163.org", Type: "course"},
			{Title: "学堂在线", URL: "https://www.xuetangx.com", Type: "course"},
		}
	case "project":
		result.Resources = []ResourceLink{
			{Title: "GitHub", URL: "https://github.com", Type: "tool"},
			{Title: "Stack Overflow", URL: "https://stackoverflow.com", Type: "article"},
			{Title: "掘金", URL: "https://juejin.cn", Type: "article"},
		}
	case "reading":
		result.Resources = []ResourceLink{
			{Title: "微信读书", URL: "https://weread.qq.com", Type: "tool"},
			{Title: "豆瓣读书", URL: "https://book.douban.com", Type: "article"},
			{Title: "得到App", URL: "https://www.dedao.cn", Type: "course"},
		}
	case "exam":
		result.Resources = []ResourceLink{
			{Title: "考试酷", URL: "https://www.examcoo.com", Type: "tool"},
			{Title: "刷题网站", URL: "https://www.nowcoder.com", Type: "tool"},
			{Title: "知乎备考经验", URL: "https://www.zhihu.com", Type: "article"},
		}
	}

	return result
}

const dateFormat = "2006-01-02"

func getTodayStr() string {
	return time.Now().Format(dateFormat)
}

func getTomorrowStr() string {
	return time.Now().AddDate(0, 0, 1).Format(dateFormat)
}

func getDayAfterTomorrowStr() string {
	return time.Now().AddDate(0, 0, 2).Format(dateFormat)
}

func callQwenAPI(apiKey, input string) (*ParseTaskResponse, error) {
	today := getTodayStr()
	tomorrow := getTomorrowStr()
	dayAfterTomorrow := getDayAfterTomorrowStr()

	prompt := fmt.Sprintf(`你是一个任务解析助手。请从用户的自然语言输入中提取任务信息,给用户输入的任务简介扩充为专业性的任务描述。

当前日期: %s
明天日期: %s  
后天日期: %s

用户输入: "%s"

请严格按照以下 JSON 格式返回（只返回纯 JSON，不要 markdown 代码块，不要任何其他文字）:
{"title":"任务标题","description":"任务描述","startDate":"YYYY-MM-DD","startTime":"HH:MM","endDate":"YYYY-MM-DD","endTime":"HH:MM","category":"分类"}

解析规则：
1. title: 提取核心任务名称
2. description: 生成详细描述，说明任务的具体内容和目标，任务描述务必给出详细指导，不少于100字
3. startDate: 开始日期，默认今天 %s
4. startTime: 开始时间，上午默认09:00，下午默认14:00，晚上默认19:00
5. endDate: 结束日期，默认等于开始日期
6. endTime: 结束时间，根据任务预估时长设定，学习类任务通常2-3小时
7. category: 只能是 study/exam/project/reading/other 之一`,
		today, tomorrow, dayAfterTomorrow, input, today)

	reqBody := QwenRequest{
		Model: "qwen-plus",
		Messages: []QwenMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	// 通义千问 API 地址
	apiURL := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("通义千问响应: %s\n", string(body))

	var qwenResp QwenResponse
	if err := json.Unmarshal(body, &qwenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if qwenResp.Error != nil {
		return nil, fmt.Errorf("API 错误: %s", qwenResp.Error.Message)
	}

	if len(qwenResp.Choices) == 0 {
		return nil, fmt.Errorf("AI 返回内容为空")
	}

	content := qwenResp.Choices[0].Message.Content
	content = extractJSON(content)

	var result ParseTaskResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v, 内容: %s", err, content)
	}

	if result.StartDate == "" {
		result.StartDate = getTodayStr()
	}
	if result.EndDate == "" {
		result.EndDate = result.StartDate
	}
	if result.StartTime == "" {
		result.StartTime = "09:00"
	}
	if result.EndTime == "" {
		result.EndTime = "18:00"
	}
	if result.Category == "" {
		result.Category = "other"
	}

	return &result, nil
}

func extractJSON(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start != -1 && end != -1 && end > start {
		return text[start : end+1]
	}
	return text
}

func mockParseTask(input string) *ParseTaskResponse {
	today := getTodayStr()

	result := &ParseTaskResponse{
		Title:       input,
		Description: "请完成: " + input,
		StartDate:   today,
		StartTime:   "09:00",
		EndDate:     today,
		EndTime:     "18:00",
		Category:    "other",
	}

	lowerInput := strings.ToLower(input)
	if containsAny(lowerInput, "学习", "复习", "作业", "课程", "章", "节", "编译", "原理") {
		result.Category = "study"
		result.Description = "学习任务: " + input + "，需要认真完成并做好笔记"
	} else if containsAny(lowerInput, "考试", "测验", "测试", "期末", "期中") {
		result.Category = "exam"
		result.Description = "考试相关: " + input + "，需要提前做好准备"
	} else if containsAny(lowerInput, "项目", "开发", "编程", "代码") {
		result.Category = "project"
		result.Description = "项目任务: " + input + "，注意代码质量和进度"
	} else if containsAny(lowerInput, "阅读", "看书", "读书", "论文") {
		result.Category = "reading"
		result.Description = "阅读任务: " + input + "，做好读书笔记"
	}

	if containsAny(lowerInput, "明天") {
		result.StartDate = getTomorrowStr()
		result.EndDate = getTomorrowStr()
	} else if containsAny(lowerInput, "后天") {
		result.StartDate = getDayAfterTomorrowStr()
		result.EndDate = getDayAfterTomorrowStr()
	}

	if containsAny(lowerInput, "上午", "早上", "早晨") {
		result.StartTime = "09:00"
		result.EndTime = "12:00"
	} else if containsAny(lowerInput, "下午") {
		result.StartTime = "14:00"
		result.EndTime = "17:00"
	} else if containsAny(lowerInput, "晚上", "晚间") {
		result.StartTime = "19:00"
		result.EndTime = "22:00"
	}

	return result
}

func mockQuiz(topic string, count int, includeEssay bool) *QuizResponse {
	questions := []MultipleChoiceQuestion{
		{
			Question: fmt.Sprintf("关于%s的基本概念，以下说法正确的是？", topic),
			Options: map[string]string{
				"A": "选项1",
				"B": "选项2",
				"C": "选项3",
				"D": "选项4",
			},
			CorrectAnswer: "A",
			Explanation:   "这是一道基础概念题，正确答案是A，因为...",
			Difficulty:    "easy",
		},
		{
			Question: fmt.Sprintf("在%s的应用场景中，最常见的做法是？", topic),
			Options: map[string]string{
				"A": "方法1",
				"B": "方法2",
				"C": "方法3",
				"D": "方法4",
			},
			CorrectAnswer: "B",
			Explanation:   "这是一道应用题，正确答案是B，因为...",
			Difficulty:    "medium",
		},
		{
			Question: fmt.Sprintf("以下关于%s的高级特性描述，错误的是？", topic),
			Options: map[string]string{
				"A": "特性1",
				"B": "特性2",
				"C": "特性3",
				"D": "特性4",
			},
			CorrectAnswer: "C",
			Explanation:   "这是一道综合分析题，正确答案是C，因为...",
			Difficulty:    "hard",
		},
	}

	// 限制题目数量
	if count < len(questions) {
		questions = questions[:count]
	}

	result := &QuizResponse{
		Questions: questions,
	}

	if includeEssay {
		result.EssayQuestion = &EssayQuestion{
			Question:        fmt.Sprintf("请详细阐述%s的原理和实际应用场景", topic),
			StudySuggestion: fmt.Sprintf("%s是一个重要的概念，其核心原理是...在实际应用中，常见于...需要注意的是...", topic),
		}
	}

	return result
}

func containsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
