// Code generated by go-swagger; DO NOT EDIT.

package deployments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeploymentsHandlerFunc turns a function with the right signature into a deployments handler
type DeploymentsHandlerFunc func(DeploymentsParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeploymentsHandlerFunc) Handle(params DeploymentsParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeploymentsHandler interface for that can handle valid deployments params
type DeploymentsHandler interface {
	Handle(DeploymentsParams, interface{}) middleware.Responder
}

// NewDeployments creates a new http.Handler for the deployments operation
func NewDeployments(ctx *middleware.Context, handler DeploymentsHandler) *Deployments {
	return &Deployments{Context: ctx, Handler: handler}
}

/*Deployments swagger:route GET /deployments deployments deployments

Returns the history of deployments

*/
type Deployments struct {
	Context *middleware.Context
	Handler DeploymentsHandler
}

func (o *Deployments) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeploymentsParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
