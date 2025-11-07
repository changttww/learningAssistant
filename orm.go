package models

import (
        "time"

        "gorm.io/gorm"
)

type BaseModel struct {
        ID        uint64         `gorm:"primaryKey;autoIncrement"`
        CreatedAt time.Time      `gorm:"precision:3;autoCreateTime"`
        UpdatedAt time.Time      `gorm:"precision:3;autoUpdateTime"`
        DeletedAt gorm.DeletedAt `gorm:"index"`
}

// ===================== Users =====================
type User struct {
        BaseModel
        Account     string `gorm:"type:varchar(64);unique;not null"`
        Email       string `gorm:"type:varchar(128);unique"`
        Phone       string `gorm:"type:varchar(32);unique"`
        DisplayName string `gorm:"type:varchar(64);not null"`
        Role        int8   `gorm:"type:tinyint;default:0;comment:0=student,1=admin"`
        AvatarURL   string `gorm:"type:varchar(256)"`
        Bio         string `gorm:"type:varchar(256)"`
        Status      int8   `gorm:"type:tinyint;default:1;comment:1=active,0=disabled"`
}

// ===================== Teams =====================
type Team struct {
        BaseModel
        Name        string `gorm:"type:varchar(64);unique;not null"`
        Description string `gorm:"type:varchar(256)"`
        OwnerUserID uint64
        Visibility  int8 `gorm:"type:tinyint;default:1;comment:1=public,2=private"`
}

// ===================== Team Members =====================
type TeamMember struct {
        BaseModel
        TeamID   uint64
        UserID   uint64
        Role     int8      `gorm:"type:tinyint;default:0"`
        JoinedAt time.Time `gorm:"precision:3;autoCreateTime"`
}

func (TeamMember) TableName() string { return "team_members" }

// ===================== Task Category =====================
type TaskCategory struct {
        BaseModel
        Name  string `gorm:"type:varchar(64);unique;not null"`
        Color string `gorm:"type:varchar(8)"`
}

// ===================== Task =====================
type Task struct {
        BaseModel
        Title           string     `gorm:"type:varchar(128);not null"`
        Description     string     `gorm:"type:text"`
        TaskType        int8       `gorm:"type:tinyint;not null;comment:1=personal,2=team"`
        CategoryID      *uint64
        CreatedBy       uint64
        OwnerUserID     *uint64
        OwnerTeamID     *uint64
        Status          int8       `gorm:"type:tinyint;default:0"`
        Priority        int8       `gorm:"type:tinyint;default:0"`
        StartAt         *time.Time `gorm:"precision:3"`
        DueAt           *time.Time `gorm:"precision:3"`
        CompletedAt     *time.Time `gorm:"precision:3"`
        EstimateMinutes *int
        EffortPoints    int `gorm:"default:0"`
}

// ===================== Task Assignees =====================
type TaskAssignee struct {
        BaseModel
        TaskID   uint64
        UserID   uint64
        IsOwner  bool  `gorm:"default:false"`
        Status   int8  `gorm:"default:0"`
        Progress int8  `gorm:"default:0"`
}

func (TaskAssignee) TableName() string { return "task_assignees" }

// ===================== Task Status History =====================
type TaskStatusHistory struct {
        BaseModel
        TaskID     uint64
        UserID     *uint64
        FromStatus int8
        ToStatus   int8
        Remark     string `gorm:"type:varchar(256)"`
}

// ===================== Learning Records =====================
type LearningRecord struct {
        BaseModel
        TaskID          uint64
        UserID          uint64
        SessionStart    time.Time `gorm:"precision:3"`
        SessionEnd      time.Time `gorm:"precision:3"`
        DurationMinutes int
        Note            string `gorm:"type:varchar(256)"`
}

// ===================== Points Ledger =====================
type PointsLedger struct {
        BaseModel
        UserID       uint64
        SourceType   int8
        SourceID     *uint64
        Delta        int
        BalanceAfter int
        Remark       string `gorm:"type:varchar(256)"`
}

// ===================== Level Rules =====================
type LevelRule struct {
        BaseModel
        Level     int    `gorm:"unique"`
        MinPoints int
        Badge     string `gorm:"type:varchar(64)"`
}

// ===================== User Profiles =====================
type UserProfile struct {
        BaseModel
        UserID          uint64 `gorm:"unique"`
        TotalPoints     int    `gorm:"default:0"`
        Level           int    `gorm:"default:1"`
        TotalStudyMins  int    `gorm:"default:0"`
        TasksCompleted  int    `gorm:"default:0"`
        TasksInProgress int    `gorm:"default:0"`
}

// ===================== Study Rooms =====================
type StudyRoom struct {
        BaseModel
        Name        string `gorm:"type:varchar(64);not null"`
        OwnerUserID uint64
        TeamID      *uint64
        IsPrivate   bool `gorm:"default:false"`
        Status      int8 `gorm:"default:1"`
}

// ===================== Study Room Members =====================
type StudyRoomMember struct {
        BaseModel
        RoomID   uint64
        UserID   uint64
        Role     int8      `gorm:"default:0"`
        JoinedAt time.Time `gorm:"precision:3;autoCreateTime"`
}

func (StudyRoomMember) TableName() string { return "study_room_members" }

// ===================== Room Sessions =====================
type RoomSession struct {
        BaseModel
        RoomID    uint64
        StartedBy uint64
        StartTime time.Time  `gorm:"precision:3"`
        EndTime   *time.Time `gorm:"precision:3"`
        Topic     string     `gorm:"type:varchar(128)"`
}

func (RoomSession) TableName() string { return "room_sessions" }

// ===================== Chat Messages =====================
type ChatMessage struct {
        BaseModel
        SessionID uint64
        RoomID    uint64
        UserID    uint64
        Content   string    `gorm:"type:text;not null"`
        MsgType   int8      `gorm:"default:0;comment:0=text,1=image,2=file,3=system"`
        SentAt    time.Time `gorm:"precision:3;autoCreateTime"`
}

func (ChatMessage) TableName() string { return "chat_messages" }
