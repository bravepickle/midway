// Middleware handlers go here
package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	//	"net/http/httptest"

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

func allowedToLogRequest(r *http.Request, body string) bool {
	// don't log requests for stubman
	if !Config.Stubman.Disabled {
		rxCond := regexp.MustCompile(`^` + prefixPathStubman)
		if rxCond.Match([]byte(r.URL.Path)) {
			return false
		}
	}

	if Config.Log.Request.Disabled {
		return false
	}

	if Config.Log.Request.Conditions.Disabled {
		return true
	}

	if Config.Log.Request.Conditions.Uri != `` {
		rxCond := regexp.MustCompile(Config.Log.Request.Conditions.Uri)
		if rxCond.Match([]byte(r.URL.Path)) {
			//			log.Println(`-------------- Allowed string `, r.URL.Path)
			return true
		} else {
			//			log.Println(`-------------- Denied string `, r.URL.Path)

			return false
		}
	}

	return true
}

func (l *CurlLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	//	nextRw := httptest.NewRecorder()

	if Debug {
		l.Request.Println(`-----`)
		l.Request.Printf("[%s] Started %s %s", start, r.Method, r.URL.Path)

		body := gencurl.CopyBody(r)
		next(rw, r)
		//		next(nextRw, r)

		res := rw.(negroni.ResponseWriter)
		l.Request.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))

		if allowedToLogRequest(r, body) {
			l.Request.Println(gencurl.FromRequestWithBody(r, body))
		}

	} else {
		body := gencurl.CopyBody(r)
		next(rw, r)
		//		next(nextRw, r)

		if allowedToLogRequest(r, body) {
			l.Request.Printf(`[%s] %s`, start, gencurl.FromRequestWithBody(r, body))
		}
	}

	//	l.Response.Printf("STATUS CODE: %d\n", nextRw.Code)
	//	l.Response.Printf("HEADERS: %s\n", nextRw.Header())
	//	l.Response.Println(`BODY: `, nextRw.Body.String())

	//	rw.Write([]byte(nextRw.Body.String()))

	//	for k, vals := range nextRw.Header() {
	//		for _, v := range vals {
	//			rw.Header().Add(k, v)
	//		}
	//	}

	//	nextRw.Header().Set(`X-TRY`, `true`)
	//	nextRw.Header().Add(`X-TRY2`, `true`)

	//	rw.Header().Set(`X-TRY`, `true`)
	//	rw.Header().Add(`X-TRY2`, `true`)

	//	rw.WriteHeader(nextRw.Code)

	//	nextRw.Flush()
}

//func (l *CurlLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
//	if Config.Log.Response.Disabled {
//		l.serveHTTPBase(rw, r, next)
//	} else {
//		l.serveHTTPResponse(rw, r, next)
//	}

//}

//// serveHTTPResponse serves response and writes response to logs
//func (l *CurlLogger) serveHTTPResponse(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
//	start := time.Now()
//	nextRw := httptest.NewRecorder()

//	if Debug {
//		l.Request.Println(`-----`)
//		l.Request.Printf("[%s] Started %s %s", start, r.Method, r.URL.Path)

//		body := gencurl.CopyBody(r)
//		next(*nextRw, r)

//		res := rw.(negroni.ResponseWriter)
//		l.Request.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))

//		l.Request.Println(gencurl.FromRequestWithBody(r, body))
//	} else {
//		body := gencurl.CopyBody(r)
//		next(*nextRw, r)
//		l.Request.Printf(`[%s] %s`, start, gencurl.FromRequestWithBody(r, body))
//	}

//	// TODO: add condition checkers
//	l.Response.Printf(`Status: %d Body: %s`, *nextRw.Code, *nextRw.Body.String())
//}

// Classic returns a new Negroni instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging in CURL format
// Static - Static File Serving
func Gateway() *negroni.Negroni {
	return negroni.New(negroni.NewRecovery(), NewLogger())
}
