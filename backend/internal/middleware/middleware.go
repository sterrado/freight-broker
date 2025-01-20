package middleware

import (
    "net/http"
    "strings"
    "freight-broker/backend/internal/services"
    "github.com/gin-gonic/gin"
)

type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func JWTAuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
            c.Abort()
            return
        }

        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
            c.Abort()
            return
        }

        claims, err := authService.ValidateToken(bearerToken[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        // Store user information in context
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)

        c.Next()
    }
}

// ErrorHandler middleware handles all errors in a consistent way
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // Only handle errors if there are any
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            // You can add custom error types and handle them differently
            switch e := err.Err.(type) {
            case *CustomError:
                c.JSON(e.StatusCode, ErrorResponse{
                    Code:    e.StatusCode,
                    Message: e.Message,
                })
            default:
                // Default error response
                c.JSON(http.StatusInternalServerError, ErrorResponse{
                    Code:    http.StatusInternalServerError,
                    Message: "An internal server error occurred",
                })
            }
        }
    }
}

// CustomError for handling specific error cases
type CustomError struct {
    StatusCode int
    Message    string
}

func (e *CustomError) Error() string {
    return e.Message
}