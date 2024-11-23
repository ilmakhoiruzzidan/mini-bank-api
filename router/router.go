package router

import (
	"github.com/gin-gonic/gin"
	"mini-bank-api/handlers"
	"mini-bank-api/repository"
	"mini-bank-api/services"
	"net/http"
)

func InitRouter(router *gin.Engine) {
	router.GET("/", homePage)

	// router customer
	customerRepository := repository.NewJSONCustomerRepository()
	customerService := services.NewCustomerService(customerRepository)
	customerHandler := handlers.NewCustomerHandler(customerService)

	router.GET("/customers", customerHandler.GetAllCustomers)

	// router auth
	authService := services.NewAuthService(customerRepository)
	authHandler := handlers.NewAuthHandler(authService)

	router.POST("/auth/login", authHandler.Login)

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
