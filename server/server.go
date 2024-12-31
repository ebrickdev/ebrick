package server

var (
	DefaultAddress = ":0"
	DefaultName    = "go.ebrick.server"
	DefaultVersion = "latest"
)

// Server defines the interface for a server with methods to initialize, start, and stop the server,
// as well as retrieve its options.
type Server interface {
	// Initialize options
	Init(...Option) error
	// Retrieve the options
	Options() Options
	// Start the server
	Start() error
	// Stop the server
	Stop() error
}
type Router interface {
}
