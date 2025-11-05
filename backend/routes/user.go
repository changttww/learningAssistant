package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// userProfileResponse describes the shape of user profile data returned to the client.
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
	TaskCompletionRate float64 `json:"task_completion_rate"`
	TasksInProgress    int     `json:"tasks_in_progress"`
	RankLabel          string  `json:"rank_label"`
	StreakDays         int     `json:"streak_days"`
	CoursesInProgress  int     `json:"courses_in_progress"`
	TotalTasks         int     `json:"total_tasks"`
}

type userAchievement struct {
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

type userAuthAccount struct {
	ID           uint64 `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	DisplayName  string `json:"display_name"`
}

var (
	userDataMu sync.RWMutex
	userData   = map[uint64]*userProfileResponse{
		1: {
			ID:        1,
			Account:   "student001",
			Display:   "张同学",
			AvatarURL: "https://avatars.dicebear.com/api/initials/%E5%BC%A0.svg",
			Bio:       "专注于教育科技领域，热衷于学习新技术，致力于打造高效的学习工具。",
			Role:      "技术项目经理",
			Status:    "online",
			BasicInfo: userBasicInfo{
				RealName: "张同学",
				Email:    "zhang@example.com",
				School:   "北京大学",
				Major:    "计算机科学与技术",
				Location: "北京市",
				JoinDate: "2023年9月",
			},
			Badges: []string{"学习达人", "团队协作", "效率先锋"},
			Preferences: userPreferencesBrief{
				Language: "zh-CN",
				Theme:    "light",
			},
		},
	}

	userStudyStats = map[uint64]*userStudyStatsResponse{
		1: {
			LevelLabel:         "学霸 Lv.4",
			CurrentLevel:       4,
			CurrentPoints:      3860,
			NextLevelPoints:    4200,
			ProgressPercent:    75,
			DistanceToNext:     340,
			TotalStudyHours:    87.5,
			TasksCompleted:     156,
			CertificatesCount:  24,
			StudyGroups:        3,
			TaskCompletionRate: 92,
			TasksInProgress:    8,
			RankLabel:          "TOP 15%",
			StreakDays:         28,
			CoursesInProgress:  3,
			TotalTasks:         102,
		},
	}

	userAchievements = map[uint64][]userAchievement{
		1: {
			{ID: 101, Title: "连续学习7天", Description: "连续7天每天学习超过2小时", AwardedAt: "2024-01-15", Type: "streak"},
			{ID: 102, Title: "完成Vue项目", Description: "完成前端高级课程项目", AwardedAt: "2024-01-12", Type: "project"},
			{ID: 103, Title: "团队协作达人", Description: "协同完成团队任务5次", AwardedAt: "2024-01-10", Type: "team"},
		},
	}

	userSkills = map[uint64]*userSkillsResponse{
		1: {
			Primary:   []string{"Vue.js", "JavaScript", "Node.js"},
			Secondary: []string{"项目管理", "UI设计", "数据分析"},
		},
	}

	userSettings = map[uint64]*userSettingsResponse{
		1: {
			Notifications: userNotificationPreferences{Email: true, SMS: false, InApp: true, Summary: true},
			Privacy:       userPrivacySettings{ShowEmail: false, ShowProfile: true, ShowStudyData: true},
			StudyHabits:   userStudyHabitSettings{DailyGoalMinutes: 120, PreferredPeriod: "evening", FocusMode: true},
		},
	}

	userStudyRecords = map[uint64][]userStudyRecord{
		1: {
			{ID: 201, Title: "Vue3 组合式API学习", Category: "前端开发", Duration: 2.5, Completed: true, RecordedAt: "2024-01-15"},
			{ID: 202, Title: "数据可视化实践", Category: "数据分析", Duration: 1.8, Completed: true, RecordedAt: "2024-01-14"},
			{ID: 203, Title: "团队项目协作", Category: "项目管理", Duration: 2.0, Completed: false, RecordedAt: "2024-01-14"},
		},
	}

	userIDCounter uint64 = 1

	authAccountsByUsername = map[string]*userAuthAccount{}
	authAccountsByEmail    = map[string]*userAuthAccount{}
)

func init() {
	userDataMu.Lock()
	defer userDataMu.Unlock()

	defaultProfile, exists := userData[1]
	if !exists {
		return
	}

	account := &userAuthAccount{
		ID:           1,
		Username:     defaultProfile.Account,
		Email:        defaultProfile.BasicInfo.Email,
		PasswordHash: hashPassword("password123"),
		DisplayName:  defaultProfile.Display,
	}

	authAccountsByUsername[normalizeUsername(account.Username)] = account
	if account.Email != "" {
		authAccountsByEmail[normalizeEmail(account.Email)] = account
	}
}

// registerUserRoutes setups all user related endpoints.
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

func parseUserID(c *gin.Context) (uint64, bool) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid user id",
		})
		return 0, false
	}
	return userID, true
}

func handleGetUserProfile(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	userDataMu.RLock()
	defer userDataMu.RUnlock()

	profile, exists := userData[userID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    profile,
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
			"message": "invalid request body",
		})
		return
	}

	userDataMu.Lock()
	defer userDataMu.Unlock()

	profile, exists := userData[userID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "user not found",
		})
		return
	}

	if req.DisplayName != "" {
		profile.Display = req.DisplayName
		updateAuthAccountDisplayNameLocked(profile.ID, req.DisplayName)
	}
	if req.Bio != "" {
		profile.Bio = req.Bio
	}
	if req.AvatarURL != "" {
		profile.AvatarURL = req.AvatarURL
	}
	if req.Role != "" {
		profile.Role = req.Role
	}
	if req.Status != "" {
		profile.Status = req.Status
	}
	if req.BasicInfo != nil {
		if req.BasicInfo.RealName != "" {
			profile.BasicInfo.RealName = req.BasicInfo.RealName
		}
		if req.BasicInfo.Email != "" && req.BasicInfo.Email != profile.BasicInfo.Email {
			updateAuthAccountEmailLocked(profile.ID, req.BasicInfo.Email)
			profile.BasicInfo.Email = req.BasicInfo.Email
		}
		if req.BasicInfo.School != "" {
			profile.BasicInfo.School = req.BasicInfo.School
		}
		if req.BasicInfo.Major != "" {
			profile.BasicInfo.Major = req.BasicInfo.Major
		}
		if req.BasicInfo.Location != "" {
			profile.BasicInfo.Location = req.BasicInfo.Location
		}
		if req.BasicInfo.JoinDate != "" {
			profile.BasicInfo.JoinDate = req.BasicInfo.JoinDate
		}
	}
	if req.Preferences != nil {
		if req.Preferences.Language != "" {
			profile.Preferences.Language = req.Preferences.Language
		}
		if req.Preferences.Theme != "" {
			profile.Preferences.Theme = req.Preferences.Theme
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "user profile updated",
		"data":    profile,
	})
}

func handleGetUserStudyStats(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	userDataMu.RLock()
	defer userDataMu.RUnlock()

	stats, exists := userStudyStats[userID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "study stats not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

func handleGetUserAchievements(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	userDataMu.RLock()
	defer userDataMu.RUnlock()

	achievements, exists := userAchievements[userID]
	if !exists {
		achievements = []userAchievement{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"items": achievements,
			"total": len(achievements),
		},
	})
}

type updateUserSkillsRequest struct {
	Primary   []string `json:"primary"`
	Secondary []string `json:"secondary"`
}

func handleGetUserSkills(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	userDataMu.RLock()
	defer userDataMu.RUnlock()

	skills, exists := userSkills[userID]
	if !exists {
		skills = &userSkillsResponse{
			Primary:   []string{},
			Secondary: []string{},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    skills,
	})
}

func handleUpdateUserSkills(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	var req updateUserSkillsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request body",
		})
		return
	}

	userDataMu.Lock()
	defer userDataMu.Unlock()

	userSkills[userID] = &userSkillsResponse{
		Primary:   append([]string{}, req.Primary...),
		Secondary: append([]string{}, req.Secondary...),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "skills updated",
		"data":    userSkills[userID],
	})
}

func handleGetUserStudyRecords(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	userDataMu.RLock()
	defer userDataMu.RUnlock()

	records := userStudyRecords[userID]
	if records == nil {
		records = []userStudyRecord{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"items": records,
			"total": len(records),
		},
	})
}

func handleGetUserSettings(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	userDataMu.RLock()
	defer userDataMu.RUnlock()

	settings, exists := userSettings[userID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "settings not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    settings,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request body",
		})
		return
	}

	userDataMu.Lock()
	defer userDataMu.Unlock()

	settings, exists := userSettings[userID]
	if !exists {
		settings = &userSettingsResponse{}
		userSettings[userID] = settings
	}

	if req.Notifications != nil {
		settings.Notifications = *req.Notifications
	}
	if req.Privacy != nil {
		settings.Privacy = *req.Privacy
	}
	if req.StudyHabits != nil {
		settings.StudyHabits = *req.StudyHabits
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "settings updated",
		"data":    settings,
	})
}

func handleGetUserNotificationPreferences(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	userDataMu.RLock()
	defer userDataMu.RUnlock()

	settings, exists := userSettings[userID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "notification preferences not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    settings.Notifications,
	})
}

func handleUpdateUserNotificationPreferences(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	var req userNotificationPreferences
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request body",
		})
		return
	}

	userDataMu.Lock()
	defer userDataMu.Unlock()

	settings, exists := userSettings[userID]
	if !exists {
		settings = &userSettingsResponse{}
		userSettings[userID] = settings
	}

	settings.Notifications = req

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "notification preferences updated",
		"data":    settings.Notifications,
	})
}

func handleUploadAvatar(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(5 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "failed to parse form data",
		})
		return
	}

	file, _, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "avatar file is required",
		})
		return
	}
	defer file.Close()

	// 由于当前未接入实际的文件存储，这里直接返回一个示例地址
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "avatar uploaded",
		"data": gin.H{
			"url": "https://avatars.dicebear.com/api/identicon/demo.svg",
		},
	})
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

func updateAuthAccountDisplayNameLocked(userID uint64, displayName string) {
	for key, account := range authAccountsByUsername {
		if account.ID == userID {
			account.DisplayName = displayName
			authAccountsByUsername[key] = account
			break
		}
	}
}

func updateAuthAccountEmailLocked(userID uint64, newEmail string) {
	var target *userAuthAccount
	for key, account := range authAccountsByUsername {
		if account.ID == userID {
			oldEmailKey := normalizeEmail(account.Email)
			if oldEmailKey != "" {
				delete(authAccountsByEmail, oldEmailKey)
			}
			account.Email = newEmail
			authAccountsByUsername[key] = account
			target = account
			break
		}
	}

	if target != nil {
		newEmailKey := normalizeEmail(newEmail)
		if newEmailKey != "" {
			authAccountsByEmail[newEmailKey] = target
		}
	}
}

type authUserSummary struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	AvatarURL   string `json:"avatar_url"`
}

func composeUserSummary(userID uint64) *authUserSummary {
	profile, exists := userData[userID]
	if !exists {
		return nil
	}

	account, hasAccount := authAccountsByUsername[normalizeUsername(profile.Account)]
	email := profile.BasicInfo.Email
	if hasAccount && account.Email != "" {
		email = account.Email
	}

	return &authUserSummary{
		ID:          profile.ID,
		Username:    profile.Account,
		DisplayName: profile.Display,
		Email:       email,
		Role:        profile.Role,
		AvatarURL:   profile.AvatarURL,
	}
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

	nUsername := normalizeUsername(username)
	nEmail := normalizeEmail(email)
	display := strings.TrimSpace(displayName)
	if display == "" {
		display = username
	}

	userDataMu.Lock()
	defer userDataMu.Unlock()

	if _, exists := authAccountsByUsername[nUsername]; exists {
		return nil, fmt.Errorf("用户名已存在")
	}
	if nEmail != "" {
		if _, exists := authAccountsByEmail[nEmail]; exists {
			return nil, fmt.Errorf("该邮箱已注册")
		}
	}

	newID := userIDCounter + 1
	userIDCounter = newID

	profile := &userProfileResponse{
		ID:        newID,
		Account:   username,
		Display:   display,
		AvatarURL: "",
		Bio:       "这位同学还没有填写个人简介。",
		Role:      "学习者",
		Status:    "online",
		BasicInfo: userBasicInfo{
			RealName: display,
			Email:    email,
			School:   "",
			Major:    "",
			Location: "",
			JoinDate: time.Now().Format("2006年1月"),
		},
		Badges: []string{},
		Preferences: userPreferencesBrief{
			Language: "zh-CN",
			Theme:    "light",
		},
	}
	userData[newID] = profile

	userStudyStats[newID] = &userStudyStatsResponse{
		LevelLabel:         "入门学员 Lv.1",
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
		TotalTasks:         0,
	}

	userAchievements[newID] = []userAchievement{}
	userSkills[newID] = &userSkillsResponse{
		Primary:   []string{},
		Secondary: []string{},
	}
	userSettings[newID] = &userSettingsResponse{
		Notifications: userNotificationPreferences{Email: true, SMS: false, InApp: true, Summary: true},
		Privacy:       userPrivacySettings{ShowEmail: false, ShowProfile: true, ShowStudyData: true},
		StudyHabits:   userStudyHabitSettings{DailyGoalMinutes: 60, PreferredPeriod: "evening", FocusMode: false},
	}
	userStudyRecords[newID] = []userStudyRecord{}

	account := &userAuthAccount{
		ID:           newID,
		Username:     username,
		Email:        email,
		PasswordHash: hashPassword(password),
		DisplayName:  display,
	}

	authAccountsByUsername[nUsername] = account
	if nEmail != "" {
		authAccountsByEmail[nEmail] = account
	}

	return composeUserSummary(newID), nil
}

func authenticateUser(identifier, password string) (*authUserSummary, error) {
	account, exists := getAuthAccount(identifier)
	if !exists {
		return nil, fmt.Errorf("用户不存在")
	}
	if !verifyPassword(account.PasswordHash, password) {
		return nil, fmt.Errorf("密码不正确")
	}
	return getUserSummary(account.ID)
}

func getUserSummary(userID uint64) (*authUserSummary, error) {
	userDataMu.RLock()
	defer userDataMu.RUnlock()

	summary := composeUserSummary(userID)
	if summary == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	return summary, nil
}

func getAuthAccount(identifier string) (*userAuthAccount, bool) {
	normalizedUsername := normalizeUsername(identifier)
	normalizedEmail := normalizeEmail(identifier)

	userDataMu.RLock()
	defer userDataMu.RUnlock()

	if account, ok := authAccountsByUsername[normalizedUsername]; ok {
		return account, true
	}
	if normalizedEmail != "" {
		if account, ok := authAccountsByEmail[normalizedEmail]; ok {
			return account, true
		}
	}
	return nil, false
}
