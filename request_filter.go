// Copyright 2021 by Chamith Udayanga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cherry

import (
	"context"
)

// Validation logic of http request can be implement using this  interface.
type RequestFilter interface {
	// DoFilter, tokens, defined roles to the http request are parsed as parameters.
	// DoFilter function must returns Authentication and EntryPointStatus variables.
	DoFilter(ctx context.Context, token string, roles []string) (authentication Authentication, status EntryPointStatus)
}
