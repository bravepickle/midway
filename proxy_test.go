// Testing & Benchmarking

package main

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TODO: test logging

func TestServerDefaultCfgSetup(t *testing.T) {
	t.Log(`>> Test running Server defaults with disabled Reverse Proxy`)
	t.Log(`Proxy: DISABLED`)
	Config.Proxy.Disabled = true
	Config.Log.Request.Output = `./test_request.log`
	Config.Log.Response.Output = `./test_response.log`
	Config.Log.Error.Output = `./test_error.log`

	t.Log(`Files to log to:`, Config.Log.Request.Output, Config.Log.Response.Output, Config.Log.Error.Output)

	mux := http.NewServeMux()
	n := Gateway()
	initRouting(mux, n)

	srv := httptest.NewServer(n)
	defer srv.Close()

	t.Log(`URL:`, srv.URL)

	testServerDisabledProxy(t, srv)
}

func testServerDisabledProxy(t *testing.T, srv *httptest.Server) {
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
		t.Log("\tChecking URI:", data.url)
		resp, err := http.Get(uri)
		if err != nil {
			t.Error("\t\tFailed to read", uri, `:`, err)
		}

		contentType := resp.Header.Get(`Content-Type`)
		if contentType != data.contentType {
			t.Errorf("\t\tFailed to assert MIME types. Expecting: \"%s\", received: \"%s\"",
				data.contentType, contentType)
		}

		if resp.StatusCode != data.code {
			t.Errorf("\t\tFailed to assert response status code. Expecting: %d, received: %d",
				data.code, resp.StatusCode)
		}

		fh, err := os.Open(Config.Log.Request.Output)
		if err != nil {
			t.Fatal(`Failed to open requests log file:`, err)
		}

		scanner := bufio.NewScanner(fh)
		found := false
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), data.url) {
				found = true
				break
			}
		}

		if !found {
			t.Error(`Failed to find in logs string:`, data.url)
		}
	}
}
