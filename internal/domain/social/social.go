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
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	ContentType   string     `json:"content_type"` // workout_plan, meal_log
	ContentID     uuid.UUID  `json:"content_id"` // References workout_plan_id or meal_log_id
	Caption       *string    `json:"caption,omitempty"`
	LikesCount    int        `json:"likes_count"`
	CommentsCount int        `json:"comments_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	User          *PostUser  `json:"user,omitempty"` // Basic user info
}

// PostUser represents basic user info in a post
type PostUser struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
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
	ID        uuid.UUID `json:"id"`
	PostID    uuid.UUID `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *PostUser `json:"user,omitempty"`
}

// CreatePostInput represents input for creating a post
type CreatePostInput struct {
	ContentType string     `json:"content_type" validate:"required,oneof=workout_plan meal_log"`
	ContentID   uuid.UUID  `json:"content_id" validate:"required"`
	Caption     *string    `json:"caption,omitempty" validate:"omitempty,max=2000"`
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
	Page     int
	PageSize int
}

// ActivityFeedItem represents an item in the activity feed
type ActivityFeedItem struct {
	Post      *Post     `json:"post"`
	IsLiked   bool      `json:"is_liked"`
	CreatedAt time.Time `json:"created_at"`
}

// FollowStats represents follow statistics
type FollowStats struct {
	FollowersCount int `json:"followers_count"`
	FollowingCount int `json:"following_count"`
}
