package model

import "time"

type AirportCode string
type FlightId string

type Flight struct {
	Id          FlightId
	Origin      AirportCode
	Destination AirportCode
	TailNumber  string
	Date        time.Time
	FlightUser  UserId
}
