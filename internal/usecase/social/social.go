package social

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"gym-pro-2026-ptit/internal/config"
	mealdomain "gym-pro-2026-ptit/internal/domain/meal"
	socialdomain "gym-pro-2026-ptit/internal/domain/social"
	"gym-pro-2026-ptit/internal/domain/user"
	workoutdomain "gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/internal/port/socialnotify"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/utils"
	"gym-pro-2026-ptit/pkg/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const maxCommentPreviewRepliesPerParent = 100

const socialCommentEditWindow = 15 * time.Minute

func roundMealLogForShare(log *mealdomain.MealLog) {
	if log == nil {
		return
	}
	log.TotalCalories = utils.RoundToTwo(log.TotalCalories)
	log.TotalProteinG = utils.RoundToTwo(log.TotalProteinG)
	log.TotalCarbsG = utils.RoundToTwo(log.TotalCarbsG)
	log.TotalFatG = utils.RoundToTwo(log.TotalFatG)
	for i := range log.Items {
		item := &log.Items[i]
		item.Calories = utils.RoundToTwo(item.Calories)
		item.ProteinG = utils.RoundToTwo(item.ProteinG)
		item.CarbsG = utils.RoundToTwo(item.CarbsG)
		item.FatG = utils.RoundToTwo(item.FatG)
		if item.Food != nil {
			item.Food.Calories = utils.RoundToTwo(item.Food.Calories)
			item.Food.ProteinG = utils.RoundToTwo(item.Food.ProteinG)
			item.Food.CarbsG = utils.RoundToTwo(item.Food.CarbsG)
			item.Food.FatG = utils.RoundToTwo(item.Food.FatG)
		}
		if item.Recipe != nil {
			item.Recipe.TotalCalories = utils.RoundToTwo(item.Recipe.TotalCalories)
			item.Recipe.TotalProteinG = utils.RoundToTwo(item.Recipe.TotalProteinG)
			item.Recipe.TotalCarbsG = utils.RoundToTwo(item.Recipe.TotalCarbsG)
			item.Recipe.TotalFatG = utils.RoundToTwo(item.Recipe.TotalFatG)
			item.Recipe.PerServingCalories = utils.RoundToTwo(item.Recipe.PerServingCalories)
			item.Recipe.PerServingProteinG = utils.RoundToTwo(item.Recipe.PerServingProteinG)
			item.Recipe.PerServingCarbsG = utils.RoundToTwo(item.Recipe.PerServingCarbsG)
			item.Recipe.PerServingFatG = utils.RoundToTwo(item.Recipe.PerServingFatG)
		}
	}
}

type SocialUseCases struct {
	postRepo           socialdomain.PostRepository
	followRepo         socialdomain.FollowRepository
	likeRepo           socialdomain.LikeRepository
	commentRepo        socialdomain.CommentRepository
	mediaAssetRepo     socialdomain.MediaAssetRepository
	preferenceRepo     socialdomain.PreferenceRepository
	reportRepo         socialdomain.ReportRepository
	blockRepo          socialdomain.BlockRepository
	notifRepo          socialdomain.InAppNotificationRepository
	userRepo           user.Repository
	mealLogRepo        mealdomain.MealLogRepository
	workoutSessionRepo workoutdomain.WorkoutSessionRepository
	workoutPlanRepo    workoutdomain.WorkoutPlanRepository
	validator          *validator.Validator
	cloudinary         config.CloudinaryConfig
	notify             socialnotify.Broadcaster
}

func NewSocialUseCases(
	cfg *config.Config,
	postRepo socialdomain.PostRepository,
	followRepo socialdomain.FollowRepository,
	likeRepo socialdomain.LikeRepository,
	commentRepo socialdomain.CommentRepository,
	mediaAssetRepo socialdomain.MediaAssetRepository,
	preferenceRepo socialdomain.PreferenceRepository,
	reportRepo socialdomain.ReportRepository,
	blockRepo socialdomain.BlockRepository,
	notifRepo socialdomain.InAppNotificationRepository,
	userRepo user.Repository,
	mealLogRepo mealdomain.MealLogRepository,
	workoutSessionRepo workoutdomain.WorkoutSessionRepository,
	workoutPlanRepo workoutdomain.WorkoutPlanRepository,
	validator *validator.Validator,
	notify socialnotify.Broadcaster,
) *SocialUseCases {
	return &SocialUseCases{
		postRepo:           postRepo,
		followRepo:         followRepo,
		likeRepo:           likeRepo,
		commentRepo:        commentRepo,
		mediaAssetRepo:     mediaAssetRepo,
		preferenceRepo:     preferenceRepo,
		reportRepo:         reportRepo,
		blockRepo:          blockRepo,
		notifRepo:          notifRepo,
		userRepo:           userRepo,
		mealLogRepo:        mealLogRepo,
		workoutSessionRepo: workoutSessionRepo,
		workoutPlanRepo:    workoutPlanRepo,
		validator:          validator,
		cloudinary:         cfg.Cloudinary,
		notify:             notify,
	}
}

type LikeResponse struct {
	LikeCount   int  `json:"like_count"`
	IsLikedByMe bool `json:"is_liked_by_me"`
}

type CreateCommentInput struct {
	Content  *string                             `json:"content,omitempty" validate:"omitempty,min=1,max=1000"`
	Media    []socialdomain.CreatePostMediaInput `json:"media,omitempty" validate:"omitempty,dive"`
	ParentID *string                             `json:"parentId" validate:"omitempty,uuid4"`
}

type UpdateCommentInput struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

type DeleteCommentResult struct {
	DeletedByRole string `json:"deletedByRole"`
}

type CommentAuthorOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatarUrl"`
}

type CommentOutput struct {
	ID               uuid.UUID           `json:"id"`
	PostID           uuid.UUID           `json:"postId"`
	ParentID         *uuid.UUID          `json:"parentId"`
	Depth            int                 `json:"depth"`
	Path             string              `json:"path"`
	DirectReplyCount int                 `json:"directReplyCount"`
	PreviewReplies   []CommentOutput     `json:"previewReplies"`
	Author           CommentAuthorOutput `json:"author"`
	Content          string              `json:"content"`
	Media            []PostMediaOutput   `json:"media"`
	IsDeleted        bool                `json:"isDeleted"`
	IsEdited         bool                `json:"isEdited"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
}

type CommentListOutput struct {
	Comments   []CommentOutput `json:"comments"`
	NextCursor *string         `json:"nextCursor"`
}

type CommentRepliesOutput struct {
	Replies    []CommentOutput `json:"replies"`
	NextCursor *string         `json:"nextCursor"`
}

type CreatePostInput struct {
	Caption        *string                             `json:"caption,omitempty" validate:"omitempty,max=2000"`
	Media          []socialdomain.CreatePostMediaInput `json:"media,omitempty" validate:"omitempty,dive"`
	ContentType    *string                             `json:"content_type,omitempty"`
	ContentTypeAlt *string                             `json:"contentType,omitempty"`
	ContentID      *string                             `json:"content_id,omitempty"`
	ContentIDAlt   *string                             `json:"contentId,omitempty"`
	Feeling        *string                             `json:"feeling,omitempty" validate:"omitempty,max=100"`
	Location       *CreatePostLocationInput            `json:"location,omitempty" validate:"omitempty"`
	Hashtags       []string                            `json:"hashtags,omitempty" validate:"omitempty,dive,max=50"`
}

type UpdatePostInput struct {
	Caption       *string                              `json:"caption,omitempty" validate:"omitempty,max=2000"`
	Feeling       *string                              `json:"feeling,omitempty" validate:"omitempty,max=100"`
	Location      *CreatePostLocationInput             `json:"location,omitempty" validate:"omitempty"`
	Hashtags      *[]string                            `json:"hashtags,omitempty" validate:"omitempty,dive,max=50"`
	Media         *[]socialdomain.CreatePostMediaInput `json:"media,omitempty" validate:"omitempty,dive"`
	ClearLocation *bool                                `json:"clear_location,omitempty"`
}

type CreatePostLocationInput struct {
	Name string `json:"name" validate:"required,max=255"`
}

type CreateMediaSignatureInput struct {
	ResourceType string `json:"resource_type" validate:"required,oneof=image video"`
	Folder       string `json:"folder" validate:"required"`
}

type MediaSignatureOutput struct {
	CloudName    string `json:"cloud_name"`
	APIKey       string `json:"api_key"`
	Timestamp    int64  `json:"timestamp"`
	Folder       string `json:"folder"`
	PublicID     string `json:"public_id"`
	Signature    string `json:"signature"`
	UploadURL    string `json:"upload_url"`
	ExpiresIn    int    `json:"expires_in"`
	ResourceType string `json:"resource_type"`
}

type ConfirmMediaInput struct {
	PublicID     string `json:"public_id" validate:"required"`
	SecureURL    string `json:"secure_url" validate:"omitempty,url"`
	ResourceType string `json:"resource_type" validate:"required,oneof=image video"`
	Bytes        int64  `json:"bytes" validate:"gte=0"`
}

type ConfirmMediaOutput struct {
	PublicID   string `json:"public_id"`
	IsOwned    bool   `json:"is_owned"`
	IsValid    bool   `json:"is_valid"`
	AssetState string `json:"asset_state"`
}

type FeedPagination struct {
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

type FeedOutput struct {
	Data       []PostOutput   `json:"data"`
	Pagination FeedPagination `json:"pagination"`
}

type NotificationRowOutput struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Meta      string    `json:"meta"`
	DayGroup  string    `json:"dayGroup"`
	IsRead    bool      `json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
	PostID    *string   `json:"postId,omitempty"`
	Kind      string    `json:"kind,omitempty"`
}

type notificationsListPagination struct {
	NextCursor *string `json:"nextCursor"`
	HasMore    bool    `json:"hasMore"`
}

type NotificationsListPayload struct {
	Data       []NotificationRowOutput     `json:"data"`
	Pagination notificationsListPagination `json:"pagination"`
}

type UnreadNotificationsCountPayload struct {
	Unread int64 `json:"unread"`
}

type MarkNotificationsReadPayload struct {
	Updated int64 `json:"updated"`
}

type UserSearchHit struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Subtitle  string    `json:"subtitle,omitempty"`
}

type SocialSearchOutput struct {
	Posts      []PostOutput    `json:"posts,omitempty"`
	Users      []UserSearchHit `json:"users,omitempty"`
	Pagination FeedPagination  `json:"pagination"`
}

type AuthorOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
}

type PostMediaOutput struct {
	Type     string `json:"type"`
	URL      string `json:"url"`
	PublicID string `json:"public_id,omitempty"`
}

type MealSummaryOutput struct {
	MealLogID     uuid.UUID `json:"meal_log_id"`
	LogDate       string    `json:"log_date"`
	MealTime      string    `json:"meal_time"`
	TotalCalories float64   `json:"total_calories"`
	TotalProteinG float64   `json:"total_protein_g"`
	TotalCarbsG   float64   `json:"total_carbs_g"`
	TotalFatG     float64   `json:"total_fat_g"`
	ItemCount     int       `json:"item_count"`
}

type WorkoutSessionSummaryOutput struct {
	SessionID           uuid.UUID  `json:"session_id"`
	PlanTitle           string     `json:"plan_title"`
	ScheduledDate       *string    `json:"scheduled_date,omitempty"`
	Status              string     `json:"status"`
	DurationMins        *int       `json:"duration_mins,omitempty"`
	TotalCaloriesBurned *int       `json:"total_calories_burned,omitempty"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
}

type PostOutput struct {
	ID                    uuid.UUID                    `json:"id"`
	Author                AuthorOutput                 `json:"author"`
	ContentType           string                       `json:"content_type"`
	ContentID             *uuid.UUID                   `json:"content_id"`
	TimeLabel             string                       `json:"time_label"`
	Caption               string                       `json:"caption"`
	Media                 []PostMediaOutput            `json:"media"`
	Feeling               *string                      `json:"feeling"`
	Location              *PostLocation                `json:"location"`
	Hashtags              []string                     `json:"hashtags"`
	LikeCount             int                          `json:"like_count"`
	CommentCount          int                          `json:"comment_count"`
	IsLikedByMe           bool                         `json:"is_liked_by_me"`
	IsInterestedByMe      bool                         `json:"is_interested_by_me"`
	IsNotInterestedByMe   bool                         `json:"is_not_interested_by_me"`
	IsEdited              bool                         `json:"is_edited"`
	SharedExercises       []interface{}                `json:"shared_exercises,omitempty"`
	MealSummary           *MealSummaryOutput           `json:"meal_summary,omitempty"`
	WorkoutSessionSummary *WorkoutSessionSummaryOutput `json:"workout_session_summary,omitempty"`
	CreatedAt             time.Time                    `json:"created_at"`
	UpdatedAt             time.Time                    `json:"updated_at"`
}

type PostLocation struct {
	Name string `json:"name"`
}

var hashtagPattern = regexp.MustCompile(`(?:^|[^[:alnum:]_])#([[:alnum:]_]+)`)

type UserProfileOutput struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	AvatarURL      *string   `json:"avatar_url,omitempty"`
	Subtitle       string    `json:"subtitle"`
	PostsCount     int64     `json:"posts_count"`
	FollowersCount int       `json:"followers_count"`
	FollowingCount int       `json:"following_count"`
	IsFollowing    bool      `json:"is_following"`
	IsMe           bool      `json:"is_me"`
}

type FollowActionOutput struct {
	UserID         uuid.UUID `json:"userId"`
	IsFollowing    bool      `json:"isFollowing"`
	FollowersCount int       `json:"followersCount"`
}

type PreferenceActionOutput struct {
	PostID              uuid.UUID `json:"post_id"`
	IsInterestedByMe    bool      `json:"is_interested_by_me"`
	IsNotInterestedByMe bool      `json:"is_not_interested_by_me"`
}

type ReportPostInput struct {
	Reason      string  `json:"reason" validate:"required,oneof=spam harassment misinformation nudity violence other"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

type ReportPostOutput struct {
	PostID      uuid.UUID `json:"post_id"`
	Reason      string    `json:"reason"`
	Description *string   `json:"description,omitempty"`
	Status      string    `json:"status"`
}

type BlockActionOutput struct {
	UserID    uuid.UUID `json:"user_id"`
	IsBlocked bool      `json:"is_blocked"`
}

func (uc *SocialUseCases) GetFeed(ctx context.Context, userID uuid.UUID, cursor string, limit int) (*FeedOutput, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	page, err := decodePageCursor(cursor)
	if err != nil {
		return nil, errors.BadRequest("invalid cursor")
	}

	items, total, err := uc.postRepo.GetFeed(ctx, userID, socialdomain.GetFeedFilter{Page: page, PageSize: limit})
	if err != nil {
		return nil, err
	}

	postIDs := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		if item.Post != nil {
			postIDs = append(postIDs, item.Post.ID)
		}
	}

	mediaByPostID, err := uc.postRepo.GetMediaByPostIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	data := make([]PostOutput, 0, len(items))
	for _, item := range items {
		if item.Post != nil {
			item.Post.Media = mediaByPostID[item.Post.ID]
		}
		data = append(data, uc.postOutputFrom(ctx, item.Post, item.IsLiked, item.IsInterested, item.IsNotInterested))
	}

	hasMore := int64(page*limit) < total
	output := &FeedOutput{
		Data:       data,
		Pagination: FeedPagination{HasMore: hasMore},
	}
	if hasMore {
		output.Pagination.NextCursor = encodePageCursor(page + 1)
	}

	return output, nil
}

func sanitizeSocialSearchQuery(q string) string {
	q = strings.TrimSpace(q)
	q = strings.ReplaceAll(q, "%", "")
	q = strings.ReplaceAll(q, "_", "")
	if len(q) > 200 {
		q = q[:200]
	}
	return q
}

func userSearchSubtitle(bio *string) string {
	if bio == nil {
		return ""
	}
	s := strings.TrimSpace(*bio)
	if s == "" {
		return ""
	}
	if len(s) > 120 {
		return s[:120] + "…"
	}
	return s
}

func (uc *SocialUseCases) Search(ctx context.Context, viewerID uuid.UUID, rawQuery, typ, cursor string, limit int) (*SocialSearchOutput, error) {
	q := sanitizeSocialSearchQuery(rawQuery)
	if len(q) < 1 {
		return nil, errors.BadRequest("query is required")
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}
	page, err := decodePageCursor(cursor)
	if err != nil {
		return nil, errors.BadRequest("invalid cursor")
	}
	typ = strings.ToLower(strings.TrimSpace(typ))
	switch typ {
	case "posts":
		items, total, err := uc.postRepo.SearchPosts(ctx, viewerID, q, socialdomain.GetFeedFilter{Page: page, PageSize: limit})
		if err != nil {
			return nil, err
		}
		postIDs := make([]uuid.UUID, 0, len(items))
		for _, item := range items {
			if item.Post != nil {
				postIDs = append(postIDs, item.Post.ID)
			}
		}
		mediaByPostID, err := uc.postRepo.GetMediaByPostIDs(ctx, postIDs)
		if err != nil {
			return nil, err
		}
		data := make([]PostOutput, 0, len(items))
		for _, item := range items {
			if item.Post != nil {
				item.Post.Media = mediaByPostID[item.Post.ID]
			}
			data = append(data, uc.postOutputFrom(ctx, item.Post, item.IsLiked, item.IsInterested, item.IsNotInterested))
		}
		hasMore := int64(page*limit) < total
		out := &SocialSearchOutput{Posts: data, Pagination: FeedPagination{HasMore: hasMore}}
		if hasMore {
			out.Pagination.NextCursor = encodePageCursor(page + 1)
		}
		return out, nil
	case "users":
		rows, total, err := uc.followRepo.SearchUsers(ctx, viewerID, q, page, limit)
		if err != nil {
			return nil, err
		}
		users := make([]UserSearchHit, 0, len(rows))
		for _, row := range rows {
			users = append(users, UserSearchHit{
				ID:        row.ID,
				Name:      row.Name,
				AvatarURL: row.AvatarURL,
				Subtitle:  userSearchSubtitle(row.Bio),
			})
		}
		hasMore := int64(page*limit) < total
		out := &SocialSearchOutput{Users: users, Pagination: FeedPagination{HasMore: hasMore}}
		if hasMore {
			out.Pagination.NextCursor = encodePageCursor(page + 1)
		}
		return out, nil
	default:
		return nil, errors.BadRequest("type must be posts or users")
	}
}

func (uc *SocialUseCases) GetPostByID(ctx context.Context, viewerID uuid.UUID, postID string) (*PostOutput, error) {
	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}

	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	isLiked, err := uc.likeRepo.Exists(ctx, post.ID, viewerID)
	if err != nil {
		return nil, err
	}

	isInterested := false
	isNotInterested := false
	pref, err := uc.preferenceRepo.GetByPostAndUser(ctx, viewerID, post.ID)
	if err != nil {
		return nil, err
	}
	if pref != nil {
		switch pref.Preference {
		case "interested":
			isInterested = true
		case "not_interested":
			isNotInterested = true
		}
	}

	mediaByPostID, err := uc.postRepo.GetMediaByPostIDs(ctx, []uuid.UUID{post.ID})
	if err != nil {
		return nil, err
	}
	post.Media = mediaByPostID[post.ID]

	out := uc.postOutputFrom(ctx, post, isLiked, isInterested, isNotInterested)
	return &out, nil
}

func (uc *SocialUseCases) loadPostByIDString(ctx context.Context, postID string) (*socialdomain.Post, error) {
	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}
	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (uc *SocialUseCases) GetPostAttachedMealLog(ctx context.Context, viewerID uuid.UUID, postID string) (*mealdomain.MealLog, error) {
	post, err := uc.loadPostByIDString(ctx, postID)
	if err != nil {
		return nil, err
	}
	ct := strings.TrimSpace(strings.ToLower(post.ContentType))
	if ct != "meal_log" || post.ContentID == nil {
		return nil, errors.NotFound("attached meal")
	}
	log, err := uc.mealLogRepo.GetByID(ctx, *post.ContentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("meal log")
		}
		return nil, errors.DatabaseError("get meal log", err)
	}
	if log.UserID != post.UserID {
		return nil, errors.NotFound("attached meal")
	}
	roundMealLogForShare(log)
	return log, nil
}

func (uc *SocialUseCases) GetPostAttachedWorkoutSession(ctx context.Context, viewerID uuid.UUID, postID string) (*workoutdomain.WorkoutSession, error) {
	post, err := uc.loadPostByIDString(ctx, postID)
	if err != nil {
		return nil, err
	}
	ct := strings.TrimSpace(strings.ToLower(post.ContentType))
	if ct != "workout_session" || post.ContentID == nil {
		return nil, errors.NotFound("attached workout")
	}
	s, err := uc.workoutSessionRepo.GetByID(ctx, *post.ContentID)
	if err != nil {
		return nil, err
	}
	if s.UserID != post.UserID {
		return nil, errors.NotFound("attached workout")
	}
	return s, nil
}

func (uc *SocialUseCases) GetUserProfile(ctx context.Context, currentUserID uuid.UUID, profileUserID string) (*UserProfileOutput, error) {
	targetID, err := uuid.Parse(profileUserID)
	if err != nil {
		return nil, errors.BadRequest("invalid user id")
	}

	u, err := uc.userRepo.GetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	stats, err := uc.followRepo.GetStats(ctx, targetID)
	if err != nil {
		return nil, err
	}

	_, totalPosts, err := uc.postRepo.GetByUserID(ctx, targetID, 1, 1)
	if err != nil {
		return nil, err
	}

	isFollowing := false
	if currentUserID != targetID {
		isFollowing, err = uc.followRepo.IsFollowing(ctx, currentUserID, targetID)
		if err != nil {
			return nil, err
		}
	}

	subtitle := fmt.Sprintf("Member since %d", u.CreatedAt.Year())
	if u.Bio != nil && strings.TrimSpace(*u.Bio) != "" {
		subtitle = *u.Bio
	}

	return &UserProfileOutput{
		ID:             u.ID,
		Name:           u.Name,
		AvatarURL:      u.AvatarURL,
		Subtitle:       subtitle,
		PostsCount:     totalPosts,
		FollowersCount: stats.FollowersCount,
		FollowingCount: stats.FollowingCount,
		IsFollowing:    isFollowing,
		IsMe:           currentUserID == targetID,
	}, nil
}

func (uc *SocialUseCases) FollowUser(ctx context.Context, currentUserID uuid.UUID, targetUserID string) (*FollowActionOutput, error) {
	targetID, err := uuid.Parse(targetUserID)
	if err != nil {
		return nil, errors.BadRequest("invalid user id")
	}

	if currentUserID == targetID {
		return nil, errors.BadRequest("cannot follow yourself")
	}

	hasBlockRelation, err := uc.followRepo.HasBlockRelation(ctx, currentUserID, targetID)
	if err != nil {
		return nil, err
	}
	if hasBlockRelation {
		return nil, errors.Forbidden("cannot follow user due to block relationship")
	}

	if _, err := uc.userRepo.GetByID(ctx, targetID); err != nil {
		return nil, err
	}

	wasFollowing, err := uc.followRepo.IsFollowing(ctx, currentUserID, targetID)
	if err != nil {
		return nil, err
	}
	if err := uc.followRepo.Follow(ctx, currentUserID, targetID); err != nil {
		return nil, err
	}
	if !wasFollowing {
		uc.tryCreateFollowNotification(ctx, currentUserID, targetID)
	}

	stats, err := uc.followRepo.GetStats(ctx, targetID)
	if err != nil {
		return nil, err
	}

	return &FollowActionOutput{
		UserID:         targetID,
		IsFollowing:    true,
		FollowersCount: stats.FollowersCount,
	}, nil
}

func (uc *SocialUseCases) UnfollowUser(ctx context.Context, currentUserID uuid.UUID, targetUserID string) (*FollowActionOutput, error) {
	targetID, err := uuid.Parse(targetUserID)
	if err != nil {
		return nil, errors.BadRequest("invalid user id")
	}

	if currentUserID == targetID {
		return nil, errors.BadRequest("cannot unfollow yourself")
	}

	if _, err := uc.userRepo.GetByID(ctx, targetID); err != nil {
		return nil, err
	}

	if err := uc.followRepo.Unfollow(ctx, currentUserID, targetID); err != nil {
		return nil, err
	}

	stats, err := uc.followRepo.GetStats(ctx, targetID)
	if err != nil {
		return nil, err
	}

	return &FollowActionOutput{
		UserID:         targetID,
		IsFollowing:    false,
		FollowersCount: stats.FollowersCount,
	}, nil
}

func (uc *SocialUseCases) MarkInterested(ctx context.Context, userID uuid.UUID, postID string) (*PreferenceActionOutput, error) {
	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}
	if _, err := uc.postRepo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	now := time.Now()
	if err := uc.preferenceRepo.Upsert(ctx, &socialdomain.PostPreference{
		UserID:     userID,
		PostID:     id,
		Preference: "interested",
		CreatedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return nil, err
	}
	return &PreferenceActionOutput{
		PostID:              id,
		IsInterestedByMe:    true,
		IsNotInterestedByMe: false,
	}, nil
}

func (uc *SocialUseCases) UnmarkInterested(ctx context.Context, userID uuid.UUID, postID string) (*PreferenceActionOutput, error) {
	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}
	if _, err := uc.postRepo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	if err := uc.preferenceRepo.Delete(ctx, userID, id, "interested"); err != nil {
		return nil, err
	}
	currentPreference, err := uc.preferenceRepo.GetByPostAndUser(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	isNotInterested := currentPreference != nil && currentPreference.Preference == "not_interested"
	return &PreferenceActionOutput{
		PostID:              id,
		IsInterestedByMe:    false,
		IsNotInterestedByMe: isNotInterested,
	}, nil
}

func (uc *SocialUseCases) MarkNotInterested(ctx context.Context, userID uuid.UUID, postID string) (*PreferenceActionOutput, error) {
	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}
	if _, err := uc.postRepo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	now := time.Now()
	if err := uc.preferenceRepo.Upsert(ctx, &socialdomain.PostPreference{
		UserID:     userID,
		PostID:     id,
		Preference: "not_interested",
		CreatedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return nil, err
	}
	return &PreferenceActionOutput{
		PostID:              id,
		IsInterestedByMe:    false,
		IsNotInterestedByMe: true,
	}, nil
}

func (uc *SocialUseCases) UnmarkNotInterested(ctx context.Context, userID uuid.UUID, postID string) (*PreferenceActionOutput, error) {
	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}
	if _, err := uc.postRepo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	if err := uc.preferenceRepo.Delete(ctx, userID, id, "not_interested"); err != nil {
		return nil, err
	}
	currentPreference, err := uc.preferenceRepo.GetByPostAndUser(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	isInterested := currentPreference != nil && currentPreference.Preference == "interested"
	return &PreferenceActionOutput{
		PostID:              id,
		IsInterestedByMe:    isInterested,
		IsNotInterestedByMe: false,
	}, nil
}

func (uc *SocialUseCases) ReportPost(ctx context.Context, userID uuid.UUID, postID string, input ReportPostInput) (*ReportPostOutput, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}
	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}
	if _, err := uc.postRepo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	description := stringPointerOrNil("")
	if input.Description != nil {
		description = stringPointerOrNil(*input.Description)
	}
	now := time.Now()
	report := &socialdomain.PostReport{
		ID:          uuid.New(),
		PostID:      id,
		ReporterID:  userID,
		Reason:      input.Reason,
		Description: description,
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.reportRepo.Upsert(ctx, report); err != nil {
		return nil, err
	}
	return &ReportPostOutput{
		PostID:      report.PostID,
		Reason:      report.Reason,
		Description: report.Description,
		Status:      report.Status,
	}, nil
}

func (uc *SocialUseCases) BlockUser(ctx context.Context, currentUserID uuid.UUID, targetUserID string) (*BlockActionOutput, error) {
	targetID, err := uuid.Parse(targetUserID)
	if err != nil {
		return nil, errors.BadRequest("invalid user id")
	}
	if currentUserID == targetID {
		return nil, errors.BadRequest("cannot block yourself")
	}
	if _, err := uc.userRepo.GetByID(ctx, targetID); err != nil {
		return nil, err
	}
	if err := uc.blockRepo.Block(ctx, currentUserID, targetID); err != nil {
		return nil, err
	}
	return &BlockActionOutput{UserID: targetID, IsBlocked: true}, nil
}

func (uc *SocialUseCases) UnblockUser(ctx context.Context, currentUserID uuid.UUID, targetUserID string) (*BlockActionOutput, error) {
	targetID, err := uuid.Parse(targetUserID)
	if err != nil {
		return nil, errors.BadRequest("invalid user id")
	}
	if currentUserID == targetID {
		return nil, errors.BadRequest("cannot unblock yourself")
	}
	if err := uc.blockRepo.Unblock(ctx, currentUserID, targetID); err != nil {
		return nil, err
	}
	return &BlockActionOutput{UserID: targetID, IsBlocked: false}, nil
}

func (uc *SocialUseCases) GetUserPosts(ctx context.Context, viewerID uuid.UUID, profileUserID, cursor string, limit int) (*FeedOutput, error) {
	targetID, err := uuid.Parse(profileUserID)
	if err != nil {
		return nil, errors.BadRequest("invalid user id")
	}

	if _, err := uc.userRepo.GetByID(ctx, targetID); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	page, err := decodePageCursor(cursor)
	if err != nil {
		return nil, errors.BadRequest("invalid cursor")
	}

	posts, total, err := uc.postRepo.GetByUserID(ctx, targetID, page, limit)
	if err != nil {
		return nil, err
	}

	postIDs := make([]uuid.UUID, 0, len(posts))
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}

	mediaByPostID, err := uc.postRepo.GetMediaByPostIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	likesByPost, err := uc.likeRepo.ExistsForPosts(ctx, viewerID, postIDs)
	if err != nil {
		return nil, err
	}

	prefsByPost, err := uc.preferenceRepo.GetByPostsAndUser(ctx, viewerID, postIDs)
	if err != nil {
		return nil, err
	}

	data := make([]PostOutput, 0, len(posts))
	for _, post := range posts {
		post.Media = mediaByPostID[post.ID]
		isLiked := likesByPost[post.ID]
		isInterested := false
		isNotInterested := false
		if pref := prefsByPost[post.ID]; pref != nil {
			switch pref.Preference {
			case "interested":
				isInterested = true
			case "not_interested":
				isNotInterested = true
			}
		}
		data = append(data, uc.postOutputFrom(ctx, &post, isLiked, isInterested, isNotInterested))
	}

	hasMore := int64(page*limit) < total
	output := &FeedOutput{
		Data:       data,
		Pagination: FeedPagination{HasMore: hasMore},
	}
	if hasMore {
		output.Pagination.NextCursor = encodePageCursor(page + 1)
	}

	return output, nil
}

func mapPost(post *socialdomain.Post, isLiked, isInterested, isNotInterested bool) PostOutput {
	caption := ""
	if post.Caption != nil {
		caption = *post.Caption
	}

	author := AuthorOutput{}
	if post.User != nil {
		author = AuthorOutput{
			ID:        post.User.ID,
			Name:      post.User.Name,
			AvatarURL: post.User.AvatarURL,
		}
	}

	media := make([]PostMediaOutput, 0, len(post.Media))
	for _, item := range post.Media {
		urlValue := ""
		if item.SecureURL != nil {
			urlValue = *item.SecureURL
		}
		media = append(media, PostMediaOutput{Type: item.ResourceType, URL: urlValue, PublicID: item.PublicID})
	}

	hashtags := post.Hashtags
	if hashtags == nil {
		hashtags = []string{}
	}

	var location *PostLocation
	if post.LocationName != nil && strings.TrimSpace(*post.LocationName) != "" {
		location = &PostLocation{Name: *post.LocationName}
	}

	return PostOutput{
		ID:                  post.ID,
		Author:              author,
		ContentType:         post.ContentType,
		ContentID:           post.ContentID,
		TimeLabel:           humanizeTime(post.CreatedAt),
		Caption:             caption,
		Media:               media,
		Feeling:             post.Feeling,
		Location:            location,
		Hashtags:            hashtags,
		LikeCount:           post.LikesCount,
		CommentCount:        post.CommentsCount,
		IsLikedByMe:         isLiked,
		IsInterestedByMe:    isInterested,
		IsNotInterestedByMe: isNotInterested,
		IsEdited:            post.UpdatedAt.After(post.CreatedAt),
		SharedExercises:     []interface{}{},
		CreatedAt:           post.CreatedAt,
		UpdatedAt:           post.UpdatedAt,
	}
}

func humanizeTime(t time.Time) string {
	d := time.Since(t)
	if d < time.Minute {
		return "JUST NOW"
	}
	if d < time.Hour {
		return fmt.Sprintf("%d MIN AGO", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d H AGO", int(d.Hours()))
	}
	return fmt.Sprintf("%d D AGO", int(d.Hours()/24))
}

func encodePageCursor(page int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(page)))
}

func decodePageCursor(cursor string) (int, error) {
	if strings.TrimSpace(cursor) == "" {
		return 1, nil
	}

	raw, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}

	page, err := strconv.Atoi(string(raw))
	if err != nil || page < 1 {
		return 0, fmt.Errorf("invalid page cursor")
	}

	return page, nil
}

func isSocialMediaFolderOwnedByUser(folder string, postsPrefix string, commentsPrefix string) bool {
	for _, p := range []string{postsPrefix, commentsPrefix} {
		if folder == p || strings.HasPrefix(folder, p+"/") {
			return true
		}
	}
	return false
}

func isSocialMediaPublicIDOwnedByUser(publicID string, userID uuid.UUID) bool {
	uid := userID.String()
	for _, p := range []string{"posts/" + uid + "/", "comments/" + uid + "/"} {
		if strings.HasPrefix(publicID, p) {
			return true
		}
	}
	return false
}

func (uc *SocialUseCases) CreateMediaSignature(ctx context.Context, userID uuid.UUID, input CreateMediaSignatureInput) (*MediaSignatureOutput, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	if strings.TrimSpace(uc.cloudinary.URL) == "" {
		return nil, errors.InternalServer("cloudinary is not configured", nil)
	}

	postsPrefix := "posts/" + userID.String()
	commentsPrefix := "comments/" + userID.String()
	folder := strings.TrimSpace(input.Folder)
	if folder == "" {
		folder = postsPrefix
	}
	if !isSocialMediaFolderOwnedByUser(folder, postsPrefix, commentsPrefix) {
		return nil, errors.Forbidden("invalid folder for current user")
	}

	cloudName, apiKey, apiSecret, err := parseCloudinaryURL(uc.cloudinary.URL)
	if err != nil {
		return nil, errors.InternalServer("invalid cloudinary configuration", err)
	}

	timestamp := time.Now().Unix()
	publicID := fmt.Sprintf("%s/%d_%s", folder, timestamp, randomHex(4))
	params := map[string]string{
		"folder":    folder,
		"public_id": publicID,
		"timestamp": strconv.FormatInt(timestamp, 10),
	}
	signature := cloudinarySignature(params, apiSecret)

	expiresAt := time.Now().Add(24 * time.Hour)
	asset := &socialdomain.SocialMediaAsset{
		PublicID:     publicID,
		UserID:       userID,
		ResourceType: input.ResourceType,
		Status:       "uploading",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ExpiresAt:    &expiresAt,
	}
	if err := uc.mediaAssetRepo.CreatePending(ctx, asset); err != nil {
		return nil, err
	}

	return &MediaSignatureOutput{
		CloudName:    cloudName,
		APIKey:       apiKey,
		Timestamp:    timestamp,
		Folder:       folder,
		PublicID:     publicID,
		Signature:    signature,
		UploadURL:    fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/%s/upload", cloudName, input.ResourceType),
		ExpiresIn:    300,
		ResourceType: input.ResourceType,
	}, nil
}

func (uc *SocialUseCases) ConfirmMedia(ctx context.Context, userID uuid.UUID, input ConfirmMediaInput) (*ConfirmMediaOutput, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	publicID := strings.TrimSpace(input.PublicID)
	if publicID == "" {
		return nil, errors.Validation("public_id is required")
	}

	if !isSocialMediaPublicIDOwnedByUser(publicID, userID) {
		return nil, errors.Forbidden("media does not belong to current user")
	}

	var bytes *int64
	if input.Bytes > 0 {
		bytes = &input.Bytes
	}

	if err := uc.mediaAssetRepo.Confirm(ctx, userID, publicID, stringPointerOrNil(input.SecureURL), bytes); err != nil {
		return nil, err
	}

	return &ConfirmMediaOutput{
		PublicID:   publicID,
		IsOwned:    true,
		IsValid:    true,
		AssetState: "ready",
	}, nil
}

func (uc *SocialUseCases) CreatePost(ctx context.Context, userID uuid.UUID, input CreatePostInput) (*PostOutput, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	contentType := normalizeContentType(pickNonEmptyStringPtr(input.ContentType, input.ContentTypeAlt))
	if !isSupportedPostContentType(contentType) {
		return nil, errors.Validation("content_type must be one of [general workout_plan meal_log workout_session]")
	}

	var contentID *uuid.UUID
	if contentType != "general" {
		idStr := pickNonEmptyString(input.ContentID, input.ContentIDAlt)
		if idStr == "" {
			return nil, errors.Validation("content_id is required when content_type is not general")
		}
		parsed, err := uuid.Parse(idStr)
		if err != nil {
			return nil, errors.Validation("content_id must be a valid UUID")
		}
		contentID = &parsed
		if err := uc.validatePostContentRef(ctx, userID, contentType, parsed); err != nil {
			return nil, err
		}
	}

	caption := normalizeCaption(input.Caption)
	hasAttachment := contentType != "general"
	if caption == nil && len(input.Media) == 0 && !hasAttachment {
		return nil, errors.Validation("caption is required when media is empty")
	}

	feeling, err := normalizeOptionalText(input.Feeling, 100, "feeling")
	if err != nil {
		return nil, err
	}

	locationName, err := normalizeLocationName(input.Location)
	if err != nil {
		return nil, err
	}

	hashtags, err := normalizeHashtags(caption, input.Hashtags)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	post := &socialdomain.Post{
		ID:            uuid.New(),
		UserID:        userID,
		ContentType:   contentType,
		ContentID:     contentID,
		Caption:       caption,
		Feeling:       feeling,
		LocationName:  locationName,
		Hashtags:      hashtags,
		LikesCount:    0,
		CommentsCount: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	media := make([]socialdomain.PostMedia, 0, len(input.Media))
	for idx, item := range input.Media {
		publicID := strings.TrimSpace(item.PublicID)
		if publicID == "" {
			return nil, errors.Validation("public_id is required")
		}
		media = append(media, socialdomain.PostMedia{
			PostID:       post.ID,
			PublicID:     publicID,
			ResourceType: item.ResourceType,
			OrderIndex:   idx,
		})
	}

	if err := uc.postRepo.CreateWithMedia(ctx, post, media); err != nil {
		return nil, err
	}

	u, err := uc.userRepo.GetByID(ctx, userID)
	if err == nil {
		post.User = &socialdomain.PostUser{ID: u.ID, Name: u.Name, AvatarURL: u.AvatarURL}
	}

	mediaByPostID, err := uc.postRepo.GetMediaByPostIDs(ctx, []uuid.UUID{post.ID})
	if err != nil {
		return nil, err
	}
	post.Media = mediaByPostID[post.ID]

	out := uc.postOutputFrom(ctx, post, false, false, false)
	return &out, nil
}

func (uc *SocialUseCases) EditPost(ctx context.Context, userID uuid.UUID, postID string, input UpdatePostInput) (*PostOutput, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	clearLoc := input.ClearLocation != nil && *input.ClearLocation
	hasMediaReplace := input.Media != nil
	if input.Caption == nil && input.Feeling == nil && input.Location == nil && input.Hashtags == nil && !hasMediaReplace && !clearLoc {
		return nil, errors.Validation("at least one field must be provided")
	}

	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}

	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if post.UserID != userID {
		return nil, errors.Forbidden("you can only edit your own post")
	}

	mediaByPostID, err := uc.postRepo.GetMediaByPostIDs(ctx, []uuid.UUID{post.ID})
	if err != nil {
		return nil, err
	}
	post.Media = mediaByPostID[post.ID]

	if input.Caption != nil {
		post.Caption = normalizeCaption(input.Caption)
	}

	if input.Feeling != nil {
		normalizedFeeling, normalizeErr := normalizeOptionalText(input.Feeling, 100, "feeling")
		if normalizeErr != nil {
			return nil, normalizeErr
		}
		post.Feeling = normalizedFeeling
	}

	if clearLoc {
		post.LocationName = nil
	} else if input.Location != nil {
		locationName, normalizeErr := normalizeLocationName(input.Location)
		if normalizeErr != nil {
			return nil, normalizeErr
		}
		post.LocationName = locationName
	}

	if input.Hashtags != nil || input.Caption != nil {
		hashtagInput := post.Hashtags
		if input.Hashtags != nil {
			hashtagInput = *input.Hashtags
		}
		hashtags, normalizeErr := normalizeHashtags(post.Caption, hashtagInput)
		if normalizeErr != nil {
			return nil, normalizeErr
		}
		post.Hashtags = hashtags
	}

	var replaceMedia []socialdomain.PostMedia
	if hasMediaReplace {
		replaceMedia = make([]socialdomain.PostMedia, 0, len(*input.Media))
		for idx, item := range *input.Media {
			publicID := strings.TrimSpace(item.PublicID)
			if publicID == "" {
				return nil, errors.Validation("public_id is required")
			}
			rt := item.ResourceType
			if rt != "image" && rt != "video" {
				return nil, errors.Validation("resource_type must be one of [image video]")
			}
			replaceMedia = append(replaceMedia, socialdomain.PostMedia{
				PostID:       post.ID,
				PublicID:     publicID,
				ResourceType: rt,
				OrderIndex:   idx,
			})
		}
		post.Media = replaceMedia
	}

	hasAttachmentEdit := post.ContentType != "general"
	if post.Caption == nil && len(post.Media) == 0 && !hasAttachmentEdit {
		return nil, errors.Validation("caption is required when media is empty")
	}

	if hasMediaReplace {
		if err := uc.postRepo.UpdateWithMediaReplace(ctx, post, true, replaceMedia); err != nil {
			return nil, err
		}
	} else {
		if err := uc.postRepo.Update(ctx, post); err != nil {
			return nil, err
		}
	}

	updatedPost, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	mediaOut, err := uc.postRepo.GetMediaByPostIDs(ctx, []uuid.UUID{updatedPost.ID})
	if err != nil {
		return nil, err
	}
	updatedPost.Media = mediaOut[updatedPost.ID]

	out := uc.postOutputFrom(ctx, updatedPost, false, false, false)
	return &out, nil
}

func (uc *SocialUseCases) DeletePost(ctx context.Context, userID uuid.UUID, postID string) error {
	id, err := uuid.Parse(postID)
	if err != nil {
		return errors.BadRequest("invalid post id")
	}

	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return errors.Forbidden("you can only delete your own post")
	}

	return uc.postRepo.Delete(ctx, id)
}

func parseCloudinaryURL(raw string) (cloudName, apiKey, apiSecret string, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", "", err
	}

	if u.Scheme != "cloudinary" {
		return "", "", "", fmt.Errorf("invalid scheme")
	}

	cloudName = u.Host
	apiKey = u.User.Username()
	apiSecret, _ = u.User.Password()

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return "", "", "", fmt.Errorf("missing cloudinary credentials")
	}

	return cloudName, apiKey, apiSecret, nil
}

func cloudinaryImageDeliveryURL(cloudinaryRawURL, publicID string) string {
	publicID = strings.TrimSpace(publicID)
	if publicID == "" {
		return ""
	}
	cloudName, _, _, err := parseCloudinaryURL(strings.TrimSpace(cloudinaryRawURL))
	if err != nil || cloudName == "" {
		return ""
	}
	parts := strings.Split(publicID, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s", cloudName, strings.Join(parts, "/"))
}

func cloudinarySignature(params map[string]string, apiSecret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if strings.TrimSpace(params[k]) != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+params[k])
	}

	h := sha1.New()
	_, _ = h.Write([]byte(strings.Join(parts, "&") + apiSecret))
	return hex.EncodeToString(h.Sum(nil))
}

func randomHex(byteLength int) string {
	b := make([]byte, byteLength)
	_, err := rand.Read(b)
	if err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	return hex.EncodeToString(b)
}

func stringPointerOrNil(v string) *string {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func pickNonEmptyString(a, b *string) string {
	if a != nil {
		s := strings.TrimSpace(*a)
		if s != "" {
			return s
		}
	}
	if b != nil {
		return strings.TrimSpace(*b)
	}
	return ""
}

func pickNonEmptyStringPtr(a, b *string) *string {
	s := pickNonEmptyString(a, b)
	if s == "" {
		return nil
	}
	return &s
}

func normalizeContentType(contentType *string) string {
	if contentType == nil {
		return "general"
	}

	trimmed := strings.ToLower(strings.TrimSpace(*contentType))
	if trimmed == "" {
		return "general"
	}

	return trimmed
}

func isSupportedPostContentType(contentType string) bool {
	switch contentType {
	case "general", "workout_plan", "meal_log", "workout_session":
		return true
	default:
		return false
	}
}

func (uc *SocialUseCases) validatePostContentRef(ctx context.Context, userID uuid.UUID, contentType string, contentID uuid.UUID) error {
	switch contentType {
	case "meal_log":
		log, err := uc.mealLogRepo.GetByID(ctx, contentID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return errors.NotFound("meal log")
			}
			return errors.DatabaseError("get meal log", err)
		}
		if log.UserID != userID {
			return errors.Forbidden("meal log does not belong to you")
		}
	case "workout_plan":
		plan, err := uc.workoutPlanRepo.GetByID(ctx, contentID)
		if err != nil {
			return err
		}
		if plan.UserID != userID {
			return errors.Forbidden("workout plan does not belong to you")
		}
	case "workout_session":
		sess, err := uc.workoutSessionRepo.GetByID(ctx, contentID)
		if err != nil {
			return err
		}
		if sess.UserID != userID {
			return errors.Forbidden("workout session does not belong to you")
		}
		if sess.Status != workoutdomain.SessionStatusCompleted {
			return errors.Validation("only completed workout sessions can be shared")
		}
	}
	return nil
}

func (uc *SocialUseCases) postOutputFrom(ctx context.Context, post *socialdomain.Post, isLiked, isInterested, isNotInterested bool) PostOutput {
	out := mapPost(post, isLiked, isInterested, isNotInterested)
	if post == nil || post.ContentID == nil {
		return out
	}
	switch post.ContentType {
	case "meal_log":
		ml, err := uc.mealLogRepo.GetByID(ctx, *post.ContentID)
		if err != nil || ml.UserID != post.UserID {
			return out
		}
		out.MealSummary = &MealSummaryOutput{
			MealLogID:     ml.ID,
			LogDate:       ml.LogDate.Format("2006-01-02"),
			MealTime:      ml.MealTime,
			TotalCalories: ml.TotalCalories,
			TotalProteinG: ml.TotalProteinG,
			TotalCarbsG:   ml.TotalCarbsG,
			TotalFatG:     ml.TotalFatG,
			ItemCount:     len(ml.Items),
		}
	case "workout_session":
		s, err := uc.workoutSessionRepo.GetByID(ctx, *post.ContentID)
		if err != nil || s.UserID != post.UserID {
			return out
		}
		out.WorkoutSessionSummary = &WorkoutSessionSummaryOutput{
			SessionID:           s.ID,
			PlanTitle:           s.Title,
			ScheduledDate:       s.ScheduledDate,
			Status:              s.Status,
			DurationMins:        s.DurationMins,
			TotalCaloriesBurned: s.TotalCaloriesBurned,
			CompletedAt:         s.CompletedAt,
		}
	}
	return out
}

func normalizeCaption(caption *string) *string {
	if caption == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*caption)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func normalizeOptionalText(value *string, maxLen int, field string) (*string, error) {
	if value == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil
	}

	if utf8.RuneCountInString(trimmed) > maxLen {
		return nil, errors.Validation(fmt.Sprintf("%s must be at most %d characters long", field, maxLen))
	}

	return &trimmed, nil
}

func normalizeLocationName(location *CreatePostLocationInput) (*string, error) {
	if location == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(location.Name)
	if trimmed == "" {
		return nil, errors.Validation("location.name is required when location is provided")
	}

	if utf8.RuneCountInString(trimmed) > 255 {
		return nil, errors.Validation("location.name must be at most 255 characters long")
	}

	return &trimmed, nil
}

func normalizeHashtags(caption *string, input []string) ([]string, error) {
	canonical := make([]string, 0)
	seen := make(map[string]struct{})

	appendTag := func(raw string) error {
		normalized, err := normalizeHashtag(raw)
		if err != nil {
			return err
		}
		if normalized == "" {
			return nil
		}
		if _, exists := seen[normalized]; exists {
			return nil
		}
		seen[normalized] = struct{}{}
		canonical = append(canonical, normalized)
		return nil
	}

	for _, item := range input {
		if err := appendTag(item); err != nil {
			return nil, err
		}
	}

	if caption != nil {
		matches := hashtagPattern.FindAllStringSubmatch(*caption, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			if err := appendTag(match[1]); err != nil {
				return nil, err
			}
		}
	}

	return canonical, nil
}

func normalizeHashtag(raw string) (string, error) {
	trimmed := strings.TrimSpace(strings.ToLower(raw))
	trimmed = strings.TrimLeft(trimmed, "#")
	if trimmed == "" {
		return "", nil
	}

	runes := make([]rune, 0, len(trimmed))
	for _, r := range trimmed {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			runes = append(runes, r)
			continue
		}
		break
	}

	normalized := string(runes)
	if normalized == "" {
		return "", nil
	}

	if utf8.RuneCountInString(normalized) > 50 {
		return "", errors.Validation("hashtags must be at most 50 characters long")
	}

	return normalized, nil
}

func (uc *SocialUseCases) LikePost(ctx context.Context, userID uuid.UUID, postID string) (*LikeResponse, error) {
	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}

	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	exists, err := uc.likeRepo.Exists(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := uc.likeRepo.Create(ctx, &socialdomain.Like{ID: uuid.New(), PostID: id, UserID: userID, CreatedAt: time.Now()}); err != nil {
			return nil, err
		}
		if err := uc.postRepo.IncrementLikesCount(ctx, id); err != nil {
			return nil, err
		}
		if post.UserID != userID {
			uc.tryCreateLikeNotification(ctx, userID, post)
		}
	}

	post, err = uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &LikeResponse{LikeCount: post.LikesCount, IsLikedByMe: true}, nil
}

func (uc *SocialUseCases) UnlikePost(ctx context.Context, userID uuid.UUID, postID string) (*LikeResponse, error) {
	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}

	if _, err := uc.postRepo.GetByID(ctx, id); err != nil {
		return nil, err
	}

	exists, err := uc.likeRepo.Exists(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		if err := uc.likeRepo.Delete(ctx, id, userID); err != nil {
			return nil, err
		}
		if err := uc.postRepo.DecrementLikesCount(ctx, id); err != nil {
			return nil, err
		}
	}

	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &LikeResponse{LikeCount: post.LikesCount, IsLikedByMe: false}, nil
}

func (uc *SocialUseCases) CreateComment(ctx context.Context, userID uuid.UUID, postID string, input CreateCommentInput) (*CommentOutput, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	content := ""
	if input.Content != nil {
		content = strings.TrimSpace(*input.Content)
	}
	if content == "" && len(input.Media) == 0 {
		return nil, errors.Validation("content or media is required")
	}
	if len(input.Media) > 1 {
		return nil, errors.Validation("only one media item is allowed")
	}
	// Backward compatibility with existing DB constraint on comments.content:
	// media-only comments are stored with a single space as placeholder content.
	if content == "" && len(input.Media) > 0 {
		content = " "
	}

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}

	post, err := uc.postRepo.GetByID(ctx, postUUID)
	if err != nil {
		return nil, err
	}

	var parentCommentID *uuid.UUID
	if input.ParentID != nil && strings.TrimSpace(*input.ParentID) != "" {
		parsedParentID, err := uuid.Parse(strings.TrimSpace(*input.ParentID))
		if err != nil {
			return nil, errors.BadRequest("invalid parentId")
		}
		parentComment, err := uc.commentRepo.GetByID(ctx, parsedParentID)
		if err != nil {
			return nil, err
		}
		if parentComment.PostID != postUUID {
			return nil, errors.BadRequest("parent comment does not belong to post")
		}
		if parentComment.DeletedAt != nil {
			return nil, errors.Conflict("cannot reply to deleted comment")
		}
		parentCommentID = &parsedParentID
	}

	now := time.Now()
	comment := &socialdomain.Comment{
		ID:              uuid.New(),
		PostID:          postUUID,
		UserID:          userID,
		ParentCommentID: parentCommentID,
		Content:         content,
		ReplyCount:      0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	media := make([]socialdomain.CommentMedia, 0, len(input.Media))
	for idx, item := range input.Media {
		publicID := strings.TrimSpace(item.PublicID)
		if publicID == "" {
			return nil, errors.Validation("public_id is required")
		}
		if item.ResourceType != "image" {
			return nil, errors.Validation("comment media only supports image")
		}
		media = append(media, socialdomain.CommentMedia{
			CommentID:    comment.ID,
			PublicID:     publicID,
			ResourceType: item.ResourceType,
			OrderIndex:   idx,
		})
	}

	if err := uc.commentRepo.CreateWithMedia(ctx, comment, media); err != nil {
		return nil, err
	}
	mediaByCommentID, err := uc.commentRepo.GetMediaByCommentIDs(ctx, []uuid.UUID{comment.ID})
	if err != nil {
		return nil, err
	}
	loadedMedia := mediaByCommentID[comment.ID]
	if loadedMedia == nil {
		loadedMedia = []socialdomain.CommentMedia{}
	}
	if err := uc.postRepo.IncrementCommentsCount(ctx, postUUID); err != nil {
		return nil, err
	}
	if parentCommentID != nil {
		if err := uc.commentRepo.IncrementReplyCount(ctx, *parentCommentID); err != nil {
			return nil, err
		}
	}

	if post.UserID != userID {
		uc.tryCreateCommentNotification(ctx, userID, post)
	}

	author := CommentAuthorOutput{ID: userID}
	if u, err := uc.userRepo.GetByID(ctx, userID); err == nil {
		author.Name = u.Name
		if u.AvatarURL != nil {
			author.AvatarURL = *u.AvatarURL
		}
	}

	depth := 0
	if parentCommentID != nil {
		depth = 1
	}
	path := buildCommentPath(parentCommentID, comment.ID)

	out := &CommentOutput{
		ID:               comment.ID,
		PostID:           comment.PostID,
		ParentID:         comment.ParentCommentID,
		Depth:            depth,
		Path:             path,
		DirectReplyCount: comment.ReplyCount,
		PreviewReplies:   make([]CommentOutput, 0),
		Author:           author,
		Content:          comment.Content,
		Media:            mapCommentMediaOutput(uc.cloudinary.URL, loadedMedia),
		IsDeleted:        comment.DeletedAt != nil,
		IsEdited:         false,
		CreatedAt:        comment.CreatedAt,
		UpdatedAt:        comment.UpdatedAt,
	}

	uc.notify.PublishCommentCreated(socialnotify.CommentCreatedPayload{
		PostID:  postUUID.String(),
		Comment: commentOutputToRealtime(*out),
	})

	return out, nil
}

func (uc *SocialUseCases) commentDepthToRoot(ctx context.Context, postID uuid.UUID, commentID uuid.UUID) (int, error) {
	d := 0
	cur := commentID
	const maxHops = 512
	for d < maxHops {
		c, err := uc.commentRepo.GetByID(ctx, cur)
		if err != nil {
			return 0, err
		}
		if c.PostID != postID {
			return 0, errors.BadRequest("comment does not belong to post")
		}
		if c.ParentCommentID == nil {
			return d, nil
		}
		d++
		cur = *c.ParentCommentID
	}
	return 0, errors.InternalServer("comment thread too deep", nil)
}

func (uc *SocialUseCases) GetPostComments(ctx context.Context, postID string, cursor string, limit int) (*CommentListOutput, error) {
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}
	if _, err := uc.postRepo.GetByID(ctx, postUUID); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	page, err := decodePageCursor(cursor)
	if err != nil {
		return nil, errors.BadRequest("invalid cursor")
	}

	comments, total, err := uc.commentRepo.GetByPostID(ctx, postUUID, socialdomain.GetCommentsFilter{Page: page, PageSize: limit})
	if err != nil {
		return nil, err
	}

	parentIDs := make([]uuid.UUID, 0, len(comments))
	for _, comment := range comments {
		parentIDs = append(parentIDs, comment.ID)
	}
	latestRepliesByParent, err := uc.commentRepo.GetLatestRepliesByParentIDs(ctx, parentIDs, maxCommentPreviewRepliesPerParent)
	if err != nil {
		return nil, err
	}

	firstLevelReplyIDs := make([]uuid.UUID, 0)
	for _, replies := range latestRepliesByParent {
		for _, reply := range replies {
			firstLevelReplyIDs = append(firstLevelReplyIDs, reply.ID)
		}
	}
	nestedRepliesByParent := make(map[uuid.UUID][]socialdomain.Comment)
	if len(firstLevelReplyIDs) > 0 {
		nestedRepliesByParent, err = uc.commentRepo.GetLatestRepliesByParentIDs(ctx, firstLevelReplyIDs, maxCommentPreviewRepliesPerParent)
		if err != nil {
			return nil, err
		}
	}

	commentIDs := make([]uuid.UUID, 0, len(comments))
	for _, comment := range comments {
		commentIDs = append(commentIDs, comment.ID)
	}
	for _, replies := range latestRepliesByParent {
		for _, reply := range replies {
			commentIDs = append(commentIDs, reply.ID)
		}
	}
	for _, nestedList := range nestedRepliesByParent {
		for _, nr := range nestedList {
			commentIDs = append(commentIDs, nr.ID)
		}
	}
	mediaByCommentID, err := uc.commentRepo.GetMediaByCommentIDs(ctx, commentIDs)
	if err != nil {
		return nil, err
	}

	data := make([]CommentOutput, 0, len(comments))
	for _, comment := range comments {
		comment.Media = mediaByCommentID[comment.ID]
		latestReplies := make([]CommentOutput, 0)
		if replies, ok := latestRepliesByParent[comment.ID]; ok {
			latestReplies = make([]CommentOutput, 0, len(replies))
			for _, reply := range replies {
				reply.Media = mediaByCommentID[reply.ID]
				nestedOut := make([]CommentOutput, 0)
				if nestedList, ok2 := nestedRepliesByParent[reply.ID]; ok2 {
					nestedOut = make([]CommentOutput, 0, len(nestedList))
					for _, nr := range nestedList {
						nr.Media = mediaByCommentID[nr.ID]
						nestedPath := buildCommentPath(&reply.ID, nr.ID)
						nestedOut = append(nestedOut, mapCommentResponse(uc.cloudinary.URL, nr, 2, nestedPath, nil))
					}
				}
				replyPath := buildCommentPath(&comment.ID, reply.ID)
				latestReplies = append(latestReplies, mapCommentResponse(uc.cloudinary.URL, reply, 1, replyPath, nestedOut))
			}
		}

		data = append(data, mapCommentResponse(uc.cloudinary.URL, comment, 0, buildCommentPath(nil, comment.ID), latestReplies))
	}

	hasMore := int64(page*limit) < total
	var nextCursor *string
	if hasMore {
		cursorValue := encodePageCursor(page + 1)
		nextCursor = &cursorValue
	}
	out := &CommentListOutput{Comments: data, NextCursor: nextCursor}
	return out, nil
}

func (uc *SocialUseCases) GetCommentReplies(ctx context.Context, postID string, commentID string, cursor string, limit int) (*CommentRepliesOutput, error) {
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}
	parentUUID, err := uuid.Parse(commentID)
	if err != nil {
		return nil, errors.BadRequest("invalid comment id")
	}

	parentComment, err := uc.commentRepo.GetByID(ctx, parentUUID)
	if err != nil {
		return nil, err
	}
	if parentComment.PostID != postUUID {
		return nil, errors.BadRequest("comment does not belong to post")
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 2000 {
		limit = 2000
	}

	page, err := decodePageCursor(cursor)
	if err != nil {
		return nil, errors.BadRequest("invalid cursor")
	}

	parentDepth, err := uc.commentDepthToRoot(ctx, postUUID, parentUUID)
	if err != nil {
		return nil, err
	}

	comments, total, err := uc.commentRepo.GetByPostID(ctx, postUUID, socialdomain.GetCommentsFilter{
		Page:            page,
		PageSize:        limit,
		ParentCommentID: &parentUUID,
	})
	if err != nil {
		return nil, err
	}

	replyRowIDs := make([]uuid.UUID, 0, len(comments))
	for _, comment := range comments {
		replyRowIDs = append(replyRowIDs, comment.ID)
	}
	nestedRepliesByParent := make(map[uuid.UUID][]socialdomain.Comment)
	if len(replyRowIDs) > 0 {
		nestedRepliesByParent, err = uc.commentRepo.GetLatestRepliesByParentIDs(ctx, replyRowIDs, maxCommentPreviewRepliesPerParent)
		if err != nil {
			return nil, err
		}
	}

	commentIDs := make([]uuid.UUID, 0, len(comments))
	for _, comment := range comments {
		commentIDs = append(commentIDs, comment.ID)
	}
	for _, nestedList := range nestedRepliesByParent {
		for _, nr := range nestedList {
			commentIDs = append(commentIDs, nr.ID)
		}
	}
	mediaByCommentID, err := uc.commentRepo.GetMediaByCommentIDs(ctx, commentIDs)
	if err != nil {
		return nil, err
	}

	data := make([]CommentOutput, 0, len(comments))
	for _, comment := range comments {
		comment.Media = mediaByCommentID[comment.ID]
		nestedOut := make([]CommentOutput, 0)
		if nestedList, ok := nestedRepliesByParent[comment.ID]; ok {
			nestedOut = make([]CommentOutput, 0, len(nestedList))
			for _, nr := range nestedList {
				nr.Media = mediaByCommentID[nr.ID]
				nestedPath := buildCommentPath(&comment.ID, nr.ID)
				nestedOut = append(nestedOut, mapCommentResponse(uc.cloudinary.URL, nr, parentDepth+2, nestedPath, nil))
			}
		}
		data = append(data, mapCommentResponse(uc.cloudinary.URL, comment, parentDepth+1, buildCommentPath(&parentUUID, comment.ID), nestedOut))
	}

	hasMore := int64(page*limit) < total
	var nextCursor *string
	if hasMore {
		cursorValue := encodePageCursor(page + 1)
		nextCursor = &cursorValue
	}
	out := &CommentRepliesOutput{Replies: data, NextCursor: nextCursor}
	return out, nil
}

func (uc *SocialUseCases) UpdateComment(ctx context.Context, userID uuid.UUID, postID string, commentID string, input UpdateCommentInput) (*CommentOutput, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}
	commentUUID, err := uuid.Parse(commentID)
	if err != nil {
		return nil, errors.BadRequest("invalid comment id")
	}

	if _, err := uc.postRepo.GetByID(ctx, postUUID); err != nil {
		return nil, err
	}

	comment, err := uc.commentRepo.GetByID(ctx, commentUUID)
	if err != nil {
		return nil, err
	}
	if comment.PostID != postUUID {
		return nil, errors.BadRequest("comment does not belong to post")
	}
	if comment.DeletedAt != nil {
		return nil, errors.NotFound("comment")
	}
	if comment.UserID != userID {
		return nil, errors.Forbidden("you can only edit your own comment")
	}
	if time.Since(comment.CreatedAt) > socialCommentEditWindow {
		return nil, errors.Forbidden("comment edit window expired")
	}

	content := strings.TrimSpace(input.Content)
	if content == "" {
		return nil, errors.Validation("content is required")
	}

	comment.Content = content
	if err := uc.commentRepo.Update(ctx, comment); err != nil {
		return nil, err
	}

	updated, err := uc.commentRepo.GetByID(ctx, commentUUID)
	if err != nil {
		return nil, err
	}
	mediaByCommentID, err := uc.commentRepo.GetMediaByCommentIDs(ctx, []uuid.UUID{commentUUID})
	if err != nil {
		return nil, err
	}
	updated.Media = mediaByCommentID[commentUUID]
	if updated.Media == nil {
		updated.Media = []socialdomain.CommentMedia{}
	}

	depth, err := uc.commentDepthToRoot(ctx, postUUID, commentUUID)
	if err != nil {
		return nil, err
	}
	path := buildCommentPath(updated.ParentCommentID, updated.ID)
	out := mapCommentResponse(uc.cloudinary.URL, *updated, depth, path, make([]CommentOutput, 0))

	uc.notify.PublishCommentUpdated(socialnotify.CommentUpdatedPayload{
		PostID:  postUUID.String(),
		Comment: commentOutputToRealtime(out),
	})

	return &out, nil
}

func (uc *SocialUseCases) DeleteComment(ctx context.Context, userID uuid.UUID, postID string, commentID string) (*DeleteCommentResult, error) {
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}
	commentUUID, err := uuid.Parse(commentID)
	if err != nil {
		return nil, errors.BadRequest("invalid comment id")
	}

	post, err := uc.postRepo.GetByID(ctx, postUUID)
	if err != nil {
		return nil, err
	}

	comment, err := uc.commentRepo.GetByID(ctx, commentUUID)
	if err != nil {
		return nil, err
	}
	if comment.PostID != postUUID {
		return nil, errors.BadRequest("comment does not belong to post")
	}
	if comment.DeletedAt != nil {
		return &DeleteCommentResult{DeletedByRole: "none"}, nil
	}

	isAuthor := comment.UserID == userID
	isPostOwner := post.UserID == userID
	if !isAuthor && !isPostOwner {
		return nil, errors.Forbidden("you cannot delete this comment")
	}

	role := "author"
	if isPostOwner && !isAuthor {
		role = "post_owner"
	}

	var parentStr *string
	if comment.ParentCommentID != nil {
		s := comment.ParentCommentID.String()
		parentStr = &s
	}

	if err := uc.commentRepo.Delete(ctx, commentUUID); err != nil {
		return nil, err
	}
	if err := uc.postRepo.DecrementCommentsCount(ctx, postUUID); err != nil {
		return nil, err
	}
	if comment.ParentCommentID != nil {
		if err := uc.commentRepo.DecrementReplyCount(ctx, *comment.ParentCommentID); err != nil {
			return nil, err
		}
	}

	uc.notify.PublishCommentDeleted(socialnotify.CommentDeletedPayload{
		PostID:          postUUID.String(),
		CommentID:       commentUUID.String(),
		ParentID:        parentStr,
		DeletedByUserID: userID.String(),
	})
	return &DeleteCommentResult{DeletedByRole: role}, nil
}

func commentOutputToRealtime(c CommentOutput) socialnotify.RealtimeComment {
	var parent *string
	if c.ParentID != nil {
		s := c.ParentID.String()
		parent = &s
	}
	previews := make([]socialnotify.RealtimeComment, 0, len(c.PreviewReplies))
	for _, p := range c.PreviewReplies {
		previews = append(previews, commentOutputToRealtime(p))
	}
	media := make([]socialnotify.RealtimeCommentMedia, 0, len(c.Media))
	for _, m := range c.Media {
		media = append(media, socialnotify.RealtimeCommentMedia{Type: m.Type, URL: m.URL})
	}
	return socialnotify.RealtimeComment{
		ID:               c.ID.String(),
		PostID:           c.PostID.String(),
		ParentID:         parent,
		Depth:            c.Depth,
		Path:             c.Path,
		DirectReplyCount: c.DirectReplyCount,
		PreviewReplies:   previews,
		Author: socialnotify.RealtimeCommentAuthor{
			ID:        c.Author.ID.String(),
			Name:      c.Author.Name,
			AvatarURL: c.Author.AvatarURL,
		},
		Content:   c.Content,
		Media:     media,
		IsDeleted: c.IsDeleted,
		IsEdited:  c.IsEdited,
		CreatedAt: c.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt: c.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func mapCommentAuthor(comment socialdomain.Comment) CommentAuthorOutput {
	author := CommentAuthorOutput{ID: comment.UserID}
	if comment.User != nil {
		author.ID = comment.User.ID
		author.Name = comment.User.Name
		if comment.User.AvatarURL != nil {
			author.AvatarURL = *comment.User.AvatarURL
		}
	}
	return author
}

func mapCommentResponse(cloudinaryRawURL string, comment socialdomain.Comment, depth int, path string, previewReplies []CommentOutput) CommentOutput {
	if previewReplies == nil {
		previewReplies = make([]CommentOutput, 0)
	}

	isEdited := comment.UpdatedAt.After(comment.CreatedAt)
	return CommentOutput{
		ID:               comment.ID,
		PostID:           comment.PostID,
		ParentID:         comment.ParentCommentID,
		Depth:            depth,
		Path:             path,
		DirectReplyCount: comment.ReplyCount,
		PreviewReplies:   previewReplies,
		Author:           mapCommentAuthor(comment),
		Content:          comment.Content,
		Media:            mapCommentMediaOutput(cloudinaryRawURL, comment.Media),
		IsDeleted:        comment.DeletedAt != nil,
		IsEdited:         isEdited,
		CreatedAt:        comment.CreatedAt,
		UpdatedAt:        comment.UpdatedAt,
	}
}

func mapCommentMediaOutput(cloudinaryRawURL string, media []socialdomain.CommentMedia) []PostMediaOutput {
	if len(media) == 0 {
		return make([]PostMediaOutput, 0)
	}
	out := make([]PostMediaOutput, 0, len(media))
	for _, item := range media {
		urlStr := ""
		if item.SecureURL != nil {
			urlStr = strings.TrimSpace(*item.SecureURL)
		}
		if urlStr == "" && item.ResourceType == "image" {
			urlStr = cloudinaryImageDeliveryURL(cloudinaryRawURL, item.PublicID)
		}
		if urlStr == "" {
			continue
		}
		out = append(out, PostMediaOutput{
			Type: item.ResourceType,
			URL:  urlStr,
		})
	}
	return out
}

func buildCommentPath(parentID *uuid.UUID, commentID uuid.UUID) string {
	if parentID == nil {
		return commentID.String()
	}
	return parentID.String() + "/" + commentID.String()
}

func notificationCategoryFromContentType(ct string) string {
	switch ct {
	case "workout_plan", "workout_session":
		return "workouts"
	case "meal_log":
		return "nutrition"
	default:
		return "social"
	}
}

func dayGroupUTC(t time.Time) string {
	now := time.Now().UTC()
	u := t.UTC()
	y1, m1, d1 := now.Date()
	y2, m2, d2 := u.Date()
	if y1 == y2 && m1 == m2 && d1 == d2 {
		return "today"
	}
	prev := now.AddDate(0, 0, -1)
	y3, m3, d3 := prev.Date()
	if y2 == y3 && m2 == m3 && d2 == d3 {
		return "yesterday"
	}
	return "earlier"
}

func socialDisplayName(name string) string {
	s := strings.TrimSpace(name)
	if s == "" {
		return "Ai đó"
	}
	return s
}

func notificationKindFromRecord(title string, postID *uuid.UUID) string {
	t := strings.TrimSpace(strings.ToLower(title))
	switch t {
	case "new like", "lượt thích":
		return "like"
	case "new comment", "bình luận mới":
		return "comment"
	case "new follower", "theo dõi mới":
		return "follow"
	}
	if strings.Contains(t, "thích") || strings.Contains(t, "liked") {
		return "like"
	}
	if strings.Contains(t, "bình luận") || strings.Contains(t, "comment") {
		return "comment"
	}
	if strings.Contains(t, "theo dõi") || strings.Contains(t, "follow") {
		return "follow"
	}
	if postID != nil {
		return "post"
	}
	return ""
}

func notificationMetaWithPostID(text string, postID uuid.UUID) string {
	return text + "|" + postID.String()
}

func effectivePostIDForNotification(n *socialdomain.InAppNotification) *uuid.UUID {
	if n == nil {
		return nil
	}
	if n.PostID != nil {
		return n.PostID
	}
	return n.RelatedPostID
}

func (uc *SocialUseCases) pushNotificationRealtime(ctx context.Context, recipient uuid.UUID, n *socialdomain.InAppNotification) {
	if uc.notify == nil || n == nil {
		return
	}
	eff := effectivePostIDForNotification(n)
	var postIDStr *string
	if eff != nil {
		s := eff.String()
		postIDStr = &s
	}
	kind := notificationKindFromRecord(n.Title, eff)
	uc.notify.PublishNotificationCreated(recipient, socialnotify.NotificationPayload{
		ID:        n.ID.String(),
		Type:      n.Type,
		Title:     n.Title,
		Meta:      n.Meta,
		DayGroup:  dayGroupUTC(n.CreatedAt),
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt.UTC().Format(time.RFC3339),
		PostID:    postIDStr,
		Kind:      kind,
	})
	unread, err := uc.notifRepo.CountUnread(ctx, recipient)
	if err != nil {
		return
	}
	uc.notify.PublishUnread(recipient, unread)
}

func (uc *SocialUseCases) tryCreateFollowNotification(ctx context.Context, followerID, followingID uuid.UUID) {
	follower, err := uc.userRepo.GetByID(ctx, followerID)
	if err != nil || follower == nil {
		return
	}
	name := socialDisplayName(follower.Name)
	n := &socialdomain.InAppNotification{
		ID:        uuid.New(),
		UserID:    followingID,
		Type:      "social",
		Title:     "Theo dõi mới",
		Meta:      name + " đã bắt đầu theo dõi bạn",
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	if err := uc.notifRepo.Create(ctx, n); err != nil {
		return
	}
	uc.pushNotificationRealtime(ctx, followingID, n)
}

func (uc *SocialUseCases) tryCreateLikeNotification(ctx context.Context, likerID uuid.UUID, post *socialdomain.Post) {
	liker, err := uc.userRepo.GetByID(ctx, likerID)
	if err != nil {
		return
	}
	pid := post.ID
	likerName := socialDisplayName(liker.Name)
	n := &socialdomain.InAppNotification{
		ID:            uuid.New(),
		UserID:        post.UserID,
		Type:          notificationCategoryFromContentType(post.ContentType),
		Title:         "Lượt thích",
		Meta:          notificationMetaWithPostID(likerName+" đã thích bài viết của bạn", pid),
		PostID:        &pid,
		RelatedPostID: &pid,
		IsRead:        false,
		CreatedAt:     time.Now(),
	}
	if err := uc.notifRepo.Create(ctx, n); err != nil {
		return
	}
	uc.pushNotificationRealtime(ctx, post.UserID, n)
}

func (uc *SocialUseCases) tryCreateCommentNotification(ctx context.Context, commenterID uuid.UUID, post *socialdomain.Post) {
	commenter, err := uc.userRepo.GetByID(ctx, commenterID)
	if err != nil {
		return
	}
	pid := post.ID
	commenterName := socialDisplayName(commenter.Name)
	n := &socialdomain.InAppNotification{
		ID:            uuid.New(),
		UserID:        post.UserID,
		Type:          notificationCategoryFromContentType(post.ContentType),
		Title:         "Bình luận mới",
		Meta:          notificationMetaWithPostID(commenterName+" đã bình luận bài viết của bạn", pid),
		PostID:        &pid,
		RelatedPostID: &pid,
		IsRead:        false,
		CreatedAt:     time.Now(),
	}
	if err := uc.notifRepo.Create(ctx, n); err != nil {
		return
	}
	uc.pushNotificationRealtime(ctx, post.UserID, n)
}

func (uc *SocialUseCases) ListNotifications(ctx context.Context, userID uuid.UUID, filter string, cursor string, limit int) (*NotificationsListPayload, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	page, err := decodePageCursor(cursor)
	if err != nil {
		return nil, errors.BadRequest("invalid cursor")
	}
	f := strings.TrimSpace(strings.ToLower(filter))
	if f != "all" && f != "workouts" && f != "social" && f != "nutrition" {
		f = "all"
	}
	rows, total, err := uc.notifRepo.ListForUser(ctx, userID, f, page, limit)
	if err != nil {
		return nil, err
	}
	out := make([]NotificationRowOutput, 0, len(rows))
	for _, r := range rows {
		eff := effectivePostIDForNotification(&r)
		var postIDStr *string
		if eff != nil {
			s := eff.String()
			postIDStr = &s
		}
		kind := notificationKindFromRecord(r.Title, eff)
		out = append(out, NotificationRowOutput{
			ID:        r.ID,
			Type:      r.Type,
			Title:     r.Title,
			Meta:      r.Meta,
			DayGroup:  dayGroupUTC(r.CreatedAt),
			IsRead:    r.IsRead,
			CreatedAt: r.CreatedAt,
			PostID:    postIDStr,
			Kind:      kind,
		})
	}
	hasMore := int64(page*limit) < total
	var next *string
	if hasMore {
		v := encodePageCursor(page + 1)
		next = &v
	}
	return &NotificationsListPayload{
		Data: out,
		Pagination: notificationsListPagination{
			NextCursor: next,
			HasMore:    hasMore,
		},
	}, nil
}

func (uc *SocialUseCases) GetNotificationsUnreadCount(ctx context.Context, userID uuid.UUID) (*UnreadNotificationsCountPayload, error) {
	n, err := uc.notifRepo.CountUnread(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &UnreadNotificationsCountPayload{Unread: n}, nil
}

func (uc *SocialUseCases) MarkNotificationsRead(ctx context.Context, userID uuid.UUID, rawIDs []string) (*MarkNotificationsReadPayload, error) {
	ids := make([]uuid.UUID, 0, len(rawIDs))
	for _, s := range rawIDs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, errors.BadRequest("invalid notification id")
		}
		ids = append(ids, id)
	}
	n, err := uc.notifRepo.MarkRead(ctx, userID, ids)
	if err != nil {
		return nil, err
	}
	if uc.notify != nil {
		unread, err := uc.notifRepo.CountUnread(ctx, userID)
		if err == nil {
			uc.notify.PublishUnread(userID, unread)
		}
	}
	return &MarkNotificationsReadPayload{Updated: n}, nil
}
