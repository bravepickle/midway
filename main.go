package main

import (
	"flag"
	"fmt"
	"net/http"
	//	"path/filepath"

	"github.com/bravepickle/midway/stubman"
	"github.com/urfave/negroni"
)

const defaultConfigPath = `./config.yaml`
const argCfgInit = `config:init`
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

	if !parseAppInput(cfgPath) {
		return
	}

	initConfig(cfgPath, &Config)

	if Debug {
		fmt.Println(`Debug enabled`)
	}

	mux := http.NewServeMux()
	n := Gateway() // Includes some default middlewares

	initStubman(mux, n)

	//	return

	// handle the rest of URIs
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ProxyRequest(w, req)
	})

	n.UseHandler(mux)

	if Debug {
		fmt.Printf("Listening to: %s\n", Config.App.String())
	}

	http.ListenAndServe(Config.App.String(), n)
}

func initStubman(mux *http.ServeMux, n *negroni.Negroni) {
	// init Stubman
	stubman.AddStubmanCrudHandlers(prefixPathStubman, mux)
	// forward all static files to directory

	//	dir := http.Dir(`/var/www/golang/src/github.com/bravepickle/midway` + prefixPathStubman + string(filepath.Separator) + stubman.StaticPath)
	//	dir := http.Dir(prefixPathStubman + string(filepath.Separator) + stubman.StaticPath)
	//	dir := http.Dir(prefixPathStubman + string(filepath.Separator) + stubman.StaticPath)
	//	dir := http.Dir(stubman.StaticPath)
	dir := http.Dir(``)
	//	fmt.Println(`Directory `, dir)
	//	return

	//	n.Use(negroni.NewStatic(http.Dir(fmt.Sprintf(`%s/%s`, prefixPathStubman, stubman.StaticPath))))
	handler := negroni.NewStatic(dir)
	//	handler.Prefix = prefixPathStubman
	n.Use(handler)

	if Debug {
		fmt.Printf("Stubman path: http://%s%s/\n", Config.App.String(), prefixPathStubman)
	}
}
