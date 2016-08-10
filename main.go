package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

const defaultConfigPath = `./config.yaml`
const argCfgInit = `config:init`

var optHelp bool
var cfgPath string
var Debug bool
var Config ConfigStruct

func init() {
	flag.BoolVar(&Debug, `debug`, false, `Enable debug mode`)
	flag.BoolVar(&optHelp, `help`, false, `Print command usage help`)
	flag.StringVar(&cfgPath, `f`, defaultConfigPath, `Path to config file in YAML format`)
}

func main() {
	flag.Parse()

	if optHelp {
		printAppUsage()
		return
	}

	if !initConfig(cfgPath, &Config) {
		return
	}

	if !parseAppInput(cfgPath) {
		return
	}

	if Debug {
		fmt.Println(`Debug enabled`)
	}

	mux := http.NewServeMux()
	n := Gateway() // Includes some default middlewares

	initRouting(mux, n)

	fmt.Printf("Listening to: %s\n", Config.App.String())
	log.Fatal(http.ListenAndServe(Config.App.String(), n))
}

func initRouting(mux *http.ServeMux, n *negroni.Negroni) {
	// favicon
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, `favicon.ico`)
	})

	// handle the rest of URIs
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		//	mux.HandleFunc("/", func(w BufferedResponseWriter, req *http.Request) {
		if !Config.Proxy.Disabled {
			ProxyRequest(w, req)
		} else {
			w.Header().Add(`X-Default-Page`, `true`)
			w.Write([]byte(fmt.Sprintf("Request %s was received at %s\n", req.URL.String(), time.Now().String())))
		}

		//		fmt.Println(w.Buffer.String())
	})

	n.UseHandler(mux)
}
