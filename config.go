// config-related operations & data
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// DbConfigStruct contains DB connection details for Stubman
type DbConfigStruct struct {
	DbName string `yaml:"dbname,omitempty"`
}

// HostPortString generates string scheme://host:port
func (t *DbConfigStruct) String() string {
	return fmt.Sprintf(`sqlite3://%s`, t.DbName)
}

// TargetConfigStruct struct contains target subsection for app config
type TargetConfigStruct struct {
	Scheme string `yaml:"scheme,omitempty"`
	Port   string `yaml:"port,omitempty"`
	Host   string `yaml:"host,omitempty"`
}

// HostPortString generates string scheme://host:port
func (t *TargetConfigStruct) String() string {
	return fmt.Sprintf(`%s://%s:%s`, t.Scheme, t.Host, t.Port)
}

// HostPortString generates string host:port
func (t *TargetConfigStruct) HostPortString() string {
	return fmt.Sprintf(`%s:%s`, t.Host, t.Port)
}

// ConfigStruct struct contains main application config
type ConfigStruct struct {
	Target TargetConfigStruct `yaml:"target,omitempty"`
	App    AppConfigStruct    `yaml:"app,omitempty"`
	Db     DbConfigStruct     `yaml:"db,omitempty"`
}

func (c *ConfigStruct) String() string {
	return fmt.Sprintf(`{app: %s, target: %s, db: %s}`, c.Target.String(), c.App.String(), c.Db.String())
}

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
