// Copyright 2021 by Chamith Udayange. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package cherry

import (
	"context"
)

type RequestFilter interface {
	DoFilter(ctx context.Context, token string, roles []string) (authentication Authentication, status EntryPointStatus)
}
