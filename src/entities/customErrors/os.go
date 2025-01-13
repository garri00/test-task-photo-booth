package customErrors

import "errors"

var (
	ErrorOsCloseFailed = errors.New("os.Close failed")
)
