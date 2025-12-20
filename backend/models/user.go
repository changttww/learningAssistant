package models

import "time"

// User 用户模型
type User struct {
	BaseModel
	Account           string `gorm:"type:varchar(64);unique;not null" json:"account"`
	Email             string `gorm:"type:varchar(128);unique" json:"email"`
	Phone             string `gorm:"type:varchar(32);unique" json:"phone"`
	DisplayName       string `gorm:"type:varchar(64);not null" json:"display_name"`
	Role              int8   `gorm:"type:tinyint;default:0;comment:0=student,1=admin" json:"role"`
	AvatarURL         string `gorm:"type:varchar(256)" json:"avatar_url"`
	Bio               string `gorm:"type:varchar(256)" json:"bio"`
	Status            int8   `gorm:"type:tinyint;default:1;comment:1=active,0=disabled" json:"status"`
	PasswordHash      string `gorm:"type:varchar(255);not null" json:"-"`
	School            string `gorm:"type:varchar(128)" json:"school"`
	Major             string `gorm:"type:varchar(128)" json:"major"`
	Location          string `gorm:"type:varchar(128)" json:"location"`
	JoinDate          string `gorm:"type:varchar(32)" json:"join_date"`
	PreferredLanguage string `gorm:"type:varchar(32);default:'zh-CN'" json:"preferred_language"`
	PreferredTheme    string `gorm:"type:varchar(32);default:'light'" json:"preferred_theme"`
}

// UserProfile 用户档案模型
type UserProfile struct {
	BaseModel
	UserID             uint64  `gorm:"unique" json:"user_id"`
	TotalPoints        int     `gorm:"default:0" json:"total_points"`
	Level              int     `gorm:"default:1" json:"level"`
	TotalStudyMins     int     `gorm:"default:0" json:"total_study_mins"`
	TasksCompleted     int     `gorm:"default:0" json:"tasks_completed"`
	TasksInProgress    int     `gorm:"default:0" json:"tasks_in_progress"`
	TaskCompletionRate float32 `gorm:"type:decimal(5,2);default:0" json:"task_completion_rate"`
	CertificatesCount  int     `gorm:"default:0" json:"certificates_count"`
	StudyGroups        int     `gorm:"default:0" json:"study_groups"`
	RankLabel          string  `gorm:"type:varchar(64);default:'TOP 100%'" json:"rank_label"`
	StreakDays         int     `gorm:"default:0" json:"streak_days"`
	CoursesInProgress  int     `gorm:"default:0" json:"courses_in_progress"`
	NextLevelPoints    int     `gorm:"default:200" json:"next_level_points"`
}

// PointsLedger 积分账本模型
type PointsSourceType int8

const (
	PointsSourceTaskCompletion PointsSourceType = 1
	PointsSourceStudyRoom      PointsSourceType = 2
	PointsSourceDailyCheckIn   PointsSourceType = 3
)

type PointsLedger struct {
	BaseModel
	UserID       uint64           `gorm:"index;not null" json:"user_id"`
	SourceType   PointsSourceType `gorm:"type:tinyint;not null;comment:1=task_completion,2=study_room_session,3=daily_check_in" json:"source_type"`
	SourceID     *uint64          `json:"source_id"`
	Delta        int              `json:"delta"`
	BalanceAfter int              `json:"balance_after"`
	Remark       string           `gorm:"type:varchar(256)" json:"remark"`
}

// LevelRule 等级规则模型
type LevelRule struct {
	BaseModel
	Level     int    `gorm:"unique" json:"level"`
	MinPoints int    `json:"min_points"`
	Badge     string `gorm:"type:varchar(64)" json:"badge"`
}

// UserBadge 用户徽章
type UserBadge struct {
	BaseModel
	UserID uint64 `gorm:"index" json:"user_id"`
	Name   string `gorm:"type:varchar(64)" json:"name"`
}

// Achievement 成就模型
type Achievement struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"unique;not null;size:64" json:"code"`
	Name        string `gorm:"size:128" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Category    string `gorm:"size:64" json:"category"`
	Icon        string `gorm:"size:255" json:"icon"`

	// JSON 形式，示例：
	// {"type":"task_count", "value":100}
	Condition string `gorm:"type:json" json:"condition"`
}

// UserAchievement 用户成就
type UserAchievement struct {
	BaseModel
	UserID        uint64    `gorm:"index:idx_user_achievement,unique;not null" json:"user_id"`
	AchievementID uint      `gorm:"index:idx_user_achievement,unique;not null" json:"achievement_id"`
	AwardedAt     time.Time `gorm:"precision:3" json:"awarded_at"`
}

// UserAchievementProgress 成就进度
type UserAchievementProgress struct {
	BaseModel
	UserID uint64 `gorm:"uniqueIndex;not null" json:"user_id"`

	TaskCreatedCount       int `json:"task_created_count"`
	TaskCompletedCount     int `json:"task_completed_count"`
	StreakDays             int `json:"streak_days"`
	StudyRoomJoinCount     int `json:"study_room_join_count"`
	StudyRoomDurationMins  int `json:"study_room_duration_mins"`
	StudyRoomNightMins     int `json:"study_room_night_mins"`
	NightSessionMaxMins    int `json:"night_session_max_mins"`
	StudyRoomChatCount     int `json:"study_room_chat_count"`
	StudyRoomLikesGiven    int `json:"study_room_likes_given"`
	StudyRoomLikesReceived int `json:"study_room_likes_received"`
	TeamTasksCompleted     int `json:"team_tasks_completed"`
	NightStudyMins         int `json:"night_study_mins"`
	MorningStudyMins       int `json:"morning_study_mins"`
	FocusModeMins          int `json:"focus_mode_mins"`
}

// UserSkill 用户技能
type UserSkill struct {
	BaseModel
	UserID   uint64 `gorm:"index" json:"user_id"`
	Name     string `gorm:"type:varchar(64)" json:"name"`
	Category string `gorm:"type:varchar(32);default:'secondary'" json:"category"`
}

// UserSetting 用户设置
type UserSetting struct {
	BaseModel
	UserID           uint64 `gorm:"uniqueIndex" json:"user_id"`
	NotifyEmail      bool   `gorm:"default:true" json:"notify_email"`
	NotifySMS        bool   `gorm:"default:false" json:"notify_sms"`
	NotifyInApp      bool   `gorm:"default:true" json:"notify_in_app"`
	NotifySummary    bool   `gorm:"default:true" json:"notify_summary"`
	ShowEmail        bool   `gorm:"default:false" json:"show_email"`
	ShowProfile      bool   `gorm:"default:true" json:"show_profile"`
	ShowStudyData    bool   `gorm:"default:true" json:"show_study_data"`
	DailyGoalMinutes int    `gorm:"default:60" json:"daily_goal_minutes"`
	PreferredPeriod  string `gorm:"type:varchar(32);default:'evening'" json:"preferred_period"`
	FocusMode        bool   `gorm:"default:false" json:"focus_mode"`
}
