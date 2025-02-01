// Package bots provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.2.0 DO NOT EDIT.
package bots

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetBots request
	GetBots(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateBotWithBody request with any body
	CreateBotWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateBot(ctx context.Context, body CreateBotJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteBot request
	DeleteBot(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetBot request
	GetBot(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetAnswers request
	GetAnswers(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateMailingWithBody request with any body
	CreateMailingWithBody(ctx context.Context, uuid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateMailing(ctx context.Context, uuid string, body CreateMailingJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// StartMailing request
	StartMailing(ctx context.Context, uuid string, entryKey string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// StartBot request
	StartBot(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// StopBot request
	StopBot(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetBots(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetBotsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateBotWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateBotRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateBot(ctx context.Context, body CreateBotJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateBotRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteBot(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteBotRequest(c.Server, uuid)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetBot(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetBotRequest(c.Server, uuid)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetAnswers(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAnswersRequest(c.Server, uuid)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateMailingWithBody(ctx context.Context, uuid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateMailingRequestWithBody(c.Server, uuid, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateMailing(ctx context.Context, uuid string, body CreateMailingJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateMailingRequest(c.Server, uuid, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) StartMailing(ctx context.Context, uuid string, entryKey string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewStartMailingRequest(c.Server, uuid, entryKey)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) StartBot(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewStartBotRequest(c.Server, uuid)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) StopBot(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewStopBotRequest(c.Server, uuid)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetBotsRequest generates requests for GetBots
func NewGetBotsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/bots")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateBotRequest calls the generic CreateBot builder with application/json body
func NewCreateBotRequest(server string, body CreateBotJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateBotRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateBotRequestWithBody generates requests for CreateBot with any type of body
func NewCreateBotRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/bots")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteBotRequest generates requests for DeleteBot
func NewDeleteBotRequest(server string, uuid string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "uuid", runtime.ParamLocationPath, uuid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/bots/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetBotRequest generates requests for GetBot
func NewGetBotRequest(server string, uuid string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "uuid", runtime.ParamLocationPath, uuid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/bots/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetAnswersRequest generates requests for GetAnswers
func NewGetAnswersRequest(server string, uuid string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "uuid", runtime.ParamLocationPath, uuid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/bots/%s/answers", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateMailingRequest calls the generic CreateMailing builder with application/json body
func NewCreateMailingRequest(server string, uuid string, body CreateMailingJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateMailingRequestWithBody(server, uuid, "application/json", bodyReader)
}

// NewCreateMailingRequestWithBody generates requests for CreateMailing with any type of body
func NewCreateMailingRequestWithBody(server string, uuid string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "uuid", runtime.ParamLocationPath, uuid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/bots/%s/mailings", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewStartMailingRequest generates requests for StartMailing
func NewStartMailingRequest(server string, uuid string, entryKey string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "uuid", runtime.ParamLocationPath, uuid)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "entryKey", runtime.ParamLocationPath, entryKey)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/bots/%s/mailings/%s/start", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewStartBotRequest generates requests for StartBot
func NewStartBotRequest(server string, uuid string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "uuid", runtime.ParamLocationPath, uuid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/bots/%s/start", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewStopBotRequest generates requests for StopBot
func NewStopBotRequest(server string, uuid string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "uuid", runtime.ParamLocationPath, uuid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/bots/%s/stop", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetBotsWithResponse request
	GetBotsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetBotsResponse, error)

	// CreateBotWithBodyWithResponse request with any body
	CreateBotWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateBotResponse, error)

	CreateBotWithResponse(ctx context.Context, body CreateBotJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateBotResponse, error)

	// DeleteBotWithResponse request
	DeleteBotWithResponse(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*DeleteBotResponse, error)

	// GetBotWithResponse request
	GetBotWithResponse(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*GetBotResponse, error)

	// GetAnswersWithResponse request
	GetAnswersWithResponse(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*GetAnswersResponse, error)

	// CreateMailingWithBodyWithResponse request with any body
	CreateMailingWithBodyWithResponse(ctx context.Context, uuid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateMailingResponse, error)

	CreateMailingWithResponse(ctx context.Context, uuid string, body CreateMailingJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateMailingResponse, error)

	// StartMailingWithResponse request
	StartMailingWithResponse(ctx context.Context, uuid string, entryKey string, reqEditors ...RequestEditorFn) (*StartMailingResponse, error)

	// StartBotWithResponse request
	StartBotWithResponse(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*StartBotResponse, error)

	// StopBotWithResponse request
	StopBotWithResponse(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*StopBotResponse, error)
}

type GetBotsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GetBots
	JSON400      *Error
	JSON401      *Error
}

// Status returns HTTPResponse.Status
func (r GetBotsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetBotsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateBotResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON401      *Error
}

// Status returns HTTPResponse.Status
func (r CreateBotResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateBotResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteBotResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON401      *Error
	JSON403      *Error
	JSON404      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteBotResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteBotResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetBotResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Bot
	JSON400      *Error
	JSON401      *Error
	JSON403      *Error
	JSON404      *Error
}

// Status returns HTTPResponse.Status
func (r GetBotResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetBotResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetAnswersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON401      *Error
	JSON403      *Error
	JSON404      *Error
}

// Status returns HTTPResponse.Status
func (r GetAnswersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAnswersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateMailingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON401      *Error
	JSON403      *Error
	JSON404      *Error
}

// Status returns HTTPResponse.Status
func (r CreateMailingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateMailingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type StartMailingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON401      *Error
	JSON403      *Error
	JSON404      *Error
}

// Status returns HTTPResponse.Status
func (r StartMailingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r StartMailingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type StartBotResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON401      *Error
	JSON403      *Error
	JSON404      *Error
}

// Status returns HTTPResponse.Status
func (r StartBotResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r StartBotResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type StopBotResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON401      *Error
	JSON403      *Error
	JSON404      *Error
}

// Status returns HTTPResponse.Status
func (r StopBotResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r StopBotResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetBotsWithResponse request returning *GetBotsResponse
func (c *ClientWithResponses) GetBotsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetBotsResponse, error) {
	rsp, err := c.GetBots(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetBotsResponse(rsp)
}

// CreateBotWithBodyWithResponse request with arbitrary body returning *CreateBotResponse
func (c *ClientWithResponses) CreateBotWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateBotResponse, error) {
	rsp, err := c.CreateBotWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateBotResponse(rsp)
}

func (c *ClientWithResponses) CreateBotWithResponse(ctx context.Context, body CreateBotJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateBotResponse, error) {
	rsp, err := c.CreateBot(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateBotResponse(rsp)
}

// DeleteBotWithResponse request returning *DeleteBotResponse
func (c *ClientWithResponses) DeleteBotWithResponse(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*DeleteBotResponse, error) {
	rsp, err := c.DeleteBot(ctx, uuid, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteBotResponse(rsp)
}

// GetBotWithResponse request returning *GetBotResponse
func (c *ClientWithResponses) GetBotWithResponse(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*GetBotResponse, error) {
	rsp, err := c.GetBot(ctx, uuid, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetBotResponse(rsp)
}

// GetAnswersWithResponse request returning *GetAnswersResponse
func (c *ClientWithResponses) GetAnswersWithResponse(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*GetAnswersResponse, error) {
	rsp, err := c.GetAnswers(ctx, uuid, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAnswersResponse(rsp)
}

// CreateMailingWithBodyWithResponse request with arbitrary body returning *CreateMailingResponse
func (c *ClientWithResponses) CreateMailingWithBodyWithResponse(ctx context.Context, uuid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateMailingResponse, error) {
	rsp, err := c.CreateMailingWithBody(ctx, uuid, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateMailingResponse(rsp)
}

func (c *ClientWithResponses) CreateMailingWithResponse(ctx context.Context, uuid string, body CreateMailingJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateMailingResponse, error) {
	rsp, err := c.CreateMailing(ctx, uuid, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateMailingResponse(rsp)
}

// StartMailingWithResponse request returning *StartMailingResponse
func (c *ClientWithResponses) StartMailingWithResponse(ctx context.Context, uuid string, entryKey string, reqEditors ...RequestEditorFn) (*StartMailingResponse, error) {
	rsp, err := c.StartMailing(ctx, uuid, entryKey, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseStartMailingResponse(rsp)
}

// StartBotWithResponse request returning *StartBotResponse
func (c *ClientWithResponses) StartBotWithResponse(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*StartBotResponse, error) {
	rsp, err := c.StartBot(ctx, uuid, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseStartBotResponse(rsp)
}

// StopBotWithResponse request returning *StopBotResponse
func (c *ClientWithResponses) StopBotWithResponse(ctx context.Context, uuid string, reqEditors ...RequestEditorFn) (*StopBotResponse, error) {
	rsp, err := c.StopBot(ctx, uuid, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseStopBotResponse(rsp)
}

// ParseGetBotsResponse parses an HTTP response from a GetBotsWithResponse call
func ParseGetBotsResponse(rsp *http.Response) (*GetBotsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetBotsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GetBots
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	}

	return response, nil
}

// ParseCreateBotResponse parses an HTTP response from a CreateBotWithResponse call
func ParseCreateBotResponse(rsp *http.Response) (*CreateBotResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateBotResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	}

	return response, nil
}

// ParseDeleteBotResponse parses an HTTP response from a DeleteBotWithResponse call
func ParseDeleteBotResponse(rsp *http.Response) (*DeleteBotResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteBotResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseGetBotResponse parses an HTTP response from a GetBotWithResponse call
func ParseGetBotResponse(rsp *http.Response) (*GetBotResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetBotResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Bot
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseGetAnswersResponse parses an HTTP response from a GetAnswersWithResponse call
func ParseGetAnswersResponse(rsp *http.Response) (*GetAnswersResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAnswersResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseCreateMailingResponse parses an HTTP response from a CreateMailingWithResponse call
func ParseCreateMailingResponse(rsp *http.Response) (*CreateMailingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateMailingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseStartMailingResponse parses an HTTP response from a StartMailingWithResponse call
func ParseStartMailingResponse(rsp *http.Response) (*StartMailingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &StartMailingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseStartBotResponse parses an HTTP response from a StartBotWithResponse call
func ParseStartBotResponse(rsp *http.Response) (*StartBotResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &StartBotResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseStopBotResponse parses an HTTP response from a StopBotWithResponse call
func ParseStopBotResponse(rsp *http.Response) (*StopBotResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &StopBotResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}
