package auth

import (
	"errors"
	"strings"

	"github.com/ebrickdev/ebrick/logger"
	"github.com/ebrickdev/ebrick/transport/http"
)

const (
	bearerPrefix = "Bearer "
)

type AuthMiddleware struct {
	authManager AuthManager
	logger      logger.Logger
}

func NewAuthMiddleware(authManager AuthManager, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authManager: authManager,
		logger:      logger,
	}
}

func (am *AuthMiddleware) TokenAuth() http.HandlerFunc {
	return func(c *http.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			am.logger.Warn("Authorization header missing")
			c.AbortWithError(http.StatusUnauthorized, errors.New("missing token"))
			return
		}

		if len(authHeader) < len(bearerPrefix) || !strings.EqualFold(authHeader[:len(bearerPrefix)], bearerPrefix) {
			am.logger.Warn("Authorization header missing Bearer prefix")
			c.AbortWithError(http.StatusUnauthorized, errors.New("invalid token format"))
			return
		}

		tokenString := strings.TrimSpace(authHeader[len(bearerPrefix):])
		if tokenString == "" {
			am.logger.Warn("Empty token after Bearer prefix removal")
			c.AbortWithError(http.StatusUnauthorized, errors.New("missing token"))
			return
		}

		principal, err := am.authManager.Authenticate(c.Request.Context(), tokenString)
		if err != nil {
			am.logger.Error("Token authentication failed", logger.Error(err))
			if errors.Is(err, ErrTokenExpired) {
				c.AbortWithError(http.StatusUnauthorized, errors.New("token expired"))
			} else {
				c.AbortWithError(http.StatusUnauthorized, errors.New("invalid or malformed token"))
			}
			return
		}

		c.Set("claims", principal)
		c.Set("user_id", principal.GetID())
		c.Next()
	}
}

func (am *AuthMiddleware) RequireRoles(requiredRoles ...string) http.HandlerFunc {
	return func(c *http.Context) {
		rawClaims, exists := c.Get("claims")
		if !exists {
			am.logger.Warn("Claims not found in context")
			c.AbortWithError(http.StatusForbidden, ErrAccessForbidden)
			return
		}

		principal, ok := rawClaims.(Principal)
		if !ok {
			am.logger.Warn("Invalid claims type in context")
			c.AbortWithError(http.StatusForbidden, ErrAccessForbidden)
			return
		}

		if !hasAnyRole(principal.GetRoles(), requiredRoles) {
			c.AbortWithError(http.StatusForbidden, ErrAccessForbidden)
			return
		}

		c.Next()
	}
}

var ErrTokenExpired = errors.New("token expired")
var ErrAccessForbidden = errors.New("access forbidden")

func hasAnyRole(userRoles, requiredRoles []string) bool {
	for _, role := range requiredRoles {
		if containsString(userRoles, role) {
			return true
		}
	}
	return false
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
