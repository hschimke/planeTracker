package database

import "github.com/hschimke/planeTracker/internal/data/model"

type FlightDatabase interface {
	GetFlightsForUser(model.User) []model.Flight
	AddFlight(model.Flight) error
	DeleteFlight(model.Flight) error
	UpdateFlight(model.Flight) error
	FillUser(*model.User)
	GetUser(model.UserId) model.User
	DeleteUser(model.User) error
	UpdateUser(model.User) error
}
