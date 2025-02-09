package auth

import (
	"context"
)

type User struct {
	ID       string
	Username string
	Password string
}

type AuthManager interface {
	Authenticate(ctx context.Context, token string) (Principal, error)
}
