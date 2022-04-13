package model

import "time"

type AirportCode string
type FlightId string
type PlaneTail string

type Flight struct {
	Id          FlightId
	Origin      AirportCode
	Destination AirportCode
	TailNumber  PlaneTail
	Date        time.Time
	FlightUser  UserId
	DateAdded   time.Time
}

type PlaneDetail struct {
	Tail    PlaneTail
	User    UserId
	Flights []Flight
	Seen    uint64
	Routes  []struct {
		Origin      AirportCode
		Destination AirportCode
		Count       uint64
	}
}
