# Learning Assistant Backend

基于 Gin + GORM 的学习助手后端服务

## 项目结构

```
backend/
├── config/          # 配置管理
├── database/        # 数据库连接和迁移
├── handlers/        # 请求处理器 (预留)
├── middleware/      # 中间件
├── models/          # 数据模型
├── routes/          # 路由配置
├── utils/           # 工具函数 (预留)
├── main.go          # 程序入口
├── go.mod           # Go 模块文件
├── .env.example     # 环境变量示例
└── README.md        # 项目说明
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并修改相应配置：

```bash
cp .env.example .env
```

### 3. 配置数据库

确保 MySQL 数据库已启动，并创建对应的数据库：

```sql
CREATE DATABASE learning_assistant CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 运行项目

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动

## API 文档

### Swagger UI 文档

项目已集成 Swagger UI，启动服务后可以通过以下方式查看完整的 API 文档：

**访问地址：** `http://localhost:8080/swagger/index.html`

**功能特性：**

- 📖 查看所有 API 接口的详细信息
- 🔍 查看请求参数、响应格式和数据模型
- ✅ 直接在页面上测试 API（点击 "Try it out" 按钮）
- 🏷️ 接口按功能模块分组（认证、任务管理、用户管理、学习室等）

**如何使用 Swagger UI 测试接口：**

1. **选择接口**：点击任意接口展开详情
2. **点击 "Try it out"**：右侧会出现可编辑的参数输入框
3. **填写参数**：
   - Path 参数：直接在路径中填写（如 `:id`, `:userId`）
   - Query 参数：填写查询参数
   - Request Body：在文本框中编辑 JSON 请求体
   - Headers：如需认证，在 `Authorization` 处填写 `Bearer {your_token}`
4. **点击 "Execute"**：发送实际请求到服务器
5. **查看响应**：页面会显示：
   - Response Code（状态码）
   - Response Body（响应内容）
   - Response Headers（响应头）
   - Request URL（实际请求地址）

**认证说明：**

部分接口需要认证，使用方式：

1. 先调用 `/api/v1/auth/login` 登录获取 token
2. 点击页面右上角的 🔒 **Authorize** 按钮
3. 在弹窗中输入：`Bearer {your_token}`（注意 Bearer 和 token 之间有空格）
4. 点击 Authorize 完成认证
5. 之后所有需要认证的接口都会自动带上 token

#### 示例：测试用户注册接口

```json
POST /api/v1/auth/register
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "display_name": "测试用户"
}
```

点击 Execute 后，Swagger UI 会直接向 `http://localhost:8080` 发送请求并显示结果。

> ⚠️ **注意：** Swagger UI 发送的是真实的 HTTP 请求，会实际修改数据库数据，请谨慎操作！

### API 端点概览

#### 认证模块

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户退出
- `POST /api/v1/auth/refresh` - 刷新访问令牌

#### 任务管理

- `GET /api/v1/tasks` - 获取任务列表
- `POST /api/v1/tasks` - 创建任务
- `GET /api/v1/tasks/:id` - 获取任务详情
- `PUT /api/v1/tasks/:id` - 更新任务
- `DELETE /api/v1/tasks/:id` - 删除任务
- `POST /api/v1/tasks/:id/complete` - 完成任务
- `GET /api/v1/tasks/statistics` - 获取任务统计

#### 用户管理

- `GET /api/v1/users/:userId` - 获取用户资料
- `GET /api/v1/users/:userId/study-stats` - 获取学习统计
- `POST /api/v1/users/:userId/check-in` - 每日签到
- `GET /api/v1/users/:userId/achievements` - 获取用户成就

#### 学习室

- `GET /api/v1/study/rooms` - 获取学习房间列表
- `POST /api/v1/study/rooms` - 创建学习房间
- `GET /api/v1/study/rooms/:roomId` - 获取房间详情
- `POST /api/v1/study/rooms/:roomId/join` - 加入学习房间

#### AI 功能

- `POST /api/v1/tasks/ai/parse` - AI 解析任务
- `POST /api/v1/tasks/ai/guidance` - AI 任务指导
- `POST /api/v1/tasks/ai/quiz` - AI 生成测验

#### 健康检查

- `GET /health` - 服务健康状态检查

> 💡 **提示：** 更多详细的接口信息、参数说明和示例请访问 Swagger UI 文档

## 数据模型

项目包含以下主要数据模型：

- **User** - 用户
- **UserProfile** - 用户档案
- **Team** - 团队
- **TeamMember** - 团队成员
- **Task** - 任务
- **TaskCategory** - 任务分类
- **TaskAssignee** - 任务分配
- **StudyRoom** - 学习室
- **ChatMessage** - 聊天消息
- **LearningRecord** - 学习记录
- **PointsLedger** - 积分账本

## 开发说明

- 使用 GORM 进行数据库操作
- 使用 Gin 作为 Web 框架
- 支持自动数据库迁移
- 包含 CORS 和日志中间件
- 遵循 RESTful API 设计规范

## 环境变量说明

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| SERVER_PORT | 服务器端口 | 8080 |
| GIN_MODE | Gin 运行模式 | debug |
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| DB_USERNAME | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | - |
| DB_DATABASE | 数据库名称 | learning_assistant |
| DB_CHARSET | 数据库字符集 | utf8mb4 |
| **QWEN_API_KEY** | **通义千问 API 密钥（AI 功能）** | - |
| JWT_SECRET | JWT 密钥（预留） | - |
| REDIS_HOST | Redis 主机（预留） | localhost |
| REDIS_PORT | Redis 端口（预留） | 6379 |

### AI 服务配置

项目集成了阿里云通义千问 AI 服务，用于以下功能：

- 🤖 **AI 解析任务** - 自动从自然语言中提取任务信息
- 📝 **AI 任务指导** - 生成学习指导和建议
- 📊 **AI 生成测验** - 根据主题自动生成测验题目

**配置步骤：**

1. **获取 API Key**
   - 访问：[阿里云 DashScope 控制台](https://dashscope.console.aliyun.com/)
   - 注册/登录阿里云账号
   - 创建 API Key

2. **配置环境变量**

   在 `.env` 文件中添加：

   ```bash
   QWEN_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxx
   ```

3. **测试 AI 功能**
   
   在 Swagger UI 中测试：
   - `POST /api/v1/tasks/ai/parse` - 测试任务解析
   - `POST /api/v1/tasks/ai/guidance` - 测试任务指导
   - `POST /api/v1/tasks/ai/quiz` - 测试测验生成

**注意事项：**

- ⚠️ 如果未配置 `QWEN_API_KEY`，AI 接口会返回默认/模拟数据，不会调用真实 AI 服务
- 💰 调用 AI API 会产生费用，请注意控制使用量
- 🔒 API Key 是敏感信息，请勿提交到代码仓库
- 🌐 需要网络连接才能访问 AI 服务

**测试示例：**

```bash
# AI 解析任务示例（正确路径）
curl -X POST http://localhost:8080/api/v1/tasks/ai/parse \
  -H "Content-Type: application/json" \
  -d '{
    "input": "明天下午3点开会，讨论项目进度"
  }'

# 返回解析结果：标题、时间、描述等结构化信息
```
