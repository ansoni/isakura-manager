// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ChannelGuide channel guide
// swagger:model ChannelGuide
type ChannelGuide struct {

	// channel name
	ChannelName string `json:"channelName,omitempty"`

	// guide
	Guide []*Guide `json:"guide"`

	// id
	ID int64 `json:"id,omitempty"`
}

// Validate validates this channel guide
func (m *ChannelGuide) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateGuide(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ChannelGuide) validateGuide(formats strfmt.Registry) error {

	if swag.IsZero(m.Guide) { // not required
		return nil
	}

	for i := 0; i < len(m.Guide); i++ {
		if swag.IsZero(m.Guide[i]) { // not required
			continue
		}

		if m.Guide[i] != nil {
			if err := m.Guide[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("guide" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ChannelGuide) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ChannelGuide) UnmarshalBinary(b []byte) error {
	var res ChannelGuide
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
