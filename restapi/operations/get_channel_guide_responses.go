// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/ansoni/isakura-manager/models"
)

// GetChannelGuideOKCode is the HTTP code returned for type GetChannelGuideOK
const GetChannelGuideOKCode int = 200

/*GetChannelGuideOK woot!

swagger:response getChannelGuideOK
*/
type GetChannelGuideOK struct {

	/*
	  In: Body
	*/
	Payload *models.ChannelGuide `json:"body,omitempty"`
}

// NewGetChannelGuideOK creates GetChannelGuideOK with default headers values
func NewGetChannelGuideOK() *GetChannelGuideOK {

	return &GetChannelGuideOK{}
}

// WithPayload adds the payload to the get channel guide o k response
func (o *GetChannelGuideOK) WithPayload(payload *models.ChannelGuide) *GetChannelGuideOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get channel guide o k response
func (o *GetChannelGuideOK) SetPayload(payload *models.ChannelGuide) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetChannelGuideOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetChannelGuideNotFoundCode is the HTTP code returned for type GetChannelGuideNotFound
const GetChannelGuideNotFoundCode int = 404

/*GetChannelGuideNotFound Not Found

swagger:response getChannelGuideNotFound
*/
type GetChannelGuideNotFound struct {
}

// NewGetChannelGuideNotFound creates GetChannelGuideNotFound with default headers values
func NewGetChannelGuideNotFound() *GetChannelGuideNotFound {

	return &GetChannelGuideNotFound{}
}

// WriteResponse to the client
func (o *GetChannelGuideNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}
