package achievement

import (
	"time"

	"learningAssistant-backend/database"
)

// UserAchievementRecord 解锁的成就信息
type UserAchievementRecord struct {
	RecordID      uint64    `json:"record_id"`
	AchievementID uint      `json:"achievement_id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Icon          string    `json:"icon"`
	AwardedAt     time.Time `json:"awarded_at"`
}

// ListUserAchievements 返回用户已获得成就
func ListUserAchievements(userID uint64) ([]UserAchievementRecord, error) {
	db := database.GetDB()
	var records []UserAchievementRecord
	err := db.Table("user_achievements AS ua").
		Select(`ua.id AS record_id,
			ua.achievement_id,
			a.code,
			a.name,
			a.description,
			a.category,
			a.icon,
			ua.awarded_at`).
		Joins("JOIN achievements AS a ON a.id = ua.achievement_id").
		Where("ua.user_id = ?", userID).
		Order("ua.awarded_at DESC").
		Scan(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}
