package points

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

const (
	taskCompletionReward      = 20
	checkInReward             = 2
	studyRoomRewardPerWindow  = 10
	studyRoomRewardWindowMins = 30
)

var (
	ErrInsufficientDuration = errors.New("insufficient_studyroom_duration")
	ErrInvalidPointsDelta   = errors.New("invalid_points_delta")
)

// AwardResult 封装积分发放结果
type AwardResult struct {
	Ledger  *models.PointsLedger
	Profile *models.UserProfile
}

// AwardTaskCompletion 任务完成加分
func AwardTaskCompletion(userID uint64, taskID uint64) (*AwardResult, error) {
	return applyPoints(userID, models.PointsSourceTaskCompletion, &taskID, taskCompletionReward, fmt.Sprintf("完成任务 #%d", taskID))
}

// AwardStudyRoomDuration 自习室在线时长加分，按 30 分钟一档
func AwardStudyRoomDuration(userID uint64, roomID *uint64, durationMinutes int) (*AwardResult, error) {
	if durationMinutes <= 0 {
		return nil, fmt.Errorf("%w: durationMinutes 必须大于 0", ErrInsufficientDuration)
	}
	blocks := durationMinutes / studyRoomRewardWindowMins
	if blocks <= 0 {
		return nil, fmt.Errorf("%w: 在线时长不足 %d 分钟，未产生积分", ErrInsufficientDuration, studyRoomRewardWindowMins)
	}

	points := blocks * studyRoomRewardPerWindow
	remark := fmt.Sprintf("自习室在线 %d 分钟", blocks*studyRoomRewardWindowMins)
	return applyPoints(userID, models.PointsSourceStudyRoom, roomID, points, remark)
}

// AwardDailyCheckIn 签到加分
func AwardDailyCheckIn(userID uint64) (*AwardResult, error) {
	return applyPoints(userID, models.PointsSourceDailyCheckIn, nil, checkInReward, "每日签到")
}

// ListLedger 获取积分账本记录
func ListLedger(userID uint64, limit int) ([]models.PointsLedger, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var records []models.PointsLedger
	if err := database.GetDB().
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func applyPoints(userID uint64, sourceType models.PointsSourceType, sourceID *uint64, delta int, remark string) (*AwardResult, error) {
	if delta <= 0 {
		return nil, fmt.Errorf("%w: 积分变动必须大于 0", ErrInvalidPointsDelta)
	}

	db := database.GetDB()
	var ledger models.PointsLedger
	var profile models.UserProfile

	err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		profile, err = loadOrCreateProfile(tx, userID)
		if err != nil {
			return err
		}

		newTotal := profile.TotalPoints + delta
		level, nextLevel, err := determineLevel(tx, newTotal)
		if err != nil {
			return err
		}

		profile.TotalPoints = newTotal
		profile.Level = level
		profile.NextLevelPoints = nextLevel

		if err := tx.Save(&profile).Error; err != nil {
			return err
		}

		ledger = models.PointsLedger{
			UserID:       userID,
			SourceType:   sourceType,
			SourceID:     sourceID,
			Delta:        delta,
			BalanceAfter: newTotal,
			Remark:       remark,
		}
		if err := tx.Create(&ledger).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &AwardResult{
		Ledger:  &ledger,
		Profile: &profile,
	}, nil
}

func loadOrCreateProfile(tx *gorm.DB, userID uint64) (models.UserProfile, error) {
	var profile models.UserProfile
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&profile).Error
	if err == nil {
		return profile, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return profile, err
	}

	profile = models.UserProfile{
		UserID:          userID,
		Level:           1,
		NextLevelPoints: 200,
		RankLabel:       "TOP 100%",
	}
	if err := tx.Create(&profile).Error; err != nil {
		return profile, err
	}
	return profile, nil
}

func determineLevel(tx *gorm.DB, totalPoints int) (int, int, error) {
	var rules []models.LevelRule
	if err := tx.Order("min_points ASC").Find(&rules).Error; err != nil {
		return 1, 200, err
	}

	level := 1
	nextThreshold := 200
	if len(rules) == 0 {
		if totalPoints > 0 {
			nextThreshold = totalPoints
		}
		return level, nextThreshold, nil
	}

	nextThreshold = rules[0].MinPoints
	for idx, rule := range rules {
		if totalPoints >= rule.MinPoints {
			level = rule.Level
			if idx+1 < len(rules) {
				nextThreshold = rules[idx+1].MinPoints
			} else {
				nextThreshold = totalPoints
			}
		} else {
			nextThreshold = rule.MinPoints
			break
		}
	}

	if nextThreshold < totalPoints {
		nextThreshold = totalPoints
	}
	if nextThreshold <= 0 {
		nextThreshold = 200
	}

	return level, nextThreshold, nil
}
