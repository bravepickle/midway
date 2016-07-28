// reverse proxy functionality
package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewUrl() *url.URL {
	return &url.URL{
		Scheme: targetSchema,
		Host:   targetHost,
	}
}

func ProxyRequest(w http.ResponseWriter, req *http.Request) {
	u := NewUrl()
	p := httputil.NewSingleHostReverseProxy(u)

	p.ServeHTTP(w, req)
}
