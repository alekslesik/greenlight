package main

import (
	"net/http"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a fixed-format JSON response from a string. Notice how we're using a raw
	// string literal (enclosed with backticks) so that we can include double-quote
	// characters in the JSON without needing to escape them? We also use the %q verb to
	// wrap the interpolated values in double-quotes.
	// js := `{"status": "available", "enviroment": %q, "version": %q}`
	// js = fmt.Sprintf(js, app.config.env, version)

	// Create a map which holds the information that we want to send in the response.
	data := map[string]string{
		"status":     "available",
		"enviroment": app.config.env,
		"version":    version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
}
