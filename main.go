package main

import (
	"encoding/json"
	"fmt"
	"log"

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

	dep_city_id := ""
	dep_date := "2004-08-15"
	arr_city_id := ""
	arr_date := "2005-08-15"

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
	sortTrips(
		dep_city_id,
		arr_city_id,
		&trips_from,
		&trips_to,
		&result,
	)

	str, err := json.MarshalIndent(&result, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", string(str))
}
