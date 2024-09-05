// Package httpport provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package httpport

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /bots)
	CreateBot(w http.ResponseWriter, r *http.Request)

	// (GET /bots/{uuid})
	GetBot(w http.ResponseWriter, r *http.Request, uuid string)

	// (GET /bots/{uuid}/answers)
	GetAnswers(w http.ResponseWriter, r *http.Request, uuid string)

	// (POST /bots/{uuid}/start)
	StartBot(w http.ResponseWriter, r *http.Request, uuid string)

	// (POST /bots/{uuid}/stop)
	StopBot(w http.ResponseWriter, r *http.Request, uuid string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// CreateBot operation middleware
func (siw *ServerInterfaceWrapper) CreateBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateBot(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetBot operation middleware
func (siw *ServerInterfaceWrapper) GetBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "uuid" -------------
	var uuid string

	err = runtime.BindStyledParameterWithLocation("simple", false, "uuid", runtime.ParamLocationPath, chi.URLParam(r, "uuid"), &uuid)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "uuid", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetBot(w, r, uuid)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetAnswers operation middleware
func (siw *ServerInterfaceWrapper) GetAnswers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "uuid" -------------
	var uuid string

	err = runtime.BindStyledParameterWithLocation("simple", false, "uuid", runtime.ParamLocationPath, chi.URLParam(r, "uuid"), &uuid)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "uuid", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAnswers(w, r, uuid)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// StartBot operation middleware
func (siw *ServerInterfaceWrapper) StartBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "uuid" -------------
	var uuid string

	err = runtime.BindStyledParameterWithLocation("simple", false, "uuid", runtime.ParamLocationPath, chi.URLParam(r, "uuid"), &uuid)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "uuid", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.StartBot(w, r, uuid)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// StopBot operation middleware
func (siw *ServerInterfaceWrapper) StopBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "uuid" -------------
	var uuid string

	err = runtime.BindStyledParameterWithLocation("simple", false, "uuid", runtime.ParamLocationPath, chi.URLParam(r, "uuid"), &uuid)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "uuid", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.StopBot(w, r, uuid)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshallingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshallingParamError) Error() string {
	return fmt.Sprintf("Error unmarshalling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshallingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/bots", wrapper.CreateBot)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/bots/{uuid}", wrapper.GetBot)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/bots/{uuid}/answers", wrapper.GetAnswers)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/bots/{uuid}/start", wrapper.StartBot)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/bots/{uuid}/stop", wrapper.StopBot)
	})

	return r
}