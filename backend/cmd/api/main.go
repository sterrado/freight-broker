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
    
    "freight-broker/backend/configs"
    "freight-broker/backend/internal/services"
    "freight-broker/backend/internal/controllers"
    "freight-broker/backend/internal/models"
    "freight-broker/backend/internal/middleware"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "github.com/gin-contrib/cors"
    _ "github.com/lib/pq"
)

func main() {
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

    if err := setupModels(db); err != nil {
        log.Fatalf("Failed to setup database models: %v", err)
    }

    authService := services.NewAuthService(config.JWTSecret)
    tmsService := services.NewTurvoService(services.TMSServiceConfig{
        APIKey:       config.TurvoAPIKey,
        ClientID:     config.ClientName,
        ClientSecret: config.ClientSecret,
        IsSandbox:    config.IsSandbox,
        TurvoUsername:  config.TurvoUsername,
        TurvoPassword:  config.TurvoPassword,
    })
    loadService := services.NewLoadService(db, tmsService)

    if err := tmsService.Authenticate(context.Background()); err != nil {
        log.Fatalf("Failed to authenticate with Turvo: %v", err)
    }

    // Initialize controllers
    authController := controllers.NewAuthController(authService)
    loadController := controllers.NewLoadController(loadService, tmsService)

    gin.SetMode(getGinMode())
    r := gin.New()

    // Add middleware
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: false,
        MaxAge:           12 * time.Hour,
    }))
    r.Use(gin.Recovery())
    r.Use(gin.Logger())
    
    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })
    
    

    api := r.Group("/api")
    {
        auth := api.Group("/auth")
        {
            auth.POST("/login", authController.Login)
        }

        protected := api.Group("")
        protected.Use(middleware.JWTAuthMiddleware(authService))
        {
            loads := protected.Group("/loads")
            {
                loads.POST("/", loadController.CreateLoad)
                loads.GET("/", loadController.ListLoads)
                loads.GET("/:id", loadController.GetLoad)
            }
        }
    }


    srv := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

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

    db.DB().SetMaxIdleConns(10)
    db.DB().SetMaxOpenConns(100)
    db.DB().SetConnMaxLifetime(time.Hour)

    return db, nil
}

func setupModels(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.Load{},
    ).Error
}

func getGinMode() string {
    mode := os.Getenv("GIN_MODE")
    if mode == "" {
        return gin.DebugMode
    }
    return mode
}