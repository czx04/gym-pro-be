package postgres

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/social"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/pkg/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// FollowRepository implementation
type followRepository struct {
	db *database.DB
}

func NewFollowRepository(db *database.DB) social.FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `
		INSERT INTO follows (follower_id, following_id)
		VALUES ($1, $2)
		ON CONFLICT (follower_id, following_id) DO NOTHING
	`

	if _, err := r.db.Exec(ctx, query, followerID, followingID); err != nil {
		return errors.DatabaseError("follow user", err)
	}

	return nil
}

func (r *followRepository) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`

	if _, err := r.db.Exec(ctx, query, followerID, followingID); err != nil {
		return errors.DatabaseError("unfollow user", err)
	}

	return nil
}

func (r *followRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)`

	var isFollowing bool
	if err := r.db.QueryRow(ctx, query, followerID, followingID).Scan(&isFollowing); err != nil {
		return false, errors.DatabaseError("check following relationship", err)
	}

	return isFollowing, nil
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
	stats := &social.FollowStats{}

	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM follows WHERE following_id = $1`, userID).Scan(&stats.FollowersCount); err != nil {
		return nil, errors.DatabaseError("count followers", err)
	}

	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM follows WHERE follower_id = $1`, userID).Scan(&stats.FollowingCount); err != nil {
		return nil, errors.DatabaseError("count following", err)
	}

	return stats, nil
}

func (r *followRepository) HasBlockRelation(ctx context.Context, userAID, userBID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM user_blocks ub
			WHERE (ub.blocker_id = $1 AND ub.blocked_id = $2)
			   OR (ub.blocker_id = $2 AND ub.blocked_id = $1)
		)
	`
	var hasRelation bool
	if err := r.db.QueryRow(ctx, query, userAID, userBID).Scan(&hasRelation); err != nil {
		return false, errors.DatabaseError("check block relation", err)
	}
	return hasRelation, nil
}

// PostRepository implementation
type postRepository struct {
	db *database.DB
}

func NewPostRepository(db *database.DB) social.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *social.Post) error {
	query := `
		INSERT INTO posts (
			id, user_id, content_type, content_id, caption, feeling, location_name, hashtags,
			likes_count, comments_count, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.Exec(ctx, query,
		post.ID,
		post.UserID,
		post.ContentType,
		post.ContentID,
		post.Caption,
		post.Feeling,
		post.LocationName,
		post.Hashtags,
		post.LikesCount,
		post.CommentsCount,
		post.CreatedAt,
		post.UpdatedAt,
	)
	if err != nil {
		return errors.DatabaseError("create post", err)
	}

	return nil
}

func (r *postRepository) CreateWithMedia(ctx context.Context, post *social.Post, media []social.PostMedia) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.DatabaseError("begin create post transaction", err)
	}
	defer tx.Rollback(ctx)

	insertPostQuery := `
		INSERT INTO posts (
			id, user_id, content_type, content_id, caption, feeling, location_name, hashtags,
			likes_count, comments_count, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = tx.Exec(ctx, insertPostQuery,
		post.ID,
		post.UserID,
		post.ContentType,
		post.ContentID,
		post.Caption,
		post.Feeling,
		post.LocationName,
		post.Hashtags,
		post.LikesCount,
		post.CommentsCount,
		post.CreatedAt,
		post.UpdatedAt,
	)
	if err != nil {
		return errors.DatabaseError("insert post", err)
	}

	for _, m := range media {
		var resourceType string
		attachAssetQuery := `
			UPDATE social_media_assets
			SET status = 'attached',
				post_id = $3,
				attached_at = NOW(),
				updated_at = NOW()
			WHERE public_id = $1
			  AND user_id = $2
			  AND status = 'ready'
			RETURNING resource_type
		`

		err := tx.QueryRow(ctx, attachAssetQuery, m.PublicID, post.UserID, post.ID).Scan(&resourceType)
		if err != nil {
			if err == pgx.ErrNoRows {
				return errors.Conflict("media asset not ready or not owned by current user")
			}
			return errors.DatabaseError("attach media asset", err)
		}

		insertPostMediaQuery := `
			INSERT INTO post_media (post_id, public_id, resource_type, order_index)
			VALUES ($1, $2, $3, $4)
		`

		_, err = tx.Exec(ctx, insertPostMediaQuery, post.ID, m.PublicID, resourceType, m.OrderIndex)
		if err != nil {
			return errors.DatabaseError("insert post media", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.DatabaseError("commit create post transaction", err)
	}

	return nil
}

func (r *postRepository) GetByID(ctx context.Context, id uuid.UUID) (*social.Post, error) {
	query := `
		SELECT
			p.id,
			p.user_id,
			p.content_type,
			p.content_id,
			p.caption,
			p.feeling,
			p.location_name,
			p.hashtags,
			p.likes_count,
			p.comments_count,
			p.created_at,
			p.updated_at,
			u.id,
			u.name,
			u.avatar_url
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.id = $1
		  AND p.deleted_at IS NULL
	`

	var post social.Post
	post.User = &social.PostUser{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.ContentType,
		&post.ContentID,
		&post.Caption,
		&post.Feeling,
		&post.LocationName,
		&post.Hashtags,
		&post.LikesCount,
		&post.CommentsCount,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.User.ID,
		&post.User.Name,
		&post.User.AvatarURL,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("post")
		}
		return nil, errors.DatabaseError("get post by id", err)
	}

	return &post, nil
}

func (r *postRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]social.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM posts WHERE user_id = $1 AND deleted_at IS NULL`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, errors.DatabaseError("count user posts", err)
	}

	query := `
		SELECT
			p.id,
			p.user_id,
			p.content_type,
			p.content_id,
			p.caption,
			p.feeling,
			p.location_name,
			p.hashtags,
			p.likes_count,
			p.comments_count,
			p.created_at,
			p.updated_at,
			u.id,
			u.name,
			u.avatar_url
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.user_id = $1
		  AND p.deleted_at IS NULL
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.DatabaseError("get posts by user id", err)
	}
	defer rows.Close()

	posts := make([]social.Post, 0)
	for rows.Next() {
		var post social.Post
		post.User = &social.PostUser{}

		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.ContentType,
			&post.ContentID,
			&post.Caption,
			&post.Feeling,
			&post.LocationName,
			&post.Hashtags,
			&post.LikesCount,
			&post.CommentsCount,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.User.ID,
			&post.User.Name,
			&post.User.AvatarURL,
		); err != nil {
			return nil, 0, errors.DatabaseError("scan user posts", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, errors.DatabaseError("iterate user posts", err)
	}

	return posts, total, nil
}

func (r *postRepository) GetFeed(ctx context.Context, userID uuid.UUID, filter social.GetFeedFilter) ([]social.ActivityFeedItem, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize

	countQuery := `
		SELECT COUNT(*)
		FROM posts p
		WHERE p.deleted_at IS NULL
		  AND NOT EXISTS (
			SELECT 1
			FROM post_preferences pp
			WHERE pp.user_id = $1
			  AND pp.post_id = p.id
			  AND pp.preference = 'not_interested'
		  )
		  AND NOT EXISTS (
			SELECT 1
			FROM user_blocks ub
			WHERE (ub.blocker_id = $1 AND ub.blocked_id = p.user_id)
			   OR (ub.blocker_id = p.user_id AND ub.blocked_id = $1)
		  )
	`

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, errors.DatabaseError("count feed posts", err)
	}

	query := `
		SELECT
			p.id,
			p.user_id,
			p.content_type,
			p.content_id,
			p.caption,
			p.feeling,
			p.location_name,
			p.hashtags,
			p.likes_count,
			p.comments_count,
			p.created_at,
			p.updated_at,
			u.id,
			u.name,
			u.avatar_url,
			EXISTS(
				SELECT 1
				FROM likes l
				WHERE l.post_id = p.id
				  AND l.user_id = $1
			) AS is_liked,
			COALESCE(pp.preference = 'interested', FALSE) AS is_interested,
			COALESCE(pp.preference = 'not_interested', FALSE) AS is_not_interested
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN post_preferences pp
			ON pp.post_id = p.id
		   AND pp.user_id = $1
		WHERE p.deleted_at IS NULL
		  AND NOT EXISTS (
			SELECT 1
			FROM post_preferences ppx
			WHERE ppx.user_id = $1
			  AND ppx.post_id = p.id
			  AND ppx.preference = 'not_interested'
		  )
		  AND NOT EXISTS (
			SELECT 1
			FROM user_blocks ub
			WHERE (ub.blocker_id = $1 AND ub.blocked_id = p.user_id)
			   OR (ub.blocker_id = p.user_id AND ub.blocked_id = $1)
		  )
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, filter.PageSize, offset)
	if err != nil {
		return nil, 0, errors.DatabaseError("get feed", err)
	}
	defer rows.Close()

	feed := make([]social.ActivityFeedItem, 0)
	for rows.Next() {
		item := social.ActivityFeedItem{Post: &social.Post{User: &social.PostUser{}}}

		if err := rows.Scan(
			&item.Post.ID,
			&item.Post.UserID,
			&item.Post.ContentType,
			&item.Post.ContentID,
			&item.Post.Caption,
			&item.Post.Feeling,
			&item.Post.LocationName,
			&item.Post.Hashtags,
			&item.Post.LikesCount,
			&item.Post.CommentsCount,
			&item.Post.CreatedAt,
			&item.Post.UpdatedAt,
			&item.Post.User.ID,
			&item.Post.User.Name,
			&item.Post.User.AvatarURL,
			&item.IsLiked,
			&item.IsInterested,
			&item.IsNotInterested,
		); err != nil {
			return nil, 0, errors.DatabaseError("scan feed item", err)
		}

		item.CreatedAt = item.Post.CreatedAt
		feed = append(feed, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, errors.DatabaseError("iterate feed", err)
	}

	return feed, total, nil
}

func (r *postRepository) GetMediaByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]social.PostMedia, error) {
	grouped := make(map[uuid.UUID][]social.PostMedia)
	if len(postIDs) == 0 {
		return grouped, nil
	}

	query := `
		SELECT pm.post_id, pm.public_id, pm.resource_type, sma.secure_url, pm.order_index
		FROM post_media pm
		LEFT JOIN social_media_assets sma ON sma.public_id = pm.public_id
		WHERE pm.post_id = $1
		ORDER BY pm.order_index ASC, pm.created_at ASC
	`

	for _, postID := range postIDs {
		rows, err := r.db.Query(ctx, query, postID)
		if err != nil {
			return nil, errors.DatabaseError("get post media", err)
		}

		items := make([]social.PostMedia, 0)
		for rows.Next() {
			var item social.PostMedia
			if err := rows.Scan(&item.PostID, &item.PublicID, &item.ResourceType, &item.SecureURL, &item.OrderIndex); err != nil {
				rows.Close()
				return nil, errors.DatabaseError("scan post media", err)
			}
			items = append(items, item)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, errors.DatabaseError("iterate post media", err)
		}
		rows.Close()

		grouped[postID] = items
	}

	return grouped, nil
}

func (r *postRepository) Update(ctx context.Context, post *social.Post) error {
	query := `
		UPDATE posts
		SET caption = $2,
			feeling = $3,
			location_name = $4,
			hashtags = $5,
			updated_at = NOW()
		WHERE id = $1
		  AND deleted_at IS NULL
	`
	result, err := r.db.Exec(ctx, query, post.ID, post.Caption, post.Feeling, post.LocationName, post.Hashtags)
	if err != nil {
		return errors.DatabaseError("update post", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound("post")
	}
	return nil
}

func (r *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE posts
		SET deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		  AND deleted_at IS NULL
	`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return errors.DatabaseError("delete post", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound("post")
	}
	return nil
}

func (r *postRepository) IncrementLikesCount(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET likes_count = likes_count + 1, updated_at = NOW() WHERE id = $1`
	if _, err := r.db.Exec(ctx, query, postID); err != nil {
		return errors.DatabaseError("increment likes count", err)
	}
	return nil
}

func (r *postRepository) DecrementLikesCount(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET likes_count = GREATEST(0, likes_count - 1), updated_at = NOW() WHERE id = $1`
	if _, err := r.db.Exec(ctx, query, postID); err != nil {
		return errors.DatabaseError("decrement likes count", err)
	}
	return nil
}

func (r *postRepository) IncrementCommentsCount(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET comments_count = comments_count + 1, updated_at = NOW() WHERE id = $1`
	if _, err := r.db.Exec(ctx, query, postID); err != nil {
		return errors.DatabaseError("increment comments count", err)
	}
	return nil
}

func (r *postRepository) DecrementCommentsCount(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET comments_count = GREATEST(0, comments_count - 1), updated_at = NOW() WHERE id = $1`
	if _, err := r.db.Exec(ctx, query, postID); err != nil {
		return errors.DatabaseError("decrement comments count", err)
	}
	return nil
}

// LikeRepository implementation
type likeRepository struct {
	db *database.DB
}

func NewLikeRepository(db *database.DB) social.LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(ctx context.Context, like *social.Like) error {
	query := `
		INSERT INTO likes (id, post_id, user_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (post_id, user_id) DO NOTHING
	`
	if _, err := r.db.Exec(ctx, query, like.ID, like.PostID, like.UserID, like.CreatedAt); err != nil {
		return errors.DatabaseError("create like", err)
	}
	return nil
}

func (r *likeRepository) Delete(ctx context.Context, postID, userID uuid.UUID) error {
	query := `DELETE FROM likes WHERE post_id = $1 AND user_id = $2`
	if _, err := r.db.Exec(ctx, query, postID, userID); err != nil {
		return errors.DatabaseError("delete like", err)
	}
	return nil
}

func (r *likeRepository) Exists(ctx context.Context, postID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = $1 AND user_id = $2)`
	var exists bool
	if err := r.db.QueryRow(ctx, query, postID, userID).Scan(&exists); err != nil {
		return false, errors.DatabaseError("check like exists", err)
	}
	return exists, nil
}

func (r *likeRepository) GetByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) ([]social.Like, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM likes WHERE post_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, postID).Scan(&total); err != nil {
		return nil, 0, errors.DatabaseError("count likes", err)
	}

	query := `
		SELECT id, post_id, user_id, created_at
		FROM likes
		WHERE post_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, postID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.DatabaseError("get likes by post", err)
	}
	defer rows.Close()

	likes := make([]social.Like, 0)
	for rows.Next() {
		var like social.Like
		if err := rows.Scan(&like.ID, &like.PostID, &like.UserID, &like.CreatedAt); err != nil {
			return nil, 0, errors.DatabaseError("scan like", err)
		}
		likes = append(likes, like)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, errors.DatabaseError("iterate likes", err)
	}

	return likes, total, nil
}

// CommentRepository implementation
type commentRepository struct {
	db *database.DB
}

func NewCommentRepository(db *database.DB) social.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *social.Comment) error {
	query := `
		INSERT INTO comments (id, post_id, user_id, parent_comment_id, content, reply_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	if _, err := r.db.Exec(ctx, query,
		comment.ID,
		comment.PostID,
		comment.UserID,
		comment.ParentCommentID,
		comment.Content,
		comment.ReplyCount,
		comment.CreatedAt,
		comment.UpdatedAt,
	); err != nil {
		return errors.DatabaseError("create comment", err)
	}
	return nil
}

func (r *commentRepository) GetByID(ctx context.Context, id uuid.UUID) (*social.Comment, error) {
	query := `
		SELECT
			c.id, c.post_id, c.user_id, c.parent_comment_id, c.content, c.reply_count,
			c.created_at, c.updated_at, c.deleted_at,
			u.id, u.name, u.avatar_url
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.id = $1
	`

	comment := &social.Comment{User: &social.PostUser{}}
	if err := r.db.QueryRow(ctx, query, id).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.ParentCommentID,
		&comment.Content,
		&comment.ReplyCount,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.DeletedAt,
		&comment.User.ID,
		&comment.User.Name,
		&comment.User.AvatarURL,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("comment")
		}
		return nil, errors.DatabaseError("get comment by id", err)
	}

	return comment, nil
}

func (r *commentRepository) GetByPostID(ctx context.Context, postID uuid.UUID, filter social.GetCommentsFilter) ([]social.Comment, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize

	countQuery := `SELECT COUNT(*) FROM comments WHERE post_id = $1`
	countArgs := []interface{}{postID}
	if filter.ParentCommentID == nil {
		countQuery = `SELECT COUNT(*) FROM comments WHERE post_id = $1 AND parent_comment_id IS NULL AND deleted_at IS NULL`
	} else {
		countQuery = `SELECT COUNT(*) FROM comments WHERE post_id = $1 AND parent_comment_id = $2 AND deleted_at IS NULL`
		countArgs = append(countArgs, *filter.ParentCommentID)
	}

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, errors.DatabaseError("count comments", err)
	}

	baseSelect := `
		SELECT
			c.id, c.post_id, c.user_id, c.parent_comment_id, c.content, c.reply_count,
			c.created_at, c.updated_at, c.deleted_at,
			u.id, u.name, u.avatar_url
		FROM comments c
		JOIN users u ON u.id = c.user_id
	`

	var query string
	var args []interface{}
	if filter.ParentCommentID == nil {
		query = baseSelect + `
			WHERE c.post_id = $1
			  AND c.parent_comment_id IS NULL
			  AND c.deleted_at IS NULL
			ORDER BY c.created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{postID, filter.PageSize, offset}
	} else {
		query = baseSelect + `
			WHERE c.post_id = $1
			  AND c.parent_comment_id = $2
			  AND c.deleted_at IS NULL
			ORDER BY c.created_at ASC
			LIMIT $3 OFFSET $4
		`
		args = []interface{}{postID, *filter.ParentCommentID, filter.PageSize, offset}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.DatabaseError("get comments by post", err)
	}
	defer rows.Close()

	comments := make([]social.Comment, 0)
	for rows.Next() {
		comment := social.Comment{User: &social.PostUser{}}
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.ParentCommentID,
			&comment.Content,
			&comment.ReplyCount,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.DeletedAt,
			&comment.User.ID,
			&comment.User.Name,
			&comment.User.AvatarURL,
		); err != nil {
			return nil, 0, errors.DatabaseError("scan comment", err)
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, errors.DatabaseError("iterate comments", err)
	}

	return comments, total, nil
}

func (r *commentRepository) GetLatestRepliesByParentIDs(ctx context.Context, parentCommentIDs []uuid.UUID, limitPerParent int) (map[uuid.UUID][]social.Comment, error) {
	result := make(map[uuid.UUID][]social.Comment)
	if len(parentCommentIDs) == 0 || limitPerParent <= 0 {
		return result, nil
	}

	query := `
		WITH ranked AS (
			SELECT
				c.id, c.post_id, c.user_id, c.parent_comment_id, c.content, c.reply_count,
				c.created_at, c.updated_at, c.deleted_at,
				u.id AS author_id, u.name, u.avatar_url,
				ROW_NUMBER() OVER (PARTITION BY c.parent_comment_id ORDER BY c.created_at DESC) AS rn
			FROM comments c
			JOIN users u ON u.id = c.user_id
			WHERE c.parent_comment_id = ANY($1)
			  AND c.deleted_at IS NULL
		)
		SELECT
			id, post_id, user_id, parent_comment_id, content, reply_count,
			created_at, updated_at, deleted_at,
			author_id, name, avatar_url
		FROM ranked
		WHERE rn <= $2
		ORDER BY parent_comment_id, created_at ASC
	`

	rows, err := r.db.Query(ctx, query, parentCommentIDs, limitPerParent)
	if err != nil {
		return nil, errors.DatabaseError("get latest replies", err)
	}
	defer rows.Close()

	for rows.Next() {
		comment := social.Comment{User: &social.PostUser{}}
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.ParentCommentID,
			&comment.Content,
			&comment.ReplyCount,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.DeletedAt,
			&comment.User.ID,
			&comment.User.Name,
			&comment.User.AvatarURL,
		); err != nil {
			return nil, errors.DatabaseError("scan latest reply", err)
		}
		if comment.ParentCommentID != nil {
			result[*comment.ParentCommentID] = append(result[*comment.ParentCommentID], comment)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, errors.DatabaseError("iterate latest replies", err)
	}

	return result, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *social.Comment) error {
	query := `
		UPDATE comments
		SET content = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(ctx, query, comment.ID, comment.Content)
	if err != nil {
		return errors.DatabaseError("update comment", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound("comment")
	}
	return nil
}

func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE comments
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return errors.DatabaseError("delete comment", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound("comment")
	}
	return nil
}

func (r *commentRepository) IncrementReplyCount(ctx context.Context, parentCommentID uuid.UUID) error {
	query := `UPDATE comments SET reply_count = reply_count + 1 WHERE id = $1`
	result, err := r.db.Exec(ctx, query, parentCommentID)
	if err != nil {
		return errors.DatabaseError("increment reply count", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound("comment")
	}
	return nil
}

func (r *commentRepository) DecrementReplyCount(ctx context.Context, parentCommentID uuid.UUID) error {
	query := `UPDATE comments SET reply_count = GREATEST(0, reply_count - 1) WHERE id = $1`
	result, err := r.db.Exec(ctx, query, parentCommentID)
	if err != nil {
		return errors.DatabaseError("decrement reply count", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound("comment")
	}
	return nil
}

// MediaAssetRepository implementation
type mediaAssetRepository struct {
	db *database.DB
}

func NewMediaAssetRepository(db *database.DB) social.MediaAssetRepository {
	return &mediaAssetRepository{db: db}
}

func (r *mediaAssetRepository) CreatePending(ctx context.Context, asset *social.SocialMediaAsset) error {
	query := `
		INSERT INTO social_media_assets (
			public_id, user_id, resource_type, status, expires_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (public_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			resource_type = EXCLUDED.resource_type,
			status = EXCLUDED.status,
			expires_at = EXCLUDED.expires_at,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Exec(ctx, query,
		asset.PublicID,
		asset.UserID,
		asset.ResourceType,
		asset.Status,
		asset.ExpiresAt,
		asset.CreatedAt,
		asset.UpdatedAt,
	)
	if err != nil {
		return errors.DatabaseError("create pending media asset", err)
	}

	return nil
}

func (r *mediaAssetRepository) Confirm(ctx context.Context, userID uuid.UUID, publicID string, secureURL *string, bytes *int64) error {
	query := `
		UPDATE social_media_assets
		SET secure_url = COALESCE($3, secure_url),
			bytes = COALESCE($4, bytes),
			status = 'ready',
			confirmed_at = COALESCE(confirmed_at, NOW()),
			updated_at = NOW()
		WHERE public_id = $1
		  AND user_id = $2
		  AND status IN ('uploading', 'ready')
	`

	result, err := r.db.Exec(ctx, query, publicID, userID, secureURL, bytes)
	if err != nil {
		return errors.DatabaseError("confirm media asset", err)
	}

	if result.RowsAffected() == 0 {
		return errors.NotFound("media asset")
	}

	return nil
}

type preferenceRepository struct {
	db *database.DB
}

func NewPreferenceRepository(db *database.DB) social.PreferenceRepository {
	return &preferenceRepository{db: db}
}

func (r *preferenceRepository) Upsert(ctx context.Context, preference *social.PostPreference) error {
	query := `
		INSERT INTO post_preferences (user_id, post_id, preference, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, post_id)
		DO UPDATE
		SET preference = EXCLUDED.preference,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(ctx, query,
		preference.UserID,
		preference.PostID,
		preference.Preference,
		preference.CreatedAt,
		preference.UpdatedAt,
	)
	if err != nil {
		return errors.DatabaseError("upsert post preference", err)
	}
	return nil
}

func (r *preferenceRepository) Delete(ctx context.Context, userID, postID uuid.UUID, preference string) error {
	query := `
		DELETE FROM post_preferences
		WHERE user_id = $1
		  AND post_id = $2
		  AND preference = $3
	`
	if _, err := r.db.Exec(ctx, query, userID, postID, preference); err != nil {
		return errors.DatabaseError("delete post preference", err)
	}
	return nil
}

func (r *preferenceRepository) GetByPostAndUser(ctx context.Context, userID, postID uuid.UUID) (*social.PostPreference, error) {
	query := `
		SELECT user_id, post_id, preference, created_at, updated_at
		FROM post_preferences
		WHERE user_id = $1
		  AND post_id = $2
	`
	var preference social.PostPreference
	err := r.db.QueryRow(ctx, query, userID, postID).Scan(
		&preference.UserID,
		&preference.PostID,
		&preference.Preference,
		&preference.CreatedAt,
		&preference.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, errors.DatabaseError("get post preference", err)
	}
	return &preference, nil
}

type reportRepository struct {
	db *database.DB
}

func NewReportRepository(db *database.DB) social.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Upsert(ctx context.Context, report *social.PostReport) error {
	query := `
		INSERT INTO post_reports (
			id, post_id, reporter_id, reason, description, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (post_id, reporter_id)
		DO UPDATE
		SET reason = EXCLUDED.reason,
			description = EXCLUDED.description,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(ctx, query,
		report.ID,
		report.PostID,
		report.ReporterID,
		report.Reason,
		report.Description,
		report.Status,
		report.CreatedAt,
		report.UpdatedAt,
	)
	if err != nil {
		return errors.DatabaseError("upsert post report", err)
	}
	return nil
}

type blockRepository struct {
	db *database.DB
}

func NewBlockRepository(db *database.DB) social.BlockRepository {
	return &blockRepository{db: db}
}

func (r *blockRepository) Block(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	query := `
		INSERT INTO user_blocks (blocker_id, blocked_id)
		VALUES ($1, $2)
		ON CONFLICT (blocker_id, blocked_id) DO NOTHING
	`
	if _, err := r.db.Exec(ctx, query, blockerID, blockedID); err != nil {
		return errors.DatabaseError("block user", err)
	}
	return nil
}

func (r *blockRepository) Unblock(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	query := `DELETE FROM user_blocks WHERE blocker_id = $1 AND blocked_id = $2`
	if _, err := r.db.Exec(ctx, query, blockerID, blockedID); err != nil {
		return errors.DatabaseError("unblock user", err)
	}
	return nil
}

func (r *blockRepository) IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_blocks WHERE blocker_id = $1 AND blocked_id = $2)`
	var blocked bool
	if err := r.db.QueryRow(ctx, query, blockerID, blockedID).Scan(&blocked); err != nil {
		return false, errors.DatabaseError("check blocked relation", err)
	}
	return blocked, nil
}
