package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const targetHost = `kernel.vm:80`
const targetSchema = `http`

func NewUrl(input *url.URL) *url.URL {
	return &url.URL{
		Scheme:   targetSchema,
		Host:     targetHost,
		Opaque:   input.Opaque,
		RawPath:  input.RawPath,
		RawQuery: input.RawQuery,
		Fragment: input.Fragment,
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		u := NewUrl(req.URL)

		p := httputil.NewSingleHostReverseProxy(u)

		p.ServeHTTP(w, req)

		fmt.Fprintf(w, "Welcome to the home page!")

	})

	n := Gateway() // Includes some default middlewares
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
