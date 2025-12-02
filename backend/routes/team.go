package routes

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/models"
)

func registerTeamRoutes(r *gin.RouterGroup) {
	r.Use(middleware.AuthMiddleware())

	r.GET("", listTeams)
	r.GET("/", listTeams)
	r.POST("", createTeam)
	r.POST("/", createTeam)
	r.POST("/join_by_name", joinTeamByName)
}

func createTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Visibility  *int8  `json:"visibility"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "团队名称不能为空"})
		return
	}

	desc := strings.TrimSpace(req.Description)
	visibility := int8(1)
	if req.Visibility != nil {
		visibility = *req.Visibility
	}

	team := models.Team{
		Name:        name,
		Description: desc,
		OwnerUserID: userID.(uint64),
		Visibility:  visibility,
	}

	db := database.GetDB()
	if err := db.Create(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "团队名称已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建团队失败"})
		return
	}

	// Ensure owner is part of the team members list
	ownerMember := models.TeamMember{
		TeamID: team.ID,
		UserID: team.OwnerUserID,
		Role:   1,
	}
	if err := db.Where("team_id = ? AND user_id = ?", team.ID, team.OwnerUserID).FirstOrCreate(&ownerMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建团队成员关系失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": team})
}

func joinTeamByName(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	var team models.Team
	if err := database.GetDB().Where("name = ?", req.Name).First(&team).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该名称的团队"})
		return
	}

	// Check if owner
	if team.OwnerUserID == uid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您是该团队的创建者"})
		return
	}

	// Check if already a member
	var count int64
	database.GetDB().Model(&models.TeamMember{}).Where("team_id = ? AND user_id = ?", team.ID, uid).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已经是该团队的成员"})
		return
	}

	member := models.TeamMember{
		TeamID: team.ID,
		UserID: uid,
		Role:   0, // Member
	}

	if err := database.GetDB().Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加入团队失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "加入成功", "data": team})
}

func listTeams(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var teams []struct {
		models.Team
		OwnerName   string `json:"owner_name"`
		MemberCount int64  `json:"member_count"`
	}

	db := database.GetDB().Table("teams").
		Select("teams.*, users.display_name as owner_name, (SELECT COUNT(*) FROM team_members WHERE team_members.team_id = teams.id) as member_count").
		Joins("LEFT JOIN users ON users.id = teams.owner_user_id")

	ownedOnly := strings.TrimSpace(strings.ToLower(c.DefaultQuery("owned_only", "true")))
	if ownedOnly == "true" || ownedOnly == "1" || ownedOnly == "" {
		db = db.Where("teams.owner_user_id = ?", userID.(uint64))
	} else {
		db = db.Where("teams.owner_user_id = ? OR teams.id IN (SELECT team_id FROM team_members WHERE user_id = ?)", userID.(uint64), userID.(uint64))
	}

	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("teams.name LIKE ? OR teams.description LIKE ?", like, like)
	}

	if err := db.Order("teams.created_at DESC").Scan(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取团队列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": teams,
		"msg":  "获取成功",
	})
}
