package routes

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

func handleGetRoomChatHistory(c *gin.Context) {
	roomIDParam := c.Param("roomId")
	roomID, err := strconv.ParseUint(roomIDParam, 10, 64)
	if err != nil || roomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "房间ID不正确"})
		return
	}
	limitParam := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 100
	}
	if limit > 200 {
		limit = 200
	}

	db := database.GetDB()
	session, hasCollaborationSession := getCollaborationSessionByRoom(db, roomID)
	if teamRoom, hasTeamRoom := getTeamChatRoomByRoom(db, roomID); hasTeamRoom {
		userID, ok := currentUserID(c)
		if !ok || teamRoom.TeamID == nil || !canAccessTeam(db, *teamRoom.TeamID, userID) {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权限访问团队聊天室"})
			return
		}
	}
	var chats []models.ChatMessage
	query := db.Where("room_id = ?", roomID)
	if hasCollaborationSession {
		query = query.Where("session_id = ?", session.ID)
	}
	if err := query.
		Order("sent_at DESC").
		Limit(limit).
		Find(&chats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取聊天记录失败"})
		return
	}

	userIDs := make([]uint64, 0, len(chats))
	for _, msg := range chats {
		userIDs = append(userIDs, msg.UserID)
	}
	nameMap := loadUserNames(userIDs)

	result := make([]map[string]interface{}, 0, len(chats))
	for _, msg := range chats {
		result = append(result, buildChatMessageResponse(msg, nameMap[msg.UserID]))
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"messages": result},
	})
}

func handlePostRoomChat(c *gin.Context) {
	roomIDParam := c.Param("roomId")
	roomID, err := strconv.ParseUint(roomIDParam, 10, 64)
	if err != nil || roomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "房间ID不正确"})
		return
	}

	var req struct {
		UserID  uint64 `json:"user_id"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数格式不正确"})
		return
	}
	req.Content = strings.TrimSpace(req.Content)
	if req.UserID == 0 {
		if uid, ok := currentUserID(c); ok {
			req.UserID = uid
		}
	}
	if req.UserID == 0 || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户ID或内容不能为空"})
		return
	}

	db := database.GetDB()
	session, hasCollaborationSession := getCollaborationSessionByRoom(db, roomID)
	sessionID := uint64(0)
	if teamRoom, hasTeamRoom := getTeamChatRoomByRoom(db, roomID); hasTeamRoom {
		if teamRoom.TeamID == nil || !canAccessTeam(db, *teamRoom.TeamID, req.UserID) {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权限访问团队聊天室"})
			return
		}
	}
	if hasCollaborationSession {
		if session.Status == models.TaskCollaborationStatusDismissed {
			c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "协作会话已结束"})
			return
		}
		if !canAccessCollaborationSession(db, &session, req.UserID) {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权限访问协作会话"})
			return
		}
		sessionID = session.ID
	}
	chat := models.ChatMessage{
		SessionID: sessionID,
		RoomID:    roomID,
		UserID:    req.UserID,
		Content:   req.Content,
		MsgType:   models.ChatMessageTypeText,
		SentAt:    time.Now(),
	}
	if err := db.Create(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存消息失败"})
		return
	}

	name := loadUserNames([]uint64{req.UserID})[req.UserID]
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"message": buildChatMessageResponse(chat, name),
		},
	})
}

func loadUserNames(ids []uint64) map[uint64]string {
	result := make(map[uint64]string)
	if len(ids) == 0 {
		return result
	}
	db := database.GetDB()
	var users []models.User
	_ = db.Where("id IN ?", ids).Select("id, display_name").Find(&users).Error
	for _, u := range users {
		result[u.ID] = u.DisplayName
	}
	return result
}
