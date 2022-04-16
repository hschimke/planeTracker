package routes

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	Origin      model.AirportCode `json:"origin"`
	Destination model.AirportCode `json:"destination"`
	TailNumber  model.PlaneTail   `json:"tail_number"`
	Date        string            `json:"date"`
	Email       model.UserId      `json:"email"`
}

type GetAllRequest struct {
	User model.UserId `json:"user"`
}

type UpdateFlightReturn AddFlightReturn

type AddFlightReturn struct {
	Id model.FlightId `json:"id"`
}

type DeleteFlightReturn struct {
	Status string `json:"status"`
}

type BulkUploadRequest struct {
	User       model.UserId `json:"user"`
	Type       string       `json:"type"`
	FlightData string       `json:"flight_data"`
}

type BuldUploadResponse struct {
	Flights []model.FlightId `json:"flights"`
}

type PlaneDetailRequest struct {
	Tail model.PlaneTail `json:"tail"`
	User model.UserId    `json:"user"`
}

type PlaneDetailResponse struct {
	Tail    model.PlaneTail `json:"tail"`
	User    model.UserId    `json:"user"`
	Flights []Flight        `json:"flights"`
	Seen    uint64          `json:"seen"`
	Routes  []struct {
		Origin      model.AirportCode `json:"origin"`
		Destination model.AirportCode `json:"destination"`
		Count       uint64            `json:"count"`
	} `json:"routes"`
}

func getAuthedEmail(ctx context.Context) model.UserId {
	email := model.UserId(ctx.Value("email").(string))
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
			Origin:      flight.Origin.ToIATA(),
			Destination: flight.Destination.ToIATA(),
			Date:        flight.Date.Format("2006-01-02"),
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

	date, dateErr := time.Parse("2006-01-02", flight.Date)
	if dateErr != nil {
		http.Error(w, dateErr.Error(), http.StatusInternalServerError)
		return
	}
	id, addErr := s.db.AddFlight(r.Context(), model.Flight{
		Origin:      flight.Origin,
		Destination: flight.Destination,
		Date:        date,
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

func (s *Server) BulkAddFlights(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var data BulkUploadRequest
	var returnData BuldUploadResponse

	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}
	switch data.Type {
	case "shortcut":
		var flights []model.Flight
		scanner := bufio.NewScanner(strings.NewReader(data.FlightData))
		for scanner.Scan() {
			line := scanner.Text()
			split := strings.Fields(line)
			if len(split) != 4 {
				http.Error(w, fmt.Sprintf("malformed line '%s' has %d fields instead of 4", line, len(split)), http.StatusInternalServerError)
				return
			}
			flightDate, dateParseErr := time.Parse("20060102", split[3])
			if dateParseErr != nil {
				http.Error(w, dateParseErr.Error(), http.StatusInternalServerError)
				return
			}

			flights = append(flights, model.Flight{
				TailNumber:  model.PlaneTail(split[0]),
				Origin:      model.AirportCode(split[1]),
				Destination: model.AirportCode(split[2]),
				Date:        flightDate,
				FlightUser:  email,
			})
		}
		if scanErr := scanner.Err(); scanErr != nil {
			http.Error(w, scanErr.Error(), http.StatusInternalServerError)
			return
		}

		for _, flight := range flights {
			id, addErr := s.db.AddFlight(r.Context(), flight)
			if addErr != nil {
				http.Error(w, addErr.Error(), http.StatusInternalServerError)
				return
			}
			returnData.Flights = append(returnData.Flights, id)
		}

		w.Header().Add("ContentType", "application/json")
		json.NewEncoder(w).Encode(&returnData)
	case "csv":
		http.Error(w, "CSV bulk upload type coming soon", http.StatusInternalServerError)
		return
	default:
		http.Error(w, "unknown bulk upload type", http.StatusInternalServerError)
		return
	}
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
	date, dateErr := time.Parse("2006-01-02", flight.Date)
	if dateErr != nil {
		http.Error(w, dateErr.Error(), http.StatusInternalServerError)
		return
	}

	delErr := s.db.DeleteFlight(r.Context(), model.Flight{
		Id:          flight.Id,
		Origin:      flight.Origin,
		Destination: flight.Destination,
		Date:        date,
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

	date, dateErr := time.Parse("2006-01-02", flight.Date)
	if dateErr != nil {
		http.Error(w, dateErr.Error(), http.StatusInternalServerError)
		return
	}

	updateErr := s.db.UpdateFlight(r.Context(), model.Flight{
		Id:          flight.Id,
		Origin:      flight.Origin,
		Destination: flight.Destination,
		Date:        date,
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

func (s *Server) GetPlaneDetail(w http.ResponseWriter, r *http.Request) {
	email := getAuthedEmail(r.Context())

	var request PlaneDetailRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	if request.User != email {
		http.Error(w, "unauthed email used", http.StatusUnauthorized)
		return
	}

	details, detailsErr := s.db.GetTailDetails(r.Context(), request.Tail, request.User)
	if detailsErr != nil {
		http.Error(w, detailsErr.Error(), http.StatusInternalServerError)
		return
	}

	response := PlaneDetailResponse{
		Tail: details.Tail,
		User: details.User,
		Seen: details.Seen,
		Routes: []struct {
			Origin      model.AirportCode "json:\"origin\""
			Destination model.AirportCode "json:\"destination\""
			Count       uint64            "json:\"count\""
		}(details.Routes),
	}

	for _, flight := range details.Flights {
		response.Flights = append(response.Flights, Flight{
			Id:          flight.Id,
			Origin:      flight.Origin.ToIATA(),
			Destination: flight.Destination.ToIATA(),
			Date:        flight.Date.Format("2006-01-02"),
			TailNumber:  flight.TailNumber,
			Email:       email,
		})
	}

	json.NewEncoder(w).Encode(response)
}
