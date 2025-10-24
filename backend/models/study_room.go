package models

import "time"

// StudyRoom 学习室模型
type StudyRoom struct {
	BaseModel
	Name        string  `gorm:"type:varchar(64);not null" json:"name"`
	OwnerUserID uint64  `json:"owner_user_id"`
	TeamID      *uint64 `json:"team_id"`
	IsPrivate   bool    `gorm:"default:false" json:"is_private"`
	Status      int8    `gorm:"default:1" json:"status"`
}

// StudyRoomMember 学习室成员模型
type StudyRoomMember struct {
	BaseModel
	RoomID   uint64    `json:"room_id"`
	UserID   uint64    `json:"user_id"`
	Role     int8      `gorm:"default:0" json:"role"`
	JoinedAt time.Time `gorm:"precision:3;autoCreateTime" json:"joined_at"`
}

// TableName 指定表名
func (StudyRoomMember) TableName() string { return "study_room_members" }

// RoomSession 房间会话模型
type RoomSession struct {
	BaseModel
	RoomID    uint64     `json:"room_id"`
	StartedBy uint64     `json:"started_by"`
	StartTime time.Time  `gorm:"precision:3" json:"start_time"`
	EndTime   *time.Time `gorm:"precision:3" json:"end_time"`
	Topic     string     `gorm:"type:varchar(128)" json:"topic"`
}

// TableName 指定表名
func (RoomSession) TableName() string { return "room_sessions" }

// ChatMessage 聊天消息模型
type ChatMessage struct {
	BaseModel
	SessionID uint64    `json:"session_id"`
	RoomID    uint64    `json:"room_id"`
	UserID    uint64    `json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	MsgType   int8      `gorm:"default:0;comment:0=text,1=image,2=file,3=system" json:"msg_type"`
	SentAt    time.Time `gorm:"precision:3;autoCreateTime" json:"sent_at"`
}

// TableName 指定表名
func (ChatMessage) TableName() string { return "chat_messages" }