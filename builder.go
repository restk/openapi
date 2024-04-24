// Copyright 2024 Arianit Uka
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package openapi

import (
	"net/http"
	"reflect"
	"strconv"
)

// Builder provides builders for building an OpenAPI spec from code
type Builder struct {
	openAPI *OpenAPI
}

// New returns an OpenAPI builder that can be used to easily generate OpenAPI specs from code.
func New(title, version string) *Builder {
	schemaPrefix := "#/components/schemas/"

	registry := NewMapRegistry(schemaPrefix, DefaultSchemaNamer)

	o := &OpenAPI{
		OpenAPI: "3.1.0",
		Info: &Info{
			Title:   title,
			Version: version,
		},
		Components: &Components{
			Schemas:         registry,
			SecuritySchemes: make(map[string]*SecurityScheme),
		},
	}

	return &Builder{
		openAPI: o,
	}
}

// Contact adds a contact
func (b *Builder) Contact() *Builder {
	contact := &Contact{}

	b.openAPI.Info.Contact = contact

	return b
}

// ContactBuilder helps build a contact
type ContactBuilder struct {
	contact *Contact
}

// Name sets the name for the contact
func (cb *ContactBuilder) Name(name string) *ContactBuilder {
	cb.contact.Name = name

	return cb
}

// URL sets the URL for the contact
func (cb *ContactBuilder) URL(url string) *ContactBuilder {
	cb.contact.URL = url

	return cb
}

// Email sets the email for the contact
func (cb *ContactBuilder) Email(email string) *ContactBuilder {
	cb.contact.Email = email

	return cb
}

// License adds a License to the OpenAPI info object
func (b *Builder) License() *Builder {
	license := &License{}
	b.openAPI.Info.License = license

	return b
}

// LicenseBuilder helps build a licensea
type LicenseBuilder struct {
	license *License
}

// Name sets the license name
func (lb *LicenseBuilder) Name(name string) *LicenseBuilder {
	lb.license.Name = name

	return lb
}

// Identifier sets the identifier
func (lb *LicenseBuilder) Identifier(identifier string) *LicenseBuilder {
	if lb.license.URL != "" {
		panic("license url and identifier cannot both be set, set only one of them")
	}

	lb.license.Identifier = identifier

	return lb
}

// URL sets the URL
func (lb *LicenseBuilder) URL(url string) *LicenseBuilder {
	if lb.license.Identifier != "" {
		panic("license url and identifier cannot both be set, set only one of them")
	}

	lb.license.URL = url

	return lb
}

// Server adds a server
func (b *Builder) Server() *ServerBuilder {
	server := &Server{}
	b.openAPI.Servers = append(b.openAPI.Servers, server)

	return &ServerBuilder{
		server: server,
	}
}

// Security sets global security
func (b *Builder) Security(securitySchema string, params []string) *Builder {
	if b.openAPI.Security == nil {
		b.openAPI.Security = []map[string][]string{}
	}

	s := map[string][]string{}
	s[securitySchema] = params

	b.openAPI.Security = append(b.openAPI.Security, s)
	return b
}

// Description sets the description for the API
func (b *Builder) Description(description string) *Builder {
	b.openAPI.Info.Description = description

	return b
}

// TermsOfService sets url for the terms of service
func (b *Builder) TermsOfService(url string) *Builder {
	b.openAPI.Info.TermsOfService = url

	return b
}

// BasicAuth adds a BasicAuth security schema
func (b *Builder) BasicAuth() *Builder {
	b.openAPI.Components.SecuritySchemes["BasicAuth"] = &SecurityScheme{
		Type:   "http",
		Scheme: "basic",
	}

	return b
}

// BearerAuth adds a BearerAuth security schema
func (b *Builder) BearerAuth() *Builder {
	b.openAPI.Components.SecuritySchemes["BearerAuth"] = &SecurityScheme{
		Type:   "http",
		Scheme: "bearer",
	}

	return b
}

// ApiKeyAuth adds a ApiKeyAuth security schema. We expect the API key to be in a Header and you specify the header as the argument to this function
func (b *Builder) ApiKeyAuth(header string) *Builder {
	b.openAPI.Components.SecuritySchemes["ApiKeyAuth"] = &SecurityScheme{
		Type: "apiKey",
		In:   "header",
		Name: header,
	}

	return b
}

// OpenID adds a OpenID security schema. url is an OpenId Connect URL to discover OAuth2
func (b *Builder) OpenID(url string) *Builder {
	b.openAPI.Components.SecuritySchemes["ApiKeyAuth"] = &SecurityScheme{
		Type:             "openIdConnect",
		OpenIDConnectURL: url,
	}

	return b
}

func (b *Builder) OAuth2() *OAuth2Builder {
	flows := &OAuthFlows{}

	b.openAPI.Components.SecuritySchemes["OAuth2"] = &SecurityScheme{
		Type:  "oauth2",
		Flows: flows,
	}

	return &OAuth2Builder{
		flows: flows,
	}
}

// OAuth2Builder helps build OAuth2
type OAuth2Builder struct {
	flows *OAuthFlows
}

// Implicit sets the Implicit flow
func (ob *OAuth2Builder) Implicit() *OAuthFlowBuilder {
	ob.flows.Implicit = &OAuthFlow{}

	return &OAuthFlowBuilder{
		flow: ob.flows.Implicit,
	}
}

// Password sets the Password flow
func (ob *OAuth2Builder) Password() *OAuthFlowBuilder {
	ob.flows.Password = &OAuthFlow{}

	return &OAuthFlowBuilder{
		flow: ob.flows.Password,
	}
}

// ClientCredentials sets the ClientCredentials flow
func (ob *OAuth2Builder) ClientCredentials() *OAuthFlowBuilder {
	ob.flows.ClientCredentials = &OAuthFlow{}

	return &OAuthFlowBuilder{
		flow: ob.flows.Password,
	}
}

// AuthorizationCode sets the AuthorizationCode flow
func (ob *OAuth2Builder) AuthorizationCode() *OAuthFlowBuilder {
	ob.flows.AuthorizationCode = &OAuthFlow{}

	return &OAuthFlowBuilder{
		flow: ob.flows.AuthorizationCode,
	}
}

// OAuthFlowBuilder helps build a oauth flow
type OAuthFlowBuilder struct {
	flow *OAuthFlow
}

// AuthorizationURL sets the authorization url for the flow
func (b *OAuthFlowBuilder) AuthorizationURL(url string) *OAuthFlowBuilder {
	b.flow.AuthorizationURL = url

	return b
}

// TokenURL sets the token url for the flow
func (b *OAuthFlowBuilder) TokenURL(url string) *OAuthFlowBuilder {
	b.flow.TokenURL = url

	return b
}

// RefreshURL sets the refresh url for the flow
func (b *OAuthFlowBuilder) RefreshURL(url string) *OAuthFlowBuilder {
	b.flow.RefreshURL = url

	return b
}

// Scopes sets the scope for the flow
func (b *OAuthFlowBuilder) Scopes(scopes map[string]string) *OAuthFlowBuilder {
	b.flow.Scopes = scopes

	return b
}

// SecurityScheme adds a custom SecurityScheme.
func (b *Builder) SecurityScheme(name string, scheme *SecurityScheme) *Builder {
	b.openAPI.Components.SecuritySchemes[name] = scheme

	return b
}

// Register registers a new Operation.
func (b *Builder) Register(op *Operation) *OperationBuilder {
	if op.Method == "" || op.Path == "" {
		panic("method and path must be specified in operation")
	}
	op.Responses = make(map[string]*Response)

	b.openAPI.AddOperation(op)

	return &OperationBuilder{
		op:      op,
		openAPI: b.openAPI,
	}
}

// FindOperationIdByTag finds the first operation with the tag and returns its id. If nothing is found, this returns an empty string.
/*
func (b *Builder) FindOperationIdByTag(tag string) string {
	for _, path := range b.openAPI.Paths {
		if path.Get != nil && path.Get.Tags != nil {
			if slices.Contains(path.Get.Tags, tag) {
				return path.Get.OperationID
			}
		}
		if path.Put != nil && path.Put.Tags != nil {
			if slices.Contains(path.Put.Tags, tag) {
				return path.Put.OperationID
			}
		}
		if path.Post != nil && path.Post.Tags != nil {
			if slices.Contains(path.Post.Tags, tag) {
				return path.Post.OperationID
			}
		}
		if path.Delete != nil && path.Delete.Tags != nil {
			if slices.Contains(path.Delete.Tags, tag) {
				return path.Delete.OperationID
			}
		}
		if path.Options != nil && path.Options.Tags != nil {
			if slices.Contains(path.Options.Tags, tag) {
				return path.Options.OperationID
			}
		}
		if path.Head != nil && path.Head.Tags != nil {
			if slices.Contains(path.Head.Tags, tag) {
				return path.Head.OperationID
			}
		}
		if path.Patch != nil && path.Patch.Tags != nil {
			if slices.Contains(path.Patch.Tags, tag) {
				return path.Patch.OperationID
			}
		}
		if path.Trace != nil && path.Trace.Tags != nil {
			if slices.Contains(path.Trace.Tags, tag) {
				return path.Trace.OperationID
			}
		}
	}

	return ""
}
*/

// OperationBuilder assists in building an openapi.Operation
type OperationBuilder struct {
	op      *Operation
	openAPI *OpenAPI
}

// Tag adds a tag
func (ob *OperationBuilder) Tag(tag string) *OperationBuilder {
	ob.op.Tags = append(ob.op.Tags, tag)

	return ob
}

// Summary sets the summary.
func (ob *OperationBuilder) Summary(summary string) *OperationBuilder {
	ob.op.Summary = summary

	return ob
}

// Description sets the description
func (ob *OperationBuilder) Description(description string) *OperationBuilder {
	ob.op.Description = description

	return ob
}

// Security sets the security for this operation
func (ob *OperationBuilder) Security(securitySchema string, params []string) *OperationBuilder {
	if ob.op.Security == nil {
		ob.op.Security = []map[string][]string{}
	}

	s := map[string][]string{}
	s[securitySchema] = params

	ob.op.Security = append(ob.op.Security, s)

	return ob
}

// Server adds a server for this operation.
func (ob *OperationBuilder) Server() *ServerBuilder {
	server := &Server{}

	ob.op.Servers = append(ob.op.Servers, server)

	return &ServerBuilder{
		server: server,
	}
}

// Response adds a response with a status code.
func (ob *OperationBuilder) Response(status int) *ResponseBuilder {
	statusStr := strconv.Itoa(status)

	if ob.op.Responses[statusStr] == nil {
		ob.op.Responses[statusStr] = &Response{}
	}

	return &ResponseBuilder{
		openAPI:            ob.openAPI,
		response:           ob.op.Responses[statusStr],
		defaultContentType: "application/json",
		nextContentType:    "",
	}
}

type ResponseBuilder struct {
	openAPI  *OpenAPI
	response *Response

	defaultContentType string
	nextContentType    string
}

// DefaultContentType sets the default content type of the Response. By default, the content type is application/json
func (rb *ResponseBuilder) DefaultContentType(contentType string) *ResponseBuilder {
	rb.defaultContentType = contentType

	return rb
}

// ContentType applies this content type to the next Body() call. This only applies to the next Body() call and any subsequent calls to Body() will default to
// DefaultContentType(). If you want to change the default content type for all Body() calls, call DefaultContentType() then call Body() without using ContentType()
func (rb *ResponseBuilder) ContentType(contentType string) *ResponseBuilder {
	rb.nextContentType = contentType

	return rb
}

// Body adds a body. f is the type that is used for the body's schema. f can be a struct, slice, map, or a basic type. For basic types, you can use our
// helper methods such as openapi.IntType, openapi.StringType, openapi.UintType, etc. (see types.go for all basic types.)
func (rb *ResponseBuilder) Body(f any) *MediaTypeBuilder {
	responseType := reflect.TypeOf(f)

	registry := rb.openAPI.Components.Schemas
	schema := registry.Schema(responseType, true, "")

	var contentType string
	var resetNextContentType bool
	if rb.nextContentType != "" {
		contentType = rb.nextContentType
		resetNextContentType = true
	} else {
		contentType = rb.defaultContentType
	}

	if rb.response.Content == nil {
		rb.response.Content = map[string]*MediaType{}
	}
	if rb.response.Content[contentType] == nil {
		rb.response.Content[contentType] = &MediaType{}
	}
	if rb.response.Content[contentType] != nil && rb.response.Content[contentType].Schema == nil {
		rb.response.Content[contentType].Schema = schema
	}

	if resetNextContentType {
		rb.nextContentType = ""
	}

	return &MediaTypeBuilder{
		openAPI:   rb.openAPI,
		mediaType: rb.response.Content[contentType],
	}
}

type MediaTypeBuilder struct {
	openAPI   *OpenAPI
	mediaType *MediaType
}

// Schema overrides the schema with the type f
func (mtb *MediaTypeBuilder) Schema(f any) {
	schemaType := reflect.TypeOf(f)

	registry := mtb.openAPI.Components.Schemas
	schema := registry.Schema(schemaType, true, "")

	mtb.mediaType.Schema = schema
}

// Example sets the example for this media type
func (mtb *MediaTypeBuilder) Example(example string) *MediaTypeBuilder {
	mtb.mediaType.Example = example

	return mtb
}

// AddExample adds an example with a name.
func (mtb *MediaTypeBuilder) AddExample(name string) *ExampleBuilder {
	example := &Example{}
	if mtb.mediaType.Examples == nil {
		mtb.mediaType.Examples = map[string]*Example{}
	}

	mtb.mediaType.Examples[name] = example

	return &ExampleBuilder{
		example: example,
	}
}

// Header adds a Header.
func (rb *ResponseBuilder) Header(name string, f any) *ParamBuilder {
	responseType := reflect.TypeOf(f)

	registry := rb.openAPI.Components.Schemas
	schema := registry.Schema(responseType, true, "")

	param := &Param{
		Schema: schema,
	}

	rb.response.Headers[name] = param
	return &ParamBuilder{
		param: param,
	}
}

// Link adds a link to the response.
func (rb *ResponseBuilder) Link(name string) *LinkBuilder {
	link := &Link{}
	rb.response.Links[name] = link

	return &LinkBuilder{
		link: link,
	}
}

// LinkBuilder assists with making links
type LinkBuilder struct {
	link *Link
}

// Set the OperationID of this link.
func (lb *LinkBuilder) OperationID(operationID string) *LinkBuilder {
	if lb.link.OperationRef != "" {
		panic("operationID and OperationRef both cannot be set on a link, pick one.")
	}
	lb.link.OperationID = operationID

	return lb
}

// OperationRef sets the OperationRef for the link.
func (lb *LinkBuilder) OperationRef(operationRef string) *LinkBuilder {
	if lb.link.OperationID != "" {
		panic("operationID and OperationRef both cannot be set on a link, pick one.")
	}

	lb.link.OperationRef = operationRef

	return lb
}

// AddParam adds a link param, expression is a OpenAPI runtime expressions (see https://swagger.io/docs/specification/links/)
func (lb *LinkBuilder) AddParam(name string, expression any) *LinkBuilder {
	lb.link.Parameters[name] = expression

	return lb
}

// RequestBody sets the link request body, expression is an OpenAPI run time expression (see https://swagger.io/docs/specification/links/)
func (lb *LinkBuilder) RequestBody(expression any) *LinkBuilder {
	lb.link.RequestBody = expression

	return lb
}

// Description sets the description of the link
func (lb *LinkBuilder) Description(description string) *LinkBuilder {
	lb.link.Description = description

	return lb
}

// Server adds a server for the link
func (lb *LinkBuilder) Server() *ServerBuilder {
	server := &Server{}

	lb.link.Server = server
	return &ServerBuilder{
		server: server,
	}
}

// ServerBuilder helps with building a Server
type ServerBuilder struct {
	server *Server
}

// URL sets the server URL
func (sb *ServerBuilder) URL(url string) *ServerBuilder {
	sb.server.URL = url

	return sb
}

// Description sets the server description
func (sb *ServerBuilder) Description(description string) *ServerBuilder {
	sb.server.Description = description

	return sb
}

// AddVariable adds a variable to the server.
func (sb *ServerBuilder) AddVariable(name string) *ServerVariableBuilder {
	variable := &ServerVariable{}

	sb.server.Variables[name] = variable

	return &ServerVariableBuilder{
		variable: variable,
	}
}

// ServerVariableBuilder helps build a server variable
type ServerVariableBuilder struct {
	variable *ServerVariable
}

// Enum sets the enum for the variable
func (svb *ServerVariableBuilder) Enum(enum []string) *ServerVariableBuilder {
	svb.variable.Enum = enum

	return svb
}

// Default sets the default value for the variable
func (svb *ServerVariableBuilder) Default(value string) *ServerVariableBuilder {
	svb.variable.Default = value

	return svb
}

// Description sets the description for the variable
func (svb *ServerVariableBuilder) Description(description string) *ServerVariableBuilder {
	svb.variable.Description = description

	return svb
}

// Callback adds a callback.
func (ob *OperationBuilder) Callback(event string, op *Operation) *OperationBuilder {
	if op.Method == "" || op.Path == "" || event == "" {
		panic("event and op.method and op.path must be specified")
	}

	if ob.op.Callbacks == nil {
		ob.op.Callbacks = map[string]map[string]*PathItem{}
	}
	if ob.op.Callbacks[event] == nil {
		ob.op.Callbacks[event] = map[string]*PathItem{}
	}

	item := ob.op.Callbacks[event][op.Path]
	if item == nil {
		item = &PathItem{}
		ob.op.Callbacks[event][op.Path] = item
	}

	op.Responses = make(map[string]*Response)

	switch op.Method {
	case http.MethodGet:
		item.Get = op
	case http.MethodPost:
		item.Post = op
	case http.MethodPut:
		item.Put = op
	case http.MethodPatch:
		item.Patch = op
	case http.MethodDelete:
		item.Delete = op
	case http.MethodHead:
		item.Head = op
	case http.MethodOptions:
		item.Options = op
	case http.MethodTrace:
		item.Trace = op
	default:
		panic("unknown method " + op.Method)
	}

	return &OperationBuilder{
		op:      op,
		openAPI: ob.openAPI,
	}
}

// Request returns a RequestBuilder which helps build a request
func (ob *OperationBuilder) Request() *RequestBuilder {
	return &RequestBuilder{
		op:                 ob.op,
		openAPI:            ob.openAPI,
		defaultContentType: "application/json",
		nextContentType:    "",
	}
}

// RequestBuilder helps build a Request
type RequestBuilder struct {
	openAPI *OpenAPI
	op      *Operation

	defaultContentType string
	nextContentType    string
}

// DefaultContentType sets the content type for all Body() calls
func (rb *RequestBuilder) DefaultContentType(contentType string) *RequestBuilder {
	rb.defaultContentType = contentType

	return rb
}

// ContentType sets the content type of the Request for the next Body() call
func (rb *RequestBuilder) ContentType(contentType string) *RequestBuilder {
	rb.nextContentType = contentType

	return rb
}

// Body sets the RequestBody
func (rb *RequestBuilder) Body(f any) *RequestBodyBuilder {
	responseType := reflect.TypeOf(f)

	registry := rb.openAPI.Components.Schemas
	ref := registry.Schema(responseType, true, "")

	var contentType string
	if rb.nextContentType != "" {
		contentType = rb.nextContentType
	} else {
		contentType = rb.defaultContentType
	}

	mediaType := &MediaType{
		Schema: ref,
	}

	if rb.op.RequestBody == nil {
		rb.op.RequestBody = &RequestBody{
			Required: true,
			Content: map[string]*MediaType{
				contentType: mediaType,
			},
		}
	} else {
		rb.op.RequestBody.Content[contentType] = mediaType
	}

	if rb.nextContentType != "" {
		rb.nextContentType = ""
	}

	return &RequestBodyBuilder{
		mediaTypeBuilder: &MediaTypeBuilder{
			openAPI:   rb.openAPI,
			mediaType: mediaType,
		},

		requestBody: rb.op.RequestBody,
	}
}

type RequestBodyBuilder struct {
	mediaTypeBuilder *MediaTypeBuilder
	requestBody      *RequestBody
}

// Description sets the Description for the RequestBody
func (rbb *RequestBodyBuilder) Description(description string) *RequestBodyBuilder {
	rbb.requestBody.Description = description

	return rbb
}

// Required makes the Body required for the Request. By default when you call Body() it is already set to true
func (rbb *RequestBodyBuilder) Required(required bool) *RequestBodyBuilder {
	rbb.requestBody.Required = required

	return rbb
}

// Example sets the example for the body
func (rbb *RequestBodyBuilder) Example(example string) *RequestBodyBuilder {
	rbb.mediaTypeBuilder.Example(example)

	return rbb
}

// AddExample adds an example for the body
func (rbb *RequestBodyBuilder) AddExample(example string) *ExampleBuilder {
	return rbb.mediaTypeBuilder.AddExample(example)
}

// QueryParam adds a query param
func (rb *RequestBuilder) QueryParam(name string, f any) *ParamBuilder {
	return rb.Param("query", name, f)
}

// PathParam adds a path param
func (rb *RequestBuilder) PathParam(name string, f any) *ParamBuilder {
	return rb.Param("path", name, f)
}

// CookieParam adds a cookie param
func (rb *RequestBuilder) CookieParam(name string, f any) *ParamBuilder {
	return rb.Param("cookie", name, f)
}

// Param adds a param. in can be "path", "query", or "cookie". See helper functions QueryParam(), PathParam() and CookieParam() for a shorthand version
func (rb *RequestBuilder) Param(in string, name string, f any) *ParamBuilder {
	paramType := reflect.TypeOf(f)

	registry := rb.openAPI.Components.Schemas
	schema := registry.Schema(paramType, true, "")

	param := &Param{
		Name:   name,
		In:     in,
		Schema: schema,
	}

	rb.op.Parameters = append(rb.op.Parameters, param)

	return &ParamBuilder{
		param: param,
	}
}

// ParamBuilder helps with building an openapi.Param
type ParamBuilder struct {
	param *Param
}

// In can be "query", "path" or "cookie". This should already be set when creating the param and is only here for a manual override.
func (pb *ParamBuilder) In(in string) *ParamBuilder {
	pb.param.In = in

	return pb
}

// Description sets the description of the param
func (pb *ParamBuilder) Description(description string) *ParamBuilder {
	pb.param.Description = description

	return pb
}

// Required sets the param as required
func (pb *ParamBuilder) Required(required bool) *ParamBuilder {
	pb.param.Required = required
	return pb
}

// Example sets the example for the param. If you want to add multiple examples, call AddExample() instead.
func (pb *ParamBuilder) Example(example string) *ParamBuilder {
	pb.param.Example = example

	return pb
}

// AddExample adds an example to a list of examples.
func (pb *ParamBuilder) AddExample(name string) *ExampleBuilder {
	example := &Example{}

	if pb.param.Examples == nil {
		pb.param.Examples = map[string]*Example{}
	}
	pb.param.Examples[name] = example

	return &ExampleBuilder{
		example: example,
	}
}

// ExampleBuilder helps build examples.
type ExampleBuilder struct {
	example *Example
}

// Ref sets the ref of the example. This references another example. $ref: '#/components/examples/objectExample'
func (eb *ExampleBuilder) Ref(ref string) {
	eb.example.Ref = ref
}

// Value sets the value for the example.
func (eb *ExampleBuilder) Value(value any) *ExampleBuilder {
	if eb.example.ExternalValue != "" {
		panic("Value and ExternalValue cannot both be set on example, pick one")
	}
	eb.example.Value = value

	return eb
}

// ExternalValue sets the ExternalValue
func (eb *ExampleBuilder) ExternalValue(uri string) *ExampleBuilder {
	if eb.example.Value != nil {
		panic("Value and ExternalValue cannot both be set on example, pick one")
	}
	eb.example.ExternalValue = uri

	return eb
}

// Summary sets a summary for an example. This is required.
func (eb *ExampleBuilder) Summary(summary string) *ExampleBuilder {
	eb.example.Summary = summary

	return eb
}

// Description sets the description for the example.
func (eb *ExampleBuilder) Description(description string) *ExampleBuilder {
	eb.example.Description = description

	return eb
}

// Style sets the style of the param. See https://swagger.io/docs/specification/serialization/
func (pb *ParamBuilder) Style(style string) *ParamBuilder {
	pb.param.Style = style

	return pb
}

// Explode explodes
func (pb *ParamBuilder) Explode(explode bool) *ParamBuilder {
	pb.param.Explode = &explode

	return pb
}

// Registry returns the Registry.
func (b *Builder) Registry() Registry {
	return b.openAPI.Components.Schemas
}

// OpenAPI returns the OpenAPI struct.
func (b *Builder) OpenAPI() *OpenAPI {
	return b.openAPI
}
