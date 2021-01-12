// Copyright 2021 by Chamith Udayanga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cherry

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// getExtractTokenFunc() returns ExtractTokenFunc. extracting ('Bearer <token>') token from
// the header of the HTTP request where the key is 'Authorization'.
func getExtractTokenFunc() ExtractTokenFunc {
	return func(ctx context.Context, r *http.Request) (token string, err error) {
		value := r.Header.Get("Authorization")
		tokenSpl := strings.Split(value, "Bearer ")
		if len(tokenSpl) < 2 {
			err = errors.New("invalid token from header")
			return
		}
		token = tokenSpl[1]
		return
	}
}
