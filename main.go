package main

import (
	"net/http"
)

const targetHost = `kernel.vm:80`
const targetSchema = `http`

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ProxyRequest(w, req)

	})

	n := Gateway() // Includes some default middlewares
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
