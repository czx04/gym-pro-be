package ws

import (
	"encoding/json"
	"sync"

	"gym-pro-2026-ptit/internal/port/socialnotify"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	mu        sync.RWMutex
	wmu       sync.Mutex
	userConns map[uuid.UUID]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{userConns: make(map[uuid.UUID]map[*websocket.Conn]struct{})}
}

func (h *Hub) Register(userID uuid.UUID, c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	m, ok := h.userConns[userID]
	if !ok {
		m = make(map[*websocket.Conn]struct{})
		h.userConns[userID] = m
	}
	m[c] = struct{}{}
}

func (h *Hub) Unregister(userID uuid.UUID, c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if m, ok := h.userConns[userID]; ok {
		delete(m, c)
		if len(m) == 0 {
			delete(h.userConns, userID)
		}
	}
}

func (h *Hub) sendJSON(userID uuid.UUID, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	h.mu.RLock()
	m := h.userConns[userID]
	conns := make([]*websocket.Conn, 0, len(m))
	for c := range m {
		conns = append(conns, c)
	}
	h.mu.RUnlock()
	h.wmu.Lock()
	defer h.wmu.Unlock()
	for _, c := range conns {
		_ = c.WriteMessage(websocket.TextMessage, data)
	}
}

func (h *Hub) WriteText(c *websocket.Conn, data []byte) error {
	h.wmu.Lock()
	defer h.wmu.Unlock()
	return c.WriteMessage(websocket.TextMessage, data)
}

func (h *Hub) PublishNotificationCreated(userID uuid.UUID, n socialnotify.NotificationPayload) {
	h.sendJSON(userID, envelope{V: 1, Type: "notification.created", Payload: notifCreatedBody{Notification: n}})
}

func (h *Hub) PublishUnread(userID uuid.UUID, unread int64) {
	h.sendJSON(userID, envelope{V: 1, Type: "notification.unread", Payload: unreadBody{UnreadCount: unread}})
}

func (h *Hub) PublishCommentCreated(p socialnotify.CommentCreatedPayload) {
	h.broadcastAll("comment.created", p)
}

func (h *Hub) PublishCommentDeleted(p socialnotify.CommentDeletedPayload) {
	h.broadcastAll("comment.deleted", p)
}

func (h *Hub) broadcastAll(eventType string, payload interface{}) {
	data, err := json.Marshal(envelope{V: 1, Type: eventType, Payload: payload})
	if err != nil {
		return
	}
	h.mu.RLock()
	conns := make([]*websocket.Conn, 0)
	for _, m := range h.userConns {
		for c := range m {
			conns = append(conns, c)
		}
	}
	h.mu.RUnlock()
	h.wmu.Lock()
	defer h.wmu.Unlock()
	for _, c := range conns {
		_ = c.WriteMessage(websocket.TextMessage, data)
	}
}

type envelope struct {
	V       int         `json:"v"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type notifCreatedBody struct {
	Notification socialnotify.NotificationPayload `json:"notification"`
}

type unreadBody struct {
	UnreadCount int64 `json:"unreadCount"`
}
