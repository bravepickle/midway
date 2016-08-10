// Testing & Benchmarking

package main

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestServerDefaultCfgSetup(t *testing.T) {
	t.Log(`>> Test running Server defaults with disabled Reverse Proxy`)
	t.Log(`Proxy: DISABLED`)
	Config.Proxy.Disabled = true
	Config.Log.Request.Output = `./test_request.log`
	Config.Log.Response.Output = `./test_response.log`
	Config.Log.Error.Output = `./test_error.log`

	Config.Log.Request.Disabled = false
	Config.Log.Request.Truncate = true

	Config.Log.Response.Disabled = false
	Config.Log.Response.Truncate = true
	Config.Log.Response.Conditions.Disabled = false
	Config.Log.Response.Conditions.Request.Disabled = false
	Config.Log.Response.Conditions.Request.Uri = "\\.html"

	Config.Log.Error.Disabled = false
	Config.Log.Error.Truncate = true

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
			url:         `/unknown.html`,
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

		time.Sleep(50000) // wait for some time before checking - file data can be not flushed yet. Async write
		assertFileContains(t, data.url, Config.Log.Request.Output)

		if ok, _ := regexp.MatchString(`\.html`, data.url); ok { // check html extensions only
			assertFileContains(t, `Request `+data.url+` was received`, Config.Log.Response.Output)
		}
	}
}

func assertFileContains(t *testing.T, txt string, path string) {
	fh, err := os.Open(path)
	if err != nil {
		t.Fatal(`Failed to open requests log file:`, err)
	}

	scanner := bufio.NewScanner(fh)
	found := false
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), txt) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf(`Failed to find in file "%s" string: %s`, path, txt)
	}
}

//func TestServerDefaultCfgSetup(t *testing.T) {
//	t.Log(`>> Test running Server defaults with disabled Reverse Proxy`)
//	t.Log(`Proxy: DISABLED`)
//	Config.Proxy.Disabled = true
//	Config.Log.Request.Output = `./test_request.log`
//	Config.Log.Response.Output = `./test_response.log`
//	Config.Log.Error.Output = `./test_error.log`

//	t.Log(`Files to log to:`, Config.Log.Request.Output, Config.Log.Response.Output, Config.Log.Error.Output)

//	mux := http.NewServeMux()
//	n := Gateway()
//	initRouting(mux, n)

//	srv := httptest.NewServer(n)
//	defer srv.Close()

//	t.Log(`URL:`, srv.URL)

//	testServerDisabledProxy(t, srv)
//}

//func testServerDisabledProxy(t *testing.T, srv *httptest.Server) {
//	var testData = []struct {
//		url         string
//		code        int
//		contentType string
//	}{
//		{
//			url:         `/favicon.ico`,
//			code:        200,
//			contentType: `image/x-icon`,
//		},
//		{
//			url:         `/unknown.ico`,
//			code:        200,
//			contentType: `text/plain; charset=utf-8`,
//		},
//	}

//	var uri string
//	for _, data := range testData {
//		uri = srv.URL + data.url
//		t.Log("\tChecking URI:", data.url)
//		resp, err := http.Get(uri)
//		if err != nil {
//			t.Error("\t\tFailed to read", uri, `:`, err)
//		}

//		contentType := resp.Header.Get(`Content-Type`)
//		if contentType != data.contentType {
//			t.Errorf("\t\tFailed to assert MIME types. Expecting: \"%s\", received: \"%s\"",
//				data.contentType, contentType)
//		}

//		if resp.StatusCode != data.code {
//			t.Errorf("\t\tFailed to assert response status code. Expecting: %d, received: %d",
//				data.code, resp.StatusCode)
//		}

//		fh, err := os.Open(Config.Log.Request.Output)
//		if err != nil {
//			t.Fatal(`Failed to open requests log file:`, err)
//		}

//		scanner := bufio.NewScanner(fh)
//		found := false
//		for scanner.Scan() {
//			if strings.Contains(scanner.Text(), data.url) {
//				found = true
//				break
//			}
//		}

//		if !found {
//			t.Error(`Failed to find in logs string:`, data.url)
//		}
//	}
//}
