package main

import (
	"net/http"

	_ "Travel/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Trips API
// @version 1.0
// @description API for managing trips

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	r.HandleFunc("/{id}/", handleDetailData)
	r.HandleFunc("/{id}", handleDetailData)
	r.HandleFunc("/", handleListData)

	http.Handle("/", r)
	http.ListenAndServe("localhost:9999", nil)
}
