// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteContentHandlerFunc turns a function with the right signature into a delete content handler
type DeleteContentHandlerFunc func(DeleteContentParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteContentHandlerFunc) Handle(params DeleteContentParams) middleware.Responder {
	return fn(params)
}

// DeleteContentHandler interface for that can handle valid delete content params
type DeleteContentHandler interface {
	Handle(DeleteContentParams) middleware.Responder
}

// NewDeleteContent creates a new http.Handler for the delete content operation
func NewDeleteContent(ctx *middleware.Context, handler DeleteContentHandler) *DeleteContent {
	return &DeleteContent{Context: ctx, Handler: handler}
}

/*DeleteContent swagger:route DELETE /content/{content_name} deleteContent

delete a piece of content from library

*/
type DeleteContent struct {
	Context *middleware.Context
	Handler DeleteContentHandler
}

func (o *DeleteContent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteContentParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
