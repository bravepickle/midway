// Testing & Benchmarking Config

package main

import (
	"bytes"
	"net/http/httptest"
	"testing"
)

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
