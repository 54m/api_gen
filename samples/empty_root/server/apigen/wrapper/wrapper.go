// Code generated by server_generator. DO NOT EDIT.
// generated version: devel
package wrapper

import (
	"fmt"

	"github.com/go-generalize/api_gen/samples/empty_root/server/apigen/wrapper/internal"
)

// APIError - API Error
type APIError struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
	err    error       `json:"-"`
}

// NewAPIError - constructor
func NewAPIError(status int, bodies ...interface{}) *APIError {
	ae := &APIError{Status: status}
	if length := len(bodies); length > 0 {
		if length == 1 {
			ae.Body = bodies[0]
		} else {
			ae.Body = bodies
		}
	}
	return ae
}

// Error - makes it compatible with `error` interface.
func (ae *APIError) Error() string {
	if ae.err == nil {
		return fmt.Sprintf("status=%d, body=%v", ae.Status, ae.Body)
	}
	return fmt.Sprintf("status=%d, body=%v, err=%+v", ae.Status, ae.Body, ae.err)
}

// SetError - sets error to APIError.err
func (ae *APIError) SetError(err error) *APIError {
	ae.err = internal.Errorf(": %w", err)
	return ae
}

// SetErrorf - sets error to APIError.err (use `%w` to wrap the error)
func (ae *APIError) SetErrorf(format string, a ...interface{}) *APIError {
	ae.err = internal.Errorf(format, a...)
	return ae
}
