package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含通用字段
type BaseModel struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `gorm:"precision:3;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"precision:3;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// GetAllModels 返回所有需要迁移的模型
func GetAllModels() []interface{} {
    return []interface{}{
        &User{},
        &UserProfile{},
        &UserBadge{},
        &Achievement{},
        &UserAchievement{},
        &UserAchievementProgress{},
        &UserSkill{},
        &UserSetting{},
        &PointsLedger{},
        &LevelRule{},
        &Team{},
        &TeamMember{},
        &TaskCategory{},
        &Task{},
        &TaskAssignee{},
        &TaskStatusHistory{},
        &LearningRecord{},
        &StudyRoom{},
        &StudyRoomMember{},
        &RoomSession{},
        &ChatMessage{},
        &StudyNote{},
    }
}
