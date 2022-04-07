package database

import "github.com/hschimke/planeTracker/internal/data/model"

type FlightDatabase interface {
	GetFlightsForUser(model.UserId) []model.Flight
	AddFlight(model.Flight) error
	DeleteFlight(model.Flight) error
	UpdateFlight(model.Flight) error
}
