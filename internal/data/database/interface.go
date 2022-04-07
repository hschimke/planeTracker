package database

import "github.com/hschimke/planeTracker/internal/data/model"

type FlightDatabase interface {
	GetFlightsForUser(model.UserId) ([]model.Flight, error)
	AddFlight(model.Flight) error
	DeleteFlight(model.Flight) error
	UpdateFlight(model.Flight, model.Flight) error
}
