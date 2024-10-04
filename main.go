package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	r.HandleFunc("/{id}/", handleDetailData)
	r.HandleFunc("/{id}", handleDetailData)
	r.HandleFunc("/", handleListData)

	http.Handle("/", r)
	http.ListenAndServe("localhost:9999", nil)
}
