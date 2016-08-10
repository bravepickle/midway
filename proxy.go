// reverse proxy functionality
package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
)

// statusLines is a cache of Status-Line strings, keyed by code (for
// HTTP/1.1) or negative code (for HTTP/1.0). This is faster than a
// map keyed by struct of two fields. This map's max size is bounded
// by 2*len(statusText), two protocol types for each known official
// status code in the statusText map.
var (
	statusMu    sync.RWMutex
	statusLines = make(map[int]string)
)

func NewUrl() *url.URL {
	return &url.URL{
		Scheme: Config.Proxy.Scheme,
		Host:   Config.Proxy.HostPortString(),
	}
}

func ProxyRequest(w http.ResponseWriter, req *http.Request) {
	u := NewUrl()
	p := httputil.NewSingleHostReverseProxy(u)

	p.ServeHTTP(w, req)
}

// BufferedResponseWriter contains buffering functionality of ResponseWriter
type BufferedResponseWriter struct {
	Request    *http.Request
	Body       *bytes.Buffer
	StatusCode int
	// initial response writer. Must be set
	Instance http.ResponseWriter
}

func (w BufferedResponseWriter) Write(p []byte) (n int, err error) {
	if _, err := w.Body.Write(p); err != nil {
		log.Println(err)
	}

	return w.Instance.Write(p)
}

func (w BufferedResponseWriter) Header() http.Header {
	return w.Instance.Header()
}

func (w BufferedResponseWriter) WriteHeader(status int) {
	w.StatusCode = status
	w.Instance.WriteHeader(status)
}

func (w BufferedResponseWriter) String() string {
	buf := bytes.NewBufferString(w.StatusLine())
	buf.WriteString(w.HeadersString())
	buf.WriteString("\r\n")
	buf.WriteString(w.Body.String())

	return buf.String()
}

// statusLine returns a response Status-Line (RFC 2616 Section 6.1)
// for the given request and response status code.
func (w BufferedResponseWriter) StatusLine() string {
	// Fast path:
	key := w.StatusCode
	proto11 := w.Request.ProtoAtLeast(1, 1)
	if !proto11 {
		key = -key
	}
	statusMu.RLock()
	line, ok := statusLines[key]
	statusMu.RUnlock()
	if ok {
		return line
	}

	// Slow path:
	proto := "HTTP/1.0"
	if proto11 {
		proto = "HTTP/1.1"
	}
	codestring := strconv.Itoa(w.StatusCode)
	text := http.StatusText(w.StatusCode)
	if text == `` {
		text = "status code " + codestring
	}
	line = proto + " " + codestring + " " + text + "\r\n"
	if ok {
		statusMu.Lock()
		defer statusMu.Unlock()
		statusLines[key] = line
	}
	return line
}

func (w BufferedResponseWriter) HeadersString() string {
	buf := &bytes.Buffer{}
	for vKey, vVals := range w.Header() {
		prefix := bytes.NewBufferString(vKey)
		prefix.WriteString(`: `)

		for _, v := range vVals {
			buf.Write(prefix.Bytes())
			buf.WriteString(v)
			buf.WriteString("\r\n")
		}
	}

	return buf.String()
}

func NewBufferedResponseWriter(rw http.ResponseWriter, r *http.Request) *BufferedResponseWriter {
	return &BufferedResponseWriter{Instance: rw, Body: &bytes.Buffer{}, Request: r, StatusCode: http.StatusOK}
}
