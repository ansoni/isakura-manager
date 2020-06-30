// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetUIContentParams creates a new GetUIContentParams object
// no default values defined in spec.
func NewGetUIContentParams() GetUIContentParams {

	return GetUIContentParams{}
}

// GetUIContentParams contains all the bound params for the get Ui content operation
// typically these are obtained from a http.Request
//
// swagger:parameters getUiContent
type GetUIContentParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	Resource string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetUIContentParams() beforehand.
func (o *GetUIContentParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rResource, rhkResource, _ := route.Params.GetOK("resource")
	if err := o.bindResource(rResource, rhkResource, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindResource binds and validates parameter Resource from path.
func (o *GetUIContentParams) bindResource(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Resource = raw

	return nil
}
