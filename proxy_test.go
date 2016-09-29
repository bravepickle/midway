// Testing & Benchmarking

package main

import (
	"bufio"
	"bytes"
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
	Config.Log.Request.Disabled = false
	Config.Log.Request.Truncate = true

	Config.Log.Response.Output = `./test_response.log`
	Config.Log.Response.Disabled = false
	Config.Log.Response.Truncate = true
	Config.Log.Response.Conditions.Disabled = true
	// FIXME: probably due to init bootstrap with regexes at the start
	//	Config.Log.Response.Conditions.Disabled = false
	//	Config.Log.Response.Conditions.Request.Disabled = false
	//	Config.Log.Response.Conditions.Request.Uri = "\\.html"

	Config.Log.Error.Output = `./test_error.log`
	Config.Log.Error.Disabled = false
	Config.Log.Error.Truncate = true

	t.Log(`Files to log to:`, Config.Log.Request.Output, Config.Log.Response.Output, Config.Log.Error.Output)

	mux := http.NewServeMux()
	n := Gateway()
	initRouting(mux, n)

	srv := httptest.NewServer(n)
	defer srv.Close()

	t.Log(`URL:`, srv.URL)

	testInitConfigFromSource(t, srv)
}

func checkInitConfigFromSource(t *testing.T) {
	//	initConfigFromSource
}

func testInitConfigFromSource(t *testing.T, srv *httptest.Server) {
	var testData = []struct {
		desc     string
		input    string
		expected ConfigStruct
	}{
		{
			desc:  `defaults`,
			input: appConfigExample,
			expected: ConfigStruct{
				App: AppConfigStruct{
					Host: ``,
					Port: `3000`,
				},
				Proxy: ProxyConfigStruct{
					Disabled: true,
					Scheme:   `http`,
					Host:     `localhost`,
					Port:     `3001`,
				},
				Log: LogConfigStruct{
					Disabled: false,
					Request: RequestLogConfigStruct{
						Disabled: false,
						Output:   ``,
						Truncate: false,
						Conditions: RequestLogCondConfigStruct{
							Disabled: true,
							Uri:      ``,
							Header:   ``,
							Method:   ``,
							Body:     ``,
						},
					},
					Response: ResponseLogConfigStruct{
						Disabled: false,
						Output:   ``,
						Truncate: false,
						Conditions: ResponseLogCondConfigStruct{
							Disabled: false,
							Header:   ``,
							Body:     ``,
							Request: RequestLogCondConfigStruct{
								Disabled: true,
								Uri:      "\\.php(\\?|$)",
								Header:   ``,
								Method:   ``,
								Body:     ``,
							},
						},
					},
					Error: ErrorLogConfigStruct{
						Disabled: false,
						Output:   ``,
						Truncate: false,
					},
				},
			},
		}, {
			desc: `custom`,
			input: `
app:
  host: example.com
  port: 3001

proxy:
  disabled: false
  scheme: https
  host: target.example.com
  port: 10001

log:
  disabled: false
  
  response:
    disabled: false
    output: ./response.log
    truncate: true
    
    conditions:
      disabled: false
      request:
        disabled: false
        uri: "\\.class.php"
        method: "GET"
        header: "X-Client: 123"
        body: "Hello"
      header: "Content-Type: text/plain"
      body: "Result"
    
  request:
    disabled: false
    output: ./request.log
    truncate: true
            
    conditions:
      disabled: false
      uri: "/users/"
      header: ""
      method: "POST"
      body: "MyRequest"
  
  error:
    disabled: false
    output: ./error.log
    truncate: true

`,
			expected: ConfigStruct{
				App: AppConfigStruct{
					Host: `example.com`,
					Port: `3001`,
				},
				Proxy: ProxyConfigStruct{
					Disabled: false,
					Scheme:   `https`,
					Host:     `target.example.com`,
					Port:     `10001`,
				},
				Log: LogConfigStruct{
					Disabled: false,
					Request: RequestLogConfigStruct{
						Disabled: false,
						Output:   `./request.log`,
						Truncate: true,
						Conditions: RequestLogCondConfigStruct{
							Disabled: false,
							Uri:      `/users/`,
							Header:   ``,
							Method:   `POST`,
							Body:     `MyRequest`,
						},
					},
					Response: ResponseLogConfigStruct{
						Disabled: false,
						Output:   `./response.log`,
						Truncate: true,
						Conditions: ResponseLogCondConfigStruct{
							Disabled: false,
							Header:   `Content-Type: text/plain`,
							Body:     `Result`,
							Request: RequestLogCondConfigStruct{
								Disabled: false,
								Uri:      "/users/",
								Header:   `X-Client: 123`,
								Method:   `GET`,
								Body:     `Hello`,
							},
						},
					},
					Error: ErrorLogConfigStruct{
						Disabled: false,
						Output:   `./error.log`,
						Truncate: true,
					},
				},
			},
		},
	}

	for _, data := range testData {
		t.Log("\tStep:", data.desc)
		var actual = &ConfigStruct{}
		buf := bytes.NewBufferString(data.input)

		initConfigFromSource(buf, actual)
		assertConfigEquals(actual, data.expected, t)
	}
}

func assertConfigEquals(actual *ConfigStruct, expected ConfigStruct, t *testing.T) {
	// app
	assertEquals(`App.Host`, actual.App.Host, expected.App.Host, t)
	assertEquals(`App.Port`, actual.App.Port, expected.App.Port, t)

	// proxy
	assertEquals(`Proxy.Disabled`, actual.Proxy.Disabled, expected.Proxy.Disabled, t)
	assertEquals(`Proxy.Scheme`, actual.Proxy.Scheme, expected.Proxy.Scheme, t)
	assertEquals(`Proxy.Host`, actual.Proxy.Host, expected.Proxy.Host, t)
	assertEquals(`Proxy.Port`, actual.Proxy.Port, expected.Proxy.Port, t)

	// log
	assertEquals(`Log.Disabled`, actual.Log.Disabled, expected.Log.Disabled, t)

	assertEquals(`Log.Request.Disabled`, actual.Log.Request.Disabled, expected.Log.Request.Disabled, t)
	assertEquals(`Log.Request.Output`, actual.Log.Request.Output, expected.Log.Request.Output, t)
	assertEquals(`Log.Request.Truncate`, actual.Log.Request.Truncate, expected.Log.Request.Truncate, t)
	assertEquals(`Log.Request.Conditions.Disabled`, actual.Log.Request.Conditions.Disabled, expected.Log.Request.Conditions.Disabled, t)
	assertEquals(`Log.Request.Conditions.Uri`, actual.Log.Request.Conditions.Uri, expected.Log.Request.Conditions.Uri, t)
	assertEquals(`Log.Request.Conditions.Method`, actual.Log.Request.Conditions.Method, expected.Log.Request.Conditions.Method, t)
	assertEquals(`Log.Request.Conditions.Header`, actual.Log.Request.Conditions.Header, expected.Log.Request.Conditions.Header, t)
	assertEquals(`Log.Request.Conditions.Body`, actual.Log.Request.Conditions.Body, expected.Log.Request.Conditions.Body, t)

	assertEquals(`Log.Response.Disabled`, actual.Log.Response.Disabled, expected.Log.Response.Disabled, t)
	assertEquals(`Log.Response.Output`, actual.Log.Response.Output, expected.Log.Response.Output, t)
	assertEquals(`Log.Response.Truncate`, actual.Log.Response.Truncate, expected.Log.Response.Truncate, t)
	assertEquals(`Log.Response.Conditions.Disabled`, actual.Log.Response.Conditions.Disabled, expected.Log.Response.Conditions.Disabled, t)
	assertEquals(`Log.Response.Conditions.Header`, actual.Log.Response.Conditions.Header, expected.Log.Response.Conditions.Header, t)
	assertEquals(`Log.Response.Conditions.Body`, actual.Log.Response.Conditions.Body, expected.Log.Response.Conditions.Body, t)
	assertEquals(`Log.Response.Conditions.Request.Disabled`, actual.Log.Request.Conditions.Disabled, expected.Log.Request.Conditions.Disabled, t)
	assertEquals(`Log.Response.Conditions.Request.Uri`, actual.Log.Request.Conditions.Uri, expected.Log.Request.Conditions.Uri, t)
	assertEquals(`Log.Response.Conditions.Request.Method`, actual.Log.Response.Conditions.Request.Method, expected.Log.Response.Conditions.Request.Method, t)
	assertEquals(`Log.Response.Conditions.Request.Header`, actual.Log.Response.Conditions.Request.Header, expected.Log.Response.Conditions.Request.Header, t)
	assertEquals(`Log.Response.Conditions.Request.Body`, actual.Log.Response.Conditions.Request.Body, expected.Log.Response.Conditions.Request.Body, t)

	assertEquals(`Log.Error.Disabled`, actual.Log.Error.Disabled, expected.Log.Error.Disabled, t)
	assertEquals(`Log.Error.Output`, actual.Log.Error.Output, expected.Log.Error.Output, t)
	assertEquals(`Log.Error.Truncate`, actual.Log.Error.Truncate, expected.Log.Error.Truncate, t)
}

func assertEquals(label string, actual interface{}, expected interface{}, t *testing.T) {
	if actual != expected {
		t.Errorf("\t\tFailed to assert that values are equal for %s: \"%v\" vs. \"%v\"", label, expected, actual)
	} else {
		t.Logf("\t\tSuccessful check for %s", label)
	}
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

		time.Sleep(time.Second / 3) // wait for some time before checking - file data can be not flushed yet. Async write
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
