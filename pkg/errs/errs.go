package errs

import "errors"

var ErrUnauthorized = errors.New("unauthorized")
var ErrNotFound = errors.New("not found")
var ErrValidation = errors.New("validation error")