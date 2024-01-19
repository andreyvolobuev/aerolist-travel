package main

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type FoundTrip struct {
	SingleDir bool
	FromTrip  Trip
	ToTrip    Trip
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
	var trips_to_q *gorm.DB
	if q != "" {
		trips_to_q = db.Debug().Where(q, args...)
	}
	if trips_to_q != nil {
		trips_to_q.Order("user_id").Find(trips_to)
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
	db.Debug().Where(q, args...).Order("user_id").Find(trips_from)
	return nil
}

func sortTrips(dep_city_id, arr_city_id string, trips_from, trips_to *[]Trip, result *[]FoundTrip) {
	did, _ := strconv.Atoi(dep_city_id)
	aid, _ := strconv.Atoi(arr_city_id)

	twoCities := (did != 0 && aid != 0)

	if twoCities {
		if len(*trips_to) != 0 {
			for _, tt := range *trips_to {
				if tt.DepCityId == did && tt.ArrCityId == aid {
					*result = append(*result, FoundTrip{SingleDir: false, FromTrip: tt})
				}
				for _, tf := range *trips_from {
					if tf.UserId > tt.UserId {
						break
					} else if tf.UserId > tt.UserId {
						continue
					} else if tf.Departure_date.Compare(*tt.Departure_date) <= 0 {
						*result = append(*result, FoundTrip{SingleDir: false, FromTrip: tf, ToTrip: tt})
					}
				}
			}
		} else {
			for _, tf := range *trips_from {
				if tf.DepCityId == did && tf.ArrCityId == aid {
					*result = append(*result, FoundTrip{SingleDir: false, FromTrip: tf})
				}
			}
		}
	} else {
		if did != 0 {
			for _, tf := range *trips_from {
				*result = append(*result, FoundTrip{SingleDir: true, FromTrip: tf})
			}
		} else if aid != 0 {
			for _, tt := range *trips_to {
				*result = append(*result, FoundTrip{SingleDir: true, ToTrip: tt})
			}
		} else {
			for _, tf := range *trips_from {
				*result = append(*result, FoundTrip{SingleDir: true, FromTrip: tf})
			}
			for _, tt := range *trips_to {
				*result = append(*result, FoundTrip{SingleDir: true, ToTrip: tt})
			}
		}
	}
}
