package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"mini-bank-api/models"
	"mini-bank-api/repository"
	"net/http"
	"os"
	"strings"
)

func getSecretKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("Error loading JWT_SECRET_KEY file")
	}
	return secretKey
}

func JWTMiddleware(customerRepo repository.CustomerRepositoryInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  "You need to be logged in",
			})
			c.Abort()
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		c.Set("accessToken", accessToken)

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"])
			}
			secretKey := []byte(getSecretKey())
			return secretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  err.Error(),
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			customer := &models.Customer{
				ID:       claims["id"].(string),
				Username: claims["username"].(string),
			}

			cust, err := customerRepo.FindCustomerByID(customer.ID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": http.StatusUnauthorized,
					"error":  "Customer not found",
				})
				c.Abort()
				return
			}

			if cust.IsLoggedOut {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": http.StatusUnauthorized,
					"error":  "You have logged out. Please log in again",
				})
				c.Abort()
				return
			}

			c.Set("customer", customer)

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Next()

	}
}
