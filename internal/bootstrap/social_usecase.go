package bootstrap

import (
	"gym-pro-2026-ptit/internal/config"
	mealdomain "gym-pro-2026-ptit/internal/domain/meal"
	socialdomain "gym-pro-2026-ptit/internal/domain/social"
	"gym-pro-2026-ptit/internal/domain/user"
	workoutdomain "gym-pro-2026-ptit/internal/domain/workout"
	"gym-pro-2026-ptit/internal/port/socialnotify"
	socialuc "gym-pro-2026-ptit/internal/usecase/social"
	"gym-pro-2026-ptit/pkg/validator"
)

func ProvideSocialUseCases(
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
	v *validator.Validator,
	b socialnotify.Broadcaster,
) *socialuc.SocialUseCases {
	return socialuc.NewSocialUseCases(cfg, postRepo, followRepo, likeRepo, commentRepo, mediaAssetRepo, preferenceRepo, reportRepo, blockRepo, notifRepo, userRepo, mealLogRepo, workoutSessionRepo, workoutPlanRepo, v, b)
}
