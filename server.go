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

// getDB connects to the database
func getDB() *gorm.DB {
	dsn := "host=localhost user=aerolist password=aerolist dbname=aerolist port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

// handleGet handles GET requests to retrieve trips
// @Summary Get Trips
// @Description Get a list of trips based on query parameters
// @Tags Trips
// @Accept  json
// @Produce  json
// @Param depCity query string false "Departure City"
// @Param arrCity query string false "Arrival City"
// @Param depDate query string false "Departure Date"
// @Param arrDate query string false "Arrival Date"
// @Success 200 {array} FoundTrip
// @Failure 400 {string} string "Bad Request"
// @Router / [GET]
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

// handlePost handles POST requests to create a trip
// @Summary Create a new Trip
// @Description Create a new trip by providing trip data in the request body
// @Tags Trips
// @Accept  json
// @Produce  json
// @Param trip body Trip true "Trip data"
// @Success 200 {object} Trip
// @Failure 400 {string} string "Invalid Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router / [POST]
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

// handlePatch handles PATCH requests to update an existing trip
// @Summary Update an existing Trip
// @Description Update a trip by providing trip data and ID
// @Tags Trips
// @Accept  json
// @Produce  json
// @Param id path int true "Trip ID"
// @Param trip body Trip true "Updated trip data"
// @Success 200 {object} Trip
// @Failure 400 {string} string "Trip does not exist"
// @Failure 500 {string} string "Internal Server Error"
// @Router /{id} [PATCH]
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

// handleDelete handles DELETE requests to delete a trip
// @Summary Delete a Trip
// @Description Delete a trip by its ID
// @Tags Trips
// @Param id path int true "Trip ID"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "Trip does not exist"
// @Failure 500 {string} string "Could not delete Trip"
// @Router /{id} [DELETE]
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
