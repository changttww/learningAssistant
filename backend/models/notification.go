package models

// Notification 通知模型
type Notification struct {
	BaseModel
	UserID       uint64 `json:"user_id" gorm:"index"` // 接收者ID
	Title        string `json:"title" gorm:"type:varchar(128)"`
	Content      string `json:"content" gorm:"type:text"`
	Type         string `json:"type" gorm:"type:varchar(32)"` // TEAM_INVITE, TEAM_APPLICATION, SYSTEM
	RelatedID    uint64 `json:"related_id"`                   // 关联ID (通常是 TeamRequestID)
	IsRead       bool   `json:"is_read" gorm:"default:false"`
	ActionStatus string `json:"action_status" gorm:"type:varchar(32);default:'PENDING'"` // PENDING, PROCESSED, NONE
}

// TeamRequest 团队请求模型
type TeamRequest struct {
	BaseModel
	TeamID    uint64 `json:"team_id" gorm:"index"`
	UserID    uint64 `json:"user_id" gorm:"index"`                             // 申请人或被邀请人
	InviterID uint64 `json:"inviter_id"`                                       // 邀请人ID (如果是主动申请则为0)
	Type      string `json:"type" gorm:"type:varchar(32)"`                     // APPLICATION (申请), INVITATION (邀请)
	Status    string `json:"status" gorm:"type:varchar(32);default:'PENDING'"` // PENDING_USER (待用户同意), PENDING_OWNER (待群主同意), APPROVED, REJECTED
}
