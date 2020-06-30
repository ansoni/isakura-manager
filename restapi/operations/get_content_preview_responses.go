// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetContentPreviewOKCode is the HTTP code returned for type GetContentPreviewOK
const GetContentPreviewOKCode int = 200

/*GetContentPreviewOK woot!

swagger:response getContentPreviewOK
*/
type GetContentPreviewOK struct {
}

// NewGetContentPreviewOK creates GetContentPreviewOK with default headers values
func NewGetContentPreviewOK() *GetContentPreviewOK {

	return &GetContentPreviewOK{}
}

// WriteResponse to the client
func (o *GetContentPreviewOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// GetContentPreviewNotFoundCode is the HTTP code returned for type GetContentPreviewNotFound
const GetContentPreviewNotFoundCode int = 404

/*GetContentPreviewNotFound Not Found

swagger:response getContentPreviewNotFound
*/
type GetContentPreviewNotFound struct {
}

// NewGetContentPreviewNotFound creates GetContentPreviewNotFound with default headers values
func NewGetContentPreviewNotFound() *GetContentPreviewNotFound {

	return &GetContentPreviewNotFound{}
}

// WriteResponse to the client
func (o *GetContentPreviewNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}