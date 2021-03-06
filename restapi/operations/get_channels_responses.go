// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/ansoni/isakura-manager/models"
)

// GetChannelsOKCode is the HTTP code returned for type GetChannelsOK
const GetChannelsOKCode int = 200

/*GetChannelsOK woot!

swagger:response getChannelsOK
*/
type GetChannelsOK struct {

	/*
	  In: Body
	*/
	Payload []*models.Channel `json:"body,omitempty"`
}

// NewGetChannelsOK creates GetChannelsOK with default headers values
func NewGetChannelsOK() *GetChannelsOK {

	return &GetChannelsOK{}
}

// WithPayload adds the payload to the get channels o k response
func (o *GetChannelsOK) WithPayload(payload []*models.Channel) *GetChannelsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get channels o k response
func (o *GetChannelsOK) SetPayload(payload []*models.Channel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetChannelsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.Channel, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
