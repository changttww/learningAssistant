
package seeder

import (
	"log"

	"gorm.io/gorm"

	"learningAssistant-backend/models"
)

// SeedKnowledgeCategories 初始化知识库分类
func SeedKnowledgeCategories(db *gorm.DB) error {
	categories := []models.KnowledgeCategory{
		{
			Name:        "编程语言",
			Description: "各种编程语言的学习，如Python、Java、Go、JavaScript等",
			Color:       "#3b82f6",
			Icon:        "code",
		},
		{
			Name:        "前端开发",
			Description: "Web前端开发技术，包括HTML、CSS、JavaScript框架等",
			Color:       "#ec4899",
			Icon:        "palette",
		},
		{
			Name:        "后端开发",
			Description: "服务器端开发技术，包括API设计、数据库、框架等",
			Color:       "#f59e0b",
			Icon:        "server",
		},
		{
			Name:        "数据科学",
			Description: "数据分析、机器学习、统计学相关知识",
			Color:       "#8b5cf6",
			Icon:        "chart-bar",
		},
		{
			Name:        "人工智能",
			Description: "AI、深度学习、NLP、计算机视觉等",
			Color:       "#10b981",
			Icon:        "sparkles",
		},
		{
			Name:        "数据库",
			Description: "SQL、NoSQL、数据库设计与优化",
			Color:       "#06b6d4",
			Icon:        "database",
		},
		{
			Name:        "DevOps",
			Description: "持续集成、容器化、云服务、运维自动化",
			Color:       "#6366f1",
			Icon:        "cloud",
		},
		{
			Name:        "算法",
			Description: "算法设计、数据结构、LeetCode刷题",
			Color:       "#ef4444",
			Icon:        "puzzle",
		},
		{
			Name:        "软技能",
			Description: "项目管理、沟通协作、时间管理等",
			Color:       "#84cc16",
			Icon:        "users",
		},
		{
			Name:        "其他",
			Description: "其他学习内容",
			Color:       "#6b7280",
			Icon:        "book",
		},
	}

	for _, category := range categories {
		var existing models.KnowledgeCategory
		result := db.Where("name = ?", category.Name).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&category).Error; err != nil {
				log.Printf("创建知识分类 %s 失败: %v", category.Name, err)
				return err
			}
			log.Printf("创建知识分类: %s", category.Name)
		}
	}

	log.Println("知识库分类初始化完成")
	return nil
}
