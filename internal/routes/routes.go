package routes

import (
	"context"
	"encoding/json"
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
	Email       string            `json:"email"`
}

type GetAllRequest struct {
	User model.UserId `json:"user,omitempty"`
}

func getAuthedEmail(ctx context.Context) model.UserId {
	email := ctx.Value("email").(model.UserId)
	return email
}

func (s *Server) GetFlightsForUser(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var user GetAllRequest

	decodeErr := json.NewDecoder(r.Body).Decode(&user)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}

	if user.User != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	flights, getErr := s.db.GetFlightsForUser(r.Context(), email)
	// TODO finish up here, doesn't work
}
func (s *Server) AddFlight(w http.ResponseWriter, r *http.Request)    {}
func (s *Server) DeleteFlight(w http.ResponseWriter, r *http.Request) {}
func (s *Server) UpdateFlight(w http.ResponseWriter, r *http.Request) {}
