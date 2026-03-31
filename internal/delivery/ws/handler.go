package ws

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"gym-pro-2026-ptit/internal/config"
	userdomain "gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/infrastructure/auth"
	"gym-pro-2026-ptit/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const maxClientMessageBytes = 4096

type Handler struct {
	hub       *Hub
	jwt       *auth.JWTManager
	userRepo  userdomain.Repository
	serverCfg *config.ServerConfig
	upgrader  websocket.Upgrader
}

func NewHandler(hub *Hub, jwt *auth.JWTManager, userRepo userdomain.Repository, cfg *config.Config) *Handler {
	h := &Handler{
		hub:       hub,
		jwt:       jwt,
		userRepo:  userRepo,
		serverCfg: &cfg.Server,
	}
	h.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     h.checkOrigin,
	}
	return h
}

func (h *Handler) checkOrigin(r *http.Request) bool {
	if len(h.serverCfg.AllowedOrigins) == 0 {
		return true
	}
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return true
	}
	for _, o := range h.serverCfg.AllowedOrigins {
		o = strings.TrimSpace(o)
		if o == "*" || strings.EqualFold(o, origin) {
			return true
		}
	}
	if strings.EqualFold(strings.TrimSpace(h.serverCfg.GinMode), "debug") {
		if isDevMobileOrigin(origin) {
			return true
		}
	}
	return false
}

func isDevMobileOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	scheme := strings.ToLower(strings.TrimSpace(u.Scheme))
	host := strings.ToLower(strings.TrimSpace(u.Hostname()))
	if host == "" {
		return false
	}

	if scheme == "exp" || scheme == "exps" {
		return true
	}

	if scheme != "http" && scheme != "https" {
		return false
	}

	if host == "localhost" || host == "127.0.0.1" {
		return true
	}

	if strings.HasPrefix(host, "10.") || strings.HasPrefix(host, "192.168.") {
		return true
	}

	if strings.HasPrefix(host, "172.") {
		parts := strings.Split(host, ".")
		if len(parts) >= 2 {
			second, err := strconv.Atoi(parts[1])
			if err == nil && second >= 16 && second <= 31 {
				return true
			}
		}
	}

	return false
}

func (h *Handler) Handle(c *gin.Context) {
	token := strings.TrimSpace(c.Query("access_token"))
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"code": "UNAUTHORIZED", "message": "missing token"}})
		return
	}
	claims, err := h.jwt.ValidateAccessToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"code": "UNAUTHORIZED", "message": "invalid token"}})
		return
	}
	userID := claims.UserID
	u, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil || u == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"code": "UNAUTHORIZED", "message": "user not found"}})
		return
	}
	_ = u

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("websocket upgrade failed", "err", err)
		return
	}
	conn.SetReadLimit(maxClientMessageBytes)
	h.hub.Register(userID, conn)
	defer func() {
		h.hub.Unregister(userID, conn)
		_ = conn.Close()
	}()

	for {
		mt, payload, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if mt != websocket.TextMessage {
			continue
		}
		var m struct {
			Type string `json:"type"`
		}
		if json.Unmarshal(payload, &m) != nil {
			continue
		}
		if m.Type == "ping" {
			_ = h.hub.WriteText(conn, []byte(`{"v":1,"type":"pong","payload":{}}`))
		}
	}
}
