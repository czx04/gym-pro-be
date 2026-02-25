# Implementation Guide

## 🎉 Congratulations!

Your Gym Pro backend project now has **complete base/skeleton code** ready for implementation!

## ✅ What's Been Created (20/21 Complete!)

### 1. **Repository Layer** ✅
- `internal/repository/postgres/user_repository.go` - User CRUD (implemented examples)
- `internal/repository/postgres/workout_repository.go` - Workout repositories with TODO markers
- `internal/repository/postgres/meal_repository.go` - Meal repositories with TODO markers
- `internal/repository/postgres/social_repository.go` - Social repositories with TODO markers
- `internal/repository/module.go` - fx module wiring all repositories

### 2. **Use Case Layer** ✅
- `internal/usecase/user/register.go` - Registration logic (complete example)
- `internal/usecase/user/login.go` - Login logic (complete example)
- `internal/usecase/user/profile.go` - Profile management (complete example)
- `internal/usecase/workout/workout_plan.go` - Workout use cases (examples with TODOs)
- `internal/usecase/module.go` - fx module with TODO markers for more use cases

### 3. **Handler Layer** ✅
- `internal/delivery/http/handler/auth_handler.go` - Auth endpoints with Swagger annotations
- `internal/delivery/http/handler/workout_handler.go` - Workout endpoints with Swagger
- `internal/delivery/http/handler/module.go` - Handler module
- **Router Updated** - Handlers wired into router (register, login, profile working!)

### 4. **OAuth2 Base** ✅ (template ready)
- `internal/infrastructure/auth/oauth.go` - Google & Facebook OAuth providers
- Complete structure with TODO markers for implementation

### 5. **Swagger Documentation** ✅
- `docs/docs.go` - Swagger initialization with all tags
- All handlers have Swagger annotations
- Ready to generate with `make swagger-gen`

### 6. **Dependency Injection** ✅
- `cmd/api/main.go` - All modules wired together
- Complete fx dependency chain:
  ```
  Config → Logger → Database → Repositories → Use Cases → Handlers → Router
  ```

## 📝 What You Need to Implement

### Priority 1: Complete Repositories (Critical)

Each repository file has methods with `// TODO:` comments. Fill in the SQL queries and logic.

#### Example Pattern:
```go
func (r *exerciseRepository) GetByID(ctx context.Context, id uuid.UUID) (*workout.Exercise, error) {
    // TODO: Query exercise by ID, parse JSONB fields
    
    // Your implementation:
    query := `
        SELECT id, name, description, category, muscle_groups, 
               equipment_needed, difficulty_level, calories_per_minute,
               video_url, thumbnail_url, is_active, created_by,
               created_at, updated_at
        FROM exercises
        WHERE id = $1 AND is_active = true
    `
    
    var e workout.Exercise
    var muscleGroupsJSON, equipmentJSON []byte
    
    err := r.db.QueryRow(ctx, query, id).Scan(
        &e.ID, &e.Name, &e.Description, &e.Category,
        &muscleGroupsJSON, &equipmentJSON,
        &e.DifficultyLevel, &e.CaloriesPerMinute,
        &e.VideoURL, &e.ThumbnailURL, &e.IsActive,
        &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
    )
    
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, errors.NotFound("exercise")
        }
        return nil, errors.DatabaseError("get exercise", err)
    }
    
    // Parse JSONB fields
    json.Unmarshal(muscleGroupsJSON, &e.MuscleGroups)
    json.Unmarshal(equipmentJSON, &e.EquipmentNeeded)
    
    return &e, nil
}
```

**Files to Complete**:
- `workout_repository.go` - ~15 methods
- `meal_repository.go` - ~15 methods
- `social_repository.go` - ~12 methods

### Priority 2: Complete Use Cases (High)

Use the patterns from `register.go` and `login.go` as templates.

#### Use Case Pattern:
```go
type SomeUseCase struct {
    repo      domain.Repository
    validator *validator.Validator
    // ... other dependencies
}

func (uc *SomeUseCase) Execute(ctx context.Context, input Input) (*Output, error) {
    // 1. Validate input
    if err := uc.validator.Validate(input); err != nil {
        return nil, errors.Validation(err.Error())
    }
    
    // 2. Check business rules
    
    // 3. Call repository
    
    // 4. Return result
}
```

**Use Cases Needed**:

**Workout Domain**:
- ListExercisesUseCase
- SearchExercisesUseCase
- GetWorkoutPlanUseCase
- UpdateWorkoutPlanUseCase
- DeleteWorkoutPlanUseCase
- ScheduleWorkoutUseCase
- BulkScheduleWorkoutUseCase
- StartWorkoutSessionUseCase
- LogExerciseSetUseCase
- CompleteWorkoutSessionUseCase
- GetWorkoutStatsUseCase

**Meal Domain**:
- CreateFoodUseCase
- CreateRecipeUseCase
- AddFoodToRecipeUseCase
- CreateMealLogUseCase
- AddItemToMealLogUseCase
- GetDailySummaryUseCase
- GetNutritionStatsUseCase

**Social Domain**:
- FollowUserUseCase
- UnfollowUserUseCase
- CreatePostUseCase
- LikePostUseCase
- UnlikePostUseCase
- CreateCommentUseCase
- GetFeedUseCase

### Priority 3: Complete Handlers (Medium)

Most handlers are placeholder. Implement them following the pattern:

```go
func (h *Handler) SomeEndpoint(c *gin.Context) {
    // 1. Get user ID if authenticated
    userID, err := middleware.GetUserID(c)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    // 2. Parse input
    var input domain.SomeInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.Error(c, errors.BadRequest("invalid request body"))
        return
    }
    
    // 3. Call use case
    result, err := h.someUC.Execute(c.Request.Context(), userID, input)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    // 4. Return response
    response.Success(c, result, "Success message")
}
```

**Handlers to Complete**:
- Workout handlers (exercises, schedules, sessions)
- Meal handlers (foods, recipes, meal logs)
- Social handlers (follow, posts, likes, comments)

### Priority 4: OAuth2 Integration (Low)

The OAuth structure is ready in `auth/oauth.go`. You need to:

1. **Create OAuth Use Case**:
```go
type OAuthLoginUseCase struct {
    userRepo user.Repository
    jwtMgr *auth.JWTManager
    googleProvider *auth.GoogleOAuthProvider
    facebookProvider *auth.FacebookOAuthProvider
}

func (uc *OAuthLoginUseCase) ExecuteGoogle(ctx context.Context, code string) (*TokenPair, error) {
    // 1. Get user info from Google
    oauthUser, err := uc.googleProvider.GetUserInfo(ctx, code)
    
    // 2. Check if user exists by oauth_provider + oauth_id
    existingUser, err := uc.userRepo.GetByOAuth(ctx, "google", oauthUser.ID)
    
    if err != nil {
        // 3. User doesn't exist - check by email
        existingUser, err = uc.userRepo.GetByEmail(ctx, oauthUser.Email)
        
        if err != nil {
            // 4. Create new user
            newUser := &user.User{
                ID: uuid.New(),
                Email: oauthUser.Email,
                Name: oauthUser.Name,
                OAuthProvider: &oauthUser.Provider,
                OAuthID: &oauthUser.ID,
                AvatarURL: oauthUser.AvatarURL,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
            }
            uc.userRepo.Create(ctx, newUser)
            existingUser = newUser
        } else {
            // 5. Link OAuth to existing user
            existingUser.OAuthProvider = &oauthUser.Provider
            existingUser.OAuthID = &oauthUser.ID
            uc.userRepo.Update(ctx, existingUser)
        }
    }
    
    // 6. Generate tokens
    accessToken, refreshToken, _ := uc.jwtMgr.GenerateTokenPair(existingUser.ID, existingUser.Email)
    
    return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken, User: existingUser}, nil
}
```

2. **Update Auth Handler** to call OAuth use case

3. **Test OAuth Flow**:
   - Visit `/api/v1/auth/oauth/google`
   - Should redirect to Google
   - After login, redirect back to callback
   - Callback should return tokens

## 🚀 Quick Start Implementation

### Day 1: Get Basic Flow Working

1. **Implement User Repository** (if not complete)
   - GetByID, GetByEmail, Create methods

2. **Test Registration & Login**:
```bash
# Start services
make docker-up
make migrate-up
make run

# Test registration
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'

# Test login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Save the access_token from response

# Test profile
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Day 2-3: Workout Domain

1. Implement Exercise repository methods
2. Create workout plan repository methods  
3. Implement corresponding use cases
4. Complete workout handlers
5. Test the flow

### Day 4-5: Meal Domain

1. Implement Food & Recipe repositories
2. Implement Meal Log repository
3. Create use cases
4. Complete handlers
5. Test nutrition tracking

### Day 6-7: Social Features

1. Implement Follow repository
2. Implement Post, Like, Comment repositories
3. Create social use cases
4. Complete handlers
5. Test social interactions

### Day 8: OAuth & Polish

1. Implement OAuth use case
2. Test Google & Facebook login
3. Generate Swagger docs: `make swagger-gen`
4. Test all endpoints
5. Write tests

## 🧪 Testing Strategy

### 1. Unit Tests

Test use cases in isolation:

```go
func TestRegisterUseCase(t *testing.T) {
    // Mock dependencies
    mockRepo := &MockUserRepository{}
    mockPwd := &MockPasswordManager{}
    mockJWT := &MockJWTManager{}
    mockValidator := &MockValidator{}
    
    uc := NewRegisterUseCase(mockRepo, mockPwd, mockJWT, mockValidator)
    
    // Setup mocks
    mockRepo.On("Exists", mock.Anything, "test@example.com").Return(false, nil)
    mockPwd.On("HashPassword", "password").Return("hashed", nil)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    mockJWT.On("GenerateTokenPair", mock.Anything, mock.Anything).Return("access", "refresh", nil)
    
    // Execute
    result, err := uc.Execute(context.Background(), user.CreateUserInput{
        Email: "test@example.com",
        Password: "password",
        Name: "Test",
    })
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 2. Integration Tests

Test repositories with real database:

```go
func TestUserRepository_Create(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewUserRepository(db)
    
    // Test
    user := &user.User{
        ID: uuid.New(),
        Email: "test@example.com",
        // ...
    }
    
    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    
    // Verify
    retrieved, err := repo.GetByID(context.Background(), user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Email, retrieved.Email)
}
```

### 3. API Tests

Test endpoints:

```go
func TestRegisterEndpoint(t *testing.T) {
    router := setupTestRouter(t)
    
    w := httptest.NewRecorder()
    body := `{"email":"test@example.com","password":"password123","name":"Test"}`
    req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 201, w.Code)
}
```

## 📚 Resources & Patterns

### JSONB Handling

```go
// Inserting JSONB
muscleGroups := []string{"chest", "triceps"}
data, _ := json.Marshal(muscleGroups)
_, err := db.Exec(ctx, "INSERT INTO exercises (muscle_groups) VALUES ($1)", data)

// Reading JSONB
var data []byte
var muscleGroups []string
err := db.QueryRow(ctx, "SELECT muscle_groups FROM exercises WHERE id = $1", id).Scan(&data)
json.Unmarshal(data, &muscleGroups)
```

### Transactions

```go
tx, err := r.db.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx)

// Do multiple operations
tx.Exec(...)
tx.Exec(...)

return tx.Commit(ctx)
```

### Pagination

```go
offset := (page - 1) * pageSize
query := `SELECT ... FROM table LIMIT $1 OFFSET $2`
rows, _ := db.Query(ctx, query, pageSize, offset)

// Count total
var total int64
db.QueryRow(ctx, "SELECT COUNT(*) FROM table").Scan(&total)
```

## ✅ Checklist

- [ ] All repository methods implemented
- [ ] All use cases implemented
- [ ] All handlers implemented
- [ ] OAuth2 working
- [ ] Swagger docs generated (`make swagger-gen`)
- [ ] Basic tests written
- [ ] README updated with your changes
- [ ] Test all endpoints manually
- [ ] Check linter errors (`make lint`)
- [ ] Format code (`make fmt`)

## 🎯 Success Criteria

Your implementation is complete when:

1. ✅ User can register, login, and update profile
2. ✅ User can browse exercises and create workout plans
3. ✅ User can schedule workouts and track sessions
4. ✅ User can log meals and view nutrition stats
5. ✅ User can follow others, create posts, and interact socially
6. ✅ OAuth2 login works with Google/Facebook
7. ✅ Swagger UI shows all endpoints at `/swagger/index.html`
8. ✅ All tests pass: `make test`

## 🆘 Getting Help

1. Check existing implemented examples (register, login)
2. Review domain models for data structures
3. Check ARCHITECTURE.md for design patterns
4. Search for similar implementations in codebase
5. Read Go documentation for pgx, Gin, fx

## 🚀 You're All Set!

Everything is wired and ready. Just fill in the `// TODO:` sections following the patterns provided. Start with user repositories and work your way up!

Happy coding! 💪
