# Cherry
Go library Back-end Services for authentication and authorization capabilities.

## Structure
 ...
#### Example code for router

```gotemplate
package http

import (
	"context"
	"fmt"
    ...
)

var (
	User      = []string{"user"}
	Admin     = []string{"admin"}
	Anonymous []string
	server    *http.Server
)

func init() {
	Anonymous = make([]string, 0)
}
func requestInitializer(ctx context.Context, r *http.Request) (ctxR context.Context) {
	ctxWithTimeout, _ := context.WithTimeout(ctx, 10*time.Second)
	uuidStr := uuid.New().String()
	ctxR = context.WithValue(ctxWithTimeout, "uuid_str", uuidStr)
	log.Trace(logTraceable.GetMessage(ctxR, "Started to process request URL:", r.URL, "Method:", r.Method))
	return
}

func InitRoutes(port int) {

	// ----------------------------------------------------------
	// Router configurations
	cherry.GetGlobalConfig().
		WithEncoder(encoder.MainEncoder).
		WithErrorEncoder(encoder.ErrorEncoder).
		WithRequestFilter(cherry.NewJwtRequestFilter(cherry.JwtConfig{
			Secret:        config.ServerConf.Jwt.Key,
			ValidDuration: time.Duration(config.ServerConf.Jwt.Duration) * time.Minute,
		})).
		WithEntryPoint(entrypoint.NewCustomEntryPoint()).
		WithRequestInitializer(requestInitializer).
		WithErrorHandler(custom_error.NewCustomLogErrorHandler())
	router := mux.NewRouter()

	// ------------------------------------------------------------
	// Facebook authentications
	authV1Router := router.PathPrefix("/oauth2/v1").Subrouter()
	authV1Router.HandleFunc("/fb/login", handler.GetFacebookLoginHandler()).Methods(http.MethodGet)
	authV1Router.HandleFunc("/fb/callback", handler.GetFacebookCallbackHandler()).Methods(http.MethodGet)
	authV1Router.HandleFunc("/fb/authenticate", handler.GetFacebookAuthenticationHandler()).Methods(http.MethodGet)

	publicV1Router := router.PathPrefix("/public/v1").Subrouter()
	// ----------------------------------------------------------
	// Home page anonymous routes can be access anyone
	// Search advertisement endpoints
	publicV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/search",
		Decoder:  decoder.SearchAdvertisementsDecoder(),
		Endpoint: endpoint.SearchAdvertisementsEndpoint(),
	})).Methods(http.MethodPost)

	// View advertisement by id
	publicV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/advertisement/{id}",
		Decoder:  decoder.GetAdvertisementDecoder(),
		Endpoint: endpoint.ViewAdvertisementsEndpoint(),
	})).Methods(http.MethodGet)

	// Get countries
	publicV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/countries",
		Endpoint: endpoint.GetCountries(),
	})).Methods(http.MethodGet)

	// Get cities
	publicV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/cities",
		Endpoint: endpoint.GetCities(),
	})).Methods(http.MethodGet)

	// Get categories
	publicV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/categories",
		Endpoint: endpoint.GetCategories(),
	})).Methods(http.MethodGet)

	// -----------------------------------------------------------
	// Generic endpoints for common usages
	commonV1Router := router.PathPrefix("/common/v1").Subrouter()
	// Get countries
	commonV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/countries",
		Endpoint: endpoint.GetCountries(),
		Role:     User,
	})).Methods(http.MethodGet)
	// Get cities
	commonV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/cities",
		Endpoint: endpoint.GetCities(),
		Role:     User,
	})).Methods(http.MethodGet)
	// Get statuses
	commonV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/statuses",
		Endpoint: endpoint.GetStatus(),
		Role:     Anonymous,
	})).Methods(http.MethodGet)
	// Get categories
	commonV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/categories",
		Endpoint: endpoint.GetCategories(),
		Role:     User,
	})).Methods(http.MethodGet)

	// -----------------------------------------------------------
	// User related endpoints
	userV1Router := router.PathPrefix("/user/v1").Subrouter()
	// Get user advertisements for dashboard viewing
	userV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/advertisements",
		Decoder:  decoder.GetAdvertisementsDecoder(),
		Endpoint: endpoint.GetAdvertisementsEndpoint(),
		Role:     User,
	})).Methods(http.MethodGet)

	// Get user advertisement by user Id(validating token) and
	// advertisement id.
	userV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/advertisement/{id}",
		Decoder:  decoder.GetAdvertisementDecoder(),
		Endpoint: endpoint.GetUserAdvertisementEndpoint(),
		Role:     User,
	})).Methods(http.MethodGet)

	// Add user advertisement to the system
	userV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/advertisement",
		Decoder:  decoder.AddAdvertisementDecoder(),
		Endpoint: endpoint.AddUserAdvertisementEndpoint(),
		Role:     User,
	})).Methods(http.MethodPost)

	// Update user advertisement by the user
	userV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/advertisement",
		Decoder:  decoder.UpdateAdvertisementDecoder(),
		Endpoint: endpoint.UpdateUserAdvertisementEndpoint(),
		Role:     User,
	})).Methods(http.MethodPut)

	// Delete advertisement
	// Need to remove from the elasticSearch cluster
	// and updated the deleted at column.
	userV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/advertisement/{id}",
		Decoder:  decoder.DeleteAdvertisementDecoder(),
		Endpoint: endpoint.DeleteUserAdvertisementEndpoint(),
		Role:     User,
	})).Methods(http.MethodDelete)

	// -----------------------------------------------------------
	// Admin related endpoints
	adminV1Router := router.PathPrefix("/admin/v1").Subrouter()
	// Get admin advertisements for dashboard viewing
	adminV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/advertisements",
		Decoder:  decoder.GetAdvertisementsDecoder(),
		Endpoint: endpoint.GetAdminAdvertisementsEndpoint(),
		Role:     Admin,
	})).Methods(http.MethodGet)

	adminV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/advertisement/{id}",
		Decoder:  decoder.GetAdvertisementDecoder(),
		Endpoint: endpoint.GetAdminAdvertisementEndpoint(),
		Role:     Admin,
	})).Methods(http.MethodGet)

	// Change state of the advertisement
	// Status : approve / reject
	// TODO: please confirm
	adminV1Router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/advertisement/{id}/status/{status}",
		Decoder:  decoder.UpdateAdvertisementStatusDecoder(),
		Endpoint: endpoint.UpdateAdvertisementStatusEndpoint(),
		Role:     Admin,
	})).Methods(http.MethodPut)

	// ---------------------------------------------------
	// This endpoint is to check whether the port is
	// up and running.
	router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/ping",
		Endpoint: endpoint.GetNopEndpoint(),
	})).Methods(http.MethodGet)

	router.Handle(cherry.NewHandler(cherry.HandlerConfig{
		Path:     "/ping",
		Endpoint: endpoint.GetNopEndpoint(),
	})).Methods(http.MethodOptions)

	router.Use(mux.CORSMethodMiddleware(router))
    ...
}

```


#### Example code for entry point

```gotemplate
package entrypoint

import (
	"context"
	"github.com/udayangaac/cherry"
	"net/http"
    ...
)

type customEntryPoint struct{}

func NewCustomEntryPoint() cherry.EntryPoint {
	return new(customEntryPoint)
}

func (e customEntryPoint) Commence(ctx context.Context, entryPointStatus cherry.EntryPointStatus, w http.ResponseWriter) (err error) {
	switch entryPointStatus.Code {
	case cherry.InvalidToken:
		err = GetError(ctx, customError.ErrInvalidToken, w)
	case cherry.EmptyToken:
		err = GetError(ctx, customError.ErrBadRequest, w)
	case cherry.TokenExpired:
		err = GetError(ctx, customError.ErrTokenExpired, w)
	case cherry.Unauthorized:
		err = GetError(ctx, customError.ErrUnauthorizedRequest, w)
	default:
		return nil
	}
	return
}

func GetError(ctx context.Context, err error, w http.ResponseWriter) error {
	encoder.ErrorEncoder(ctx, err, w)
	return err
}

```
