// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	stunnerv1 "github.com/l7mp/stunner/pkg/apis/v1"
	"github.com/oapi-codegen/runtime"
)

// V1Config Config provides a STUNner config. Schema is defined in https://github.com/l7mp/stunner/tree/main/pkg/apis/v1
type V1Config = stunnerv1.StunnerConfig

// ListV1ConfigsParams defines parameters for ListV1Configs.
type ListV1ConfigsParams struct {
	// Watch Watch for changes to the described resources and return them as a stream of add, update, and remove notifications.
	Watch *bool `form:"watch,omitempty" json:"watch,omitempty"`
}

// ListV1ConfigsNamespaceParams defines parameters for ListV1ConfigsNamespace.
type ListV1ConfigsNamespaceParams struct {
	// Watch Watch for changes to the described resources and return them as a stream of add, update, and remove notifications.
	Watch *bool `form:"watch,omitempty" json:"watch,omitempty"`
}

// GetV1ConfigNamespaceNameParams defines parameters for GetV1ConfigNamespaceName.
type GetV1ConfigNamespaceNameParams struct {
	// Watch Watch for changes to the described resources and return them as a stream of add, update, and remove notifications.
	Watch *bool `form:"watch,omitempty" json:"watch,omitempty"`
}

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
	// ListV1Configs request
	ListV1Configs(ctx context.Context, params *ListV1ConfigsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListV1ConfigsNamespace request
	ListV1ConfigsNamespace(ctx context.Context, namespace string, params *ListV1ConfigsNamespaceParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetV1ConfigNamespaceName request
	GetV1ConfigNamespaceName(ctx context.Context, namespace string, name string, params *GetV1ConfigNamespaceNameParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) ListV1Configs(ctx context.Context, params *ListV1ConfigsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListV1ConfigsRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListV1ConfigsNamespace(ctx context.Context, namespace string, params *ListV1ConfigsNamespaceParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListV1ConfigsNamespaceRequest(c.Server, namespace, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetV1ConfigNamespaceName(ctx context.Context, namespace string, name string, params *GetV1ConfigNamespaceNameParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetV1ConfigNamespaceNameRequest(c.Server, namespace, name, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewListV1ConfigsRequest generates requests for ListV1Configs
func NewListV1ConfigsRequest(server string, params *ListV1ConfigsParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/configs")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Watch != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "watch", runtime.ParamLocationQuery, *params.Watch); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewListV1ConfigsNamespaceRequest generates requests for ListV1ConfigsNamespace
func NewListV1ConfigsNamespaceRequest(server string, namespace string, params *ListV1ConfigsNamespaceParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespace", runtime.ParamLocationPath, namespace)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/configs/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Watch != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "watch", runtime.ParamLocationQuery, *params.Watch); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetV1ConfigNamespaceNameRequest generates requests for GetV1ConfigNamespaceName
func NewGetV1ConfigNamespaceNameRequest(server string, namespace string, name string, params *GetV1ConfigNamespaceNameParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespace", runtime.ParamLocationPath, namespace)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "name", runtime.ParamLocationPath, name)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/configs/%s/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Watch != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "watch", runtime.ParamLocationQuery, *params.Watch); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
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
	// ListV1ConfigsWithResponse request
	ListV1ConfigsWithResponse(ctx context.Context, params *ListV1ConfigsParams, reqEditors ...RequestEditorFn) (*ListV1ConfigsResponse, error)

	// ListV1ConfigsNamespaceWithResponse request
	ListV1ConfigsNamespaceWithResponse(ctx context.Context, namespace string, params *ListV1ConfigsNamespaceParams, reqEditors ...RequestEditorFn) (*ListV1ConfigsNamespaceResponse, error)

	// GetV1ConfigNamespaceNameWithResponse request
	GetV1ConfigNamespaceNameWithResponse(ctx context.Context, namespace string, name string, params *GetV1ConfigNamespaceNameParams, reqEditors ...RequestEditorFn) (*GetV1ConfigNamespaceNameResponse, error)
}

type ListV1ConfigsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r ListV1ConfigsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListV1ConfigsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListV1ConfigsNamespaceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r ListV1ConfigsNamespaceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListV1ConfigsNamespaceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetV1ConfigNamespaceNameResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *V1Config
}

// Status returns HTTPResponse.Status
func (r GetV1ConfigNamespaceNameResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetV1ConfigNamespaceNameResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ListV1ConfigsWithResponse request returning *ListV1ConfigsResponse
func (c *ClientWithResponses) ListV1ConfigsWithResponse(ctx context.Context, params *ListV1ConfigsParams, reqEditors ...RequestEditorFn) (*ListV1ConfigsResponse, error) {
	rsp, err := c.ListV1Configs(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListV1ConfigsResponse(rsp)
}

// ListV1ConfigsNamespaceWithResponse request returning *ListV1ConfigsNamespaceResponse
func (c *ClientWithResponses) ListV1ConfigsNamespaceWithResponse(ctx context.Context, namespace string, params *ListV1ConfigsNamespaceParams, reqEditors ...RequestEditorFn) (*ListV1ConfigsNamespaceResponse, error) {
	rsp, err := c.ListV1ConfigsNamespace(ctx, namespace, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListV1ConfigsNamespaceResponse(rsp)
}

// GetV1ConfigNamespaceNameWithResponse request returning *GetV1ConfigNamespaceNameResponse
func (c *ClientWithResponses) GetV1ConfigNamespaceNameWithResponse(ctx context.Context, namespace string, name string, params *GetV1ConfigNamespaceNameParams, reqEditors ...RequestEditorFn) (*GetV1ConfigNamespaceNameResponse, error) {
	rsp, err := c.GetV1ConfigNamespaceName(ctx, namespace, name, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetV1ConfigNamespaceNameResponse(rsp)
}

// ParseListV1ConfigsResponse parses an HTTP response from a ListV1ConfigsWithResponse call
func ParseListV1ConfigsResponse(rsp *http.Response) (*ListV1ConfigsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListV1ConfigsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseListV1ConfigsNamespaceResponse parses an HTTP response from a ListV1ConfigsNamespaceWithResponse call
func ParseListV1ConfigsNamespaceResponse(rsp *http.Response) (*ListV1ConfigsNamespaceResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListV1ConfigsNamespaceResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetV1ConfigNamespaceNameResponse parses an HTTP response from a GetV1ConfigNamespaceNameWithResponse call
func ParseGetV1ConfigNamespaceNameResponse(rsp *http.Response) (*GetV1ConfigNamespaceNameResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetV1ConfigNamespaceNameResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest V1Config
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case rsp.StatusCode == 200:
		// Content-type (application/json;stream=watch) unsupported

	}

	return response, nil
}
