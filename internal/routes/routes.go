package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/hschimke/planeTracker/internal/data/database"
	"github.com/hschimke/planeTracker/internal/data/model"
)

type Server struct {
	db database.FlightDatabase
}

func NewServer(db database.FlightDatabase) *Server {
	return &Server{db: db}
}

type Flight struct {
	Id          model.FlightId    `json:"id"`
	Origin      model.AirportCode `json:"origin,omitempty"`
	Destination model.AirportCode `json:"destination,omitempty"`
	TailNumber  string            `json:"tail_number,omitempty"`
	Date        time.Time         `json:"date"`
	Email       model.UserId      `json:"email"`
}

type GetAllRequest struct {
	User model.UserId `json:"user,omitempty"`
}

type UpdateFlightReturn AddFlightReturn

type AddFlightReturn struct {
	Id model.FlightId `json:"id,omitempty"`
}

type DeleteFlightReturn struct {
	Status string `json:"status,omitempty"`
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

	userFlights, getErr := s.db.GetFlightsForUser(r.Context(), email)
	if getErr != nil {
		http.Error(w, getErr.Error(), http.StatusInternalServerError)
		return
	}

	returnFlights := make([]Flight, 0, len(userFlights))

	for _, flight := range userFlights {
		returnFlights = append(returnFlights, Flight{
			Id:          flight.Id,
			Origin:      flight.Origin,
			Destination: flight.Destination,
			Date:        flight.Date,
			TailNumber:  flight.TailNumber,
			Email:       email,
		})
	}

	w.Header().Add("ContentType", "application/json")
	json.NewEncoder(w).Encode(&returnFlights)
}

func (s *Server) AddFlight(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var flight Flight
	decodeErr := json.NewDecoder(r.Body).Decode(&flight)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if flight.Email != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	id, addErr := s.db.AddFlight(r.Context(), model.Flight{
		Origin:      flight.Origin,
		Destination: flight.Destination,
		Date:        flight.Date,
		FlightUser:  flight.Email,
		TailNumber:  flight.TailNumber,
	})

	if addErr != nil {
		http.Error(w, addErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("ContentType", "application/json")
	json.NewEncoder(w).Encode(AddFlightReturn{
		Id: id,
	})
}

func (s *Server) DeleteFlight(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var flight Flight
	decodeErr := json.NewDecoder(r.Body).Decode(&flight)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if flight.Email != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	delErr := s.db.DeleteFlight(r.Context(), model.Flight{
		Origin:      flight.Origin,
		Destination: flight.Destination,
		Date:        flight.Date,
		FlightUser:  flight.Email,
		TailNumber:  flight.TailNumber,
	})

	w.Header().Add("ContentType", "application/json")

	if delErr != nil {
		json.NewEncoder(w).Encode(DeleteFlightReturn{
			Status: "false",
		})
		http.Error(w, delErr.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(DeleteFlightReturn{
		Status: "true",
	})
}

func (s *Server) UpdateFlight(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var flight Flight
	decodeErr := json.NewDecoder(r.Body).Decode(&flight)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if flight.Email != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	updateErr := s.db.UpdateFlight(r.Context(), model.Flight{
		Origin:      flight.Origin,
		Destination: flight.Destination,
		Date:        flight.Date,
		FlightUser:  flight.Email,
		TailNumber:  flight.TailNumber,
	})

	if updateErr != nil {
		http.Error(w, updateErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("ContentType", "application/json")
	json.NewEncoder(w).Encode(UpdateFlightReturn{
		Id: flight.Id,
	})
}
