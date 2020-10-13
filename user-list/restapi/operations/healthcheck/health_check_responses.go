// Code generated by go-swagger; DO NOT EDIT.

package healthcheck

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"user_list/models"
)

// HealthCheckOKCode is the HTTP code returned for type HealthCheckOK
const HealthCheckOKCode int = 200

/*HealthCheckOK user response

swagger:response healthCheckOK
*/
type HealthCheckOK struct {

	/*
	  In: Body
	*/
	Payload *models.HealthCheckStatus `json:"body,omitempty"`
}

// NewHealthCheckOK creates HealthCheckOK with default headers values
func NewHealthCheckOK() *HealthCheckOK {

	return &HealthCheckOK{}
}

// WithPayload adds the payload to the health check o k response
func (o *HealthCheckOK) WithPayload(payload *models.HealthCheckStatus) *HealthCheckOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the health check o k response
func (o *HealthCheckOK) SetPayload(payload *models.HealthCheckStatus) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *HealthCheckOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*HealthCheckDefault unexpected error

swagger:response healthCheckDefault
*/
type HealthCheckDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewHealthCheckDefault creates HealthCheckDefault with default headers values
func NewHealthCheckDefault(code int) *HealthCheckDefault {
	if code <= 0 {
		code = 500
	}

	return &HealthCheckDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the health check default response
func (o *HealthCheckDefault) WithStatusCode(code int) *HealthCheckDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the health check default response
func (o *HealthCheckDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the health check default response
func (o *HealthCheckDefault) WithPayload(payload *models.Error) *HealthCheckDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the health check default response
func (o *HealthCheckDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *HealthCheckDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
