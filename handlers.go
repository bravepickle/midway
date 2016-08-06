// Middleware handlers go here
package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync/atomic"
	"time"
	//	"net/http/httptest"

	"github.com/bravepickle/gencurl"
	"github.com/urfave/negroni"
)

var idLogNum uint64

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
	return log.New(file, `[CURL] `, 0)
}

func allowedToLogRequest(r *http.Request, body string) bool {
	if Config.Log.Request.Disabled {
		return false
	}

	if Config.Log.Request.Conditions.Disabled {
		return true
	}

	if Config.Log.Request.Conditions.Uri != `` {
		rxCond := regexp.MustCompile(Config.Log.Request.Conditions.Uri)
		if !rxCond.Match([]byte(r.RequestURI)) {
			return false
		}
	}

	if Config.Log.Request.Conditions.Method != `` {
		rxCond := regexp.MustCompile(Config.Log.Request.Conditions.Method)
		if !rxCond.Match([]byte(r.Method)) {
			return false
		}
	}

	if Config.Log.Request.Conditions.Header != `` {
		rxCond := regexp.MustCompile(Config.Log.Request.Conditions.Header)
		if !containsHeader(r.Header, rxCond) {
			return false
		}
	}

	if Config.Log.Request.Conditions.Body != `` {
		rxCond := regexp.MustCompile(Config.Log.Request.Conditions.Body)
		if !rxCond.Match([]byte(body)) {
			return false
		}
	}

	return true
}

// containsHeader check if header exists
func containsHeader(headers http.Header, rxCond *regexp.Regexp) bool {
	for vKey, vVals := range headers {
		prefix := bytes.NewBufferString(vKey)
		prefix.WriteString(`: `)

		for _, v := range vVals {
			h := bytes.NewBuffer(prefix.Bytes())
			h.WriteString(v)

			if rxCond.Match(h.Bytes()) {
				return true
			}
		}
	}

	return false
}

func (l *CurlLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	idNum := atomic.AddUint64(&idLogNum, 1)

	if Debug {
		//		l.Request.Println(`-----`)
		l.Request.Printf("[%d][%s] Started %s %s", idNum, start, r.Method, r.URL.Path)

		body := gencurl.CopyBody(r)
		next(rw, r)

		res := rw.(negroni.ResponseWriter)
		l.Request.Printf("[%d] Completed %v %s in %v\n", idNum, res.Status(), http.StatusText(res.Status()), time.Since(start))

		if allowedToLogRequest(r, body) {
			l.Request.Printf("[%d] %s\n", idNum, gencurl.FromRequestWithBody(r, body))
		}

		// TODO: add response status, headers, body in plain format in log

	} else {
		body := gencurl.CopyBody(r)
		next(rw, r)

		if allowedToLogRequest(r, body) {
			l.Request.Printf(`[%d][%s] %s`, idNum, start, gencurl.FromRequestWithBody(r, body))
		}

		// TODO: add response status, headers, body in plain format in log
	}
}

// Classic returns a new Negroni instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging in CURL format
func Gateway() *negroni.Negroni {
	return negroni.New(negroni.NewRecovery(), NewLogger())
}
