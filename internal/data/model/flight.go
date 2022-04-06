package model

import "time"

type AirportCode string

type Flight struct {
	Origin      AirportCode
	Destination AirportCode
	TailNumber  string
	Date        time.Time
	FlightUser  User
}
