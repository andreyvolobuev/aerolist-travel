package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func HandleGetTrips(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(str))
}

func main() {

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", HandleGetTrips)
	http.ListenAndServe("localhost:9000", nil)
}
