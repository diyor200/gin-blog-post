package middleware

import (
	"github.com/diyor200/gin-middleware-blogpost/internal/controller"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func CheckUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, ok := bearerToken(ctx.Request)
		if !ok {
			log.Println("not authorized")
			ctx.Next()
			return
		}
		log.Println("incomed token", token)
		userID, err := controller.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			ctx.Next()
			return
		}
		log.Println("authorized --- userID=== ", userID)
		ctx.Set("user_id", userID)
		ctx.Next()
		return
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	token := strings.Split(header, " ")[1]

	if token == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}
	return "", false
}
