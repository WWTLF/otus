// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	Core "user_list/CORE"
	"user_list/HandlersImpl/UserHandlers"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	log "github.com/sirupsen/logrus"

	"user_list/models"
	"user_list/restapi/operations"
	"user_list/restapi/operations/healthcheck"
	"user_list/restapi/operations/user"
)

//go:generate swagger generate server --target ../../otus1 --name UserList --spec ../otus55-users-1.0.0-resolved.yaml --principal interface{}

func configureFlags(api *operations.UserListAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.UserListAPI) http.Handler {
	Core.GetInstance().DBInit()

	api.ServeError = errors.ServeError
	api.Logger = log.Debugf
	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	api.UseRedoc()
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	api.HealthcheckHealthCheckHandler = healthcheck.HealthCheckHandlerFunc(func(params healthcheck.HealthCheckParams) middleware.Responder {
		status := "OK"
		return healthcheck.NewHealthCheckOK().WithPayload(&models.HealthCheckStatus{Status: &status})
	})

	api.HealthcheckReadinessHealthCheckHandler = healthcheck.ReadinessHealthCheckHandlerFunc(func(params healthcheck.ReadinessHealthCheckParams) middleware.Responder {
		if Core.GetInstance().DB != nil {
			status := "OK"
			return healthcheck.NewReadinessHealthCheckOK().WithPayload(&models.HealthCheckStatus{Status: &status})
		} else {
			var code int32 = 500
			message := "DB connection is not ready"
			return healthcheck.NewReadinessHealthCheckDefault(500).WithPayload(&models.Error{Code: &code, Message: &message})
		}
	})

	api.UserCreateUserHandler = user.CreateUserHandlerFunc(UserHandlers.CreateUser)
	api.UserDeleteUserHandler = user.DeleteUserHandlerFunc(UserHandlers.DeleteUser)
	api.UserFindUserByIDHandler = user.FindUserByIDHandlerFunc(UserHandlers.FindUserById)
	api.UserUpdateUserHandler = user.UpdateUserHandlerFunc(UserHandlers.UpdateUser)

	api.PreServerShutdown = func() {
		log.Info("Shutting down...")
		Core.GetInstance().DBClose()
	}

	api.ServerShutdown = func() {
		log.Info("Shutting down...Done")
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
