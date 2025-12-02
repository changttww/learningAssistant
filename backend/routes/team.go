package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"learningAssistant-backend/database"
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/models"
)

func registerTeamRoutes(r *gin.RouterGroup) {
	r.Use(middleware.AuthMiddleware())

	r.GET("", listTeams)
	r.GET("/", listTeams)
}

func listTeams(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var teams []models.Team
	db := database.GetDB().Model(&models.Team{})

	ownedOnly := strings.TrimSpace(strings.ToLower(c.DefaultQuery("owned_only", "true")))
	if ownedOnly == "true" || ownedOnly == "1" || ownedOnly == "" {
		db = db.Where("owner_user_id = ?", userID.(uint64))
	} else {
		db = db.Where("owner_user_id = ? OR id IN (SELECT team_id FROM team_members WHERE user_id = ?)", userID.(uint64), userID.(uint64))
	}

	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", like, like)
	}

	if err := db.Order("created_at DESC").Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取团队列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": teams,
		"msg":  "获取成功",
	})
}
