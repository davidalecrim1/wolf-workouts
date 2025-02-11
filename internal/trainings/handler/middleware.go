package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const JWT_CLAIMS_KEY = "jwt_claims"

func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			slog.Error("No authorization header found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			slog.Error("Invalid authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := validateTokenAndExtractClaims(token, config)
		if err != nil {
			slog.Error("Invalid token", "error", err, "token", token)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set(JWT_CLAIMS_KEY, claims)
		c.Next()
	}
}

func validateTokenAndExtractClaims(token string, config *config.Config) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
		}

		_, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}

		return config.GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}
