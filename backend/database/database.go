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
	if shouldSeedDemoData() {
		if err := seeder.SeedDemoData(DB); err != nil {
			log.Fatal("Failed to seed demo data:", err)
		}
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
