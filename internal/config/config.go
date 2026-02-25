package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	Logger   LoggerConfig
	RateLimit RateLimitConfig
	Pagination PaginationConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port           string
	Host           string
	GinMode        string
	AllowedOrigins []string
}

// DatabaseConfig holds database configuration
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

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret              string
	AccessTokenExpire   time.Duration
	RefreshTokenExpire  time.Duration
}

// OAuthConfig holds OAuth providers configuration
type OAuthConfig struct {
	Google   OAuthProviderConfig
	Facebook OAuthProviderConfig
}

// OAuthProviderConfig holds individual OAuth provider configuration
type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level             string
	Encoding          string
	OutputPaths       []string
	ErrorOutputPaths  []string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
}

// PaginationConfig holds pagination configuration
type PaginationConfig struct {
	DefaultPageSize int
	MaxPageSize     int
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Try to read config file, but don't fail if it doesn't exist
	_ = viper.ReadInConfig()

	// Set defaults
	setDefaults()

	config := &Config{
		Server: ServerConfig{
			Port:           viper.GetString("SERVER_PORT"),
			Host:           viper.GetString("SERVER_HOST"),
			GinMode:        viper.GetString("GIN_MODE"),
			AllowedOrigins: viper.GetStringSlice("ALLOWED_ORIGINS"),
		},
		Database: DatabaseConfig{
			Host:               viper.GetString("DB_HOST"),
			Port:               viper.GetString("DB_PORT"),
			User:               viper.GetString("DB_USER"),
			Password:           viper.GetString("DB_PASSWORD"),
			DBName:             viper.GetString("DB_NAME"),
			SSLMode:            viper.GetString("DB_SSL_MODE"),
			MaxConnections:     viper.GetInt("DB_MAX_CONNECTIONS"),
			MaxIdleConnections: viper.GetInt("DB_MAX_IDLE_CONNECTIONS"),
			MaxLifetimeMinutes: viper.GetInt("DB_MAX_LIFETIME_MINUTES"),
		},
		JWT: JWTConfig{
			Secret:              viper.GetString("JWT_SECRET"),
			AccessTokenExpire:   viper.GetDuration("JWT_ACCESS_TOKEN_EXPIRE"),
			RefreshTokenExpire:  viper.GetDuration("JWT_REFRESH_TOKEN_EXPIRE"),
		},
		OAuth: OAuthConfig{
			Google: OAuthProviderConfig{
				ClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
				ClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
				RedirectURL:  viper.GetString("GOOGLE_REDIRECT_URL"),
			},
			Facebook: OAuthProviderConfig{
				ClientID:     viper.GetString("FACEBOOK_APP_ID"),
				ClientSecret: viper.GetString("FACEBOOK_APP_SECRET"),
				RedirectURL:  viper.GetString("FACEBOOK_REDIRECT_URL"),
			},
		},
		Logger: LoggerConfig{
			Level:             viper.GetString("LOG_LEVEL"),
			Encoding:          viper.GetString("LOG_ENCODING"),
			OutputPaths:       viper.GetStringSlice("LOG_OUTPUT_PATHS"),
			ErrorOutputPaths:  viper.GetStringSlice("LOG_ERROR_OUTPUT_PATHS"),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: viper.GetInt("RATE_LIMIT_REQUESTS_PER_MINUTE"),
		},
		Pagination: PaginationConfig{
			DefaultPageSize: viper.GetInt("DEFAULT_PAGE_SIZE"),
			MaxPageSize:     viper.GetInt("MAX_PAGE_SIZE"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate validates the configuration
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

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("ALLOWED_ORIGINS", []string{"http://localhost:3000"})

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("DB_MAX_CONNECTIONS", 20)
	viper.SetDefault("DB_MAX_IDLE_CONNECTIONS", 5)
	viper.SetDefault("DB_MAX_LIFETIME_MINUTES", 30)

	viper.SetDefault("JWT_ACCESS_TOKEN_EXPIRE", "15m")
	viper.SetDefault("JWT_REFRESH_TOKEN_EXPIRE", "168h") // 7 days

	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_ENCODING", "json")
	viper.SetDefault("LOG_OUTPUT_PATHS", []string{"stdout"})
	viper.SetDefault("LOG_ERROR_OUTPUT_PATHS", []string{"stderr"})

	viper.SetDefault("RATE_LIMIT_REQUESTS_PER_MINUTE", 100)

	viper.SetDefault("DEFAULT_PAGE_SIZE", 20)
	viper.SetDefault("MAX_PAGE_SIZE", 100)
}
