package seeder

import (
	"log"
	"time"

	"gorm.io/gorm"

	"learningAssistant-backend/models"
)

// SeedKnowledgeEntries 初始化知识库演示条目
func SeedKnowledgeEntries(db *gorm.DB) error {
	log.Println("开始初始化知识库演示数据...")

	// 为用户ID=1添加演示知识条目
	entries := []models.KnowledgeBaseEntry{
		// 前端开发知识
		{
			UserID:      1,
			SourceType:  3, // 手动添加
			Title:       "Vue 3 Composition API",
			Content:     "Vue 3 的 Composition API 是一种新的组织组件逻辑的方式。使用 setup() 函数、ref、reactive、computed 等来管理组件状态。相比 Options API，它提供了更好的代码复用和类型推断支持。",
			Summary:     "Vue 3 新特性，用于组织组件逻辑的API",
			Category:    "前端开发",
			SubCategory: "Vue.js",
			Level:       3, // 已掌握
			Status:      1,
		},
		{
			UserID:      1,
			SourceType:  3,
			Title:       "JavaScript Promise 异步编程",
			Content:     "Promise 是 JavaScript 处理异步操作的核心机制。它有三种状态：pending、fulfilled、rejected。使用 then()、catch()、finally() 处理结果，async/await 是 Promise 的语法糖。",
			Summary:     "JavaScript异步编程的核心概念",
			Category:    "前端开发",
			SubCategory: "JavaScript",
			Level:       3,
			Status:      1,
		},
		{
			UserID:      1,
			SourceType:  1, // 从任务创建
			Title:       "CSS Flexbox 布局",
			Content:     "Flexbox 是一种一维布局模型，主要用于在容器中排列子元素。关键属性包括：display: flex, flex-direction, justify-content, align-items, flex-wrap 等。",
			Summary:     "CSS弹性盒子布局模型",
			Category:    "前端开发",
			SubCategory: "CSS",
			Level:       2, // 熟悉
			Status:      1,
		},
		{
			UserID:      1,
			SourceType:  3,
			Title:       "React Hooks 使用",
			Content:     "React Hooks 允许在函数组件中使用状态和其他 React 特性。常用 Hooks：useState、useEffect、useContext、useRef、useMemo、useCallback。",
			Summary:     "React函数组件的状态管理方式",
			Category:    "前端开发",
			SubCategory: "React",
			Level:       2,
			Status:      1,
		},

		// 后端开发知识
		{
			UserID:      1,
			SourceType:  3,
			Title:       "Go语言并发编程 - Goroutine",
			Content:     "Goroutine 是 Go 语言的轻量级线程。使用 go 关键字启动，通过 channel 进行通信。select 语句用于多路复用 channel 操作。",
			Summary:     "Go语言并发模型的核心概念",
			Category:    "后端开发",
			SubCategory: "Go",
			Level:       2,
			Status:      1,
		},
		{
			UserID:      1,
			SourceType:  2, // 从笔记创建
			Title:       "RESTful API 设计原则",
			Content:     "REST 是一种架构风格。关键原则：使用HTTP方法(GET/POST/PUT/DELETE)、无状态、资源用URI表示、使用正确的状态码、版本控制。",
			Summary:     "Web API设计的最佳实践",
			Category:    "后端开发",
			SubCategory: "API设计",
			Level:       3,
			Status:      1,
		},
		{
			UserID:      1,
			SourceType:  3,
			Title:       "MySQL索引优化",
			Content:     "索引是提高查询性能的关键。常见索引类型：B+树索引、哈希索引、全文索引。需要注意：避免全表扫描、合理使用复合索引、定期分析执行计划。",
			Summary:     "数据库查询性能优化核心技术",
			Category:    "数据库",
			SubCategory: "MySQL",
			Level:       2,
			Status:      1,
		},

		// 算法知识
		{
			UserID:      1,
			SourceType:  3,
			Title:       "二分查找算法",
			Content:     "二分查找是一种在有序数组中查找元素的高效算法。时间复杂度 O(log n)。关键点：确定边界条件、处理中点计算溢出、理解循环不变量。",
			Summary:     "高效的查找算法",
			Category:    "算法",
			SubCategory: "查找算法",
			Level:       3,
			Status:      1,
		},
		{
			UserID:      1,
			SourceType:  3,
			Title:       "动态规划基础",
			Content:     "动态规划是解决最优化问题的方法。核心思想：将问题分解为子问题，存储子问题的解避免重复计算。关键步骤：定义状态、状态转移方程、边界条件。",
			Summary:     "解决最优化问题的编程范式",
			Category:    "算法",
			SubCategory: "动态规划",
			Level:       1, // 了解
			Status:      1,
		},

		// 软技能
		{
			UserID:      1,
			SourceType:  3,
			Title:       "Git 版本控制",
			Content:     "Git 是分布式版本控制系统。常用命令：git clone/add/commit/push/pull/branch/merge。了解工作区、暂存区、本地仓库、远程仓库的概念。",
			Summary:     "代码版本管理工具",
			Category:    "开发工具",
			SubCategory: "Git",
			Level:       3,
			Status:      1,
		},
		{
			UserID:      1,
			SourceType:  3,
			Title:       "敏捷开发方法论",
			Content:     "敏捷开发强调迭代、快速反馈和适应变化。Scrum 是常用框架：Sprint计划、每日站会、Sprint回顾。看板方法可视化工作流程。",
			Summary:     "现代软件开发流程方法",
			Category:    "软技能",
			SubCategory: "项目管理",
			Level:       2,
			Status:      1,
		},

		// 人工智能
		{
			UserID:      1,
			SourceType:  3,
			Title:       "机器学习基础概念",
			Content:     "机器学习是让计算机从数据中学习的技术。三种类型：监督学习、无监督学习、强化学习。常见算法：线性回归、决策树、神经网络。",
			Summary:     "AI技术的核心领域",
			Category:    "人工智能",
			SubCategory: "机器学习",
			Level:       1,
			Status:      1,
		},
		{
			UserID:      1,
			SourceType:  3,
			Title:       "大语言模型 LLM 原理",
			Content:     "大语言模型基于 Transformer 架构。核心概念：自注意力机制、位置编码、预训练与微调。GPT、BERT、LLaMA 是代表模型。",
			Summary:     "ChatGPT等AI助手的技术基础",
			Category:    "人工智能",
			SubCategory: "NLP",
			Level:       1,
			Status:      1,
		},
	}

	for i := range entries {
		entries[i].CreatedAt = time.Now().Add(-time.Duration(i*2) * 24 * time.Hour)
		entries[i].UpdatedAt = entries[i].CreatedAt
	}

	for _, entry := range entries {
		// 检查是否已存在
		var existing models.KnowledgeBaseEntry
		result := db.Where("user_id = ? AND title = ?", entry.UserID, entry.Title).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&entry).Error; err != nil {
				log.Printf("创建知识条目失败 [%s]: %v", entry.Title, err)
			} else {
				log.Printf("创建知识条目: %s (分类: %s, 等级: %d)", entry.Title, entry.Category, entry.Level)
			}
		}
	}

	log.Println("知识库演示数据初始化完成")
	return nil
}
