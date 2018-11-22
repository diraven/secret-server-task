// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/diraven/secret-server-task/server/models"
	"github.com/diraven/secret-server-task/server/storage"
	"github.com/diraven/secret-server-task/server/storage/json_memory"
	"log"
	"net/http"

	"github.com/diraven/secret-server-task/server/restapi/operations"
	"github.com/diraven/secret-server-task/server/restapi/operations/secret"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

//go:generate swagger generate server --target .. --name Secret --spec ../swagger.yaml

// Prepare our secrets storage.
var secretStorage storage.Secret

// Initialize secrets storage.
func init() {
	var err error
	if secretStorage, err = json_memory.NewJSONMemory("./secrets.json"); err != nil {
		log.Fatal(err)
	}
}

//noinspection GoUnusedParameter
func configureFlags(api *operations.SecretAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SecretAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UrlformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	//api.XMLProducer = runtime.XMLProducer()

	// Adding secrets.
	api.SecretAddSecretHandler = secret.AddSecretHandlerFunc(func(params secret.AddSecretParams) middleware.Responder {
		// params verification and requirements are enforced beforehand by go-swagger according to the
		// swagger API specification. They are guaranteed to be valid, so we just use them.

		var secretObject *models.Secret
		var err error

		// Try to put the secret into our secret storage.
		if secretObject, err = secretStorage.Put(params.Secret, params.ExpireAfterViews, params.ExpireAfter); err != nil {
			// Log error.
			log.Println(err)

			// Notify user about the error.
			return secret.NewAddSecretMethodNotAllowed()
		}

		// Return newly created secret.
		return secret.NewAddSecretOK().WithPayload(secretObject)
	})

	// Getting secrets.
	api.SecretGetSecretByHashHandler = secret.GetSecretByHashHandlerFunc(func(params secret.GetSecretByHashParams) middleware.Responder {
		// params verification and requirements are enforced beforehand by go-swagger according to the
		// swagger API specification. They are guaranteed to be valid, so we just use them.

		var secretObject *models.Secret
		var err error

		// Try to get the secret from our storage.
		if secretObject, err = secretStorage.Get(params.Hash); err != nil {
			// Log error.
			log.Println(err)

			// Notify user about the error (since API specification does not allow to return specific errors, just
			// a secret not found error will be returned).
			return secret.NewGetSecretByHashNotFound()
		}

		// If secret object was retrieved successfully:
		if secretObject != nil {
			// Return it.
			return secret.NewGetSecretByHashOK().WithPayload(secretObject)
		}

		// Log failed attempt to retrieve secret.
		// Can be used by fail2ban to ban offending IPs trying to brute force hashes.
		log.Println("Hash Not Found: " + params.HTTPRequest.RemoteAddr + ": " + params.Hash)

		// Hash not found. Notify user.
		return secret.NewGetSecretByHashNotFound()
	})

	// Server shutdown handler.
	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
//noinspection GoUnusedParameter
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
//noinspection GoUnusedParameter
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
