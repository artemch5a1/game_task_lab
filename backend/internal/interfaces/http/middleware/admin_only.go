package middleware

import (
	"net/http"
	"strings"

	appauth "example/web-service-gin/internal/application/abstraction/auth"
	"example/web-service-gin/internal/constants"
	specifictype "example/web-service-gin/internal/domain/specific_type"

	"github.com/gin-gonic/gin"
)

func RequireAdmin(verifier appauth.TokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			return
		}

		_, role, err := verifier.Verify(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			return
		}

		if role != specifictype.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": constants.ErrForbidden})
			return
		}

		c.Next()
	}
}

