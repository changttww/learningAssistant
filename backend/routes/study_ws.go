package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

var studyHubRegistry = newStudyHubStore()

type wsEnvelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type wsMemberState struct {
	UserID      uint64  `json:"user_id"`
	DisplayName string  `json:"display_name"`
	PeerID      string  `json:"peer_id"`
	IsBusy      bool    `json:"is_busy"`
	PartnerID   *uint64 `json:"partner_id,omitempty"`
}

type studyHubStore struct {
	mu   sync.Mutex
	hubs map[uint64]*studyRoomHub
}

func newStudyHubStore() *studyHubStore {
	return &studyHubStore{
		hubs: make(map[uint64]*studyRoomHub),
	}
}

func (m *studyHubStore) getHub(roomID uint64) *studyRoomHub {
	m.mu.Lock()
	defer m.mu.Unlock()
	if hub, ok := m.hubs[roomID]; ok {
		return hub
	}
	hub := newStudyRoomHub(roomID)
	m.hubs[roomID] = hub
	return hub
}

func (m *studyHubStore) totalOnline() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	total := 0
	for _, hub := range m.hubs {
		total += len(hub.clients)
	}
	return total
}

func (m *studyHubStore) activeMinutes() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	total := 0
	now := time.Now()
	for _, hub := range m.hubs {
		for _, cl := range hub.clients {
			if !cl.sessionStart.IsZero() {
				mins := int(now.Sub(cl.sessionStart).Minutes())
				if mins < 0 {
					mins = 0
				}
				total += mins
			}
		}
	}
	return total
}

type studyRoomHub struct {
	roomID   uint64
	clients  map[uint64]*studyClient
	pending  map[uint64]uint64 // target -> caller
	busy     map[uint64]uint64 // user -> partner
	mu       sync.Mutex
	upgrader websocket.Upgrader
}

func newStudyRoomHub(roomID uint64) *studyRoomHub {
	return &studyRoomHub{
		roomID:  roomID,
		clients: make(map[uint64]*studyClient),
		pending: make(map[uint64]uint64),
		busy:    make(map[uint64]uint64),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

type studyClient struct {
	userID       uint64
	displayName  string
	peerID       string
	conn         *websocket.Conn
	hub          *studyRoomHub
	send         chan wsEnvelope
	sessionStart time.Time
	recordID     uint64
}

func (h *studyRoomHub) handleWebSocket(c *gin.Context) {
	roomIDParam := c.Param("roomId")
	roomID, err := strconv.ParseUint(roomIDParam, 10, 64)
	if err != nil || roomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "房间ID不正确"})
		return
	}

	userIDParam := c.Query("user_id")
	displayName := strings.TrimSpace(c.Query("display_name"))
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少用户ID"})
		return
	}
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户ID无效"})
		return
	}
	if displayName == "" {
		displayName = "学习者"
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("websocket upgrade failed:", err)
		return
	}

	client := &studyClient{
		userID:      userID,
		displayName: displayName,
		conn:        conn,
		hub:         h,
		send:        make(chan wsEnvelope, 16),
	}

	h.registerClient(client)
	h.startSession(client)
	go client.writeLoop()
	go client.readLoop()

	client.send <- wsEnvelope{Type: "state", Data: mustMarshal(h.buildStatePayload())}
	h.broadcast(wsEnvelope{Type: "member_joined", Data: mustMarshal(h.memberState(client))}, client.userID)
}

func (h *studyRoomHub) registerClient(client *studyClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.userID] = client
}

func (h *studyRoomHub) unregisterClient(client *studyClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client.userID)
	delete(h.pending, client.userID)
	if partner, ok := h.busy[client.userID]; ok {
		delete(h.busy, client.userID)
		delete(h.busy, partner)
	}
	h.finishSession(client)
}

func (h *studyRoomHub) startSession(client *studyClient) {
	client.sessionStart = time.Now()
	record := models.LearningRecord{
		TaskID:          h.roomID,
		UserID:          client.userID,
		SessionStart:    client.sessionStart,
		DurationMinutes: 0,
		Note:            fmt.Sprintf("room:%d", h.roomID),
	}
	if err := database.GetDB().Create(&record).Error; err == nil {
		client.recordID = record.ID
	}
}

func (h *studyRoomHub) finishSession(client *studyClient) {
	if client.recordID == 0 || client.sessionStart.IsZero() {
		return
	}
	end := time.Now()
	duration := int(end.Sub(client.sessionStart).Minutes())
	if duration < 1 {
		duration = 1
	}
	db := database.GetDB()
	update := map[string]interface{}{
		"session_end":      end,
		"duration_minutes": duration,
	}
	_ = db.Model(&models.LearningRecord{}).
		Where("id = ?", client.recordID).
		Updates(update).Error

	_ = db.Model(&models.StudyRoom{}).
		Where("id = ?", h.roomID).
		UpdateColumn("focus_minutes_today", gorm.Expr("focus_minutes_today + ?", duration)).Error

	_ = db.Model(&models.UserProfile{}).
		Where("user_id = ?", client.userID).
		UpdateColumn("total_study_mins", gorm.Expr("total_study_mins + ?", duration)).Error
}

func (h *studyRoomHub) memberState(client *studyClient) wsMemberState {
	h.mu.Lock()
	defer h.mu.Unlock()
	member := wsMemberState{
		UserID:      client.userID,
		DisplayName: client.displayName,
		PeerID:      client.peerID,
	}
	if partner, ok := h.busy[client.userID]; ok {
		member.IsBusy = true
		member.PartnerID = &partner
	}
	return member
}

func (h *studyRoomHub) buildStatePayload() map[string]interface{} {
	h.mu.Lock()
	defer h.mu.Unlock()
	members := make([]wsMemberState, 0, len(h.clients))
	for _, cl := range h.clients {
		member := wsMemberState{
			UserID:      cl.userID,
			DisplayName: cl.displayName,
			PeerID:      cl.peerID,
		}
		if partner, ok := h.busy[cl.userID]; ok {
			member.IsBusy = true
			member.PartnerID = &partner
		}
		members = append(members, member)
	}
	return map[string]interface{}{
		"room_id": h.roomID,
		"members": members,
	}
}

func (h *studyRoomHub) broadcast(msg wsEnvelope, exclude ...uint64) {
	h.mu.Lock()
	targets := make([]*studyClient, 0, len(h.clients))
	for id, client := range h.clients {
		skip := false
		for _, ex := range exclude {
			if id == ex {
				skip = true
				break
			}
		}
		if !skip {
			targets = append(targets, client)
		}
	}
	h.mu.Unlock()

	for _, client := range targets {
		select {
		case client.send <- msg:
		default:
			go client.conn.Close()
		}
	}
}

func (h *studyRoomHub) currentOnline() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients)
}

func (c *studyClient) readLoop() {
	defer func() {
		c.hub.unregisterClient(c)
		c.conn.Close()
		c.hub.broadcast(wsEnvelope{Type: "member_left", Data: mustMarshal(map[string]uint64{"user_id": c.userID})}, c.userID)
		c.hub.broadcast(wsEnvelope{Type: "state", Data: mustMarshal(c.hub.buildStatePayload())}, c.userID)
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		var env wsEnvelope
		if err := json.Unmarshal(message, &env); err != nil {
			continue
		}
		c.hub.handleEvent(c, env)
	}
}

func (c *studyClient) writeLoop() {
	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			return
		}
	}
}

func (h *studyRoomHub) handleEvent(client *studyClient, env wsEnvelope) {
	switch env.Type {
	case "register_peer":
		var payload struct {
			PeerID string `json:"peer_id"`
		}
		_ = json.Unmarshal(env.Data, &payload)
		h.mu.Lock()
		client.peerID = strings.TrimSpace(payload.PeerID)
		h.mu.Unlock()
		h.broadcast(wsEnvelope{Type: "state", Data: mustMarshal(h.buildStatePayload())}, 0)

	case "call_request":
		var payload struct {
			TargetID uint64 `json:"target_id"`
		}
		if err := json.Unmarshal(env.Data, &payload); err != nil {
			return
		}
		h.handleCallRequest(client, payload.TargetID)

	case "call_accept":
		var payload struct {
			FromID uint64 `json:"from_id"`
		}
		if err := json.Unmarshal(env.Data, &payload); err != nil {
			return
		}
		h.handleCallAccept(client, payload.FromID)

	case "call_reject":
		var payload struct {
			FromID uint64 `json:"from_id"`
			Reason string `json:"reason"`
		}
		_ = json.Unmarshal(env.Data, &payload)
		if payload.Reason == "" {
			payload.Reason = "rejected"
		}
		h.handleCallReject(client, payload.FromID, payload.Reason)

	case "call_end":
		var payload struct {
			PartnerID uint64 `json:"partner_id"`
		}
		_ = json.Unmarshal(env.Data, &payload)
		h.handleCallEnd(client, payload.PartnerID)

	case "chat":
		var payload struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(env.Data, &payload); err != nil {
			return
		}
		content := strings.TrimSpace(payload.Content)
		if content == "" {
			return
		}
		now := time.Now()
		db := database.GetDB()
		chat := models.ChatMessage{
			SessionID: 0,
			RoomID:    h.roomID,
			UserID:    client.userID,
			Content:   content,
			MsgType:   0,
			SentAt:    now,
		}
		if err := db.Create(&chat).Error; err != nil {
			log.Println("store chat failed:", err)
		}
		sentAt := now.Format(time.RFC3339)
		msg := map[string]interface{}{
			"id":           chat.ID,
			"user_id":      client.userID,
			"display_name": client.displayName,
			"content":      content,
			"sent_at":      sentAt,
		}
		h.broadcast(wsEnvelope{Type: "chat", Data: mustMarshal(msg)}, 0)

	case "direct_chat":
		var payload struct {
			TargetID uint64 `json:"target_id"`
			Content  string `json:"content"`
		}
		if err := json.Unmarshal(env.Data, &payload); err != nil {
			return
		}
		content := strings.TrimSpace(payload.Content)
		if content == "" || payload.TargetID == 0 {
			return
		}
		target := h.clients[payload.TargetID]
		if target == nil {
			return
		}
		now := time.Now().Format(time.RFC3339)
		msg := map[string]interface{}{
			"from_id":      client.userID,
			"to_id":        payload.TargetID,
			"display_name": client.displayName,
			"content":      content,
			"sent_at":      now,
		}
		if target != nil {
			target.send <- wsEnvelope{Type: "direct_chat", Data: mustMarshal(msg)}
		}
		client.send <- wsEnvelope{Type: "direct_chat", Data: mustMarshal(msg)}

	case "state_request":
		client.send <- wsEnvelope{Type: "state", Data: mustMarshal(h.buildStatePayload())}
	}
}

func (h *studyRoomHub) handleCallRequest(caller *studyClient, targetID uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if caller.peerID == "" {
		caller.send <- wsEnvelope{Type: "call_denied", Data: mustMarshal(map[string]string{"reason": "missing_peer_id"})}
		return
	}
	if _, busy := h.busy[caller.userID]; busy {
		caller.send <- wsEnvelope{Type: "call_denied", Data: mustMarshal(map[string]string{"reason": "caller_busy"})}
		return
	}
	if partner, ok := h.busy[targetID]; ok && partner != 0 {
		caller.send <- wsEnvelope{Type: "call_denied", Data: mustMarshal(map[string]string{"reason": "target_busy"})}
		return
	}
	if pendingCaller, ok := h.pending[targetID]; ok && pendingCaller != caller.userID {
		caller.send <- wsEnvelope{Type: "call_denied", Data: mustMarshal(map[string]string{"reason": "target_busy"})}
		return
	}
	target := h.clients[targetID]
	if target == nil {
		caller.send <- wsEnvelope{Type: "call_denied", Data: mustMarshal(map[string]string{"reason": "target_offline"})}
		return
	}
	h.pending[targetID] = caller.userID
	target.send <- wsEnvelope{Type: "incoming_call", Data: mustMarshal(map[string]interface{}{
		"from_id":   caller.userID,
		"from_name": caller.displayName,
		"from_peer": caller.peerID,
	})}
}

func (h *studyRoomHub) handleCallAccept(callee *studyClient, fromID uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	callerID, ok := h.pending[callee.userID]
	if !ok || callerID != fromID {
		return
	}
	caller := h.clients[callerID]
	if caller == nil {
		delete(h.pending, callee.userID)
		return
	}
	if callee.peerID == "" {
		callee.send <- wsEnvelope{Type: "call_denied", Data: mustMarshal(map[string]string{"reason": "missing_peer_id"})}
		delete(h.pending, callee.userID)
		return
	}
	delete(h.pending, callee.userID)
	h.busy[callerID] = callee.userID
	h.busy[callee.userID] = callerID

	callData := map[string]interface{}{
		"caller_id":      caller.userID,
		"caller_name":    caller.displayName,
		"caller_peer_id": caller.peerID,
		"callee_id":      callee.userID,
		"callee_name":    callee.displayName,
		"callee_peer_id": callee.peerID,
	}
	caller.send <- wsEnvelope{Type: "call_start", Data: mustMarshal(callData)}
	callee.send <- wsEnvelope{Type: "call_start", Data: mustMarshal(callData)}
	h.broadcast(wsEnvelope{Type: "state", Data: mustMarshal(h.buildStatePayload())}, 0)
}

func (h *studyRoomHub) handleCallReject(callee *studyClient, fromID uint64, reason string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	callerID, ok := h.pending[callee.userID]
	if !ok || callerID != fromID {
		return
	}
	delete(h.pending, callee.userID)
	if caller := h.clients[callerID]; caller != nil {
		caller.send <- wsEnvelope{Type: "call_denied", Data: mustMarshal(map[string]string{"reason": reason})}
	}
}

func (h *studyRoomHub) handleCallEnd(client *studyClient, partnerID uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if currentPartner, ok := h.busy[client.userID]; !ok || currentPartner != partnerID {
		return
	}
	delete(h.busy, client.userID)
	delete(h.busy, partnerID)
	if partner := h.clients[partnerID]; partner != nil {
		partner.send <- wsEnvelope{Type: "call_ended", Data: mustMarshal(map[string]uint64{"partner_id": client.userID})}
	}
	client.send <- wsEnvelope{Type: "call_ended", Data: mustMarshal(map[string]uint64{"partner_id": partnerID})}
	h.broadcast(wsEnvelope{Type: "state", Data: mustMarshal(h.buildStatePayload())}, 0)
}

func mustMarshal(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
