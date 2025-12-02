package routes

import (
	"learningAssistant-backend/database"
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerTeamRoutes(rg *gin.RouterGroup) {
	rg.Use(middleware.AuthMiddleware())
	rg.GET("/", handleGetMyTeams)
	rg.GET("", handleGetMyTeams)
	rg.POST("/", handleCreateTeam)
	rg.POST("", handleCreateTeam)
	rg.POST("/join-by-name", handleJoinTeamByName)
}

// handleCreateTeam 创建团队
func handleCreateTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if team name exists
	var count int64
	database.DB.Model(&models.Team{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team name already exists"})
		return
	}

	team := models.Team{
		Name:        req.Name,
		Description: req.Description,
		OwnerUserID: userID.(uint64),
		Visibility:  1, // Default public
	}

	tx := database.DB.Begin()

	if err := tx.Create(&team).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
		return
	}

	member := models.TeamMember{
		TeamID: team.ID,
		UserID: userID.(uint64),
		Role:   1, // Owner/Admin role
	}

	if err := tx.Create(&member).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add owner to team members"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, team)
}

// handleGetMyTeams 获取当前用户加入的团队列表
func handleGetMyTeams(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var members []models.TeamMember
	if err := database.DB.Where("user_id = ?", userID).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team memberships"})
		return
	}

	if len(members) == 0 {
		c.JSON(http.StatusOK, []models.Team{})
		return
	}

	var teamIDs []uint64
	for _, m := range members {
		teamIDs = append(teamIDs, m.TeamID)
	}

	var teams []models.Team
	if err := database.DB.Where("id IN ?", teamIDs).Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// handleJoinTeamByName 通过名称加入团队
func handleJoinTeamByName(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		TeamName string `json:"team_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team name is required"})
		return
	}

	var team models.Team
	if err := database.DB.Where("name = ?", req.TeamName).First(&team).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Check if already a member
	var count int64
	database.DB.Model(&models.TeamMember{}).Where("team_id = ? AND user_id = ?", team.ID, userID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already a member of this team"})
		return
	}

	member := models.TeamMember{
		TeamID: team.ID,
		UserID: userID.(uint64),
		Role:   0, // Default role
	}

	if err := database.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join team"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined team successfully", "team": team})
}
