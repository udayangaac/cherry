// Copyright 2021 by Chamith Udayanga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cherry

import (
	"context"
	"net/http"
)

// DecodeRequestFunc extracting data from HTTP request returning the extracted values
// from the request or returning the decode error.
type DecodeRequestFunc func(context.Context, *http.Request) (interface{}, error)

// EncodeResponseFunc response data is written to the http.ResponseWriter and if eny
// error occurred it must be return as a error.
type EncodeResponseFunc func(context.Context, http.ResponseWriter, interface{}) error

// End
type Endpoint func(context.Context, Authentication, interface{}) (interface{}, error)

type ErrorEncoder func(context.Context, error, http.ResponseWriter)

type ExtractTokenFunc func(ctx context.Context, r *http.Request) (token string, err error)

type InitializeRequestFunc func(ctx context.Context, r *http.Request) (ctxR context.Context)
