package cherry

import (
	"context"
	"errors"
	"net/http"
)

type defaultEntryPoint struct{}

// Create and defaultEntryPoint instance for EntryPoint interface
// Can be used for JWT authentication.
func NewDefaultEntryPoint() EntryPoint {
	return new(defaultEntryPoint)
}


// Commence is the function where the response generate for authentication, authorization errors
// If there is any error when the authentication, function need to returns the error.
// 'status' value is the same value that is return by request_filter.go > RequestFilter.DoFilter.
func (e defaultEntryPoint) Commence(ctx context.Context, EntryPointStatus EntryPointStatus, w http.ResponseWriter) (err error) {

	statusCode := http.StatusOK
	switch status.Code {
	case InvalidToken:
		statusCode = http.StatusBadRequest
		err = errors.New("invalid token")
	case EmptyToken:
		statusCode = http.StatusBadRequest
		err = errors.New("bad request")
	case TokenExpired:
		statusCode = http.StatusBadRequest
		err = errors.New("token expired")
	case Unauthorized:
		statusCode = http.StatusUnauthorized
		err = errors.New("unauthorised request")
	default:
		return nil
	}

	// Build the authentication/authorization error response
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	w.Write(body)
	return
}
