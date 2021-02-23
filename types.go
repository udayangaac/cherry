// Copyright 2021 by Chamith Udayanga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package cherry

import (
	"context"
	"net/http"
)

type DecodeRequestFunc func(context.Context, *http.Request) (interface{}, error)

type Endpoint func(context.Context, Authentication, interface{}) (interface{}, error)

type EncodeResponseFunc func(context.Context, http.ResponseWriter, interface{}) error

type ErrorEncoder func(context.Context, error, http.ResponseWriter)

type ExtractTokenFunc func(ctx context.Context, r *http.Request) (token string, err error)

type InitializeRequestFunc func(ctx context.Context, r *http.Request) (ctxR context.Context)
