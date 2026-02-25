# Architecture Documentation

## Overview

This project follows **Clean Architecture** principles, ensuring separation of concerns and maintainability.

## Layers

### 1. Domain Layer (`internal/domain/`)

The core business logic layer containing:

- **Entities**: Core business models (User, Workout, Meal, Social)
- **Repository Interfaces**: Define data access contracts
- **Business Rules**: Domain-specific validation and logic

**Key Principle**: This layer has NO dependencies on external frameworks or libraries.

### 2. Use Case Layer (`internal/usecase/`)

Contains application-specific business logic:

- Orchestrates domain entities
- Implements business workflows
- Calls repository interfaces
- Independent of delivery mechanism (HTTP, gRPC, etc.)

**Example Use Cases**:
- `CreateWorkoutPlan`
- `ScheduleWorkout`
- `LogMealConsumption`
- `FollowUser`

### 3. Delivery Layer (`internal/delivery/`)

Handles external communication:

- **HTTP Handlers**: Gin route handlers
- **Middleware**: Authentication, logging, CORS, etc.
- **Request/Response DTOs**: Input validation and response formatting

**Responsibilities**:
- Receive HTTP requests
- Validate input
- Call use cases
- Return formatted responses

### 4. Infrastructure Layer (`internal/infrastructure/`)

External concerns:

- **Database**: PostgreSQL connection and pooling
- **Authentication**: JWT and OAuth2 implementation
- **Logger**: Zap logger wrapper
- **External APIs**: Third-party integrations

### 5. Repository Layer (`internal/repository/`)

Data access implementation:

- Implements domain repository interfaces
- Handles database queries
- Manages transactions
- PostgreSQL-specific logic

## Dependency Flow

```
HTTP Request → Handler → Use Case → Repository → Database
              ↓           ↓           ↓
          Middleware   Domain     Infrastructure
```

**Dependency Rule**: Dependencies point inward. Inner layers know nothing about outer layers.

## Key Design Patterns

### 1. Dependency Injection (Uber-go/fx)

All dependencies are injected through fx modules:

```go
fx.New(
    config.Module,
    logger.Module,
    database.Module,
    auth.Module,
    router.Module,
)
```

### 2. Repository Pattern

Domain defines interfaces, infrastructure provides implementations:

```go
// Domain
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}

// Infrastructure
type postgresUserRepository struct {
    db *database.DB
}
```

### 3. Middleware Chain

Modular middleware for cross-cutting concerns:

```
Request → Recovery → Logger → CORS → Rate Limit → Auth → Handler
```

## Data Flow Examples

### Workout Tracking Flow

1. **Exercise Library** (Pre-populated by admin)
   - System contains exercises with details
   - Users browse and search

2. **Create Workout Plan**
   - User selects exercises from library
   - Configures sets, reps, rest time
   - Saves as template or one-time use

3. **Schedule Workout**
   - User picks a date and time
   - Can create recurring schedules
   - Bulk schedule for week/month

4. **Execute Workout Session**
   - User starts session on scheduled date
   - Logs actual performance per exercise
   - System calculates calories burned
   - Completes session with notes and mood

5. **View Progress**
   - Historical data shows improvement
   - Compare planned vs actual performance
   - Statistics and charts

### Meal Tracking Flow

1. **Food Library**
   - System provides common foods
   - Users can add custom foods

2. **Create Recipe** (Optional)
   - Combine multiple foods
   - System auto-calculates nutrition

3. **Daily Meal Logging**
   - User logs meals as consumed
   - Select date and meal time
   - Add foods or recipes
   - Adjust portions

4. **Track Nutrition**
   - View daily totals
   - Compare with calorie target
   - Weekly/monthly statistics

## Security Considerations

### Authentication Flow

1. **Registration/Login**
   ```
   Client → POST /auth/register → Hash Password → Store in DB → Return JWT
   ```

2. **OAuth Flow**
   ```
   Client → GET /auth/oauth/google → Google Auth → Callback → Create/Link User → Return JWT
   ```

3. **Protected Routes**
   ```
   Client → Add Bearer Token → Middleware validates JWT → Extract User ID → Allow Request
   ```

### Password Security

- Passwords hashed with bcrypt (cost 10)
- Never stored in plain text
- Never returned in API responses

### JWT Security

- Access tokens: Short-lived (15 minutes)
- Refresh tokens: Long-lived (7 days)
- Tokens signed with HMAC-SHA256
- Secret key from environment variables

## Error Handling

### Error Types

```go
errors.BadRequest()      // 400
errors.Unauthorized()    // 401
errors.Forbidden()       // 403
errors.NotFound()        // 404
errors.Conflict()        // 409
errors.Validation()      // 422
errors.InternalServer()  // 500
```

### Error Response Format

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": {
      "email": "email must be a valid email address"
    }
  }
}
```

## Database Design Principles

### Normalization

- Tables follow 3NF (Third Normal Form)
- Minimal data redundancy
- Foreign keys enforce referential integrity

### Indexing Strategy

- Primary keys on all tables
- Foreign keys indexed
- Frequently queried fields indexed
- Composite indexes for common queries

### Performance Optimization

- Connection pooling (20 max connections)
- Prepared statements for repeated queries
- Pagination for large result sets
- Selective field fetching

## Testing Strategy

### Unit Tests

Test individual functions in isolation:

```go
func TestHashPassword(t *testing.T) {
    pm := NewPasswordManager()
    hash, err := pm.HashPassword("password123")
    assert.NoError(t, err)
    assert.True(t, pm.VerifyPassword(hash, "password123"))
}
```

### Integration Tests

Test repository implementations with test database:

```go
func TestUserRepository_Create(t *testing.T) {
    // Setup test database
    // Create user
    // Assert user was created
    // Cleanup
}
```

### Handler Tests

Test HTTP endpoints with httptest:

```go
func TestHealthCheck(t *testing.T) {
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/health", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
}
```

## Configuration Management

### Environment-based Config

- Development: `.env` file
- Production: Environment variables
- Default values for optional settings

### Configuration Layers

1. **Defaults**: Hardcoded safe defaults
2. **File**: `.env` or config file
3. **Environment**: OS environment variables (highest priority)

## Logging Standards

### Log Levels

- **Debug**: Detailed diagnostic information
- **Info**: General informational messages
- **Warn**: Warning messages (non-critical)
- **Error**: Error messages (needs attention)
- **Fatal**: Critical errors (application exits)

### Structured Logging

```go
log.Info("User created",
    zap.String("user_id", user.ID.String()),
    zap.String("email", user.Email),
)
```

### Log Context

Each request includes:
- Request ID
- User ID (if authenticated)
- HTTP method and path
- Response status
- Latency

## API Versioning

Current API version: **v1**

Base path: `/api/v1`

Future versions will use: `/api/v2`, `/api/v3`, etc.

### Versioning Strategy

- Major version in URL path
- Backward compatibility within major version
- Deprecation warnings before breaking changes
- Documentation for migration guides

## Deployment Architecture

### Production Setup

```
[Load Balancer]
      ↓
[API Server 1] [API Server 2] [API Server 3]
      ↓              ↓              ↓
[PostgreSQL Primary]  ←→  [PostgreSQL Replica]
      ↓
[Redis Cache]
```

### Scaling Considerations

- **Horizontal Scaling**: Multiple API server instances
- **Database Read Replicas**: Separate read/write
- **Caching Layer**: Redis for sessions and frequently accessed data
- **CDN**: Static assets and media files
- **Queue System**: Background jobs (future enhancement)

## Monitoring and Observability

### Health Checks

- `/health` - Basic health check
- Database connection status
- External service availability

### Metrics (Future)

- Request rate
- Response times
- Error rates
- Database query performance
- Active connections

### Logging

- Structured JSON logs
- Centralized log aggregation
- Error tracking and alerting

## Best Practices

1. **Always use context.Context** for cancellation and timeouts
2. **Never return passwords** in API responses
3. **Validate all inputs** using validator tags
4. **Log errors with context** for debugging
5. **Use transactions** for multi-step operations
6. **Handle errors gracefully** with proper status codes
7. **Document APIs** with Swagger annotations
8. **Write tests** for critical business logic
9. **Follow Go conventions** and idioms
10. **Keep functions small** and focused

## Future Enhancements

1. **Redis Integration**: Session management and caching
2. **WebSocket Support**: Real-time notifications
3. **File Upload**: Profile pictures, meal photos
4. **Email Service**: Verification, notifications
5. **Push Notifications**: Mobile push via FCM
6. **Analytics Dashboard**: Admin panel
7. **Rate Limiting**: Per-user API quotas
8. **API Gateway**: Centralized routing and auth
9. **Microservices**: Split into smaller services
10. **GraphQL**: Alternative API interface
