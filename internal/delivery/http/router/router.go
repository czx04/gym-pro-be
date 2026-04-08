package router

import (
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/delivery/http/handler"
	"gym-pro-2026-ptit/internal/delivery/http/middleware"
	"gym-pro-2026-ptit/internal/delivery/http/websocket"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router wraps Gin engine
type Router struct {
	engine *gin.Engine
}

// New creates a new router
func New(
	cfg *config.Config,
	authMiddleware middleware.AuthMiddleware,
	wsHub *websocket.Hub,
	authHandler *handler.AuthHandler,
	workoutHandler *handler.WorkoutHandler,
	exerciseHandler *handler.ExerciseHandler,
	foodHandler *handler.FoodHandler,
	recipeHandler *handler.RecipeHandler,
	mealLogHandler *handler.MealLogHandler,
	mealDailyHandler *handler.MealDailyHandler,
	userHandler *handler.UserHandler,
	socialHandler *handler.SocialHandler,
) *Router {
	gin.SetMode(cfg.Server.GinMode)

	engine := gin.New()

	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.CORSMiddleware(&cfg.Server))
	engine.Use(middleware.ErrorHandlerMiddleware())

	engine.GET("/health", healthCheckHandler)
	engine.GET("/ping", pingHandler)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := engine.Group("/api/v1")
	{
		wsHub.RegisterRoutes(v1)

		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register/request", authHandler.RegisterRequestOTP)
			authRoutes.POST("/register/verify", authHandler.VerifyOTP)

			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.RefreshToken)
			authRoutes.POST("/reset-password/request", authHandler.ResetPasswordRequestOTP)
			authRoutes.POST("/reset-password/verify", authHandler.VerifyOTPForgotPassword)
			authRoutes.POST("/reset-password", authHandler.ResetPassword)
		}

		authenticated := v1.Group("")
		authenticated.Use(gin.HandlerFunc(authMiddleware))
		{
			users := authenticated.Group("/users")
			{
				users.GET("/me", authHandler.GetMe)
				users.PUT("/me", authHandler.UpdateMe)
				users.GET("/me/weight-history", userHandler.GetMyWeightHistory)
				users.GET("/me/meal-streak", userHandler.GetMealStreak)
				users.GET("/me/workout-stats", userHandler.GetMyWorkoutStats)
				users.POST("/me/email/request-otp", userHandler.RequestChangeEmailOTP)
				users.POST("/me/email/verify", userHandler.VerifyChangeEmailOTP)
				users.POST("/me/push-token", userHandler.RegisterPushToken)
				users.PUT("/me/avatar", userHandler.UploadAvatarImage)
				users.GET("/nutrition-target", userHandler.GetUserNutritionTarget)
				users.PUT("/nutrition-target", userHandler.UpdateUserNutritionTarget)
				users.DELETE("/me", userHandler.DeleteAccount)
			}

			exercises := authenticated.Group("/exercises")
			{
				exercises.GET("", exerciseHandler.ListExercises)
				exercises.GET("/:id/stats", exerciseHandler.GetExerciseStats)
				exercises.GET("/:id", exerciseHandler.GetExercise)
			}

			workoutPlans := authenticated.Group("/workout-plans")
			{
				workoutPlans.POST("", workoutHandler.CreateWorkoutPlan)
				workoutPlans.GET("", workoutHandler.ListWorkoutPlans)
				workoutPlans.GET("/:id", workoutHandler.GetWorkoutPlan)
				workoutPlans.PUT("/:id", workoutHandler.UpdateWorkoutPlan)
				workoutPlans.DELETE("/:id", workoutHandler.DeleteWorkoutPlan)
			}

			workoutSessions := authenticated.Group("/workout-sessions")
			{
				workoutSessions.GET("/scheduled-dates", workoutHandler.GetScheduledDates)
				workoutSessions.GET("/weekly-summary", workoutHandler.GetWeeklyWorkoutSummary)
				workoutSessions.GET("", workoutHandler.GetSessionsByDate)
				workoutSessions.GET("/:id", workoutHandler.GetSessionByID)
				workoutSessions.DELETE("/:id", workoutHandler.DeleteWorkoutSession)
				workoutSessions.POST("", workoutHandler.CreateWorkoutSession)
				workoutSessions.PATCH("/:id/finish", workoutHandler.FinishWorkoutSession)
			}

			foods := authenticated.Group("/foods")
			{
				foods.GET("", foodHandler.ListFoods)
				foods.GET("/:id", foodHandler.GetFood)
				foods.GET("/search", foodHandler.SearchFoods)
				foods.POST("", foodHandler.CreateFood)
				foods.PUT("/:id", foodHandler.UpdateFood)
				foods.DELETE("/:id", foodHandler.DeleteFood)
				foods.POST("/scan", foodHandler.ScanFood)
			}

			recipes := authenticated.Group("/recipes")
			{
				recipes.POST("", recipeHandler.CreateRecipe)
				recipes.GET("", recipeHandler.ListRecipes)
				recipes.GET("/:id", recipeHandler.GetRecipe)
				recipes.PUT("/:id", recipeHandler.UpdateRecipe)
				recipes.DELETE("/:id", recipeHandler.DeleteRecipe)
			}

			mealLogs := authenticated.Group("/meal-logs")
			{
				mealLogs.POST("", mealLogHandler.CreateMealLog)
				mealLogs.GET("/stats", mealLogHandler.GetNutritionStats)
				mealLogs.GET("/logged-dates", mealLogHandler.ListLoggedDates)
				mealLogs.GET("/date/:date", mealLogHandler.GetMealLogsByDate)
				mealLogs.GET("/:id", mealLogHandler.GetMealLog)
				mealLogs.PUT("/:id", mealLogHandler.UpdateMealLog)
				mealLogs.DELETE("/:id", mealLogHandler.DeleteMealLog)
			}

			mealDaily := authenticated.Group("/meal-daily")
			{
				mealDaily.GET("/date/:date", mealDailyHandler.GetMealDailyTargetByDate)
			}

			social := authenticated.Group("/social")
			{
				social.GET("/notifications", socialHandler.SocialNotifications)
				social.POST("/notifications", socialHandler.SocialNotificationsWrite)
				social.GET("/search", socialHandler.Search)
				social.GET("/feed", socialHandler.GetFeed)
				social.PUT("/posts/:postId", socialHandler.UpdatePost)
				social.PUT("/posts/:postId/preference", socialHandler.SetPostPreference)
				social.PUT("/users/:userId/follow", socialHandler.SetFollowState)
				social.PUT("/users/:userId/block", socialHandler.SetBlockState)
				social.POST("/posts", socialHandler.CreatePost)
				social.DELETE("/posts/:postId", socialHandler.DeletePost)
				social.GET("/posts/:postId/attachment", socialHandler.GetPostAttachment)
				social.GET("/posts/:postId", socialHandler.GetPostByID)
				social.POST("/posts/:postId/reports", socialHandler.ReportPost)
				social.GET("/posts/:postId/comments", socialHandler.GetPostComments)
				social.GET("/posts/:postId/comments/:commentId/replies", socialHandler.GetCommentReplies)
				social.POST("/posts/:postId/comments", socialHandler.CreateComment)
				social.PUT("/posts/:postId/comments/:commentId", socialHandler.UpdateComment)
				social.DELETE("/posts/:postId/comments/:commentId", socialHandler.DeleteComment)
				social.GET("/users/:userId/profile", socialHandler.GetUserProfile)
				social.GET("/users/:userId/posts", socialHandler.GetUserPosts)
				social.POST("/media/signature", socialHandler.CreateMediaSignature)
				social.POST("/media/confirm", socialHandler.ConfirmMedia)
			}
		}
	}

	return &Router{engine: engine}
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "Service is healthy",
	})
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
