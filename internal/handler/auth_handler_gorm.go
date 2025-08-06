package handler

import (
	"net/http"

	"go-debt-tracker/internal/models"
	"go-debt-tracker/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandlerGORM struct {
	authService *service.AuthServiceGORM
}

func NewAuthHandlerGORM(authService *service.AuthServiceGORM) *AuthHandlerGORM {
	return &AuthHandlerGORM{authService: authService}
}

func (h *AuthHandlerGORM) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    response.User,
		"sharing_summary": response.SharingSummary,
	})
}

func (h *AuthHandlerGORM) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
} 