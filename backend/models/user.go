package models

// User 用户模型
type User struct {
	BaseModel
	Account     string `gorm:"type:varchar(64);unique;not null" json:"account"`
	Email       string `gorm:"type:varchar(128);unique" json:"email"`
	Phone       string `gorm:"type:varchar(32);unique" json:"phone"`
	DisplayName string `gorm:"type:varchar(64);not null" json:"display_name"`
	Role        int8   `gorm:"type:tinyint;default:0;comment:0=student,1=admin" json:"role"`
	AvatarURL   string `gorm:"type:varchar(256)" json:"avatar_url"`
	Bio         string `gorm:"type:varchar(256)" json:"bio"`
	Status      int8   `gorm:"type:tinyint;default:1;comment:1=active,0=disabled" json:"status"`
}

// UserProfile 用户档案模型
type UserProfile struct {
	BaseModel
	UserID          uint64 `gorm:"unique" json:"user_id"`
	TotalPoints     int    `gorm:"default:0" json:"total_points"`
	Level           int    `gorm:"default:1" json:"level"`
	TotalStudyMins  int    `gorm:"default:0" json:"total_study_mins"`
	TasksCompleted  int    `gorm:"default:0" json:"tasks_completed"`
	TasksInProgress int    `gorm:"default:0" json:"tasks_in_progress"`
}

// PointsLedger 积分账本模型
type PointsLedger struct {
	BaseModel
	UserID       uint64  `json:"user_id"`
	SourceType   int8    `json:"source_type"`
	SourceID     *uint64 `json:"source_id"`
	Delta        int     `json:"delta"`
	BalanceAfter int     `json:"balance_after"`
	Remark       string  `gorm:"type:varchar(256)" json:"remark"`
}

// LevelRule 等级规则模型
type LevelRule struct {
	BaseModel
	Level     int    `gorm:"unique" json:"level"`
	MinPoints int    `json:"min_points"`
	Badge     string `gorm:"type:varchar(64)" json:"badge"`
}