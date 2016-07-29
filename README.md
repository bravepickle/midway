# Request Midway App

Serves as middleware between requests and responses, providing formatted logging and stubbing capabilities

## Usage
Command usage

```shell
$ midway -help

Web middleware app to log, proxy requests etc.

Usage: midway [options] [arg]

Options:
  -debug
    	Enable debug mode
  -f string
    	Path to config file in YAML format (default "./config.yaml")
  -help
    	Print command usage help

Arguments:
  config:init
	initialize example config for running application. If file exists, then it will be reset to defaults

Example:
  midway config:init

```

*Stubman* - response stubbing functionality
