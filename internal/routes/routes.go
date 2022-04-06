package routes

import "github.com/hschimke/planeTracker/internal/data/database"

type Server struct {
	db *database.FlightDatabase
}
