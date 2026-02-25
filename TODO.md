# Implementation TODO List

This file tracks all remaining TODO items in the codebase that need implementation.

## 🔴 High Priority (Required for Basic Functionality)

### User Repository
- ✅ `Create()` - Complete
- ✅ `GetByID()` - Complete
- ✅ `GetByEmail()` - Complete
- ✅ `GetByOAuth()` - Complete
- ✅ `Update()` - Complete
- ✅ `Delete()` - Complete
- ✅ `Exists()` - Complete
- ⏳ `UpdateProfile()` - Needs dynamic query for selective updates
- ⏳ `UpdatePassword()` - Complete but can be enhanced

### Workout Repositories

#### Exercise Repository (9 methods)
- ⏳ `Create()` - Insert exercise with JSONB fields
- ⏳ `GetByID()` - Query with JSONB parsing
- ⏳ `List()` - Pagination
- ⏳ `Search()` - Dynamic filters (category, muscle group, equipment, query)
- ⏳ `Update()` - Update exercise
- ⏳ `Delete()` - Soft delete (is_active = false)

#### WorkoutPlan Repository (9 methods)
- ⏳ `Create()` - Insert workout plan
- ⏳ `GetByID()` - Query with JOIN exercises
- ⏳ `GetByUserID()` - List user's plans with pagination
- ⏳ `Update()` - Update plan
- ⏳ `Delete()` - Delete plan
- ⏳ `AddExercise()` - Add exercise to plan
- ⏳ `UpdateExercise()` - Update exercise config
- ⏳ `RemoveExercise()` - Remove exercise from plan
- ⏳ `GetExercises()` - Query plan exercises

#### WorkoutSchedule Repository (8 methods)
- ⏳ `Create()` - Insert schedule
- ⏳ `GetByID()` - Query schedule
- ⏳ `GetByUserID()` - List with filters
- ⏳ `GetByDateRange()` - Query by date range
- ⏳ `Update()` - Update schedule
- ⏳ `Delete()` - Delete schedule
- ⏳ `MarkCompleted()` - Mark as completed
- ⏳ `BulkCreate()` - Batch insert schedules

#### WorkoutSession Repository (9 methods)
- ⏳ `Create()` - Insert session
- ⏳ `GetByID()` - Query with exercises
- ⏳ `GetByUserID()` - List with pagination
- ⏳ `Update()` - Update session
- ⏳ `Complete()` - Mark completed with stats
- ⏳ `AddExerciseLog()` - Log exercise with JSONB sets
- ⏳ `GetExercises()` - Query session exercises
- ⏳ `GetStats()` - Calculate workout statistics
- ⏳ `Delete()` - Delete session

### Meal Repositories

#### Food Repository (7 methods)
- ⏳ `Create()` - Insert food
- ⏳ `GetByID()` - Query food
- ⏳ `List()` - List with pagination
- ⏳ `Search()` - Search with filters
- ⏳ `Update()` - Update food
- ⏳ `Delete()` - Delete user food only
- ⏳ `GetByBarcode()` - Query by barcode

#### Recipe Repository (10 methods)
- ⏳ `Create()` - Insert recipe
- ⏳ `GetByID()` - Query with foods
- ⏳ `GetByUserID()` - List user's recipes
- ⏳ `Update()` - Update recipe
- ⏳ `Delete()` - Delete recipe
- ⏳ `AddFood()` - Add food to recipe
- ⏳ `UpdateFood()` - Update food quantity
- ⏳ `RemoveFood()` - Remove food
- ⏳ `GetFoods()` - Query recipe foods
- ⏳ `RecalculateNutrition()` - Sum and update nutrition

#### MealLog Repository (12 methods)
- ⏳ `Create()` - Insert meal log
- ⏳ `GetByID()` - Query with items
- ⏳ `GetByUserID()` - List with filters
- ⏳ `GetByDate()` - Query by date
- ⏳ `Update()` - Update log
- ⏳ `Delete()` - Delete log
- ⏳ `AddItem()` - Add food/recipe item
- ⏳ `UpdateItem()` - Update item quantity
- ⏳ `RemoveItem()` - Remove item
- ⏳ `GetItems()` - Query log items
- ⏳ `RecalculateNutrition()` - Sum nutrition
- ⏳ `GetDailySummary()` - Calculate daily summary
- ⏳ `GetStats()` - Calculate period statistics

### Social Repositories

#### Follow Repository (6 methods)
- ⏳ `Follow()` - Create follow relationship
- ⏳ `Unfollow()` - Delete relationship
- ⏳ `IsFollowing()` - Check if following
- ⏳ `GetFollowers()` - List followers
- ⏳ `GetFollowing()` - List following
- ⏳ `GetStats()` - Count followers/following

#### Post Repository (10 methods)
- ⏳ `Create()` - Insert post
- ⏳ `GetByID()` - Query post
- ⏳ `GetByUserID()` - List user's posts
- ⏳ `GetFeed()` - Activity feed (complex query)
- ⏳ `Update()` - Update post
- ⏳ `Delete()` - Delete post
- ⏳ `IncrementLikesCount()` - Increment counter
- ⏳ `DecrementLikesCount()` - Decrement counter
- ⏳ `IncrementCommentsCount()` - Increment counter
- ⏳ `DecrementCommentsCount()` - Decrement counter

#### Like Repository (4 methods)
- ⏳ `Create()` - Insert like
- ⏳ `Delete()` - Delete like
- ⏳ `Exists()` - Check if liked
- ⏳ `GetByPostID()` - List likes

#### Comment Repository (5 methods)
- ⏳ `Create()` - Insert comment
- ⏳ `GetByID()` - Query comment
- ⏳ `GetByPostID()` - List post comments
- ⏳ `Update()` - Update comment
- ⏳ `Delete()` - Delete comment

## 🟡 Medium Priority (Extend Functionality)

### Use Cases to Implement

#### Workout Use Cases (~15 use cases)
- ⏳ ListExercisesUseCase
- ⏳ SearchExercisesUseCase
- ⏳ GetExerciseUseCase
- ⏳ GetWorkoutPlanUseCase
- ⏳ UpdateWorkoutPlanUseCase
- ⏳ DeleteWorkoutPlanUseCase
- ⏳ ListWorkoutPlansUseCase
- ⏳ ScheduleWorkoutUseCase
- ⏳ BulkScheduleWorkoutUseCase
- ⏳ GetSchedulesUseCase
- ⏳ StartWorkoutSessionUseCase
- ⏳ LogExerciseSetUseCase
- ⏳ CompleteWorkoutSessionUseCase
- ⏳ GetSessionHistoryUseCase
- ⏳ GetWorkoutStatsUseCase

#### Meal Use Cases (~12 use cases)
- ⏳ CreateFoodUseCase
- ⏳ ListFoodsUseCase
- ⏳ SearchFoodsUseCase
- ⏳ UpdateFoodUseCase
- ⏳ CreateRecipeUseCase
- ⏳ AddFoodToRecipeUseCase
- ⏳ RecipeNutritionCalculatorUseCase
- ⏳ CreateMealLogUseCase
- ⏳ AddItemToMealLogUseCase
- ⏳ GetMealHistoryUseCase
- ⏳ GetDailySummaryUseCase
- ⏳ GetNutritionStatsUseCase

#### Social Use Cases (~8 use cases)
- ⏳ FollowUserUseCase
- ⏳ UnfollowUserUseCase
- ⏳ CreatePostUseCase
- ⏳ LikePostUseCase
- ⏳ UnlikePostUseCase
- ⏳ CreateCommentUseCase
- ⏳ DeleteCommentUseCase
- ⏳ GetFeedUseCase

### HTTP Handlers

#### Workout Handlers (~20 handlers)
- ✅ `CreateWorkoutPlan()` - Wired
- ⏳ `ListWorkoutPlans()` - Placeholder
- ⏳ `GetWorkoutPlan()` - Placeholder
- ⏳ `UpdateWorkoutPlan()` - Placeholder
- ⏳ `DeleteWorkoutPlan()` - Placeholder
- ✅ `AddExerciseToWorkout()` - Wired
- ⏳ `UpdateExerciseInWorkout()` - Placeholder
- ⏳ `RemoveExerciseFromWorkout()` - Placeholder
- ⏳ Exercise handlers (list, get, search)
- ⏳ Schedule handlers (create, bulk, list, calendar, update, delete)
- ⏳ Session handlers (start, log set, complete, history, stats)

#### Meal Handlers (~15 handlers)
- ⏳ All food handlers
- ⏳ All recipe handlers
- ⏳ All meal log handlers
- ⏳ Statistics handlers

#### Social Handlers (~10 handlers)
- ⏳ Follow/unfollow handlers
- ⏳ Post handlers
- ⏳ Like handlers
- ⏳ Comment handlers
- ⏳ Feed handler

## 🟢 Low Priority (Nice to Have)

### OAuth2 Implementation
- ⏳ Implement `OAuthLoginUseCase`
- ⏳ Complete Google OAuth flow in handler
- ⏳ Complete Facebook OAuth flow in handler
- ⏳ Add OAuth providers to fx module
- ⏳ Test OAuth callbacks

### Swagger Documentation
- ⏳ Add annotations to all remaining handlers
- ⏳ Run `make swagger-gen` to generate docs
- ⏳ Test Swagger UI
- ⏳ Document all response models

### Testing
- ⏳ Unit tests for all use cases
- ⏳ Integration tests for repositories
- ⏳ HTTP handler tests
- ⏳ E2E tests for main flows

### Additional Features
- ⏳ File upload for images (profile, meal photos)
- ⏳ Email verification
- ⏳ Password reset flow
- ⏳ Admin panel endpoints
- ⏳ Analytics dashboard
- ⏳ Notification system

## 📌 How to Use This File

1. **Pick a section** to work on (start with repositories)
2. **Find TODO comments** in the code files
3. **Implement the functionality** following patterns
4. **Mark as complete** by changing ⏳ to ✅
5. **Test your implementation**
6. **Move to next item**

## 🎯 Estimated Effort

- **Repositories**: 2-3 days (42 methods)
- **Use Cases**: 3-4 days (35 use cases)
- **Handlers**: 2-3 days (45 handlers)
- **OAuth2**: 1 day
- **Testing**: 2-3 days
- **Total**: 10-14 days for complete implementation

## ✅ Progress Tracking

Track your progress:
- Total TODO items: ~120
- Repositories: 0/42 (0%)
- Use Cases: 3/38 (8%) - Register, Login, Profile complete
- Handlers: 4/49 (8%) - Register, Login, GetMe, UpdateMe complete
- OAuth2: 0/4 (0%)
- Tests: 0/30 (0%)

**Overall Progress: 7/163 (4%)**

Keep going! The foundation is solid! 💪
