package main

import (
	"net/http"
)

type fooHandler struct {
	Message string
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(f.Message))
	if err != nil {
		return
	}
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("bar called"))
	if err != nil {
		return
	}
}

func main() {
	http.Handle("/foo", &fooHandler{Message: "foo called"})
	http.HandleFunc("/bar", barHandler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		return
	}
}
