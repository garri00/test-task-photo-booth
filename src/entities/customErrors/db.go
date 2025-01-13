package customErrors

import "errors"

var ErrNoRowsFindToDelete = errors.New("no row found to delete")
