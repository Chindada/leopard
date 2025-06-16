package client

import "errors"

var (
	ErrUIDRequired      = errors.New("uid is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrHostRequired     = errors.New("host is required")
	ErrPortRequired     = errors.New("port is required")
)
