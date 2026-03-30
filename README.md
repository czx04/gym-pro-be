# Gym Pro API - Backend

A production-ready Go backend API for a fitness tracking mobile application, built with Clean Architecture principles.

## Features

- 🏋️ **Workout Management**: Exercise library, workout plans, scheduling, and session tracking
- 🍎 **Meal Tracking**: Food library, recipes, and daily meal logging with nutrition tracking
- 👥 **Social Features**: Follow users, share workouts/meals, likes, and comments
- 🔐 **Authentication**: JWT + OAuth2 (Google, Facebook)
- 📊 **Analytics**: Progress tracking, statistics, and calorie target adherence
- 🐳 **Docker Ready**: Containerized development and deployment
- 📝 **API Documentation**: Auto-generated Swagger/OpenAPI docs

## Tech Stack

- **Framework**: Gin Web Framework
- **Dependency Injection**: Uber-go/fx
- **Logger**: Uber Zap (structured logging)
- **Database**: PostgreSQL with pgx driver
- **Migrations**: golang-migrate
- **Authentication**: JWT + OAuth2
- **Validation**: go-playground/validator
- **Documentation**: Swaggo (Swagger/OpenAPI)
- **Containerization**: Docker & Docker Compose

## Architecture

The project follows **Clean Architecture** principles:

```
gym-pro-2026-ptit/
├── cmd/api/              # Application entry point
├── internal/             # Private application code
│   ├── config/          # Configuration management
│   ├── domain/          # Business entities & interfaces
│   │   ├── user/
│   │   ├── workout/
│   │   ├── meal/
│   │   └── social/
│   ├── usecase/         # Business logic (to be implemented)
│   ├── delivery/        # HTTP handlers & middleware
│   ├── repository/      # Data access layer (to be implemented)
│   └── infrastructure/  # External services (auth, logger, database)
├── pkg/                 # Public shared utilities
│   ├── errors/
│   ├── response/
│   └── validator/
├── migrations/          # Database migrations
├── docs/                # Swagger documentation
└── docker/              # Docker configuration
```

## Prerequisites

- Go 1.24+
- PostgreSQL 16+
- Docker & Docker Compose (optional)
- Make

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd gym-pro-2026-ptit
```

### 2. Install dependencies

```bash
make deps
make dev-setup
```

This will:
- Download Go dependencies
- Install golang-migrate tool
- Install swag tool

### 3. Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

Update `.env` with your configuration:

```env
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=gymadmin
DB_PASSWORD=secret123
DB_NAME=gym_pro_db

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this
JWT_ACCESS_TOKEN_EXPIRE=15m
JWT_REFRESH_TOKEN_EXPIRE=7d

# OAuth2 (Configure with your credentials)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
FACEBOOK_APP_ID=your-facebook-app-id
FACEBOOK_APP_SECRET=your-facebook-app-secret
```

### 4. Start PostgreSQL

```bash
make docker-up
```

This starts PostgreSQL and Redis containers.

### 5. Run database migrations

```bash
make migrate-up
```

### 6. Run the application

```bash
make run
```

The API server will start at `http://localhost:8080`

## API Documentation

Once the server is running, access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## Available Make Commands

### Development
- `make run` - Run the application locally
- `make build` - Build the application
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage report
- `make clean` - Clean build artifacts

### Code Quality
- `make fmt` - Format code
- `make lint` - Run linter

### Dependencies
- `make deps` - Download dependencies
- `make deps-upgrade` - Upgrade dependencies

### Docker
- `make docker-up` - Start Docker containers
- `make docker-down` - Stop Docker containers
- `make docker-logs` - Show Docker logs
- `make docker-build` - Build Docker image
- `make docker-rebuild` - Rebuild and restart containers

### Database
- `make migrate-up` - Run migrations
- `make migrate-down` - Rollback migrations
- `make migrate-create name=<name>` - Create new migration
- `make db-reset` - Reset database (drop, create, migrate)
- `make db-psql` - Connect to database with psql

### Swagger
- `make swagger-gen` - Generate Swagger documentation
- `make swagger-fmt` - Format Swagger comments

### Setup
- `make dev-setup` - Setup development environment
- `make dev-start` - Start all development services

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/oauth/google` - Google OAuth
- `GET /api/v1/auth/oauth/facebook` - Facebook OAuth

### User Management
- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update profile
- `GET /api/v1/users/:id` - Get user by ID

### Workouts
- **Exercises**: `GET /api/v1/exercises`
- **Workout Plans**: CRUD operations at `/api/v1/workout-plans`
- **Schedules**: `/api/v1/workout-schedules`
- **Sessions**: `/api/v1/workout-sessions` (tracking)

### Meals
- **Foods**: `/api/v1/foods` (library + custom foods)
- **Recipes**: `/api/v1/recipes`
- **Meal Logs**: `/api/v1/meal-logs` (daily tracking)
- **Statistics**: `/api/v1/meal-logs/stats/*`

### Social
- **Follow**: `/api/v1/social/follow/*`
- **Posts**: `/api/v1/social/posts`
- **Feed**: `/api/v1/social/feed`
- **Likes/Comments**: `/api/v1/social/posts/:id/like`, `/api/v1/social/posts/:id/comments`

### Health
- `GET /health` - Health check
- `GET /ping` - Ping endpoint

## Database Schema

The database includes the following main tables:

- **users** - User accounts and profiles
- **exercises** - Pre-populated exercise library
- **workout_plans** - User workout plans
- **workout_plan_exercises** - Exercises in plans
- **workout_schedules** - Scheduled workouts
- **workout_sessions** - Actual workout tracking
- **workout_session_exercises** - Per-exercise performance tracking
- **foods** - Food library (system + user custom)
- **recipes** - User recipes
- **recipe_foods** - Foods in recipes
- **meal_logs** - Daily meal consumption logs
- **meal_log_items** - Foods/recipes in meal logs
- **follows** - User follow relationships
- **posts** - Shared workouts/meals
- **likes** - Post likes
- **comments** - Post comments

## Project Status

### ✅ Completed
- Project structure and setup
- Configuration management
- Logger integration
- Database connection
- Database migrations
- Domain models and interfaces
- JWT authentication
- Middleware (auth, logging, CORS, rate limiting, error handling)
- Router with all endpoint definitions
- Docker setup
- Makefile with development commands

### 🚧 To Be Implemented
- PostgreSQL repositories
- OAuth2 integration (Google, Facebook)
- Use cases (business logic)
- HTTP handlers for all endpoints
- Swagger documentation generation

## Development Guidelines

1. **Clean Architecture**: Follow separation of concerns
2. **Error Handling**: Use custom error types from `pkg/errors`
3. **Logging**: Use structured logging with Zap
4. **Validation**: Validate all inputs using validator tags
5. **Testing**: Write unit and integration tests
6. **Documentation**: Add Swagger annotations to handlers

## Production Deployment

### Build for Production

```bash
make prod-build
```

### Docker Deployment

1. Build the Docker image:
```bash
make docker-build
```

2. Update `docker-compose.yml` to enable the API service

3. Start all services:
```bash
docker-compose up -d
```

## Environment Variables

See `.env.example` for all available configuration options.

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests and linting: `make test && make lint`
4. Format code: `make fmt`
5. Submit a pull request

## License

MIT

## Contact

For questions or support, contact the development team.

Cloudflared

cloudflared-windows-amd64.exe tunnel --url http://localhost:8080

cloudflared tunnel --url http://localhost:8080