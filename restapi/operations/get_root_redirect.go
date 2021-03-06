// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetRootRedirectHandlerFunc turns a function with the right signature into a get root redirect handler
type GetRootRedirectHandlerFunc func(GetRootRedirectParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetRootRedirectHandlerFunc) Handle(params GetRootRedirectParams) middleware.Responder {
	return fn(params)
}

// GetRootRedirectHandler interface for that can handle valid get root redirect params
type GetRootRedirectHandler interface {
	Handle(GetRootRedirectParams) middleware.Responder
}

// NewGetRootRedirect creates a new http.Handler for the get root redirect operation
func NewGetRootRedirect(ctx *middleware.Context, handler GetRootRedirectHandler) *GetRootRedirect {
	return &GetRootRedirect{Context: ctx, Handler: handler}
}

/*GetRootRedirect swagger:route GET / getRootRedirect

app redirect

*/
type GetRootRedirect struct {
	Context *middleware.Context
	Handler GetRootRedirectHandler
}

func (o *GetRootRedirect) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetRootRedirectParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
