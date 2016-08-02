package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bravepickle/midway/stubman"
	"github.com/urfave/negroni"
)

const defaultConfigPath = `./config.yaml`
const argCfgInit = `config:init`
const argDbInit = `db:init`
const argDbImport = `db:import`
const prefixPathStubman = `/stubman`

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

	err := initStubman(mux, n)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	//	return

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

	if Debug {
		fmt.Printf("Listening to: %s\n", Config.App.String())
	}

	http.ListenAndServe(Config.App.String(), n)
}

// init Stubman
func initStubman(mux *http.ServeMux, n *negroni.Negroni) error {
	db := stubman.NewDb(Config.Db.DbName, true)
	err := db.Init()
	if err != nil {
		return err
	}

	db.MakeDefault()

	stubman.AddStubmanCrudHandlers(prefixPathStubman, mux)

	// forward all static files to directory
	n.Use(negroni.NewStatic(http.Dir(``)))

	if Debug {
		fmt.Printf("Stubman path: http://%s%s/\n", Config.App.String(), prefixPathStubman)
	}

	return nil
}
