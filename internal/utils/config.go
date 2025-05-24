package utils

import (
	"os"
	"strconv"
)

type Config struct {
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBName     string
    JWTSecret  string
    Port       string
}

func LoadConfig() *Config {
    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnvAsInt("DB_PORT", 5432),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "1984"),
        DBName:     getEnv("DB_NAME", "nethernode"),
        JWTSecret:  getEnv("JWT_SECRET", "nethernodeserver"),
        Port:       getEnv("PORT", "8080"),
    }
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return defaultValue
}
