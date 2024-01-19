package main

import (
	"time"

	"gorm.io/datatypes"
)

type Country struct {
	ID               uint   `gorm:"primaryKey"`
	Capital          string `gorm:"size:64"`
	CodeFips         string `gorm:"size:2;column:codeFips"`
	CodeIso2Country  string `gorm:"size:2;column:codeIso2Country"`
	CodeIso3Country  string `gorm:"size:3;column:codeIso3Country"`
	Continent        string `gorm:"size:2"`
	Name             string `gorm:"size:64"`
	NumericIso       int    `gorm:"size:2;column:numericIso"`
	Wikipedia_link   string
	Flag             string
	Bounds           datatypes.JSON
	Cases            datatypes.JSON
	NameTranslations datatypes.JSON
}

func (c *Country) TableName() string {
	return "travel_country"
}

type City struct {
	ID           uint   `gorm:"primaryKey"`
	GMT          string `gorm:"column:GMT;size:12"`
	CountryId    int
	Country      Country
	CodeIataCity string `gorm:"size:3;column:codeIataCity"`
	GeonameId    int    `gorm:"column:geonameId"`
	Lat          float64
	Lon          float64
	Admin_name   string `gorm:"size:64"`
	Name         string `gorm:"size:64"`
	Time_zone    string `gorm:"size:64"`
	Cases        datatypes.JSON
	// Departures []Trip `gorm:"foreignKey:DepCityId"`
	// Arrivals   []Trip `gorm:"foreignKey:ArrCityId"`
}

func (c *City) TableName() string {
	return "travel_city"
}

type Trip struct {
	//gorm.Model
	ID        uint `gorm:"primaryKey"`
	UserId    uint
	DepCityId int
	// DepCity        City
	DepartureDate *time.Time
	ArrCityId     int
	// ArrCity        City
	Text       string
	DistanceKm float32
	IsVerified bool
}

func (t *Trip) TableName() string {
	return "travel_trip"
}

type TripViewRequest struct {
	// gorm.Model
	ID           uint `gorm:"primaryKey"`
	IssuerId     int
	TripId       int
	Trip         Trip
	JoinedWith   Trip
	Text         string
	DateApproved *time.Time
	Approved     bool
	// issuer = models.ForeignKey("user.Profile", on_delete=models.CASCADE, related_name="trip_view_requests")
	// trip = models.ForeignKey("Trip", on_delete=models.CASCADE, related_name="view_request")
	// joined_with = models.ForeignKey("TripViewRequest", on_delete=models.CASCADE, related_name="joined", null=True, blank=True)
}

func (t *TripViewRequest) TableName() string {
	return "travel_tripviewrequest"
}
