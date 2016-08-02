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
type CurlLogger struct {
	Response *log.Logger
	Request  *log.Logger
	Error    *log.Logger
}

// NewLogger returns a new Logger instance. If nil is returned then it should be skipped
func NewLogger() *CurlLogger {
	if Config.Log.Disabled {
		return nil
	}

	l := &CurlLogger{}

	if !Config.Log.Request.Disabled {
		if Config.Log.Request.Output == `` {
			l.Request = newFileLog(os.Stdout)
		} else {
			f, err := os.OpenFile(Config.Log.Request.Output, os.O_APPEND|os.O_WRONLY|os.O_EXCL, 0664)
			if err != nil {
				log.Fatalln(err)
			}

			l.Request = newFileLog(f)
		}
	}

	if !Config.Log.Response.Disabled {
		if Config.Log.Response.Output == `` {
			l.Response = newFileLog(os.Stdout)
		} else {
			f, err := os.OpenFile(Config.Log.Response.Output, os.O_APPEND|os.O_WRONLY|os.O_EXCL, 0664)
			if err != nil {
				log.Fatalln(err)
			}

			l.Response = newFileLog(f)
		}
	}

	if Config.Log.ErrorLog == `` {
		l.Error = newFileLog(os.Stderr)
	} else {
		f, err := os.OpenFile(Config.Log.ErrorLog, os.O_APPEND|os.O_WRONLY|os.O_EXCL, 0664)
		if err != nil {
			log.Fatalln(err)
		}
		l.Error = newFileLog(f)
	}

	return l
}

func newFileLog(file *os.File) *log.Logger {
	return log.New(file, "[CURL] ", 0)
}

func (l *CurlLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	if Debug {
		l.Request.Println(`-----`)
		l.Request.Printf("[%s] Started %s %s", start, r.Method, r.URL.Path)

		body := gencurl.CopyBody(r)
		next(rw, r)

		res := rw.(negroni.ResponseWriter)
		l.Request.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))

		l.Request.Println(gencurl.FromRequestWithBody(r, body))
	} else {
		body := gencurl.CopyBody(r)
		next(rw, r)
		l.Request.Printf(`[%s] %s`, start, gencurl.FromRequestWithBody(r, body))
	}

}

// Classic returns a new Negroni instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging in CURL format
// Static - Static File Serving
func Gateway() *negroni.Negroni {
	return negroni.New(negroni.NewRecovery(), NewLogger())
	//	return negroni.New(NewLogger(), NewProxy())
}
