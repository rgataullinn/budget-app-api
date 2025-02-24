package middleware

import (
	"fmt"
	"net/http"
	"os"
	"personal-finance-api/db"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Auth")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Print(claims["sub"].(float64))
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, err := db.GetUser(int(claims["sub"].(float64)))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("userId", user.Id)

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

}
