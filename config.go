// config-related operations & data
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

var rxList map[string]*regexp.Regexp

// ConfigStruct struct contains main application config
type ConfigStruct struct {
	Proxy ProxyConfigStruct `yaml:"proxy,omitempty"`
	App   AppConfigStruct   `yaml:"app,omitempty"`
	Log   LogConfigStruct   `yaml:"log,omitempty"`
}

func (c *ConfigStruct) String() string {
	return fmt.Sprintf(`{app: %s, proxy: %s, log: %s}`, c.Proxy.String(), c.App.String(), c.Log.String())
}

// ProxyConfigStruct struct contains reverse proxy configuration
type ProxyConfigStruct struct {
	Disabled bool   `yaml:"disabled,omitempty"`
	Scheme   string `yaml:"scheme,omitempty"`
	Port     string `yaml:"port,omitempty"`
	Host     string `yaml:"host,omitempty"`
}

// HostPortString generates string scheme://host:port
func (t *ProxyConfigStruct) String() string {
	if t.Disabled {
		return fmt.Sprint(`disabled`)
	} else {
		return fmt.Sprintf(`%s://%s:%s`, t.Scheme, t.Host, t.Port)
	}
}

// HostPortString generates string host:port
func (t *ProxyConfigStruct) HostPortString() string {
	return fmt.Sprintf(`%s:%s`, t.Host, t.Port)
}

type RequestLogConfigStruct struct {
	Disabled   bool                       `yaml:"disabled,omitempty"`
	Output     string                     `yaml:"output,omitempty"`
	Conditions RequestLogCondConfigStruct `yaml:"conditions,omitempty"`
	Truncate   bool                       `yaml:"truncate,omitempty"`
}

type RequestLogCondConfigStruct struct {
	Disabled bool   `yaml:"disabled,omitempty"`
	Method   string `yaml:"method,omitempty"`
	Uri      string `yaml:"uri,omitempty"`
	Header   string `yaml:"header,omitempty"`
	Body     string `yaml:"body,omitempty"`
}

type ResponseLogConfigStruct struct {
	Disabled   bool                        `yaml:"disabled,omitempty"`
	Output     string                      `yaml:"output,omitempty"`
	Conditions ResponseLogCondConfigStruct `yaml:"conditions,omitempty"`
	Truncate   bool                        `yaml:"truncate,omitempty"`
}

type ErrorLogConfigStruct struct {
	Disabled bool   `yaml:"disabled,omitempty"`
	Output   string `yaml:"output,omitempty"`
	Truncate bool   `yaml:"truncate,omitempty"`
}

type ResponseLogCondConfigStruct struct {
	Disabled bool                       `yaml:"disabled,omitempty"`
	Request  RequestLogCondConfigStruct `yaml:"request,omitempty"`
	Header   string                     `yaml:"header,omitempty"`
	Body     string                     `yaml:"body,omitempty"`
}

// LogConfigStruct contains DB connection details
type LogConfigStruct struct {
	Disabled bool                    `yaml:"disabled,omitempty"`
	Request  RequestLogConfigStruct  `yaml:"request,omitempty"`
	Response ResponseLogConfigStruct `yaml:"response,omitempty"`
	Error    ErrorLogConfigStruct    `yaml:"error,omitempty"`
}

// String generates string
func (t *LogConfigStruct) String() string {
	return fmt.Sprintf(`{error: %s, request: %s, response: %s}`, t.Error.Output, t.Request.Output, t.Response.Output)
}

// AppConfigStruct contain common application settings
type AppConfigStruct struct {
	Port string
	Host string
}

func (t *AppConfigStruct) String() string {
	return fmt.Sprintf(`%s:%s`, t.Host, t.Port)
}

// initConfig parses config from file and puts it to config struct
func initConfig(cfgPath string, config *ConfigStruct) bool {
	if cfgPath == `` {
		cfgPath = defaultConfigPath
	}

	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open config: %s. Use %s %s to init config\n",
			cfgPath, os.Args[0], argCfgInit)
		return false
	}

	cfgFileString, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %s\n", err.Error())
		return false
	}

	err = yaml.Unmarshal(cfgFileString, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config: %s\n", err.Error())
		return false
	}

	rxList = make(map[string]*regexp.Regexp)
	prepareRequestCond()
	prepareResponseCond()

	return true
}

func saveToFile(str string, cfgPath string) (bool, error) {
	file, err := os.Create(cfgPath)
	if err != nil {
		return false, err
	}

	_, err = file.WriteString(str)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ---------------------- RegEx expressions preparation

func prepareRequestCond() {
	if Config.Log.Request.Disabled {
		return
	}

	prepareRegExRequestCond(&Config.Log.Request.Conditions)
}

// rxListCompile compile pattern that should be used by rxMatch() later on
func rxListCompile(pattern string) {
	if pattern != `` {
		rxList[pattern] = regexp.MustCompile(pattern)
	}
}

// rxMatch search in precompiled conditions text pattern match
func rxMatch(pattern string, value string) bool {
	if pattern == `` {
		return true
	}

	rxCond, ok := rxList[pattern]
	if !ok {
		return false // if cannot find pattern, then just skip it
	}

	if !rxCond.Match([]byte(value)) {
		return false
	}

	return true
}

func prepareRegExRequestCond(reqCond *RequestLogCondConfigStruct) {
	if reqCond.Disabled {
		return
	}

	rxListCompile(reqCond.Uri)
	rxListCompile(reqCond.Method)
	rxListCompile(reqCond.Header)
	rxListCompile(reqCond.Body)
}

func prepareResponseCond() {
	if Config.Log.Response.Disabled || Config.Log.Response.Conditions.Disabled {
		return
	}

	if !Config.Log.Response.Conditions.Request.Disabled {
		prepareRegExRequestCond(&Config.Log.Response.Conditions.Request)
	}

	rxListCompile(Config.Log.Response.Conditions.Header)
	rxListCompile(Config.Log.Response.Conditions.Body)
}
