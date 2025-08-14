// Package customerrrors defines custom error types used across the application.
package customerrrors

import "errors"

var ErrUnexpected = errors.New("unexpected error occurred")
