package pmail

import "errors"

var (
	ErrInvalidEmail = errors.New("email is not valid (missing headers or body)")
)
