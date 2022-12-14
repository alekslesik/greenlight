package main

import (
	"fmt"
	"net/http"
	"time"
	"os"
	"os/signal"
	"syscall"
)

func (app *application) serve() error {
	// Declare a HTTP server using the same settings as in our main() function.
	// Use the httprouter instance returned by app.routes() as the server handler.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start a background goroutine.
	go func() {
		// Create a quit channel which carries os.Signal values.
		//! We need to use a buffered channel here because signal.Notify()
		// does not wait for a receiver to be available when sending a signal to
		// the quit channel. If we had used a regular (non-buffered) channel
		// here instead, a signal could be ‘missed’ if our quit channel is not
		// ready to receive at the exact moment that the signal is sent. By using
		// a buffered channel, we avoid this problem and ensure that we never
		// miss a signal.
		quit := make(chan os.Signal, 1)

		// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and
		// relay them to the quit channel. Any other signals will not be caught by
		// signal.Notify() and will retain their default behavior.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Read the signal from the quit channel. This code will block until a signal is received.
		s := <-quit

		// Log a message to say that the signal has been caught. Notice that we also
		// call the String() method on the signal to get the signal name and include it
		// in the log entry properties.
		app.logger.PrintInfo("caught signal", map[string]string{
			"signal" : s.String(),
		})

		// Exit the application with a 0 (success) status code.
		os.Exit(0)
	}()

	// Again, we use the PrintInfo() method to write a "starting server" message at the
	// INFO level. But this time we pass a map containing additional properties (the
	// operating environment and server address) as the final parameter.

	// Start the HTTP server.
	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	// Start the server as normal, returning any error.
	return srv.ListenAndServe()
}