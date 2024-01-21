package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"

	"example/model"
)

func RequireAuth (context *gin.Context) {
	// Get the cookie off the request
	tokenStr, err := context.Cookie("Authorization")

	if err != nil {
		context.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode and validate token
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	
		return []byte(os.Getenv("JWT_PRIVATE_KEY")), nil
	})

	log.Printf("alg is %s\n", token.Header)
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check for expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			context.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find user with token sub
		user, err := model.FindUserById(uint(claims["id"].(float64)))
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
		}

		// Attach to req
		context.Set("user", user)

		// Continue
		context.Next()
	} else {
		context.AbortWithStatus(http.StatusUnauthorized)
	}
	

}