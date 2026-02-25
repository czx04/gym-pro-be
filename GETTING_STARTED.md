# Getting Started Guide

## 🎉 Project Setup Complete!

Your Gym Pro backend base project has been successfully created with a solid foundation following Clean Architecture principles.

## ✅ What's Implemented

### Core Infrastructure (100%)
- ✅ **Project Structure**: Clean Architecture layout with proper separation of concerns
- ✅ **Configuration**: Viper-based config with environment variables support
- ✅ **Logger**: Uber Zap structured logging with fx integration
- ✅ **Database**: PostgreSQL connection pooling with pgx driver
- ✅ **Migrations**: Complete database schema (users, workouts, meals, social)
- ✅ **Error Handling**: Standardized error types and response formatting
- ✅ **Validation**: go-playground/validator integration

### Authentication & Security (95%)
- ✅ **JWT**: Token generation, validation, and refresh logic
- ✅ **Password Hashing**: bcrypt implementation
- ✅ **Auth Middleware**: JWT authentication for protected routes
- ⏳ **OAuth2**: Structure ready, needs Google & Facebook integration

### HTTP Layer (100%)
- ✅ **Router**: Gin router with all API endpoint definitions
- ✅ **Middleware**: Auth, logging, CORS, rate limiting, error handling, recovery
- ✅ **Response Utilities**: Standardized JSON responses and pagination
- ✅ **Health Checks**: `/health` and `/ping` endpoints

### Domain Models (100%)
- ✅ **User Domain**: Complete models and repository interfaces
- ✅ **Workout Domain**: Exercise, WorkoutPlan, Schedule, Session models
- ✅ **Meal Domain**: Food, Recipe, MealLog models
- ✅ **Social Domain**: Follow, Post, Like, Comment models

### Development Tools (100%)
- ✅ **Makefile**: Comprehensive commands for dev workflow
- ✅ **Docker**: docker-compose.yml with PostgreSQL and Redis
- ✅ **Dockerfile**: Multi-stage build for production
- ✅ **.gitignore**: Proper Git ignore rules
- ✅ **Documentation**: README, ARCHITECTURE, and this guide

## 🚀 Quick Start (5 minutes)

### 1. Install Go Tools

```bash
# Install development tools
make dev-setup
```

This installs:
- `golang-migrate` for database migrations
- `swag` for Swagger documentation generation

### 2. Start Database

```bash
# Start PostgreSQL and Redis
make docker-up

# Wait for database to be ready (5 seconds)
```

### 3. Run Migrations

```bash
# Create database tables
make migrate-up
```

Expected output:
```
Applying migrations...
✓ 000001_create_users_table.up.sql
✓ 000002_create_workout_tables.up.sql
✓ 000003_create_meal_tables.up.sql
✓ 000004_create_social_tables.up.sql
Done!
```

### 4. Start the API Server

```bash
# Run the application
make run
```

Expected output:
```
🚀 Gym Pro API Server Starting...
[INFO] Starting HTTP server (host: 0.0.0.0, port: 8080)
[INFO] HTTP server started successfully (address: 0.0.0.0:8080)
```

### 5. Test the Server

Open your browser or use curl:

```bash
# Health check
curl http://localhost:8080/health

# Response:
# {"status":"ok","message":"Service is healthy"}

# Ping
curl http://localhost:8080/ping

# Response:
# {"message":"pong"}
```

## 📝 Next Steps - What to Implement

### Phase 1: Complete Repositories (Priority: HIGH)

Implement PostgreSQL repositories for all domains:

**Location**: `internal/repository/postgres/`

#### User Repository
```go
// File: internal/repository/postgres/user_repository.go
type userRepository struct {
    db *database.DB
}

func (r *userRepository) Create(ctx context.Context, user *user.User) error {
    query := `INSERT INTO users (id, email, password_hash, name, ...) VALUES ($1, $2, $3, $4, ...)`
    _, err := r.db.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.Name, ...)
    return err
}

// Implement all methods from user.Repository interface
```

**Tasks**:
- [ ] Create `user_repository.go`
- [ ] Create `workout_repository.go` (Exercise, WorkoutPlan, Schedule, Session)
- [ ] Create `meal_repository.go` (Food, Recipe, MealLog)
- [ ] Create `social_repository.go` (Follow, Post, Like, Comment)
- [ ] Create repository module for fx: `internal/repository/module.go`

### Phase 2: Implement Use Cases (Priority: HIGH)

Create business logic layer:

**Location**: `internal/usecase/`

#### Example: User Use Case
```go
// File: internal/usecase/user/register.go
type RegisterUseCase struct {
    userRepo user.Repository
    passwordMgr *auth.PasswordManager
    jwtMgr *auth.JWTManager
    validator *validator.Validator
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input user.CreateUserInput) (*TokenPair, error) {
    // 1. Validate input
    if err := uc.validator.Validate(input); err != nil {
        return nil, errors.Validation(err.Error())
    }
    
    // 2. Check if user exists
    exists, _ := uc.userRepo.Exists(ctx, input.Email)
    if exists {
        return nil, errors.Conflict("email already registered")
    }
    
    // 3. Hash password
    hash, err := uc.passwordMgr.HashPassword(input.Password)
    if err != nil {
        return nil, errors.InternalServer("failed to hash password", err)
    }
    
    // 4. Create user
    newUser := &user.User{
        ID: uuid.New(),
        Email: input.Email,
        PasswordHash: hash,
        Name: input.Name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := uc.userRepo.Create(ctx, newUser); err != nil {
        return nil, errors.DatabaseError("create user", err)
    }
    
    // 5. Generate tokens
    accessToken, refreshToken, err := uc.jwtMgr.GenerateTokenPair(newUser.ID, newUser.Email)
    if err != nil {
        return nil, errors.InternalServer("failed to generate tokens", err)
    }
    
    return &TokenPair{
        AccessToken: accessToken,
        RefreshToken: refreshToken,
        User: newUser,
    }, nil
}
```

**Tasks**:
- [ ] User: Register, Login, UpdateProfile, GetProfile
- [ ] Workout: CreatePlan, ScheduleWorkout, StartSession, LogExercise
- [ ] Meal: CreateFood, CreateRecipe, LogMeal, GetNutritionStats
- [ ] Social: FollowUser, CreatePost, LikePost, CommentPost
- [ ] Create usecase modules for fx

### Phase 3: Implement HTTP Handlers (Priority: HIGH)

Create Gin HTTP handlers:

**Location**: `internal/delivery/http/handler/`

#### Example: Auth Handler
```go
// File: internal/delivery/http/handler/auth_handler.go
type AuthHandler struct {
    registerUC *usecase.RegisterUseCase
    loginUC *usecase.LoginUseCase
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.CreateUserInput true "Registration request"
// @Success 201 {object} response.Response{data=TokenResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    var input user.CreateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.Error(c, errors.BadRequest("invalid request body"))
        return
    }
    
    result, err := h.registerUC.Execute(c.Request.Context(), input)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    response.Created(c, result, "User registered successfully")
}
```

**Tasks**:
- [ ] `auth_handler.go` - Register, Login, Refresh, OAuth
- [ ] `user_handler.go` - Profile management
- [ ] `workout_handler.go` - Workout management
- [ ] `meal_handler.go` - Meal tracking
- [ ] `social_handler.go` - Social features
- [ ] Wire handlers in router (replace placeholder handlers)

### Phase 4: OAuth2 Integration (Priority: MEDIUM)

**Location**: `internal/infrastructure/auth/oauth.go`

```go
type OAuthProvider interface {
    GetAuthURL(state string) string
    GetUserInfo(code string) (*OAuthUserInfo, error)
}

type GoogleOAuth struct {
    config *oauth2.Config
}

type FacebookOAuth struct {
    config *oauth2.Config
}
```

**Tasks**:
- [ ] Implement Google OAuth provider
- [ ] Implement Facebook OAuth provider
- [ ] Create OAuth handlers
- [ ] Test OAuth flow

### Phase 5: Swagger Documentation (Priority: MEDIUM)

**Tasks**:
- [ ] Add Swagger annotations to all handlers
- [ ] Run `make swagger-gen` to generate docs
- [ ] Test Swagger UI at `/swagger/index.html`
- [ ] Document all request/response models

### Phase 6: Testing (Priority: MEDIUM)

**Tasks**:
- [ ] Write unit tests for use cases
- [ ] Write integration tests for repositories
- [ ] Write HTTP handler tests
- [ ] Aim for >70% code coverage
- [ ] Run `make test-coverage`

## 📚 Development Workflow

### Daily Development

```bash
# 1. Start database
make docker-up

# 2. Run migrations (if new migrations)
make migrate-up

# 3. Start development server
make run

# 4. In another terminal, watch logs
make docker-logs
```

### Making Changes

```bash
# 1. Create a new feature branch
git checkout -b feature/your-feature

# 2. Make your changes

# 3. Format code
make fmt

# 4. Run tests
make test

# 5. Run linter (if golangci-lint installed)
make lint

# 6. Commit changes
git add .
git commit -m "Add your feature"
```

### Database Changes

```bash
# 1. Create new migration
make migrate-create name=add_new_table

# 2. Edit migration files in migrations/
# - Edit .up.sql to add changes
# - Edit .down.sql to revert changes

# 3. Apply migration
make migrate-up

# 4. Test rollback
make migrate-down
make migrate-up
```

### Adding New Dependencies

```bash
# 1. Add dependency
go get github.com/some/package

# 2. Update go.mod and go.sum
make deps

# 3. Vendor (optional)
go mod vendor
```

## 🐛 Troubleshooting

### Database Connection Failed

```bash
# Check if PostgreSQL is running
make docker-logs

# Restart PostgreSQL
make docker-down
make docker-up

# Check connection
make db-psql
```

### Migration Errors

```bash
# Check current migration version
migrate -path migrations -database "postgresql://..." version

# Force to specific version if stuck
make migrate-force version=4

# Reset database completely
make db-reset
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in .env
SERVER_PORT=8081
```

## 📖 Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Uber-go/fx](https://uber-go.github.io/fx/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 🤝 Getting Help

1. Check `README.md` for overview
2. Read `ARCHITECTURE.md` for design details
3. Review example implementations in the codebase
4. Ask the development team

## 🎯 Success Criteria

Your implementation is complete when:

- [ ] All TODO endpoints return real data (not placeholder responses)
- [ ] All tests pass with >70% coverage
- [ ] Swagger documentation is complete
- [ ] OAuth2 login works for Google and Facebook
- [ ] Can create, schedule, and track workouts
- [ ] Can log meals and view nutrition statistics
- [ ] Social features (follow, post, like, comment) work
- [ ] Application runs successfully in Docker

## 🚀 You're Ready!

The foundation is solid. Start with implementing repositories, then use cases, then handlers. Follow the examples provided and refer to the domain models for the data structures.

Good luck! 💪
