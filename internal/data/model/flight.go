package model

import "time"

type AirportCode string

type Flight struct {
	Id          string
	Origin      AirportCode
	Destination AirportCode
	TailNumber  string
	Date        time.Time
	FlightUser  UserId
}
