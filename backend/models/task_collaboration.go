package models

import (
	"time"

	"gorm.io/datatypes"
)

const (
	TaskCollaborationStatusActive    int8 = 1
	TaskCollaborationStatusDismissed int8 = 2

	ChatMessageTypeText          int8 = 0
	ChatMessageTypeSystem        int8 = 3
	ChatMessageTypeKnowledgeCard int8 = 4
)

// TaskCollaborationSession stores the task-specific semantics for a temporary room.
type TaskCollaborationSession struct {
	BaseModel
	TaskID        uint64         `gorm:"index;not null" json:"task_id"`
	TeamID        *uint64        `gorm:"index" json:"team_id"`
	RoomID        uint64         `gorm:"uniqueIndex;not null" json:"room_id"`
	CreatedBy     uint64         `gorm:"index;not null" json:"created_by"`
	Status        int8           `gorm:"type:tinyint;default:1;index" json:"status"`
	DismissedAt   *time.Time     `gorm:"precision:3" json:"dismissed_at"`
	Minutes       datatypes.JSON `gorm:"type:json" json:"minutes"`
	MinutesHash   string         `gorm:"type:varchar(64);index" json:"minutes_hash"`
	LastMinutesAt *time.Time     `gorm:"precision:3" json:"last_minutes_at"`
}

func (TaskCollaborationSession) TableName() string { return "task_collaboration_sessions" }

// TaskCollaborationParticipant records the intended participants derived from the task.
type TaskCollaborationParticipant struct {
	BaseModel
	SessionID uint64 `gorm:"index;not null" json:"session_id"`
	UserID    uint64 `gorm:"index;not null" json:"user_id"`
	Role      int8   `gorm:"type:tinyint;default:0" json:"role"`
}

func (TaskCollaborationParticipant) TableName() string {
	return "task_collaboration_participants"
}
