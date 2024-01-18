package main

import "time"

type Trip struct {
	ID             uint
	User_id        uint
	Departure_id   uint
	Dep_city_id    uint
	Departure_date *time.Time
	Arrival_id     uint
	Arr_city_id    uint
	Arrival_date   *time.Time
	Text           string
	Distance_km    float32
	Is_verified    bool
	Date_created   *time.Time
	Date_posted    *time.Time
	Date_edited    *time.Time
}

func (t *Trip) TableName() string {
	return "travel_trip"
}
