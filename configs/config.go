package configs

import (
    "fmt"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    TurvoAPIKey   string
    ClientName    string
    ClientSecret  string
    IsSandbox     bool
}

func LoadConfig() (*Config, error) {
    // Only load .env file if not in production
    if os.Getenv("GO_ENV") != "production" {
        if err := godotenv.Load(); err != nil {
            // Don't return error as .env file is optional
            fmt.Println("Warning: .env file not found")
        }
    }

    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "freight_broker"),
        TurvoAPIKey:   getEnv("TURVO_API_KEY", ""),
        ClientName:    getEnv("CLIENT_NAME", ""),
        ClientSecret:  getEnv("CLIENT_SECRET", ""),
        IsSandbox:     getEnv("ENVIRONMENT", "sandbox") == "sandbox",
    }, nil
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}