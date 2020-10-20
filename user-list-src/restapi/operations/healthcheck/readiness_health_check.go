// Code generated by go-swagger; DO NOT EDIT.

package healthcheck

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ReadinessHealthCheckHandlerFunc turns a function with the right signature into a readiness health check handler
type ReadinessHealthCheckHandlerFunc func(ReadinessHealthCheckParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ReadinessHealthCheckHandlerFunc) Handle(params ReadinessHealthCheckParams) middleware.Responder {
	return fn(params)
}

// ReadinessHealthCheckHandler interface for that can handle valid readiness health check params
type ReadinessHealthCheckHandler interface {
	Handle(ReadinessHealthCheckParams) middleware.Responder
}

// NewReadinessHealthCheck creates a new http.Handler for the readiness health check operation
func NewReadinessHealthCheck(ctx *middleware.Context, handler ReadinessHealthCheckHandler) *ReadinessHealthCheck {
	return &ReadinessHealthCheck{Context: ctx, Handler: handler}
}

/*ReadinessHealthCheck swagger:route GET /health/readiness healthcheck readinessHealthCheck

health check

Active health check readiness status

*/
type ReadinessHealthCheck struct {
	Context *middleware.Context
	Handler ReadinessHealthCheckHandler
}

func (o *ReadinessHealthCheck) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewReadinessHealthCheckParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}