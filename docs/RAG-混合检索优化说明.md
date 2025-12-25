# RAG 知识库检索优化说明

## 优化前后对比

### 优化前的问题
1. **简单 LIKE 查询**：无法理解语义相似性，"学习方法" 搜不到 "怎么高效学习"
2. **本地 Embedding 效果差**：基于字符哈希，无真正的语义理解
3. **单一检索策略**：要么向量、要么关键词，无法互补

### 优化后的方案（业界标准）

```
用户Query 
    ↓
Query预处理（停用词过滤、分词）
    ↓
┌─────────────────────────────────────┐
│         并行执行两种检索              │
├─────────────────┬───────────────────┤
│   向量检索       │    BM25检索       │
│  (语义相似度)     │   (关键词匹配)     │
└─────────────────┴───────────────────┘
    ↓
分数融合（线性加权：α*向量 + (1-α)*BM25）
    ↓
Top-K 结果 + 引用信息
    ↓
LLM 生成回答
```

## 新增功能

### 1. 混合检索 (Hybrid Search)
- **向量检索**：捕捉语义相似性（"机器学习" ≈ "人工智能"）
- **BM25 检索**：精确关键词匹配，处理专有名词
- **分数融合**：默认 60% 向量 + 40% BM25

### 2. BM25 算法
- 业界标准的关键词检索算法（Elasticsearch 默认使用）
- 考虑词频 (TF)、逆文档频率 (IDF)、文档长度归一化

### 3. 中文分词支持
- Unigram + Bigram 分词
- 停用词过滤
- 同义词扩展（可扩展）

### 4. 搜索结果增强
- 返回匹配分数
- 返回匹配的关键词
- 返回高亮摘要

## 如何进一步提升效果

### 方案1：使用真实 Embedding API（推荐）

修改 `backend/routes/knowledge_base.go` 中的 `initRAGServices()`:

```go
func initRAGServices() {
    // 使用 Qwen Embedding API（阿里云）
    embeddingService := rag.NewQwenEmbeddingService("your-api-key")
    
    // 或使用环境变量
    // export QWEN_API_KEY=your-api-key
    // embeddingService := rag.NewQwenEmbeddingService("")
    
    ragService = rag.NewRAGService(embeddingService)
    aiAnalysisService = rag.NewAIAnalysisService("your-llm-api-key")
    hybridSearchService = rag.NewHybridSearchService(embeddingService)
}
```

**支持的 Embedding 服务**：
- 阿里云 Qwen Embedding (text-embedding-v3)
- OpenAI Ada Embedding
- 智谱 GLM Embedding
- 本地部署的 sentence-transformers

### 方案2：使用向量数据库（大规模场景）

对于 >10000 条知识，建议使用专业向量数据库：
- **Milvus**：开源，性能好
- **Qdrant**：Rust 编写，轻量级
- **Pinecone**：云服务，免运维
- **PostgreSQL + pgvector**：适合已有 PG 的项目

### 方案3：添加 Reranking（高级）

在检索后添加重排序模型：
```
召回 Top-50 → Reranker 模型 → 精排 Top-5
```

推荐模型：
- Cohere Rerank
- bge-reranker-large
- jina-reranker-v1

## 配置建议

| 场景 | 向量权重(α) | 说明 |
|------|------------|------|
| 语义搜索为主 | 0.7-0.8 | "类似xxx的内容" |
| 精确搜索为主 | 0.3-0.4 | 搜索专有名词、代码 |
| 均衡 | 0.5-0.6 | 通用场景（默认） |

## 测试检索效果

```bash
# 测试搜索 API
curl -X GET "http://localhost:8080/api/v1/knowledge-base/search?q=学习方法" \
  -H "Authorization: Bearer your-token"

# 测试 RAG 问答
curl -X POST "http://localhost:8080/api/v1/knowledge-base/chat" \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{"query": "总结一下我学过的数学知识", "limit": 5}'
```
