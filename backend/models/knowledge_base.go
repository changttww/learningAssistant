package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// KnowledgeBaseEntry 知识库条目 - 存储用户学习内容
type KnowledgeBaseEntry struct {
	BaseModel
	UserID       uint64         `gorm:"index" json:"user_id"`
	TeamID       *uint64        `gorm:"index" json:"team_id"` // 关联的团队ID（如果是团队共享知识）
	SourceType   int8           `gorm:"type:tinyint;comment:1=task_note,2=study_note,3=quiz_answer" json:"source_type"`
	SourceID     uint64         `json:"source_id"` // 关联的任务ID或笔记ID
	TaskID       *uint64        `json:"task_id"`   // 关联的任务
	NoteID       *uint64        `json:"note_id"`   // 关联的笔记
	Title        string         `gorm:"type:varchar(256);not null;index" json:"title"`
	Content      string         `gorm:"type:longtext" json:"content"`
	Summary      string         `gorm:"type:text" json:"summary"`               // AI生成的摘要
	Keywords     datatypes.JSON `gorm:"type:json" json:"keywords"`              // 关键词列表 []string
	Category     string         `gorm:"type:varchar(64);index" json:"category"` // 知识分类
	SubCategory  string         `gorm:"type:varchar(64)" json:"sub_category"`   // 子分类
	Level        int8           `gorm:"type:tinyint;default:0" json:"level"`    // 掌握等级 0-4
	Tags         datatypes.JSON `gorm:"type:json" json:"tags"`                  // 标签 []string
	Status       int8           `gorm:"type:tinyint;default:0" json:"status"`   // 0=draft, 1=published, 2=archived
	ViewCount    int            `gorm:"default:0" json:"view_count"`            // 查看次数
	LastReviewAt *time.Time     `gorm:"precision:3" json:"last_review_at"`      // 最后复习时间
	// 显示属性 - AI分析生成
	DisplayColor string `gorm:"type:varchar(32)" json:"display_color"` // 显示颜色(如: #3b82f6 或 渐变色)
	DisplayIcon  string `gorm:"type:varchar(64)" json:"display_icon"`  // 显示图标(iconify图标名)
	Subject      string `gorm:"type:varchar(64)" json:"subject"`       // 学科领域(数学/物理/语文等)
}

// KnowledgeVectorCache 向量缓存 - 存储embedding向量用于快速检索
type KnowledgeVectorCache struct {
	BaseModel
	EntryID     uint64 `gorm:"uniqueIndex;not null" json:"entry_id"`
	ContentHash string `gorm:"index;type:varchar(64)" json:"content_hash"` // 内容哈希，用于变化检测
	Vector      Vector `gorm:"type:json" json:"vector"`                    // 向量数据（浮点数组）
	VectorDim   int    `json:"vector_dim"`                                 // 向量维度
	VectorModel string `gorm:"type:varchar(64)" json:"vector_model"`       // 使用的模型
}

// Vector 自定义向量类型
type Vector []float32

// Value 实现 driver.Valuer 接口
func (v Vector) Value() (driver.Value, error) {
	return json.Marshal(v)
}

// Scan 实现 sql.Scanner 接口
func (v *Vector) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, &v)
}

// KnowledgeRelation 知识点关联 - 建立知识点间的关系
type KnowledgeRelation struct {
	BaseModel
	UserID        uint64  `gorm:"index" json:"user_id"`
	SourceEntryID uint64  `json:"source_entry_id"`                   // 源知识点
	TargetEntryID uint64  `json:"target_entry_id"`                   // 目标知识点
	RelationType  int8    `gorm:"type:tinyint" json:"relation_type"` // 1=prerequisite, 2=related, 3=extends, 4=conflict
	Strength      float32 `json:"strength"`                          // 关联强度 0-1
}

// UserKnowledgeStats 用户知识库统计
type UserKnowledgeStats struct {
	BaseModel
	UserID        uint64     `gorm:"uniqueIndex;not null" json:"user_id"`
	TotalEntries  int        `json:"total_entries"`  // 总知识条目数
	MasteredCount int        `json:"mastered_count"` // 掌握的知识点数
	LearningCount int        `json:"learning_count"` // 学习中的知识点数
	ToLearnCount  int        `json:"to_learn_count"` // 待学习的知识点数
	LastUpdateAt  *time.Time `gorm:"precision:3" json:"last_update_at"`
}

// KnowledgeCategory 知识分类预设
type KnowledgeCategory struct {
	BaseModel
	Name        string `gorm:"type:varchar(64);unique" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Color       string `gorm:"type:varchar(8)" json:"color"`
	Icon        string `gorm:"type:varchar(64)" json:"icon"`
}

// TaskQuizRecord 任务生成的测验记录 - 用于加入知识库
type TaskQuizRecord struct {
	BaseModel
	UserID           uint64         `gorm:"index" json:"user_id"`
	TaskID           uint64         `json:"task_id"`
	QuizContent      datatypes.JSON `gorm:"type:json" json:"quiz_content"` // 生成的测验内容
	Answers          datatypes.JSON `gorm:"type:json" json:"answers"`      // 用户答案
	Score            int            `json:"score"`                         // 得分
	Status           int8           `gorm:"type:tinyint" json:"status"`    // 0=draft, 1=submitted, 2=reviewed
	KnowledgeEntryID *uint64        `json:"knowledge_entry_id"`            // 关联到知识库的条目ID
}

// TableName 指定表名
func (KnowledgeBaseEntry) TableName() string   { return "knowledge_base_entries" }
func (KnowledgeVectorCache) TableName() string { return "knowledge_vector_caches" }
func (KnowledgeRelation) TableName() string    { return "knowledge_relations" }
func (UserKnowledgeStats) TableName() string   { return "user_knowledge_stats" }
func (KnowledgeCategory) TableName() string    { return "knowledge_categories" }
func (TaskQuizRecord) TableName() string       { return "task_quiz_records" }
