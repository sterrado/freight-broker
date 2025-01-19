package middleware

import (
    "net/http"
    
    "freight-broker/internal/interfaces"
    "github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

// TokenMiddleware ensures TMS token is valid before processing requests
func TokenMiddleware(tmsService interfaces.TMSService) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !tmsService.IsTokenValid() {
            if err := tmsService.RefreshToken(c.Request.Context()); err != nil {
                c.JSON(http.StatusUnauthorized, ErrorResponse{
                    Code:    http.StatusUnauthorized,
                    Message: "Failed to refresh TMS token",
                })
                c.Abort()
                return
            }
        }
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