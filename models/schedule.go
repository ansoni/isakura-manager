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

// Schedule schedule
// swagger:model Schedule
type Schedule struct {

	// append date
	AppendDate bool `json:"appendDate,omitempty"`

	// append time
	AppendTime bool `json:"appendTime,omitempty"`

	// filter
	Filter string `json:"filter,omitempty"`

	// folder
	Folder string `json:"folder,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// prepend date
	PrependDate bool `json:"prependDate,omitempty"`

	// prepend time
	PrependTime bool `json:"prependTime,omitempty"`

	// searches
	Searches []*ScheduleSearchesItems0 `json:"searches"`

	// watch
	Watch bool `json:"watch,omitempty"`
}

// Validate validates this schedule
func (m *Schedule) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSearches(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Schedule) validateSearches(formats strfmt.Registry) error {

	if swag.IsZero(m.Searches) { // not required
		return nil
	}

	for i := 0; i < len(m.Searches); i++ {
		if swag.IsZero(m.Searches[i]) { // not required
			continue
		}

		if m.Searches[i] != nil {
			if err := m.Searches[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("searches" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Schedule) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Schedule) UnmarshalBinary(b []byte) error {
	var res Schedule
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ScheduleSearchesItems0 schedule searches items0
// swagger:model ScheduleSearchesItems0
type ScheduleSearchesItems0 struct {

	// replace
	Replace string `json:"replace,omitempty"`

	// search
	Search string `json:"search,omitempty"`
}

// Validate validates this schedule searches items0
func (m *ScheduleSearchesItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ScheduleSearchesItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ScheduleSearchesItems0) UnmarshalBinary(b []byte) error {
	var res ScheduleSearchesItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}