package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/urfave/negroni"
)

const defaultConfigPath = `./config.yaml`
const argCfgInit = `config:init`

var optHelp bool
var cfgPath string
var Debug bool
var ProfileHost string
var Config ConfigStruct

func init() {
	flag.BoolVar(&Debug, `debug`, false, `Enable debug mode`)
	flag.BoolVar(&optHelp, `help`, false, `Print command usage help`)
	flag.StringVar(&cfgPath, `f`, defaultConfigPath, `Path to config file in YAML format`)
	flag.StringVar(&ProfileHost, `prof`, `localhost:6060`, `Host:port for profiling. Available only for debug mode`)
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

		go func() {
			fmt.Printf("Profiling enabled at: http://%s/debug/pprof/Warning! Under high load should be handled carefully, memory leaks possible\n", ProfileHost)
			log.Fatal(http.ListenAndServe(ProfileHost, nil))
		}()
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
		if !Config.Proxy.Disabled {
			ProxyRequest(w, req)
		} else {
			w.Header().Add(`X-Default-Page`, `true`)
			w.Write([]byte(fmt.Sprintf("Request %s was received at %s\n", req.URL.String(), time.Now().String())))
		}
	})

	n.UseHandler(mux)
}
