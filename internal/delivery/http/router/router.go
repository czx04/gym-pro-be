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
		//v1.Use(middleware.RateLimitMiddleware(&cfg.RateLimit))

		wsHub.RegisterRoutes(v1)

		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register/request", authHandler.RegisterRequestOTP)
			authRoutes.POST("/register/verify", authHandler.VerifyOTP)

			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.RefreshToken)
			authRoutes.GET("/oauth/google", authHandler.GoogleOAuth)
			authRoutes.GET("/oauth/google/callback", authHandler.GoogleOAuthCallback)
			authRoutes.GET("/oauth/facebook", authHandler.FacebookOAuth)
			authRoutes.GET("/oauth/facebook/callback", authHandler.FacebookOAuthCallback)
			authRoutes.POST("/reset-password/request", authHandler.ResetPasswordRequestOTP)
			authRoutes.POST("/reset-password/verify", authHandler.VerifyOTPForgotPassword)
			authRoutes.POST("/reset-password", authHandler.ResetPassword)
		}

		authenticated := v1.Group("")
		authenticated.Use(gin.HandlerFunc(authMiddleware))
		{
			// User routes
			users := authenticated.Group("/users")
			{
				users.GET("/me", authHandler.GetMe)
				users.PUT("/me", authHandler.UpdateMe)
				users.GET("/me/weight-history", userHandler.GetMyWeightHistory)
				users.GET("/me/meal-streak", userHandler.GetMealStreak)
				users.POST("/me/push-token", userHandler.RegisterPushToken)
				users.DELETE("/me/push-token", userHandler.DeletePushToken)
				users.PUT("/me/avatar", userHandler.UploadAvatarImage)
				users.GET("/:id", placeholderHandler("Get user by ID"))
				users.GET("/nutrition-target", userHandler.GetUserNutritionTarget)
				users.PUT("/nutrition-target", userHandler.UpdateUserNutritionTarget)
				// Backward compatibility (deprecated): use PUT /users/me/avatar
				users.PUT("/avatar", userHandler.UploadAvatarImage)
				users.DELETE("/me", userHandler.DeleteAccount)
			}

			// Exercise routes
			exercises := authenticated.Group("/exercises")
			{
				exercises.GET("", exerciseHandler.ListExercises)
				exercises.GET("/:id/stats", exerciseHandler.GetExerciseStats)
				exercises.GET("/:id", exerciseHandler.GetExercise)
			}

			// Workout Plan routes
			workoutPlans := authenticated.Group("/workout-plans")
			{
				workoutPlans.POST("", workoutHandler.CreateWorkoutPlan)
				workoutPlans.GET("", workoutHandler.ListWorkoutPlans)
				workoutPlans.GET("/:id", workoutHandler.GetWorkoutPlan)
				workoutPlans.PUT("/:id", workoutHandler.UpdateWorkoutPlan)
				workoutPlans.DELETE("/:id", workoutHandler.DeleteWorkoutPlan)

				// Exercise management
				workoutPlans.POST("/:id/exercises", placeholderHandler("Update exercise in plan"))
				workoutPlans.PUT("/:id/exercises/:exerciseId", placeholderHandler("Update exercise in plan"))
				workoutPlans.DELETE("/:id/exercises/:exerciseId", placeholderHandler("Remove exercise from plan"))
			}

			// Workout Session routes (calendar & tracking)
			workoutSessions := authenticated.Group("/workout-sessions")
			{
				workoutSessions.GET("/scheduled-dates", workoutHandler.GetScheduledDates)
				workoutSessions.GET("/weekly-summary", workoutHandler.GetWeeklyWorkoutSummary)
				workoutSessions.GET("", workoutHandler.GetSessionsByDate)
				workoutSessions.GET("/:id", workoutHandler.GetSessionByID)
				workoutSessions.DELETE("/:id", workoutHandler.DeleteWorkoutSession)
				workoutSessions.POST("", workoutHandler.CreateWorkoutSession)
				workoutSessions.PATCH("/:id", workoutHandler.UpdateWorkoutSession)
				workoutSessions.PATCH("/:id/exercise-sets/:setId", workoutHandler.UpdateSessionSet)
				workoutSessions.PATCH("/:id/finish", workoutHandler.FinishWorkoutSession)
			}

			// Food routes
			foods := authenticated.Group("/foods")
			{
				foods.GET("", foodHandler.ListFoods)
				foods.GET("/:id", foodHandler.GetFood)
				foods.GET("/search", foodHandler.SearchFoods)
				foods.POST("", foodHandler.CreateFood)
				foods.PUT("/:id", foodHandler.UpdateFood)
				foods.DELETE("/:id", foodHandler.DeleteFood)
				foods.POST("/scan", foodHandler.ScanFood)
				foods.POST("/sync-vectors", foodHandler.SyncVectors)
			}

			// Recipe routes
			recipes := authenticated.Group("/recipes")
			{
				recipes.POST("", recipeHandler.CreateRecipe)
				recipes.GET("", recipeHandler.ListRecipes)
				recipes.GET("/:id", recipeHandler.GetRecipe)
				recipes.PUT("/:id", recipeHandler.UpdateRecipe)
				recipes.DELETE("/:id", recipeHandler.DeleteRecipe)

				// Food management in recipes (Foods are managed during recipe Create/Update as requested)
				recipes.POST("/:id/foods", placeholderHandler("Add food to recipe"))
				recipes.PUT("/:id/foods/:foodId", placeholderHandler("Update food in recipe"))
				recipes.DELETE("/:id/foods/:foodId", placeholderHandler("Remove food from recipe"))
			}

			// Meal Log routes
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

			// Meal Daily routes
			mealDaily := authenticated.Group("/meal-daily")
			{
				mealDaily.GET("/date/:date", mealDailyHandler.GetMealDailyTargetByDate)
			}

			// Social routes
			social := authenticated.Group("/social")
			{
				social.GET("/feed", socialHandler.GetFeed)
				social.POST("/users/:userId/follow", socialHandler.FollowUser)
				social.DELETE("/users/:userId/follow", socialHandler.UnfollowUser)
				social.POST("/posts", socialHandler.CreatePost)
				social.GET("/posts/:postId", socialHandler.GetPostByID)
				social.POST("/posts/:postId/likes", socialHandler.LikePost)
				social.DELETE("/posts/:postId/likes", socialHandler.UnlikePost)
				social.GET("/posts/:postId/comments", socialHandler.GetPostComments)
				social.GET("/posts/:postId/comments/:commentId/replies", socialHandler.GetCommentReplies)
				social.POST("/posts/:postId/comments", socialHandler.CreateComment)
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

func placeholderHandler(description string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Endpoint not yet implemented: " + description,
		})
	}
}
