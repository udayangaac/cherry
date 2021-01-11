// Copyright 2021 by Chamith Udayange. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cherry

// Authentication The information related to the authentication process.
// It can be get from the endpoint. (See ~/types.go > Endpoint)
type Authentication struct {
	ID    int         `json:"id"`
	Roles []string    `json:"roles"`
	Data  interface{} `json:"data"`
}
