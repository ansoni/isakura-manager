// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetChannelsGuideHandlerFunc turns a function with the right signature into a get channels guide handler
type GetChannelsGuideHandlerFunc func(GetChannelsGuideParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetChannelsGuideHandlerFunc) Handle(params GetChannelsGuideParams) middleware.Responder {
	return fn(params)
}

// GetChannelsGuideHandler interface for that can handle valid get channels guide params
type GetChannelsGuideHandler interface {
	Handle(GetChannelsGuideParams) middleware.Responder
}

// NewGetChannelsGuide creates a new http.Handler for the get channels guide operation
func NewGetChannelsGuide(ctx *middleware.Context, handler GetChannelsGuideHandler) *GetChannelsGuide {
	return &GetChannelsGuide{Context: ctx, Handler: handler}
}

/*GetChannelsGuide swagger:route GET /channels/guide getChannelsGuide

get channel guide

*/
type GetChannelsGuide struct {
	Context *middleware.Context
	Handler GetChannelsGuideHandler
}

func (o *GetChannelsGuide) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetChannelsGuideParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
