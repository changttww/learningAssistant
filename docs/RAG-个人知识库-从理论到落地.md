---
md041: false
md051: false
---

# 基于 RAG 的个人知识库实现：从理论到落地

> 适用工程：`learningAssistant`（Go 后端 + Vue3 前端）
>
> 本文目标：把“RAG 的理论体系”与“本项目如何落地实现”打通，形成可对外展示的实训项目亮点文档。内容包含：设计动机、端到端链路、关键数据结构/接口、落地代码位置、已完成的产品化能力（技能雷达/学习趋势/知识点统计）、以及可扩展方向。

---

## 目录

- [1. 项目背景与目标](#1-项目背景与目标)
- [2. RAG 理论：从概念到可工程化的最小闭环](#2-rag-理论从概念到可工程化的最小闭环)
  - [2.1 RAG 的核心定义](#21-rag-的核心定义)
  - [2.2 RAG 的端到端流程（Ingest → Embed → Retrieve → Generate）](#22-rag-的端到端流程ingest--embed--retrieve--generate)
  - [2.3 为什么 RAG 适合个人知识库？](#23-为什么-rag-适合个人知识库)
  - [2.4 工程化难点与关键策略](#24-工程化难点与关键策略)
- [3. 本项目的 RAG 总体架构](#3-本项目的-rag-总体架构)
  - [3.1 后端模块划分（Go）](#31-后端模块划分go)
  - [3.2 前端模块划分（Vue3）](#32-前端模块划分vue3)
  - [3.3 两条主链路：问答型 RAG vs 分析型 RAG](#33-两条主链路问答型-rag-vs-分析型-rag)
- [4. 数据模型与知识库构建：个人知识如何进入系统](#4-数据模型与知识库构建个人知识如何进入系统)
  - [4.1 知识条目与元数据（Knowledge Entry）](#41-知识条目与元数据knowledge-entry)
  - [4.2 Chunking（分段）策略建议](#42-chunking分段策略建议)
  - [4.3 Embedding（向量化）与索引](#43-embedding向量化与索引)
  - [4.4（实训重点）本项目的向量化与余弦相似度检索：从理论到代码一一对应](#44实训重点本项目的向量化与余弦相似度检索从理论到代码一一对应)
  - [4.5（答辩速查）RAG 理论步骤 → 本项目代码落点对照表](#45答辩速查rag-理论步骤--本项目代码落点对照表)
  - [4.6（经验沉淀）关键参数与工程取舍：阈值、维度、批量、降级](#46经验沉淀关键参数与工程取舍阈值维度批量降级)
- [5. 检索（Retrieval）：让知识库“可被找出来”](#5-检索retrieval让知识库可被找出来)
  - [5.1 查询向量化与 TopK 召回](#51-查询向量化与-topk-召回)
  - [5.2 过滤与多租户隔离（user_id）](#52-过滤与多租户隔离user_id)
  - [5.3 重排与融合（可选增强）](#53-重排与融合可选增强)
- [6. 生成（Generation）：让知识库“会说话”](#6-生成generation让知识库会说话)
  - [6.1 Prompt 组装与上下文管理](#61-prompt-组装与上下文管理)
  - [6.2 引用溯源（可解释性）](#62-引用溯源可解释性)
  - [6.3 安全与合规（必须做的护栏）](#63-安全与合规必须做的护栏)
  - [6.4（答辩展示）检索增强生成的流程时序（文字版）](#64答辩展示检索增强生成的流程时序文字版)
  - [6.5（可直接套用）Prompt 结构模板：系统指令 + 检索片段 + 输出格式](#65可直接套用prompt-结构模板系统指令--检索片段--输出格式)
- [7. 分析型 RAG（本项目亮点）：从知识库生成结构化学习画像](#7-分析型-rag本项目亮点从知识库生成结构化学习画像)
  - [7.1 为什么“结构化分析”是实训亮点](#71-为什么结构化分析是实训亮点)
  - [7.2 知识分布（knowledge_distribution）](#72-知识分布knowledge_distribution)
  - [7.3 技能雷达（skill_radar）：固定 6 维度 + 分类映射兜底](#73-技能雷达skill_radar固定-6-维度--分类映射兜底)
  - [7.4 学习趋势（learning_trends）：弃用时长，改为三计数指标](#74-学习趋势learning_trends弃用时长改为三计数指标)
  - [7.5 范围聚合与补零：30 天（日）/ 90 天（周）/ 年度（月）](#75-range-aggregation)
- [8. 前端落地展示（Home 首页）：从接口到 ECharts 可视化](#8-前端落地展示home-首页从接口到-echarts-可视化)
  - [8.1 接口解包与兼容（{code,data}）](#81-接口解包与兼容codedata)
  - [8.2 总知识点数：从知识分布求和（替代学习时长）](#82-总知识点数从知识分布求和替代学习时长)
  - [8.3 趋势图：仅做 label 格式化（聚合交给后端）](#83-趋势图仅做-label-格式化聚合交给后端)
- [9. 代码与接口索引（工程内定位）](#9-代码与接口索引工程内定位)
- [10. 可扩展方向（下一阶段可做成加分项）](#10-可扩展方向下一阶段可做成加分项)

---

## 1. 项目背景与目标

在学习辅助场景中，用户的知识碎片来源多样：课堂笔记、错题总结、项目经验、阅读摘录、搜索资料等。传统“笔记系统”的痛点是：**积累容易、检索困难、复用不强**。因此本项目将“个人知识库”作为核心资产，并引入 RAG：

- **让知识库可检索**：不仅按标题/标签查，还能按语义相似度查；
- **让知识库可生成**：从个人知识中生成人设化总结、学习建议、答疑；
- **让知识库可度量**：沉淀为学习画像（知识分布、技能雷达、学习趋势）。

> 实训亮点：本项目不仅做了“问答型 RAG”，还做了“分析型 RAG / 学习画像生成”，并且已经产品化落地到首页数据面板。

---

## 2. RAG 理论：从概念到可工程化的最小闭环

### 2.1 RAG 的核心定义

RAG（Retrieval-Augmented Generation）即“检索增强生成”。它把大模型的回答过程拆成两部分：

1) **检索（Retrieval）**：从知识库中找出与问题相关的内容片段；
2) **生成（Generation）**：将检索到的内容作为上下文（context），交给大模型生成答案/总结/推理结果。

本质上，RAG 是让模型在回答时“带着证据说话”。

### 2.2 RAG 的端到端流程（Ingest → Embed → Retrieve → Generate）

1) **Ingest（入库）**
   - 输入来源：用户新增知识条目、同步笔记、文档解析等。
   - 数据处理：清洗、去噪、结构化元数据（标签/分类/来源/时间/用户）。

2) **Chunking（分段）**
   - 把长文本切成更适合检索的小段（chunk），减少“召回无关内容”的概率。

3) **Embedding（向量化）**
   - 将 chunk 转为向量表示（embedding），用于语义相似度检索。

4) **Index（索引/向量库）**
   - 将向量写入向量索引（可以是独立向量数据库，也可以是本地/数据库表存储 + 相似度计算）。

5) **Retrieve（召回）**
   - 把用户 query 也向量化；
   - 用相似度（cosine/dot）做 TopK 召回；
   - 同时按 `user_id`、`category` 等进行过滤（多租户隔离）。

6) **Augment + Generate（组装上下文并生成）**
   - 将召回内容拼接成 prompt context；
   - 交给 LLM 输出答案/总结/建议；
   - 可选：输出引用来源（知识条目 ID/片段位置）。

### 2.3 为什么 RAG 适合个人知识库？

- **避免幻觉**：大模型易编造，RAG 让回答基于个人知识库事实。
- **个性化**：知识库是“我的内容”，回答自然更贴合用户背景。
- **持续进化**：用户越用越积累，知识库越强，系统越聪明。
- **可解释**：可以展示“我用到了哪些笔记/知识点”。

### 2.4 工程化难点与关键策略

1) **多租户隔离**：必须严格按 `user_id` 查询/检索，保证隐私与数据隔离。
2) **增量更新**：新增/编辑/删除知识条目，需要同步更新 embedding 与索引。
3) **可追溯与可解释**：生成结果最好返回引用来源，便于验证。
4) **成本与性能**：向量化与检索的耗时、模型调用成本需要控制。
5) **口径统一**：统计画像指标必须“来自知识库事实”，避免前端 mock 或口径漂移。

---

## 3. 本项目的 RAG 总体架构

### 3.1 后端模块划分（Go）

本项目后端采用 Go（Gin + GORM），RAG 主要落在：

- `backend/services/rag/`
  - `rag_service.go`：RAG 主流程服务（检索 + 上下文组装 + 调用模型）。
  - `embedding.go`：向量化/embedding 相关封装。
  - `ai_analysis.go`：知识库分析与学习画像生成（分布/雷达/趋势）。

路由层：

- `backend/routes/knowledge_base.go`：知识库分析、技能雷达、趋势等 API。
- 其他 AI/聊天相关：如 `backend/routes/ai.go`、`backend/routes/study_chat.go`、`backend/routes/study_ws.go` 等（具体以工程为准）。

### 3.2 前端模块划分（Vue3）

- API 封装：`front/src/api/modules/knowledge.js`
  - 负责调用 `/knowledge-base/...` 等接口
- 首页展示：`front/src/views/Home.vue`
  - ECharts 渲染知识分布饼图、技能雷达、学习趋势
  - 总知识点数（由知识分布求和）
- 知识图谱（新增）：`front/src/views/KnowledgeGraph.vue`
- RAG 问答机器人（新增）：`front/src/views/KnowledgeChat.vue`

### 3.3 两条主链路：问答型 RAG vs 分析型 RAG

- **问答型 RAG**：用户发起问题 → 检索知识库 → 生成答案。
- **分析型 RAG（本项目亮点）**：对用户知识库做结构化分析 → 产出可视化指标（分布/雷达/趋势）→ 首页展示学习画像。

> 分析型 RAG 的价值：把“知识库”从纯内容仓库升级为“学习画像引擎”，这是实训项目更容易展示的亮点。

---

## 4. 数据模型与知识库构建：个人知识如何进入系统

### 4.1 知识条目与元数据（Knowledge Entry）

知识条目是 RAG 的事实来源。虽然字段以代码为准，但本项目整体依赖如下关键属性：

- `user_id`：条目归属（多租户隔离的核心）
- `title`/`content`：主要文本
- `category`：分类（用于分布统计、技能映射、过滤检索）
- `tags`/`keywords`：标签（用于检索增强与分类兜底）
- `created_at`：用于趋势统计

模型位置：`backend/models/knowledge_base.go`

### 4.2 Chunking（分段）策略建议

为了让检索更准确，入库文本通常需要分段。推荐策略：

- 按语义段落切：优先按标题/小节/段落分割
- 控制 chunk 长度：例如 200～800 字（或 300～1,000 tokens）
- chunk 重叠：如 10%～20% overlap，避免跨段信息丢失
- chunk 元数据：保存 `entry_id`、`category`、`tags`、`created_at`、`user_id`

> 注：本工程已经具备 embedding/service 层结构，chunking 可以作为下一阶段重要增强点（也是展示 RAG 专业度的加分项）。

### 4.3 Embedding（向量化）与索引

代码位置：`backend/services/rag/embedding.go`

落地要点：

1) **向量生成**：将文本 chunk 调用 embedding 模型得到向量。
2) **索引存储**：可存入向量数据库或本地表（实现方式以工程为准）。
3) **增量更新**：条目新增/编辑触发 embedding 更新；删除触发索引删除。

---

## 4.4（实训重点）本项目的向量化与余弦相似度检索：从理论到代码一一对应

这一节专门把 RAG 的 **Embed（向量化）** 与 **Retrieve（检索）** 两步，映射到本项目已经落地的代码实现，便于答辩时“拿着代码讲原理”。

### 4.4.1 EmbeddingService 抽象：把理论中的“向量化能力”做成可替换接口

对应理论：RAG 流程中的 **Embedding（向量化）**。

本项目在 `backend/services/rag/rag_service.go` 定义了 `EmbeddingService` 接口：

- `GenerateEmbedding(text) -> Vector`
- `GenerateEmbeddings(texts) -> []Vector`
- `CosineSimilarity(vec1, vec2) -> float32`

意义：

- 把“向量生成”和“相似度度量”抽象出来，使 RAG 主流程不绑定某个具体模型/供应商。
- 也为本地开发提供可降级实现（没有 API Key 也能跑通链路）。

### 4.4.2 QwenEmbeddingService：Embedding 的工程落地（真实模型 + 本地 mock）

对应代码：`backend/services/rag/embedding.go`

1) **真实 embedding 调用**

- `NewQwenEmbeddingService()` 默认使用 DashScope 的 `text-embedding-v3`
- `GenerateEmbeddings()` 支持批量向量化，并做了“单次最多 10 条”的分批控制（防止请求过大）
- `callQwenAPI()` 通过 HTTP POST 请求 embedding API，解析响应中的 `embedding` 数组。

这对应理论中的：

- “把文本映射到高维向量空间”
- “向量在同一语义空间中可比较（计算相似度）”。

1) **本地开发 mock 向量**

当 `QWEN_API_KEY` 缺失时，`GenerateEmbedding()` 会走 `mockEmbedding(text)`：

- 固定维度 `dim := 1536`（与 Qwen embedding 维度一致）
- 基于字符哈希 + `sin()` 生成伪随机、可复现实验的向量

意义：

- 保障“RAG 链路可演示/可联调”，即使没有真实模型 Key 也能跑通。
- 对实训来说，能展示工程的健壮性：线上用真实 embedding，线下可 mock。

> 答辩表达建议：强调我们不仅“接了模型”，还做了接口抽象与可降级策略，保证工程可运行、易扩展。

1) **LocalEmbeddingService：完全离线的 embedding 兜底**

同样在 `embedding.go` 提供 `LocalEmbeddingService`：

- `simpleEmbedding()` 使用 256 维词频/字符哈希特征向量
- 并对向量做归一化（L2 norm）

意义：

- 用于单元测试/无网环境/无 Key 环境
- 体现“embedding 的本质是把文本映射到可比较的向量空间”，而不是绑定某个供应商。

### 4.4.3 我们是否实现了“余弦相似度”？实现了，而且是检索核心

对应代码：`backend/services/rag/embedding.go`

- `QwenEmbeddingService.CosineSimilarity(vec1, vec2)`
- `LocalEmbeddingService.CosineSimilarity(vec1, vec2)`

实现原理（对应理论）：

- 余弦相似度衡量两个向量夹角：
  - 夹角越小（方向越接近）→ 语义越相似
  - 与向量长度无关，更适合语义向量比较

实现公式（代码体现）：

- 点积：`dotProduct += vec1[i] * vec2[i]`
- 模长：`norm = sqrt(sum(vec[i]^2))`
- `similarity = dot / (||a|| * ||b||)`

### 4.4.4 向量缓存（Vector Cache）：把“Embed 结果”落到数据库，提升检索性能

对应代码：`backend/services/rag/rag_service.go`

在 `AddDocument()` 中：

- 对入库文本先做清洗：`stripHTMLTags()` 去 HTML、压空白
- 生成摘要与关键词（用于结构化字段与降级检索）
- 调用 `GenerateEmbedding(cleanTitle + " " + summary)` 生成向量
- 写入 `models.KnowledgeVectorCache`（向量缓存表）

关键字段（从代码可见的使用方式）：

- `EntryID`：对应知识条目
- `ContentHash`：内容哈希（便于判断更新/去重）
- `Vector`：向量本体
- `VectorDim`：向量维度
- `VectorModel`：向量模型标识（例如 `qwen-embedding`）

对应理论中的：

- “索引/向量库（Index）”
- 本项目采用“数据库向量缓存表”的方式，属于轻量级向量索引实现（适合实训与中小规模数据）。

### 4.4.5 向量检索（Retrieve）：TopK + 阈值过滤 + 多租户隔离 + 降级策略

对应代码：`backend/services/rag/rag_service.go`

#### 1) 多租户隔离（RAG 必做）

在 `vectorSearch(userID, queryVector, limit)`：

- 先查 `knowledge_base_entries`：`Where("user_id = ? AND status = 1", userID)`
- 再用这些 entry_id 去取向量缓存 `knowledge_vector_caches`

这直接对应理论中的：

- “以 user_id 做过滤（tenant isolation）”
- 确保检索只发生在“我的知识库”。

#### 2) 相似度计算（余弦）与阈值过滤

- 对每条缓存向量计算：`similarity := CosineSimilarity(queryVector, cache.Vector)`
- 过滤阈值：`if similarity > 0.3 { ... }`

对应理论中的：

- “向量空间相似度召回”
- “阈值用于减少低相关噪声”。

> 实训表达建议：阈值 0.3 是工程经验值，可在日志/实验中调参；后续可做动态阈值或 TopK 不限阈值。

#### 3) TopK 召回

当前实现逻辑：

- 先计算相似度并收集候选
- 若无候选，返回空数组
- 然后取前 `limit` 条作为召回结果

> 注意：当前代码里“按相似度排序”部分有注释提示为简化版：
> `// 按相似度排序 // 这里简化处理，实际可用更完善的排序`
>
> 这点可以在答辩中如实说明，并作为后续优化点：对候选按 score 排序后再取 TopK。

#### 4) 降级策略：向量检索失败/无结果 → 关键词检索

在 `SearchKnowledge()` 中：

- 先尝试向量检索（Embedding + Cosine）
- 若失败或结果为空，则降级到关键词 LIKE：
  - `title LIKE ? OR content LIKE ? OR keywords LIKE ? OR category LIKE ?`

对应理论中的：

- “混合检索（Dense + Sparse）”的简化版本
- 保障系统稳定性：即使 embedding 不可用，也能查到内容。

### 4.5（答辩速查）RAG 理论步骤 → 本项目代码落点对照表

下面这张表可以直接用在实训答辩 PPT：把 RAG 理论的每一步，映射到本项目的真实代码与数据结构，证明我们“不是只会讲概念”，而是完成了工程闭环。

| RAG 理论步骤 | 这一步在做什么 | 本项目落点（文件/函数） | 数据/产物 | 备注（我们怎么做的） |
| --- | --- | --- | --- | --- |
| Ingest（入库） | 接收用户知识内容、清洗与结构化 | `backend/services/rag/rag_service.go` → `AddDocument()` | `KnowledgeBaseEntry`（标题/内容/摘要/关键词/分类） | `stripHTMLTags()` 去 HTML；生成 `summary/keywords/category` 用于后续检索与统计 |
| Embed（向量化） | 文本 → 向量（语义空间表达） | `backend/services/rag/embedding.go` → `GenerateEmbedding/GenerateEmbeddings` | `models.Vector` | 默认走 `QwenEmbeddingService`；无 Key 时可 `mockEmbedding()` 或 `LocalEmbeddingService` 兜底 |
| Index（索引/向量库） | 存储向量，供检索快速使用 | `AddDocument()` 写 `models.KnowledgeVectorCache` | `KnowledgeVectorCache`（EntryID/Vector/Dim/Model/Hash） | 采用“数据库向量缓存表”作为轻量索引实现，适合实训与中小数据量 |
| Retrieve（召回） | 用相似度从知识库找 TopK 相关内容 | `rag_service.go` → `SearchKnowledge()` + `vectorSearch()` | `[]KnowledgeBaseEntry` | 先向量检索（余弦相似度）→ 无结果/失败再降级关键词 LIKE（稳定性策略） |
| Similarity（相似度度量） | 计算 query 与候选的相似程度 | `embedding.go` → `CosineSimilarity()` | `float32 similarity` | 明确实现余弦相似度：点积/模长/归一化；并设置阈值（默认 `>0.3`） |
| TopK（选取结果） | 取最相关的 K 条作为上下文 | `vectorSearch(userID, queryVector, limit)` | TopK entries | 当前实现有“简化排序”注释，可作为优化点：按 score 排序后取 TopK |
| Filter（过滤） | 过滤非本人数据 / 非发布状态 | `vectorSearch()` DB 查询条件 | 过滤后的候选集合 | 通过 `user_id` + `status=1` 实现多租户隔离与数据可用性控制 |
| Generate（生成） | 使用召回片段增强 prompt，再让模型生成 | （工程存在 AI/聊天路由）`backend/routes/ai.go` / `study_chat.go` / `study_ws.go` | 生成回答/总结 | 当前文档重点展示“检索增强能力 + 分析型 RAG”，问答生成链路可在答辩时展示调用接口与 prompt 组装 |
| Analytics（分析型 RAG） | 用知识库事实生成结构化学习画像 | `backend/services/rag/ai_analysis.go` | `knowledge_distribution/skill_radar/learning_trends` | 已产品化落地首页：分布、技能雷达、趋势（30/90/年度聚合与补零） |

### 4.6（经验沉淀）关键参数与工程取舍：阈值、维度、批量、降级

这一部分用于解释“为什么我们这么实现”，体现工程取舍与性能/质量思维（非常适合实训亮点表达）。

#### 1) 相似度阈值：`similarity > 0.3`

- **代码位置**：`backend/services/rag/rag_service.go` → `vectorSearch()`
- **理论对应**：Retrieve 阶段的“过滤低相关候选”，减少噪声。
- **工程取舍**：
  - 阈值过低：召回更多但噪声大（会把无关条目也当上下文）
  - 阈值过高：召回变少，容易查不到内容
- **我们的策略**：先用经验阈值 0.3（可调参），结合“无结果降级关键词检索”保证可用性。

#### 2) 向量维度选择：1536（真实 embedding） vs 256（本地离线 embedding）

- **代码位置**：
  - `QwenEmbeddingService.mockEmbedding()`：`dim := 1536`（与真实模型一致）
  - `LocalEmbeddingService.simpleEmbedding()`：`dim := 256`
- **理论对应**：Embedding 向量维度越高表达能力越强，但存储/计算成本越高。
- **我们的策略**：
  - 线上/真实环境：用 Qwen embedding（高维、语义能力强）
  - 线下/无 Key 环境：用本地 embedding 提供“可运行/可演示/可测试”的链路兜底

#### 3) 批量 embedding 的分批控制：每批最多 10 条

- **代码位置**：`backend/services/rag/embedding.go` → `GenerateEmbeddings()`
- **理论对应**：Ingest/Embed 的工程化处理（吞吐、成本、稳定性）。
- **工程取舍**：
  - 一次向量化太多：请求体大、超时/失败概率增加
  - 分批太小：请求次数多、吞吐降低
- **我们的策略**：固定 10 条/批，优先保证稳定性与可控性（可按实际模型限额再调）。

#### 4) 向量索引的实现方式：向量缓存表（轻量向量库）

- **代码落点**：`models.KnowledgeVectorCache`（通过 `AddDocument()` 写入）
- **理论对应**：Index（向量库/索引）。
- **工程取舍**：
  - 专用向量数据库：性能/功能强，但引入成本与运维复杂
  - 数据库缓存表：实现简单、易部署，适合实训与中小数据
- **我们的策略**：先用缓存表实现最小闭环；后续可无缝升级到专业向量库（保持 `EmbeddingService`/检索层接口即可）。

#### 5) 降级策略：向量检索失败/无结果 → 关键词检索

- **代码位置**：`rag_service.go` → `SearchKnowledge()`
- **理论对应**：混合检索（Dense + Sparse）的工程化版本。
- **我们的策略**：
  - 优先向量检索保证语义能力
  - 兜底关键词 LIKE 保证稳定性与可用性

---

## 5. 检索（Retrieval）：让知识库“可被找出来”

### 5.1 查询向量化与 TopK 召回

最小可用检索链路：

- query → embedding
- 与知识库向量做相似度计算
- 取 TopK 作为上下文

### 5.2 过滤与多租户隔离（user_id）

个人知识库必须保证：

- 检索时 **只检索该用户的向量/条目**
- 统计分析时 **只统计该用户的条目**

这也是本项目在统计口径上反复对齐的原因：指标必须来自“我的知识库”。

### 5.3 重排与融合（可选增强）

为了提升检索质量，可做：

- 关键词召回（BM25） + 向量召回（Dense）融合
- 引入 reranker（重排模型）提升 TopK 相关性
- 按 `category/tags` 做偏好加权（个性化）

---

## 6. 生成（Generation）：让知识库“会说话”

### 6.1 Prompt 组装与上下文管理

RAG 的 prompt 一般包含：

- **系统指令**：语气、边界、格式
- **用户问题**
- **检索片段集合**（含来源信息）
- **输出格式**：Markdown、要点列表、引用要求等

代码承载：`backend/services/rag/rag_service.go`（以工程为准）。

### 6.2 引用溯源（可解释性）

推荐输出结构：

- `answer`：模型输出
- `citations`：引用列表（entry_id、标题、片段摘要）

这样可以在前端展示“本次回答参考了哪些知识点”，增强可信度。

### 6.3 安全与合规（必须做的护栏）

- 不返回跨用户数据
- 对敏感信息做脱敏/过滤
- 防 prompt 注入：对上下文片段做边界封装，限制其影响系统指令

### 6.4（答辩展示）检索增强生成的流程时序（文字版）

这一节可以直接作为答辩讲稿：用“时序/流水线”的方式解释 RAG（Retrieve + Generate）如何在本项目中执行。虽然不同业务路由（学习聊天/AI 助手/学习空间）入口不同，但核心链路一致。

**参与方（逻辑组件）**：

- 前端页面（如聊天页/学习空间/首页触发的 AI 分析）
- 后端路由（`backend/routes/*`：如 `ai.go`、`study_chat.go`、`study_ws.go` 等）
- RAG 服务（`backend/services/rag/rag_service.go`）
- Embedding 服务（`backend/services/rag/embedding.go`）
- 数据库（知识条目表 `knowledge_base_entries` + 向量缓存表 `knowledge_vector_caches`）
- 大模型服务（生成模型，具体供应商按配置）

**时序步骤（Retrieve → Augment → Generate）**：

1) **用户发起问题/请求**
   - 输入：`query`（自然语言问题），以及当前用户身份（`user_id`）

2) **后端进入 RAG 检索阶段（Retrieve）**
   - 调用 `EmbeddingService.GenerateEmbedding(query)`（见 `embedding.go`）
   - 得到 `queryVector` 后，进入 `vectorSearch(userID, queryVector, limit)`（见 `rag_service.go`）

3) **向量召回（TopK + Filter）**
   - 过滤：只从 `user_id` + `status=1` 的知识条目范围取向量缓存，保证多租户隔离
   - 相似度：对每个缓存向量计算 `CosineSimilarity(queryVector, cache.Vector)`
   - 阈值：过滤 `similarity > 0.3` 的候选（可调参）
   - TopK：取最多 `limit` 条作为召回结果

4) **降级策略（Fallback）**
   - 若向量召回失败/结果为空：`SearchKnowledge()` 自动降级到关键词 LIKE 检索（title/content/keywords/category）
   - 目的：保证“有网无网、有 Key 无 Key、数据量大小”情况下都可用、可演示

5) **上下文组装（Augment）**
   - 将召回的知识条目（title/summary/content 等）拼成“检索片段集合”
   - 组装为 prompt 的 Context 区域（见下节 6.5 模板）

6) **调用生成模型（Generate）**
   - 输入：系统指令 + 用户问题 + 检索片段集合 + 输出格式要求
   - 输出：回答/总结/建议（可进一步扩展为 citations）

7) **返回前端渲染**
   - 前端展示：模型回答
   - 可选：展示引用（来源条目标题/ID），增强可信度

> 对应理论映射：
>
> - 2~4 步是 Retrieve（召回）
> - 5 步是 Augment（上下文增强）
> - 6 步是 Generate（生成）

### 6.5（可直接套用）Prompt 结构模板：系统指令 + 检索片段 + 输出格式

这一节给出“可直接用于实现/答辩展示”的 prompt 模板，说明我们为什么要做检索片段封装，以及如何保证模型回答可控。

> 注意：本项目代码中 Generation 入口可能分布在多个路由/服务层；模板用于描述“我们生成阶段的标准做法/规范”。

#### 6.5.1 Prompt 的推荐结构

- **System（系统指令）**：约束模型角色、输出格式、禁止编造
- **User（用户问题）**：原始 query
- **Context（检索片段）**：来自我们知识库召回的 TopK 内容
- **Output format（输出格式）**：要求 Markdown/要点/引用等

#### 6.5.2 模板（示例，可复制到实现中）

System：

- 你是一名学习助手，只能基于给定的“知识库片段”回答。
- 如果知识库片段不足以回答，请明确说明“不足以回答”，并给出你需要的补充信息。
- 不允许编造不存在的概念、数据或来源。

User：

- 问题：`{user_query}`

Context（RAG 检索结果，来自系统内部知识库）：

- [片段 1] entry_id=`{id1}` title=`{t1}`
  - summary: `{summary1}`
  - content: `{chunk_or_content1}`
- [片段 2] entry_id=`{id2}` title=`{t2}`
  - summary: `{summary2}`
  - content: `{chunk_or_content2}`
- ...（最多 TopK 条）

Output format：

1) 先给结论（不超过 5 行）
2) 再给依据：引用你使用到的片段编号（如引用 [片段 2]）
3) 最后给建议/下一步学习行动（3 条以内）

#### 6.5.3 为什么模板里要强调“引用片段编号/来源”？

对应理论：RAG 的核心优势之一是 **可解释性**。

- 若回答能引用 `[片段 X]`，用户可追溯到具体知识条目，增强可信度。
- 也便于后续做“citation 输出结构化字段”（例如 `citations: [{entry_id, title}]`）。

---

## 7. 分析型 RAG（本项目亮点）：从知识库生成结构化学习画像

### 7.1 为什么“结构化分析”是实训亮点

和常见的“做一个能聊天的 RAG”不同，本项目将知识库进一步产品化：

- 把知识库条目转为 **可视化指标**
- 在首页展示“学习画像”，让用户看到自己学习结构与趋势

这类能力更容易体现：

- 数据处理能力
- 口径设计能力
- 前后端联动与工程落地能力

代码核心：`backend/services/rag/ai_analysis.go`

### 7.2 知识分布（knowledge_distribution）

定义：按分类汇总知识条目数量（count）和占比（percentage）。

用途：

- 首页饼图展示
- 也是“总知识点数”的可解释来源（对 count 求和）

### 7.3 技能雷达（skill_radar）：固定 6 维度 + 分类映射兜底

本项目技能雷达的关键设计（已落地）：

- **固定 6 维**：保证前端雷达图稳定，不因分类变化导致显示异常
- **从知识库条目计算**：不再依赖模拟数据
- **分类归一化策略**：
  - category 别名映射
  - 标签/标题/关键词兜底

解决了典型问题：

- 用户确有编程知识，但“编程/考试技巧”为 0（归类漏算）

### 7.4 学习趋势（learning_trends）：弃用时长，改为三计数指标

本项目已将趋势口径重构为三条可解释曲线：

- `done_tasks`：当期完成任务数
- `new_notes`：当期创建笔记数
- `new_knowledge`：当期新增知识点数

原因：

- `study_hours` 口径容易产生歧义（跨天、计时来源不统一）
- 计数口径更可解释、数据链路更清晰

<!-- markdownlint-disable MD033 -->

<a id="75-range-aggregation"></a>

<!-- markdownlint-enable MD033 -->

### 7.5 范围聚合与补零：30 天（日）/ 90 天（周）/ 年度（月） (range-aggregation)

本项目已落地“范围聚合交给后端”策略：

- `range=30`：按日聚合（`YYYY-MM-DD`）
- `range=90`：按周聚合（`YYYY-Www`）
- `range=year`：按月聚合（`YYYY-MM`，12 个桶）

并在后端：

- **补齐缺失桶（补零）**，保证 x 轴连续，避免前端抽样导致“末尾点丢失”

对应路由：`/knowledge-base/trends?range=...`

---

## 8. 前端落地展示（Home 首页）：从接口到 ECharts 可视化

### 8.1 接口解包与兼容（{code,data}）

本项目后端接口常见包装结构为 `{code, data}`，前端需要兼容：

- 数组响应解包：`unwrapArrayResponse`
- 报告响应解包：`unwrapReportResponse`

落地位置：`front/src/views/Home.vue`

### 8.2 总知识点数：从知识分布求和（替代学习时长）

实训展示点：**指标来源可解释**。

- 总知识点数不再来自学习统计（也不再展示“总学习时长”）
- 直接从 `knowledge_distribution[].count` 求和得到

优势：

- 与知识库事实一致
- 不依赖不稳定的时长口径

### 8.3 趋势图：仅做 label 格式化（聚合交给后端）

前端目前做的事情很“轻”：

- 请求后端已聚合数据 `/knowledge-base/trends?range=...`
- 根据 range 格式化横轴 label：
  - year：`YYYY-MM` → `M月`
  - 90：`YYYY-Www` → `Wxx`
  - 30：`YYYY-MM-DD` → `M/D`

这样避免：

- 前端切片伪造范围
- 前端抽样导致末尾点丢失

---

## 9. 代码与接口索引（工程内定位）

### 后端

- RAG 服务与分析：
  - `backend/services/rag/rag_service.go`
  - `backend/services/rag/embedding.go`
  - `backend/services/rag/ai_analysis.go`
- 知识库路由：
  - `backend/routes/knowledge_base.go`
- 数据模型：
  - `backend/models/knowledge_base.go`

### 前端

- 知识库 API：
  - `front/src/api/modules/knowledge.js`
- 首页（可视化与指标展示）：
  - `front/src/views/Home.vue`
- 知识图谱（新增）：
  - `front/src/views/KnowledgeGraph.vue`
- RAG 问答机器人（新增）：
  - `front/src/views/KnowledgeChat.vue`

常用接口（以实际路由为准）：

- `GET /knowledge-base/analysis`
- `GET /knowledge-base/skill-radar`
- `GET /knowledge-base/trends?range=30|90|year`
- `GET /knowledge-base/graph` - 获取知识图谱数据（新增）
- `POST /knowledge-base/chat` - RAG 问答（带引用溯源，新增）

---

## 10. 可扩展方向（下一阶段可做成加分项）

1) **~~回答引用溯源（citations）~~** ✅ 已实现
   - 回答返回引用条目 ID，可在前端展示"参考来源"。
   - 代码位置：`backend/routes/knowledge_base.go` → `ragChat()`
   - 前端展示：`front/src/views/KnowledgeChat.vue`

2) **知识图谱可视化** ✅ 已实现
   - 展示知识点之间的关联关系（同分类、显式关系）
   - 代码位置：`backend/services/rag/rag_service.go` → `GetKnowledgeGraph()`
   - 前端展示：`front/src/views/KnowledgeGraph.vue`（ECharts Graph）

3) **混合检索（BM25 + 向量）与重排**
   - 提升召回准确率，降低"答非所问"。

4) **增量 embedding 与队列化**
   - 用异步任务处理向量化，避免阻塞主流程。

5) **知识库质量指标**
   - 统计"重复条目率、过期条目率、低质量条目率"，进一步产品化。

6) **个性化学习建议生成**
   - 将知识分布/趋势作为 prompt 输入，生成"本周复盘 + 下周计划"。

---

## 附：本项目实训亮点总结（可用于答辩/展示）

- **RAG 双形态落地**：既能做问答型 RAG，也能做分析型 RAG（学习画像）。
- **RAG 问答带引用溯源**：用户提问后，展示回答的同时显示"参考了哪些知识点"，体现可解释性。
- **知识图谱可视化**：用 ECharts Graph 展示知识点关联，适合答辩演示。
- **指标口径可解释**：趋势弃用时长，改为可解释的行为计数；总知识点数来自知识分布求和。
- **前后端职责清晰**：聚合与补零交给后端，前端专注展示与格式化。
- **产品闭环**：知识库不是"存完就结束"，而是反哺首页画像与学习趋势，驱动用户持续使用。
