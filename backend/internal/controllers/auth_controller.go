// internal/controllers/auth_controller.go

package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"freight-broker/backend/internal/services"
)

type AuthController struct {
    authService *services.AuthService
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string `json:"token"`
}

func NewAuthController(authService *services.AuthService) *AuthController {
    return &AuthController{
        authService: authService,
    }
}

func (c *AuthController) Login(ctx *gin.Context) {
    var req LoginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // In a real app, we would validate credentials against a DB
    if !validateCredentials(req.Username, req.Password) {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    // Generate JWT token
    token, err := c.authService.GenerateToken("user123", req.Username, "broker")
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
        return
    }

    ctx.JSON(http.StatusOK, LoginResponse{Token: token})
}

func validateCredentials(username, password string) bool {
    // Not a real app implementation
    return username == "admin" && password == "password"
}