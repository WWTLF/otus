// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
	Core "user_list/CORE"
	"user_list/HandlersImpl/UserHandlers"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	log "github.com/sirupsen/logrus"

	"user_list/models"
	"user_list/restapi/operations"
	"user_list/restapi/operations/healthcheck"
	"user_list/restapi/operations/instruments"
	"user_list/restapi/operations/user"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:generate swagger generate server --target ../../otus1 --name UserList --spec ../otus55-users-1.0.0-resolved.yaml --principal interface{}
var (
	metrics_rq_latancy = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "app_request_latency_seconds",
			Help: "Application Request Latency",
		},
		[]string{"method", "endpoint"},
	)

	metrics_rq_counter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "app_request_count",
		Help: "Application Request Count",
	}, []string{"method", "endpoint", "http_status"})
)

type CustomResponder func(http.ResponseWriter, runtime.Producer)

func (c CustomResponder) WriteResponse(w http.ResponseWriter, p runtime.Producer) {
	c(w, p)
}

func NewCustomResponder(r *http.Request, h http.Handler) middleware.Responder {
	return CustomResponder(func(w http.ResponseWriter, _ runtime.Producer) {
		h.ServeHTTP(w, r)
	})
}

func configureFlags(api *operations.UserListAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
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

	api.InstrumentsGetMetricsHandler = instruments.GetMetricsHandlerFunc(func(p instruments.GetMetricsParams) middleware.Responder {
		// some logic here
		return NewCustomResponder(p.HTTPRequest, promhttp.Handler())
	})

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

// LoggingResponseWriter is a wrapper around an http.ResponseWriter which captures the
// status code written to the response, so that it can be logged.
type LoggingResponseWriter struct {
	wrapped    http.ResponseWriter
	StatusCode int
	// Response content could also be captured here, but I was only interested in logging the response status code
}

func NewLoggingResponseWriter(wrapped http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{wrapped: wrapped}
}

func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.wrapped.Header()
}

func (lrw *LoggingResponseWriter) Write(content []byte) (int, error) {
	return lrw.wrapped.Write(content)
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.StatusCode = statusCode
	lrw.wrapped.WriteHeader(statusCode)
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		endpoint := r.URL.Path
		if strings.Contains(endpoint, "/user/") {
			endpoint = "/user"
		}
		t0 := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
		w2 := NewLoggingResponseWriter(w)
		handler.ServeHTTP(w2, r)
		total := time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)) - t0
		metrics_rq_counter.With(prometheus.Labels{"method": r.Method, "endpoint": endpoint, "http_status": fmt.Sprintf("%d", w2.StatusCode)}).Inc()
		metrics_rq_latancy.With(prometheus.Labels{"method": r.Method, "endpoint": endpoint}).Observe(float64(total / 1000))
	})
}
