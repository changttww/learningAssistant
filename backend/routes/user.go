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
)

type userProfileResponse struct {
	ID         uint64                `json:"id"`
	Account    string                `json:"account"`
	Display    string                `json:"display_name"`
	AvatarURL  string                `json:"avatar_url"`
	Bio        string                `json:"bio"`
	Role       string                `json:"role"`
	Status     string                `json:"status"`
	BasicInfo  userBasicInfo         `json:"basic_info"`
	Badges     []string              `json:"badges"`
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
	LevelLabel        string  `json:"level_label"`
	CurrentLevel      int     `json:"current_level"`
	CurrentPoints     int     `json:"current_points"`
	NextLevelPoints   int     `json:"next_level_points"`
	ProgressPercent   int     `json:"progress_percent"`
	DistanceToNext    int     `json:"distance_to_next"`
	TotalStudyHours   float64 `json:"total_study_hours"`
	TasksCompleted    int     `json:"tasks_completed"`
	CertificatesCount int     `json:"certificates_count"`
	StudyGroups       int     `json:"study_groups"`
	TaskCompletionRate int    `json:"task_completion_rate"`
	TasksInProgress   int     `json:"tasks_in_progress"`
	RankLabel         string  `json:"rank_label"`
	StreakDays        int     `json:"streak_days"`
	CoursesInProgress int     `json:"courses_in_progress"`
}

type userAchievementResponse struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AwardedAt   string `json:"awarded_at"`
	Type        string `json:"type"`
}

type userSkillsResponse struct {
	Primary   []string `json:"primary"`
	Secondary []string `json:"secondary"`
}

type userSettingsResponse struct {
	Notifications userNotificationPreferences `json:"notifications"`
	Privacy       userPrivacySettings          `json:"privacy"`
	StudyHabits   userStudyHabitSettings       `json:"study_habits"`
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

type userStudyRecord struct {
	ID        uint64  `json:"id"`
	Title     string  `json:"title"`
	Category  string  `json:"category"`
	Duration  float64 `json:"duration"`
	Completed bool    `json:"completed"`
	RecordedAt string `json:"recorded_at"`
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
	router.PUT("/:userId", handleUpdateUserProfile)
	router.GET("/:userId/study-stats", handleGetUserStudyStats)
	router.GET("/:userId/achievements", handleGetUserAchievements)
	router.GET("/:userId/skills", handleGetUserSkills)
	router.PUT("/:userId/skills", handleUpdateUserSkills)
	router.GET("/:userId/study-records", handleGetUserStudyRecords)
	router.GET("/:userId/settings", handleGetUserSettings)
	router.PUT("/:userId/settings", handleUpdateUserSettings)
	router.GET("/:userId/notification-preferences", handleGetUserNotificationPreferences)
	router.PUT("/:userId/notification-preferences", handleUpdateUserNotificationPreferences)
	router.POST("/avatar", handleUploadAvatar)
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

type updateUserProfileRequest struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	BasicInfo   *struct {
		RealName string `json:"real_name"`
		Email    string `json:"email"`
		School   string `json:"school"`
		Major    string `json:"major"`
		Location string `json:"location"`
		JoinDate string `json:"join_date"`
	} `json:"basic_info"`
	Preferences *struct {
		Language string `json:"language"`
		Theme    string `json:"theme"`
	} `json:"preferences"`
}

func handleUpdateUserProfile(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	var req updateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数格式不正确",
		})
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			return err
		}

		if strings.TrimSpace(req.DisplayName) != "" {
			user.DisplayName = strings.TrimSpace(req.DisplayName)
		}
		if req.Bio != "" {
			user.Bio = req.Bio
		}
		if req.AvatarURL != "" {
			user.AvatarURL = req.AvatarURL
		}
		if req.Role != "" {
			user.Role = mapRoleToCode(req.Role)
		}
		if req.Status != "" {
			user.Status = mapStatusToCode(req.Status)
		}
		if req.Preferences != nil {
			if req.Preferences.Language != "" {
				user.PreferredLanguage = req.Preferences.Language
			}
			if req.Preferences.Theme != "" {
				user.PreferredTheme = req.Preferences.Theme
			}
		}
		if req.BasicInfo != nil {
			if req.BasicInfo.RealName != "" {
				user.DisplayName = req.BasicInfo.RealName
			}
			if req.BasicInfo.Email != "" && !strings.EqualFold(user.Email, req.BasicInfo.Email) {
				if err := ensureEmailUnique(tx, req.BasicInfo.Email, user.ID); err != nil {
					return err
				}
				user.Email = req.BasicInfo.Email
			}
			if req.BasicInfo.School != "" {
				user.School = req.BasicInfo.School
			}
			if req.BasicInfo.Major != "" {
				user.Major = req.BasicInfo.Major
			}
			if req.BasicInfo.Location != "" {
				user.Location = req.BasicInfo.Location
			}
			if req.BasicInfo.JoinDate != "" {
				user.JoinDate = req.BasicInfo.JoinDate
			}
		}

		return tx.Save(&user).Error
	})

	if err != nil {
		if errorsIsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		if strings.Contains(err.Error(), "邮箱已存在") {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新用户信息失败"})
		return
	}

	handleGetUserProfile(c)
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

func handleGetUserAchievements(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	db := database.GetDB()
	var achievements []models.UserAchievement
	if err := db.Where("user_id = ?", userID).Order("awarded_at DESC").Find(&achievements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取成就失败"})
		return
	}

	items := make([]userAchievementResponse, 0, len(achievements))
	for _, ach := range achievements {
		items = append(items, userAchievementResponse{
			ID:          ach.ID,
			Title:       ach.Title,
			Description: ach.Description,
			AwardedAt:   ach.AwardedAt.Format("2006-01-02"),
			Type:        ach.Type,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"items": items,
			"total": len(items),
		},
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

func handleUpdateUserSkills(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	var req updateUserSkillsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数格式不正确"})
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserSkill{}).Error; err != nil {
			return err
		}

		for _, name := range req.Primary {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			if err := tx.Create(&models.UserSkill{
				UserID:    userID,
				Name:      name,
				Category:  "primary",
			}).Error; err != nil {
				return err
			}
		}
		for _, name := range req.Secondary {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			if err := tx.Create(&models.UserSkill{
				UserID:    userID,
				Name:      name,
				Category:  "secondary",
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新技能失败"})
		return
	}

	handleGetUserSkills(c)
}

func handleGetUserStudyRecords(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	db := database.GetDB()
	var records []models.LearningRecord
	if err := db.Where("user_id = ?", userID).Order("session_start DESC").Limit(30).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取学习记录失败"})
		return
	}

	if len(records) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data": gin.H{
				"items": []userStudyRecord{},
				"total": 0,
			},
		})
		return
	}

	taskIDs := make([]uint64, 0, len(records))
	for _, rec := range records {
		taskIDs = append(taskIDs, rec.TaskID)
	}

	taskMap := map[uint64]models.Task{}
	if len(taskIDs) > 0 {
		var tasks []models.Task
		if err := db.Where("id IN ?", taskIDs).Find(&tasks).Error; err == nil {
			for _, task := range tasks {
				taskMap[task.ID] = task
			}
		}
	}

	categoryIDs := []uint64{}
	for _, task := range taskMap {
		if task.CategoryID != nil {
			categoryIDs = append(categoryIDs, *task.CategoryID)
		}
	}

	categoryMap := map[uint64]string{}
	if len(categoryIDs) > 0 {
		var categories []models.TaskCategory
		if err := db.Where("id IN ?", categoryIDs).Find(&categories).Error; err == nil {
			for _, cat := range categories {
				categoryMap[cat.ID] = cat.Name
			}
		}
	}

	items := make([]userStudyRecord, 0, len(records))
	for _, rec := range records {
		task := taskMap[rec.TaskID]
		category := "自定义学习"
		if task.CategoryID != nil {
			if name, ok := categoryMap[*task.CategoryID]; ok {
				category = name
			}
		}
		if task.Title != "" {
			category = category
		}

		durationHours := float64(rec.DurationMinutes) / 60.0
		durationHours = math.Round(durationHours*10) / 10
		if durationHours == 0 {
			durationHours = float64(rec.DurationMinutes) / 60.0
		}

		item := userStudyRecord{
			ID:        rec.ID,
			Title:     task.Title,
			Category:  category,
			Duration:  durationHours,
			Completed: rec.SessionEnd.After(rec.SessionStart),
			RecordedAt: rec.SessionEnd.Format("2006-01-02"),
		}
		if item.Title == "" {
			item.Title = "学习记录"
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"items": items,
			"total": len(items),
		},
	})
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
	Privacy       *userPrivacySettings          `json:"privacy"`
	StudyHabits   *userStudyHabitSettings       `json:"study_habits"`
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

func handleGetUserNotificationPreferences(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	settings, err := ensureUserSettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户通知设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": userNotificationPreferences{
			Email:   settings.NotifyEmail,
			SMS:     settings.NotifySMS,
			InApp:   settings.NotifyInApp,
			Summary: settings.NotifySummary,
		},
	})
}

func handleUpdateUserNotificationPreferences(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	var req userNotificationPreferences
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数格式不正确"})
		return
	}

	settings, err := ensureUserSettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户通知设置失败"})
		return
	}

	settings.NotifyEmail = req.Email
	settings.NotifySMS = req.SMS
	settings.NotifyInApp = req.InApp
	settings.NotifySummary = req.Summary

	if err := database.GetDB().Save(settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新通知设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "notification preferences updated",
		"data": userNotificationPreferences{
			Email:   settings.NotifyEmail,
			SMS:     settings.NotifySMS,
			InApp:   settings.NotifyInApp,
			Summary: settings.NotifySummary,
		},
	})
}

func handleUploadAvatar(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "暂未开放头像上传功能",
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
		LevelLabel:        levelLabel,
		CurrentLevel:      profile.Level,
		CurrentPoints:     profile.TotalPoints,
		NextLevelPoints:   nextLevel,
		ProgressPercent:   progress,
		DistanceToNext:    distance,
		TotalStudyHours:   totalHours,
		TasksCompleted:    profile.TasksCompleted,
		CertificatesCount: profile.CertificatesCount,
		StudyGroups:       profile.StudyGroups,
		TaskCompletionRate: rate,
		TasksInProgress:   profile.TasksInProgress,
		RankLabel:         profile.RankLabel,
		StreakDays:        profile.StreakDays,
		CoursesInProgress: profile.CoursesInProgress,
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
		LevelLabel:        "学霸 Lv.1",
		CurrentLevel:      1,
		CurrentPoints:     0,
		NextLevelPoints:   200,
		ProgressPercent:   0,
		DistanceToNext:    200,
		TotalStudyHours:   0,
		TasksCompleted:    0,
		CertificatesCount: 0,
		StudyGroups:       0,
		TaskCompletionRate: 0,
		TasksInProgress:   0,
		RankLabel:         "TOP 100%",
		StreakDays:        0,
		CoursesInProgress: 0,
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
			UserID:           user.ID,
			TotalPoints:      0,
			Level:            1,
			TotalStudyMins:   0,
			TasksCompleted:   0,
			TasksInProgress:  0,
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
