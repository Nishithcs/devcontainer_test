package config

import (
	"clusterix-code/internal/constants"
	"clusterix-code/internal/utils/di"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server           ServerConfig
	Database         DatabaseConfig
	Auth             AuthConfig
	Logger           LoggerConfig
	RabbitMQ         RabbitMQConfig
	ExternalServices ExternalServicesConfig
	Redis            RedisConfig
	MongoDB          MongoDBConfig
}

type ExternalServicesConfig struct {
	Auth AuthServiceConfig
}

type AuthServiceConfig struct {
	BaseURL     string
	MasterToken string
}

type ServerConfig struct {
	Port            int
	Environment     string
	AllowedOrigins  []string
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	MaxConns int
	Timeout  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type MongoDBConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Database   string
	AuthSource string
}

type AuthConfig struct {
	JWTSecret string
}

type LoggerConfig struct {
	LogLevel string
	Format   string
}

type RabbitMQConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DefaultExchange string
	Protocol        string
	Queues          []QueueConfig
}

type QueueConfig struct {
	Name         string
	ShouldCreate bool
	Bindings     []QueueBindingConfig
}

type QueueBindingConfig struct {
	ExchangeName string
	RoutingKey   string
}

func Provider(c *di.Container) (*Config, error) {

	return NewConfig()
}

func NewConfig() (*Config, error) {
	rabbitMQQueues := []QueueConfig{
		{
			Name:         constants.AUTH_USER_CREATED_QUEUE,
			ShouldCreate: true,
			Bindings: []QueueBindingConfig{
				{
					ExchangeName: constants.CLUSTERIX_CODE_V1_EXCHANGE,
					RoutingKey:   constants.AUTH_USER_CREATED_QUEUE,
				},
			},
		},
		{
			Name:         constants.AUTH_USER_UPDATED_QUEUE,
			ShouldCreate: true,
			Bindings: []QueueBindingConfig{
				{
					ExchangeName: constants.CLUSTERIX_CODE_V1_EXCHANGE,
					RoutingKey:   constants.AUTH_USER_UPDATED_QUEUE,
				},
			},
		},
		{
			Name:         constants.AUTH_USER_DELETED_QUEUE,
			ShouldCreate: true,
			Bindings: []QueueBindingConfig{
				{
					ExchangeName: constants.CLUSTERIX_CODE_V1_EXCHANGE,
					RoutingKey:   constants.AUTH_USER_DELETED_QUEUE,
				},
			},
		},
		{
			Name:         constants.WORKSPACE_LOG_HANDLER_QUEUE,
			ShouldCreate: true,
			Bindings: []QueueBindingConfig{
				{
					ExchangeName: constants.CLUSTERIX_CODE_V1_EXCHANGE,
					RoutingKey:   constants.WORKSPACE_LOG_HANDLER_QUEUE,
				},
			},
		},
	}

	return &Config{
		Server: ServerConfig{
			Port:            getEnvAsInt("SERVER_PORT", 8080),
			Environment:     GetEnv("APP_ENV", "development"),
			AllowedOrigins:  getEnvAsSlice("ALLOWED_ORIGINS", []string{"*"}),
			ShutdownTimeout: getEnvAsDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "postgres"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     GetEnv("DB_USER", "postgres"),
			Password: GetEnv("DB_PASSWORD", "postgres"),
			Name:     GetEnv("DB_NAME", "postgres"),
			SSLMode:  GetEnv("DB_SSLMODE", "disable"),
			MaxConns: getEnvAsInt("DB_MAX_CONNS", 10),
			Timeout:  getEnvAsDuration("DB_TIMEOUT", 5*time.Second),
		},
		Auth: AuthConfig{
			JWTSecret: GetEnv("AUTH_JWT_SECRET", "your-secret-key"),
		},
		Logger: LoggerConfig{
			LogLevel: GetEnv("LOG_LEVEL", "info"),
			Format:   GetEnv("LOG_FORMAT", "json"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:            GetEnv("RABBITMQ_HOST", "rabbitmq"),
			Port:            getEnvAsInt("RABBITMQ_PORT", 5672),
			User:            GetEnv("RABBITMQ_USER", "guest"),
			Password:        GetEnv("RABBITMQ_PASS", "guest"),
			DefaultExchange: GetEnv("RABBITMQ_DEFAULT_EXCHANGE", "innoscripta"),
			Protocol:        GetEnv("RABBITMQ_PROTOCOL", "amqps"),
			Queues:          rabbitMQQueues,
		},
		ExternalServices: ExternalServicesConfig{
			Auth: AuthServiceConfig{
				BaseURL:     GetEnv("AUTH_API_URL", ""),
				MasterToken: GetEnv("AUTH_API_MASTER_TOKEN", ""),
			},
		},
		Redis: RedisConfig{
			Host:     GetEnv("REDIS_HOST", "redis"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Username: GetEnv("REDIS_USERNAME", ""),
			Password: GetEnv("REDIS_PASSWORD", ""),
		},
		MongoDB: MongoDBConfig{
			Host:       GetEnv("MONGO_HOST", "localhost"),
			Port:       getEnvAsInt("MONGO_PORT", 27017),
			Username:   GetEnv("MONGO_USERNAME", "admin"),
			Password:   GetEnv("MONGO_PASSWORD", "password"),
			Database:   GetEnv("MONGO_DB", "clusterix"),
			AuthSource: GetEnv("MONGO_AUTH_SOURCE", "admin"),
		},
	}, nil
}

// GetEnv Helper functions for environment variables
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return defaultValue
}
