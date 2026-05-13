package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

const roomKindTeamChat = "team_chat"

func registerTeamChatRoutes(r *gin.RouterGroup) {
	r.POST("/:id/chat-room", ensureTeamChatRoom)
	r.GET("/:id/chat-room", ensureTeamChatRoom)
}

func ensureTeamChatRoom(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || teamID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	userID := userIDValue.(uint64)
	db := database.GetDB()
	if !canAccessTeam(db, teamID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "您不是该团队成员"})
		return
	}

	var room models.StudyRoom
	err = db.Where("team_id = ? AND room_kind = ?", teamID, roomKindTeamChat).
		Order("id ASC").
		First(&room).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取团队聊天室失败"})
		return
	}
	if err == gorm.ErrRecordNotFound {
		var team models.Team
		if err := db.First(&team, teamID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "团队不存在"})
			return
		}
		room = models.StudyRoom{
			Name:        fmt.Sprintf("%s 团队聊天室", team.Name),
			OwnerUserID: team.OwnerUserID,
			TeamID:      &team.ID,
			Description: "团队持久沟通聊天室",
			Tags:        "team-chat",
			RoomKind:    roomKindTeamChat,
			IsPrivate:   true,
			Status:      1,
		}
		if err := db.Create(&room).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建团队聊天室失败"})
			return
		}
	}
	ensureTeamChatMembers(db, &room)
	memberCount, _ := countStudyRoomMembers(db, room.ID)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"room": buildStudyRoomListItem(&room, memberCount)}})
}

func canAccessTeam(db *gorm.DB, teamID, userID uint64) bool {
	var count int64
	if err := db.Model(&models.TeamMember{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Count(&count).Error; err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	var team models.Team
	return db.Select("id, owner_user_id").First(&team, teamID).Error == nil && team.OwnerUserID == userID
}

func ensureTeamChatMembers(db *gorm.DB, room *models.StudyRoom) {
	if room.TeamID == nil {
		return
	}
	var members []models.TeamMember
	_ = db.Where("team_id = ?", *room.TeamID).Find(&members).Error
	for _, teamMember := range members {
		role := int8(0)
		if teamMember.Role == 1 || teamMember.UserID == room.OwnerUserID {
			role = 1
		}
		member := models.StudyRoomMember{
			RoomID:   room.ID,
			UserID:   teamMember.UserID,
			Role:     role,
			JoinedAt: time.Now(),
		}
		_ = db.Where("room_id = ? AND user_id = ?", room.ID, teamMember.UserID).FirstOrCreate(&member).Error
	}
}

func getTeamChatRoomByRoom(db *gorm.DB, roomID uint64) (models.StudyRoom, bool) {
	var room models.StudyRoom
	if err := db.Where("id = ? AND room_kind = ?", roomID, roomKindTeamChat).First(&room).Error; err != nil {
		return models.StudyRoom{}, false
	}
	return room, true
}

func canAccessTeamChatRoom(db *gorm.DB, roomID, userID uint64) bool {
	room, ok := getTeamChatRoomByRoom(db, roomID)
	if !ok || room.TeamID == nil {
		return false
	}
	return canAccessTeam(db, *room.TeamID, userID)
}
