package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mini-bank-api/models"
	"mini-bank-api/repository"
	"mini-bank-api/utils"
	"net/http"
	"strings"
)

func JWTMiddleware(customerRepo repository.CustomerRepositoryInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := utils.GetSecretKey()
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

			secretKey := []byte(secretKey)
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Next()

	}
}
