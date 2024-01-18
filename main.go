package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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

	dep_city_id := "574"
	arr_city_id := "1906"
	dep_date := "2009-08-15"
	arr_date := "2021-08-15"

	var dep_date_time time.Time
	if dep_date != "" {
		dep_date_time, _ = time.Parse(time.DateOnly, dep_date)
	}

	var arr_date_time time.Time
	if arr_date != "" {
		arr_date_time, _ = time.Parse(time.DateOnly, arr_date)
	}

	var travellers_q *gorm.DB
	if dep_city_id != "" && arr_city_id != "" {
		travellers_q = db.Debug().Model(&Trip{}).Distinct("user_id").Where(
			"dep_city_id = ? AND departure_date >= ?",
			dep_city_id,
			dep_date_time,
		)
	}

	//var final_travellers_q *gorm.DB
	final_travellers_q := db.Debug().Where(
		"user_id IN (?) AND arr_city_id = ? AND departure_date <= ?",
		travellers_q,
		arr_city_id,
		arr_date_time,
	)
	final_travellers_q.Find(&trips_to)
	db.Debug().Where(
		"user_id IN (?) AND dep_city_id = ? AND departure_date >= ?",
		final_travellers_q.Distinct("user_id"),
		dep_city_id,
		dep_date_time,
	).Find(&trips_from)

	fmt.Printf("%+v\n", trips_to)
	fmt.Println("____---_____", dep_date, travellers_q, final_travellers_q)
	fmt.Printf("%+v\n", trips_from)

	str, err := json.MarshalIndent(&trips_to, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", string(str))
}
