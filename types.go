// Copyright 2021 by Chamith Udayanga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cherry

import (
	"context"
	"net/http"
)

// DecodeRequestFunc extracts data from HTTP request returning the extracted
// values or returning decode error.
type DecodeRequestFunc func(context.Context, *http.Request) (interface{}, error)

// Endpoint Authentication, response can be captured and success response must be
// return.
type Endpoint func(context.Context, Authentication, interface{}) (interface{}, error)

// EncodeResponseFunc encodes the data returning from the Endpoint.
type EncodeResponseFunc func(context.Context, http.ResponseWriter, interface{}) error

// ErrorEncoder encodes the error responses.
type ErrorEncoder func(context.Context, error, http.ResponseWriter)

// ExtractTokenFunc extracts the token from the request and must be return as a string
type ExtractTokenFunc func(ctx context.Context, r *http.Request) (token string, err error)

// InitializeRequestFunc, The context of the function stack for a certain request can be
// changed here.
type InitializeRequestFunc func(ctx context.Context, r *http.Request) (ctxR context.Context)
