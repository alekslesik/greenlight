package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Update the routes() method to return a http.Handler instead of a *httprouter.Router.
func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the relevant methods, URL patterns and handler functions for our/v1/healthcheck
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST" respectively
	// router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.requirePermission("movies:read", app.healthcheckHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies",  app.requirePermission("movies:read", app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write",app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read",app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))
	//Register a new user; POST {"name":"", "email":"", "password":""}
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	// Activate a specific user; PUT {"token":""}
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	// Generate a new authentication token; POST {"email": "", "password": ""}'
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Wrap the router with the panic recovery middleware.
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
