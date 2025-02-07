// File: auth/core.go
package auth

import "net/http"

// User represents an authenticated user.
type User struct {
	ID       string
	Username string
	// Additional fields (roles, email, etc.) can be added here.
}

// Authenticator defines the core method for authenticating an HTTP request.
type Authenticator interface {
	Authenticate(r *http.Request) (*User, error)
}
