package cherry

import (
	"context"
	"net/http"
)

// Deprecated
func GetDefaultDecoder(reqPtr interface{}) DecodeRequestFunc {
	return func(ctx context.Context, req *http.Request) (i interface{}, err error) {
		err = NewRequestReader().Read(ctx, req, reqPtr)
		if err != nil {
			return
		}
		return reqPtr, nil
	}
}
