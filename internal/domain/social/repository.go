package social

import (
	"context"

	"github.com/google/uuid"
)

// FollowRepository defines the interface for follow data access
type FollowRepository interface {
	// Follow creates a follow relationship
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error

	// Unfollow removes a follow relationship
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error

	// IsFollowing checks if a user is following another user
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)

	// GetFollowers retrieves followers of a user
	GetFollowers(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]PostUser, int64, error)

	// GetFollowing retrieves users that a user is following
	GetFollowing(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]PostUser, int64, error)

	// GetStats retrieves follow statistics
	GetStats(ctx context.Context, userID uuid.UUID) (*FollowStats, error)
}

// PostRepository defines the interface for post data access
type PostRepository interface {
	// Create creates a new post
	Create(ctx context.Context, post *Post) error

	// CreateWithMedia creates a post and attaches media atomically
	CreateWithMedia(ctx context.Context, post *Post, media []PostMedia) error

	// GetByID retrieves a post by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Post, error)

	// GetByUserID retrieves posts by a user
	GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]Post, int64, error)

	// GetFeed retrieves activity feed for a user (posts from followed users)
	GetFeed(ctx context.Context, userID uuid.UUID, filter GetFeedFilter) ([]ActivityFeedItem, int64, error)

	// GetMediaByPostIDs retrieves media grouped by post IDs
	GetMediaByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]PostMedia, error)

	// Update updates a post
	Update(ctx context.Context, post *Post) error

	// Delete deletes a post
	Delete(ctx context.Context, id uuid.UUID) error

	// IncrementLikesCount increments likes count
	IncrementLikesCount(ctx context.Context, postID uuid.UUID) error

	// DecrementLikesCount decrements likes count
	DecrementLikesCount(ctx context.Context, postID uuid.UUID) error

	// IncrementCommentsCount increments comments count
	IncrementCommentsCount(ctx context.Context, postID uuid.UUID) error

	// DecrementCommentsCount decrements comments count
	DecrementCommentsCount(ctx context.Context, postID uuid.UUID) error
}

// MediaAssetRepository defines the interface for social media assets
type MediaAssetRepository interface {
	// CreatePending creates a pending media asset after signature is issued
	CreatePending(ctx context.Context, asset *SocialMediaAsset) error

	// Confirm sets asset status to ready and stores metadata
	Confirm(ctx context.Context, userID uuid.UUID, publicID string, secureURL *string, bytes *int64) error
}

// LikeRepository defines the interface for like data access
type LikeRepository interface {
	// Create creates a like
	Create(ctx context.Context, like *Like) error

	// Delete deletes a like
	Delete(ctx context.Context, postID, userID uuid.UUID) error

	// Exists checks if a like exists
	Exists(ctx context.Context, postID, userID uuid.UUID) (bool, error)

	// GetByPostID retrieves likes for a post
	GetByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) ([]Like, int64, error)
}

// CommentRepository defines the interface for comment data access
type CommentRepository interface {
	// Create creates a comment
	Create(ctx context.Context, comment *Comment) error

	// GetByID retrieves a comment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Comment, error)

	// GetByPostID retrieves comments for a post
	GetByPostID(ctx context.Context, postID uuid.UUID, filter GetCommentsFilter) ([]Comment, int64, error)

	// Update updates a comment
	Update(ctx context.Context, comment *Comment) error

	// Delete deletes a comment
	Delete(ctx context.Context, id uuid.UUID) error
}
