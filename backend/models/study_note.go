package models

type StudyNote struct {
	BaseModel
	UserID  uint64  `json:"user_id"`
	TaskID  *uint64 `json:"task_id"`
	Title   string  `gorm:"type:varchar(256);not null" json:"title"`
	Content string  `gorm:"type:longtext" json:"content"`
}
