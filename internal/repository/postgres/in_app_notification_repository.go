package postgres

import (
	"context"
	"strings"

	"gym-pro-2026-ptit/internal/domain/social"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/pkg/errors"

	"github.com/google/uuid"
)

type inAppNotificationRepository struct {
	db *database.DB
}

func NewInAppNotificationRepository(db *database.DB) social.InAppNotificationRepository {
	return &inAppNotificationRepository{db: db}
}

func (r *inAppNotificationRepository) Create(ctx context.Context, n *social.InAppNotification) error {
	if n == nil {
		return errors.BadRequest("notification is nil")
	}
	q := `
		INSERT INTO in_app_notifications (id, user_id, type, title, meta, post_id, related_post_id, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, q, n.ID, n.UserID, n.Type, n.Title, n.Meta, n.PostID, n.RelatedPostID, n.IsRead, n.CreatedAt)
	if err != nil {
		return errors.DatabaseError("create in_app_notification", err)
	}
	return nil
}

func (r *inAppNotificationRepository) ListForUser(ctx context.Context, userID uuid.UUID, filter string, page, limit int) ([]social.InAppNotification, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	filter = strings.TrimSpace(strings.ToLower(filter))
	if filter != "all" && filter != "workouts" && filter != "social" && filter != "nutrition" {
		filter = "all"
	}

	var countQ, listQ string
	var args []interface{}
	if filter != "all" {
		countQ = `SELECT COUNT(*) FROM in_app_notifications WHERE user_id = $1 AND type = $2`
		listQ = `
			SELECT id, user_id, type, title, meta, post_id, related_post_id, is_read, created_at
			FROM in_app_notifications
			WHERE user_id = $1 AND type = $2
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4`
		args = []interface{}{userID, filter, limit, offset}
	} else {
		countQ = `SELECT COUNT(*) FROM in_app_notifications WHERE user_id = $1`
		listQ = `
			SELECT id, user_id, type, title, meta, post_id, related_post_id, is_read, created_at
			FROM in_app_notifications
			WHERE user_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3`
		args = []interface{}{userID, limit, offset}
	}

	var countArgs []interface{}
	if filter != "all" {
		countArgs = []interface{}{userID, filter}
	} else {
		countArgs = []interface{}{userID}
	}

	var total int64
	if err := r.db.QueryRow(ctx, countQ, countArgs...).Scan(&total); err != nil {
		return nil, 0, errors.DatabaseError("count in_app_notifications", err)
	}

	rows, err := r.db.Query(ctx, listQ, args...)
	if err != nil {
		return nil, 0, errors.DatabaseError("list in_app_notifications", err)
	}
	defer rows.Close()

	out := make([]social.InAppNotification, 0)
	for rows.Next() {
		var n social.InAppNotification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Meta, &n.PostID, &n.RelatedPostID, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, 0, errors.DatabaseError("scan in_app_notification", err)
		}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, errors.DatabaseError("iterate in_app_notifications", err)
	}
	return out, total, nil
}

func (r *inAppNotificationRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	q := `SELECT COUNT(*) FROM in_app_notifications WHERE user_id = $1 AND is_read = false`
	var n int64
	if err := r.db.QueryRow(ctx, q, userID).Scan(&n); err != nil {
		return 0, errors.DatabaseError("count unread notifications", err)
	}
	return n, nil
}

func (r *inAppNotificationRepository) MarkRead(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	q := `
		UPDATE in_app_notifications
		SET is_read = true
		WHERE user_id = $1 AND id = ANY($2::uuid[]) AND is_read = false
	`
	ct, err := r.db.Exec(ctx, q, userID, ids)
	if err != nil {
		return 0, errors.DatabaseError("mark notifications read", err)
	}
	return ct.RowsAffected(), nil
}
