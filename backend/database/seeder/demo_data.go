package seeder

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"learningAssistant-backend/models"
)

type demoSeedContext struct {
	Users                map[string]*models.User
	Teams                map[string]*models.Team
	TaskCategories       map[string]*models.TaskCategory
	Tasks                map[string]*models.Task
	StudyRooms           map[string]*models.StudyRoom
	Sessions             map[string]*models.RoomSession
	Achievements         map[string]models.Achievement
	userAchievementCodes map[string][]string
}

// SeedDemoData 填充演示数据，便于前端联调与示例展示
func SeedDemoData(db *gorm.DB) error {
	ctx := &demoSeedContext{
		Users:                map[string]*models.User{},
		Teams:                map[string]*models.Team{},
		TaskCategories:       map[string]*models.TaskCategory{},
		Tasks:                map[string]*models.Task{},
		StudyRooms:           map[string]*models.StudyRoom{},
		Sessions:             map[string]*models.RoomSession{},
		Achievements:         map[string]models.Achievement{},
		userAchievementCodes: map[string][]string{},
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := seedDemoUsers(tx, ctx); err != nil {
			return err
		}
		if err := seedDemoTeams(tx, ctx); err != nil {
			return err
		}
		if err := seedDemoTaskCategories(tx, ctx); err != nil {
			return err
		}
		if err := seedDemoTasks(tx, ctx); err != nil {
			return err
		}
		if err := seedDemoLearningRecords(tx, ctx); err != nil {
			return err
		}
		if err := seedDemoStudyRooms(tx, ctx); err != nil {
			return err
		}
		if err := seedDemoPointsLedger(tx, ctx); err != nil {
			return err
		}
		if err := seedDemoUserAchievements(tx, ctx); err != nil {
			return err
		}
		return nil
	})
}

type userSeed struct {
	Key              string
	Password         string
	User             models.User
	Profile          models.UserProfile
	Settings         models.UserSetting
	Skills           []models.UserSkill
	Badges           []string
	Progress         models.UserAchievementProgress
	AchievementCodes []string
}

func seedDemoUsers(tx *gorm.DB, ctx *demoSeedContext) error {
	seeds := []userSeed{
		{
			Key:      "mentor",
			Password: "Passw0rd!",
			User: models.User{
				Account:           "mentor_lee",
				Email:             "mentor.lee@example.com",
				Phone:             "13000000001",
				DisplayName:       "李老师",
				Role:              1,
				AvatarURL:         "https://placehold.co/96x96/2563eb/ffffff.png?text=L",
				Bio:               "时间管理教练，主张把每天的黄金时间留给最重要的事。",
				Status:            1,
				School:            "华南理工大学",
				Major:             "教育学",
				Location:          "广州",
				JoinDate:          "2023-09",
				PreferredLanguage: "zh-CN",
				PreferredTheme:    "light",
			},
			Profile: models.UserProfile{
				TotalPoints:        1280,
				Level:              5,
				TotalStudyMins:     5420,
				TasksCompleted:     42,
				TasksInProgress:    5,
				TaskCompletionRate: 89.0,
				CertificatesCount:  3,
				StudyGroups:        2,
				RankLabel:          "TOP 12%",
				StreakDays:         18,
				CoursesInProgress:  2,
				NextLevelPoints:    5000,
			},
			Settings: models.UserSetting{
				NotifyEmail:      true,
				NotifySMS:        false,
				NotifyInApp:      true,
				NotifySummary:    true,
				ShowEmail:        false,
				ShowProfile:      true,
				ShowStudyData:    true,
				DailyGoalMinutes: 120,
				PreferredPeriod:  "morning",
				FocusMode:        true,
			},
			Skills: []models.UserSkill{
				{Name: "Golang", Category: "primary"},
				{Name: "Productivity", Category: "secondary"},
				{Name: "数据分析", Category: "secondary"},
			},
			Badges: []string{"早起之星", "学习引导员"},
			Progress: models.UserAchievementProgress{
				TaskCreatedCount:       55,
				TaskCompletedCount:     42,
				StreakDays:             18,
				StudyRoomJoinCount:     34,
				StudyRoomDurationMins:  3600,
				StudyRoomNightMins:     260,
				NightSessionMaxMins:    120,
				StudyRoomChatCount:     42,
				StudyRoomLikesGiven:    30,
				StudyRoomLikesReceived: 28,
				TeamTasksCompleted:     14,
			},
			AchievementCodes: []string{"task_diligent_20", "streak_15", "studyroom_active_30", "studyroom_focus_20h", "team_contributor_10"},
		},
		{
			Key:      "focus",
			Password: "Focus123!",
			User: models.User{
				Account:           "focus_tan",
				Email:             "focus.tan@example.com",
				Phone:             "13000000002",
				DisplayName:       "Tan",
				Role:              0,
				AvatarURL:         "https://placehold.co/96x96/f59e0b/ffffff.png?text=T",
				Bio:               "准备秋招，主攻数据分析与模型评估。",
				Status:            1,
				School:            "上海交通大学",
				Major:             "信息管理",
				Location:          "上海",
				JoinDate:          "2024-02",
				PreferredLanguage: "zh-CN",
				PreferredTheme:    "dark",
			},
			Profile: models.UserProfile{
				TotalPoints:        620,
				Level:              4,
				TotalStudyMins:     2400,
				TasksCompleted:     18,
				TasksInProgress:    6,
				TaskCompletionRate: 75.0,
				CertificatesCount:  1,
				StudyGroups:        1,
				RankLabel:          "TOP 35%",
				StreakDays:         7,
				CoursesInProgress:  1,
				NextLevelPoints:    1000,
			},
			Settings: models.UserSetting{
				NotifyEmail:      true,
				NotifySMS:        false,
				NotifyInApp:      true,
				NotifySummary:    true,
				ShowEmail:        false,
				ShowProfile:      true,
				ShowStudyData:    true,
				DailyGoalMinutes: 90,
				PreferredPeriod:  "evening",
				FocusMode:        false,
			},
			Skills: []models.UserSkill{
				{Name: "Python", Category: "primary"},
				{Name: "SQL", Category: "primary"},
				{Name: "写作", Category: "secondary"},
			},
			Badges: []string{"效率提升"},
			Progress: models.UserAchievementProgress{
				TaskCreatedCount:       26,
				TaskCompletedCount:     18,
				StreakDays:             7,
				StudyRoomJoinCount:     12,
				StudyRoomDurationMins:  1500,
				StudyRoomNightMins:     300,
				NightSessionMaxMins:    90,
				StudyRoomChatCount:     18,
				StudyRoomLikesGiven:    12,
				StudyRoomLikesReceived: 15,
				TeamTasksCompleted:     6,
			},
			AchievementCodes: []string{"task_starter_5", "streak_7", "studyroom_focus_5h", "studyroom_first_join", "team_joiner"},
		},
		{
			Key:      "rookie",
			Password: "Hello123!",
			User: models.User{
				Account:           "rookie_wu",
				Email:             "rookie.wu@example.com",
				Phone:             "13000000003",
				DisplayName:       "小吴",
				Role:              0,
				AvatarURL:         "https://placehold.co/96x96/10b981/ffffff.png?text=W",
				Bio:               "大三学生，刚开始系统化学习，目标保持稳定节奏。",
				Status:            1,
				School:            "中山大学",
				Major:             "统计学",
				Location:          "深圳",
				JoinDate:          "2024-05",
				PreferredLanguage: "zh-CN",
				PreferredTheme:    "light",
			},
			Profile: models.UserProfile{
				TotalPoints:        120,
				Level:              2,
				TotalStudyMins:     560,
				TasksCompleted:     4,
				TasksInProgress:    3,
				TaskCompletionRate: 57.0,
				CertificatesCount:  0,
				StudyGroups:        0,
				RankLabel:          "TOP 78%",
				StreakDays:         3,
				CoursesInProgress:  1,
				NextLevelPoints:    200,
			},
			Settings: models.UserSetting{
				NotifyEmail:      true,
				NotifySMS:        false,
				NotifyInApp:      true,
				NotifySummary:    true,
				ShowEmail:        false,
				ShowProfile:      true,
				ShowStudyData:    true,
				DailyGoalMinutes: 45,
				PreferredPeriod:  "afternoon",
				FocusMode:        false,
			},
			Skills: []models.UserSkill{
				{Name: "线性代数", Category: "secondary"},
				{Name: "英语听力", Category: "secondary"},
			},
			Badges: []string{"保持专注"},
			Progress: models.UserAchievementProgress{
				TaskCreatedCount:       10,
				TaskCompletedCount:     4,
				StreakDays:             3,
				StudyRoomJoinCount:     4,
				StudyRoomDurationMins:  320,
				StudyRoomNightMins:     45,
				NightSessionMaxMins:    45,
				StudyRoomChatCount:     6,
				StudyRoomLikesGiven:    4,
				StudyRoomLikesReceived: 2,
				TeamTasksCompleted:     1,
			},
			AchievementCodes: []string{},
		},
	}

	for _, seed := range seeds {
		user := seed.User
		user.PasswordHash = hashPassword(seed.Password)
		user.Account = strings.TrimSpace(user.Account)

		var existing models.User
		err := tx.Where("account = ?", user.Account).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&user).Error; err != nil {
				return fmt.Errorf("create user %s: %w", user.Account, err)
			}
			existing = user
		} else if err != nil {
			return err
		} else {
			updateMap := map[string]interface{}{
				"email":              user.Email,
				"phone":              user.Phone,
				"display_name":       user.DisplayName,
				"role":               user.Role,
				"avatar_url":         user.AvatarURL,
				"bio":                user.Bio,
				"status":             user.Status,
				"password_hash":      user.PasswordHash,
				"school":             user.School,
				"major":              user.Major,
				"location":           user.Location,
				"join_date":          user.JoinDate,
				"preferred_language": user.PreferredLanguage,
				"preferred_theme":    user.PreferredTheme,
			}
			if err := tx.Model(&existing).Updates(updateMap).Error; err != nil {
				return fmt.Errorf("update user %s: %w", user.Account, err)
			}
		}

		if err := upsertUserProfile(tx, existing.ID, seed.Profile); err != nil {
			return err
		}
		if err := upsertUserSettings(tx, existing.ID, seed.Settings); err != nil {
			return err
		}
		if err := upsertUserProgress(tx, existing.ID, seed.Progress); err != nil {
			return err
		}
		if err := ensureUserSkills(tx, existing.ID, seed.Skills); err != nil {
			return err
		}
		if err := ensureUserBadges(tx, existing.ID, seed.Badges); err != nil {
			return err
		}

		ctx.Users[seed.Key] = &existing
		ctx.userAchievementCodes[seed.Key] = seed.AchievementCodes
	}

	return nil
}

func upsertUserProfile(tx *gorm.DB, userID uint64, profile models.UserProfile) error {
	profile.UserID = userID
	var existing models.UserProfile
	err := tx.Where("user_id = ?", userID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(&profile).Error; err != nil {
			return fmt.Errorf("create profile for user %d: %w", userID, err)
		}
		return nil
	} else if err != nil {
		return err
	}

	updateMap := map[string]interface{}{
		"total_points":         profile.TotalPoints,
		"level":                profile.Level,
		"total_study_mins":     profile.TotalStudyMins,
		"tasks_completed":      profile.TasksCompleted,
		"tasks_in_progress":    profile.TasksInProgress,
		"task_completion_rate": profile.TaskCompletionRate,
		"certificates_count":   profile.CertificatesCount,
		"study_groups":         profile.StudyGroups,
		"rank_label":           profile.RankLabel,
		"streak_days":          profile.StreakDays,
		"courses_in_progress":  profile.CoursesInProgress,
		"next_level_points":    profile.NextLevelPoints,
	}
	return tx.Model(&existing).Updates(updateMap).Error
}

func upsertUserSettings(tx *gorm.DB, userID uint64, settings models.UserSetting) error {
	settings.UserID = userID
	var existing models.UserSetting
	err := tx.Where("user_id = ?", userID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(&settings).Error; err != nil {
			return fmt.Errorf("create settings for user %d: %w", userID, err)
		}
		return nil
	} else if err != nil {
		return err
	}

	updateMap := map[string]interface{}{
		"notify_email":       settings.NotifyEmail,
		"notify_sms":         settings.NotifySMS,
		"notify_in_app":      settings.NotifyInApp,
		"notify_summary":     settings.NotifySummary,
		"show_email":         settings.ShowEmail,
		"show_profile":       settings.ShowProfile,
		"show_study_data":    settings.ShowStudyData,
		"daily_goal_minutes": settings.DailyGoalMinutes,
		"preferred_period":   settings.PreferredPeriod,
		"focus_mode":         settings.FocusMode,
	}
	return tx.Model(&existing).Updates(updateMap).Error
}

func upsertUserProgress(tx *gorm.DB, userID uint64, progress models.UserAchievementProgress) error {
	progress.UserID = userID
	var existing models.UserAchievementProgress
	err := tx.Where("user_id = ?", userID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(&progress).Error; err != nil {
			return fmt.Errorf("create achievement progress for user %d: %w", userID, err)
		}
		return nil
	} else if err != nil {
		return err
	}

	updateMap := map[string]interface{}{
		"task_created_count":        progress.TaskCreatedCount,
		"task_completed_count":      progress.TaskCompletedCount,
		"streak_days":               progress.StreakDays,
		"study_room_join_count":     progress.StudyRoomJoinCount,
		"study_room_duration_mins":  progress.StudyRoomDurationMins,
		"study_room_night_mins":     progress.StudyRoomNightMins,
		"night_session_max_mins":    progress.NightSessionMaxMins,
		"study_room_chat_count":     progress.StudyRoomChatCount,
		"study_room_likes_given":    progress.StudyRoomLikesGiven,
		"study_room_likes_received": progress.StudyRoomLikesReceived,
		"team_tasks_completed":      progress.TeamTasksCompleted,
	}
	return tx.Model(&existing).Updates(updateMap).Error
}

func ensureUserSkills(tx *gorm.DB, userID uint64, skills []models.UserSkill) error {
	for _, skill := range skills {
		skill.UserID = userID
		var existing models.UserSkill
		err := tx.Where("user_id = ? AND name = ?", userID, skill.Name).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&skill).Error; err != nil {
				return fmt.Errorf("add skill %s for user %d: %w", skill.Name, userID, err)
			}
		} else if err != nil {
			return err
		} else {
			if err := tx.Model(&existing).Update("category", skill.Category).Error; err != nil {
				return fmt.Errorf("update skill %s for user %d: %w", skill.Name, userID, err)
			}
		}
	}
	return nil
}

func ensureUserBadges(tx *gorm.DB, userID uint64, badges []string) error {
	for _, badge := range badges {
		var existing models.UserBadge
		err := tx.Where("user_id = ? AND name = ?", userID, badge).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			record := models.UserBadge{
				UserID: userID,
				Name:   badge,
			}
			if err := tx.Create(&record).Error; err != nil {
				return fmt.Errorf("add badge %s for user %d: %w", badge, userID, err)
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

type teamSeed struct {
	Key            string
	Name           string
	Description    string
	OwnerUserKey   string
	Visibility     int8
	Members        []teamMemberSeed
	MemberJoinDays int
}

type teamMemberSeed struct {
	UserKey string
	Role    int8
}

func seedDemoTeams(tx *gorm.DB, ctx *demoSeedContext) error {
	seeds := []teamSeed{
		{
			Key:          "sunrise",
			Name:         "晨光学习团",
			Description:  "每天早上同步当天目标，互相监督打卡。",
			OwnerUserKey: "mentor",
			Visibility:   1,
			Members: []teamMemberSeed{
				{UserKey: "mentor", Role: 1},
				{UserKey: "focus", Role: 1},
				{UserKey: "rookie", Role: 0},
			},
			MemberJoinDays: 10,
		},
	}

	for _, seed := range seeds {
		owner, ok := ctx.Users[seed.OwnerUserKey]
		if !ok {
			return fmt.Errorf("team owner %s not found", seed.OwnerUserKey)
		}

		var team models.Team
		err := tx.Where("name = ?", seed.Name).First(&team).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			team = models.Team{
				Name:        seed.Name,
				Description: seed.Description,
				OwnerUserID: owner.ID,
				Visibility:  seed.Visibility,
			}
			if err := tx.Create(&team).Error; err != nil {
				return fmt.Errorf("create team %s: %w", seed.Name, err)
			}
		} else if err != nil {
			return err
		} else {
			if err := tx.Model(&team).Updates(map[string]interface{}{
				"description":   seed.Description,
				"owner_user_id": owner.ID,
				"visibility":    seed.Visibility,
			}).Error; err != nil {
				return fmt.Errorf("update team %s: %w", seed.Name, err)
			}
		}

		joinedAt := time.Now().AddDate(0, 0, -seed.MemberJoinDays)
		for _, memberSeed := range seed.Members {
			user, ok := ctx.Users[memberSeed.UserKey]
			if !ok {
				return fmt.Errorf("team member %s not found", memberSeed.UserKey)
			}
			var member models.TeamMember
			err := tx.Where("team_id = ? AND user_id = ?", team.ID, user.ID).First(&member).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				record := models.TeamMember{
					TeamID:   team.ID,
					UserID:   user.ID,
					Role:     memberSeed.Role,
					JoinedAt: joinedAt,
				}
				if err := tx.Create(&record).Error; err != nil {
					return fmt.Errorf("add member %s to team %s: %w", memberSeed.UserKey, seed.Name, err)
				}
			} else if err != nil {
				return err
			} else {
				if err := tx.Model(&member).Update("role", memberSeed.Role).Error; err != nil {
					return fmt.Errorf("update member %s role: %w", memberSeed.UserKey, err)
				}
			}
		}
		ctx.Teams[seed.Key] = &team
	}

	return nil
}

type taskCategorySeed struct {
	Key   string
	Name  string
	Color string
}

func seedDemoTaskCategories(tx *gorm.DB, ctx *demoSeedContext) error {
	seeds := []taskCategorySeed{
		{Key: "algo", Name: "算法刷题", Color: "#10b981"},
		{Key: "paper", Name: "论文阅读", Color: "#6366f1"},
		{Key: "english", Name: "英语训练", Color: "#f59e0b"},
		{Key: "team", Name: "团队项目", Color: "#06b6d4"},
	}

	for _, seed := range seeds {
		var category models.TaskCategory
		err := tx.Where("name = ?", seed.Name).First(&category).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			category = models.TaskCategory{
				Name:  seed.Name,
				Color: seed.Color,
			}
			if err := tx.Create(&category).Error; err != nil {
				return fmt.Errorf("create task category %s: %w", seed.Name, err)
			}
		} else if err != nil {
			return err
		} else {
			if err := tx.Model(&category).Update("color", seed.Color).Error; err != nil {
				return fmt.Errorf("update task category %s: %w", seed.Name, err)
			}
		}
		ctx.TaskCategories[seed.Key] = &category
	}
	return nil
}

type taskSeed struct {
	Key             string
	Title           string
	Description     string
	TaskType        int8
	CategoryKey     string
	CreatedBy       string
	OwnerUser       *string
	OwnerTeam       *string
	Status          int8
	Priority        int8
	StartInDays     int
	DueInDays       int
	CompleteInDays  *int
	EstimateMinutes *int
	EffortPoints    int
	Assignees       []taskAssigneeSeed
	Histories       []taskHistorySeed
}

type taskAssigneeSeed struct {
	UserKey  string
	IsOwner  bool
	Status   int8
	Progress int8
}

type taskHistorySeed struct {
	FromStatus int8
	ToStatus   int8
	Remark     string
	DayOffset  int
}

func seedDemoTasks(tx *gorm.DB, ctx *demoSeedContext) error {
	now := time.Now()
	mentor := "mentor"
	seeds := []taskSeed{
		{
			Key:             "algo-notes",
			Title:           "刷算法题：二叉树专题",
			Description:     "练习二叉树前中后序遍历与层序遍历，总结模板。",
			TaskType:        1,
			CategoryKey:     "algo",
			CreatedBy:       mentor,
			OwnerUser:       &mentor,
			Status:          2,
			Priority:        2,
			StartInDays:     -12,
			DueInDays:       -9,
			CompleteInDays:  ptrInt(-9),
			EstimateMinutes: ptrInt(180),
			EffortPoints:    120,
			Assignees: []taskAssigneeSeed{
				{UserKey: mentor, IsOwner: true, Status: 2, Progress: 100},
			},
			Histories: []taskHistorySeed{
				{FromStatus: 0, ToStatus: 1, Remark: "拆分子任务，开始练习", DayOffset: -11},
				{FromStatus: 1, ToStatus: 2, Remark: "完成专题刷题", DayOffset: -9},
			},
		},
		{
			Key:             "paper-review",
			Title:           "阅读 Transformers 论文并做总结",
			Description:     "完成阅读笔记，梳理 Attention 公式与实验对比。",
			TaskType:        1,
			CategoryKey:     "paper",
			CreatedBy:       mentor,
			OwnerUser:       &mentor,
			Status:          1,
			Priority:        2,
			StartInDays:     -3,
			DueInDays:       2,
			EstimateMinutes: ptrInt(200),
			EffortPoints:    160,
			Assignees: []taskAssigneeSeed{
				{UserKey: mentor, IsOwner: true, Status: 1, Progress: 40},
			},
			Histories: []taskHistorySeed{
				{FromStatus: 0, ToStatus: 1, Remark: "分章节阅读", DayOffset: -3},
			},
		},
		{
			Key:             "english-listening",
			Title:           "英语听力练习 30 分钟",
			Description:     "精听 VOA，复盘错题，整理生词。",
			TaskType:        1,
			CategoryKey:     "english",
			CreatedBy:       "focus",
			OwnerUser:       strPtr("focus"),
			Status:          1,
			Priority:        1,
			StartInDays:     -1,
			DueInDays:       1,
			EstimateMinutes: ptrInt(40),
			EffortPoints:    50,
			Assignees: []taskAssigneeSeed{
				{UserKey: "focus", IsOwner: true, Status: 1, Progress: 30},
			},
			Histories: []taskHistorySeed{
				{FromStatus: 0, ToStatus: 1, Remark: "设定材料与目标", DayOffset: -1},
			},
		},
		{
			Key:             "team-kanban",
			Title:           "团队：构建学习进度看板",
			Description:     "设计看板字段，完成接口联调，输出演示版。",
			TaskType:        2,
			CategoryKey:     "team",
			CreatedBy:       mentor,
			OwnerUser:       &mentor,
			OwnerTeam:       strPtr("sunrise"),
			Status:          1,
			Priority:        2,
			StartInDays:     -5,
			DueInDays:       4,
			EstimateMinutes: ptrInt(240),
			EffortPoints:    220,
			Assignees: []taskAssigneeSeed{
				{UserKey: mentor, IsOwner: true, Status: 1, Progress: 40},
				{UserKey: "focus", IsOwner: false, Status: 1, Progress: 30},
			},
			Histories: []taskHistorySeed{
				{FromStatus: 0, ToStatus: 1, Remark: "立项并确认字段", DayOffset: -5},
			},
		},
		{
			Key:            "team-retro",
			Title:          "团队：周度复盘与分工",
			Description:    "汇总学习数据，确定下周优先事项。",
			TaskType:       2,
			CategoryKey:    "team",
			CreatedBy:      "focus",
			OwnerUser:      strPtr("focus"),
			OwnerTeam:      strPtr("sunrise"),
			Status:         2,
			Priority:       1,
			StartInDays:    -9,
			DueInDays:      -7,
			CompleteInDays: ptrInt(-6),
			EffortPoints:   140,
			Assignees: []taskAssigneeSeed{
				{UserKey: "focus", IsOwner: true, Status: 2, Progress: 100},
				{UserKey: "rookie", IsOwner: false, Status: 2, Progress: 100},
			},
			Histories: []taskHistorySeed{
				{FromStatus: 0, ToStatus: 1, Remark: "收集成员反馈", DayOffset: -8},
				{FromStatus: 1, ToStatus: 2, Remark: "输出复盘结论", DayOffset: -6},
			},
		},
		{
			Key:             "vocab",
			Title:           "背诵考研英语单词第5单元",
			Description:     "使用艾宾浩斯曲线复习，记录遗忘点。",
			TaskType:        1,
			CategoryKey:     "english",
			CreatedBy:       "rookie",
			OwnerUser:       strPtr("rookie"),
			Status:          0,
			Priority:        1,
			StartInDays:     0,
			DueInDays:       1,
			EstimateMinutes: ptrInt(50),
			EffortPoints:    60,
			Assignees: []taskAssigneeSeed{
				{UserKey: "rookie", IsOwner: true, Status: 0, Progress: 0},
			},
		},
	}

	for _, seed := range seeds {
		creator, ok := ctx.Users[seed.CreatedBy]
		if !ok {
			return fmt.Errorf("creator %s not found for task %s", seed.CreatedBy, seed.Title)
		}

		category, ok := ctx.TaskCategories[seed.CategoryKey]
		if !ok {
			return fmt.Errorf("category %s not found for task %s", seed.CategoryKey, seed.Title)
		}

		var ownerUserID *uint64
		if seed.OwnerUser != nil {
			user, ok := ctx.Users[*seed.OwnerUser]
			if !ok {
				return fmt.Errorf("owner user %s not found for task %s", *seed.OwnerUser, seed.Title)
			}
			ownerUserID = &user.ID
		}

		var ownerTeamID *uint64
		if seed.OwnerTeam != nil {
			team, ok := ctx.Teams[*seed.OwnerTeam]
			if !ok {
				return fmt.Errorf("owner team %s not found for task %s", *seed.OwnerTeam, seed.Title)
			}
			ownerTeamID = &team.ID
		}

		startAt := now.AddDate(0, 0, seed.StartInDays)
		var dueAt *time.Time
		if seed.DueInDays != 0 {
			d := now.AddDate(0, 0, seed.DueInDays)
			dueAt = &d
		}
		var completedAt *time.Time
		if seed.CompleteInDays != nil {
			c := now.AddDate(0, 0, *seed.CompleteInDays)
			completedAt = &c
		}

		taskData := models.Task{
			Title:           seed.Title,
			Description:     seed.Description,
			TaskType:        seed.TaskType,
			CategoryID:      &category.ID,
			CreatedBy:       creator.ID,
			OwnerUserID:     ownerUserID,
			OwnerTeamID:     ownerTeamID,
			Status:          seed.Status,
			Priority:        seed.Priority,
			StartAt:         &startAt,
			DueAt:           dueAt,
			CompletedAt:     completedAt,
			EstimateMinutes: seed.EstimateMinutes,
			EffortPoints:    seed.EffortPoints,
		}

		var task models.Task
		err := tx.Where("title = ? AND created_by = ?", seed.Title, creator.ID).First(&task).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			task = taskData
			if err := tx.Create(&task).Error; err != nil {
				return fmt.Errorf("create task %s: %w", seed.Title, err)
			}
		} else if err != nil {
			return err
		} else {
			if err := tx.Model(&task).Updates(taskData).Error; err != nil {
				return fmt.Errorf("update task %s: %w", seed.Title, err)
			}
		}

		for _, assigneeSeed := range seed.Assignees {
			user, ok := ctx.Users[assigneeSeed.UserKey]
			if !ok {
				return fmt.Errorf("assignee %s not found for task %s", assigneeSeed.UserKey, seed.Title)
			}
			var assignee models.TaskAssignee
			err := tx.Where("task_id = ? AND user_id = ?", task.ID, user.ID).First(&assignee).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				record := models.TaskAssignee{
					TaskID:   task.ID,
					UserID:   user.ID,
					IsOwner:  assigneeSeed.IsOwner,
					Status:   assigneeSeed.Status,
					Progress: assigneeSeed.Progress,
				}
				if err := tx.Create(&record).Error; err != nil {
					return fmt.Errorf("add assignee %s for task %s: %w", assigneeSeed.UserKey, seed.Title, err)
				}
			} else if err != nil {
				return err
			} else {
				if err := tx.Model(&assignee).Updates(map[string]interface{}{
					"is_owner": assigneeSeed.IsOwner,
					"status":   assigneeSeed.Status,
					"progress": assigneeSeed.Progress,
				}).Error; err != nil {
					return fmt.Errorf("update assignee %s for task %s: %w", assigneeSeed.UserKey, seed.Title, err)
				}
			}
		}

		for _, history := range seed.Histories {
			createdAt := now.AddDate(0, 0, history.DayOffset)
			var record models.TaskStatusHistory
			err := tx.Where("task_id = ? AND from_status = ? AND to_status = ? AND remark = ?", task.ID, history.FromStatus, history.ToStatus, history.Remark).
				First(&record).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				entry := models.TaskStatusHistory{
					TaskID:     task.ID,
					FromStatus: history.FromStatus,
					ToStatus:   history.ToStatus,
					Remark:     history.Remark,
					BaseModel: models.BaseModel{
						CreatedAt: createdAt,
						UpdatedAt: createdAt,
					},
				}
				if err := tx.Create(&entry).Error; err != nil {
					return fmt.Errorf("add history for task %s: %w", seed.Title, err)
				}
			} else if err != nil {
				return err
			}
		}

		ctx.Tasks[seed.Key] = &task
	}

	return nil
}

type learningRecordSeed struct {
	TaskKey       string
	UserKey       string
	StartHoursAgo int
	DurationMins  int
	Note          string
}

func seedDemoLearningRecords(tx *gorm.DB, ctx *demoSeedContext) error {
	now := time.Now()
	seeds := []learningRecordSeed{
		// mentor 连续 7 天内的学习记录，覆盖算法、论文和团队任务
		{TaskKey: "algo-notes", UserKey: "mentor", StartHoursAgo: 24*10 + 6, DurationMins: 80, Note: "完成后序遍历总结"},
		{TaskKey: "paper-review", UserKey: "mentor", StartHoursAgo: 24*6 + 9, DurationMins: 90, Note: "阅读并记录注意力公式"},
		{TaskKey: "team-kanban", UserKey: "mentor", StartHoursAgo: 24*5 + 8, DurationMins: 75, Note: "定义接口字段与权限"},
		{TaskKey: "algo-notes", UserKey: "mentor", StartHoursAgo: 24*4 + 6, DurationMins: 95, Note: "刷题复盘并整理模板"},
		{TaskKey: "paper-review", UserKey: "mentor", StartHoursAgo: 24*3 + 5, DurationMins: 85, Note: "补充实验对比记录"},
		{TaskKey: "team-kanban", UserKey: "mentor", StartHoursAgo: 24*2 + 4, DurationMins: 80, Note: "接口联调与测试"},
		{TaskKey: "paper-review", UserKey: "mentor", StartHoursAgo: 24 + 3, DurationMins: 70, Note: "完善摘要与引言"},
		{TaskKey: "algo-notes", UserKey: "mentor", StartHoursAgo: 12, DurationMins: 60, Note: "早间错题回顾与笔记"},

		// focus 在最近 7 天内的听力与团队协作记录
		{TaskKey: "team-kanban", UserKey: "focus", StartHoursAgo: 24*4 + 9, DurationMins: 95, Note: "组件拆分与样式联调"},
		{TaskKey: "english-listening", UserKey: "focus", StartHoursAgo: 24*3 + 2, DurationMins: 45, Note: "精听训练并整理生词"},
		{TaskKey: "team-retro", UserKey: "focus", StartHoursAgo: 24*2 + 10, DurationMins: 60, Note: "整理会议结论与分工"},
		{TaskKey: "english-listening", UserKey: "focus", StartHoursAgo: 20, DurationMins: 35, Note: "晚间听力复盘"},
		{TaskKey: "english-listening", UserKey: "focus", StartHoursAgo: 4, DurationMins: 30, Note: "睡前泛听放松"},

		// rookie 以背单词为主的练习记录
		{TaskKey: "team-retro", UserKey: "rookie", StartHoursAgo: 24*5 + 11, DurationMins: 50, Note: "整理个人复盘输出"},
		{TaskKey: "vocab", UserKey: "rookie", StartHoursAgo: 24*3 + 1, DurationMins: 30, Note: "第5单元首轮记忆"},
		{TaskKey: "vocab", UserKey: "rookie", StartHoursAgo: 24 + 2, DurationMins: 25, Note: "艾宾浩斯复习清单"},
		{TaskKey: "vocab", UserKey: "rookie", StartHoursAgo: 6, DurationMins: 20, Note: "快速巩固易错词"},
	}

	for _, seed := range seeds {
		task, ok := ctx.Tasks[seed.TaskKey]
		if !ok {
			return fmt.Errorf("task %s not found for learning record", seed.TaskKey)
		}
		user, ok := ctx.Users[seed.UserKey]
		if !ok {
			return fmt.Errorf("user %s not found for learning record", seed.UserKey)
		}

		start := now.Add(-time.Duration(seed.StartHoursAgo) * time.Hour)
		end := start.Add(time.Duration(seed.DurationMins) * time.Minute)

		var record models.LearningRecord
		err := tx.Where("task_id = ? AND user_id = ? AND session_start = ?", task.ID, user.ID, start).
			First(&record).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			entry := models.LearningRecord{
				TaskID:          task.ID,
				UserID:          user.ID,
				SessionStart:    start,
				SessionEnd:      end,
				DurationMinutes: seed.DurationMins,
				Note:            seed.Note,
			}
			if err := tx.Create(&entry).Error; err != nil {
				return fmt.Errorf("create learning record for task %s: %w", seed.TaskKey, err)
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

type studyRoomSeed struct {
	Key                 string
	Name                string
	Description         string
	OwnerUserKey        string
	TeamKey             *string
	Tags                []string
	MaxMembers          int
	IsPrivate           bool
	AccessCode          string
	Status              int8
	FocusMinutesToday   int
	LastSessionInMinute *int
	Members             []roomMemberSeed
	Sessions            []roomSessionSeed
}

type roomMemberSeed struct {
	UserKey string
	Role    int8
}

type roomSessionSeed struct {
	Key             string
	StartedBy       string
	StartMinutesAgo int
	DurationMinutes int
	Topic           string
	Messages        []chatMessageSeed
}

type chatMessageSeed struct {
	UserKey   string
	OffsetMin int
	Content   string
	MsgType   int8
}

func seedDemoStudyRooms(tx *gorm.DB, ctx *demoSeedContext) error {
	now := time.Now()
	last60 := 60
	last720 := 720
	seeds := []studyRoomSeed{
		{
			Key:                 "morning",
			Name:                "晨间深度工作室",
			Description:         "早起 90 分钟深度工作冲刺区",
			OwnerUserKey:        "mentor",
			TeamKey:             strPtr("sunrise"),
			Tags:                []string{"早起", "深度工作"},
			MaxMembers:          12,
			IsPrivate:           false,
			Status:              1,
			FocusMinutesToday:   180,
			LastSessionInMinute: &last60,
			Members: []roomMemberSeed{
				{UserKey: "mentor", Role: 1},
				{UserKey: "focus", Role: 0},
			},
			Sessions: []roomSessionSeed{
				{
					Key:             "morning-today",
					StartedBy:       "mentor",
					StartMinutesAgo: 75,
					DurationMinutes: 90,
					Topic:           "晨间读论文与写总结",
					Messages: []chatMessageSeed{
						{UserKey: "mentor", OffsetMin: 5, Content: "今天先过一遍实验配置，留 20 分钟总结。", MsgType: 0},
						{UserKey: "focus", OffsetMin: 25, Content: "我在准备看板接口，稍后分享进度。", MsgType: 0},
					},
				},
			},
		},
		{
			Key:                 "night",
			Name:                "夜猫冲刺房",
			Description:         "夜间复盘与补课，限制 8 人。",
			OwnerUserKey:        "focus",
			Tags:                []string{"夜间", "冲刺"},
			MaxMembers:          8,
			IsPrivate:           true,
			AccessCode:          "nightowl",
			Status:              1,
			FocusMinutesToday:   120,
			LastSessionInMinute: &last720,
			Members: []roomMemberSeed{
				{UserKey: "focus", Role: 1},
				{UserKey: "rookie", Role: 0},
			},
			Sessions: []roomSessionSeed{
				{
					Key:             "night-yesterday",
					StartedBy:       "focus",
					StartMinutesAgo: 840,
					DurationMinutes: 80,
					Topic:           "夜间刷题与英语听力",
					Messages: []chatMessageSeed{
						{UserKey: "focus", OffsetMin: 10, Content: "今晚想把听力练习做完，欢迎陪跑。", MsgType: 0},
						{UserKey: "rookie", OffsetMin: 40, Content: "我在背单词，有空一起讨论记忆法。", MsgType: 0},
					},
				},
			},
		},
	}

	for _, seed := range seeds {
		owner, ok := ctx.Users[seed.OwnerUserKey]
		if !ok {
			return fmt.Errorf("room owner %s not found", seed.OwnerUserKey)
		}

		var teamID *uint64
		if seed.TeamKey != nil {
			team, ok := ctx.Teams[*seed.TeamKey]
			if !ok {
				return fmt.Errorf("team %s not found for room %s", *seed.TeamKey, seed.Name)
			}
			teamID = &team.ID
		}

		var lastSessionStarted *time.Time
		if seed.LastSessionInMinute != nil {
			t := now.Add(-time.Duration(*seed.LastSessionInMinute) * time.Minute)
			lastSessionStarted = &t
		}

		roomData := models.StudyRoom{
			Name:               seed.Name,
			OwnerUserID:        owner.ID,
			TeamID:             teamID,
			Description:        seed.Description,
			Tags:               strings.Join(seed.Tags, ","),
			MaxMembers:         seed.MaxMembers,
			IsPrivate:          seed.IsPrivate,
			Status:             seed.Status,
			AccessCode:         "",
			FocusMinutesToday:  seed.FocusMinutesToday,
			LastSessionStarted: lastSessionStarted,
		}
		if seed.IsPrivate && seed.AccessCode != "" {
			roomData.AccessCode = hashPassword(seed.AccessCode)
		}

		var room models.StudyRoom
		err := tx.Where("name = ?", seed.Name).First(&room).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			room = roomData
			if err := tx.Create(&room).Error; err != nil {
				return fmt.Errorf("create study room %s: %w", seed.Name, err)
			}
		} else if err != nil {
			return err
		} else {
			if err := tx.Model(&room).Updates(roomData).Error; err != nil {
				return fmt.Errorf("update study room %s: %w", seed.Name, err)
			}
		}

		for _, memberSeed := range seed.Members {
			user, ok := ctx.Users[memberSeed.UserKey]
			if !ok {
				return fmt.Errorf("room member %s not found", memberSeed.UserKey)
			}
			var member models.StudyRoomMember
			err := tx.Where("room_id = ? AND user_id = ?", room.ID, user.ID).First(&member).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				record := models.StudyRoomMember{
					RoomID:   room.ID,
					UserID:   user.ID,
					Role:     memberSeed.Role,
					JoinedAt: now.AddDate(0, 0, -7),
				}
				if err := tx.Create(&record).Error; err != nil {
					return fmt.Errorf("add room member %s: %w", memberSeed.UserKey, err)
				}
			} else if err != nil {
				return err
			} else {
				if err := tx.Model(&member).Update("role", memberSeed.Role).Error; err != nil {
					return fmt.Errorf("update room member %s: %w", memberSeed.UserKey, err)
				}
			}
		}

		for _, sessionSeed := range seed.Sessions {
			start := now.Add(-time.Duration(sessionSeed.StartMinutesAgo) * time.Minute)
			end := start.Add(time.Duration(sessionSeed.DurationMinutes) * time.Minute)
			starter, ok := ctx.Users[sessionSeed.StartedBy]
			if !ok {
				return fmt.Errorf("session starter %s not found", sessionSeed.StartedBy)
			}

			var session models.RoomSession
			err := tx.Where("room_id = ? AND start_time = ?", room.ID, start).First(&session).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				entry := models.RoomSession{
					RoomID:    room.ID,
					StartedBy: starter.ID,
					StartTime: start,
					EndTime:   &end,
					Topic:     sessionSeed.Topic,
				}
				if err := tx.Create(&entry).Error; err != nil {
					return fmt.Errorf("create room session %s: %w", sessionSeed.Key, err)
				}
				session = entry
			} else if err != nil {
				return err
			}

			for _, msgSeed := range sessionSeed.Messages {
				user, ok := ctx.Users[msgSeed.UserKey]
				if !ok {
					return fmt.Errorf("chat user %s not found", msgSeed.UserKey)
				}
				sentAt := start.Add(time.Duration(msgSeed.OffsetMin) * time.Minute)
				var msg models.ChatMessage
				err := tx.Where("session_id = ? AND user_id = ? AND content = ?", session.ID, user.ID, msgSeed.Content).
					First(&msg).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					entry := models.ChatMessage{
						SessionID: session.ID,
						RoomID:    room.ID,
						UserID:    user.ID,
						Content:   msgSeed.Content,
						MsgType:   msgSeed.MsgType,
						SentAt:    sentAt,
					}
					if err := tx.Create(&entry).Error; err != nil {
						return fmt.Errorf("create chat message in %s: %w", sessionSeed.Key, err)
					}
				} else if err != nil {
					return err
				}
			}

			ctx.Sessions[sessionSeed.Key] = &session
		}

		ctx.StudyRooms[seed.Key] = &room
	}
	return nil
}

type ledgerSeed struct {
	UserKey      string
	SourceType   models.PointsSourceType
	TaskKey      string
	Delta        int
	BalanceAfter int
	Remark       string
	DayOffset    int
}

func seedDemoPointsLedger(tx *gorm.DB, ctx *demoSeedContext) error {
	now := time.Now()
	seeds := []ledgerSeed{
		{UserKey: "mentor", SourceType: models.PointsSourceTaskCompletion, TaskKey: "algo-notes", Delta: 80, BalanceAfter: 1200, Remark: "完成算法专题任务", DayOffset: -9},
		{UserKey: "mentor", SourceType: models.PointsSourceStudyRoom, TaskKey: "", Delta: 30, BalanceAfter: 1230, Remark: "晨间自习打卡", DayOffset: -1},
		{UserKey: "mentor", SourceType: models.PointsSourceDailyCheckIn, TaskKey: "", Delta: 20, BalanceAfter: 1250, Remark: "每日签到", DayOffset: 0},
		{UserKey: "focus", SourceType: models.PointsSourceTaskCompletion, TaskKey: "team-retro", Delta: 60, BalanceAfter: 520, Remark: "完成团队复盘", DayOffset: -6},
		{UserKey: "focus", SourceType: models.PointsSourceStudyRoom, TaskKey: "", Delta: 25, BalanceAfter: 545, Remark: "夜间学习房间打卡", DayOffset: -1},
		{UserKey: "rookie", SourceType: models.PointsSourceStudyRoom, TaskKey: "", Delta: 15, BalanceAfter: 120, Remark: "参与夜猫冲刺房", DayOffset: -1},
	}

	for _, seed := range seeds {
		user, ok := ctx.Users[seed.UserKey]
		if !ok {
			return fmt.Errorf("ledger user %s not found", seed.UserKey)
		}

		var sourceID *uint64
		if seed.TaskKey != "" {
			task, ok := ctx.Tasks[seed.TaskKey]
			if !ok {
				return fmt.Errorf("ledger task %s not found", seed.TaskKey)
			}
			sourceID = &task.ID
		}

		createdAt := now.AddDate(0, 0, seed.DayOffset)
		var ledger models.PointsLedger
		err := tx.Where("user_id = ? AND source_type = ? AND delta = ? AND remark = ?", user.ID, seed.SourceType, seed.Delta, seed.Remark).
			First(&ledger).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			record := models.PointsLedger{
				UserID:       user.ID,
				SourceType:   seed.SourceType,
				SourceID:     sourceID,
				Delta:        seed.Delta,
				BalanceAfter: seed.BalanceAfter,
				Remark:       seed.Remark,
				BaseModel: models.BaseModel{
					CreatedAt: createdAt,
					UpdatedAt: createdAt,
				},
			}
			if err := tx.Create(&record).Error; err != nil {
				return fmt.Errorf("create ledger for user %s: %w", seed.UserKey, err)
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func seedDemoUserAchievements(tx *gorm.DB, ctx *demoSeedContext) error {
	if len(ctx.Achievements) == 0 {
		var achievements []models.Achievement
		if err := tx.Find(&achievements).Error; err != nil {
			return fmt.Errorf("load achievements: %w", err)
		}
		for _, item := range achievements {
			ctx.Achievements[item.Code] = item
		}
	}

	now := time.Now()
	for userKey, codes := range ctx.userAchievementCodes {
		user, ok := ctx.Users[userKey]
		if !ok {
			return fmt.Errorf("user %s not found for achievements", userKey)
		}
		for idx, code := range codes {
			achievement, ok := ctx.Achievements[code]
			if !ok {
				continue
			}
			awardedAt := now.Add(-time.Duration(idx+1) * time.Hour)
			record := models.UserAchievement{
				UserID:        user.ID,
				AchievementID: achievement.ID,
				AwardedAt:     awardedAt,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}, {Name: "achievement_id"}},
				DoNothing: true,
			}).Create(&record).Error; err != nil {
				return fmt.Errorf("seed user achievement %s for %s: %w", code, userKey, err)
			}
		}
	}
	return nil
}

func ptrInt(value int) *int {
	return &value
}

func strPtr(value string) *string {
	return &value
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}
