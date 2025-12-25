# RAG 知识助手技术文档

> 基于检索增强生成（Retrieval-Augmented Generation）的个人知识库问答系统

## 📋 目录

- [系统概述](#系统概述)
- [技术架构](#技术架构)
- [核心组件](#核心组件)
- [数据模型](#数据模型)
- [API 接口](#api-接口)
- [配置说明](#配置说明)
- [使用指南](#使用指南)
- [性能优化](#性能优化)
- [常见问题](#常见问题)

---

## 系统概述

### 什么是 RAG？

RAG（Retrieval-Augmented Generation）是一种结合信息检索和文本生成的技术架构。它通过以下步骤工作：

1. **检索（Retrieval）**：根据用户查询，从知识库中检索相关文档
2. **增强（Augmented）**：将检索到的文档作为上下文
3. **生成（Generation）**：结合上下文，由大语言模型生成准确的回答

### 系统特点

| 特性 | 说明 |
|------|------|
| 🎯 **语义检索** | 基于 Qwen Embedding 的 1024 维向量，理解语义而非简单关键词匹配 |
| 🔀 **混合检索** | Vector + BM25 双路召回，兼顾语义理解和关键词精确匹配 |
| 📚 **知识图谱** | 可视化展示知识点之间的关联关系 |
| 🎓 **学习追踪** | 掌握等级、复习提醒、学习统计 |
| 💬 **引用溯源** | 每个回答都标注知识来源，可追溯验证 |

---

## 技术架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         前端层 (Vue 3)                           │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────────┐   │
│  │  知识图谱      │  │  知识问答      │  │  知识库管理        │   │
│  │  ECharts      │  │  Markdown     │  │  CRUD 操作        │   │
│  │  Force Graph  │  │  marked.js    │  │                   │   │
│  └───────────────┘  └───────────────┘  └───────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │ HTTP/WebSocket
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         后端层 (Go/Gin)                          │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                      RAG 服务层                          │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │   │
│  │  │  Embedding  │  │  Hybrid     │  │  LLM        │      │   │
│  │  │  Service    │  │  Search     │  │  Service    │      │   │
│  │  │             │  │             │  │             │      │   │
│  │  │  Qwen       │  │  Vector     │  │  Qwen       │      │   │
│  │  │  text-      │  │  +          │  │  qwen-plus  │      │   │
│  │  │  embedding  │  │  BM25       │  │             │      │   │
│  │  │  -v3        │  │             │  │             │      │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘      │   │
│  └─────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    知识管理层                            │   │
│  │  • 文档 CRUD      • 分类管理      • 关系管理              │   │
│  │  • 向量缓存       • 统计分析      • 图谱生成              │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │ GORM
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        数据层 (MySQL)                            │
│  ┌───────────────────┐  ┌───────────────────────────────────┐  │
│  │ knowledge_base    │  │ knowledge_vector_caches           │  │
│  │ _entries          │  │                                   │  │
│  │                   │  │  • entry_id (FK)                  │  │
│  │  • id             │  │  • content_hash                   │  │
│  │  • user_id        │  │  • vector (1024维)                │  │
│  │  • title          │  │  • vector_model                   │  │
│  │  • content        │  │                                   │  │
│  │  • category       │  └───────────────────────────────────┘  │
│  │  • level          │  ┌───────────────────────────────────┐  │
│  │  • ...            │  │ knowledge_relations               │  │
│  └───────────────────┘  └───────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### 技术栈

| 层级 | 技术 | 版本 | 用途 |
|------|------|------|------|
| 前端 | Vue 3 | 3.x | 响应式 UI 框架 |
| 前端 | ECharts | 5.x | 知识图谱可视化 |
| 前端 | marked | 4.x | Markdown 渲染 |
| 前端 | Tailwind CSS | 3.x | 样式框架 |
| 后端 | Go | 1.21+ | 主开发语言 |
| 后端 | Gin | 1.9+ | Web 框架 |
| 后端 | GORM | 2.x | ORM 框架 |
| 数据库 | MySQL | 8.0+ | 关系型数据库 |
| AI | Qwen Embedding | text-embedding-v3 | 文本向量化 |
| AI | Qwen LLM | qwen-plus | 回答生成 |

---

## 核心组件

### 1. Embedding 服务

负责将文本转换为高维向量，用于语义相似度计算。

```go
// EmbeddingService 向量化服务接口
type EmbeddingService interface {
    // 生成文本向量
    GenerateEmbedding(text string) (models.Vector, error)
    // 计算向量相似度
    CosineSimilarity(vec1, vec2 models.Vector) float32
    // 批量生成向量
    GenerateEmbeddings(texts []string) ([]models.Vector, error)
}
```

**Qwen Embedding 配置：**
- 模型：`text-embedding-v3`
- 维度：1024
- API：DashScope 兼容模式
- 超时：30 秒

### 2. 混合检索服务

结合向量检索和 BM25 关键词检索，提升召回质量。

```
用户查询: "什么是牛顿第二定律"
         │
         ├─────────────────┬─────────────────┐
         ▼                 ▼                 │
   ┌───────────┐    ┌───────────┐           │
   │ Vector 检索│    │ BM25 检索 │           │
   │ (权重 0.7) │    │ (权重 0.3) │           │
   └───────────┘    └───────────┘           │
         │                 │                 │
         │  语义匹配        │  关键词匹配      │
         │  "力与加速度"    │  "牛顿" "定律"   │
         │                 │                 │
         └────────┬────────┘                 │
                  ▼                          │
         ┌───────────────┐                   │
         │  RRF 分数融合  │                   │
         │  倒数排名融合   │                   │
         └───────────────┘                   │
                  │                          │
                  ▼                          │
         ┌───────────────┐                   │
         │  Top-K 结果    │◄──────────────────┘
         │  排序输出       │
         └───────────────┘
```

**BM25 参数：**
- k1 = 1.2（词频饱和参数）
- b = 0.75（文档长度归一化参数）

**融合公式：**
```
final_score = α × vector_score + (1-α) × bm25_score
其中 α = 0.7
```

### 3. RAG 问答流程

```
┌──────────────────────────────────────────────────────────────┐
│                        RAG 问答流程                           │
└──────────────────────────────────────────────────────────────┘

Step 1: 查询理解
┌─────────────────┐
│  用户提问        │  "帮我解释一下相对论的基本概念"
└────────┬────────┘
         │
         ▼
Step 2: 向量化
┌─────────────────┐
│  Query Embedding │  [0.023, -0.041, 0.087, ...]  (1024维)
└────────┬────────┘
         │
         ▼
Step 3: 混合检索
┌─────────────────┐
│  Hybrid Search   │  检索 Top-5 相关文档
│  Vector + BM25   │
└────────┬────────┘
         │
         ▼
Step 4: 上下文构建
┌─────────────────┐
│  Context Build   │  拼接检索到的文档内容
│                  │  添加引用标记
└────────┬────────┘
         │
         ▼
Step 5: LLM 生成
┌─────────────────┐
│  Qwen qwen-plus  │  生成带引用的回答
└────────┬────────┘
         │
         ▼
Step 6: 响应返回
┌─────────────────┐
│  Response        │  Markdown 格式回答
│  + Citations     │  + 引用来源列表
└─────────────────┘
```

### 4. 知识图谱

基于 ECharts Force Graph 实现的交互式知识网络可视化。

**节点类型：**
| 属性 | 说明 |
|------|------|
| ID | 知识条目唯一标识 |
| Name | 标题（截断至 20 字符） |
| Category | 学科分类 |
| Level | 掌握等级 (0-4) |
| Value | 节点大小 = ViewCount + 10 |
| Color | 分类对应颜色 |

**边类型：**
| 类型 | 代码 | 说明 |
|------|------|------|
| 前置 | 1 | A 是 B 的前置知识 |
| 相关 | 2 | A 与 B 相关 |
| 扩展 | 3 | A 扩展自 B |
| 冲突 | 4 | A 与 B 存在冲突 |
| 同分类 | 5 | A 与 B 属于同一分类（自动生成） |

---

## 数据模型

### 知识条目 (knowledge_base_entries)

```sql
CREATE TABLE knowledge_base_entries (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    created_at      DATETIME(3),
    updated_at      DATETIME(3),
    deleted_at      DATETIME(3),
    
    user_id         BIGINT UNSIGNED NOT NULL,      -- 所属用户
    source_type     TINYINT NOT NULL,              -- 来源类型: 1=任务, 2=笔记, 3=手动
    source_id       BIGINT UNSIGNED,               -- 来源 ID
    
    title           VARCHAR(500) NOT NULL,         -- 标题
    content         TEXT,                          -- 内容（纯文本）
    summary         TEXT,                          -- 摘要（前200字）
    keywords        JSON,                          -- 关键词数组
    
    category        VARCHAR(50),                   -- 一级分类
    sub_category    VARCHAR(50),                   -- 二级分类
    subject         VARCHAR(50),                   -- 学科
    
    level           TINYINT DEFAULT 0,             -- 掌握等级 0-4
    status          TINYINT DEFAULT 1,             -- 状态: 1=正常
    view_count      INT DEFAULT 0,                 -- 浏览次数
    last_review_at  DATETIME(3),                   -- 最后复习时间
    
    display_color   VARCHAR(20),                   -- 显示颜色
    display_icon    VARCHAR(50),                   -- 显示图标
    
    INDEX idx_user_status (user_id, status),
    INDEX idx_category (category)
);
```

### 向量缓存 (knowledge_vector_caches)

```sql
CREATE TABLE knowledge_vector_caches (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    created_at      DATETIME(3),
    updated_at      DATETIME(3),
    deleted_at      DATETIME(3),
    
    entry_id        BIGINT UNSIGNED NOT NULL,      -- 关联知识条目
    content_hash    VARCHAR(32) NOT NULL,          -- 内容 MD5 哈希
    vector          JSON NOT NULL,                 -- 向量数据 (1024维)
    vector_dim      INT NOT NULL,                  -- 向量维度
    vector_model    VARCHAR(50),                   -- 使用的模型名称
    
    UNIQUE INDEX idx_entry_id (entry_id)
);
```

### 知识关系 (knowledge_relations)

```sql
CREATE TABLE knowledge_relations (
    id               BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id          BIGINT UNSIGNED NOT NULL,
    source_entry_id  BIGINT UNSIGNED NOT NULL,     -- 源知识点
    target_entry_id  BIGINT UNSIGNED NOT NULL,     -- 目标知识点
    relation_type    TINYINT NOT NULL,             -- 关系类型 1-5
    strength         FLOAT DEFAULT 0.5,            -- 关系强度 0-1
    
    INDEX idx_source (source_entry_id),
    INDEX idx_target (target_entry_id)
);
```

---

## API 接口

### 知识库管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/knowledge/entries` | 获取知识条目列表 |
| GET | `/api/knowledge/entries/:id` | 获取单个知识条目 |
| POST | `/api/knowledge/entries` | 创建知识条目 |
| PUT | `/api/knowledge/entries/:id` | 更新知识条目 |
| DELETE | `/api/knowledge/entries/:id` | 删除知识条目 |
| POST | `/api/knowledge/sync` | 同步知识库 |

### RAG 问答

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/knowledge/chat` | RAG 问答 |
| GET | `/api/knowledge/search` | 搜索知识库 |

**问答请求示例：**
```json
{
    "query": "什么是牛顿第二定律？",
    "limit": 5
}
```

**问答响应示例：**
```json
{
    "answer": "牛顿第二定律表明...",
    "citations": [
        {
            "id": 123,
            "title": "物理-力学基础",
            "category": "物理",
            "summary": "介绍力学的基本概念...",
            "similarity": 0.89
        }
    ],
    "query": "什么是牛顿第二定律？"
}
```

### 知识图谱

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/knowledge/graph` | 获取知识图谱数据 |
| GET | `/api/knowledge/stats` | 获取知识统计 |

---

## 配置说明

### 环境变量

在 `backend/.env` 文件中配置：

```properties
# 服务器配置
SERVER_PORT=8080
GIN_MODE=debug

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_DATABASE=learning_assistant
DB_CHARSET=utf8mb4

# AI 配置 - 通义千问 API
# 从 https://dashscope.console.aliyun.com/ 获取
QWEN_API_KEY=sk-xxxxxxxxxxxxxxxx
```

### 获取 API Key

1. 访问 [阿里云 DashScope 控制台](https://dashscope.console.aliyun.com/)
2. 注册/登录账号
3. 创建 API Key
4. 将 Key 填入 `.env` 文件的 `QWEN_API_KEY`

### 服务启动

```bash
# 后端
cd backend
go build -o learningAssistant-backend.exe .
./learningAssistant-backend.exe

# 前端
cd front
npm install
npm run dev
```

启动成功后，日志会显示：
```
[RAG] 使用 Qwen Embedding API (text-embedding-v3)
```

---

## 使用指南

### 1. 知识入库

知识自动从以下来源同步：
- ✅ 完成的学习任务
- ✅ 创建的学习笔记
- ✅ 手动添加的知识点

点击"🔄 同步知识库"按钮手动触发同步。

### 2. 知识问答

在知识问答页面：
1. 输入问题
2. 系统检索相关知识
3. AI 生成回答
4. 查看引用来源

### 3. 知识图谱

在知识图谱页面：
- 🔵 蓝色节点：数学
- 🟣 紫色节点：物理
- 🟢 绿色节点：化学/生物
- 🟡 黄色节点：语文
- 🔴 红色节点：英语
- ...

节点大小反映浏览次数，连线表示知识关联。

### 4. 掌握等级

| 等级 | 名称 | 说明 |
|------|------|------|
| 0 | 待学习 | 尚未开始学习 |
| 1 | 了解 | 初步了解概念 |
| 2 | 熟悉 | 能够理解和应用 |
| 3 | 掌握 | 熟练掌握 |
| 4 | 精通 | 融会贯通 |

---

## 性能优化

### 当前优化措施

1. **向量缓存**
   - 知识条目的向量只在内容变化时重新生成
   - 使用 MD5 哈希检测内容变化

2. **预编译正则**
   - HTML 标签清理使用预编译正则
   - 避免重复编译开销

3. **批量查询**
   - 一次性获取用户所有向量缓存
   - 减少数据库往返

### 未来优化建议

| 优化项 | 说明 | 优先级 |
|--------|------|--------|
| 向量数据库 | 使用 Milvus/Qdrant 替代 MySQL 存储向量 | 高 |
| Reranking | 添加重排序模型提升精度 | 中 |
| 增量索引 | 实时更新向量，无需全量同步 | 中 |
| 缓存层 | Redis 缓存热点查询 | 低 |
| 分块策略 | 长文档自动分块，提升召回 | 低 |

---

## 常见问题

### Q1: 向量生成失败怎么办？

**症状**：同步知识库时，部分条目没有生成向量

**排查步骤**：
1. 检查 `QWEN_API_KEY` 是否正确配置
2. 检查网络是否能访问 `dashscope.aliyuncs.com`
3. 查看后端日志中的错误信息

### Q2: 搜索结果不相关怎么办？

**可能原因**：
- 知识库内容太少
- 查询过于笼统
- 向量未更新

**解决方案**：
1. 添加更多知识条目
2. 使用更具体的查询词
3. 点击"同步知识库"重新生成向量

### Q3: 知识图谱显示不出来？

**排查步骤**：
1. 确保用户有知识条目（至少 1 条）
2. 打开浏览器控制台查看错误
3. 检查 API 返回的数据是否正确

### Q4: 如何更换 Embedding 模型？

修改 `backend/services/rag/embedding.go`：
```go
const (
    qwenEmbeddingModel = "text-embedding-v3"  // 可更换为其他模型
)
```

---

## 版本历史

| 版本 | 日期 | 更新内容 |
|------|------|----------|
| v1.0 | 2025-12-25 | 初始版本，实现基础 RAG 功能 |
| v1.1 | 2025-12-25 | 添加混合检索 (Vector + BM25) |
| v1.2 | 2025-12-25 | 集成 Qwen Embedding API |

---

## 参考资料

- [Retrieval-Augmented Generation (RAG)](https://arxiv.org/abs/2005.11401)
- [BM25 算法详解](https://en.wikipedia.org/wiki/Okapi_BM25)
- [阿里云 DashScope 文档](https://help.aliyun.com/zh/dashscope/)
- [Qwen Embedding 模型](https://help.aliyun.com/zh/dashscope/developer-reference/text-embedding-api-details)

---

> 📝 **文档维护**：本文档随代码更新同步维护  
> 📧 **问题反馈**：请在 GitHub Issues 中提交
