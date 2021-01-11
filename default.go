// Copyright 2021 by Cherry . All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author: Chamith Udayanga. (https://github.com/udayangaac)
package cherry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// defaultErrorEncoder executes when endpoint,response encoder or request decoder
// return a error. Returning the error response with HTTP status code 500 (Internal Server Error)
// Also a custom error encoder can be defined (See ~/config.go > globalConfig.WithErrorEncoder )
func defaultErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	if m, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := m.MarshalJSON(); marshalErr == nil {
			contentType, body = "application/json; charset=utf-8", jsonBody
		}
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(body)
}

// warningEncoder generates a warning message for undefined response encoder.
// There must be at least one encoder, this encoder gives a waring message when encoder is not defined.
func warningEncoder(ctx context.Context, w http.ResponseWriter, data interface{}) (err error) {
	epError, _ := data.(string)
	contentType, body := "text/plain; charset=utf-8",
		[]byte(fmt.Sprintf("Define the encoders at least success encoder. %v", epError))
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(body)
	return
}

// getWarningEndpoint generates a warning message for undefined endpoints
// The missing for different HTTP method of defined route in
// error message can be find.
func getWarningEndpoint(route string) Endpoint {
	return func(ctx context.Context, authentication Authentication, req interface{}) (resp interface{}, err error) {
		resp = fmt.Sprintf("Missing endpoint/s for Route{%v}", route)
		return
	}
}

// NopDecoder not implement any operation. Can be used when the request does
// not need to be processed.
func NopDecoder(ctx context.Context, r *http.Request) (data interface{}, err error) {
	return nil, nil
}

// DefaultJSONEncoder encode the data to JSON format and returning with HTTP status code 200 (OK)
// It can be used as success response's encoder for JSON responses.
func DefaultJSONEncoder(ctx context.Context, w http.ResponseWriter, data interface{}) (err error) {
	bodyByteArr := make([]byte,0)
	bodyByteArr,err =  json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bodyByteArr)
	return
}
