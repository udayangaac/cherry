// Copyright 2021 by Chamith Udayanga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cherry

// Default settings initialization
type globalConfig struct {
	dec                   DecodeRequestFunc
	enc                   EncodeResponseFunc
	errEnc                ErrorEncoder
	errHandler            ErrorHandler
	entryPoint            EntryPoint
	requestFilter         RequestFilter
	initializeRequestFunc InitializeRequestFunc
}

var gConf *globalConfig

func init() {
	gConf = &globalConfig{}
}

func GetGlobalConfig() *globalConfig {
	return gConf
}

func (g *globalConfig) WithEncoder(encoder EncodeResponseFunc) *globalConfig {
	g.enc = encoder
	return g
}

func (g *globalConfig) WithErrorEncoder(encoder ErrorEncoder) *globalConfig {
	g.errEnc = encoder
	return g
}

func (g *globalConfig) WithErrorHandler(errorHandler ErrorHandler) *globalConfig {
	g.errHandler = errorHandler
	return g
}

func (g *globalConfig) WithEntryPoint(entryPoint EntryPoint) *globalConfig {
	g.entryPoint = entryPoint
	return g
}

func (g *globalConfig) WithRequestFilter(requestFilter RequestFilter) *globalConfig {
	g.requestFilter = requestFilter
	return g
}

func (g *globalConfig) WithRequestInitializer(initializer InitializeRequestFunc) *globalConfig {
	g.initializeRequestFunc = initializer
	return g
}
