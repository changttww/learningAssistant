package rag

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

// HybridSearchService 混合检索服务
// 业界标准：向量检索 + BM25关键词检索 + 重排序
type HybridSearchService struct {
	embeddingService EmbeddingService
	// BM25 参数
	k1 float64 // 词频饱和参数，通常 1.2-2.0
	b  float64 // 文档长度归一化参数，通常 0.75
}

// NewHybridSearchService 创建混合检索服务
func NewHybridSearchService(embeddingService EmbeddingService) *HybridSearchService {
	return &HybridSearchService{
		embeddingService: embeddingService,
		k1:               1.5,
		b:                0.75,
	}
}

// HybridSearchResult 混合检索结果
type HybridSearchResult struct {
	Entry          models.KnowledgeBaseEntry
	VectorScore    float32  // 向量相似度分数 [0,1]
	BM25Score      float32  // BM25 分数
	FinalScore     float32  // 融合后的最终分数
	MatchedTerms   []string // 匹配到的关键词
	MatchHighlight string   // 高亮摘要
}

// Search 执行混合检索
// alpha: 向量分数权重 (0-1)，1-alpha 为 BM25 权重
func (h *HybridSearchService) Search(userID uint64, query string, limit int, alpha float32) ([]HybridSearchResult, error) {
	if limit <= 0 {
		limit = 10
	}
	if alpha < 0 || alpha > 1 {
		alpha = 0.7 // 默认向量权重 70%
	}

	db := database.GetDB()

	// 1. 获取用户所有知识条目
	var entries []models.KnowledgeBaseEntry
	if err := db.Where("user_id = ? AND status = 1", userID).Find(&entries).Error; err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return []HybridSearchResult{}, nil
	}

	// 2. Query 预处理和扩展
	processedQuery := h.preprocessQuery(query)
	queryTerms := h.tokenize(processedQuery)

	// 3. 并行执行两种检索
	vectorScores := h.vectorSearch(entries, query)
	bm25Scores := h.bm25Search(entries, queryTerms)

	// 4. 分数融合 (Reciprocal Rank Fusion 或线性加权)
	results := h.fuseScores(entries, vectorScores, bm25Scores, queryTerms, alpha)

	// 5. 排序并返回 Top-K
	sort.Slice(results, func(i, j int) bool {
		return results[i].FinalScore > results[j].FinalScore
	})

	if len(results) > limit {
		results = results[:limit]
	}

	// 6. 过滤低分结果 - 提高阈值，避免返回不相关内容
	filtered := make([]HybridSearchResult, 0, len(results))
	for _, r := range results {
		// 相似度阈值：至少 35% 才认为相关
		// 低于此阈值的结果很可能是噪音
		if r.FinalScore >= 0.35 {
			filtered = append(filtered, r)
		}
	}

	return filtered, nil
}

// preprocessQuery 查询预处理
func (h *HybridSearchService) preprocessQuery(query string) string {
	// 1. 转小写
	query = strings.ToLower(query)

	// 2. 移除停用词
	stopWords := map[string]bool{
		"的": true, "了": true, "是": true, "在": true, "我": true, "有": true,
		"和": true, "与": true, "这": true, "那": true, "它": true, "他": true,
		"她": true, "们": true, "什么": true, "怎么": true, "如何": true, "为什么": true,
		"哪些": true, "哪个": true, "请": true, "帮我": true, "告诉": true, "一下": true,
		"能": true, "可以": true, "需要": true, "想": true, "要": true, "会": true,
		"the": true, "a": true, "an": true, "is": true, "are": true, "was": true,
		"were": true, "be": true, "been": true, "being": true, "have": true, "has": true,
		"do": true, "does": true, "did": true, "will": true, "would": true, "could": true,
		"should": true, "may": true, "might": true, "must": true, "can": true,
		"what": true, "how": true, "why": true, "when": true, "where": true, "which": true,
	}

	words := h.tokenize(query)
	filtered := make([]string, 0, len(words))
	for _, w := range words {
		if !stopWords[w] && len(w) >= 2 {
			filtered = append(filtered, w)
		}
	}

	return strings.Join(filtered, " ")
}

// tokenize 分词（支持中英文）
func (h *HybridSearchService) tokenize(text string) []string {
	var tokens []string

	// 使用正则分割
	re := regexp.MustCompile(`[\s\p{P}]+`)
	words := re.Split(text, -1)

	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}

		// 检查是否包含中文
		hasChinese := false
		for _, r := range word {
			if unicode.Is(unicode.Han, r) {
				hasChinese = true
				break
			}
		}

		if hasChinese {
			// 中文按字符分割 (简单的 unigram + bigram)
			runes := []rune(word)
			// Unigram
			for _, r := range runes {
				if unicode.Is(unicode.Han, r) {
					tokens = append(tokens, string(r))
				}
			}
			// Bigram
			for i := 0; i < len(runes)-1; i++ {
				if unicode.Is(unicode.Han, runes[i]) && unicode.Is(unicode.Han, runes[i+1]) {
					tokens = append(tokens, string(runes[i:i+2]))
				}
			}
		} else {
			// 英文直接作为 token
			if len(word) >= 2 {
				tokens = append(tokens, strings.ToLower(word))
			}
		}
	}

	return tokens
}

// vectorSearch 向量检索
func (h *HybridSearchService) vectorSearch(entries []models.KnowledgeBaseEntry, query string) map[uint64]float32 {
	scores := make(map[uint64]float32)

	// 生成查询向量
	queryVector, err := h.embeddingService.GenerateEmbedding(query)
	if err != nil || len(queryVector) == 0 {
		return scores
	}

	db := database.GetDB()

	// 获取所有条目的向量缓存
	entryIDs := make([]uint64, len(entries))
	for i, e := range entries {
		entryIDs[i] = e.ID
	}

	var caches []models.KnowledgeVectorCache
	db.Where("entry_id IN ?", entryIDs).Find(&caches)

	// 创建 entryID -> vector 映射
	vectorMap := make(map[uint64]models.Vector)
	for _, cache := range caches {
		vectorMap[cache.EntryID] = cache.Vector
	}

	// 计算相似度
	for _, entry := range entries {
		if vec, ok := vectorMap[entry.ID]; ok && len(vec) > 0 {
			similarity := h.embeddingService.CosineSimilarity(queryVector, vec)
			scores[entry.ID] = similarity
		}
	}

	return scores
}

// bm25Search BM25 检索
func (h *HybridSearchService) bm25Search(entries []models.KnowledgeBaseEntry, queryTerms []string) map[uint64]float32 {
	scores := make(map[uint64]float32)

	if len(queryTerms) == 0 {
		return scores
	}

	// 构建文档集合
	docs := make([][]string, len(entries))
	totalLength := 0
	for i, entry := range entries {
		// 合并标题、内容、关键词作为文档
		docText := entry.Title + " " + entry.Content + " " + entry.Category
		docs[i] = h.tokenize(docText)
		totalLength += len(docs[i])
	}

	avgDocLen := float64(totalLength) / float64(len(docs))

	// 计算 IDF
	idf := h.calculateIDF(docs, queryTerms)

	// 计算每个文档的 BM25 分数
	for i, entry := range entries {
		score := h.calculateBM25(docs[i], queryTerms, idf, avgDocLen)
		scores[entry.ID] = float32(score)
	}

	return scores
}

// calculateIDF 计算逆文档频率
func (h *HybridSearchService) calculateIDF(docs [][]string, terms []string) map[string]float64 {
	idf := make(map[string]float64)
	n := float64(len(docs))

	for _, term := range terms {
		df := 0
		for _, doc := range docs {
			for _, t := range doc {
				if t == term {
					df++
					break
				}
			}
		}
		// IDF 公式: log((N - df + 0.5) / (df + 0.5) + 1)
		idf[term] = math.Log((n-float64(df)+0.5)/(float64(df)+0.5) + 1)
	}

	return idf
}

// calculateBM25 计算单个文档的 BM25 分数
func (h *HybridSearchService) calculateBM25(doc []string, queryTerms []string, idf map[string]float64, avgDocLen float64) float64 {
	score := 0.0
	docLen := float64(len(doc))

	// 统计词频
	tf := make(map[string]int)
	for _, t := range doc {
		tf[t]++
	}

	for _, term := range queryTerms {
		termFreq := float64(tf[term])
		if termFreq == 0 {
			continue
		}

		// BM25 公式
		numerator := termFreq * (h.k1 + 1)
		denominator := termFreq + h.k1*(1-h.b+h.b*(docLen/avgDocLen))
		score += idf[term] * (numerator / denominator)
	}

	return score
}

// fuseScores 融合向量和 BM25 分数
func (h *HybridSearchService) fuseScores(
	entries []models.KnowledgeBaseEntry,
	vectorScores map[uint64]float32,
	bm25Scores map[uint64]float32,
	queryTerms []string,
	alpha float32,
) []HybridSearchResult {

	results := make([]HybridSearchResult, 0, len(entries))

	// 归一化 BM25 分数
	maxBM25 := float32(0)
	for _, score := range bm25Scores {
		if score > maxBM25 {
			maxBM25 = score
		}
	}

	for _, entry := range entries {
		vectorScore := vectorScores[entry.ID]
		bm25Score := bm25Scores[entry.ID]

		// 归一化 BM25
		normalizedBM25 := float32(0)
		if maxBM25 > 0 {
			normalizedBM25 = bm25Score / maxBM25
		}

		// 线性加权融合
		finalScore := alpha*vectorScore + (1-alpha)*normalizedBM25

		// 找到匹配的关键词
		matchedTerms := h.findMatchedTerms(entry, queryTerms)

		// 生成高亮摘要
		highlight := h.generateHighlight(entry.Content, queryTerms)

		if finalScore > 0 || len(matchedTerms) > 0 {
			results = append(results, HybridSearchResult{
				Entry:          entry,
				VectorScore:    vectorScore,
				BM25Score:      normalizedBM25,
				FinalScore:     finalScore,
				MatchedTerms:   matchedTerms,
				MatchHighlight: highlight,
			})
		}
	}

	return results
}

// findMatchedTerms 找到匹配的关键词
func (h *HybridSearchService) findMatchedTerms(entry models.KnowledgeBaseEntry, queryTerms []string) []string {
	content := strings.ToLower(entry.Title + " " + entry.Content + " " + entry.Category)
	matched := make([]string, 0)

	for _, term := range queryTerms {
		if strings.Contains(content, term) {
			matched = append(matched, term)
		}
	}

	return matched
}

// generateHighlight 生成高亮摘要
func (h *HybridSearchService) generateHighlight(content string, queryTerms []string) string {
	if len(content) == 0 || len(queryTerms) == 0 {
		return ""
	}

	// 找到第一个匹配词的位置
	contentLower := strings.ToLower(content)
	firstMatchPos := -1
	for _, term := range queryTerms {
		pos := strings.Index(contentLower, term)
		if pos != -1 && (firstMatchPos == -1 || pos < firstMatchPos) {
			firstMatchPos = pos
		}
	}

	if firstMatchPos == -1 {
		// 没有匹配，返回开头
		if len(content) > 150 {
			return content[:150] + "..."
		}
		return content
	}

	// 提取匹配位置前后的上下文
	start := firstMatchPos - 50
	if start < 0 {
		start = 0
	}
	end := firstMatchPos + 100
	if end > len(content) {
		end = len(content)
	}

	excerpt := content[start:end]
	if start > 0 {
		excerpt = "..." + excerpt
	}
	if end < len(content) {
		excerpt = excerpt + "..."
	}

	return excerpt
}

// QueryExpansion 查询扩展（同义词扩展）
func (h *HybridSearchService) QueryExpansion(query string) []string {
	// 学科领域同义词映射
	synonyms := map[string][]string{
		"数学":     {"数学", "算术", "代数", "几何", "math", "mathematics"},
		"物理":     {"物理", "力学", "电学", "光学", "physics"},
		"化学":     {"化学", "元素", "分子", "chemistry"},
		"编程":     {"编程", "程序", "代码", "开发", "programming", "coding"},
		"英语":     {"英语", "英文", "english"},
		"函数":     {"函数", "方程", "公式", "function"},
		"算法":     {"算法", "algorithm", "数据结构"},
		"学习":     {"学习", "复习", "掌握", "理解"},
		"python": {"python", "py", "编程"},
		"java":   {"java", "编程", "开发"},
	}

	expandedTerms := []string{query}

	queryLower := strings.ToLower(query)
	for key, syns := range synonyms {
		if strings.Contains(queryLower, key) {
			expandedTerms = append(expandedTerms, syns...)
		}
	}

	return expandedTerms
}
