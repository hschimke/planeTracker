package database

import "github.com/hschimke/planeTracker/internal/data/model"

type FlightDatabase interface {
	GetFlightsForUser(model.User)
	AddFlight(model.Flight)
	DeleteFlight(model.Flight)
	UpdateFlight(model.Flight)
	GetUser(model.User)
	DeleteUser(model.User)
	UpdateUser(model.User)
}
