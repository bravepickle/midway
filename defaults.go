// string defaults for configs are kept here

package main

const appConfigExample = `# common application settings
app:
  host: localhost # any host will be checked
  port: 3000 # port to listen to

# target server to proxy to requests in gateway mode
target:
  scheme: http
  host: kernel.vm
  port: 80

# DB settings
db:
  dbname: ./data.sqlite
`
