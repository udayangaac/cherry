// Copyright 2021 by Chamith Udayange. All rights reserved.
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

type LogErrorHandler struct {
	logger *log.Logger
}

func NewLogErrorHandler(logger *log.Logger) *LogErrorHandler {
	return &LogErrorHandler{
		logger: logger,
	}
}

func (h *LogErrorHandler) Handle(ctx context.Context, err error) {
	h.logger.Println(err)
}
