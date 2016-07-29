package main

import (
	"flag"
	"fmt"
	"net/http"
)

const defaultConfigPath = `./config.yaml`
const argCfgInit = `config:init`

var cfgPath string
var Debug bool
var Config ConfigStruct

func init() {
	flag.BoolVar(&Debug, `debug`, false, `Enable debug mode`)
	flag.StringVar(&cfgPath, `f`, defaultConfigPath, `Path to config file in YAML format`)
}

func main() {
	flag.Parse()

	if !parseAppInput(cfgPath) {
		return
	}

	initConfig(cfgPath, &Config)

	if Debug {
		fmt.Println(`Debug enabled`)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ProxyRequest(w, req)
	})

	n := Gateway() // Includes some default middlewares
	n.UseHandler(mux)

	if Debug {
		fmt.Printf("Listening to: \"%s\"\n", Config.App.String())
	}

	http.ListenAndServe(Config.App.String(), n)
}
