package models

// StudyBuddy 学习伙伴关系
type StudyBuddy struct {
	BaseModel
	UserID  uint64 `gorm:"index;not null" json:"user_id"`
	BuddyID uint64 `gorm:"index;not null" json:"buddy_id"`
	Remark  string `gorm:"type:varchar(128)" json:"remark"`
	Tags    string `gorm:"type:varchar(128)" json:"tags"`
}

func (StudyBuddy) TableName() string { return "study_buddies" }
