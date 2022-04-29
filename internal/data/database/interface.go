package database

import (
	"context"

	"github.com/hschimke/planeTracker/internal/data/model"
)

type FlightDatabase interface {
	GetFlightsForUser(context.Context, model.UserId) ([]model.Flight, error)
	AddFlight(context.Context, model.Flight, bool) (model.FlightId, error)
	DeleteFlight(context.Context, model.Flight) error
	UpdateFlight(context.Context, model.Flight) error
	GetTailDetails(context.Context, model.PlaneTail, model.UserId) (model.PlaneDetail, error)

	GetPassengersForUser(context.Context, model.UserId) ([]model.Passenger, error)
	AddPassengerForUser(context.Context, model.UserId, model.Passenger) error
	RemovePassengerForUser(context.Context, model.UserId, model.UserId) error
	SetDefaultStatusForPassengerOfUser(context.Context, model.UserId, model.UserId, bool) error
	AddPassengerToFlight(context.Context, model.FlightId, model.UserId, model.UserId) error
	RemovePassengerFromFlight(context.Context, model.FlightId, model.UserId, model.UserId) error
	GetFlightsAsPassenger(context.Context, model.UserId) ([]model.Flight, error)
	GetPassengersForFlightUser(context.Context, model.FlightId, model.UserId) ([]model.UserId, error)
}
