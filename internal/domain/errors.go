package domain

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrValid    = errors.New("validation error")
)
