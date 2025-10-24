package models

import "time"

// Team 团队模型
type Team struct {
	BaseModel
	Name        string `gorm:"type:varchar(64);unique;not null" json:"name"`
	Description string `gorm:"type:varchar(256)" json:"description"`
	OwnerUserID uint64 `json:"owner_user_id"`
	Visibility  int8   `gorm:"type:tinyint;default:1;comment:1=public,2=private" json:"visibility"`
}

// TeamMember 团队成员模型
type TeamMember struct {
	BaseModel
	TeamID   uint64    `json:"team_id"`
	UserID   uint64    `json:"user_id"`
	Role     int8      `gorm:"type:tinyint;default:0" json:"role"`
	JoinedAt time.Time `gorm:"precision:3;autoCreateTime" json:"joined_at"`
}

// TableName 指定表名
func (TeamMember) TableName() string { return "team_members" }