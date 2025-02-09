package auth

// Principal defines the interface for authenticated user information.
type Principal interface {
	GetID() string
	GetEmail() string
	GetRoles() []string
	GetClaims() map[string]interface{}
}
