package postgres

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/social"
	"gym-pro-2026-ptit/internal/infrastructure/database"

	"github.com/google/uuid"
)

// FollowRepository implementation
type followRepository struct {
	db *database.DB
}

func NewFollowRepository(db *database.DB) social.FollowRepository {
	return &followRepository{db: db}
}

// TODO: Implement all FollowRepository methods
func (r *followRepository) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	// TODO: Insert into follows table
	// Handle conflict (already following) gracefully
	// Use ON CONFLICT DO NOTHING or check first
	return nil
}

func (r *followRepository) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	// TODO: Delete from follows table
	return nil
}

func (r *followRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	// TODO: Check if follow relationship exists
	return false, nil
}

func (r *followRepository) GetFollowers(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]social.PostUser, int64, error) {
	// TODO: Query followers with pagination
	// JOIN users table to get user details
	// Query: SELECT users.* FROM follows JOIN users ON follows.follower_id = users.id WHERE follows.following_id = $1
	return nil, 0, nil
}

func (r *followRepository) GetFollowing(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]social.PostUser, int64, error) {
	// TODO: Query following with pagination
	// JOIN users table to get user details
	// Query: SELECT users.* FROM follows JOIN users ON follows.following_id = users.id WHERE follows.follower_id = $1
	return nil, 0, nil
}

func (r *followRepository) GetStats(ctx context.Context, userID uuid.UUID) (*social.FollowStats, error) {
	// TODO: Count followers and following
	// Query 1: SELECT COUNT(*) FROM follows WHERE following_id = $1 (followers)
	// Query 2: SELECT COUNT(*) FROM follows WHERE follower_id = $1 (following)
	return nil, nil
}

// PostRepository implementation
type postRepository struct {
	db *database.DB
}

func NewPostRepository(db *database.DB) social.PostRepository {
	return &postRepository{db: db}
}

// TODO: Implement all PostRepository methods
func (r *postRepository) Create(ctx context.Context, post *social.Post) error {
	// TODO: Insert into posts table
	// content_type: 'workout_plan' or 'meal_log'
	// content_id: references workout_plan_id or meal_log_id
	return nil
}

func (r *postRepository) GetByID(ctx context.Context, id uuid.UUID) (*social.Post, error) {
	// TODO: Query post with user details
	// JOIN users table
	return nil, nil
}

func (r *postRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]social.Post, int64, error) {
	// TODO: Query user's posts with pagination
	// Include user details
	return nil, 0, nil
}

func (r *postRepository) GetFeed(ctx context.Context, userID uuid.UUID, filter social.GetFeedFilter) ([]social.ActivityFeedItem, int64, error) {
	// TODO: Get activity feed for user
	// Query posts from users that current user follows
	// Include:
	// - Post details
	// - User who posted
	// - Whether current user has liked the post
	// Complex query:
	// SELECT posts.*, users.*, EXISTS(SELECT 1 FROM likes WHERE post_id = posts.id AND user_id = $1) as is_liked
	// FROM posts
	// JOIN users ON posts.user_id = users.id
	// WHERE posts.user_id IN (SELECT following_id FROM follows WHERE follower_id = $1)
	// ORDER BY posts.created_at DESC
	return nil, 0, nil
}

func (r *postRepository) Update(ctx context.Context, post *social.Post) error {
	// TODO: Update post (mainly caption)
	return nil
}

func (r *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: Delete post (cascade will delete likes and comments)
	return nil
}

func (r *postRepository) IncrementLikesCount(ctx context.Context, postID uuid.UUID) error {
	// TODO: UPDATE posts SET likes_count = likes_count + 1 WHERE id = $1
	return nil
}

func (r *postRepository) DecrementLikesCount(ctx context.Context, postID uuid.UUID) error {
	// TODO: UPDATE posts SET likes_count = likes_count - 1 WHERE id = $1
	// Ensure likes_count doesn't go below 0
	return nil
}

func (r *postRepository) IncrementCommentsCount(ctx context.Context, postID uuid.UUID) error {
	// TODO: UPDATE posts SET comments_count = comments_count + 1 WHERE id = $1
	return nil
}

func (r *postRepository) DecrementCommentsCount(ctx context.Context, postID uuid.UUID) error {
	// TODO: UPDATE posts SET comments_count = comments_count - 1 WHERE id = $1
	return nil
}

// LikeRepository implementation
type likeRepository struct {
	db *database.DB
}

func NewLikeRepository(db *database.DB) social.LikeRepository {
	return &likeRepository{db: db}
}

// TODO: Implement all LikeRepository methods
func (r *likeRepository) Create(ctx context.Context, like *social.Like) error {
	// TODO: Insert into likes table
	// Use ON CONFLICT to handle duplicate likes
	return nil
}

func (r *likeRepository) Delete(ctx context.Context, postID, userID uuid.UUID) error {
	// TODO: Delete from likes table
	return nil
}

func (r *likeRepository) Exists(ctx context.Context, postID, userID uuid.UUID) (bool, error) {
	// TODO: Check if like exists
	return false, nil
}

func (r *likeRepository) GetByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) ([]social.Like, int64, error) {
	// TODO: Query likes for a post with pagination
	return nil, 0, nil
}

// CommentRepository implementation
type commentRepository struct {
	db *database.DB
}

func NewCommentRepository(db *database.DB) social.CommentRepository {
	return &commentRepository{db: db}
}

// TODO: Implement all CommentRepository methods
func (r *commentRepository) Create(ctx context.Context, comment *social.Comment) error {
	// TODO: Insert into comments table
	return nil
}

func (r *commentRepository) GetByID(ctx context.Context, id uuid.UUID) (*social.Comment, error) {
	// TODO: Query comment with user details
	return nil, nil
}

func (r *commentRepository) GetByPostID(ctx context.Context, postID uuid.UUID, filter social.GetCommentsFilter) ([]social.Comment, int64, error) {
	// TODO: Query comments for a post with pagination
	// JOIN users table for user details
	// Order by created_at DESC (newest first) or ASC (oldest first)
	return nil, 0, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *social.Comment) error {
	// TODO: Update comment content
	return nil
}

func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: Delete comment
	return nil
}
