package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-bank-api/services"
	"mini-bank-api/utils"
	"net/http"
)

type AuthHandler struct {
	authService services.AuthServiceInterface
}

func NewAuthHandler(authService services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (handler *AuthHandler) Login(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid credentials", err.Error())
		return
	}

	accessToken, err := handler.authService.Login(request.Username, request.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	utils.LoginResponse(c, accessToken, "login success")
}

func (handler *AuthHandler) Logout(c *gin.Context) {
	accessToken, exists := c.Get("accessToken")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "No access token")
		return
	}

	tokenString, ok := accessToken.(string)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid token format", "")
		return
	}

	err := handler.authService.Logout(tokenString)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Failed to logout", err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "logout success")
}

func (handler *AuthHandler) GetCurrentUser(c *gin.Context) {
	accessToken, exists := c.Get("accessToken")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Please login first", "")
		return
	}
	tokenString, ok := accessToken.(string)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid token format", "")
		return
	}

	currentUser, err := handler.authService.GetCurrentUserInfo(tokenString)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "You're not authorized to access this resource", err.Error())
	}

	utils.SuccessResponse(c, currentUser, "success retrieve current user")

}
