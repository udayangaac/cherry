// Copyright 2021 by Chamith Udayanga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cherry

import (
	"context"
	"log"
)

type ErrorHandler interface {
	Handle(ctx context.Context, err error)
}

type logErrorHandler struct {
	logger *log.Logger
}

func NewLogErrorHandler(logger *log.Logger) ErrorHandler {
	return &logErrorHandler{
		logger: logger,
	}
}

func (h *logErrorHandler) Handle(ctx context.Context, err error) {
	h.logger.Println(err)
}
