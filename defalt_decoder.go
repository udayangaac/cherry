package cherry

import (
	"context"
	"net/http"
)

func GetDefaultDecoder(reqPtr interface{}) DecodeRequestFunc {
	return func(ctx context.Context, req *http.Request) (i interface{}, err error) {
		err = NewReader().Read(ctx, req, reqPtr)
		if err != nil {
			return
		}
		return reqPtr, nil
	}
}
