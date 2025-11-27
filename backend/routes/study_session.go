package routes

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

const (
	sessionTimeout = 2 * time.Minute
)

type startStudySessionRequest struct {
	UserID   uint64  `json:"user_id"`
	Source   string  `json:"source"`
	SourceID *uint64 `json:"source_id"`
	Note     string  `json:"note"`
}

type pingStudySessionRequest struct {
	SessionID uint64 `json:"session_id"`
	UserID    uint64 `json:"user_id"`
}

type endStudySessionRequest struct {
	SessionID uint64 `json:"session_id"`
	UserID    uint64 `json:"user_id"`
}

func registerStudySessionRoutes(router *gin.RouterGroup) {
	router.POST("/start", handleStartStudySession)
	router.POST("/ping", handlePingStudySession)
	router.POST("/end", handleEndStudySession)
	router.POST("/aggregate/daily", handleAggregateDailyStudyStats)
}

func handleStartStudySession(c *gin.Context) {
	var req startStudySessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求格式错误"})
		return
	}
	if req.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少用户ID"})
		return
	}

	now := time.Now()
	db := database.GetDB()
	_, _ = autoCloseExpiredSessions(db, now)

	session := models.StudySession{
		UserID:     req.UserID,
		Source:     normalizeSessionSource(req.Source),
		SourceID:   req.SourceID,
		StartTime:  now,
		LastPingAt: now,
		Note:       strings.TrimSpace(req.Note),
	}

	if err := db.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建学习会话失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"session_id": session.ID,
			"start_time": session.StartTime.Format(time.RFC3339),
			"source":     session.Source,
		},
	})
}

func handlePingStudySession(c *gin.Context) {
	var req pingStudySessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求格式错误"})
		return
	}
	if req.SessionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 session_id"})
		return
	}

	db := database.GetDB()
	now := time.Now()
	_, _ = autoCloseExpiredSessions(db, now)

	var session models.StudySession
	if err := db.First(&session, req.SessionID).Error; err != nil {
		status := http.StatusInternalServerError
		msg := "加载学习会话失败"
		if errorsIsNotFound(err) {
			status = http.StatusNotFound
			msg = "会话不存在"
		}
		c.JSON(status, gin.H{"code": status, "message": msg})
		return
	}
	if req.UserID != 0 && session.UserID != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "会话归属不匹配"})
		return
	}

	if session.EndTime != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "session ended",
			"data": gin.H{
				"session_id": session.ID,
				"end_time":   session.EndTime.Format(time.RFC3339),
				"duration":   session.DurationMinutes,
				"ended":      true,
			},
		})
		return
	}

	if now.Sub(session.LastPingAt) > sessionTimeout {
		_ = finalizeStudySession(db, &session, session.LastPingAt)
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "session timeout, closed automatically",
			"data": gin.H{
				"session_id": session.ID,
				"end_time":   session.LastPingAt.Format(time.RFC3339),
				"duration":   session.DurationMinutes,
				"ended":      true,
			},
		})
		return
	}

	if err := db.Model(&models.StudySession{}).
		Where("id = ? AND end_time IS NULL", session.ID).
		Update("last_ping_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新心跳失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "pong",
		"data": gin.H{
			"session_id": session.ID,
			"last_ping":  now.Format(time.RFC3339),
			"ended":      false,
		},
	})
}

func handleEndStudySession(c *gin.Context) {
	var req endStudySessionRequest
	_ = c.ShouldBindJSON(&req)

	if req.SessionID == 0 {
		if idStr := c.Query("session_id"); idStr != "" {
			req.SessionID = parseUint64(idStr)
		}
	}
	if req.SessionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 session_id"})
		return
	}

	db := database.GetDB()
	now := time.Now()
	_, _ = autoCloseExpiredSessions(db, now)

	var session models.StudySession
	if err := db.First(&session, req.SessionID).Error; err != nil {
		status := http.StatusInternalServerError
		msg := "加载学习会话失败"
		if errorsIsNotFound(err) {
			status = http.StatusNotFound
			msg = "会话不存在"
		}
		c.JSON(status, gin.H{"code": status, "message": msg})
		return
	}
	if req.UserID != 0 && session.UserID != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "会话归属不匹配"})
		return
	}

	if session.EndTime == nil {
		if err := finalizeStudySession(db, &session, now); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "结束学习会话失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "session ended",
		"data": gin.H{
			"session_id": session.ID,
			"start_time": session.StartTime.Format(time.RFC3339),
			"end_time":   session.EndTime.Format(time.RFC3339),
			"duration":   session.DurationMinutes,
			"ended":      true,
		},
	})
}

func finalizeStudySession(db *gorm.DB, session *models.StudySession, endTime time.Time) error {
	if session.EndTime != nil {
		return nil
	}
	if endTime.IsZero() {
		endTime = time.Now()
	}
	if endTime.Before(session.StartTime) {
		endTime = session.StartTime
	}

	duration := calculateDurationMinutes(session.StartTime, endTime)
	update := map[string]interface{}{
		"end_time":         endTime,
		"duration_minutes": duration,
		"last_ping_at":     endTime,
	}
	if err := db.Model(&models.StudySession{}).
		Where("id = ? AND end_time IS NULL", session.ID).
		Updates(update).Error; err != nil {
		return err
	}

	session.EndTime = &endTime
	session.DurationMinutes = duration
	session.LastPingAt = endTime
	return nil
}

func autoCloseExpiredSessions(db *gorm.DB, now time.Time) (int, error) {
	cutoff := now.Add(-sessionTimeout)
	var stale []models.StudySession
	if err := db.Where("end_time IS NULL AND last_ping_at < ?", cutoff).
		Find(&stale).Error; err != nil {
		return 0, err
	}

	closed := 0
	for i := range stale {
		if err := finalizeStudySession(db, &stale[i], stale[i].LastPingAt); err != nil {
			log.Printf("auto close session %d failed: %v", stale[i].ID, err)
			continue
		}
		closed++
	}
	return closed, nil
}

func normalizeSessionSource(source string) string {
	val := strings.TrimSpace(strings.ToLower(source))
	switch val {
	case "study_room", "task", "focus", "video":
		return val
	default:
		return "unknown"
	}
}

func calculateDurationMinutes(start, end time.Time) int {
	if end.Before(start) {
		return 0
	}
	duration := int(end.Sub(start).Minutes())
	if duration == 0 && !end.Equal(start) {
		return 1
	}
	return duration
}

func parseUint64(val string) uint64 {
	num, _ := strconv.ParseUint(strings.TrimSpace(val), 10, 64)
	return num
}

type aggregateDailyRequest struct {
	Date   string `json:"date"`
	UserID uint64 `json:"user_id"`
}

func handleAggregateDailyStudyStats(c *gin.Context) {
	var req aggregateDailyRequest
	_ = c.ShouldBindJSON(&req)

	targetDay := normalizeDay(time.Now())
	if strings.TrimSpace(req.Date) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(req.Date))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "date 格式应为 YYYY-MM-DD"})
			return
		}
		targetDay = normalizeDay(parsed.In(time.Now().Location()))
	}

	db := database.GetDB()
	autoClosed, _ := autoCloseExpiredSessions(db, time.Now())

	updated, err := aggregateDailyStats(db, targetDay, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "聚合学习统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"date":          targetDay.Format("2006-01-02"),
			"updated_users": updated,
			"auto_closed":   autoClosed,
		},
	})
}

type dailyAggregation struct {
	Minutes               int
	SessionCount          int
	NightMinutes          int
	MorningMinutes        int
	FocusModeMinutes      int
	StudyRoomMinutes      int
	StudyRoomNightMinutes int
}

func aggregateDailyStats(db *gorm.DB, day time.Time, userFilter uint64) (int, error) {
	dayStart := normalizeDay(day)
	dayEnd := dayStart.Add(24 * time.Hour)

	query := db.Where("end_time IS NOT NULL AND end_time >= ? AND end_time < ?", dayStart, dayEnd)
	if userFilter != 0 {
		query = query.Where("user_id = ?", userFilter)
	}

	var sessions []models.StudySession
	if err := query.Find(&sessions).Error; err != nil {
		return 0, err
	}

	aggMap := make(map[uint64]*dailyAggregation)
	for i := range sessions {
		session := sessions[i]
		if session.EndTime == nil {
			continue
		}
		agg := aggMap[session.UserID]
		if agg == nil {
			agg = &dailyAggregation{}
			aggMap[session.UserID] = agg
		}

		duration := session.DurationMinutes
		if duration == 0 {
			duration = calculateDurationMinutes(session.StartTime, *session.EndTime)
		}
		agg.Minutes += duration
		agg.SessionCount++

		night := nightMinutesForDay(dayStart, session.StartTime, *session.EndTime)
		morning := minutesInWindow(session.StartTime, *session.EndTime, dayStart, dayStart.Add(8*time.Hour))
		agg.NightMinutes += night
		agg.MorningMinutes += morning

		switch session.Source {
		case "focus":
			agg.FocusModeMinutes += duration
		case "study_room":
			agg.StudyRoomMinutes += duration
			agg.StudyRoomNightMinutes += night
		}
	}

	updated := 0
	for userID, stat := range aggMap {
		record := models.DailyStudyStat{
			UserID:                userID,
			Date:                  dayStart,
			Minutes:               stat.Minutes,
			SessionCount:          stat.SessionCount,
			NightMinutes:          stat.NightMinutes,
			MorningMinutes:        stat.MorningMinutes,
			FocusModeMinutes:      stat.FocusModeMinutes,
			StudyRoomMinutes:      stat.StudyRoomMinutes,
			StudyRoomNightMinutes: stat.StudyRoomNightMinutes,
		}

		onConflict := clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}},
			DoUpdates: clause.AssignmentColumns([]string{"minutes", "session_count", "night_minutes", "morning_minutes", "focus_mode_minutes", "study_room_minutes", "study_room_night_minutes", "updated_at"}),
		}
		if err := db.Clauses(onConflict).Create(&record).Error; err != nil {
			return updated, err
		}

		totals, err := loadUserStudyTotals(db, userID)
		if err != nil {
			return updated, err
		}
		if err := upsertUserProfileStudyMinutes(db, userID, totals.Minutes); err != nil {
			return updated, err
		}
		if err := updateStudyAchievementProgress(db, userID, totals); err != nil {
			return updated, err
		}
		updated++
	}

	return updated, nil
}

type studyTotals struct {
	Minutes               int
	NightMinutes          int
	MorningMinutes        int
	FocusModeMinutes      int
	StudyRoomMinutes      int
	StudyRoomNightMinutes int
}

func loadUserStudyTotals(db *gorm.DB, userID uint64) (studyTotals, error) {
	var totals studyTotals
	err := db.Model(&models.DailyStudyStat{}).
		Select("COALESCE(SUM(minutes),0) as minutes, COALESCE(SUM(night_minutes),0) as night_minutes, COALESCE(SUM(morning_minutes),0) as morning_minutes, COALESCE(SUM(focus_mode_minutes),0) as focus_mode_minutes, COALESCE(SUM(study_room_minutes),0) as study_room_minutes, COALESCE(SUM(study_room_night_minutes),0) as study_room_night_minutes").
		Where("user_id = ?", userID).
		Scan(&totals).Error
	return totals, err
}

func upsertUserProfileStudyMinutes(db *gorm.DB, userID uint64, totalMinutes int) error {
	var profile models.UserProfile
	err := db.Where("user_id = ?", userID).First(&profile).Error
	if errorsIsNotFound(err) {
		profile = models.UserProfile{
			UserID:          userID,
			Level:           1,
			NextLevelPoints: 200,
			RankLabel:       "TOP 100%",
			TotalStudyMins:  totalMinutes,
		}
		return db.Create(&profile).Error
	}
	if err != nil {
		return err
	}

	return db.Model(&models.UserProfile{}).
		Where("id = ?", profile.ID).
		Update("total_study_mins", totalMinutes).Error
}

func updateStudyAchievementProgress(db *gorm.DB, userID uint64, totals studyTotals) error {
	var progress models.UserAchievementProgress
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&progress).Error
	if errorsIsNotFound(err) {
		progress = models.UserAchievementProgress{UserID: userID}
		if err := db.Create(&progress).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"study_room_duration_mins": totals.StudyRoomMinutes,
		"study_room_night_mins":    totals.StudyRoomNightMinutes,
		"night_study_mins":         totals.NightMinutes,
		"morning_study_mins":       totals.MorningMinutes,
		"focus_mode_mins":          totals.FocusModeMinutes,
	}

	return db.Model(&models.UserAchievementProgress{}).
		Where("id = ?", progress.ID).
		Updates(updates).Error
}

func normalizeDay(t time.Time) time.Time {
	loc := t.Location()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

func minutesInWindow(start, end, windowStart, windowEnd time.Time) int {
	if end.Before(windowStart) || start.After(windowEnd) {
		return 0
	}
	s := start
	if s.Before(windowStart) {
		s = windowStart
	}
	e := end
	if e.After(windowEnd) {
		e = windowEnd
	}
	if e.Before(s) {
		return 0
	}
	return int(e.Sub(s).Minutes())
}

func nightMinutesForDay(dayStart, start, end time.Time) int {
	earlyStart := dayStart
	earlyEnd := dayStart.Add(2 * time.Hour)
	lateStart := dayStart.Add(22 * time.Hour)
	lateEnd := dayStart.Add(24 * time.Hour)

	return minutesInWindow(start, end, earlyStart, earlyEnd) +
		minutesInWindow(start, end, lateStart, lateEnd)
}
