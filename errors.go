package pmail

import "errors"

var (
	ErrInvalidEmail  = errors.New("email is not valid (missing headers or body)")
	ErrPartHasNoBody = errors.New("email part has no body (or is already consumed)")
)
