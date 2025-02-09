package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ebrickdev/ebrick/logger"
	"github.com/ebrickdev/ebrick/transport/httpserver"
)

const (
	// bearerPrefix is the expected prefix in the Authorization header.
	bearerPrefix = "Bearer "
)

// AuthMiddleware provides middleware for token authentication and role-based authorization.
type AuthMiddleware struct {
	authManager AuthManager
	logger      logger.Logger
}

// NewAuthMiddleware creates a new AuthMiddleware instance.
func NewAuthMiddleware(authManager AuthManager, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authManager: authManager,
		logger:      logger,
	}
}

// TokenAuth is middleware that validates a JWT token from the Authorization header.
func (am *AuthMiddleware) TokenAuth() httpserver.HandlerFunc {
	return func(c httpserver.Context) {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			am.logger.Warn("Authorization header missing")
			abortWithError(c, http.StatusUnauthorized, "missing token")
			return
		}

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			am.logger.Warn("Authorization header missing Bearer prefix")
			abortWithError(c, http.StatusUnauthorized, "invalid token format")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			am.logger.Warn("Empty token after Bearer prefix removal")
			abortWithError(c, http.StatusUnauthorized, "missing token")
			return
		}

		principal, err := am.authManager.Authenticate(c.Request().Context(), tokenString)
		if err != nil {
			am.logger.Error("Token authentication failed", logger.Error(err))
			if isTokenExpired(err) {
				abortWithError(c, http.StatusUnauthorized, "token expired")
			} else {
				abortWithError(c, http.StatusUnauthorized, "invalid or malformed token")
			}
			return
		}

		// Store authentication claims and identifier in the context for later use.
		c.Set("claims", principal)
		c.Set("user_id", principal.GetID())
		c.Next()
	}
}

// RequireRoles is middleware that restricts access to endpoints based on user roles.
func (am *AuthMiddleware) RequireRoles(requiredRoles ...string) httpserver.HandlerFunc {
	return func(c httpserver.Context) {
		rawClaims, exists := c.Get("claims")
		if !exists {
			abortWithError(c, http.StatusForbidden, "access forbidden")
			return
		}

		principal, ok := rawClaims.(Principal)
		if !ok {
			abortWithError(c, http.StatusForbidden, "access forbidden")
			return
		}

		if !hasAnyRole(principal.GetRoles(), requiredRoles) {
			abortWithError(c, http.StatusForbidden, "access forbidden")
			return
		}

		c.Next()
	}
}

// abortWithError sends a JSON error response and aborts further processing.
func abortWithError(c httpserver.Context, statusCode int, message string) {
	c.JSON(statusCode, httpserver.H{"error": message})
	c.Abort()
}

// isTokenExpired checks if the provided error indicates an expired token.
func isTokenExpired(err error) bool {
	return errors.Is(err, ErrTokenExpired)
}

// ErrTokenExpired is returned when a token has expired.
var ErrTokenExpired = errors.New("token expired")

// hasAnyRole returns true if any of the requiredRoles is present in userRoles.
func hasAnyRole(userRoles, requiredRoles []string) bool {
	for _, role := range requiredRoles {
		if containsString(userRoles, role) {
			return true
		}
	}
	return false
}

// containsString checks if a slice contains the specified string.
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
