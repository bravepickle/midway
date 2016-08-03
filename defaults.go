// string defaults for configs are kept here

package main

const appConfigExample = `# common application settings
app:
  host: null # any host will be checked
  port: 3000 # port to listen to

# proxy server to proxy to requests in gateway mode
proxy:
  disabled: true
  scheme: http
  host: localhost
  port: 80

# Logging settings
log:
  disabled: false
  
  # logging response
  response:
    disabled: false
    # if blank string then will use Stdout
    #output: ./response.log 
    output: ""
    
    # RegEx strings that specify conditions when response should be logged. If contain empty string then allow all. 
    # Filters use AND as glue between conditionals
    conditions:
      disabled: true # disable conditions check
      uri: ""
      header: ""
      body: ""
    
  # logging request
  request:
    disabled: false
    #output: ./request.log
    output: ""
    conditions:
      disabled: true # disable conditions check
      uri: ""
      header: ""
      method: ""
      body: ""
  
  # if blank string then will use Stderr
  #error_log: ./error.log
  error_log: ""

stubman:
  disabled: false
  # DB settings
  db:
    dbname: ./data.sqlite
`
