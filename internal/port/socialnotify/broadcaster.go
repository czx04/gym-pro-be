package socialnotify

import "github.com/google/uuid"

type NotificationPayload struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Meta      string `json:"meta"`
	DayGroup  string `json:"dayGroup"`
	IsRead    bool   `json:"isRead"`
	CreatedAt string `json:"createdAt"`
}

type RealtimeCommentAuthor struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
}

type RealtimeCommentMedia struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type RealtimeComment struct {
	ID               string                 `json:"id"`
	PostID           string                 `json:"postId"`
	ParentID         *string                `json:"parentId"`
	Depth            int                    `json:"depth"`
	Path             string                 `json:"path"`
	DirectReplyCount int                    `json:"directReplyCount"`
	PreviewReplies   []RealtimeComment      `json:"previewReplies"`
	Author           RealtimeCommentAuthor  `json:"author"`
	Content          string                 `json:"content"`
	Media            []RealtimeCommentMedia `json:"media"`
	IsDeleted        bool                   `json:"isDeleted"`
	CreatedAt        string                 `json:"createdAt"`
}

type CommentCreatedPayload struct {
	PostID  string          `json:"postId"`
	Comment RealtimeComment `json:"comment"`
}

type CommentDeletedPayload struct {
	PostID          string  `json:"postId"`
	CommentID       string  `json:"commentId"`
	ParentID        *string `json:"parentId,omitempty"`
	DeletedByUserID string  `json:"deletedByUserId"`
}

type Broadcaster interface {
	PublishNotificationCreated(userID uuid.UUID, n NotificationPayload)
	PublishUnread(userID uuid.UUID, unread int64)
	PublishCommentCreated(p CommentCreatedPayload)
	PublishCommentDeleted(p CommentDeletedPayload)
}
