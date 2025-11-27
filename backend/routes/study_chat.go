package routes

import (
	"net/http"
	"sort"
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
	var chats []models.ChatMessage
	if err := db.Where("room_id = ?", roomID).
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
		result = append(result, map[string]interface{}{
			"id":           msg.ID,
			"user_id":      msg.UserID,
			"display_name": nameMap[msg.UserID],
			"content":      msg.Content,
			"sent_at":      msg.SentAt,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i]["sent_at"].(time.Time).Before(result[j]["sent_at"].(time.Time))
	})

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
	if req.UserID == 0 || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户ID或内容不能为空"})
		return
	}

	db := database.GetDB()
	chat := models.ChatMessage{
		SessionID: 0,
		RoomID:    roomID,
		UserID:    req.UserID,
		Content:   req.Content,
		MsgType:   0,
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
			"message": map[string]interface{}{
				"id":           chat.ID,
				"user_id":      chat.UserID,
				"display_name": name,
				"content":      chat.Content,
				"sent_at":      chat.SentAt,
			},
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
