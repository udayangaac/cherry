// Copyright 2021 by Chamith Udayanga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cherry

import (
	"log"
	"net/http"
	"os"
)

type Server struct {
	dec                   DecodeRequestFunc
	enc                   EncodeResponseFunc
	endpoint              Endpoint
	errEnc                ErrorEncoder
	errHandler            ErrorHandler
	entryPoint            EntryPoint
	requestFilter         RequestFilter
	extractTokenFunc      ExtractTokenFunc
	role                  []string
	initializeRequestFunc InitializeRequestFunc
}

// Configuration specified to there
type HandlerConfig struct {
	Path     string
	Decoder  DecodeRequestFunc
	Encoder  EncodeResponseFunc
	Role     []string
	Endpoint Endpoint
}

// NewServer creates a reference of the server instance with
// defined configurations.
func NewServer(config HandlerConfig) (path string, s *Server) {

	// Initialize the server with default configuration
	// For specific routes.
	ser := Server{
		dec:              NopDecoder,
		enc:              warningEncoder,
		endpoint:         getWarningEndpoint(config.Path),
		errEnc:           defaultErrorEncoder,
		errHandler:       NewLogErrorHandler(log.New(os.Stderr, "", log.LstdFlags)),
		entryPoint:       NewDefaultEntryPoint(),
		extractTokenFunc: getExtractTokenFunc(),
		role:             config.Role,
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
		ser.endpoint = config.Endpoint
	}

	// Assigning custom configurations
	// Global (For all routes).
	if gConf.entryPoint != nil {
		ser.entryPoint = gConf.entryPoint
	}
	if gConf.enc != nil {
		ser.enc = gConf.enc
	}

	if gConf.requestFilter != nil {
		ser.requestFilter = gConf.requestFilter
	}

	if gConf.initializeRequestFunc != nil {
		ser.initializeRequestFunc = gConf.initializeRequestFunc
	}

	if gConf.errHandler != nil {
		ser.errHandler = gConf.errHandler
	}

	return config.Path, &ser
}

// All defined function related the http request are executing inside
// this function. No need to implement custom ServeHTTP functions.
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		token string
	)

	ctx := r.Context()

	if s.initializeRequestFunc != nil {
		ctx = s.initializeRequestFunc(ctx, r)
	}

	epStatus := EntryPointStatus{}
	auth := Authentication{}

	if s.role == nil {
		goto requestProcessStart
	}

	token, err = s.extractTokenFunc(ctx, r)
	if err != nil {
		epStatus = EntryPointStatus{
			Code: InvalidToken,
			Desc: "Invalid or empty token",
		}
	}
	if s.requestFilter != nil && err == nil {
		auth, epStatus = s.requestFilter.DoFilter(ctx, token, s.role)
	}

	err = s.entryPoint.Commence(ctx, epStatus, w)
	if err != nil {
		s.errHandler.Handle(ctx, err)
		return
	}

requestProcessStart:
	request, err := s.dec(ctx, r)
	if err != nil {
		s.errEnc(ctx, err, w)
		s.errHandler.Handle(ctx, err)
		return
	}

	response, err := s.endpoint(ctx, auth, request)
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
