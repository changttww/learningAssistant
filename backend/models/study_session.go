package models

import "time"

// StudySession 学习会话原始记录
type StudySession struct {
	BaseModel
	UserID          uint64     `gorm:"index;not null" json:"user_id"`
	Source          string     `gorm:"type:varchar(32);not null" json:"source"`
	SourceID        *uint64    `json:"source_id"`
	StartTime       time.Time  `gorm:"precision:3;not null" json:"start_time"`
	EndTime         *time.Time `gorm:"precision:3" json:"end_time"`
	LastPingAt      time.Time  `gorm:"precision:3;not null" json:"last_ping_at"`
	DurationMinutes int        `gorm:"default:0" json:"duration_minutes"`
	Note            string     `gorm:"type:varchar(256)" json:"note"`
}

// DailyStudyStat 日学习聚合
type DailyStudyStat struct {
	BaseModel
	UserID                uint64    `gorm:"uniqueIndex:idx_user_date;not null" json:"user_id"`
	Date                  time.Time `gorm:"type:date;uniqueIndex:idx_user_date;not null" json:"date"`
	Minutes               int       `gorm:"default:0" json:"minutes"`
	SessionCount          int       `gorm:"default:0" json:"session_count"`
	NightMinutes          int       `gorm:"default:0" json:"night_minutes"`
	MorningMinutes        int       `gorm:"default:0" json:"morning_minutes"`
	FocusModeMinutes      int       `gorm:"default:0" json:"focus_mode_minutes"`
	StudyRoomMinutes      int       `gorm:"default:0" json:"study_room_minutes"`
	StudyRoomNightMinutes int       `gorm:"default:0" json:"study_room_night_minutes"`
}
