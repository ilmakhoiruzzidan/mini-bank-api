package router

import (
	"github.com/gin-gonic/gin"
	"mini-bank-api/handlers"
	"mini-bank-api/middlewares"
	"mini-bank-api/repository"
	"mini-bank-api/services"
	"net/http"
)

func InitRouter(router *gin.Engine) {
	router.GET("/", homePage)

	// DI customer
	customerRepository := repository.NewJSONCustomerRepository()
	customerService := services.NewCustomerService(customerRepository)
	customerHandler := handlers.NewCustomerHandler(customerService)

	// DI auth
	historyRepository := repository.NewJSONHistoryRepository()
	authService := services.NewAuthService(customerRepository, historyRepository)
	authHandler := handlers.NewAuthHandler(authService)

	// public route
	router.POST("/auth/login", authHandler.Login)

	// protected (need auth)
	protected := router.Group("/")
	protected.Use(middlewares.JWTMiddleware())
	{
		protected.POST("/auth/logout", authHandler.Logout)
		protected.GET("/customers", customerHandler.GetAllCustomers)

	}

	err := router.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func homePage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
