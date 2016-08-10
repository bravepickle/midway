// Testing & Benchmarking

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkServerWriteToFiles(t *testing.B) {
	t.Log(`>> Test running Server defaults with disabled Reverse Proxy`)
	t.Log(`Proxy: DISABLED`)
	Config.Proxy.Disabled = true
	//	Config.Log.Request.Output = `./test_request.log`
	//	Config.Log.Response.Output = `./test_response.log`
	//	Config.Log.Error.Output = `./test_error.log`

	Config.Log.Request.Disabled = false
	Config.Log.Request.Truncate = true

	Config.Log.Response.Disabled = false
	Config.Log.Response.Truncate = true

	Config.Log.Error.Disabled = false
	Config.Log.Error.Truncate = true

	t.Logf(`Files to log to: "%s", "%s", "%s"`, Config.Log.Request.Output, Config.Log.Response.Output, Config.Log.Error.Output)

	mux := http.NewServeMux()
	n := Gateway()
	initRouting(mux, n)

	srv := httptest.NewServer(n)
	defer srv.Close()

	t.Log(`URL:`, srv.URL)

	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		benchServerDisabledProxy(t, srv)
	}
}

func benchServerDisabledProxy(t *testing.B, srv *httptest.Server) {
	var testData = []struct {
		url         string
		code        int
		contentType string
	}{
		{
			url:         `/favicon.ico`,
			code:        200,
			contentType: `image/x-icon`,
		},
		{
			url:         `/unknown.ico`,
			code:        200,
			contentType: `text/plain; charset=utf-8`,
		},
	}

	var uri string
	for _, data := range testData {
		uri = srv.URL + data.url
		//		t.Log(uri)
		resp, err := http.Get(uri)
		if err != nil {
			t.Fatal("\t\tFailed to read", uri, `:`, err)
		}

		contentType := resp.Header.Get(`Content-Type`)
		if contentType != data.contentType {
			t.Fatalf("\t\tFailed to assert MIME types. Expecting: \"%s\", received: \"%s\"",
				data.contentType, contentType)
		}

		if resp.StatusCode != data.code {
			t.Fatalf("\t\tFailed to assert response status code. Expecting: %d, received: %d",
				data.code, resp.StatusCode)
		}
	}
}
