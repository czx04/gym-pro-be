package social

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"gym-pro-2026-ptit/internal/config"
	socialdomain "gym-pro-2026-ptit/internal/domain/social"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/pkg/errors"
	"gym-pro-2026-ptit/pkg/validator"

	"github.com/google/uuid"
)

type SocialUseCases struct {
	postRepo       socialdomain.PostRepository
	followRepo     socialdomain.FollowRepository
	mediaAssetRepo socialdomain.MediaAssetRepository
	userRepo       user.Repository
	validator      *validator.Validator
	cloudinary     config.CloudinaryConfig
}

func NewSocialUseCases(
	cfg *config.Config,
	postRepo socialdomain.PostRepository,
	followRepo socialdomain.FollowRepository,
	mediaAssetRepo socialdomain.MediaAssetRepository,
	userRepo user.Repository,
	validator *validator.Validator,
) *SocialUseCases {
	return &SocialUseCases{
		postRepo:       postRepo,
		followRepo:     followRepo,
		mediaAssetRepo: mediaAssetRepo,
		userRepo:       userRepo,
		validator:      validator,
		cloudinary:     cfg.Cloudinary,
	}
}

type CreatePostInput struct {
	Caption     *string                             `json:"caption,omitempty" validate:"omitempty,max=2000"`
	Media       []socialdomain.CreatePostMediaInput `json:"media,omitempty" validate:"omitempty,dive"`
	ContentType *string                             `json:"content_type,omitempty" validate:"omitempty,oneof=workout_plan meal_log"`
	ContentID   *uuid.UUID                          `json:"content_id,omitempty"`
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

type AuthorOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
}

type PostMediaOutput struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type PostOutput struct {
	ID              uuid.UUID         `json:"id"`
	Author          AuthorOutput      `json:"author"`
	StreakText      string            `json:"streak_text"`
	TimeLabel       string            `json:"time_label"`
	Caption         string            `json:"caption"`
	Media           []PostMediaOutput `json:"media"`
	LikeCount       int               `json:"like_count"`
	CommentCount    int               `json:"comment_count"`
	IsLikedByMe     bool              `json:"is_liked_by_me"`
	SharedExercises []interface{}     `json:"shared_exercises,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
}

type UserProfileOutput struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	AvatarURL      *string   `json:"avatar_url,omitempty"`
	Subtitle       string    `json:"subtitle"`
	StreakValue    int       `json:"streak_value"`
	PostsCount     int64     `json:"posts_count"`
	FollowersCount int       `json:"followers_count"`
	IsFollowing    bool      `json:"is_following"`
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
		data = append(data, mapPost(item.Post, item.IsLiked))
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

func (uc *SocialUseCases) GetPostByID(ctx context.Context, _ uuid.UUID, postID string) (*PostOutput, error) {
	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.BadRequest("invalid post id")
	}

	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	mediaByPostID, err := uc.postRepo.GetMediaByPostIDs(ctx, []uuid.UUID{post.ID})
	if err != nil {
		return nil, err
	}
	post.Media = mediaByPostID[post.ID]

	out := mapPost(post, false)
	return &out, nil
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
		StreakValue:    0,
		PostsCount:     totalPosts,
		FollowersCount: stats.FollowersCount,
		IsFollowing:    isFollowing,
	}, nil
}

func (uc *SocialUseCases) GetUserPosts(ctx context.Context, _ uuid.UUID, profileUserID, cursor string, limit int) (*FeedOutput, error) {
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

	data := make([]PostOutput, 0, len(posts))
	for _, post := range posts {
		post.Media = mediaByPostID[post.ID]
		data = append(data, mapPost(&post, false))
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

func mapPost(post *socialdomain.Post, isLiked bool) PostOutput {
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
		media = append(media, PostMediaOutput{Type: item.ResourceType, URL: urlValue})
	}

	return PostOutput{
		ID:              post.ID,
		Author:          author,
		StreakText:      "0 DAY STREAK",
		TimeLabel:       humanizeTime(post.CreatedAt),
		Caption:         caption,
		Media:           media,
		LikeCount:       post.LikesCount,
		CommentCount:    post.CommentsCount,
		IsLikedByMe:     isLiked,
		SharedExercises: []interface{}{},
		CreatedAt:       post.CreatedAt,
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

func (uc *SocialUseCases) CreateMediaSignature(ctx context.Context, userID uuid.UUID, input CreateMediaSignatureInput) (*MediaSignatureOutput, error) {
	if err := uc.validator.Validate(input); err != nil {
		return nil, errors.Validation(err.Error())
	}

	if strings.TrimSpace(uc.cloudinary.URL) == "" {
		return nil, errors.InternalServer("cloudinary is not configured", nil)
	}

	ownerPrefix := "posts/" + userID.String()
	folder := strings.TrimSpace(input.Folder)
	if folder == "" {
		folder = ownerPrefix
	}
	if !strings.HasPrefix(folder, ownerPrefix) {
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

	ownerPrefix := "posts/" + userID.String() + "/"
	isOwned := strings.HasPrefix(publicID, ownerPrefix)
	if !isOwned {
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

	contentType := "workout_plan"
	if input.ContentType != nil {
		contentType = *input.ContentType
	}

	contentID := uuid.Nil
	if input.ContentID != nil {
		contentID = *input.ContentID
	}

	now := time.Now()
	post := &socialdomain.Post{
		ID:            uuid.New(),
		UserID:        userID,
		ContentType:   contentType,
		ContentID:     contentID,
		Caption:       input.Caption,
		LikesCount:    0,
		CommentsCount: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	media := make([]socialdomain.PostMedia, 0, len(input.Media))
	for idx, item := range input.Media {
		media = append(media, socialdomain.PostMedia{
			PostID:       post.ID,
			PublicID:     strings.TrimSpace(item.PublicID),
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

	out := mapPost(post, false)
	return &out, nil
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
