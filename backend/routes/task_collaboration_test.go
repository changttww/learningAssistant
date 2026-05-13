package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

func setupTaskCollaborationTest(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	t.Setenv("QWEN_API_KEY", "")
	gin.SetMode(gin.TestMode)

	dsn := "file:" + strings.NewReplacer("/", "_", " ", "_").Replace(t.Name()) + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	database.DB = db
	if err := db.AutoMigrate(models.GetAllModels()...); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	r := gin.New()
	api := r.Group("/api")
	tasks := api.Group("/tasks")
	registerTaskRoutes(tasks)
	teams := api.Group("/teams")
	registerTeamRoutes(teams)
	study := api.Group("/study")
	registerStudyRoutes(study)
	{
		rooms := study.Group("/rooms")
		rooms.GET("/:roomId/chat/history", handleGetRoomChatHistory)
		rooms.POST("/:roomId/chat", handlePostRoomChat)
	}
	return r, db
}

func authRequest(method, path string, userID uint64, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-token-"+jsonNumber(userID)+"-test")
	return req
}

func jsonNumber(v uint64) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func decodeBody(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var payload map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response %q: %v", rr.Body.String(), err)
	}
	return payload
}

func seedCollaborationTask(t *testing.T, db *gorm.DB) models.Task {
	t.Helper()
	ownerID := uint64(1)
	users := []models.User{
		{BaseModel: models.BaseModel{ID: ownerID}, Account: "owner", Email: "owner@example.com", Phone: "10000000001", PasswordHash: "x", DisplayName: "负责人"},
		{BaseModel: models.BaseModel{ID: 2}, Account: "partner", Email: "partner@example.com", Phone: "10000000002", PasswordHash: "x", DisplayName: "协作者"},
		{BaseModel: models.BaseModel{ID: 3}, Account: "reviewer", Email: "reviewer@example.com", Phone: "10000000003", PasswordHash: "x", DisplayName: "评审"},
	}
	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("create user %d: %v", user.ID, err)
		}
	}
	team := models.Team{Name: "Spec Team", OwnerUserID: ownerID}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	if err := db.Create(&models.TeamMember{TeamID: team.ID, UserID: ownerID}).Error; err != nil {
		t.Fatalf("create owner member: %v", err)
	}
	if err := db.Create(&models.TeamMember{TeamID: team.ID, UserID: 2}).Error; err != nil {
		t.Fatalf("create collaborator member: %v", err)
	}
	if err := db.Create(&models.TeamMember{TeamID: team.ID, UserID: 3}).Error; err != nil {
		t.Fatalf("create reviewer member: %v", err)
	}
	task := models.Task{
		Title:       "Build task room",
		Description: "Discuss task collaboration",
		TaskType:    2,
		CreatedBy:   ownerID,
		OwnerUserID: &ownerID,
		OwnerTeamID: &team.ID,
	}
	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("create task: %v", err)
	}
	if err := db.Create(&models.TaskAssignee{TaskID: task.ID, UserID: 2}).Error; err != nil {
		t.Fatalf("create assignee: %v", err)
	}
	return task
}

func TestCreateTaskCollaborationSessionReusesActiveRoomAndAddsTeamParticipants(t *testing.T) {
	r, db := setupTaskCollaborationTest(t)
	task := seedCollaborationTask(t, db)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, authRequest(http.MethodPost, "/api/tasks/"+jsonNumber(task.ID)+"/collaboration-sessions", 1, nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first create status = %d body=%s", first.Code, first.Body.String())
	}
	second := httptest.NewRecorder()
	r.ServeHTTP(second, authRequest(http.MethodPost, "/api/tasks/"+jsonNumber(task.ID)+"/collaboration-sessions", 1, nil))
	if second.Code != http.StatusOK {
		t.Fatalf("second create status = %d body=%s", second.Code, second.Body.String())
	}

	firstData := decodeBody(t, first)["data"].(map[string]interface{})["session"].(map[string]interface{})
	secondData := decodeBody(t, second)["data"].(map[string]interface{})["session"].(map[string]interface{})
	if firstData["id"] != secondData["id"] {
		t.Fatalf("expected active session to be reused, got first=%v second=%v", firstData["id"], secondData["id"])
	}
	if firstData["room_id"] != secondData["room_id"] {
		t.Fatalf("expected active room to be reused, got first=%v second=%v", firstData["room_id"], secondData["room_id"])
	}

	var participantCount int64
	if err := db.Model(&models.TaskCollaborationParticipant{}).
		Where("session_id = ?", uint64(firstData["id"].(float64))).
		Count(&participantCount).Error; err != nil {
		t.Fatalf("count participants: %v", err)
	}
	if participantCount != 3 {
		t.Fatalf("expected all team members as participants, got %d", participantCount)
	}

	var systemCount int64
	if err := db.Model(&models.ChatMessage{}).
		Where("session_id = ? AND msg_type = ?", uint64(firstData["id"].(float64)), int8(3)).
		Count(&systemCount).Error; err != nil {
		t.Fatalf("count system messages: %v", err)
	}
	if systemCount != 1 {
		t.Fatalf("expected exactly one invitation system message, got %d", systemCount)
	}

	detail := httptest.NewRecorder()
	r.ServeHTTP(detail, authRequest(http.MethodGet, "/api/tasks/collaboration-sessions/"+jsonNumber(uint64(firstData["id"].(float64))), 1, nil))
	if detail.Code != http.StatusOK {
		t.Fatalf("detail status = %d body=%s", detail.Code, detail.Body.String())
	}
	detailSession := decodeBody(t, detail)["data"].(map[string]interface{})["session"].(map[string]interface{})
	participants, ok := detailSession["participants"].([]interface{})
	if !ok || len(participants) != 3 {
		t.Fatalf("expected participant display data for three team members, got %#v", detailSession["participants"])
	}
}

func TestTaskCollaborationDismissalAllowsNewActiveRoom(t *testing.T) {
	r, _ := setupTaskCollaborationTest(t)
	task := seedCollaborationTask(t, database.GetDB())

	first := httptest.NewRecorder()
	r.ServeHTTP(first, authRequest(http.MethodPost, "/api/tasks/"+jsonNumber(task.ID)+"/collaboration-sessions", 1, nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first create status = %d body=%s", first.Code, first.Body.String())
	}
	firstSession := decodeBody(t, first)["data"].(map[string]interface{})["session"].(map[string]interface{})
	firstID := uint64(firstSession["id"].(float64))

	dismiss := httptest.NewRecorder()
	r.ServeHTTP(dismiss, authRequest(http.MethodPost, "/api/tasks/collaboration-sessions/"+jsonNumber(firstID)+"/dismiss", 1, nil))
	if dismiss.Code != http.StatusOK {
		t.Fatalf("dismiss status = %d body=%s", dismiss.Code, dismiss.Body.String())
	}

	second := httptest.NewRecorder()
	r.ServeHTTP(second, authRequest(http.MethodPost, "/api/tasks/"+jsonNumber(task.ID)+"/collaboration-sessions", 1, nil))
	if second.Code != http.StatusOK {
		t.Fatalf("second create status = %d body=%s", second.Code, second.Body.String())
	}
	secondSession := decodeBody(t, second)["data"].(map[string]interface{})["session"].(map[string]interface{})
	if firstSession["id"] == secondSession["id"] {
		t.Fatalf("expected new active room after dismissal, got %v", secondSession["id"])
	}
}

func TestTaskCollaborationDismissalDoesNotReactivateOlderActiveRooms(t *testing.T) {
	r, db := setupTaskCollaborationTest(t)
	task := seedCollaborationTask(t, db)

	oldRoom := models.StudyRoom{
		Name:        "旧协作房间",
		OwnerUserID: 1,
		TeamID:      task.OwnerTeamID,
		RoomKind:    "task_collaboration",
		IsPrivate:   true,
		Status:      1,
	}
	if err := db.Create(&oldRoom).Error; err != nil {
		t.Fatalf("create old room: %v", err)
	}
	oldSession := models.TaskCollaborationSession{
		TaskID:    task.ID,
		TeamID:    task.OwnerTeamID,
		RoomID:    oldRoom.ID,
		CreatedBy: 1,
		Status:    models.TaskCollaborationStatusActive,
	}
	if err := db.Create(&oldSession).Error; err != nil {
		t.Fatalf("create old active session: %v", err)
	}
	newRoom := models.StudyRoom{
		Name:        "新协作房间",
		OwnerUserID: 1,
		TeamID:      task.OwnerTeamID,
		RoomKind:    "task_collaboration",
		IsPrivate:   true,
		Status:      1,
	}
	if err := db.Create(&newRoom).Error; err != nil {
		t.Fatalf("create new room: %v", err)
	}
	newSession := models.TaskCollaborationSession{
		TaskID:    task.ID,
		TeamID:    task.OwnerTeamID,
		RoomID:    newRoom.ID,
		CreatedBy: 1,
		Status:    models.TaskCollaborationStatusActive,
	}
	if err := db.Create(&newSession).Error; err != nil {
		t.Fatalf("create new active session: %v", err)
	}

	create := httptest.NewRecorder()
	r.ServeHTTP(create, authRequest(http.MethodPost, "/api/tasks/"+jsonNumber(task.ID)+"/collaboration-sessions", 1, nil))
	if create.Code != http.StatusOK {
		t.Fatalf("create status = %d body=%s", create.Code, create.Body.String())
	}
	session := decodeBody(t, create)["data"].(map[string]interface{})["session"].(map[string]interface{})
	activeID := uint64(session["id"].(float64))
	if activeID != newSession.ID {
		t.Fatalf("expected newest active session to be canonical, got %d want %d", activeID, newSession.ID)
	}

	dismiss := httptest.NewRecorder()
	r.ServeHTTP(dismiss, authRequest(http.MethodPost, "/api/tasks/collaboration-sessions/"+jsonNumber(activeID)+"/dismiss", 1, nil))
	if dismiss.Code != http.StatusOK {
		t.Fatalf("dismiss status = %d body=%s", dismiss.Code, dismiss.Body.String())
	}

	var activeCount int64
	if err := db.Model(&models.TaskCollaborationSession{}).
		Where("task_id = ? AND status = ?", task.ID, models.TaskCollaborationStatusActive).
		Count(&activeCount).Error; err != nil {
		t.Fatalf("count active sessions: %v", err)
	}
	if activeCount != 0 {
		t.Fatalf("expected no active sessions after dismissal, got %d", activeCount)
	}
}

func TestTaskCollaborationSessionRejectsUnauthorizedAndHidesFromStudyRooms(t *testing.T) {
	r, db := setupTaskCollaborationTest(t)
	task := seedCollaborationTask(t, db)

	unauthorized := httptest.NewRecorder()
	r.ServeHTTP(unauthorized, authRequest(http.MethodPost, "/api/tasks/"+jsonNumber(task.ID)+"/collaboration-sessions", 99, nil))
	if unauthorized.Code != http.StatusNotFound && unauthorized.Code != http.StatusForbidden {
		t.Fatalf("expected unauthorized rejection, got %d body=%s", unauthorized.Code, unauthorized.Body.String())
	}

	create := httptest.NewRecorder()
	r.ServeHTTP(create, authRequest(http.MethodPost, "/api/tasks/"+jsonNumber(task.ID)+"/collaboration-sessions", 1, nil))
	if create.Code != http.StatusOK {
		t.Fatalf("create status = %d body=%s", create.Code, create.Body.String())
	}

	list := httptest.NewRecorder()
	r.ServeHTTP(list, authRequest(http.MethodGet, "/api/study/rooms", 1, nil))
	if list.Code != http.StatusOK {
		t.Fatalf("list rooms status = %d body=%s", list.Code, list.Body.String())
	}
	payload := decodeBody(t, list)
	rooms := payload["data"].(map[string]interface{})["rooms"].([]interface{})
	if len(rooms) != 0 {
		t.Fatalf("expected task collaboration rooms hidden from study discovery, got %d", len(rooms))
	}
}

func TestTaskCollaborationMessagesKnowledgeCardsAndDismissal(t *testing.T) {
	r, db := setupTaskCollaborationTest(t)
	task := seedCollaborationTask(t, db)

	create := httptest.NewRecorder()
	r.ServeHTTP(create, authRequest(http.MethodPost, "/api/tasks/"+jsonNumber(task.ID)+"/collaboration-sessions", 1, nil))
	if create.Code != http.StatusOK {
		t.Fatalf("create status = %d body=%s", create.Code, create.Body.String())
	}
	session := decodeBody(t, create)["data"].(map[string]interface{})["session"].(map[string]interface{})
	sessionID := uint64(session["id"].(float64))
	roomID := uint64(session["room_id"].(float64))

	message := httptest.NewRecorder()
	r.ServeHTTP(message, authRequest(http.MethodPost, "/api/study/rooms/"+jsonNumber(roomID)+"/chat", 1, map[string]interface{}{
		"content": "我来同步接口约定",
	}))
	if message.Code != http.StatusOK {
		t.Fatalf("post text status = %d body=%s", message.Code, message.Body.String())
	}

	entry := models.KnowledgeBaseEntry{
		UserID:   1,
		Title:    "接口约定",
		Content:  "请求字段和错误码约定",
		Summary:  "接口字段说明",
		Category: "backend",
		Status:   1,
	}
	if err := db.Create(&entry).Error; err != nil {
		t.Fatalf("create knowledge: %v", err)
	}

	card := httptest.NewRecorder()
	r.ServeHTTP(card, authRequest(http.MethodPost, "/api/tasks/collaboration-sessions/"+jsonNumber(sessionID)+"/knowledge-cards", 1, map[string]interface{}{
		"knowledge_entry_id": entry.ID,
	}))
	if card.Code != http.StatusOK {
		t.Fatalf("post knowledge card status = %d body=%s", card.Code, card.Body.String())
	}

	detail := httptest.NewRecorder()
	r.ServeHTTP(detail, authRequest(http.MethodGet, "/api/tasks/collaboration-sessions/"+jsonNumber(sessionID)+"/knowledge-cards/"+jsonNumber(entry.ID), 2, nil))
	if detail.Code != http.StatusOK {
		t.Fatalf("knowledge card detail status = %d body=%s", detail.Code, detail.Body.String())
	}
	detailEntry := decodeBody(t, detail)["data"].(map[string]interface{})["entry"].(map[string]interface{})
	if detailEntry["content"] != entry.Content {
		t.Fatalf("expected shared note detail content, got %#v", detailEntry)
	}

	history := httptest.NewRecorder()
	r.ServeHTTP(history, authRequest(http.MethodGet, "/api/study/rooms/"+jsonNumber(roomID)+"/chat/history", 1, nil))
	if history.Code != http.StatusOK {
		t.Fatalf("history status = %d body=%s", history.Code, history.Body.String())
	}
	messages := decodeBody(t, history)["data"].(map[string]interface{})["messages"].([]interface{})
	if len(messages) != 3 {
		t.Fatalf("expected system, text, and knowledge messages, got %d", len(messages))
	}
	if messages[0].(map[string]interface{})["message_type"] != "knowledge_card" {
		t.Fatalf("expected latest message to be knowledge card, got %#v", messages[0])
	}

	dismiss := httptest.NewRecorder()
	r.ServeHTTP(dismiss, authRequest(http.MethodPost, "/api/tasks/collaboration-sessions/"+jsonNumber(sessionID)+"/dismiss", 1, nil))
	if dismiss.Code != http.StatusOK {
		t.Fatalf("dismiss status = %d body=%s", dismiss.Code, dismiss.Body.String())
	}

	afterDismiss := httptest.NewRecorder()
	r.ServeHTTP(afterDismiss, authRequest(http.MethodPost, "/api/study/rooms/"+jsonNumber(roomID)+"/chat", 1, map[string]interface{}{
		"content": "还能写吗",
	}))
	if afterDismiss.Code != http.StatusConflict {
		t.Fatalf("expected conflict after dismissal, got %d body=%s", afterDismiss.Code, afterDismiss.Body.String())
	}

	var sessionModel models.TaskCollaborationSession
	if err := db.First(&sessionModel, sessionID).Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if sessionModel.DismissedAt == nil || time.Since(*sessionModel.DismissedAt) > time.Minute {
		t.Fatalf("expected recent dismissed_at, got %#v", sessionModel.DismissedAt)
	}
}

func TestTaskCollaborationMinutesSaveAndTeamKnowledge(t *testing.T) {
	r, db := setupTaskCollaborationTest(t)
	task := seedCollaborationTask(t, db)

	create := httptest.NewRecorder()
	r.ServeHTTP(create, authRequest(http.MethodPost, "/api/tasks/"+jsonNumber(task.ID)+"/collaboration-sessions", 1, nil))
	if create.Code != http.StatusOK {
		t.Fatalf("create status = %d body=%s", create.Code, create.Body.String())
	}
	session := decodeBody(t, create)["data"].(map[string]interface{})["session"].(map[string]interface{})
	sessionID := uint64(session["id"].(float64))
	roomID := uint64(session["room_id"].(float64))

	empty := httptest.NewRecorder()
	r.ServeHTTP(empty, authRequest(http.MethodPost, "/api/tasks/collaboration-sessions/"+jsonNumber(sessionID)+"/minutes/generate", 1, nil))
	if empty.Code != http.StatusBadRequest {
		t.Fatalf("expected empty-content rejection, got %d body=%s", empty.Code, empty.Body.String())
	}

	message := httptest.NewRecorder()
	r.ServeHTTP(message, authRequest(http.MethodPost, "/api/study/rooms/"+jsonNumber(roomID)+"/chat", 1, map[string]interface{}{
		"content": "李同学负责接口联调，下一步补充错误码",
	}))
	if message.Code != http.StatusOK {
		t.Fatalf("post text status = %d body=%s", message.Code, message.Body.String())
	}
	if err := db.Create(&models.ChatMessage{
		SessionID: sessionID,
		RoomID:    roomID + 999,
		UserID:    1,
		Content:   "其他房间的同号会话消息不应进入纪要",
		MsgType:   models.ChatMessageTypeText,
		SentAt:    time.Now(),
	}).Error; err != nil {
		t.Fatalf("seed cross-room message: %v", err)
	}

	generate := httptest.NewRecorder()
	r.ServeHTTP(generate, authRequest(http.MethodPost, "/api/tasks/collaboration-sessions/"+jsonNumber(sessionID)+"/minutes/generate", 1, nil))
	if generate.Code != http.StatusOK {
		t.Fatalf("generate status = %d body=%s", generate.Code, generate.Body.String())
	}
	minutes := decodeBody(t, generate)["data"].(map[string]interface{})["minutes"]
	if summary, _ := minutes.(map[string]interface{})["summary"].(string); strings.Contains(summary, "其他房间") {
		t.Fatalf("minutes leaked cross-room message: %s", summary)
	}

	save := httptest.NewRecorder()
	r.ServeHTTP(save, authRequest(http.MethodPost, "/api/tasks/collaboration-sessions/"+jsonNumber(sessionID)+"/minutes/save", 1, map[string]interface{}{
		"minutes": minutes,
	}))
	if save.Code != http.StatusOK {
		t.Fatalf("save status = %d body=%s", save.Code, save.Body.String())
	}
	duplicate := httptest.NewRecorder()
	r.ServeHTTP(duplicate, authRequest(http.MethodPost, "/api/tasks/collaboration-sessions/"+jsonNumber(sessionID)+"/minutes/save", 1, map[string]interface{}{
		"minutes": minutes,
	}))
	if duplicate.Code != http.StatusOK {
		t.Fatalf("duplicate save status = %d body=%s", duplicate.Code, duplicate.Body.String())
	}
	if dup := decodeBody(t, duplicate)["data"].(map[string]interface{})["duplicate"]; dup != true {
		t.Fatalf("expected duplicate save marker, got %#v", dup)
	}

	var updatedTask models.Task
	if err := db.First(&updatedTask, task.ID).Error; err != nil {
		t.Fatalf("load task: %v", err)
	}
	if len(updatedTask.Comments) == 0 {
		t.Fatalf("expected minutes saved into task comments")
	}
	var comments []TaskComment
	if err := json.Unmarshal(updatedTask.Comments, &comments); err != nil {
		t.Fatalf("decode task comments: %v", err)
	}
	if len(comments) == 0 || strings.Contains(comments[len(comments)-1].Content, `{"summary"`) {
		t.Fatalf("expected readable minutes comment, got %#v", comments)
	}

	knowledge := httptest.NewRecorder()
	r.ServeHTTP(knowledge, authRequest(http.MethodPost, "/api/tasks/collaboration-sessions/"+jsonNumber(sessionID)+"/team-knowledge", 1, map[string]interface{}{
		"title":   "接口联调纪要",
		"content": `{"summary":"接口联调讨论结论","synchronized_knowledge":["接口约定"],"action_items":[{"owner":"李同学","action":"补充错误码"}],"blockers":[],"next_steps":["继续联调"]}`,
	}))
	if knowledge.Code != http.StatusOK {
		t.Fatalf("team knowledge status = %d body=%s", knowledge.Code, knowledge.Body.String())
	}
	var count int64
	if err := db.Model(&models.KnowledgeBaseEntry{}).Where("team_id = ? AND task_id = ?", *task.OwnerTeamID, task.ID).Count(&count).Error; err != nil {
		t.Fatalf("count team knowledge: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one team knowledge entry, got %d", count)
	}
	var entry models.KnowledgeBaseEntry
	if err := db.Where("team_id = ? AND task_id = ?", *task.OwnerTeamID, task.ID).First(&entry).Error; err != nil {
		t.Fatalf("load team knowledge: %v", err)
	}
	if strings.Contains(entry.Content, `{"summary"`) || !strings.Contains(entry.Content, "接口联调讨论结论") {
		t.Fatalf("expected readable team knowledge content, got %q", entry.Content)
	}
}

func TestPersistentTeamChatRoomCreationReuseAndAuthorization(t *testing.T) {
	r, db := setupTaskCollaborationTest(t)
	task := seedCollaborationTask(t, db)
	teamID := *task.OwnerTeamID

	first := httptest.NewRecorder()
	r.ServeHTTP(first, authRequest(http.MethodPost, "/api/teams/"+jsonNumber(teamID)+"/chat-room", 1, nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first team chat status = %d body=%s", first.Code, first.Body.String())
	}
	firstRoom := decodeBody(t, first)["data"].(map[string]interface{})["room"].(map[string]interface{})

	second := httptest.NewRecorder()
	r.ServeHTTP(second, authRequest(http.MethodPost, "/api/teams/"+jsonNumber(teamID)+"/chat-room", 2, nil))
	if second.Code != http.StatusOK {
		t.Fatalf("second team chat status = %d body=%s", second.Code, second.Body.String())
	}
	secondRoom := decodeBody(t, second)["data"].(map[string]interface{})["room"].(map[string]interface{})
	if firstRoom["id"] != secondRoom["id"] {
		t.Fatalf("expected persistent team chat room reuse, got first=%v second=%v", firstRoom["id"], secondRoom["id"])
	}

	roomID := uint64(firstRoom["id"].(float64))
	message := httptest.NewRecorder()
	r.ServeHTTP(message, authRequest(http.MethodPost, "/api/study/rooms/"+jsonNumber(roomID)+"/chat", 2, map[string]interface{}{
		"content": "团队长期沟通消息",
	}))
	if message.Code != http.StatusOK {
		t.Fatalf("team chat message status = %d body=%s", message.Code, message.Body.String())
	}
	history := httptest.NewRecorder()
	r.ServeHTTP(history, authRequest(http.MethodGet, "/api/study/rooms/"+jsonNumber(roomID)+"/chat/history", 1, nil))
	if history.Code != http.StatusOK {
		t.Fatalf("team chat history status = %d body=%s", history.Code, history.Body.String())
	}
	messages := decodeBody(t, history)["data"].(map[string]interface{})["messages"].([]interface{})
	if len(messages) != 1 || messages[0].(map[string]interface{})["content"] != "团队长期沟通消息" {
		t.Fatalf("expected persisted team chat message, got %#v", messages)
	}

	forbidden := httptest.NewRecorder()
	r.ServeHTTP(forbidden, authRequest(http.MethodPost, "/api/teams/"+jsonNumber(teamID)+"/chat-room", 99, nil))
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("expected non-member rejection, got %d body=%s", forbidden.Code, forbidden.Body.String())
	}
}
