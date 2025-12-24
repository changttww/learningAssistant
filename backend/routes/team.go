package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	r.POST("/:id/invite", inviteMember)
	r.GET("/:id/members", listTeamMembers)
	r.GET("/:id/activities", listTeamActivities)
	r.GET("/:id/requests", listTeamRequests)
	r.POST("/:id/requests/:requestId/handle", handleTeamRequest)
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

	// Check if already applied
	var reqCount int64
	database.GetDB().Model(&models.TeamRequest{}).Where("team_id = ? AND user_id = ? AND status IN ('PENDING_USER', 'PENDING_OWNER')", team.ID, uid).Count(&reqCount)
	if reqCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已经申请加入该团队或已被邀请"})
		return
	}

	// Create application request
	request := models.TeamRequest{
		TeamID:    team.ID,
		UserID:    uid,
		InviterID: 0,
		Type:      "APPLICATION",
		Status:    "PENDING_OWNER",
	}

	if err := database.GetDB().Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "申请加入失败"})
		return
	}

	// Notify Owner with detailed info
	var applicant models.User
	database.GetDB().First(&applicant, uid)

	// Create related data JSON
	relatedData := fmt.Sprintf(`{"applicant_name":"%s","team_name":"%s","team_id":%d,"application_type":"direct"}`, applicant.DisplayName, team.Name, team.ID)

	notification := models.Notification{
		UserID:      team.OwnerUserID,
		Title:       "新的入队申请",
		Content:     applicant.DisplayName + " 申请加入团队: " + team.Name,
		Type:        "TEAM_APPLICATION",
		RelatedID:   request.ID,
		RelatedData: relatedData,
	}
	database.GetDB().Create(&notification)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "申请已发送，等待团队创建者审核"})
}

func inviteMember(c *gin.Context) {
	teamID := c.Param("id")
	var req struct {
		Account string `json:"account" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	// Check permission (must be member)
	var member models.TeamMember
	if err := database.GetDB().Where("team_id = ? AND user_id = ?", teamID, uid).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "您不是该团队成员"})
		return
	}

	// Find target user
	var targetUser models.User
	if err := database.GetDB().Where("account = ?", req.Account).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该用户"})
		return
	}

	// Check if target is already member
	var count int64
	database.GetDB().Model(&models.TeamMember{}).Where("team_id = ? AND user_id = ?", teamID, targetUser.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该用户已经是团队成员"})
		return
	}

	// Check if already invited/applied
	var reqCount int64
	database.GetDB().Model(&models.TeamRequest{}).Where("team_id = ? AND user_id = ? AND status IN ('PENDING_USER', 'PENDING_OWNER')", teamID, targetUser.ID).Count(&reqCount)
	if reqCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该用户已被邀请或已申请"})
		return
	}

	// Create invitation request
	request := models.TeamRequest{
		TeamID:    member.TeamID,
		UserID:    targetUser.ID,
		InviterID: uid,
		Type:      "INVITATION",
		Status:    "PENDING_USER",
	}

	if err := database.GetDB().Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "邀请失败"})
		return
	}

	// Notify Target User with detailed info
	var team models.Team
	database.GetDB().First(&team, teamID)

	// Get inviter's display name
	var inviter models.User
	database.GetDB().First(&inviter, uid)

	// Create related data JSON
	relatedData := fmt.Sprintf(`{"inviter_name":"%s","team_name":"%s","team_id":%d}`, inviter.DisplayName, team.Name, team.ID)

	notification := models.Notification{
		UserID:      targetUser.ID,
		Title:       "团队邀请",
		Content:     inviter.DisplayName + " 邀请您加入团队: " + team.Name,
		Type:        "TEAM_INVITE",
		RelatedID:   request.ID,
		RelatedData: relatedData,
	}
	database.GetDB().Create(&notification)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "邀请已发送"})
}

func listTeamRequests(c *gin.Context) {
	teamID := c.Param("id")
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	// Check if owner
	var team models.Team
	if err := database.GetDB().First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "团队不存在"})
		return
	}
	if team.OwnerUserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有团队创建者可以查看申请列表"})
		return
	}

	var requests []struct {
		models.TeamRequest
		UserName    string `json:"user_name"`
		InviterName string `json:"inviter_name"`
	}

	db := database.GetDB().Table("team_requests").
		Select("team_requests.*, u1.display_name as user_name, u2.display_name as inviter_name").
		Joins("LEFT JOIN users u1 ON u1.id = team_requests.user_id").
		Joins("LEFT JOIN users u2 ON u2.id = team_requests.inviter_id").
		Where("team_requests.team_id = ? AND team_requests.status = 'PENDING_OWNER'", teamID)

	if err := db.Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取申请列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": requests})
}

func handleTeamRequest(c *gin.Context) {
	requestID := c.Param("requestId")
	var req struct {
		Action string `json:"action" binding:"required"` // ACCEPT, REJECT
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	var request models.TeamRequest
	if err := database.GetDB().First(&request, requestID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "请求不存在"})
		return
	}

	if request.Status != "PENDING_USER" && request.Status != "PENDING_OWNER" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求状态无效"})
		return
	}

	var team models.Team
	database.GetDB().First(&team, request.TeamID)

	// Logic for handling
	if request.Status == "PENDING_USER" {
		// User handling invitation
		if request.UserID != uid {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权处理此请求"})
			return
		}
		if req.Action == "REJECT" {
			request.Status = "REJECTED"
			database.GetDB().Save(&request)
			database.GetDB().Model(&models.Notification{}).Where("user_id = ? AND related_id = ? AND type IN ?", uid, request.ID, []string{"TEAM_INVITE", "TEAM_APPLICATION"}).Updates(map[string]interface{}{"action_status": "REJECTED", "is_read": true})
			// Notify Inviter? Maybe later.
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "已拒绝"})
			return
		}

		// If Accepted
		// Check if inviter is owner
		if request.InviterID == team.OwnerUserID {
			// Direct join
			member := models.TeamMember{TeamID: team.ID, UserID: uid, Role: 0}
			database.GetDB().Create(&member)
			request.Status = "APPROVED"
			database.GetDB().Save(&request)
			database.GetDB().Model(&models.Notification{}).Where("user_id = ? AND related_id = ? AND type IN ?", uid, request.ID, []string{"TEAM_INVITE", "TEAM_APPLICATION"}).Updates(map[string]interface{}{"action_status": "ACCEPTED", "is_read": true})
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "已加入团队"})
		} else {
			// Needs owner approval
			request.Status = "PENDING_OWNER"
			database.GetDB().Save(&request)
			// Notify Owner with detailed info
			var invitee models.User
			database.GetDB().First(&invitee, request.UserID)

			var inviter models.User
			database.GetDB().First(&inviter, request.InviterID)

			// Create related data JSON
			relatedData := fmt.Sprintf(`{"invitee_name":"%s","inviter_name":"%s","team_name":"%s","team_id":%d,"application_type":"invitation"}`, invitee.DisplayName, inviter.DisplayName, team.Name, team.ID)

			notification := models.Notification{
				UserID:      team.OwnerUserID,
				Title:       "新的入队申请",
				Content:     invitee.DisplayName + " 接受了 " + inviter.DisplayName + " 的邀请，等待您的批准: " + team.Name,
				Type:        "TEAM_APPLICATION",
				RelatedID:   request.ID,
				RelatedData: relatedData,
			}
			database.GetDB().Create(&notification)
			database.GetDB().Model(&models.Notification{}).Where("user_id = ? AND related_id = ? AND type IN ?", uid, request.ID, []string{"TEAM_INVITE", "TEAM_APPLICATION"}).Updates(map[string]interface{}{"action_status": "ACCEPTED", "is_read": true})
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "已接受邀请，等待团队队长审核"})
		}
	} else if request.Status == "PENDING_OWNER" {
		// Owner handling application
		if team.OwnerUserID != uid {
			c.JSON(http.StatusForbidden, gin.H{"error": "只有群主可以处理此请求"})
			return
		}

		if req.Action == "REJECT" {
			request.Status = "REJECTED"
			database.GetDB().Save(&request)
			// Notify User with detailed info
			var rejecter models.User
			database.GetDB().First(&rejecter, uid)

			// Create related data JSON
			relatedData := fmt.Sprintf(`{"rejecter_name":"%s","team_name":"%s","team_id":%d}`, rejecter.DisplayName, team.Name, team.ID)

			notification := models.Notification{
				UserID:      request.UserID,
				Title:       "申请被拒绝",
				Content:     "您的入队申请被 " + rejecter.DisplayName + " 拒绝，您加入团队 " + team.Name + " 的申请未通过",
				Type:        "SYSTEM",
				RelatedID:   request.ID,
				RelatedData: relatedData,
			}
			database.GetDB().Create(&notification)
			database.GetDB().Model(&models.Notification{}).Where("user_id = ? AND related_id = ? AND type IN ?", uid, request.ID, []string{"TEAM_INVITE", "TEAM_APPLICATION"}).Updates(map[string]interface{}{"action_status": "REJECTED", "is_read": true})
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "已拒绝"})
			return
		}

		// Approve
		member := models.TeamMember{TeamID: team.ID, UserID: request.UserID, Role: 0}
		database.GetDB().Create(&member)
		request.Status = "APPROVED"
		database.GetDB().Save(&request)

		// Notify User with detailed info
		var approver models.User
		database.GetDB().First(&approver, uid)

		// Create related data JSON
		relatedData := fmt.Sprintf(`{"approver_name":"%s","team_name":"%s","team_id":%d}`, approver.DisplayName, team.Name, team.ID)

		notification := models.Notification{
			UserID:      request.UserID,
			Title:       "申请通过",
			Content:     "您的入队申请已被 " + approver.DisplayName + " 批准，您已成功加入团队: " + team.Name,
			Type:        "SYSTEM",
			RelatedID:   request.ID,
			RelatedData: relatedData,
		}
		database.GetDB().Create(&notification)
		database.GetDB().Model(&models.Notification{}).Where("user_id = ? AND related_id = ? AND type IN ?", uid, request.ID, []string{"TEAM_INVITE", "TEAM_APPLICATION"}).Updates(map[string]interface{}{"action_status": "APPROVED", "is_read": true})
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "已批准"})
	}
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

func listTeamMembers(c *gin.Context) {
	teamID := c.Param("id")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	var members []struct {
		UserID      uint64 `json:"user_id"`
		Account     string `json:"account"`
		Nickname    string `json:"nickname"`
		Avatar      string `json:"avatar"`
		Role        int8   `json:"role"`
		TotalPoints int    `json:"total_points"`
	}

	db := database.GetDB()
	err := db.Table("team_members").
		Select("team_members.user_id, users.account, users.display_name as nickname, users.avatar_url as avatar, team_members.role, COALESCE(user_profiles.total_points, 0) as total_points").
		Joins("JOIN users ON users.id = team_members.user_id").
		Joins("LEFT JOIN user_profiles ON user_profiles.user_id = users.id").
		Where("team_members.team_id = ?", teamID).
		Scan(&members).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team members"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": members})
}

func listTeamActivities(c *gin.Context) {
	teamID := c.Param("id")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	type Activity struct {
		UserName   string    `json:"user_name"`
		UserAvatar string    `json:"user_avatar"`
		Action     string    `json:"action"`
		TaskTitle  string    `json:"task_title"`
		Time       time.Time `json:"time"`
	}

	var activities []Activity
	var tasks []models.Task

	// Fetch recent tasks for the team
	db := database.GetDB()
	if err := db.Where("owner_team_id = ?", teamID).Order("updated_at desc").Limit(20).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team activities"})
		return
	}

	for _, task := range tasks {
		var user models.User
		// Use CreatedBy for now.
		userId := task.CreatedBy

		db.Select("display_name, avatar_url").First(&user, userId)

		action := "更新了任务"
		if task.Status == 2 {
			action = "完成了任务"
		} else if task.UpdatedAt.Sub(task.CreatedAt) < 5*time.Second {
			action = "创建了任务"
		}

		activities = append(activities, Activity{
			UserName:   user.DisplayName,
			UserAvatar: user.AvatarURL,
			Action:     action,
			TaskTitle:  task.Title,
			Time:       task.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": activities})
}
