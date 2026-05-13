package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

type collaborationSessionResponse struct {
	ID           uint64                         `json:"id"`
	TaskID       uint64                         `json:"task_id"`
	TeamID       *uint64                        `json:"team_id"`
	RoomID       uint64                         `json:"room_id"`
	Status       int8                           `json:"status"`
	StatusLabel  string                         `json:"status_label"`
	IsActive     bool                           `json:"is_active"`
	CreatedBy    uint64                         `json:"created_by"`
	DismissedAt  *time.Time                     `json:"dismissed_at"`
	Minutes      interface{}                    `json:"minutes,omitempty"`
	Participants []collaborationParticipantView `json:"participants,omitempty"`
	CreatedAt    time.Time                      `json:"created_at"`
}

type collaborationParticipantView struct {
	UserID      uint64 `json:"user_id"`
	DisplayName string `json:"display_name"`
	Role        int8   `json:"role"`
}

type knowledgeCardPayload struct {
	KnowledgeEntryID uint64   `json:"knowledge_entry_id"`
	Title            string   `json:"title"`
	Summary          string   `json:"summary"`
	Excerpt          string   `json:"excerpt"`
	Category         string   `json:"category"`
	SourceUserID     uint64   `json:"source_user_id"`
	Tags             []string `json:"tags"`
}

type minutesPayload struct {
	Summary               string              `json:"summary"`
	SynchronizedKnowledge []string            `json:"synchronized_knowledge"`
	ActionItems           []minutesActionItem `json:"action_items"`
	Blockers              []string            `json:"blockers"`
	NextSteps             []string            `json:"next_steps"`
	GeneratedAt           time.Time           `json:"generated_at"`
}

type minutesActionItem struct {
	Owner  string `json:"owner"`
	Action string `json:"action"`
}

func registerTaskCollaborationRoutes(r *gin.RouterGroup) {
	r.POST("/:id/collaboration-sessions", createTaskCollaborationSession)
	r.GET("/:id/collaboration-sessions", listTaskCollaborationSessions)
	r.GET("/collaboration-sessions/:sessionId", getTaskCollaborationSession)
	r.POST("/collaboration-sessions/:sessionId/dismiss", dismissTaskCollaborationSession)
	r.POST("/collaboration-sessions/:sessionId/knowledge-cards", sendTaskCollaborationKnowledgeCard)
	r.GET("/collaboration-sessions/:sessionId/knowledge-cards/:knowledgeEntryId", getTaskCollaborationKnowledgeCard)
	r.POST("/collaboration-sessions/:sessionId/minutes/generate", generateTaskCollaborationMinutes)
	r.POST("/collaboration-sessions/:sessionId/minutes/save", saveTaskCollaborationMinutes)
	r.POST("/collaboration-sessions/:sessionId/team-knowledge", saveTaskCollaborationTeamKnowledge)
}

func createTaskCollaborationSession(c *gin.Context) {
	task, userID, ok := loadTaskForCollaboration(c)
	if !ok {
		return
	}

	db := database.GetDB()
	normalizeTaskCollaborationLifecycle(db, task.ID)
	var active models.TaskCollaborationSession
	if err := db.Where("task_id = ? AND status = ?", task.ID, models.TaskCollaborationStatusActive).
		Order("id DESC").
		First(&active).Error; err == nil {
		dismissOlderActiveCollaborationSessions(db, task.ID, active.ID)
		ensureCollaborationParticipants(db, &active, &task, userID)
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"session": buildCollaborationSessionResponse(&active)}})
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询协作房间失败"})
		return
	}

	room := models.StudyRoom{
		Name:        "任务协作：" + task.Title,
		OwnerUserID: userID,
		TeamID:      task.OwnerTeamID,
		Description: "从团队任务发起的持久协作房间",
		Tags:        "task-collaboration",
		RoomKind:    "task_collaboration",
		IsPrivate:   true,
		Status:      1,
	}
	if err := db.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建协作房间失败"})
		return
	}

	session := models.TaskCollaborationSession{
		TaskID:    task.ID,
		TeamID:    task.OwnerTeamID,
		RoomID:    room.ID,
		CreatedBy: userID,
		Status:    models.TaskCollaborationStatusActive,
	}
	if err := db.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建协作会话失败"})
		return
	}

	participants := ensureCollaborationParticipants(db, &session, &task, userID)
	createCollaborationInviteMessage(db, &session, userID, participants)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"session": buildCollaborationSessionResponse(&session)}})
}

func ensureCollaborationParticipants(db *gorm.DB, session *models.TaskCollaborationSession, task *models.Task, requesterID uint64) []uint64 {
	participants := deriveTaskCollaborationParticipants(db, task, requesterID)
	for _, uid := range participants {
		role := int8(0)
		if uid == requesterID {
			role = 1
		}
		participant := models.TaskCollaborationParticipant{SessionID: session.ID, UserID: uid, Role: role}
		_ = db.Where("session_id = ? AND user_id = ?", session.ID, uid).FirstOrCreate(&participant).Error
		member := models.StudyRoomMember{RoomID: session.RoomID, UserID: uid, Role: role, JoinedAt: time.Now()}
		_ = db.Where("room_id = ? AND user_id = ?", session.RoomID, uid).FirstOrCreate(&member).Error
	}
	return participants
}

func listTaskCollaborationSessions(c *gin.Context) {
	task, _, ok := loadTaskForCollaboration(c)
	if !ok {
		return
	}
	normalizeTaskCollaborationLifecycle(database.GetDB(), task.ID)
	var sessions []models.TaskCollaborationSession
	if err := database.GetDB().Where("task_id = ?", task.ID).Order("id DESC").Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取协作会话失败"})
		return
	}
	items := make([]collaborationSessionResponse, 0, len(sessions))
	for i := range sessions {
		items = append(items, buildCollaborationSessionResponse(&sessions[i]))
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"sessions": items}})
}

func getTaskCollaborationSession(c *gin.Context) {
	session, _, ok := loadSessionForCollaboration(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"session": buildCollaborationSessionResponse(&session)}})
}

func dismissTaskCollaborationSession(c *gin.Context) {
	session, _, ok := loadSessionForCollaboration(c)
	if !ok {
		return
	}
	now := time.Now()
	updates := map[string]interface{}{
		"status":       models.TaskCollaborationStatusDismissed,
		"dismissed_at": now,
	}
	db := database.GetDB()
	if err := db.Model(&session).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解散协作会话失败"})
		return
	}
	_ = db.Model(&models.TaskCollaborationSession{}).
		Where("task_id = ? AND status = ?", session.TaskID, models.TaskCollaborationStatusActive).
		Updates(updates).Error
	session.Status = models.TaskCollaborationStatusDismissed
	session.DismissedAt = &now
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"session": buildCollaborationSessionResponse(&session)}})
}

func dismissOlderActiveCollaborationSessions(db *gorm.DB, taskID, keepSessionID uint64) {
	now := time.Now()
	_ = db.Model(&models.TaskCollaborationSession{}).
		Where("task_id = ? AND status = ? AND id <> ?", taskID, models.TaskCollaborationStatusActive, keepSessionID).
		Updates(map[string]interface{}{
			"status":       models.TaskCollaborationStatusDismissed,
			"dismissed_at": now,
		}).Error
}

func normalizeTaskCollaborationLifecycle(db *gorm.DB, taskID uint64) {
	var latest models.TaskCollaborationSession
	if err := db.Where("task_id = ?", taskID).Order("id DESC").First(&latest).Error; err != nil {
		return
	}
	if latest.Status == models.TaskCollaborationStatusDismissed {
		now := time.Now()
		_ = db.Model(&models.TaskCollaborationSession{}).
			Where("task_id = ? AND status = ?", taskID, models.TaskCollaborationStatusActive).
			Updates(map[string]interface{}{
				"status":       models.TaskCollaborationStatusDismissed,
				"dismissed_at": now,
			}).Error
		return
	}
	if latest.Status == models.TaskCollaborationStatusActive {
		dismissOlderActiveCollaborationSessions(db, taskID, latest.ID)
	}
}

func sendTaskCollaborationKnowledgeCard(c *gin.Context) {
	session, userID, ok := loadSessionForCollaboration(c)
	if !ok {
		return
	}
	if session.Status == models.TaskCollaborationStatusDismissed {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "协作会话已结束"})
		return
	}
	var req struct {
		KnowledgeEntryID uint64 `json:"knowledge_entry_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.KnowledgeEntryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "知识点ID不能为空"})
		return
	}
	var entry models.KnowledgeBaseEntry
	if err := database.GetDB().
		Where("id = ? AND user_id = ?", req.KnowledgeEntryID, userID).
		First(&entry).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "只能共享自己的知识库内容"})
		return
	}
	payload := knowledgeCardPayload{
		KnowledgeEntryID: entry.ID,
		Title:            entry.Title,
		Summary:          entry.Summary,
		Excerpt:          firstNonEmpty(entry.Summary, truncateRunes(entry.Content, 160)),
		Category:         entry.Category,
		SourceUserID:     entry.UserID,
	}
	content, _ := json.Marshal(payload)
	chat := models.ChatMessage{
		SessionID: session.ID,
		RoomID:    session.RoomID,
		UserID:    userID,
		Content:   string(content),
		MsgType:   models.ChatMessageTypeKnowledgeCard,
		SentAt:    time.Now(),
	}
	if err := database.GetDB().Create(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存知识卡片失败"})
		return
	}
	name := loadUserNames([]uint64{userID})[userID]
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"message": buildChatMessageResponse(chat, name)}})
}

func getTaskCollaborationKnowledgeCard(c *gin.Context) {
	session, _, ok := loadSessionForCollaboration(c)
	if !ok {
		return
	}
	entryID, err := strconv.ParseUint(c.Param("knowledgeEntryId"), 10, 64)
	if err != nil || entryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "笔记ID不能为空"})
		return
	}
	var count int64
	if err := database.GetDB().Model(&models.ChatMessage{}).
		Where("session_id = ? AND room_id = ? AND msg_type = ? AND content LIKE ?", session.ID, session.RoomID, models.ChatMessageTypeKnowledgeCard, fmt.Sprintf("%%\"knowledge_entry_id\":%d%%", entryID)).
		Count(&count).Error; err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "共享笔记不存在"})
		return
	}
	var entry models.KnowledgeBaseEntry
	if err := database.GetDB().First(&entry, entryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "笔记不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"entry": gin.H{
		"id":             entry.ID,
		"title":          entry.Title,
		"summary":        entry.Summary,
		"content":        entry.Content,
		"category":       entry.Category,
		"source_user_id": entry.UserID,
	}}})
}

func generateTaskCollaborationMinutes(c *gin.Context) {
	session, _, ok := loadSessionForCollaboration(c)
	if !ok {
		return
	}
	var task models.Task
	_ = database.GetDB().First(&task, session.TaskID).Error
	var messages []models.ChatMessage
	if err := database.GetDB().Where("session_id = ? AND room_id = ? AND msg_type <> ?", session.ID, session.RoomID, models.ChatMessageTypeSystem).
		Order("sent_at ASC").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取会话消息失败"})
		return
	}
	if len(messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "当前讨论内容不足，无法生成纪要"})
		return
	}
	minutes := buildRuleBasedMinutes(&task, messages)
	if apiKey := getQwenAPIKey(); apiKey != "" {
		if aiMinutes, err := generateAIMinutes(apiKey, &task, messages); err == nil {
			minutes = aiMinutes
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"minutes": minutes}})
}

func saveTaskCollaborationMinutes(c *gin.Context) {
	session, userID, ok := loadSessionForCollaboration(c)
	if !ok {
		return
	}
	var req struct {
		Minutes json.RawMessage `json:"minutes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Minutes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "纪要内容不能为空"})
		return
	}
	hash := sha256.Sum256(req.Minutes)
	hashValue := hex.EncodeToString(hash[:])
	if session.MinutesHash == hashValue {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"duplicate": true, "session": buildCollaborationSessionResponse(&session)}})
		return
	}
	now := time.Now()
	if err := database.GetDB().Model(&session).Updates(map[string]interface{}{
		"minutes":         datatypes.JSON(req.Minutes),
		"minutes_hash":    hashValue,
		"last_minutes_at": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存纪要失败"})
		return
	}
	appendMinutesComment(session.TaskID, userID, req.Minutes)
	session.Minutes = datatypes.JSON(req.Minutes)
	session.MinutesHash = hashValue
	session.LastMinutesAt = &now
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"duplicate": false, "session": buildCollaborationSessionResponse(&session)}})
}

func saveTaskCollaborationTeamKnowledge(c *gin.Context) {
	session, userID, ok := loadSessionForCollaboration(c)
	if !ok {
		return
	}
	if session.TeamID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该任务不属于团队"})
		return
	}
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	_ = c.ShouldBindJSON(&req)
	if strings.TrimSpace(req.Title) == "" {
		req.Title = fmt.Sprintf("任务协作纪要 #%d", session.ID)
	}
	if strings.TrimSpace(req.Content) == "" {
		req.Content = string(session.Minutes)
	}
	req.Content = formatReadableMinutes([]byte(req.Content))
	if strings.TrimSpace(req.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "团队知识内容不能为空"})
		return
	}
	entry := models.KnowledgeBaseEntry{
		UserID:     userID,
		TeamID:     session.TeamID,
		SourceType: 1,
		SourceID:   session.TaskID,
		TaskID:     &session.TaskID,
		Title:      req.Title,
		Content:    req.Content,
		Summary:    truncateRunes(req.Content, 180),
		Category:   "团队协作",
		Status:     1,
	}
	if err := database.GetDB().Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存团队知识失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"entry": entry}})
}

func loadTaskForCollaboration(c *gin.Context) (models.Task, uint64, bool) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || taskID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return models.Task{}, 0, false
	}
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return models.Task{}, 0, false
	}
	task, err := findAccessibleTask(database.GetDB(), taskID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或无权限访问"})
		return models.Task{}, 0, false
	}
	if task.TaskType != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅团队任务支持协作会话"})
		return models.Task{}, 0, false
	}
	return task, userID, true
}

func loadSessionForCollaboration(c *gin.Context) (models.TaskCollaborationSession, uint64, bool) {
	sessionID, err := strconv.ParseUint(c.Param("sessionId"), 10, 64)
	if err != nil || sessionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的协作会话ID"})
		return models.TaskCollaborationSession{}, 0, false
	}
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return models.TaskCollaborationSession{}, 0, false
	}
	var session models.TaskCollaborationSession
	if err := database.GetDB().First(&session, sessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "协作会话不存在"})
		return models.TaskCollaborationSession{}, 0, false
	}
	if !canAccessCollaborationSession(database.GetDB(), &session, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问协作会话"})
		return models.TaskCollaborationSession{}, 0, false
	}
	return session, userID, true
}

func currentUserID(c *gin.Context) (uint64, bool) {
	if value, exists := c.Get("user_id"); exists {
		if id, ok := value.(uint64); ok && id > 0 {
			return id, true
		}
	}
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if !strings.HasPrefix(token, "mock-token-") {
		return 0, false
	}
	raw := strings.TrimPrefix(token, "mock-token-")
	parts := strings.Split(raw, "-")
	if len(parts) == 0 {
		return 0, false
	}
	id, err := strconv.ParseUint(parts[0], 10, 64)
	return id, err == nil && id > 0
}

func findAccessibleTask(db *gorm.DB, taskID, userID uint64) (models.Task, error) {
	var task models.Task
	err := db.Where(`
		id = ? AND (
			created_by = ?
			OR owner_user_id = ?
			OR id IN (SELECT task_id FROM task_assignees WHERE user_id = ?)
			OR (
				task_type = 2
				AND owner_team_id IN (SELECT team_id FROM team_members WHERE user_id = ?)
			)
		)`, taskID, userID, userID, userID, userID).First(&task).Error
	return task, err
}

func deriveTaskCollaborationParticipants(db *gorm.DB, task *models.Task, requesterID uint64) []uint64 {
	seen := map[uint64]bool{}
	add := func(id uint64) {
		if id > 0 {
			seen[id] = true
		}
	}
	add(requesterID)
	add(task.CreatedBy)
	if task.OwnerUserID != nil {
		add(*task.OwnerUserID)
	}
	if task.ParentID != nil {
		var parent models.Task
		if err := db.First(&parent, *task.ParentID).Error; err == nil {
			add(parent.CreatedBy)
			if parent.OwnerUserID != nil {
				add(*parent.OwnerUserID)
			}
		}
	}
	if task.OwnerTeamID != nil {
		var members []models.TeamMember
		_ = db.Where("team_id = ?", *task.OwnerTeamID).Find(&members).Error
		for _, member := range members {
			add(member.UserID)
		}
	}
	var assignees []models.TaskAssignee
	_ = db.Where("task_id = ?", task.ID).Find(&assignees).Error
	for _, assignee := range assignees {
		add(assignee.UserID)
	}
	result := make([]uint64, 0, len(seen))
	for id := range seen {
		result = append(result, id)
	}
	return result
}

func canAccessCollaborationSession(db *gorm.DB, session *models.TaskCollaborationSession, userID uint64) bool {
	var count int64
	if err := db.Model(&models.TaskCollaborationParticipant{}).
		Where("session_id = ? AND user_id = ?", session.ID, userID).
		Count(&count).Error; err == nil && count > 0 {
		return true
	}
	_, err := findAccessibleTask(db, session.TaskID, userID)
	return err == nil
}

func getCollaborationSessionByRoom(db *gorm.DB, roomID uint64) (models.TaskCollaborationSession, bool) {
	var session models.TaskCollaborationSession
	if err := db.Where("room_id = ?", roomID).First(&session).Error; err != nil {
		return models.TaskCollaborationSession{}, false
	}
	return session, true
}

func buildCollaborationSessionResponse(session *models.TaskCollaborationSession) collaborationSessionResponse {
	var minutes interface{}
	if len(session.Minutes) > 0 {
		_ = json.Unmarshal(session.Minutes, &minutes)
	}
	label := "进行中"
	if session.Status == models.TaskCollaborationStatusDismissed {
		label = "已结束"
	}
	return collaborationSessionResponse{
		ID:           session.ID,
		TaskID:       session.TaskID,
		TeamID:       session.TeamID,
		RoomID:       session.RoomID,
		Status:       session.Status,
		StatusLabel:  label,
		IsActive:     session.Status == models.TaskCollaborationStatusActive,
		CreatedBy:    session.CreatedBy,
		DismissedAt:  session.DismissedAt,
		Minutes:      minutes,
		Participants: loadCollaborationParticipants(session.ID),
		CreatedAt:    session.CreatedAt,
	}
}

func buildChatMessageResponse(msg models.ChatMessage, displayName string) map[string]interface{} {
	messageType := "text"
	content := interface{}(msg.Content)
	if msg.MsgType == models.ChatMessageTypeSystem {
		messageType = "system"
	}
	if msg.MsgType == models.ChatMessageTypeKnowledgeCard {
		messageType = "knowledge_card"
		var payload knowledgeCardPayload
		if err := json.Unmarshal([]byte(msg.Content), &payload); err == nil {
			content = payload
		}
	}
	return map[string]interface{}{
		"id":           msg.ID,
		"user_id":      msg.UserID,
		"display_name": displayName,
		"content":      content,
		"message_type": messageType,
		"msg_type":     msg.MsgType,
		"session_id":   msg.SessionID,
		"sent_at":      msg.SentAt,
	}
}

func buildRuleBasedMinutes(task *models.Task, messages []models.ChatMessage) minutesPayload {
	var texts []string
	var knowledge []string
	for _, msg := range messages {
		if msg.MsgType == models.ChatMessageTypeKnowledgeCard {
			var card knowledgeCardPayload
			if err := json.Unmarshal([]byte(msg.Content), &card); err == nil && card.Title != "" {
				knowledge = append(knowledge, card.Title)
			}
			continue
		}
		if strings.TrimSpace(msg.Content) != "" {
			texts = append(texts, strings.TrimSpace(msg.Content))
		}
	}
	summary := "本次围绕任务进行了同步讨论。"
	if len(texts) > 0 {
		summary = truncateRunes(strings.Join(texts, "；"), 220)
	}
	next := []string{"根据本次讨论推进任务并在任务评论中同步进展"}
	if task.Title != "" {
		next = []string{"继续推进任务：" + task.Title}
	}
	return minutesPayload{
		Summary:               summary,
		SynchronizedKnowledge: knowledge,
		ActionItems:           []minutesActionItem{{Owner: "待确认", Action: "根据讨论结果补充负责人和截止时间"}},
		Blockers:              []string{},
		NextSteps:             next,
		GeneratedAt:           time.Now(),
	}
}

func generateAIMinutes(apiKey string, task *models.Task, messages []models.ChatMessage) (minutesPayload, error) {
	lines := make([]string, 0, len(messages))
	for _, msg := range messages {
		if msg.MsgType == models.ChatMessageTypeKnowledgeCard {
			var card knowledgeCardPayload
			if err := json.Unmarshal([]byte(msg.Content), &card); err == nil {
				lines = append(lines, fmt.Sprintf("用户%d共享知识卡片：%s - %s", msg.UserID, card.Title, firstNonEmpty(card.Summary, card.Excerpt)))
			}
			continue
		}
		lines = append(lines, fmt.Sprintf("用户%d：%s", msg.UserID, msg.Content))
	}
	prompt := fmt.Sprintf(`你是团队任务协作纪要助手。请只输出 JSON，不要输出 Markdown。
任务标题：%s
任务描述：%s
讨论记录：
%s

输出字段：
{
  "summary": "100字以内讨论摘要",
  "synchronized_knowledge": ["成员同步的知识点"],
  "action_items": [{"owner":"负责人，无法判断则写待确认","action":"要做什么"}],
  "blockers": ["阻塞点，没有则为空数组"],
  "next_steps": ["下一步"]
}`, task.Title, task.Description, strings.Join(lines, "\n"))

	reply, err := callQwenForChat(apiKey, prompt, nil)
	if err != nil {
		return minutesPayload{}, err
	}
	var minutes minutesPayload
	if err := json.Unmarshal([]byte(extractJSONObject(reply)), &minutes); err != nil {
		return minutesPayload{}, err
	}
	if minutes.GeneratedAt.IsZero() {
		minutes.GeneratedAt = time.Now()
	}
	if strings.TrimSpace(minutes.Summary) == "" {
		return minutesPayload{}, fmt.Errorf("ai minutes missing summary")
	}
	return minutes, nil
}

func extractJSONObject(value string) string {
	trimmed := strings.TrimSpace(value)
	if strings.HasPrefix(trimmed, "```") {
		trimmed = strings.TrimPrefix(trimmed, "```json")
		trimmed = strings.TrimPrefix(trimmed, "```")
		trimmed = strings.TrimSuffix(trimmed, "```")
		trimmed = strings.TrimSpace(trimmed)
	}
	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start >= 0 && end >= start {
		return trimmed[start : end+1]
	}
	return trimmed
}

func appendMinutesComment(taskID, userID uint64, raw json.RawMessage) {
	db := database.GetDB()
	var task models.Task
	if err := db.First(&task, taskID).Error; err != nil {
		return
	}
	var comments []TaskComment
	if len(task.Comments) > 0 {
		_ = json.Unmarshal(task.Comments, &comments)
	}
	comments = append(comments, TaskComment{
		UserID:    userID,
		Content:   "AI 协作纪要：\n" + formatReadableMinutes(raw),
		CreatedAt: time.Now(),
	})
	encoded, _ := json.Marshal(comments)
	_ = db.Model(&task).Update("comments", datatypes.JSON(encoded)).Error
}

func loadCollaborationParticipants(sessionID uint64) []collaborationParticipantView {
	var participants []models.TaskCollaborationParticipant
	_ = database.GetDB().Where("session_id = ?", sessionID).Order("role DESC, user_id ASC").Find(&participants).Error
	ids := make([]uint64, 0, len(participants))
	for _, participant := range participants {
		ids = append(ids, participant.UserID)
	}
	names := loadUserNames(ids)
	result := make([]collaborationParticipantView, 0, len(participants))
	for _, participant := range participants {
		result = append(result, collaborationParticipantView{
			UserID:      participant.UserID,
			DisplayName: firstNonEmpty(names[participant.UserID], fmt.Sprintf("用户 %d", participant.UserID)),
			Role:        participant.Role,
		})
	}
	return result
}

func createCollaborationInviteMessage(db *gorm.DB, session *models.TaskCollaborationSession, inviterID uint64, participantIDs []uint64) {
	names := loadUserNames(participantIDs)
	displayNames := make([]string, 0, len(participantIDs))
	for _, id := range participantIDs {
		displayNames = append(displayNames, firstNonEmpty(names[id], fmt.Sprintf("用户 %d", id)))
	}
	content := "你邀请了 " + strings.Join(displayNames, "、") + " 进入任务协作"
	chat := models.ChatMessage{
		SessionID: session.ID,
		RoomID:    session.RoomID,
		UserID:    inviterID,
		Content:   content,
		MsgType:   models.ChatMessageTypeSystem,
		SentAt:    time.Now(),
	}
	_ = db.Create(&chat).Error
}

func formatReadableMinutes(raw []byte) string {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return ""
	}
	var minutes minutesPayload
	if err := json.Unmarshal([]byte(trimmed), &minutes); err == nil {
		return renderMinutesPayload(minutes)
	}
	var generic map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &generic); err != nil {
		return trimmed
	}
	return renderGenericMinutes(generic)
}

func renderMinutesPayload(minutes minutesPayload) string {
	var lines []string
	if strings.TrimSpace(minutes.Summary) != "" {
		lines = append(lines, "## 讨论摘要", minutes.Summary)
	}
	lines = appendStringSection(lines, "同步内容", minutes.SynchronizedKnowledge)
	if len(minutes.ActionItems) > 0 {
		lines = append(lines, "## 行动项")
		for _, item := range minutes.ActionItems {
			lines = append(lines, fmt.Sprintf("- %s：%s", firstNonEmpty(item.Owner, "待确认"), item.Action))
		}
	}
	lines = appendStringSection(lines, "阻塞点", minutes.Blockers)
	lines = appendStringSection(lines, "下一步", minutes.NextSteps)
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n")
}

func renderGenericMinutes(value map[string]interface{}) string {
	var lines []string
	if summary, ok := value["summary"].(string); ok && strings.TrimSpace(summary) != "" {
		lines = append(lines, "## 讨论摘要", summary)
	}
	for _, key := range []struct {
		name  string
		title string
	}{
		{"synchronized_knowledge", "同步内容"},
		{"blockers", "阻塞点"},
		{"next_steps", "下一步"},
	} {
		if items := interfaceStringList(value[key.name]); len(items) > 0 {
			lines = appendStringSection(lines, key.title, items)
		}
	}
	if items, ok := value["action_items"].([]interface{}); ok && len(items) > 0 {
		lines = append(lines, "## 行动项")
		for _, raw := range items {
			if item, ok := raw.(map[string]interface{}); ok {
				lines = append(lines, fmt.Sprintf("- %s：%s", firstNonEmpty(fmt.Sprint(item["owner"]), "待确认"), fmt.Sprint(item["action"])))
			}
		}
	}
	if len(lines) == 0 {
		encoded, _ := json.MarshalIndent(value, "", "  ")
		return string(encoded)
	}
	return strings.Join(lines, "\n")
}

func appendStringSection(lines []string, title string, items []string) []string {
	if len(items) == 0 {
		return lines
	}
	lines = append(lines, "## "+title)
	for _, item := range items {
		if strings.TrimSpace(item) != "" {
			lines = append(lines, "- "+strings.TrimSpace(item))
		}
	}
	return lines
}

func interfaceStringList(value interface{}) []string {
	raw, ok := value.([]interface{})
	if !ok {
		return nil
	}
	items := make([]string, 0, len(raw))
	for _, item := range raw {
		if text := strings.TrimSpace(fmt.Sprint(item)); text != "" {
			items = append(items, text)
		}
	}
	return items
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func truncateRunes(value string, max int) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= max {
		return string(runes)
	}
	return string(runes[:max]) + "..."
}
