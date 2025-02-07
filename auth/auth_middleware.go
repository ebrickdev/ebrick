package auth

import (
	"net/http"

	"github.com/ebrickdev/ebrick/transport/httpserver"
)

func AuthMiddleware(authenticator Authenticator) httpserver.HandlerFunc {
	return func(ctx httpserver.Context) {
		user, err := authenticator.Authenticate(ctx.Request())
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, httpserver.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		// Attach the user to the context for downstream handlers.
		ctx.Set("user", user)
		ctx.Next()
	}
}
