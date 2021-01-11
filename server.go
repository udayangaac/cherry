// Copyright 2021 by Chamith Udayange. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package cherry

import (
	"log"
	"net/http"
	"os"
)

type Server struct {
	dec        DecodeRequestFunc
	enc        EncodeResponseFunc
	ep         Endpoint
	errEnc     ErrorEncoder
	errHandler ErrorHandler
	entryPoint EntryPoint
	rf         RequestFilter
	et         ExtractTokenFunc
	r          []string
	ir         InitializeRequestFunc
}

// Configuration specified to there
type HandlerConfig struct {
	Path     string
	Decoder  DecodeRequestFunc
	Encoder  EncodeResponseFunc
	Role     []string
	Endpoint Endpoint
}

func NewServer(config HandlerConfig) (path string, s *Server) {

	// Initialize the server with default configuration
	// For specific routes.
	ser := Server{
		dec:        NopDecoder,
		enc:        warningEncoder,
		ep:         getWarningEndpoint(config.Path),
		errEnc:     defaultErrorEncoder,
		errHandler: NewLogErrorHandler(log.New(os.Stderr, "", log.LstdFlags)),
		entryPoint: NewDefaultEntryPoint(),
		et:         getExtractTokenFunc(),
		r:          config.Role,
	}

	// Assigning custom configurations
	// Specific for the route.
	if config.Decoder != nil {
		ser.dec = config.Decoder
	}
	if config.Encoder != nil {
		ser.enc = config.Encoder
	}
	if config.Endpoint != nil {
		ser.ep = config.Endpoint
	}

	// Assigning custom configurations
	// Global (For all routes).
	if gConf.entryPoint != nil {
		ser.entryPoint = gConf.entryPoint
	}
	if gConf.enc != nil {
		ser.enc = gConf.enc
	}
	if gConf.ir != nil {
		ser.ir = gConf.ir
	}

	return config.Path, &ser
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		token string
	)

	ctx := r.Context()

	if s.ir != nil {
		ctx = s.ir(ctx, r)
	}

	epStatus := EntryPointStatus{}
	auth := Authentication{}

	if s.r == nil {
		goto startRP
	}

	token, err = s.et(ctx, r)
	if err != nil {
		epStatus = EntryPointStatus{
			Code: InvalidToken,
			Desc: "Invalid or empty token",
		}
	}
	if s.rf != nil && err == nil {
		auth, epStatus = s.rf.DoFilter(ctx, token, s.r)
	}

	err = s.entryPoint.Commence(ctx, epStatus, w)
	if err != nil {
		s.errHandler.Handle(ctx, err)
		return
	}

startRP:
	request, err := s.dec(ctx, r)
	if err != nil {
		s.errEnc(ctx, err, w)
		s.errHandler.Handle(ctx, err)
		return
	}

	response, err := s.ep(ctx, auth, request)
	if err != nil {
		s.errEnc(ctx, err, w)
		s.errHandler.Handle(ctx, err)
		return
	}

	if err := s.enc(ctx, w, response); err != nil {
		s.errEnc(ctx, err, w)
		s.errHandler.Handle(ctx, err)
		return
	}
}
