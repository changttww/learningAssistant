package rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

// RelationMiningService AI关系挖掘服务
type RelationMiningService struct {
	embeddingService EmbeddingService
	apiKey           string
}

// RelationCandidate 关系候选
type RelationCandidate struct {
	EntryID    uint64  `json:"entry_id"`
	Title      string  `json:"title"`
	Summary    string  `json:"summary"`
	Category   string  `json:"category"`
	Similarity float32 `json:"similarity"`
}

// AIRelationResult AI推理返回的关系结果
type AIRelationResult struct {
	SourceID     uint64 `json:"source_id"`
	TargetID     uint64 `json:"target_id"`
	RelationType int8   `json:"relation_type"` // 1=前置, 2=相关, 3=扩展
	Reason       string `json:"reason"`
}

// MiningResult 挖掘结果
type MiningResult struct {
	NewEntryID     uint64             `json:"new_entry_id"`
	RelationsFound int                `json:"relations_found"`
	Relations      []AIRelationResult `json:"relations"`
}

// NewRelationMiningService 创建关系挖掘服务
func NewRelationMiningService(embeddingService EmbeddingService) *RelationMiningService {
	apiKey := os.Getenv("QWEN_API_KEY")
	return &RelationMiningService{
		embeddingService: embeddingService,
		apiKey:           apiKey,
	}
}

// MineRelationsForEntry 为新知识点挖掘与现有知识点的关系
// 这是AI Auto-Linking的核心方法
func (s *RelationMiningService) MineRelationsForEntry(userID uint64, entryID uint64) (*MiningResult, error) {
	db := database.GetDB()

	// 1. 获取新知识点信息
	var newEntry models.KnowledgeBaseEntry
	if err := db.First(&newEntry, entryID).Error; err != nil {
		return nil, fmt.Errorf("获取知识点失败: %w", err)
	}

	// 2. 获取新知识点的向量
	var newCache models.KnowledgeVectorCache
	if err := db.Where("entry_id = ?", entryID).First(&newCache).Error; err != nil {
		return nil, fmt.Errorf("获取向量缓存失败: %w", err)
	}

	// 3. 候选召回：获取用户其他知识点，计算相似度，取Top 8
	candidates, err := s.getCandidates(userID, entryID, newCache.Vector, 8)
	if err != nil {
		return nil, fmt.Errorf("召回候选失败: %w", err)
	}

	if len(candidates) == 0 {
		return &MiningResult{
			NewEntryID:     entryID,
			RelationsFound: 0,
			Relations:      []AIRelationResult{},
		}, nil
	}

	// 4. AI推理：判断关系类型
	relations, err := s.inferRelations(newEntry, candidates)
	if err != nil {
		fmt.Printf("[RelationMining] AI推理失败，使用降级策略: %v\n", err)
		// 降级策略：基于相似度自动建立"相关"关系
		relations = s.fallbackRelations(newEntry, candidates)
	}

	// 5. 落库：将关系写入数据库（带防环检测）
	savedCount := 0
	skippedCount := 0
	for _, rel := range relations {
		// 检查是否已存在该关系（正向）
		var existingCount int64
		db.Model(&models.KnowledgeRelation{}).
			Where("user_id = ? AND source_entry_id = ? AND target_entry_id = ?",
				userID, rel.SourceID, rel.TargetID).
			Count(&existingCount)

		if existingCount > 0 {
			continue // 关系已存在，跳过
		}

		// 【防环检测】检查是否存在反向的定向关系（前置/扩展）
		// 只对定向关系(1=前置, 3=扩展)做防环检测，"相关"关系(2)是双向的，无需检测
		if rel.RelationType == 1 || rel.RelationType == 3 {
			var reverseCount int64
			db.Model(&models.KnowledgeRelation{}).
				Where("user_id = ? AND source_entry_id = ? AND target_entry_id = ? AND relation_type IN (1, 3)",
					userID, rel.TargetID, rel.SourceID).
				Count(&reverseCount)

			if reverseCount > 0 {
				// 存在反向定向关系，禁止插入，避免形成 A->B 且 B->A 的互指环
				fmt.Printf("[RelationMining] 防环检测：拒绝插入 %d->%d (类型:%d)，因为已存在反向关系 %d->%d\n",
					rel.SourceID, rel.TargetID, rel.RelationType, rel.TargetID, rel.SourceID)
				skippedCount++
				continue
			}
		}

		// 通过防环检测，创建关系
		relation := &models.KnowledgeRelation{
			UserID:        userID,
			SourceEntryID: rel.SourceID,
			TargetEntryID: rel.TargetID,
			RelationType:  rel.RelationType,
			Strength:      0.8, // AI推断的关系强度较高
		}
		if err := db.Create(relation).Error; err != nil {
			fmt.Printf("[RelationMining] 保存关系失败: %v\n", err)
			continue
		}
		savedCount++
	}

	if skippedCount > 0 {
		fmt.Printf("[RelationMining] 防环检测拦截了 %d 个可能形成互指的关系\n", skippedCount)
	}

	fmt.Printf("[RelationMining] 为知识点 %d(%s) 挖掘到 %d 个关系，保存 %d 个\n",
		entryID, newEntry.Title, len(relations), savedCount)

	return &MiningResult{
		NewEntryID:     entryID,
		RelationsFound: savedCount,
		Relations:      relations,
	}, nil
}

// getCandidates 向量召回候选知识点
func (s *RelationMiningService) getCandidates(userID uint64, excludeID uint64, queryVector models.Vector, limit int) ([]RelationCandidate, error) {
	db := database.GetDB()

	// 获取用户其他知识点的向量缓存
	var caches []models.KnowledgeVectorCache
	subQuery := db.Table("knowledge_base_entries").
		Select("id").
		Where("user_id = ? AND status = 1 AND id != ?", userID, excludeID)

	if err := db.Where("entry_id IN (?)", subQuery).Find(&caches).Error; err != nil {
		return nil, err
	}

	// 计算相似度并排序
	type scoredCandidate struct {
		entryID    uint64
		similarity float32
	}
	var scored []scoredCandidate

	for _, cache := range caches {
		similarity := s.embeddingService.CosineSimilarity(queryVector, cache.Vector)
		// 相似度阈值：0.4以上才考虑有关系
		if similarity >= 0.4 {
			scored = append(scored, scoredCandidate{
				entryID:    cache.EntryID,
				similarity: similarity,
			})
		}
	}

	// 按相似度降序排序
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].similarity > scored[i].similarity {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// 取Top N
	if len(scored) > limit {
		scored = scored[:limit]
	}

	// 获取知识点详情
	var candidates []RelationCandidate
	for _, s := range scored {
		var entry models.KnowledgeBaseEntry
		if err := db.First(&entry, s.entryID).Error; err != nil {
			continue
		}
		candidates = append(candidates, RelationCandidate{
			EntryID:    entry.ID,
			Title:      entry.Title,
			Summary:    entry.Summary,
			Category:   entry.Category,
			Similarity: s.similarity,
		})
	}

	return candidates, nil
}

// inferRelations 调用AI推理关系类型
func (s *RelationMiningService) inferRelations(newEntry models.KnowledgeBaseEntry, candidates []RelationCandidate) ([]AIRelationResult, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("未配置API Key")
	}

	// 构造候选知识点描述
	var candidateDescs []string
	for i, c := range candidates {
		desc := fmt.Sprintf("%d. [ID:%d] %s - %s (分类:%s, 相似度:%.2f)",
			i+1, c.EntryID, c.Title, truncateSummary(c.Summary, 50), c.Category, c.Similarity)
		candidateDescs = append(candidateDescs, desc)
	}

	// 构造Prompt
	prompt := fmt.Sprintf(`你是一个知识图谱构建专家。请分析以下新知识点与候选知识点之间的逻辑关系。

【新知识点】
ID: %d
标题: %s
摘要: %s
分类: %s

【候选知识点列表】
%s

【关系类型说明】
1 = 前置关系 (prerequisite)：学习新知识点之前必须先掌握候选知识点。例如：学"深度学习"前需要掌握"线性代数"
2 = 相关关系 (related)：两个知识点属于同一层级的相关概念。例如："Python"和"Java"都是编程语言
3 = 扩展关系 (extends)：新知识点是候选知识点的具体应用、实战项目或进阶内容。例如："贪吃蛇项目"是"Python基础"的扩展

【任务】
请判断新知识点与每个候选知识点是否存在上述关系。只输出确实存在关系的配对。

【输出格式】
请直接输出JSON数组，不要有任何其他文字：
[
  {"source_id": 候选知识点ID, "target_id": %d, "relation_type": 关系类型数字, "reason": "简短理由"},
  ...
]

如果是前置关系，source_id是前置知识点，target_id是新知识点。
如果是扩展关系，source_id是基础知识点，target_id是扩展知识点（新知识点）。
如果没有明确关系，返回空数组 []`,
		newEntry.ID, newEntry.Title, truncateSummary(newEntry.Summary, 100), newEntry.Category,
		strings.Join(candidateDescs, "\n"),
		newEntry.ID,
	)

	// 调用AI
	qwenURL := os.Getenv("QWEN_CHAT_URL")
	if qwenURL == "" {
		qwenURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	}

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
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResult struct {
		Output struct {
			Text string `json:"text"`
		} `json:"output"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		return nil, err
	}

	if apiResult.Output.Text == "" {
		return nil, fmt.Errorf("AI返回为空")
	}

	// 解析JSON结果
	jsonStr := strings.TrimSpace(apiResult.Output.Text)
	jsonStr = strings.TrimPrefix(jsonStr, "```json")
	jsonStr = strings.TrimPrefix(jsonStr, "```")
	jsonStr = strings.TrimSuffix(jsonStr, "```")
	jsonStr = strings.TrimSpace(jsonStr)

	var relations []AIRelationResult
	if err := json.Unmarshal([]byte(jsonStr), &relations); err != nil {
		fmt.Printf("[RelationMining] 解析AI结果失败: %v, 原始内容: %s\n", err, jsonStr)
		return nil, err
	}

	// 验证关系类型
	var validRelations []AIRelationResult
	for _, rel := range relations {
		if rel.RelationType >= 1 && rel.RelationType <= 3 {
			validRelations = append(validRelations, rel)
		}
	}

	return validRelations, nil
}

// fallbackRelations 降级策略：基于相似度建立相关关系
func (s *RelationMiningService) fallbackRelations(newEntry models.KnowledgeBaseEntry, candidates []RelationCandidate) []AIRelationResult {
	var relations []AIRelationResult

	for _, c := range candidates {
		// 相似度>0.6的建立"相关"关系
		if c.Similarity >= 0.6 {
			relations = append(relations, AIRelationResult{
				SourceID:     c.EntryID,
				TargetID:     newEntry.ID,
				RelationType: 2, // 相关
				Reason:       fmt.Sprintf("基于语义相似度(%.2f)自动关联", c.Similarity),
			})
		}
	}

	// 同分类的建立弱关联
	for _, c := range candidates {
		if c.Category == newEntry.Category && c.Similarity >= 0.5 && c.Similarity < 0.6 {
			relations = append(relations, AIRelationResult{
				SourceID:     c.EntryID,
				TargetID:     newEntry.ID,
				RelationType: 2, // 相关
				Reason:       "同分类知识点",
			})
		}
	}

	return relations
}

// MineAllRelations 批量挖掘用户所有知识点的关系
func (s *RelationMiningService) MineAllRelations(userID uint64) (int, error) {
	db := database.GetDB()

	// 获取用户所有知识点
	var entries []models.KnowledgeBaseEntry
	if err := db.Where("user_id = ? AND status = 1", userID).Find(&entries).Error; err != nil {
		return 0, fmt.Errorf("获取知识条目失败: %w", err)
	}

	totalRelations := 0
	for i, entry := range entries {
		result, err := s.MineRelationsForEntry(userID, entry.ID)
		if err != nil {
			fmt.Printf("[RelationMining] 处理知识点 %d 失败: %v\n", entry.ID, err)
			continue
		}
		totalRelations += result.RelationsFound

		// 进度日志
		if (i+1)%10 == 0 {
			fmt.Printf("[RelationMining] 已处理 %d/%d 个知识点\n", i+1, len(entries))
		}

		// 避免API限流
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("[RelationMining] 批量挖掘完成，共发现 %d 个关系\n", totalRelations)
	return totalRelations, nil
}

// truncateSummary 截断摘要
func truncateSummary(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
