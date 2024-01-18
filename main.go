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
	var trip Trip
	db.First(&trip, 1)
	str, err := json.MarshalIndent(&trip, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%#v\n", trip)
	fmt.Printf("%s\n", string(str))
}
