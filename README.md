# Request Midway App

Serves as middleware between requests and responses, providing formatted logging and stubbing capabilities. Can be used as requests debugger or API gateway for Microservices. 

Main purpose: tool for developers to investigate how requests are working or stub some requests. Supports onl TCP protocol: HTTP

Based on (negroni)[https://github.com/urfave/negroni] and fork of (gencurl)[https://github.com/bravepickle/gencurl] packages

## Usage
Command usage

```
$ midway -help

Web middleware app to log, proxy requests etc.

Usage: ./midway [options] [arg]

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
  ./midway config:init

```

Default URL: http://localhost:3000


## Features
- log requests in CLI CURL format
- can work as a reverse proxy

## Notes
- output STDOUT, STDERR streams to screen greatly impacts performance! Avoid under heavy load, use files or redirect streams to files instead
- use carefully filter conditions with RegEx, they affect performance greatly

## TODOs
- logger customizations: output source, naming, disabling etc.
- modes enable-disable
- allow redirect to another paths requests in proxy mode
- customize (change, add) request and response headers
- responses logger (disabled by default)
- separate file for error logging and access logs
- uniquely identify each request with responses, errors and keep those IDs in logs. Uniqueness can be partial
- cover with tests, benchmarks, profiling
- proxy section: disable flag, 404 page by default
- login use, credentials in config
- use SQLite as DB engine, add db:init command
- import/export stub data in GUI
- add to README examples how to setup configs for NGINX & Apache wit Midway App
- add multiplexer hendler: send copy of request to another address asynchoneously, without waiting for response (use goroutines?). Test it output
- add to config optional preconditions to start logging requests and responses, using RegEx
- log time taken for handling requests
- init RegEx expressions once and compile them with reuse
- add, set, remove headers for responses and requests
- init all RegEx rxCond compiled sources in array on initConfig and read from them - check on race conditions if affected

