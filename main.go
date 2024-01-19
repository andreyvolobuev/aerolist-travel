package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=aerolist password=aerolist dbname=aerolist port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var trips_to []Trip
	var trips_from []Trip

	userQuery := "574.2004-08-15.1906.2008-08-15"
	splittedQuery := strings.Split(userQuery, ".")

	dep_city_id := splittedQuery[0]
	arr_city_id := splittedQuery[2]
	dep_date := splittedQuery[1]
	arr_date := splittedQuery[3]

	err = findTrips(
		db,
		dep_city_id,
		dep_date,
		arr_city_id,
		arr_date,
		&trips_from,
		&trips_to,
	)
	if err != nil {
		log.Println("Could not find trips")
	}

	var result []FoundTrip
	err = sortTrips(
		dep_city_id,
		arr_city_id,
		&trips_from,
		&trips_to,
		&result,
	)
	if err != nil {
		log.Println("Could not sort trips")
	}

	str, err := json.MarshalIndent(&result, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", string(str))
}
