// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetContentHandlerFunc turns a function with the right signature into a get content handler
type GetContentHandlerFunc func(GetContentParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetContentHandlerFunc) Handle(params GetContentParams) middleware.Responder {
	return fn(params)
}

// GetContentHandler interface for that can handle valid get content params
type GetContentHandler interface {
	Handle(GetContentParams) middleware.Responder
}

// NewGetContent creates a new http.Handler for the get content operation
func NewGetContent(ctx *middleware.Context, handler GetContentHandler) *GetContent {
	return &GetContent{Context: ctx, Handler: handler}
}

/*GetContent swagger:route GET /content getContent

get existing content we have downloaded

*/
type GetContent struct {
	Context *middleware.Context
	Handler GetContentHandler
}

func (o *GetContent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetContentParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
