# Request Midway App

Serves as middleware between requests and responses, providing formatted logging and stubbing capabilities. Can be used as requests debugger or API gateway for Microservices. 

Main purpose: tool for developers to investigate how requests are working or stub some requests. Supports onl TCP protocol: HTTP

Based on (negroni)[https://github.com/urfave/negroni] and fork of (gencurl)[https://github.com/bravepickle/gencurl] packages

## Usage
Command usage

```
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
  -httptest.serve string
    	if non-empty, httptest.NewServer serves on this address and blocks

Arguments:
  config:init
    	initialize example config for running application. If file exists, then it will be reset to defaults
  stubman:db:init
    	initialize DB. If it exists, then DB will be reset
  stubman:db:import <file.sql>
    	import data from SQL file to DB. Second argument must be present with file path

Example:
  midway config:init

```

Default URL: http://localhost:3000

It is optimize-wise to serve static files by web server before Midway App, especially in high-load project. /stubman/static should be served by Midway app

## Features
- log requests in CLI CURL format
- can work as a reverse proxy
- can work as stub service using Stubman
- uses (SQLite)[http://github.com/mattn/go-sqlite3] as DB storage for stubman

## Stubman
*Stubman* - response stubbing functionality with GUI interface
http://localhost:3000/stubman - entrypoint for handling stub service requests and responses


## TODOs
- logger customizations: output source, naming, disabling etc.
- modes enable-disable
- allow redirect to another paths requests in proxy mode
- customize (change, add) request and response headers
- responses logger (disabled by default)
- separate file for error logging and access logs
- uniquely identify each request with responses, errors and keep those IDs in logs. Uniqueness can be partial
- cover with tests, benchmarks, profiling
- flag to disable stubman in whole and GUI side only
- proxy section: disable flag, 404 page by default
- switch stubman GUI themes 
- login use, credentials in config
- use SQLite as DB engine, add db:init command
- import/export stub data in GUI
- add to README examples how to setup configs for NGINX & Apache wit Midway App
- add multiplexer hendler: send copy of request to another address asynchoneously, without waiting for response (use goroutines?). Test it output
- add to config optional preconditions to start logging requests and responses, using RegEx
- in Stubman add button to generate CURL request for given stub
- log time taken for handling requests

