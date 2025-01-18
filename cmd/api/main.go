package main

import (
    "fmt"
    "log"
    "freight-broker/configs"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

func main() {
    config, err := configs.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    db, err := setupDatabase(config)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    r := gin.Default()
    // Setup routes and middleware here

    r.Run(":8080")
}

func setupDatabase(config *configs.Config) (*gorm.DB, error) {
    dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
    
    return gorm.Open("postgres", dbURL)
}