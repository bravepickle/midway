// Middleware handlers go here
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bravepickle/gencurl"
	"github.com/urfave/negroni"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type CurlLogger negroni.Logger

// NewLogger returns a new Logger instance
func NewLogger() *CurlLogger {
	return &CurlLogger{log.New(os.Stdout, "[CURL] ", 0)}
}

func (l *CurlLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	l.Println(`-----`)
	l.Printf("Started %s %s", r.Method, r.URL.Path)

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	l.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))
	l.Println(gencurl.FromRequest(r))
	l.Println(`-----`)
}

// Classic returns a new Negroni instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging in CURL format
// Static - Static File Serving
func Gateway() *negroni.Negroni {
	return negroni.New(negroni.NewRecovery(), NewLogger())
}
