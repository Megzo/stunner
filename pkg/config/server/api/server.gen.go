// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	stunnerv1 "github.com/l7mp/stunner/pkg/apis/v1"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// V1Config Config provides a STUNner config. Schema is defined in https://github.com/l7mp/stunner/tree/main/pkg/apis/v1
type V1Config = stunnerv1.StunnerConfig

// V1ConfigList ConfigList is a list of Configs.
type V1ConfigList struct {
	// Items Items is the list of Config objects in the list.
	Items []V1Config `json:"items"`

	// Version version defines the versioned schema of this object.
	Version string `json:"version"`
}

// V1Error API error.
type V1Error struct {
	// Code Error code.
	Code int32 `json:"code"`

	// Message Error message.
	Message *string `json:"message,omitempty"`
}

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

	// Node Name of the node the client runs on.
	Node *string `form:"node,omitempty" json:"node,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /api/v1/configs)
	ListV1Configs(w http.ResponseWriter, r *http.Request, params ListV1ConfigsParams)

	// (GET /api/v1/configs/{namespace})
	ListV1ConfigsNamespace(w http.ResponseWriter, r *http.Request, namespace string, params ListV1ConfigsNamespaceParams)

	// (GET /api/v1/configs/{namespace}/{name})
	GetV1ConfigNamespaceName(w http.ResponseWriter, r *http.Request, namespace string, name string, params GetV1ConfigNamespaceNameParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// ListV1Configs operation middleware
func (siw *ServerInterfaceWrapper) ListV1Configs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params ListV1ConfigsParams

	// ------------- Optional query parameter "watch" -------------

	err = runtime.BindQueryParameter("form", true, false, "watch", r.URL.Query(), &params.Watch)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "watch", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListV1Configs(w, r, params)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// ListV1ConfigsNamespace operation middleware
func (siw *ServerInterfaceWrapper) ListV1ConfigsNamespace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "namespace" -------------
	var namespace string

	err = runtime.BindStyledParameter("simple", false, "namespace", mux.Vars(r)["namespace"], &namespace)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "namespace", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params ListV1ConfigsNamespaceParams

	// ------------- Optional query parameter "watch" -------------

	err = runtime.BindQueryParameter("form", true, false, "watch", r.URL.Query(), &params.Watch)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "watch", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListV1ConfigsNamespace(w, r, namespace, params)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetV1ConfigNamespaceName operation middleware
func (siw *ServerInterfaceWrapper) GetV1ConfigNamespaceName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "namespace" -------------
	var namespace string

	err = runtime.BindStyledParameter("simple", false, "namespace", mux.Vars(r)["namespace"], &namespace)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "namespace", Err: err})
		return
	}

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameter("simple", false, "name", mux.Vars(r)["name"], &name)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetV1ConfigNamespaceNameParams

	// ------------- Optional query parameter "watch" -------------

	err = runtime.BindQueryParameter("form", true, false, "watch", r.URL.Query(), &params.Watch)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "watch", Err: err})
		return
	}

	// ------------- Optional query parameter "node" -------------

	err = runtime.BindQueryParameter("form", true, false, "node", r.URL.Query(), &params.Node)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "node", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetV1ConfigNamespaceName(w, r, namespace, name, params)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
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

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
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
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = mux.NewRouter()
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

	r.HandleFunc(options.BaseURL+"/api/v1/configs", wrapper.ListV1Configs).Methods("GET")

	r.HandleFunc(options.BaseURL+"/api/v1/configs/{namespace}", wrapper.ListV1ConfigsNamespace).Methods("GET")

	r.HandleFunc(options.BaseURL+"/api/v1/configs/{namespace}/{name}", wrapper.GetV1ConfigNamespaceName).Methods("GET")

	return r
}

type ListV1ConfigsRequestObject struct {
	Params ListV1ConfigsParams
}

type ListV1ConfigsResponseObject interface {
	VisitListV1ConfigsResponse(w http.ResponseWriter) error
}

type ListV1Configs200JSONResponse V1ConfigList

func (response ListV1Configs200JSONResponse) VisitListV1ConfigsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ListV1Configs500JSONResponse V1Error

func (response ListV1Configs500JSONResponse) VisitListV1ConfigsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type ListV1ConfigsNamespaceRequestObject struct {
	Namespace string `json:"namespace"`
	Params    ListV1ConfigsNamespaceParams
}

type ListV1ConfigsNamespaceResponseObject interface {
	VisitListV1ConfigsNamespaceResponse(w http.ResponseWriter) error
}

type ListV1ConfigsNamespace200JSONResponse V1ConfigList

func (response ListV1ConfigsNamespace200JSONResponse) VisitListV1ConfigsNamespaceResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ListV1ConfigsNamespace400JSONResponse V1Error

func (response ListV1ConfigsNamespace400JSONResponse) VisitListV1ConfigsNamespaceResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type ListV1ConfigsNamespace500JSONResponse V1Error

func (response ListV1ConfigsNamespace500JSONResponse) VisitListV1ConfigsNamespaceResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetV1ConfigNamespaceNameRequestObject struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Params    GetV1ConfigNamespaceNameParams
}

type GetV1ConfigNamespaceNameResponseObject interface {
	VisitGetV1ConfigNamespaceNameResponse(w http.ResponseWriter) error
}

type GetV1ConfigNamespaceName200JSONResponse V1Config

func (response GetV1ConfigNamespaceName200JSONResponse) VisitGetV1ConfigNamespaceNameResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetV1ConfigNamespaceName400JSONResponse V1Error

func (response GetV1ConfigNamespaceName400JSONResponse) VisitGetV1ConfigNamespaceNameResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetV1ConfigNamespaceName500JSONResponse V1Error

func (response GetV1ConfigNamespaceName500JSONResponse) VisitGetV1ConfigNamespaceNameResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (GET /api/v1/configs)
	ListV1Configs(ctx context.Context, request ListV1ConfigsRequestObject) (ListV1ConfigsResponseObject, error)

	// (GET /api/v1/configs/{namespace})
	ListV1ConfigsNamespace(ctx context.Context, request ListV1ConfigsNamespaceRequestObject) (ListV1ConfigsNamespaceResponseObject, error)

	// (GET /api/v1/configs/{namespace}/{name})
	GetV1ConfigNamespaceName(ctx context.Context, request GetV1ConfigNamespaceNameRequestObject) (GetV1ConfigNamespaceNameResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHttpHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHttpMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// ListV1Configs operation middleware
func (sh *strictHandler) ListV1Configs(w http.ResponseWriter, r *http.Request, params ListV1ConfigsParams) {
	var request ListV1ConfigsRequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ListV1Configs(ctx, request.(ListV1ConfigsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListV1Configs")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListV1ConfigsResponseObject); ok {
		if err := validResponse.VisitListV1ConfigsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// ListV1ConfigsNamespace operation middleware
func (sh *strictHandler) ListV1ConfigsNamespace(w http.ResponseWriter, r *http.Request, namespace string, params ListV1ConfigsNamespaceParams) {
	var request ListV1ConfigsNamespaceRequestObject

	request.Namespace = namespace
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ListV1ConfigsNamespace(ctx, request.(ListV1ConfigsNamespaceRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListV1ConfigsNamespace")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListV1ConfigsNamespaceResponseObject); ok {
		if err := validResponse.VisitListV1ConfigsNamespaceResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetV1ConfigNamespaceName operation middleware
func (sh *strictHandler) GetV1ConfigNamespaceName(w http.ResponseWriter, r *http.Request, namespace string, name string, params GetV1ConfigNamespaceNameParams) {
	var request GetV1ConfigNamespaceNameRequestObject

	request.Namespace = namespace
	request.Name = name
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetV1ConfigNamespaceName(ctx, request.(GetV1ConfigNamespaceNameRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetV1ConfigNamespaceName")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetV1ConfigNamespaceNameResponseObject); ok {
		if err := validResponse.VisitGetV1ConfigNamespaceNameResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xXUW/bNhD+KwS3R0eKmw0D9LSuKIagQzYk7vpQ5IGmzhI7iWSPJ6dBoP8+HElLbuw0",
	"BVIUKNo3ijzefXf3fWf6TmrXe2fBUpDVnQy6hV7F5XZZvHB2Yxr+qCFoNJ6Ms7KSaV94dFtTQxBKXK1e",
	"X1hAoeNJIa6iH2GCqGFjLNTCWNES+VCVZWOoHdaFdn3Z/db7MtBgLWBJCFD2ytjS/9eUyptQbpdyIT+c",
	"NO6Ebj3ISmbb7bK4SquMcc/qxPTeITFsq/qPLsmF9IpaWcmHMOxHHsfFXIW/TKCHKsFnnKsSHa/cRqT9",
	"UHBAdB6QDMSqGoI+HPo55212QS3ccyLc+h1oClzB3Sn7nTz9jLCRlfypnFtZ5j6WcxPHhcwlVIjqlr+3",
	"gCGGv48mH+TeJVB5D2qRfDNAak3I8BjR1CI0tonVQ3g/GIRaVm+naDvk19OF5EGmar9EdHgI6fk/5wL4",
	"6LCk2tVweCH6EXzGNzYOe0WyksbS2bMZq7EEDSDH7iEE1TzoKR8/nmeEc5gcmxm7cYf+Vy1M+qkVKd8p",
	"CyIm/MG7EPWFkJQ1oFp3IPqhI3Pi0ZHTrsufTAxgJ6vXlxciAG6NBrFxKN7A+nL1QvRQGyWMbSBw6JiK",
	"oY5RXr68WsWQbH4IRu9ENlFGLovT4pTr5jxY5Y2s5FlxWpxlicXOsJLK7bJM1+NWA3S8ueQy71HcKNJt",
	"jrnPftV1ghUdvNIQtcU0UOzkvJaVZBn+u8zSizhQ9UCAQVZv78d8E4NwurpVXBJGwERPZmuoBUJwA2pu",
	"gOUvGjBKsBeKexIIQfWsA1XXCzH4WhEssm3vtiCsI7MxOiKMeA1Hfj8A3srFbjjFbOUij14uTmbO2rkO",
	"lJXjeM0MC97ZkBj/7PQ0Ed8S2FhP5X2XA5XvQtL07PCzJkQccJGmHxfq71fc5V+/bMik8iPRzi0BWtVF",
	"/gLuND9G03t8Ku8mNoxP5pYIHjR3a6bYIwy72Nk9RrXJMM1MEI0iuFG3kX03rZkBmSB4lEAgqCe+xF+s",
	"iS52L+o8dggHOEKheUT9oP+T6P/LV6L/H6reMaD4FmSX1k9TH5NwT347cXy2KP+ESZOT0njxLcqSMX0p",
	"ON/PgPhkGa2rIS50Z8CSwMEGkV4/x8KxufxUpb7KNPoxiY5MonEh02ZS84CdrKTkhpBqjin8uXg1rAEt",
	"EIT8+t0XVXoWF/w/i5SmSvD7/Hf+O1gYNzMiP4fleD3+HwAA//9wiPtGLw8AAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
