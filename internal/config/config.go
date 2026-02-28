package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	Logger   LoggerConfig
	RateLimit RateLimitConfig
	Pagination PaginationConfig
}

// ServerConfig 
type ServerConfig struct {
	Port           string
	Host           string
	GinMode        string
	AllowedOrigins []string
}

// DatabaseConfig 
type DatabaseConfig struct {
	Host               string
	Port               string
	User               string
	Password           string
	DBName             string
	SSLMode            string
	MaxConnections     int
	MaxIdleConnections int
	MaxLifetimeMinutes int
}

// JWTConfig 
type JWTConfig struct {
	Secret              string
	AccessTokenExpire   time.Duration
	RefreshTokenExpire  time.Duration
}

// OAuthConfig 
type OAuthConfig struct {
	Google   OAuthProviderConfig
	Facebook OAuthProviderConfig
}

// OAuthProviderConfig 
type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// LoggerConfig 
type LoggerConfig struct {
	Level             string
	Encoding          string
	OutputPaths       []string
	ErrorOutputPaths  []string
}

// RateLimitConfig 
type RateLimitConfig struct {
	RequestsPerMinute int
}

// PaginationConfig 
type PaginationConfig struct {
	DefaultPageSize int
	MaxPageSize     int
}

// Load configuration 
func Load() (*Config, error) {
	// Load .env 
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	getEnv := func(key, defaultVal string) string {
		if val := os.Getenv(key); val != "" {
			return val
		}
		return defaultVal
	}

	getEnvInt := func(key string, defaultVal int) int {
		if val := os.Getenv(key); val != "" {
			if intVal, err := strconv.Atoi(val); err == nil {
				return intVal
			}
		}
		return defaultVal
	}

	getEnvDuration := func(key string, defaultVal time.Duration) time.Duration {
		if val := os.Getenv(key); val != "" {
			if duration, err := time.ParseDuration(val); err == nil {
				return duration
			}
		}
		return defaultVal
	}

	// Parse allowed origins
	allowedOriginsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	allowedOrigins := strings.Split(allowedOriginsStr, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	// Parse log paths
	logOutputPaths := strings.Split(getEnv("LOG_OUTPUT_PATHS", "stdout"), ",")
	logErrorPaths := strings.Split(getEnv("LOG_ERROR_OUTPUT_PATHS", "stderr"), ",")
	for i := range logOutputPaths {
		logOutputPaths[i] = strings.TrimSpace(logOutputPaths[i])
	}
	for i := range logErrorPaths {
		logErrorPaths[i] = strings.TrimSpace(logErrorPaths[i])
	}

	config := &Config{
		Server: ServerConfig{
			Port:           getEnv("SERVER_PORT", "8080"),
			Host:           getEnv("SERVER_HOST", "0.0.0.0"),
			GinMode:        getEnv("GIN_MODE", "debug"),
			AllowedOrigins: allowedOrigins,
		},
		Database: DatabaseConfig{
			Host:               getEnv("DB_HOST", "localhost"),
			Port:               getEnv("DB_PORT", "5432"),
			User:               getEnv("DB_USER", "gymadmin"),
			Password:           getEnv("DB_PASSWORD", "secret123"),
			DBName:             getEnv("DB_NAME", "gym_pro_db"),
			SSLMode:            getEnv("DB_SSL_MODE", "disable"),
			MaxConnections:     getEnvInt("DB_MAX_CONNECTIONS", 20),
			MaxIdleConnections: getEnvInt("DB_MAX_IDLE_CONNECTIONS", 5),
			MaxLifetimeMinutes: getEnvInt("DB_MAX_LIFETIME_MINUTES", 30),
		},
		JWT: JWTConfig{
			Secret:              getEnv("JWT_SECRET", ""),
			AccessTokenExpire:   getEnvDuration("JWT_ACCESS_TOKEN_EXPIRE", 15*time.Minute),
			RefreshTokenExpire:  getEnvDuration("JWT_REFRESH_TOKEN_EXPIRE", 168*time.Hour),
		},
		OAuth: OAuthConfig{
			Google: OAuthProviderConfig{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
			},
			Facebook: OAuthProviderConfig{
				ClientID:     getEnv("FACEBOOK_APP_ID", ""),
				ClientSecret: getEnv("FACEBOOK_APP_SECRET", ""),
				RedirectURL:  getEnv("FACEBOOK_REDIRECT_URL", ""),
			},
		},
		Logger: LoggerConfig{
			Level:             getEnv("LOG_LEVEL", "debug"),
			Encoding:          getEnv("LOG_ENCODING", "json"),
			OutputPaths:       logOutputPaths,
			ErrorOutputPaths:  logErrorPaths,
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 100),
		},
		Pagination: PaginationConfig{
			DefaultPageSize: getEnvInt("DEFAULT_PAGE_SIZE", 20),
			MaxPageSize:     getEnvInt("MAX_PAGE_SIZE", 100),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate config
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

// Build database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

