// config-related operations & data
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

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
}

type RequestLogCondConfigStruct struct {
	Disabled bool   `yaml:"disabled,omitempty"`
	Method   string `yaml:"method,omitempty"`
	Uri      string `yaml:"uri,omitempty"`
	Header   string `yaml:"header,omitempty"`
	Body     string `yaml:"body,omitempty"`
}

type ResponseLogConfigStruct struct {
	Disabled   bool                       `yaml:"disabled,omitempty"`
	Output     string                     `yaml:"output,omitempty"`
	Conditions RequestLogCondConfigStruct `yaml:"conditions,omitempty"`
}

type ResponseLogCondConfigStruct struct {
	Disabled bool   `yaml:"disabled,omitempty"`
	Uri      string `yaml:"uri,omitempty"`
	Header   string `yaml:"header,omitempty"`
	Body     string `yaml:"body,omitempty"`
}

// LogConfigStruct contains DB connection details
type LogConfigStruct struct {
	Disabled bool                    `yaml:"disabled,omitempty"`
	Request  RequestLogConfigStruct  `yaml:"request,omitempty"`
	Response ResponseLogConfigStruct `yaml:"response,omitempty"`
	ErrorLog string                  `yaml:"error_log,omitempty"`
}

// String generates string
func (t *LogConfigStruct) String() string {
	return fmt.Sprintf(`{error: %s, request: %s, response: %s}`, t.ErrorLog, t.Request.Output, t.Response.Output)
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
