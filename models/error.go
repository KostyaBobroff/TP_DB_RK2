package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Error error
// swagger:model Error
type Error struct {

	// Текстовое описание ошибки.
	// В процессе проверки API никаких проверок на содерижимое данного описание не делается.
	//
	// Read Only: true
	Message string `json:"message,omitempty"`
}

// Validate validates this error
func (m *Error) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalJSON interface implementation
func (m *Error) MarshalJSON() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalJSONJSON interface implementation
func (m *Error) UnmarshalJSONJSON(b []byte) error {
	var res Error
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
