package social

import (
	"time"

	"github.com/google/uuid"
)

// Follow represents a follower relationship
type Follow struct {
	FollowerID  uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// Post represents a shared workout or meal log
type Post struct {
	ID            uuid.UUID   `json:"id"`
	UserID        uuid.UUID   `json:"user_id"`
	ContentType   string      `json:"content_type"`
	ContentID     *uuid.UUID  `json:"content_id,omitempty"`
	Caption       *string     `json:"caption,omitempty"`
	Feeling       *string     `json:"feeling,omitempty"`
	LocationName  *string     `json:"location_name,omitempty"`
	Hashtags      []string    `json:"hashtags,omitempty"`
	LikesCount    int         `json:"likes_count"`
	CommentsCount int         `json:"comments_count"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DeletedAt     *time.Time  `json:"deleted_at,omitempty"`
	User          *PostUser   `json:"user,omitempty"` // Basic user info
	Media         []PostMedia `json:"media,omitempty"`
}

type PostMedia struct {
	PostID       uuid.UUID `json:"post_id"`
	PublicID     string    `json:"public_id"`
	ResourceType string    `json:"resource_type"`
	SecureURL    *string   `json:"secure_url,omitempty"`
	OrderIndex   int       `json:"order_index"`
}

type SocialMediaAsset struct {
	PublicID     string     `json:"public_id"`
	UserID       uuid.UUID  `json:"user_id"`
	ResourceType string     `json:"resource_type"`
	SecureURL    *string    `json:"secure_url,omitempty"`
	Bytes        *int64     `json:"bytes,omitempty"`
	Status       string     `json:"status"`
	PostID       *uuid.UUID `json:"post_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	ConfirmedAt  *time.Time `json:"confirmed_at,omitempty"`
	AttachedAt   *time.Time `json:"attached_at,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// PostUser represents basic user info in a post
type PostUser struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
}

type PostLocation struct {
	Name string `json:"name"`
}

// Like represents a like on a post
type Like struct {
	ID        uuid.UUID `json:"id"`
	PostID    uuid.UUID `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Comment represents a comment on a post
type Comment struct {
	ID              uuid.UUID  `json:"id"`
	PostID          uuid.UUID  `json:"post_id"`
	UserID          uuid.UUID  `json:"user_id"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id,omitempty"`
	Content         string     `json:"content"`
	ReplyCount      int        `json:"reply_count"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	User            *PostUser  `json:"user,omitempty"`
}

// CreatePostInput represents input for creating a post
type CreatePostInput struct {
	ContentType string                 `json:"content_type,omitempty"`
	ContentID   *uuid.UUID             `json:"content_id,omitempty"`
	Caption     *string                `json:"caption,omitempty" validate:"omitempty,max=2000"`
	Media       []CreatePostMediaInput `json:"media,omitempty" validate:"omitempty,dive"`
	Feeling     *string                `json:"feeling,omitempty" validate:"omitempty,max=100"`
	Location    *CreatePostLocation    `json:"location,omitempty" validate:"omitempty"`
	Hashtags    []string               `json:"hashtags,omitempty" validate:"omitempty,dive,max=50"`
}

type CreatePostLocation struct {
	Name string `json:"name" validate:"required,max=255"`
}

type CreatePostMediaInput struct {
	PublicID     string `json:"public_id" validate:"required"`
	ResourceType string `json:"resource_type" validate:"required,oneof=image video"`
}

// CreateCommentInput represents input for creating a comment
type CreateCommentInput struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// UpdateCommentInput represents input for updating a comment
type UpdateCommentInput struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// GetFeedFilter represents filters for getting feed
type GetFeedFilter struct {
	Page     int
	PageSize int
}

// GetCommentsFilter represents filters for getting comments
type GetCommentsFilter struct {
	Page            int
	PageSize        int
	ParentCommentID *uuid.UUID
}

// ActivityFeedItem represents an item in the activity feed
type ActivityFeedItem struct {
	Post            *Post     `json:"post"`
	IsLiked         bool      `json:"is_liked"`
	IsInterested    bool      `json:"is_interested_by_me"`
	IsNotInterested bool      `json:"is_not_interested_by_me"`
	CreatedAt       time.Time `json:"created_at"`
}

// FollowStats represents follow statistics
type FollowStats struct {
	FollowersCount int `json:"followers_count"`
	FollowingCount int `json:"following_count"`
}

type PostPreference struct {
	UserID     uuid.UUID `json:"user_id"`
	PostID     uuid.UUID `json:"post_id"`
	Preference string    `json:"preference"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostReport struct {
	ID          uuid.UUID `json:"id"`
	PostID      uuid.UUID `json:"post_id"`
	ReporterID  uuid.UUID `json:"reporter_id"`
	Reason      string    `json:"reason"`
	Description *string   `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserBlock struct {
	BlockerID uuid.UUID `json:"blocker_id"`
	BlockedID uuid.UUID `json:"blocked_id"`
	CreatedAt time.Time `json:"created_at"`
}
