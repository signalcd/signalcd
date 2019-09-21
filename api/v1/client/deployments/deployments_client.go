// Code generated by go-swagger; DO NOT EDIT.

package deployments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new deployments API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for deployments API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
CurrentDeployment returns the currently active deployment
*/
func (a *Client) CurrentDeployment(params *CurrentDeploymentParams, authInfo runtime.ClientAuthInfoWriter) (*CurrentDeploymentOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCurrentDeploymentParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "currentDeployment",
		Method:             "GET",
		PathPattern:        "/deployments/current",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &CurrentDeploymentReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*CurrentDeploymentOK), nil

}

/*
Deployments returns the history of deployments
*/
func (a *Client) Deployments(params *DeploymentsParams, authInfo runtime.ClientAuthInfoWriter) (*DeploymentsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeploymentsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "deployments",
		Method:             "GET",
		PathPattern:        "/deployments",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &DeploymentsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DeploymentsOK), nil

}

/*
SetCurrentDeployment schedules a new deployment
*/
func (a *Client) SetCurrentDeployment(params *SetCurrentDeploymentParams, authInfo runtime.ClientAuthInfoWriter) (*SetCurrentDeploymentOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetCurrentDeploymentParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "setCurrentDeployment",
		Method:             "POST",
		PathPattern:        "/deployments/current",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &SetCurrentDeploymentReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*SetCurrentDeploymentOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
