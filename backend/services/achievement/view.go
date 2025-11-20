package achievement

import (
	"errors"
	"sort"
	"time"

	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

// ProgressSnapshot 用于前端展示的进度快照
type ProgressSnapshot struct {
	TaskCreatedCount       int `json:"task_created_count"`
	TaskCompletedCount     int `json:"task_completed_count"`
	StreakDays             int `json:"streak_days"`
	StudyRoomJoinCount     int `json:"study_room_join_count"`
	StudyRoomDurationMins  int `json:"study_room_duration_minutes"`
	StudyRoomNightMins     int `json:"study_room_night_minutes"`
	NightSessionMaxMins    int `json:"night_session_max_minutes"`
	StudyRoomChatCount     int `json:"study_room_chat_count"`
	StudyRoomLikesGiven    int `json:"study_room_likes_given"`
	StudyRoomLikesReceived int `json:"study_room_likes_received"`
	TeamTasksCompleted     int `json:"team_tasks_completed"`
}

// AchievementView 带进度的成就视图
type AchievementView struct {
	ID            uint              `json:"id"`
	Code          string            `json:"code"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Category      string            `json:"category"`
	ConditionType string            `json:"condition_type"`
	Icon          string            `json:"icon"`
	TargetValue   float64           `json:"target_value"`
	CurrentValue  float64           `json:"current_value"`
	Completed     bool              `json:"completed"`
	AwardedAt     string            `json:"awarded_at"`
	History       []AchievementView `json:"history,omitempty"`
}

type Overview struct {
	Unlocked []AchievementView `json:"unlocked"`
	Upcoming []AchievementView `json:"upcoming"`
	Progress ProgressSnapshot  `json:"progress"`
}

// BuildOverview 返回已解锁与待解锁的成就视图及当前进度
func BuildOverview(userID uint64) (*Overview, error) {
	if err := EnsureAchievementsForUser(userID); err != nil {
		return nil, err
	}

	db := database.GetDB()

	progress, err := loadProgress(db, userID)
	if err != nil {
		return nil, err
	}
	snapshot := ProgressSnapshot{
		TaskCreatedCount:       progress.TaskCreatedCount,
		TaskCompletedCount:     progress.TaskCompletedCount,
		StreakDays:             progress.StreakDays,
		StudyRoomJoinCount:     progress.StudyRoomJoinCount,
		StudyRoomDurationMins:  progress.StudyRoomDurationMins,
		StudyRoomNightMins:     progress.StudyRoomNightMins,
		NightSessionMaxMins:    progress.NightSessionMaxMins,
		StudyRoomChatCount:     progress.StudyRoomChatCount,
		StudyRoomLikesGiven:    progress.StudyRoomLikesGiven,
		StudyRoomLikesReceived: progress.StudyRoomLikesReceived,
		TeamTasksCompleted:     progress.TeamTasksCompleted,
	}

	var achievements []models.Achievement
	if err := db.Find(&achievements).Error; err != nil {
		return nil, err
	}

	unlockedMap, err := findUnlocked(db, userID)
	if err != nil {
		return nil, err
	}

	type unlockedWrap struct {
		view AchievementView
		time time.Time
	}

	unlockedByType := make(map[string][]unlockedWrap)
	upcomingByType := make(map[string]AchievementView)

	for _, ach := range achievements {
		cond, err := parseCondition(ach.Condition)
		if err != nil {
			continue
		}
		target := cond.Value
		current := currentValueForCondition(cond.Type, progress)
		record := AchievementView{
			ID:            ach.ID,
			Code:          ach.Code,
			Name:          ach.Name,
			Description:   ach.Description,
			Category:      ach.Category,
			ConditionType: cond.Type,
			Icon:          ach.Icon,
			TargetValue:   target,
			CurrentValue:  current,
			Completed:     false,
			AwardedAt:     "",
		}

		if awardedAt, ok := unlockedMap[ach.ID]; ok {
			record.Completed = true
			record.AwardedAt = awardedAt.Format("2006-01-02")
			unlockedByType[cond.Type] = append(unlockedByType[cond.Type], unlockedWrap{
				view: record,
				time: awardedAt,
			})
		} else {
			// 同一类型只保留下一个目标值最低且尚未达成的成就
			existing, exists := upcomingByType[cond.Type]
			if !exists || record.TargetValue < existing.TargetValue {
				upcomingByType[cond.Type] = record
			}
		}
	}

	unlocked := make([]AchievementView, 0, len(unlockedByType))
	for _, items := range unlockedByType {
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].time.After(items[j].time)
		})
		primary := items[0].view
		if len(items) > 1 {
			for _, wrap := range items[1:] {
				primary.History = append(primary.History, wrap.view)
			}
		}
		unlocked = append(unlocked, primary)
	}

	upcoming := make([]AchievementView, 0, len(upcomingByType))
	for _, item := range upcomingByType {
		upcoming = append(upcoming, item)
	}

	return &Overview{
		Unlocked: unlocked,
		Upcoming: upcoming,
		Progress: snapshot,
	}, nil
}

func loadProgress(db *gorm.DB, userID uint64) (*models.UserAchievementProgress, error) {
	var progress models.UserAchievementProgress
	err := db.Where("user_id = ?", userID).First(&progress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			progress = models.UserAchievementProgress{
				UserID: userID,
			}
			if err := db.Create(&progress).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &progress, nil
}

func findUnlocked(db *gorm.DB, userID uint64) (map[uint]time.Time, error) {
	var records []models.UserAchievement
	if err := db.Where("user_id = ?", userID).Find(&records).Error; err != nil {
		return nil, err
	}
	result := make(map[uint]time.Time, len(records))
	for _, record := range records {
		result[record.AchievementID] = record.AwardedAt
	}
	return result, nil
}

func currentValueForCondition(condType string, progress *models.UserAchievementProgress) float64 {
	switch condType {
	case "task_total_completed":
		return float64(progress.TaskCompletedCount)
	case "task_total_created":
		return float64(progress.TaskCreatedCount)
	case "streak_task_completion":
		return float64(progress.StreakDays)
	case "studyroom_join_count":
		return float64(progress.StudyRoomJoinCount)
	case "studyroom_duration_hours":
		return minutesToHours(progress.StudyRoomDurationMins)
	case "studyroom_night_hours":
		return minutesToHours(progress.StudyRoomNightMins)
	case "studyroom_chat_count":
		return float64(progress.StudyRoomChatCount)
	case "studyroom_likes_given":
		return float64(progress.StudyRoomLikesGiven)
	case "studyroom_likes_received":
		return float64(progress.StudyRoomLikesReceived)
	case "team_task_completed":
		return float64(progress.TeamTasksCompleted)
	default:
		return 0
	}
}
