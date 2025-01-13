package customErrors

import "errors"

var (
	ErrorBodyCloseFailed = errors.New("Body.Close() failed")
)
