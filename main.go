package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	METHOD_NOT_ALLOWED = "[%s] METHOD IS NOT ALLOWED!"

	// headers
	CONTENT_TYPE     = "Content-Type"
	APPLICATION_JSON = "application/json"
)

func handleGet(w http.ResponseWriter, r *http.Request) {
	dsn := "host=localhost user=aerolist password=aerolist dbname=aerolist port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var trips_to []Trip
	var trips_from []Trip
	var result []FoundTrip

	dep_city := r.URL.Query().Get("dep_city")
	dep_date := r.URL.Query().Get("dep_date")
	arr_city := r.URL.Query().Get("arr_city")
	arr_date := r.URL.Query().Get("arr_date")

	findTrips(
		db,
		dep_city,
		dep_date,
		arr_city,
		arr_date,
		&trips_from,
		&trips_to,
	)
	sortTrips(
		dep_city,
		arr_city,
		&trips_from,
		&trips_to,
		&result,
	)
	str, err := json.MarshalIndent(&result, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Fprint(w, string(str))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "{\"status\": \"post_ok\"}")
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "{\"status\": \"delete_ok\"}")
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)

	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodDelete:
		handleDelete(w, r)
	default:
		log.Printf(METHOD_NOT_ALLOWED, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", handleHTTP)
	http.ListenAndServe("localhost:9000", nil)
}
