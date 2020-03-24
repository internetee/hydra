// Code generated by go-swagger; DO NOT EDIT.

package public

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/ory/hydra/internal/httpclient/models"
)

// WellKnownReader is a Reader for the WellKnown structure.
type WellKnownReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WellKnownReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWellKnownOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewWellKnownInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewWellKnownOK creates a WellKnownOK with default headers values
func NewWellKnownOK() *WellKnownOK {
	return &WellKnownOK{}
}

/*WellKnownOK handles this case with default header values.

JSONWebKeySet
*/
type WellKnownOK struct {
	Payload *models.JSONWebKeySet
}

func (o *WellKnownOK) Error() string {
	return fmt.Sprintf("[GET /.well-known/jwks.json][%d] wellKnownOK  %+v", 200, o.Payload)
}

func (o *WellKnownOK) GetPayload() *models.JSONWebKeySet {
	return o.Payload
}

func (o *WellKnownOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.JSONWebKeySet)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWellKnownInternalServerError creates a WellKnownInternalServerError with default headers values
func NewWellKnownInternalServerError() *WellKnownInternalServerError {
	return &WellKnownInternalServerError{}
}

/*WellKnownInternalServerError handles this case with default header values.

genericError
*/
type WellKnownInternalServerError struct {
	Payload *models.GenericError
}

func (o *WellKnownInternalServerError) Error() string {
	return fmt.Sprintf("[GET /.well-known/jwks.json][%d] wellKnownInternalServerError  %+v", 500, o.Payload)
}

func (o *WellKnownInternalServerError) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *WellKnownInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
