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

func JWTMiddleware(repositoryToken repository.TokenRepositoryInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := utils.GetSecretKey()
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "please login first", "")
			c.Abort()
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		c.Set("accessToken", accessToken)
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}
			return []byte(secretKey), nil
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
			if repositoryToken.IsTokenRevoked(accessToken) {
				utils.ErrorResponse(c, http.StatusUnauthorized, "Token is revoked", "")
				c.Abort()
				return
			}

			customer := &models.Customer{
				ID:       claims["id"].(string),
				Username: claims["username"].(string),
			}
			c.Set("customer", customer)
		} else {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", "")
			c.Abort()
			return
		}

		c.Next()
	}
}
