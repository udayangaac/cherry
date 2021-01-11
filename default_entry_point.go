package cherry

import (
	"context"
	"errors"
	"net/http"
)

type defaultEntryPoint struct{}

func NewDefaultEntryPoint() EntryPoint {
	return new(defaultEntryPoint)
}

func (e defaultEntryPoint) Commence(ctx context.Context, status EntryPointStatus, w http.ResponseWriter) (err error) {

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

	// Build the auth error response
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	w.Write(body)
	return
}