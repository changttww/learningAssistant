package models

import (
	"time"

	"gorm.io/datatypes"
)

// TaskCategory 任务分类模型
type TaskCategory struct {
	BaseModel
	Name  string `gorm:"type:varchar(64);unique;not null" json:"name"`
	Color string `gorm:"type:varchar(8)" json:"color"`
}

// Task 任务模型
type Task struct {
	BaseModel
	Title           string         `gorm:"type:varchar(128);not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	TaskType        int8           `gorm:"type:tinyint;not null;comment:1=personal,2=team" json:"task_type"`
	CategoryID      *uint64        `json:"category_id"`
	Category        *TaskCategory  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedBy       uint64         `json:"created_by"`
	OwnerUserID     *uint64        `json:"owner_user_id"`
	OwnerTeamID     *uint64        `json:"owner_team_id"`
	Status          int8           `gorm:"type:tinyint;default:0" json:"status"`
	Priority        int8           `gorm:"type:tinyint;default:0" json:"priority"`
	StartAt         *time.Time     `gorm:"precision:3" json:"start_at"`
	DueAt           *time.Time     `gorm:"precision:3" json:"due_at"`
	CompletedAt     *time.Time     `gorm:"precision:3" json:"completed_at"`
	EstimateMinutes *int           `json:"estimate_minutes"`
	EffortPoints    int            `gorm:"default:0" json:"effort_points"`
	Subtasks        datatypes.JSON `gorm:"type:json" json:"subtasks"`
	Comments        datatypes.JSON `gorm:"type:json" json:"comments"`
}

// TaskAssignee 任务分配模型
type TaskAssignee struct {
	BaseModel
	TaskID   uint64 `json:"task_id"`
	UserID   uint64 `json:"user_id"`
	IsOwner  bool   `gorm:"default:false" json:"is_owner"`
	Status   int8   `gorm:"default:0" json:"status"`
	Progress int8   `gorm:"default:0" json:"progress"`
}

// TableName 指定表名
func (TaskAssignee) TableName() string { return "task_assignees" }

// TaskStatusHistory 任务状态历史模型
type TaskStatusHistory struct {
	BaseModel
	TaskID     uint64  `json:"task_id"`
	UserID     *uint64 `json:"user_id"`
	FromStatus int8    `json:"from_status"`
	ToStatus   int8    `json:"to_status"`
	Remark     string  `gorm:"type:varchar(256)" json:"remark"`
}

// LearningRecord 学习记录模型
type LearningRecord struct {
	BaseModel
	TaskID          uint64    `json:"task_id"`
	UserID          uint64    `json:"user_id"`
	SessionStart    time.Time `gorm:"precision:3" json:"session_start"`
	SessionEnd      time.Time `gorm:"precision:3" json:"session_end"`
	DurationMinutes int       `json:"duration_minutes"`
	Note            string    `gorm:"type:varchar(256)" json:"note"`
}
