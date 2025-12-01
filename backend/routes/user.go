package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
	"learningAssistant-backend/services/achievement"
	"learningAssistant-backend/services/points"
)

type userProfileResponse struct {
	ID          uint64               `json:"id"`
	Account     string               `json:"account"`
	Display     string               `json:"display_name"`
	AvatarURL   string               `json:"avatar_url"`
	Bio         string               `json:"bio"`
	Role        string               `json:"role"`
	Status      string               `json:"status"`
	BasicInfo   userBasicInfo        `json:"basic_info"`
	Badges      []string             `json:"badges"`
	Preferences userPreferencesBrief `json:"preferences"`
}

type userBasicInfo struct {
	RealName string `json:"real_name"`
	Email    string `json:"email"`
	School   string `json:"school"`
	Major    string `json:"major"`
	Location string `json:"location"`
	JoinDate string `json:"join_date"`
}

type userPreferencesBrief struct {
	Language string `json:"language"`
	Theme    string `json:"theme"`
}

type userStudyStatsResponse struct {
	LevelLabel         string  `json:"level_label"`
	CurrentLevel       int     `json:"current_level"`
	CurrentPoints      int     `json:"current_points"`
	NextLevelPoints    int     `json:"next_level_points"`
	ProgressPercent    int     `json:"progress_percent"`
	DistanceToNext     int     `json:"distance_to_next"`
	TotalStudyHours    float64 `json:"total_study_hours"`
	TasksCompleted     int     `json:"tasks_completed"`
	CertificatesCount  int     `json:"certificates_count"`
	StudyGroups        int     `json:"study_groups"`
	TaskCompletionRate int     `json:"task_completion_rate"`
	TasksInProgress    int     `json:"tasks_in_progress"`
	RankLabel          string  `json:"rank_label"`
	StreakDays         int     `json:"streak_days"`
	CoursesInProgress  int     `json:"courses_in_progress"`
}

type userAchievementResponse struct {
	ID            uint64 `json:"id"`
	AchievementID uint   `json:"achievement_id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Category      string `json:"category"`
	Icon          string `json:"icon"`
	AwardedAt     string `json:"awarded_at"`
}

type userSkillsResponse struct {
	Primary   []string `json:"primary"`
	Secondary []string `json:"secondary"`
}

type userNotificationPreferences struct {
	Email   bool `json:"email"`
	SMS     bool `json:"sms"`
	InApp   bool `json:"in_app"`
	Summary bool `json:"weekly_summary"`
}

type userPrivacySettings struct {
	ShowEmail     bool `json:"show_email"`
	ShowProfile   bool `json:"show_profile"`
	ShowStudyData bool `json:"show_study_data"`
}

type userStudyHabitSettings struct {
	DailyGoalMinutes int    `json:"daily_goal_minutes"`
	PreferredPeriod  string `json:"preferred_period"`
	FocusMode        bool   `json:"focus_mode"`
}

type userSettingsResponse struct {
	Notifications userNotificationPreferences `json:"notifications"`
	Privacy       userPrivacySettings         `json:"privacy"`
	StudyHabits   userStudyHabitSettings      `json:"study_habits"`
}

type authUserSummary struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	AvatarURL   string `json:"avatar_url"`
}

func registerUserRoutes(router *gin.RouterGroup) {
	router.GET("/:userId", handleGetUserProfile)
	router.GET("/:userId/study-stats", handleGetUserStudyStats)
	router.POST("/:userId/check-in", handleUserDailyCheckIn)
	router.GET("/:userId/points/ledger", handleGetUserPointsLedger)
	router.GET("/:userId/achievements", handleGetUserAchievements)
	router.GET("/:userId/skills", handleGetUserSkills)
	router.GET("/:userId/settings", handleGetUserSettings)
	router.PUT("/:userId/settings", handleUpdateUserSettings)
	router.GET("/:userId/notification-preferences", handleGetUserNotificationPreferences)
	router.PUT("/:userId/notification-preferences", handleUpdateUserNotificationPreferences)
	router.POST("/avatar", handleUploadAvatar)

	// 学习伙伴
	router.GET("/:userId/buddies", handleListStudyBuddies)
	router.POST("/:userId/buddies", handleAddStudyBuddy)
	router.PUT("/:userId/buddies/:buddyId", handleUpdateStudyBuddy)
	router.DELETE("/:userId/buddies/:buddyId", handleDeleteStudyBuddy)
}

func handleGetUserProfile(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	db := database.GetDB()
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if errorsIsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	var badges []models.UserBadge
	if err := db.Where("user_id = ?", userID).Order("created_at ASC").Find(&badges).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户徽章失败"})
		return
	}

	response := buildUserProfileResponse(&user, badges)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}
func handleGetUserStudyStats(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	db := database.GetDB()
	var profile models.UserProfile
	err := db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errorsIsNotFound(err) {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "success",
				"data":    defaultStudyStats(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取学习统计失败"})
		return
	}

	response := buildStudyStatsResponse(&profile)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}

func handleUserDailyCheckIn(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	db := database.GetDB()
	var user models.User
	if err := db.Select("id").First(&user, userID).Error; err != nil {
		if errorsIsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询用户失败"})
		return
	}

	var updatedProfile models.UserProfile
	err := db.Transaction(func(tx *gorm.DB) error {
		var profile models.UserProfile
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&profile).Error
		if err != nil {
			if !errorsIsNotFound(err) {
				return err
			}
			profile = models.UserProfile{
				UserID:             userID,
				Level:              1,
				NextLevelPoints:    200,
				RankLabel:          "TOP 100%",
				StreakDays:         1,
				TaskCompletionRate: 0,
			}
			if err := tx.Create(&profile).Error; err != nil {
				return err
			}
			updatedProfile = profile
			return nil
		}

		profile.StreakDays++
		if err := tx.Model(&models.UserProfile{}).
			Where("id = ?", profile.ID).
			Update("streak_days", profile.StreakDays).Error; err != nil {
			return err
		}
		updatedProfile = profile
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "签到失败"})
		return
	}

	pointResult, err := points.AwardDailyCheckIn(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "签到积分发放失败"})
		return
	}
	updatedProfile = *pointResult.Profile

	if err := achievement.ProcessEvent(achievement.Event{
		Type:   achievement.EventStreakUpdated,
		UserID: userID,
		Value:  updatedProfile.StreakDays,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "签到事件处理失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "check-in success",
		"data":    buildPointsAwardData(pointResult, updatedProfile.StreakDays),
	})
}

func handleGetUserPointsLedger(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}
	if !ensureUserExists(c, userID) {
		return
	}

	limit := 50
	if query := strings.TrimSpace(c.Query("limit")); query != "" {
		if parsed, err := strconv.Atoi(query); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	records, err := points.ListLedger(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取积分记录失败"})
		return
	}

	profile, err := loadUserProfileSnapshot(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取积分概要失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"items":   records,
			"total":   len(records),
			"summary": buildPointsSummary(profile),
		},
	})
}

func handleGetUserAchievements(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	overview, err := achievement.BuildOverview(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取成就失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    overview,
	})
}

func handleGetUserSkills(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	db := database.GetDB()
	var skills []models.UserSkill
	if err := db.Where("user_id = ?", userID).Order("created_at ASC").Find(&skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取技能失败"})
		return
	}

	response := userSkillsResponse{
		Primary:   []string{},
		Secondary: []string{},
	}
	for _, skill := range skills {
		if skill.Category == "primary" {
			response.Primary = append(response.Primary, skill.Name)
		} else {
			response.Secondary = append(response.Secondary, skill.Name)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}

type updateUserSkillsRequest struct {
	Primary   []string `json:"primary"`
	Secondary []string `json:"secondary"`
}

func handleGetUserSettings(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	settings, err := ensureUserSettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    buildSettingsResponse(settings),
	})
}

type updateUserSettingsRequest struct {
	Notifications *userNotificationPreferences `json:"notifications"`
	Privacy       *userPrivacySettings         `json:"privacy"`
	StudyHabits   *userStudyHabitSettings      `json:"study_habits"`
}

func handleUpdateUserSettings(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	var req updateUserSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数格式不正确"})
		return
	}

	settings, err := ensureUserSettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户设置失败"})
		return
	}

	if req.Notifications != nil {
		settings.NotifyEmail = req.Notifications.Email
		settings.NotifySMS = req.Notifications.SMS
		settings.NotifyInApp = req.Notifications.InApp
		settings.NotifySummary = req.Notifications.Summary
	}
	if req.Privacy != nil {
		settings.ShowEmail = req.Privacy.ShowEmail
		settings.ShowProfile = req.Privacy.ShowProfile
		settings.ShowStudyData = req.Privacy.ShowStudyData
	}
	if req.StudyHabits != nil {
		if req.StudyHabits.DailyGoalMinutes > 0 {
			settings.DailyGoalMinutes = req.StudyHabits.DailyGoalMinutes
		}
		if strings.TrimSpace(req.StudyHabits.PreferredPeriod) != "" {
			settings.PreferredPeriod = req.StudyHabits.PreferredPeriod
		}
		settings.FocusMode = req.StudyHabits.FocusMode
	}

	if err := database.GetDB().Save(settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新用户设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "settings updated",
		"data":    buildSettingsResponse(settings),
	})
}

func parseUserID(c *gin.Context) (uint64, bool) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不正确",
		})
		return 0, false
	}
	return userID, true
}

func ensureUserExists(c *gin.Context, userID uint64) bool {
	var user models.User
	if err := database.GetDB().Select("id").First(&user, userID).Error; err != nil {
		if errorsIsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询用户失败"})
		}
		return false
	}
	return true
}

func loadUserProfileSnapshot(userID uint64) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := database.GetDB().Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if !errorsIsNotFound(err) {
			return nil, err
		}
		// 返回默认档案用于展示
		profile = models.UserProfile{
			UserID:          userID,
			Level:           1,
			NextLevelPoints: 200,
			RankLabel:       "TOP 100%",
		}
	}
	return &profile, nil
}

func buildPointsAwardData(result *points.AwardResult, streakDays int) gin.H {
	data := gin.H{
		"points_added":      result.Ledger.Delta,
		"balance_after":     result.Ledger.BalanceAfter,
		"current_level":     result.Profile.Level,
		"next_level_points": result.Profile.NextLevelPoints,
		"ledger":            result.Ledger,
	}
	if streakDays > 0 {
		data["streak_days"] = streakDays
	}
	return data
}

func buildPointsSummary(profile *models.UserProfile) gin.H {
	nextLevel := profile.NextLevelPoints
	if nextLevel <= 0 {
		nextLevel = 200
	}
	current := profile.TotalPoints
	progress := 0
	if nextLevel > 0 {
		progress = int(math.Round(float64(current) / float64(nextLevel) * 100))
		if progress > 100 {
			progress = 100
		}
		if progress < 0 {
			progress = 0
		}
	}
	distance := nextLevel - current
	if distance < 0 {
		distance = 0
	}

	levelLabel := fmt.Sprintf("学霸 Lv.%d", profile.Level)
	if profile.Level <= 0 {
		levelLabel = "学霸 Lv.1"
	}

	return gin.H{
		"current_points":    current,
		"next_level_points": nextLevel,
		"distance_to_next":  distance,
		"progress_percent":  progress,
		"current_level":     profile.Level,
		"level_label":       levelLabel,
	}
}

func handleListStudyBuddies(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}
	db := database.GetDB()
	var buddies []models.StudyBuddy
	if err := db.Where("user_id = ?", userID).Find(&buddies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取学习伙伴失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": buddies})
}

func handleAddStudyBuddy(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}
	var req struct {
		BuddyID uint64 `json:"buddy_id"`
		Remark  string `json:"remark"`
		Tags    string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数格式不正确"})
		return
	}
	if req.BuddyID == 0 || req.BuddyID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "伙伴ID无效"})
		return
	}
	db := database.GetDB()
	var count int64
	db.Model(&models.StudyBuddy{}).
		Where("user_id = ? AND buddy_id = ?", userID, req.BuddyID).
		Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已存在", "data": nil})
		return
	}
	buddy := models.StudyBuddy{
		UserID:  userID,
		BuddyID: req.BuddyID,
		Remark:  strings.TrimSpace(req.Remark),
		Tags:    strings.TrimSpace(req.Tags),
	}
	if err := db.Create(&buddy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "添加学习伙伴失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": buddy})
}

func handleUpdateStudyBuddy(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}
	buddyID, err := strconv.ParseUint(c.Param("buddyId"), 10, 64)
	if err != nil || buddyID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "伙伴ID无效"})
		return
	}
	var req struct {
		Remark string `json:"remark"`
		Tags   string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数格式不正确"})
		return
	}
	db := database.GetDB()
	update := map[string]interface{}{
		"remark": strings.TrimSpace(req.Remark),
		"tags":   strings.TrimSpace(req.Tags),
	}
	if err := db.Model(&models.StudyBuddy{}).
		Where("user_id = ? AND buddy_id = ?", userID, buddyID).
		Updates(update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func handleDeleteStudyBuddy(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}
	buddyID, err := strconv.ParseUint(c.Param("buddyId"), 10, 64)
	if err != nil || buddyID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "伙伴ID无效"})
		return
	}
	db := database.GetDB()
	if err := db.Where("user_id = ? AND buddy_id = ?", userID, buddyID).
		Delete(&models.StudyBuddy{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func buildUserProfileResponse(user *models.User, badges []models.UserBadge) userProfileResponse {
	badgeNames := make([]string, 0, len(badges))
	for _, badge := range badges {
		badgeNames = append(badgeNames, badge.Name)
	}

	return userProfileResponse{
		ID:        user.ID,
		Account:   user.Account,
		Display:   user.DisplayName,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Role:      roleLabel(user.Role),
		Status:    statusLabel(user.Status),
		BasicInfo: userBasicInfo{
			RealName: user.DisplayName,
			Email:    user.Email,
			School:   user.School,
			Major:    user.Major,
			Location: user.Location,
			JoinDate: user.JoinDate,
		},
		Badges: badgeNames,
		Preferences: userPreferencesBrief{
			Language: user.PreferredLanguage,
			Theme:    user.PreferredTheme,
		},
	}
}

func buildStudyStatsResponse(profile *models.UserProfile) userStudyStatsResponse {
	nextLevel := profile.NextLevelPoints
	if nextLevel <= 0 {
		nextLevel = 200
	}
	progress := 0
	if nextLevel > 0 {
		progress = int(math.Round(float64(profile.TotalPoints) / float64(nextLevel) * 100))
		if progress > 100 {
			progress = 100
		}
		if progress < 0 {
			progress = 0
		}
	}
	distance := nextLevel - profile.TotalPoints
	if distance < 0 {
		distance = 0
	}
	totalHours := float64(profile.TotalStudyMins) / 60.0
	totalHours = math.Round(totalHours*10) / 10

	levelLabel := fmt.Sprintf("学霸 Lv.%d", profile.Level)
	if profile.Level <= 0 {
		levelLabel = "学霸 Lv.1"
	}

	rate := int(math.Round(float64(profile.TaskCompletionRate)))

	return userStudyStatsResponse{
		LevelLabel:         levelLabel,
		CurrentLevel:       profile.Level,
		CurrentPoints:      profile.TotalPoints,
		NextLevelPoints:    nextLevel,
		ProgressPercent:    progress,
		DistanceToNext:     distance,
		TotalStudyHours:    totalHours,
		TasksCompleted:     profile.TasksCompleted,
		CertificatesCount:  profile.CertificatesCount,
		StudyGroups:        profile.StudyGroups,
		TaskCompletionRate: rate,
		TasksInProgress:    profile.TasksInProgress,
		RankLabel:          profile.RankLabel,
		StreakDays:         profile.StreakDays,
		CoursesInProgress:  profile.CoursesInProgress,
	}
}

func buildSettingsResponse(settings *models.UserSetting) userSettingsResponse {
	return userSettingsResponse{
		Notifications: userNotificationPreferences{
			Email:   settings.NotifyEmail,
			SMS:     settings.NotifySMS,
			InApp:   settings.NotifyInApp,
			Summary: settings.NotifySummary,
		},
		Privacy: userPrivacySettings{
			ShowEmail:     settings.ShowEmail,
			ShowProfile:   settings.ShowProfile,
			ShowStudyData: settings.ShowStudyData,
		},
		StudyHabits: userStudyHabitSettings{
			DailyGoalMinutes: settings.DailyGoalMinutes,
			PreferredPeriod:  settings.PreferredPeriod,
			FocusMode:        settings.FocusMode,
		},
	}
}

func ensureUserSettings(userID uint64) (*models.UserSetting, error) {
	db := database.GetDB()
	var settings models.UserSetting
	err := db.Where("user_id = ?", userID).First(&settings).Error
	if err == nil {
		return &settings, nil
	}
	if !errorsIsNotFound(err) {
		return nil, err
	}

	settings = models.UserSetting{UserID: userID}
	if err := db.Create(&settings).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

func ensureEmailUnique(db *gorm.DB, email string, userID uint64) error {
	if strings.TrimSpace(email) == "" {
		return nil
	}
	var count int64
	if err := db.Model(&models.User{}).
		Where("LOWER(email) = ?", strings.ToLower(email)).
		Where("id <> ?", userID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("邮箱已存在")
	}
	return nil
}

func defaultStudyStats() userStudyStatsResponse {
	return userStudyStatsResponse{
		LevelLabel:         "学霸 Lv.1",
		CurrentLevel:       1,
		CurrentPoints:      0,
		NextLevelPoints:    200,
		ProgressPercent:    0,
		DistanceToNext:     200,
		TotalStudyHours:    0,
		TasksCompleted:     0,
		CertificatesCount:  0,
		StudyGroups:        0,
		TaskCompletionRate: 0,
		TasksInProgress:    0,
		RankLabel:          "TOP 100%",
		StreakDays:         0,
		CoursesInProgress:  0,
	}
}

func roleLabel(role int8) string {
	switch role {
	case 1:
		return "管理员"
	case 2:
		return "教师"
	default:
		return "学生"
	}
}

func mapRoleToCode(role string) int8 {
	role = strings.TrimSpace(role)
	switch role {
	case "管理员":
		return 1
	case "教师":
		return 2
	default:
		return 0
	}
}

func statusLabel(status int8) string {
	if status == 1 {
		return "online"
	}
	return "offline"
}

func mapStatusToCode(status string) int8 {
	switch strings.ToLower(status) {
	case "online", "active":
		return 1
	default:
		return 0
	}
}

func errorsIsNotFound(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}

func composeUserSummary(userID uint64) (*authUserSummary, error) {
	return composeUserSummaryWithDB(database.GetDB(), userID)
}

func composeUserSummaryWithDB(db *gorm.DB, userID uint64) (*authUserSummary, error) {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &authUserSummary{
		ID:          user.ID,
		Username:    user.Account,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        roleLabel(user.Role),
		AvatarURL:   user.AvatarURL,
	}, nil
}

func registerUserAccount(username, email, password, displayName string) (*authUserSummary, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("用户名不能为空")
	}
	password = strings.TrimSpace(password)
	if password == "" {
		return nil, fmt.Errorf("密码不能为空")
	}
	display := strings.TrimSpace(displayName)
	if display == "" {
		display = username
	}

	db := database.GetDB()
	var result *authUserSummary

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := ensureAccountUnique(tx, username); err != nil {
			return err
		}
		if strings.TrimSpace(email) != "" {
			if err := ensureEmailUnique(tx, email, 0); err != nil {
				return fmt.Errorf("该邮箱已注册")
			}
		}

		user := models.User{
			Account:           username,
			Email:             strings.TrimSpace(email),
			DisplayName:       display,
			Role:              0,
			AvatarURL:         "",
			Bio:               "这位同学还没有填写个人简介。",
			Status:            1,
			PasswordHash:      hashPassword(password),
			School:            "",
			Major:             "",
			Location:          "",
			JoinDate:          time.Now().Format("2006年1月"),
			PreferredLanguage: "zh-CN",
			PreferredTheme:    "light",
		}

		createTx := tx
		omitCols := []string{}
		if strings.TrimSpace(user.Phone) == "" {
			omitCols = append(omitCols, "Phone")
		}
		if strings.TrimSpace(user.Email) == "" {
			omitCols = append(omitCols, "Email")
		}
		if len(omitCols) > 0 {
			createTx = createTx.Omit(omitCols...)
		}
		if err := createTx.Create(&user).Error; err != nil {
			return err
		}

		profile := models.UserProfile{
			UserID:             user.ID,
			TotalPoints:        0,
			Level:              1,
			TotalStudyMins:     0,
			TasksCompleted:     0,
			TasksInProgress:    0,
			TaskCompletionRate: 0,
			CertificatesCount:  0,
			StudyGroups:        0,
			RankLabel:          "TOP 100%",
			StreakDays:         0,
			CoursesInProgress:  0,
			NextLevelPoints:    200,
		}
		if err := tx.Create(&profile).Error; err != nil {
			return err
		}

		settings := models.UserSetting{
			UserID:           user.ID,
			NotifyEmail:      true,
			NotifySMS:        false,
			NotifyInApp:      true,
			NotifySummary:    true,
			ShowEmail:        false,
			ShowProfile:      true,
			ShowStudyData:    true,
			DailyGoalMinutes: 60,
			PreferredPeriod:  "evening",
			FocusMode:        false,
		}
		if err := tx.Create(&settings).Error; err != nil {
			return err
		}

		summary, err := composeUserSummaryWithDB(tx, user.ID)
		if err != nil {
			return err
		}
		result = summary
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func ensureAccountUnique(db *gorm.DB, account string) error {
	var count int64
	if err := db.Model(&models.User{}).
		Where("LOWER(account) = ?", strings.ToLower(account)).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("用户名已存在")
	}
	return nil
}

func authenticateUser(identifier, password string) (*authUserSummary, error) {
	normalized := strings.TrimSpace(identifier)
	if normalized == "" {
		return nil, fmt.Errorf("请输入用户名或邮箱")
	}

	db := database.GetDB()
	var user models.User
	if err := db.Where("LOWER(account) = ? OR LOWER(email) = ?", strings.ToLower(normalized), strings.ToLower(normalized)).
		First(&user).Error; err != nil {
		if errorsIsNotFound(err) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	if !verifyPassword(user.PasswordHash, password) {
		return nil, fmt.Errorf("密码不正确")
	}

	if summary, err := composeUserSummary(user.ID); err == nil {
		return summary, nil
	} else {
		return nil, err
	}
}

func getUserSummary(userID uint64) (*authUserSummary, error) {
	return composeUserSummary(userID)
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func verifyPassword(hash, password string) bool {
	return hash == hashPassword(password)
}
