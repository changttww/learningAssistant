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

## API 端点

### 健康检查
- `GET /health` - 服务健康状态检查

### API v1
- `GET /api/v1/users` - 用户相关接口 (待实现)
- `GET /api/v1/tasks` - 任务相关接口 (待实现)
- `GET /api/v1/teams` - 团队相关接口 (待实现)
- `GET /api/v1/study-rooms` - 学习室相关接口 (待实现)

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