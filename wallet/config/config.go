package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all service configuration settings.
type Config struct {
	// App settings
	AppPort string
	AppEnv  string

	// Database settings
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBMaxConns        int32
	DBMinConns        int32
	DBMaxIdleTime     time.Duration
	DBMaxConnLifetime time.Duration

	// Server settings
	ServerReadTimeout     time.Duration
	ServerWriteTimeout    time.Duration
	ServerShutdownTimeout time.Duration
}

// Load reads configuration from config.env file and environment variables.
func Load(envPath ...string) (*Config, error) {
	path := "config.env"
	if len(envPath) > 0 && envPath[0] != "" {
		path = envPath[0]
	}

	// Try reading config.env if it exists
	if err := loadEnvFile(path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to parse env file %s: %w", path, err)
	}

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "production"),

		DBHost:            getEnv("DB_HOST", "postgres"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "wallet_user"),
		DBPassword:        getEnv("DB_PASSWORD", "wallet_password"),
		DBName:            getEnv("DB_NAME", "wallet_db"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		DBMaxConns:        int32(getEnvAsInt("DB_MAX_CONNS", 100)),
		DBMinConns:        int32(getEnvAsInt("DB_MIN_CONNS", 20)),
		DBMaxIdleTime:     getEnvAsDuration("DB_MAX_IDLE_TIME", 5*time.Minute),
		DBMaxConnLifetime: getEnvAsDuration("DB_MAX_CONN_LIFETIME", time.Hour),

		ServerReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second),
		ServerWriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		ServerShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
	}

	return cfg, nil
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)

			// Only set if not already set in OS environment
			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, val)
			}
		}
	}

	return scanner.Err()
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	valStr := getEnv(key, "")
	if val, err := time.ParseDuration(valStr); err == nil {
		return val
	}
	return defaultVal
}
