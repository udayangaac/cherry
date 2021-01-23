// Copyright 2021 by Chamith Udayanga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cherry

import (
	"context"
	"net/http"
)

// Entry point status can be used to
// create authentication and authorization error responses
type EntryPointStatus struct {
	Code int
	Desc string
}

const (
	InvalidToken = 1
	EmptyToken   = 2
	TokenExpired = 3
	Unauthorized = 4
)

type EntryPoint interface {
	// Commence can be used to create response based on status(EntryPointStatus).
	// If response need to be terminated, err must be return and response must be written to the
	// http.ResponseWriter.
	Commence(ctx context.Context, status EntryPointStatus, w http.ResponseWriter) (err error)
}
