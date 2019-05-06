// Code generated by go-swagger; DO NOT EDIT.

package pipeline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/metalmatze/cd/api/v1/models"
)

// PipelineOKCode is the HTTP code returned for type PipelineOK
const PipelineOKCode int = 200

/*PipelineOK OK

swagger:response pipelineOK
*/
type PipelineOK struct {

	/*
	  In: Body
	*/
	Payload *models.Pipeline `json:"body,omitempty"`
}

// NewPipelineOK creates PipelineOK with default headers values
func NewPipelineOK() *PipelineOK {

	return &PipelineOK{}
}

// WithPayload adds the payload to the pipeline o k response
func (o *PipelineOK) WithPayload(payload *models.Pipeline) *PipelineOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the pipeline o k response
func (o *PipelineOK) SetPayload(payload *models.Pipeline) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PipelineOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PipelineBadRequestCode is the HTTP code returned for type PipelineBadRequest
const PipelineBadRequestCode int = 400

/*PipelineBadRequest bad request

swagger:response pipelineBadRequest
*/
type PipelineBadRequest struct {
}

// NewPipelineBadRequest creates PipelineBadRequest with default headers values
func NewPipelineBadRequest() *PipelineBadRequest {

	return &PipelineBadRequest{}
}

// WriteResponse to the client
func (o *PipelineBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// PipelineInternalServerErrorCode is the HTTP code returned for type PipelineInternalServerError
const PipelineInternalServerErrorCode int = 500

/*PipelineInternalServerError internal server error

swagger:response pipelineInternalServerError
*/
type PipelineInternalServerError struct {
}

// NewPipelineInternalServerError creates PipelineInternalServerError with default headers values
func NewPipelineInternalServerError() *PipelineInternalServerError {

	return &PipelineInternalServerError{}
}

// WriteResponse to the client
func (o *PipelineInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}