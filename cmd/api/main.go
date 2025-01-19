package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    "net/http"
    
    "freight-broker/configs"
    "freight-broker/internal/services"
    "freight-broker/internal/controllers"
    "freight-broker/internal/models"
    "freight-broker/internal/middleware"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

func main() {
    // Load configuration
    config, err := configs.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize database
    db, err := setupDatabase(config)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Auto-migrate database models
    if err := setupModels(db); err != nil {
        log.Fatalf("Failed to setup database models: %v", err)
    }

    // Initialize services
    tmsService := services.NewTurvoService(services.TMSServiceConfig{
        APIKey:       config.TurvoAPIKey,
        ClientID:     config.ClientName,
        ClientSecret: config.ClientSecret,
        IsSandbox:    config.IsSandbox,
    })

    // Initial authentication
    if err := tmsService.Authenticate(context.Background()); err != nil {
        log.Fatalf("Failed to authenticate with Turvo: %v", err)
    }

    // Initialize controllers
    loadController := controllers.NewLoadController(db, tmsService)

    // Setup Gin router
    gin.SetMode(getGinMode())
    r := gin.New() // Use New() instead of Default() for custom middleware

    // Add middleware
    r.Use(gin.Recovery())
    r.Use(gin.Logger())
    r.Use(middleware.ErrorHandler())
    
    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })

    // Setup API routes
    api := r.Group("/api")
    {
        // Apply TMS token middleware to all API routes
        api.Use(middleware.TokenMiddleware(tmsService))

        loads := api.Group("/loads")
        {
            loads.POST("/", loadController.CreateLoad)
            loads.GET("/", loadController.ListLoads)
            loads.GET("/:id", loadController.GetLoad)
            loads.PUT("/:id", loadController.UpdateLoad)
            loads.DELETE("/:id", loadController.DeleteLoad)
        }
    }

    // Configure server
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Start server in a goroutine
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the server
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    // Give outstanding requests 5 seconds to complete
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exiting")
}

func setupDatabase(config *configs.Config) (*gorm.DB, error) {
    dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
    
    db, err := gorm.Open("postgres", dbURL)
    if err != nil {
        return nil, err
    }

    // Set connection pool settings
    db.DB().SetMaxIdleConns(10)
    db.DB().SetMaxOpenConns(100)
    db.DB().SetConnMaxLifetime(time.Hour)

    return db, nil
}

func setupModels(db *gorm.DB) error {
    // Auto-migrate your models here
    return db.AutoMigrate(
        &models.Load{},
        // Add other models here
    ).Error
}

func getGinMode() string {
    mode := os.Getenv("GIN_MODE")
    if mode == "" {
        return gin.DebugMode
    }
    return mode
}