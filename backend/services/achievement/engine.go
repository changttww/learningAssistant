package achievement

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

// AchievementEventType 支持的触发事件类型
type AchievementEventType string

const (
	EventTaskCreated      AchievementEventType = "task_created"
	EventTaskCompleted    AchievementEventType = "task_completed"
	EventStreakUpdated    AchievementEventType = "streak_updated"
	EventStudyRoomJoin    AchievementEventType = "studyroom_join"
	EventTeamTaskFinished AchievementEventType = "team_task_finished"
)

// Event 成就事件
type Event struct {
	Type     AchievementEventType   `json:"type"`
	UserID   uint64                 `json:"user_id"`
	Value    int                    `json:"value"`
	Metadata map[string]interface{} `json:"metadata"`
}

// ProcessEvent 处理单个成就事件
func ProcessEvent(evt Event) error {
	if evt.UserID == 0 {
		return fmt.Errorf("achievement event missing user id")
	}

	db := database.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		progress, err := ensureProgressForUpdate(tx, evt.UserID)
		if err != nil {
			return err
		}

		if err := syncProgressFromProfile(tx, progress); err != nil {
			return err
		}

		applyEvent(progress, evt)

		if err := tx.Save(progress).Error; err != nil {
			return err
		}

		return evaluateAchievements(tx, progress)
	})
}

// EnsureAchievementsForUser 主动检查并解锁成就
func EnsureAchievementsForUser(userID uint64) error {
	if userID == 0 {
		return fmt.Errorf("invalid user id")
	}

	db := database.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		progress, err := ensureProgressForUpdate(tx, userID)
		if err != nil {
			return err
		}
		if err := syncProgressFromProfile(tx, progress); err != nil {
			return err
		}
		if err := tx.Save(progress).Error; err != nil {
			return err
		}
		return evaluateAchievements(tx, progress)
	})
}

func ensureProgressForUpdate(tx *gorm.DB, userID uint64) (*models.UserAchievementProgress, error) {
	var progress models.UserAchievementProgress
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&progress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			progress = models.UserAchievementProgress{UserID: userID}
			if err := tx.Create(&progress).Error; err != nil {
				return nil, err
			}
			return &progress, nil
		}
		return nil, err
	}
	return &progress, nil
}

func syncProgressFromProfile(tx *gorm.DB, progress *models.UserAchievementProgress) error {
	var profile models.UserProfile
	err := tx.Where("user_id = ?", progress.UserID).First(&profile).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		if profile.TasksCompleted > progress.TaskCompletedCount {
			progress.TaskCompletedCount = profile.TasksCompleted
		}
		if profile.StreakDays > progress.StreakDays {
			progress.StreakDays = profile.StreakDays
		}
	}
	return nil
}

func applyEvent(progress *models.UserAchievementProgress, evt Event) {
	switch evt.Type {
	case EventTaskCreated:
		progress.TaskCreatedCount++
	case EventTaskCompleted:
		progress.TaskCompletedCount++
	case EventStreakUpdated:
		streak := evt.Value
		if streak == 0 {
			streak = metaInt(evt.Metadata, "streak_days")
		}
		if streak < 0 {
			streak = 0
		}
		progress.StreakDays = streak
	case EventStudyRoomJoin:
		progress.StudyRoomJoinCount++
		progress.StudyRoomDurationMins += metaInt(evt.Metadata, "duration_minutes")

		nightMinutes := metaInt(evt.Metadata, "night_minutes")
		progress.StudyRoomNightMins += nightMinutes

		sessionNight := metaInt(evt.Metadata, "session_night_minutes")
		if sessionNight > progress.NightSessionMaxMins {
			progress.NightSessionMaxMins = sessionNight
		}

		progress.StudyRoomChatCount += metaInt(evt.Metadata, "chat_messages")
		progress.StudyRoomLikesGiven += metaInt(evt.Metadata, "likes_given")
		progress.StudyRoomLikesReceived += metaInt(evt.Metadata, "likes_received")
	case EventTeamTaskFinished:
		progress.TeamTasksCompleted++
	}
}

type achievementCondition struct {
	Type  string
	Value float64
	Raw   map[string]interface{}
}

func parseCondition(raw string) (achievementCondition, error) {
	result := achievementCondition{
		Raw: map[string]interface{}{},
	}
	if raw == "" {
		return result, fmt.Errorf("empty condition")
	}
	if err := json.Unmarshal([]byte(raw), &result.Raw); err != nil {
		return result, err
	}

	typeVal, ok := result.Raw["type"].(string)
	if !ok || typeVal == "" {
		return result, fmt.Errorf("condition missing type")
	}
	result.Type = typeVal
	if value, ok := result.Raw["value"]; ok {
		result.Value = toFloat(value)
	}
	return result, nil
}

func (c achievementCondition) extraString(key string) string {
	if value, ok := c.Raw[key]; ok {
		return fmt.Sprint(value)
	}
	return ""
}

func toFloat(value interface{}) float64 {
	switch v := value.(type) {
	case float32:
		return float64(v)
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		f, _ := v.Float64()
		return f
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	default:
		return 0
	}
}

func evaluateAchievements(tx *gorm.DB, progress *models.UserAchievementProgress) error {
	var achievements []models.Achievement
	if err := tx.Find(&achievements).Error; err != nil {
		return err
	}

	var unlockedIDs []uint
	if err := tx.Model(&models.UserAchievement{}).
		Where("user_id = ?", progress.UserID).
		Pluck("achievement_id", &unlockedIDs).Error; err != nil {
		return err
	}
	existing := make(map[uint]struct{}, len(unlockedIDs))
	for _, id := range unlockedIDs {
		existing[id] = struct{}{}
	}

	for _, ach := range achievements {
		if _, ok := existing[ach.ID]; ok {
			continue
		}
		cond, err := parseCondition(ach.Condition)
		if err != nil {
			continue
		}
		if !checkCondition(progress, cond) {
			continue
		}

		record := models.UserAchievement{
			UserID:        progress.UserID,
			AchievementID: ach.ID,
			AwardedAt:     time.Now(),
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&record).Error; err != nil {
			return err
		}
	}

	return nil
}

func checkCondition(progress *models.UserAchievementProgress, cond achievementCondition) bool {
	switch cond.Type {
	case "task_total_completed":
		return float64(progress.TaskCompletedCount) >= cond.Value
	case "task_total_created":
		return float64(progress.TaskCreatedCount) >= cond.Value
	case "streak_task_completion":
		return float64(progress.StreakDays) >= cond.Value
	case "studyroom_join_count":
		return float64(progress.StudyRoomJoinCount) >= cond.Value
	case "studyroom_duration_hours":
		return minutesToHours(progress.StudyRoomDurationMins) >= cond.Value
	case "studyroom_night_hours":
		mode := cond.extraString("mode")
		total := minutesToHours(progress.StudyRoomNightMins)
		single := minutesToHours(progress.NightSessionMaxMins)
		switch mode {
		case "single":
			return single >= cond.Value
		case "single_or_total":
			return single >= cond.Value || total >= cond.Value
		default:
			return total >= cond.Value
		}
	case "studyroom_chat_count":
		return float64(progress.StudyRoomChatCount) >= cond.Value
	case "studyroom_likes_given":
		return float64(progress.StudyRoomLikesGiven) >= cond.Value
	case "studyroom_likes_received":
		return float64(progress.StudyRoomLikesReceived) >= cond.Value
	case "team_task_completed":
		return float64(progress.TeamTasksCompleted) >= cond.Value
	default:
		return false
	}
}

func minutesToHours(mins int) float64 {
	if mins <= 0 {
		return 0
	}
	return float64(mins) / 60.0
}

func metaInt(meta map[string]interface{}, key string) int {
	if meta == nil {
		return 0
	}
	value, ok := meta[key]
	if !ok {
		return 0
	}
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		i, _ := v.Int64()
		return int(i)
	case string:
		n, _ := strconv.Atoi(v)
		return n
	default:
		return 0
	}
}
