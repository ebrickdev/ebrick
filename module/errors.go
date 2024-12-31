package module

import "errors"

var (
	ErrModulePathNotFound = errors.New("module path not found")
	ErrModuleNotFound     = errors.New("module not found")
	ErrInvalidModuleType  = errors.New("invalid module type")
)
