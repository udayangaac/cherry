package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/udayangaac/cherry"
	"log"
	"net/http"
	"time"
)

func main() {

	router := mux.NewRouter()

	jwtConfig := cherry.JwtConfig{
		Secret:        "1qaz2wsx!",
		ValidDuration: 60 * time.Second,
	}

	cherry.GetGlobalConfig().
		WithEncoder(getEncoder()).
		WithRequestFilter(cherry.NewJwtRequestFilter(jwtConfig))

	router.Handle(cherry.NewServer(cherry.HandlerConfig{
		Path:     "/authenticate",
		Endpoint: authenticateEndpoint(jwtConfig),
	})).Methods(http.MethodGet)

	router.Handle(cherry.NewServer(cherry.HandlerConfig{
		Path:     "/hello",
		Role:     []string{},
		Endpoint: getHelloEndpoint(),
	})).Methods(http.MethodGet)

	log.Print(http.ListenAndServe(":8082", router))
}

type SampleStruct struct {
	Sample string
}

func authenticateEndpoint(config cherry.JwtConfig) cherry.Endpoint {
	return func(ctx context.Context, authentication cherry.Authentication, request interface{}) (response interface{}, err error) {
		response, err = cherry.GenerateJwtToke(cherry.Authentication{
			ID:    1,
			Roles: []string{"user"},
			Data:  SampleStruct{Sample: "Sample Data"},
		}, config)
		return
	}
}

func getHelloEndpoint() cherry.Endpoint {
	return func(ctx context.Context, authentication cherry.Authentication, request interface{}) (response interface{}, err error) {
		return "Hello world", nil
	}
}


func getEncoder() cherry.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, data interface{}) (err error) {
		dataStr, _ := data.(string)
		contentType, body := "text/plain; charset=utf-8", []byte(dataStr)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", contentType)
		_, err = w.Write(body)
		return
	}
}
