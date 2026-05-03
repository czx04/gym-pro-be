package bootstrap

import (
	"context"
	"gym-pro-2026-ptit/internal/delivery/http/router"
	"gym-pro-2026-ptit/internal/delivery/http/websocket"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/internal/port/socialnotify"

	"github.com/google/uuid"

	"go.uber.org/fx"
)

type socialBroadcaster struct {
	hub *websocket.Hub
}

func ProvideSocialBroadcaster(h *websocket.Hub) socialnotify.Broadcaster {
	return &socialBroadcaster{hub: h}
}

func (b *socialBroadcaster) PublishNotificationCreated(userID uuid.UUID, n socialnotify.NotificationPayload) {
	_ = b.hub.SendJSONToUser(userID, map[string]any{
		"v":    1,
		"type": "notification.created",
		"payload": map[string]any{
			"notification": n,
		},
	})
}

func (b *socialBroadcaster) PublishUnread(userID uuid.UUID, unread int64) {
	_ = b.hub.SendJSONToUser(userID, map[string]any{
		"v":    1,
		"type": "notification.unread",
		"payload": map[string]any{
			"unreadCount": unread,
		},
	})
}

func (b *socialBroadcaster) PublishCommentCreated(p socialnotify.CommentCreatedPayload) {
	_ = b.hub.BroadcastJSON(map[string]any{
		"v":       1,
		"type":    "comment.created",
		"payload": p,
	})
}

func (b *socialBroadcaster) PublishCommentUpdated(p socialnotify.CommentUpdatedPayload) {
	_ = b.hub.BroadcastJSON(map[string]any{
		"v":       1,
		"type":    "comment.updated",
		"payload": p,
	})
}

func (b *socialBroadcaster) PublishCommentDeleted(p socialnotify.CommentDeletedPayload) {
	_ = b.hub.BroadcastJSON(map[string]any{
		"v":       1,
		"type":    "comment.deleted",
		"payload": p,
	})
}

func RegisterSocialWebSocketHook(lc fx.Lifecycle, _ *router.Router) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("social realtime broadcaster wired to /api/v1/ws")
			return nil
		},
	})
}
