package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pintarlabs/license-server/internal/config"
	"github.com/pintarlabs/license-server/internal/middleware"
)

func RequireAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			middleware.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing or invalid token")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			middleware.ErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN", err.Error())
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole != role {
			middleware.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "You do not have permission to access this resource")
			return
		}
		c.Next()
	}
}
