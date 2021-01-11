// Copyright 2021 by Chamith Udayange. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package cherry

import (
	"context"
	"net/http"
)

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
	Commence(ctx context.Context, status EntryPointStatus, w http.ResponseWriter) (err error)
}
