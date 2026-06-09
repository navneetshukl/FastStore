package server

import (
	"errors"
)

var (
	invalid_command error = errors.New("invalid command")
	server_error    error = errors.New("server error")
)
