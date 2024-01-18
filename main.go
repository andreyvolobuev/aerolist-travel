package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type FoundTrip struct {
	Direct   bool
	FromTrip Trip
	ToTrip   Trip
}

func findTrips(db *gorm.DB, dep_city_id, dep_date, arr_city_id, arr_date string, trips_from, trips_to *[]Trip) error {
	q := ""
	args := make([]interface{}, 0)
	if dep_city_id != "" {
		q += "dep_city_id = ?"
		args = append(args, dep_city_id)
	}
	if dep_date != "" {
		if q != "" {
			q += " AND "
		}
		q += "departure_date >= ?"
		dep_date_time, err := time.Parse(time.DateOnly, dep_date)
		if err != nil {
			return err
		}
		args = append(args, dep_date_time)
	}
	var sub_q *gorm.DB
	if q != "" {
		sub_q = db.Model(&Trip{}).Distinct("user_id").Where(q, args...)
	}
	q = ""
	args = make([]interface{}, 0)
	if sub_q != nil {
		q += "user_id IN (?)"
		args = append(args, sub_q)
	}
	if arr_city_id != "" {
		if q != "" {
			q += " AND "
		}
		q += "arr_city_id = ?"
		args = append(args, arr_city_id)
	}
	if arr_date != "" {
		if q != "" {
			q += " AND "
		}
		q += "departure_date <= ?"
		arr_date_time, err := time.Parse(time.DateOnly, arr_date)
		if err != nil {
			return err
		}
		args = append(args, arr_date_time)
	}
	var trips_to_q *gorm.DB
	if q != "" {
		trips_to_q = db.Debug().Where(q, args...)
	}
	if trips_to_q != nil {
		trips_to_q.Find(trips_to)
	}
	q = ""
	args = make([]interface{}, 0)
	if trips_to_q != nil {
		q += "user_id IN (?)"
		args = append(args, trips_to_q.Distinct("user_id"))
	}
	if dep_city_id != "" {
		if q != "" {
			q += " AND "
		}
		q += "dep_city_id = ?"
		args = append(args, dep_city_id)
	}
	if dep_date != "" {
		if q != "" {
			q += " AND "
		}
		q += "departure_date >= ?"
		dep_date_time, err := time.Parse(time.DateOnly, dep_date)
		if err != nil {
			return err
		}
		args = append(args, dep_date_time)
	}
	db.Debug().Where(q, args...).Find(trips_from)
	return nil
}

func sortTrips(dep_city_id, arr_city_id string, trips_from, trips_to *[]Trip, result *[]FoundTrip) error {
	did, _ := strconv.Atoi(dep_city_id)
	aid, _ := strconv.Atoi(arr_city_id)
	for _, tf := range *trips_from {
		if tf.DepCityId == did && tf.ArrCityId == aid {
			*result = append(*result, FoundTrip{Direct: true, FromTrip: tf})
			continue
		}
		for _, tt := range *trips_to {
			if tf.Departure_date.Compare(*tt.Departure_date) <= 0 && tf.UserId == tt.UserId {
				*result = append(*result, FoundTrip{Direct: false, FromTrip: tf, ToTrip: tt})
			}
		}
	}
	return nil
}

func main() {
	dsn := "host=localhost user=aerolist password=aerolist dbname=aerolist port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var trips_to []Trip
	var trips_from []Trip

	dep_city_id := "574"
	arr_city_id := "1906"
	dep_date := "2009-08-15"
	arr_date := "2021-08-15"

	err = findTrips(db, dep_city_id, dep_date, arr_city_id, arr_date, &trips_from, &trips_to)
	if err != nil {
		log.Println("Could not find trips")
	}
	var result []FoundTrip
	sortTrips(dep_city_id, arr_city_id, &trips_from, &trips_to, &result)

	str, err := json.MarshalIndent(&result, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", string(str))
}
