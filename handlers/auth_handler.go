package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-bank-api/services"
	"net/http"
)

type AuthHandler struct {
	Service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{
		Service: service,
	}
}

func (handler *AuthHandler) Login(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	token, err := handler.Service.Login(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "login success",
		"accessToken": token,
	})
}
