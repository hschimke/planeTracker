package routes

import (
	"encoding/json"
	"net/http"

	"github.com/hschimke/planeTracker/internal/data/model"
)

/*
func (s *Server) AddFlight(w http.ResponseWriter, r *http.Request) {}
	email := getAuthedEmail(r.Context())

	var flight AddFlightRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&flight)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if flight.Email != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}
*/

type statusResponse struct {
	Status bool `json:"status"`
}

type passenger struct {
	PassengerId      model.UserId `json:"passenger_id"`
	DefaultCompanion bool         `json:"default_companion"`
}

type getPassengersForUserRequest struct {
	User model.UserId `json:"user"`
}
type getPassengersForUserResponse struct {
	Passengers []passenger `json:"passengers"`
}

func (s *Server) GetPassengersForUser(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var request getPassengersForUserRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if request.User != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	passengers, dbErr := s.db.GetPassengersForUser(r.Context(), request.User)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
	}

	passResp := make([]passenger, 0, len(passengers))
	for _, pass := range passengers {
		passResp = append(passResp, passenger{
			PassengerId:      pass.PassengerId,
			DefaultCompanion: pass.DefaultCompanion,
		})
	}

	w.Header().Add("ContentType", "application/json")
	json.NewEncoder(w).Encode(getPassengersForUserResponse{
		Passengers: passResp,
	})
}

type addPassengerForUserRequest struct {
	User      model.UserId `json:"user"`
	Passenger passenger    `json:"passenger"`
}
type addPassengerForUserResponse = statusResponse

func (s *Server) AddPassengerForUser(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var request addPassengerForUserRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if request.User != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	addErr := s.db.AddPassengerForUser(r.Context(), request.User, model.Passenger{PassengerId: request.Passenger.PassengerId, DefaultCompanion: request.Passenger.DefaultCompanion})
	if addErr != nil {
		http.Error(w, addErr.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("ContentType", "application/json")
	json.NewEncoder(w).Encode(addPassengerForUserResponse{
		Status: true,
	})
}

type removePassengerForUserRequest struct {
	User      model.UserId `json:"user"`
	Passenger model.UserId `json:"passenger"`
}
type removePassengerForUserResponse = statusResponse

func (s *Server) RemovePassengerForUser(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var request removePassengerForUserRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if request.User != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	removeErr := s.db.RemovePassengerForUser(r.Context(), request.User, request.Passenger)
	if removeErr != nil {
		http.Error(w, removeErr.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("ContentType", "application/json")
	json.NewEncoder(w).Encode(removePassengerForUserResponse{
		Status: true,
	})
}

type setDefaultStatusForPassengerOfUserRequest struct {
	User          model.UserId `json:"user"`
	Passenger     model.UserId `json:"passenger"`
	DefaultStatus bool         `json:"default_status"`
}
type setDefaultStatusForPassengerOfUserResponse = statusResponse

func (s *Server) SetDefaultStatusForPassengerOfUser(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var request setDefaultStatusForPassengerOfUserRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if request.User != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	dbErr := s.db.SetDefaultStatusForPassengerOfUser(r.Context(), request.User, request.Passenger, request.DefaultStatus)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("ContentType", "application/json")
	json.NewEncoder(w).Encode(setDefaultStatusForPassengerOfUserResponse{
		Status: true,
	})
}

type addPassengerToFlightRequest struct {
	Flight    model.FlightId `json:"flight"`
	User      model.UserId   `json:"user"`
	Passenger model.UserId   `json:"passenger"`
}
type addPassengerToFlightResponse = statusResponse

func (s *Server) AddPassengerToFlight(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var request addPassengerToFlightRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if request.User != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	addErr := s.db.AddPassengerToFlight(r.Context(), request.Flight, request.User, request.Passenger)
	if addErr != nil {
		http.Error(w, addErr.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("ContentType", "application/json")
	json.NewEncoder(w).Encode(addPassengerToFlightResponse{
		Status: true,
	})
}

type removePassengerFromFlightRequest = addPassengerToFlightRequest
type removePassengerFromFlightResponse = addPassengerToFlightResponse

func (s *Server) RemovePassengerFromFlight(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var request removePassengerFromFlightRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if request.User != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	dbErr := s.db.RemovePassengerFromFlight(r.Context(), request.Flight, request.User, request.Passenger)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("ContentType", "application/json")
	json.NewEncoder(w).Encode(removePassengerFromFlightResponse{
		Status: true,
	})
}

type getFlightsAsPassengerRequest struct {
	Passenger model.UserId `json:"passenger"`
}
type getFlightsAsPassengerResponse struct {
	Flights []Flight `json:"flights"`
}

func (s *Server) GetFlightsAsPassenger(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var request getFlightsAsPassengerRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if request.Passenger != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	flights, dbErr := s.db.GetFlightsAsPassenger(r.Context(), request.Passenger)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
	}

	var response getFlightsAsPassengerResponse
	response.Flights = make([]Flight, 0, len(flights))
	for _, flight := range flights {
		response.Flights = append(response.Flights, Flight{
			Id:          flight.Id,
			Origin:      flight.Origin.ToIATA(),
			Destination: flight.Destination.ToIATA(),
			TailNumber:  flight.TailNumber,
			Date:        flight.Date.Format("2006-01-02"),
			Email:       flight.FlightUser,
		})
	}

	w.Header().Add("ContentType", "application/json")
	json.NewEncoder(w).Encode(response)
}
