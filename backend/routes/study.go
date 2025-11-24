package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

type studySummaryResponse struct {
	ActiveRooms     int     `json:"active_rooms"`
	OnlineUsers     int     `json:"online_users"`
	TodayStudyHours float64 `json:"today_study_hours"`
	StreakDays      int     `json:"streak_days"`
}

func registerStudyRoutes(router *gin.RouterGroup) {
	rooms := router.Group("/rooms")
	{
		rooms.GET("", handleListStudyRooms)
		rooms.GET("/:roomId", handleGetStudyRoomDetail)
		rooms.POST("", handleCreateStudyRoom)
		rooms.POST("/:roomId/join", handleJoinStudyRoom)
	}

	router.GET("/summary", handleStudySummary)
	router.GET("/records", handleStudyRecords)
}

func registerStudyWebsocketRoutes(router *gin.RouterGroup) {
	router.GET("/rooms/:roomId/ws", func(c *gin.Context) {
		roomIDParam := c.Param("roomId")
		roomID, err := strconv.ParseUint(roomIDParam, 10, 64)
		if err != nil || roomID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "房间ID不正确"})
			return
		}
		hub := studyHubRegistry.getHub(roomID)
		hub.handleWebSocket(c)
	})
}

type createStudyRoomRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RoomType    string   `json:"room_type"`
	Password    string   `json:"password"`
	MaxMembers  int      `json:"max_members"`
	Tags        []string `json:"tags"`
	OwnerUserID uint64   `json:"owner_user_id"`
	TeamID      *uint64  `json:"team_id"`
}

type studyRoomListItem struct {
	ID                 uint64   `json:"id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	Tags               []string `json:"tags"`
	Status             string   `json:"status"`
	StatusCode         int8     `json:"status_code"`
	StatusClass        string   `json:"status_class"`
	RoomType           string   `json:"room_type"`
	IsPrivate          bool     `json:"is_private"`
	CurrentUsers       int      `json:"current_users"`
	MaxUsers           int      `json:"max_users"`
	OwnerUserID        uint64   `json:"owner_user_id"`
	TeamID             *uint64  `json:"team_id"`
	StudyTime          string   `json:"study_time"`
	FocusMinutesToday  int      `json:"focus_minutes_today"`
	LastSessionStarted *string  `json:"last_session_started"`
}

func handleListStudyRooms(c *gin.Context) {
	db := database.GetDB()
	var rooms []models.StudyRoom
	if err := db.Order("created_at DESC").
		Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取学习房间列表失败",
		})
		return
	}

	memberCounts, err := loadStudyRoomMemberCounts(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "加载房间成员数据失败",
		})
		return
	}

	response := make([]studyRoomListItem, 0, len(rooms))
	for _, room := range rooms {
		response = append(response, buildStudyRoomListItem(&room, memberCounts[room.ID]))
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"rooms": response,
		},
	})
}

func handleGetStudyRoomDetail(c *gin.Context) {
	roomIDParam := c.Param("roomId")
	roomID, err := strconv.ParseUint(roomIDParam, 10, 64)
	if err != nil || roomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "房间ID不正确",
		})
		return
	}

	db := database.GetDB()
	var room models.StudyRoom
	if err := db.First(&room, roomID).Error; err != nil {
		if errorsIsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "房间不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取房间详情失败",
		})
		return
	}

	memberCount, err := countStudyRoomMembers(db, room.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取房间成员信息失败",
		})
		return
	}

	response := buildStudyRoomListItem(&room, memberCount)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"room": response,
		},
	})
}

func handleCreateStudyRoom(c *gin.Context) {
	var req createStudyRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数格式不正确",
		})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "房间名称不能为空",
		})
		return
	}
	if len([]rune(req.Name)) < 2 || len([]rune(req.Name)) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "房间名称长度需在2-20个字符之间",
		})
		return
	}

	description := strings.TrimSpace(req.Description)
	if len([]rune(description)) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "房间描述不能超过100个字符",
		})
		return
	}

	roomType := strings.ToLower(strings.TrimSpace(req.RoomType))
	if roomType != "private" {
		roomType = "public"
	}

	password := strings.TrimSpace(req.Password)
	if roomType == "private" {
		if password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "私密房间需要设置密码",
			})
			return
		}
		if len(password) < 6 || len(password) > 12 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "房间密码长度需在6-12位之间",
			})
			return
		}
	} else {
		password = ""
	}

	maxMembers := req.MaxMembers
	if maxMembers < 0 {
		maxMembers = 0
	}

	tags := normalizeStudyRoomTags(req.Tags)
	ownerID := req.OwnerUserID
	if ownerID == 0 {
		ownerID = 1
	}

	room := models.StudyRoom{
		Name:              req.Name,
		OwnerUserID:       ownerID,
		TeamID:            req.TeamID,
		Description:       description,
		Tags:              strings.Join(tags, ","),
		MaxMembers:        maxMembers,
		IsPrivate:         roomType == "private",
		Status:            1,
		FocusMinutesToday: 0,
	}
	if password != "" {
		room.AccessCode = hashPassword(password)
	}

	db := database.GetDB()
	if err := db.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建学习房间失败",
		})
		return
	}

	response := buildStudyRoomListItem(&room, 0)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "房间创建成功",
		"data": gin.H{
			"room": response,
		},
	})
}

func loadStudyRoomMemberCounts(db *gorm.DB) (map[uint64]int64, error) {
	type memberAggregate struct {
		RoomID uint64
		Total  int64
	}
	var aggregates []memberAggregate
	if err := db.Model(&models.StudyRoomMember{}).
		Select("room_id, COUNT(*) AS total").
		Group("room_id").
		Find(&aggregates).Error; err != nil {
		return nil, err
	}
	result := make(map[uint64]int64, len(aggregates))
	for _, agg := range aggregates {
		result[agg.RoomID] = agg.Total
	}
	return result, nil
}

func countStudyRoomMembers(db *gorm.DB, roomID uint64) (int64, error) {
	var total int64
	if err := db.Model(&models.StudyRoomMember{}).
		Where("room_id = ?", roomID).
		Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func buildStudyRoomListItem(room *models.StudyRoom, memberCount int64) studyRoomListItem {
	var lastStarted *string
	if room.LastSessionStarted != nil {
		formatted := room.LastSessionStarted.Format(time.RFC3339)
		lastStarted = &formatted
	}

	return studyRoomListItem{
		ID:                 room.ID,
		Name:               room.Name,
		Description:        room.Description,
		Tags:               parseStudyRoomTags(room.Tags),
		Status:             studyRoomStatusLabel(room.Status),
		StatusCode:         room.Status,
		StatusClass:        studyRoomStatusClass(room.Status),
		RoomType:           studyRoomType(room.IsPrivate),
		IsPrivate:          room.IsPrivate,
		CurrentUsers:       int(memberCount),
		MaxUsers:           room.MaxMembers,
		OwnerUserID:        room.OwnerUserID,
		TeamID:             room.TeamID,
		StudyTime:          formatStudyDuration(room.FocusMinutesToday),
		FocusMinutesToday:  room.FocusMinutesToday,
		LastSessionStarted: lastStarted,
	}
}

func normalizeStudyRoomTags(input []string) []string {
	result := make([]string, 0, len(input))
	seen := make(map[string]struct{}, len(input))
	for _, tag := range input {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			continue
		}
		lower := strings.ToLower(trimmed)
		if _, exists := seen[lower]; exists {
			continue
		}
		result = append(result, trimmed)
		seen[lower] = struct{}{}
		if len(result) >= 3 {
			break
		}
	}
	return result
}

func parseStudyRoomTags(tagString string) []string {
	if strings.TrimSpace(tagString) == "" {
		return []string{}
	}
	raw := strings.Split(tagString, ",")
	return normalizeStudyRoomTags(raw)
}

func studyRoomStatusLabel(status int8) string {
	switch status {
	case 2:
		return "待开始"
	case 3:
		return "已结束"
	default:
		return "进行中"
	}
}

func studyRoomStatusClass(status int8) string {
	switch status {
	case 2:
		return "bg-yellow-100 text-yellow-800"
	case 3:
		return "bg-gray-100 text-gray-800"
	default:
		return "bg-green-100 text-green-800"
	}
}

func studyRoomType(isPrivate bool) string {
	if isPrivate {
		return "private"
	}
	return "public"
}

func formatStudyDuration(minutes int) string {
	if minutes <= 0 {
		return "0h"
	}
	hours := minutes / 60
	remain := minutes % 60
	if hours == 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	if remain == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, remain)
}

func handleJoinStudyRoom(c *gin.Context) {
	roomIDParam := c.Param("roomId")
	roomID, err := strconv.ParseUint(roomIDParam, 10, 64)
	if err != nil || roomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "房间ID不正确"})
		return
	}

	var payload struct {
		Password string `json:"password"`
		UserID   uint64 `json:"user_id"`
	}
	_ = c.ShouldBindJSON(&payload)

	db := database.GetDB()
	var room models.StudyRoom
	if err := db.First(&room, roomID).Error; err != nil {
		status := http.StatusInternalServerError
		msg := "加载房间失败"
		if errorsIsNotFound(err) {
			status = http.StatusNotFound
			msg = "房间不存在"
		}
		c.JSON(status, gin.H{"code": status, "message": msg})
		return
	}

	if room.IsPrivate {
		if strings.TrimSpace(payload.Password) == "" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "私密房间需要密码"})
			return
		}
		if !verifyPassword(room.AccessCode, payload.Password) {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "房间密码错误"})
			return
		}
	}

	if payload.UserID != 0 {
		member := models.StudyRoomMember{
			RoomID: room.ID,
			UserID: payload.UserID,
		}
		_ = db.Where("room_id = ? AND user_id = ?", room.ID, payload.UserID).
			FirstOrCreate(&member).Error
	}

	online := studyHubRegistry.getHub(room.ID)
	memberCount, _ := countStudyRoomMembers(db, room.ID)
	response := buildStudyRoomListItem(&room, memberCount)
	response.CurrentUsers = int(memberCount)
	if online != nil {
		response.CurrentUsers = int(memberCount) + online.currentOnline()
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "加入成功",
		"data": gin.H{
			"room": response,
		},
	})
}

func handleStudySummary(c *gin.Context) {
	db := database.GetDB()

	var activeCount int64
	db.Model(&models.StudyRoom{}).Where("status = ?", 1).Count(&activeCount)

	todayStart := time.Now().Truncate(24 * time.Hour)
	var totalMinutes int64
	db.Model(&models.StudyRoom{}).
		Where("updated_at >= ?", todayStart).
		Select("SUM(focus_minutes_today)").Scan(&totalMinutes)
	// 加上当前在线未结束会话的时长
	activeMins := studyHubRegistry.activeMinutes()
	totalMinutes += int64(activeMins)

	userIDParam := c.Query("user_id")
	var streak int
	if userIDParam != "" {
		if userID, err := strconv.ParseUint(userIDParam, 10, 64); err == nil {
			var profile models.UserProfile
			if err := db.Where("user_id = ?", userID).First(&profile).Error; err == nil {
				streak = profile.StreakDays
			}
		}
	}
	if streak == 0 {
		streak = 1
	}

	summary := studySummaryResponse{
		ActiveRooms:     int(activeCount),
		OnlineUsers:     studyHubRegistry.totalOnline(),
		TodayStudyHours: float64(totalMinutes) / 60.0,
		StreakDays:      streak,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    summary,
	})
}

func handleStudyRecords(c *gin.Context) {
	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少用户ID"})
		return
	}
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户ID不正确"})
		return
	}

	db := database.GetDB()
	var records []models.LearningRecord
	if err := db.Where("user_id = ?", userID).
		Order("session_end DESC").
		Limit(20).
		Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取学习记录失败"})
		return
	}

	payload := make([]map[string]interface{}, 0, len(records))
	for _, rec := range records {
		payload = append(payload, map[string]interface{}{
			"id":          rec.ID,
			"title":       fmt.Sprintf("任务 %d", rec.TaskID),
			"duration":    rec.DurationMinutes,
			"recorded_at": rec.SessionStart.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"records": payload,
		},
	})
}
