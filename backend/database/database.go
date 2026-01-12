package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"learningAssistant-backend/config"
	"learningAssistant-backend/database/seeder"
	"learningAssistant-backend/models"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	var err error

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.AppConfig.Database.Username,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Database,
		config.AppConfig.Database.Charset,
	)

	// 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// 自动迁移数据库表
	if err := AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if err := seeder.SeedLevelRules(DB); err != nil {
		log.Fatal("Failed to seed level rules:", err)
	}
	if err := seeder.SeedAchievements(DB); err != nil {
		log.Fatal("Failed to seed achievements:", err)
	}
	if err := seeder.SeedTaskCategories(DB); err != nil {
		log.Println("Warning: Failed to seed task categories:", err)
	}
	if shouldSeedDemoData() {
		if err := seeder.SeedDemoData(DB); err != nil {
			log.Fatal("Failed to seed demo data:", err)
		}
	}

	// 初始化知识库分类
	if err := seeder.SeedKnowledgeCategories(DB); err != nil {
		log.Println("Warning: Failed to seed knowledge categories:", err)
	}

	// 初始化知识库演示数据
	if err := seeder.SeedKnowledgeEntries(DB); err != nil {
		log.Println("Warning: Failed to seed knowledge entries:", err)
	}
}

// AutoMigrate 自动迁移数据库表，从model中获取数据结构，并创建对应的表结构
func AutoMigrate() error {
	allModels := models.GetAllModels()

	for _, model := range allModels {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
	}

	log.Println("Database migration completed successfully")

	// 兼容旧表结构：study_notes 曾存在 (user_id, task_id, origin) 的唯一索引，
	// 会导致同一任务无法创建多篇笔记（与当前前端需求不一致）。这里自动移除该索引。
	if DB.Migrator().HasIndex(&models.StudyNote{}, "idx_user_task_origin") {
		if err := DB.Migrator().DropIndex(&models.StudyNote{}, "idx_user_task_origin"); err != nil {
			return fmt.Errorf("failed to drop index idx_user_task_origin on study_notes: %w", err)
		}
		log.Println("Dropped legacy unique index idx_user_task_origin on study_notes")
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

func shouldSeedDemoData() bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv("SEED_DEMO_DATA")))
	return value == "true" || value == "1" || value == "yes"
}
