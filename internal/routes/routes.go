package routes

import (
	"github.com/hschimke/planeTracker/internal/data/database"
	"github.com/hschimke/planeTracker/internal/data/model"
	"net/http"
	"time"
)

type Server struct {
	db database.FlightDatabase
}

func NewServer(db database.FlightDatabase) *Server {
	return &Server{db: db}
}

type Flight struct {
	Origin      model.AirportCode `json:"origin,omitempty"`
	Destination model.AirportCode `json:"destination,omitempty"`
	TailNumber  string            `json:"tail_number,omitempty"`
	Date        time.Time         `json:"date"`
}

func (s *Server) GetFlightsForUser(w http.ResponseWriter, r *http.Request) {}
func (s *Server) AddFlight(w http.ResponseWriter, r *http.Request)         {}
func (s *Server) DeleteFlight(w http.ResponseWriter, r *http.Request)      {}
func (s *Server) UpdateFlight(w http.ResponseWriter, r *http.Request)      {}
