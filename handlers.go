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
			l.Request = newFileLog(openOrCreateFile(Config.Log.Request.Output, Config.Log.Request.Truncate))
		}
	}

	if !Config.Log.Response.Disabled {
		if Config.Log.Response.Output == `` {
			l.Response = newFileLog(os.Stdout)
		} else {
			l.Response = newFileLog(openOrCreateFile(Config.Log.Response.Output, Config.Log.Response.Truncate))
		}
	}

	if !Config.Log.Error.Disabled {
		if Config.Log.Error.Output == `` {
			l.Error = newFileLog(os.Stderr)
		} else {
			file := openOrCreateFile(Config.Log.Error.Output, Config.Log.Error.Truncate)
			l.Error = newFileLog(file)
			log.SetOutput(file) // reset default logger
		}
	}

	return l
}

// openOrCreateFile Open existing file or create new and return pointer
func openOrCreateFile(path string, truncate bool) (f *os.File) {
	var err error
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(path)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatalln(err)
		}
	} else {
		//		f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_EXCL, 0664)
		f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			log.Fatalln(err)
		}

		if truncate {
			err = os.Truncate(path, 0)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	return f
}

func newFileLog(file *os.File) *log.Logger {
	return log.New(file, `[CURL] `, 0)
}

func allowedToLogRequest(r *http.Request, body string) bool {
	if Config.Log.Request.Disabled {
		return false
	}

	return allowedToLogRequestCond(&Config.Log.Request.Conditions, r, body)
}

func allowedToLogRequestCond(reqCond *RequestLogCondConfigStruct, r *http.Request, body string) bool {
	if reqCond.Disabled {
		return true
	}

	if reqCond.Uri != `` {
		rxCond := regexp.MustCompile(reqCond.Uri)
		if !rxCond.Match([]byte(r.RequestURI)) {
			return false
		}
	}

	if reqCond.Method != `` {
		rxCond := regexp.MustCompile(reqCond.Method)
		if !rxCond.Match([]byte(r.Method)) {
			return false
		}
	}

	if reqCond.Header != `` {
		rxCond := regexp.MustCompile(reqCond.Header)
		if !containsHeader(r.Header, rxCond) {
			return false
		}
	}

	if reqCond.Body != `` {
		rxCond := regexp.MustCompile(reqCond.Body)
		if !rxCond.Match([]byte(body)) {
			return false
		}
	}

	return true
}

func allowedToLogResponse(rw *BufferedResponseWriter, requestBody string) bool {
	if Config.Log.Response.Disabled {
		return false
	}

	// TODO: implement this
	if Config.Log.Response.Conditions.Disabled {
		return true
	}

	if !Config.Log.Response.Conditions.Request.Disabled &&
		!allowedToLogRequestCond(&Config.Log.Response.Conditions.Request, rw.Request, requestBody) {
		return false
	}

	if Config.Log.Response.Conditions.Header != `` {
		rxCond := regexp.MustCompile(Config.Log.Response.Conditions.Header)
		if !containsHeader(rw.Header(), rxCond) {
			return false
		}
	}

	if Config.Log.Response.Conditions.Body != `` {
		rxCond := regexp.MustCompile(Config.Log.Response.Conditions.Body)
		if !rxCond.Match([]byte(rw.Body.String())) {
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
	mw := NewBufferedResponseWriter(rw, r) // TODO: it should not be initialized if response output is disabled

	if Debug {
		body := gencurl.CopyBody(r)
		logRequest := allowedToLogRequest(r, body)

		if logRequest {
			l.Request.Printf("[%d][%s] Started %s %s", idNum, start, r.Method, r.URL.Path)
		}

		next(mw, r)
		//		next(rw, r)

		res := rw.(negroni.ResponseWriter)

		if logRequest {
			l.Request.Printf("[%d] Completed %v %s in %v\n", idNum, res.Status(), http.StatusText(res.Status()), time.Since(start))
			l.Request.Printf("[%d] %s\n", idNum, gencurl.FromRequestWithBody(r, body))
		}

		// TODO: response body can be copied properly without losing send to end user
		if allowedToLogResponse(mw, body) {
			// TODO: status code, headers, body goes here
			l.Response.Printf("[%d] ----\n%s\n\n", idNum, mw.String())
		}

		// TODO: add response status, headers, body in plain format in log

	} else {
		body := gencurl.CopyBody(r)
		next(mw, r)

		if allowedToLogRequest(r, body) {
			l.Request.Printf("[%d][%s] %s\n", idNum, start, gencurl.FromRequestWithBody(r, body))
		}

		// TODO: response body can be copied properly without losing send to end user
		if allowedToLogResponse(mw, body) {
			// TODO: status code, headers, body goes here
			l.Response.Printf("[%d] ----\n%s\n\n", idNum, mw.String())
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
