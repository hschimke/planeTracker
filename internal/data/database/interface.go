package database

import (
	"context"
	"github.com/hschimke/planeTracker/internal/data/model"
)

type FlightDatabase interface {
	GetFlightsForUser(context.Context, model.UserId) ([]model.Flight, error)
	AddFlight(context.Context, model.Flight) error
	DeleteFlight(context.Context, model.Flight) error
	UpdateFlight(context.Context, model.Flight) error
}
