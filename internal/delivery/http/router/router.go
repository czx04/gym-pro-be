package router

import (
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/delivery/http/handler"
	"gym-pro-2026-ptit/internal/delivery/http/middleware"

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
	authHandler *handler.AuthHandler,
	workoutHandler *handler.WorkoutHandler,
	exerciseHandler *handler.ExerciseHandler,
	foodHandler *handler.FoodHandler,
	recipeHandler *handler.RecipeHandler,
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
				users.GET("/:id", placeholderHandler("Get user by ID"))
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

			// Workout Schedule routes
			workoutSchedules := authenticated.Group("/workout-schedules")
			{
				workoutSchedules.POST("", placeholderHandler("Schedule workout"))
				workoutSchedules.POST("/bulk", placeholderHandler("Bulk schedule workouts"))
				workoutSchedules.GET("", placeholderHandler("List schedules"))
				workoutSchedules.GET("/calendar/:year/:month", placeholderHandler("Calendar view"))
				workoutSchedules.PUT("/:id", placeholderHandler("Update schedule"))
				workoutSchedules.DELETE("/:id", placeholderHandler("Delete schedule"))
			}

			// Workout Session routes (calendar & tracking)
			workoutSessions := authenticated.Group("/workout-sessions")
			{
				workoutSessions.GET("/scheduled-dates", workoutHandler.GetScheduledDates)
				workoutSessions.GET("", workoutHandler.GetSessionsByDate)
				workoutSessions.GET("/:id", workoutHandler.GetSessionByID)
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
				mealLogs.POST("", placeholderHandler("Create meal log"))
				mealLogs.GET("", placeholderHandler("Get meal log history"))
				mealLogs.GET("/date/:date", placeholderHandler("Get logs by date"))
				mealLogs.GET("/:id", placeholderHandler("Get meal log"))
				mealLogs.PUT("/:id", placeholderHandler("Update meal log"))
				mealLogs.DELETE("/:id", placeholderHandler("Delete meal log"))

				// Item management
				mealLogs.POST("/:id/items", placeholderHandler("Add item to meal log"))
				mealLogs.PUT("/:id/items/:itemId", placeholderHandler("Update item"))
				mealLogs.DELETE("/:id/items/:itemId", placeholderHandler("Remove item"))

				// Statistics
				mealLogs.GET("/stats/daily", placeholderHandler("Daily nutrition stats"))
				mealLogs.GET("/stats/weekly", placeholderHandler("Weekly nutrition stats"))
				mealLogs.GET("/stats/monthly", placeholderHandler("Monthly nutrition stats"))
			}

			// Social routes
			social := authenticated.Group("/social")
			{
				// Follow management
				social.POST("/follow/:userId", placeholderHandler("Follow user"))
				social.DELETE("/follow/:userId", placeholderHandler("Unfollow user"))
				social.GET("/followers", placeholderHandler("Get followers"))
				social.GET("/following", placeholderHandler("Get following list"))

				// Post management
				social.POST("/posts", placeholderHandler("Create post"))
				social.GET("/posts", placeholderHandler("Get user posts"))
				social.GET("/feed", placeholderHandler("Get activity feed"))
				social.DELETE("/posts/:id", placeholderHandler("Delete post"))

				// Likes
				social.POST("/posts/:id/like", placeholderHandler("Like post"))
				social.DELETE("/posts/:id/like", placeholderHandler("Unlike post"))

				// Comments
				social.POST("/posts/:id/comments", placeholderHandler("Add comment"))
				social.GET("/posts/:id/comments", placeholderHandler("Get comments"))
				social.PUT("/comments/:id", placeholderHandler("Update comment"))
				social.DELETE("/comments/:id", placeholderHandler("Delete comment"))
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
