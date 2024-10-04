package main

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type FindTripQuery struct {
	DepCity string `schema:"dep_city"`
	DepDate string `schema:"dep_date"`
	ArrCity string `schema:"arr_city"`
	ArrDate string `scheme:"arr_date"`
}

type FoundTrip struct {
	SingleDir bool  `json:"singleDir"`
	FromTrip  *Trip `json:"fromTrip"`
	ToTrip    *Trip `json:"toTrip"`
}

func findTrips(db *gorm.DB, query FindTripQuery, trips_from, trips_to *[]Trip) error {
	q := ""
	args := make([]interface{}, 0)
	if query.DepCity != "" {
		q += "dep_city_id = ?"
		args = append(args, query.DepCity)
	}
	if query.DepDate != "" {
		if q != "" {
			q += " AND "
		}
		q += "departure_date >= ?"
		dep_date_time, err := time.Parse(time.DateOnly, query.DepDate)
		if err != nil {
			return err
		}
		args = append(args, dep_date_time)
	}
	if query.ArrDate != "" {
		if q != "" {
			q += " AND "
		}
		q += "departure_date <= ?"
		arr_date_time, err := time.Parse(time.DateOnly, query.ArrDate)
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
	if query.ArrCity != "" {
		if q != "" {
			q += " AND "
		}
		q += "arr_city_id = ?"
		args = append(args, query.ArrCity)
	}
	if query.ArrDate != "" {
		if q != "" {
			q += " AND "
		}
		q += "departure_date <= ?"
		arr_date_time, err := time.Parse(time.DateOnly, query.ArrDate)
		if err != nil {
			return err
		}
		args = append(args, arr_date_time)
	}
	if query.DepDate != "" {
		if q != "" {
			q += " AND "
		}
		q += "departure_date >= ?"
		dep_date_time, err := time.Parse(time.DateOnly, query.DepDate)
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
	if query.DepCity != "" {
		if q != "" {
			q += " AND "
		}
		q += "dep_city_id = ?"
		args = append(args, query.DepCity)
	}
	if query.DepDate != "" {
		if q != "" {
			q += " AND "
		}
		q += "departure_date >= ?"
		dep_date_time, err := time.Parse(time.DateOnly, query.DepDate)
		if err != nil {
			return err
		}
		args = append(args, dep_date_time)
	}
	if query.ArrDate != "" {
		if q != "" {
			q += " AND "
		}
		q += "departure_date <= ?"
		arr_date_time, err := time.Parse(time.DateOnly, query.ArrDate)
		if err != nil {
			return err
		}
		args = append(args, arr_date_time)
	}
	db.Debug().Where(q, args...).Not(trips_to_q).Order("user_id").Find(trips_from)
	return nil
}

func sortTrips(query FindTripQuery, trips_from, trips_to *[]Trip, result *[]FoundTrip) {
	did, _ := strconv.Atoi(query.DepCity)
	aid, _ := strconv.Atoi(query.ArrCity)

	twoCities := (did != 0 && aid != 0)

	if twoCities {
		if len(*trips_to) != 0 {
			for _, tt := range *trips_to {
				if tt.DepCityId == did && tt.ArrCityId == aid {
					*result = append(*result, FoundTrip{SingleDir: false, FromTrip: &tt})
					continue
				}
				for _, tf := range *trips_from {
					if tf.UserId > tt.UserId {
						break
					} else if tf.UserId > tt.UserId {
						continue
					} else if tf.DepartureDate.Compare(*tt.DepartureDate) <= 0 {
						*result = append(*result, FoundTrip{SingleDir: false, FromTrip: &tf, ToTrip: &tt})
					}
				}
			}
		} else {
			for _, tf := range *trips_from {
				if tf.DepCityId == did && tf.ArrCityId == aid {
					*result = append(*result, FoundTrip{SingleDir: false, FromTrip: &tf})
				}
			}
		}
	} else {
		if did != 0 {
			for _, tf := range *trips_from {
				*result = append(*result, FoundTrip{SingleDir: true, FromTrip: &tf})
			}
		} else if aid != 0 {
			for _, tt := range *trips_to {
				*result = append(*result, FoundTrip{SingleDir: true, ToTrip: &tt})
			}
		} else {
			for _, tf := range *trips_from {
				*result = append(*result, FoundTrip{SingleDir: true, FromTrip: &tf})
			}
			for _, tt := range *trips_to {
				*result = append(*result, FoundTrip{SingleDir: true, ToTrip: &tt})
			}
		}
	}
}

func createTrip(db *gorm.DB, trip *Trip) error {
	currentTime := time.Now()

	trip.DateEdited = &currentTime
	trip.DateCreated = &currentTime

	result := db.Model(&Trip{}).Create(&trip)
	return result.Error
}

func updateTrip(db *gorm.DB, trip *Trip) error {
	currentTime := time.Now()
	trip.DateEdited = &currentTime

	result := db.Save(trip)
	return result.Error
}

func deleteTrip(db *gorm.DB, trip *Trip) error {
	result := db.Delete(&Trip{}, trip.ID)
	return result.Error
}
