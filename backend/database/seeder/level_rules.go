package seeder

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"learningAssistant-backend/models"
)

type levelRuleSeed struct {
	Level     int
	MinPoints int
	Badge     string
}

var levelRuleSeeds = []levelRuleSeed{
	{Level: 1, MinPoints: 0, Badge: "Lv.1"},
	{Level: 2, MinPoints: 100, Badge: "Lv.2"},
	{Level: 3, MinPoints: 200, Badge: "Lv.3"},
	{Level: 4, MinPoints: 500, Badge: "Lv.4"},
	{Level: 5, MinPoints: 1000, Badge: "Lv.5"},
	{Level: 6, MinPoints: 5000, Badge: "Lv.6"},
}

// SeedLevelRules 将等级规则写入数据库
func SeedLevelRules(db *gorm.DB) error {
	for _, seed := range levelRuleSeeds {
		record := models.LevelRule{
			Level:     seed.Level,
			MinPoints: seed.MinPoints,
			Badge:     seed.Badge,
		}
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "level"}},
			DoUpdates: clause.AssignmentColumns([]string{"min_points", "badge"}),
		}).Create(&record).Error; err != nil {
			return fmt.Errorf("seed level rule %d: %w", seed.Level, err)
		}
	}
	return nil
}
