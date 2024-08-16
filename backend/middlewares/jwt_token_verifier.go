package middlewares

import (
	"net/http"
	"strings"

	"github.com/daiki-kim/tweet-app/backend/pkg/auth"
	"github.com/gin-gonic/gin"
)

func JwtTokenVerifier() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get header
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		if authorizationHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// get token
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			return
		}

		// get token string
		tokenString := strings.TrimSpace(bearerToken[1])
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			return
		}

		// Validate token
		// TODO: serch about how to use verify key 2024-08-16
		claims, err := auth.ValidateToken(tokenString, auth.TokenVerifyKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// set email to context
		ctx.Set("email", claims.Email)

		ctx.Next()
	}
}
