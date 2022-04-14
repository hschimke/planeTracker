package model

import (
	"strings"
	"time"

	"github.com/hschimke/planeTracker/internal/data/airports"
)

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

func (flight *Flight) Capitalize() {
	flight.Origin = AirportCode(strings.ToUpper(string(flight.Origin)))
	flight.Destination = AirportCode(strings.ToUpper(string(flight.Destination)))
	flight.TailNumber = PlaneTail(strings.ToUpper(string(flight.TailNumber)))
}

func (flight *Flight) Normalize() {
	flight.Capitalize()

	flight.Origin = conformCodeIfPossible(flight.Origin)
	flight.Destination = conformCodeIfPossible(flight.Destination)
}

func conformCodeIfPossible(code AirportCode) AirportCode {
	airportData := airports.AirporData
	if _, found := airportData.ByICAO[string(code)]; found {
		return code
	}
	if airport, found := airportData.ByIATA[string(code)]; found {
		return AirportCode(airport.ICAO)
	}
	return code
}

func (code AirportCode) ToIATA() AirportCode {
	if IATACode := airports.AirporData.ByICAO[string(code)].IATA; IATACode == "" {
		return code
	}

	return AirportCode(airports.AirporData.ByICAO[string(code)].IATA)
}
