package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	METHOD_NOT_ALLOWED = "[%s] METHOD IS NOT ALLOWED!"

	// headers
	CONTENT_TYPE     = "Content-Type"
	APPLICATION_JSON = "application/json"
)

var decoder = schema.NewDecoder()

func getDB() *gorm.DB {
	dsn := "host=localhost user=aerolist password=aerolist dbname=aerolist port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	db := getDB()

	var trips_to []Trip
	var trips_from []Trip
	var result []FoundTrip
	var query FindTripQuery

	err := decoder.Decode(&query, r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	findTrips(
		db,
		query,
		&trips_from,
		&trips_to,
	)
	sortTrips(
		query,
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
	db := getDB()

	var trip Trip

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &trip)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = createTrip(db, &trip)
	if err != nil {
		http.Error(w, "Could not create Trip", http.StatusInternalServerError)
	}

	str, err := json.MarshalIndent(&trip, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Fprint(w, string(str))
}

func handlePatch(w http.ResponseWriter, r *http.Request) {
	db := getDB()

	vars := mux.Vars(r)
	id := vars["id"]

	var trip Trip
	result := db.First(&trip, id)

	if result.Error != nil {
		http.Error(w, "Trip does not exist", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &trip)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = updateTrip(db, &trip)
	if err != nil {
		http.Error(w, "Could not update Trip", http.StatusInternalServerError)
	}

	str, err := json.MarshalIndent(&trip, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Fprint(w, string(str))
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	db := getDB()

	vars := mux.Vars(r)
	id := vars["id"]

	var trip Trip
	result := db.First(&trip, id)

	if result.Error != nil {
		http.Error(w, "Trip does not exist", http.StatusBadRequest)
		return
	}

	err := deleteTrip(db, &trip)
	if err != nil {
		http.Error(w, "Could not delete Trip", http.StatusInternalServerError)
	}

	fmt.Fprint(w, string("ok"))
}

func handleListData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)

	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodGet:
		handleGet(w, r)
	default:
		log.Printf(METHOD_NOT_ALLOWED, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleDetailData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)

	switch r.Method {
	case http.MethodPut:
		handlePatch(w, r)
	case http.MethodDelete:
		handleDelete(w, r)
	default:
		log.Printf(METHOD_NOT_ALLOWED, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
