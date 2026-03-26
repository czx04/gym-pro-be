package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	userdomain "gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	tokenQueryKey       = "token"
)

// Service allows other modules to hook into WebSocket lifecycle events.
type Service interface {
	OnConnect(ctx context.Context, userID uuid.UUID, conn *ws.Conn)
	OnDisconnect(ctx context.Context, userID uuid.UUID, conn *ws.Conn)
	OnMessage(ctx context.Context, userID uuid.UUID, messageType int, payload []byte)
}

type connection struct {
	conn *ws.Conn
	mu   sync.Mutex
}

// Hub manages WebSocket connections grouped by user ID.
type Hub struct {
	jwtManager *auth.JWTManager
	userRepo   userdomain.Repository
	upgrader   ws.Upgrader

	mu          sync.RWMutex
	connections map[uuid.UUID]map[*ws.Conn]*connection
	services    map[string]Service
}

// NewHub creates a WebSocket hub that authenticates users via JWT.
func NewHub(jwtManager *auth.JWTManager, userRepo userdomain.Repository) *Hub {
	return &Hub{
		jwtManager: jwtManager,
		userRepo:   userRepo,
		upgrader: ws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
		connections: make(map[uuid.UUID]map[*ws.Conn]*connection),
		services:    make(map[string]Service),
	}
}

// RegisterRoutes mounts WebSocket routes into a Gin group.
func (h *Hub) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/ws", h.Handle)
}

// RegisterService registers a service that listens to WebSocket events.
func (h *Hub) RegisterService(name string, service Service) {
	if service == nil || name == "" {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.services[name] = service
}

// Handle authenticates, upgrades and keeps a connection alive.
func (h *Hub) Handle(c *gin.Context) {
	userID, err := h.authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("websocket upgrade failed", "err", err)
		return
	}

	h.addConnection(userID, conn)
	h.notifyConnect(c.Request.Context(), userID, conn)
	logger.Info("websocket connected", "user_id", userID)

	defer func() {
		h.removeConnection(userID, conn)
		h.notifyDisconnect(context.Background(), userID, conn)
		_ = conn.Close()
		logger.Info("websocket disconnected", "user_id", userID)
	}()

	conn.SetReadLimit(64 * 1024)
	_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(_ string) error {
		return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	for {
		messageType, payload, readErr := conn.ReadMessage()
		if readErr != nil {
			if ws.IsUnexpectedCloseError(readErr, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				logger.Warn("websocket read failed", "err", readErr, "user_id", userID)
			}
			return
		}
		h.notifyMessage(c.Request.Context(), userID, messageType, payload)
	}
}

// SendToUser sends raw bytes to all connections of a user.
func (h *Hub) SendToUser(userID uuid.UUID, messageType int, payload []byte) error {
	conns := h.userConnections(userID)
	if len(conns) == 0 {
		return fmt.Errorf("no active websocket connection for user %s", userID)
	}

	var sendErr error
	for _, c := range conns {
		if err := c.write(messageType, payload); err != nil {
			sendErr = errors.Join(sendErr, err)
		}
	}
	return sendErr
}

// SendJSONToUser marshals and sends JSON payload to all connections of a user.
func (h *Hub) SendJSONToUser(userID uuid.UUID, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return h.SendToUser(userID, ws.TextMessage, data)
}

// Broadcast sends raw bytes to all active connections.
func (h *Hub) Broadcast(messageType int, payload []byte) error {
	conns := h.allConnections()
	var sendErr error
	for _, c := range conns {
		if err := c.write(messageType, payload); err != nil {
			sendErr = errors.Join(sendErr, err)
		}
	}
	return sendErr
}

// BroadcastJSON marshals and sends JSON payload to all active connections.
func (h *Hub) BroadcastJSON(payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return h.Broadcast(ws.TextMessage, data)
}

// CloseUser closes and removes all connections for a user.
func (h *Hub) CloseUser(userID uuid.UUID) {
	conns := h.userConnections(userID)
	for _, c := range conns {
		_ = c.close()
	}

	h.mu.Lock()
	delete(h.connections, userID)
	h.mu.Unlock()
}

// Shutdown closes all active connections and clears internal state.
func (h *Hub) Shutdown() {
	conns := h.allConnections()
	for _, c := range conns {
		_ = c.close()
	}

	h.mu.Lock()
	h.connections = make(map[uuid.UUID]map[*ws.Conn]*connection)
	h.mu.Unlock()
}

func (h *Hub) authenticate(c *gin.Context) (uuid.UUID, error) {
	token := h.extractToken(c)
	if token == "" {
		return uuid.Nil, errors.New("missing token")
	}

	claims, err := h.jwtManager.ValidateAccessToken(token)
	if err != nil {
		return uuid.Nil, err
	}

	user, err := h.userRepo.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		return uuid.Nil, err
	}
	if user == nil {
		return uuid.Nil, errors.New("user not found")
	}

	return claims.UserID, nil
}

func (h *Hub) extractToken(c *gin.Context) string {
	authHeader := c.GetHeader(authorizationHeader)
	if strings.HasPrefix(authHeader, bearerPrefix) {
		if token := strings.TrimPrefix(authHeader, bearerPrefix); token != "" {
			return token
		}
	}
	return c.Query(tokenQueryKey)
}

func (h *Hub) addConnection(userID uuid.UUID, wsConn *ws.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.connections[userID]; !ok {
		h.connections[userID] = make(map[*ws.Conn]*connection)
	}
	h.connections[userID][wsConn] = &connection{conn: wsConn}
}

func (h *Hub) removeConnection(userID uuid.UUID, wsConn *ws.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	userConns, ok := h.connections[userID]
	if !ok {
		return
	}
	delete(userConns, wsConn)
	if len(userConns) == 0 {
		delete(h.connections, userID)
	}
}

func (h *Hub) userConnections(userID uuid.UUID) []*connection {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userConns, ok := h.connections[userID]
	if !ok {
		return nil
	}

	result := make([]*connection, 0, len(userConns))
	for _, conn := range userConns {
		result = append(result, conn)
	}
	return result
}

func (h *Hub) allConnections() []*connection {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]*connection, 0)
	for _, userConns := range h.connections {
		for _, conn := range userConns {
			result = append(result, conn)
		}
	}
	return result
}

func (h *Hub) notifyConnect(ctx context.Context, userID uuid.UUID, conn *ws.Conn) {
	for _, svc := range h.snapshotServices() {
		svc.OnConnect(ctx, userID, conn)
	}
}

func (h *Hub) notifyDisconnect(ctx context.Context, userID uuid.UUID, conn *ws.Conn) {
	for _, svc := range h.snapshotServices() {
		svc.OnDisconnect(ctx, userID, conn)
	}
}

func (h *Hub) notifyMessage(ctx context.Context, userID uuid.UUID, messageType int, payload []byte) {
	for _, svc := range h.snapshotServices() {
		svc.OnMessage(ctx, userID, messageType, payload)
	}
}

func (h *Hub) snapshotServices() []Service {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]Service, 0, len(h.services))
	for _, svc := range h.services {
		result = append(result, svc)
	}
	return result
}

func (c *connection) write(messageType int, payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(messageType, payload)
}

func (c *connection) close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.Close()
}
