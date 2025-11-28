package seeder

import (
	"learningAssistant-backend/models"

	"gorm.io/gorm"
)

// SeedTaskCategories 初始化任务分类数据
func SeedTaskCategories(db *gorm.DB) error {
	categories := []models.TaskCategory{
		{Name: "学习", Color: "#3B82F6"},
		{Name: "工作", Color: "#10B981"},
		{Name: "运动", Color: "#F59E0B"},
		{Name: "娱乐", Color: "#8B5CF6"},
		{Name: "生活", Color: "#EF4444"},
		{Name: "阅读", Color: "#06B6D4"},
		{Name: "项目", Color: "#F97316"},
		{Name: "其他", Color: "#6B7280"},
	}

	for _, category := range categories {
		var existingCategory models.TaskCategory
		err := db.Where("name = ?", category.Name).First(&existingCategory).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&category).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
