package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Config содержит все конфигурационные параметры приложения
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
}

// ServerConfig содержит параметры HTTP сервера
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// DatabaseConfig содержит параметры подключения к базе данных
type DatabaseConfig struct {
	Driver   string `json:"driver"`
	DSN      string `json:"dsn"`
	LogLevel string `json:"log_level"`
	MaxIdle  int    `json:"max_idle"`
	MaxOpen  int    `json:"max_open"`
}

// JWTConfig содержит параметры для JWT токенов
type JWTConfig struct {
	SecretKey            string        `json:"secret_key"`
	AccessTokenDuration  time.Duration `json:"access_token_duration"`
	RefreshTokenDuration time.Duration `json:"refresh_token_duration"`
	Issuer               string        `json:"issuer"`
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "localhost"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", "30s"),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", "30s"),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", "60s"),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite3"),
			DSN:      getEnv("DB_DSN", "./courseforge.db"),
			LogLevel: getEnv("DB_LOG_LEVEL", "info"),
			MaxIdle:  getIntEnv("DB_MAX_IDLE", 10),
			MaxOpen:  getIntEnv("DB_MAX_OPEN", 100),
		},
		JWT: JWTConfig{
			SecretKey:            getEnv("JWT_SECRET_KEY", "courseforge-secret-key-change-in-production"),
			AccessTokenDuration:  getDurationEnv("JWT_ACCESS_TOKEN_DURATION", "24h"),
			RefreshTokenDuration: getDurationEnv("JWT_REFRESH_TOKEN_DURATION", "168h"), // 7 дней
			Issuer:               getEnv("JWT_ISSUER", "courseforge"),
		},
	}
}

// GetServerAddress возвращает адрес сервера в формате host:port
func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

// GetDatabaseDSN возвращает строку подключения к базе данных
func (c *Config) GetDatabaseDSN() string {
	return c.Database.DSN
}

// IsProduction проверяет, запущено ли приложение в продакшене
func (c *Config) IsProduction() bool {
	return getEnv("APP_ENV", "development") == "production"
}

// IsDevelopment проверяет, запущено ли приложение в режиме разработки
func (c *Config) IsDevelopment() bool {
	return getEnv("APP_ENV", "development") == "development"
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	// Проверяем обязательные поля
	if c.JWT.SecretKey == "" {
		return fmt.Errorf("JWT secret key is required")
	}

	if c.Database.DSN == "" {
		return fmt.Errorf("database DSN is required")
	}
	return nil
}

// Вспомогательные функции для получения переменных окружения

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if int64Value, err := strconv.ParseInt(value, 10, 64); err == nil {
			return int64Value
		}
		log.Printf("Warning: invalid int64 value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		log.Printf("Warning: invalid boolean value for %s: %s, using default: %t", key, value, defaultValue)
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	log.Printf("Warning: invalid duration value for %s: %s, using default: %s", key, value, defaultValue)
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

// LoadFromFile загружает конфигурацию из JSON файла (дополнительно к переменным окружения)
func LoadFromFile(filename string) (*Config, error) {
	config := Load() // Сначала загружаем из переменных окружения

	// Если файл существует, дополняем конфигурацию из него
	if _, err := os.Stat(filename); err == nil {
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, fmt.Errorf("failed to decode config file: %w", err)
		}
	}

	return config, nil
}

// SaveToFile сохраняет текущую конфигурацию в JSON файл
func (c *Config) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}
