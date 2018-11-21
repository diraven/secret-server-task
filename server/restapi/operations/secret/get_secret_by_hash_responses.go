// Code generated by go-swagger; DO NOT EDIT.

package secret

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	models "github.com/diraven/secret-server-task/server/models"
	"github.com/go-openapi/runtime"
)

// GetSecretByHashOKCode is the HTTP code returned for type GetSecretByHashOK
const GetSecretByHashOKCode int = 200

/*GetSecretByHashOK successful operation

swagger:response getSecretByHashOK
*/
type GetSecretByHashOK struct {

	/*
	  In: Body
	*/
	Payload *models.Secret `json:"body,omitempty"`
}

// NewGetSecretByHashOK creates GetSecretByHashOK with default headers values
func NewGetSecretByHashOK() *GetSecretByHashOK {

	return &GetSecretByHashOK{}
}

// WithPayload adds the payload to the get secret by hash o k response
func (o *GetSecretByHashOK) WithPayload(payload *models.Secret) *GetSecretByHashOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get secret by hash o k response
func (o *GetSecretByHashOK) SetPayload(payload *models.Secret) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSecretByHashOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetSecretByHashNotFoundCode is the HTTP code returned for type GetSecretByHashNotFound
const GetSecretByHashNotFoundCode int = 404

/*GetSecretByHashNotFound Secret not found

swagger:response getSecretByHashNotFound
*/
type GetSecretByHashNotFound struct {
}

// NewGetSecretByHashNotFound creates GetSecretByHashNotFound with default headers values
func NewGetSecretByHashNotFound() *GetSecretByHashNotFound {

	return &GetSecretByHashNotFound{}
}

// WriteResponse to the client
func (o *GetSecretByHashNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}
