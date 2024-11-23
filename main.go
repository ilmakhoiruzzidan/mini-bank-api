package main

import (
	"github.com/gin-gonic/gin"
	"mini-bank-api/initializers"
	"mini-bank-api/router"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	// instate router
	r := gin.Default()

	router.InitRouter(r)
}
