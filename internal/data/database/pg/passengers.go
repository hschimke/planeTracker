package pg

import (
	"context"
	"fmt"

	"github.com/hschimke/planeTracker/internal/data/model"
)

const (
	getUsersPassengersSql        string = "SELECT passenger_id, default_passenger FROM passengers WHERE user_id = $1"
	addPassengerForUserSql       string = "INSERT INTO passengers (user_id, passenger_id, default_passenger) values ($1, $2, $3)"
	removePassengerForUserSql    string = "DELETE FROM passengers WHERE user_id = $1 AND passenger_id = $2"
	setDefaultPassengerStatusSql string = "UPDATE passengers SET default_passenger = $3 WHERE user_id = $1 AND passenger_id = $2"
	addPassengerToFlightSql      string = "INSERT INTO flight_passengers (flight_id, passenger_id) VALUES ($1, $2)"
	removePassengerFromFlightSql string = "DELETE FROM flight_passengers WHERE flight_id = $1 AND passenger_id = $2"
	getFlightsAsPassengerSql     string = "SELECT id, origin, destination, tail, flight_date, added, count(passenger_id) as cidc FROM flights LEFT JOIN flight_passengers ON flights.id = flight_passengers.flight_id WHERE id IN (SELECT flight_id FROM flight_passengers WHERE passenger_id = $1) GROUP BY id"
	getPassengersForFlightUser   string = "SELECT passenger_id FROM flight_passengers WHERE flight_id = $1"
)

func (p *PostgresDatabase) GetPassengersForUser(ctx context.Context, user model.UserId) ([]model.Passenger, error) {
	query, queryErr := p.db.Query(ctx, getUsersPassengersSql, user)
	if queryErr != nil {
		return nil, queryErr
	}
	defer query.Close()

	var passengerList []model.Passenger

	for query.Next() {
		var passenger model.Passenger
		scanErr := query.Scan(&passenger.PassengerId, &passenger.DefaultCompanion)
		if scanErr != nil {
			return nil, scanErr
		}
		passengerList = append(passengerList, passenger)
	}

	return passengerList, nil
}

func (p *PostgresDatabase) AddPassengerForUser(ctx context.Context, user model.UserId, passenger model.Passenger) error {
	_, insertErr := p.db.Exec(ctx, addPassengerForUserSql, user, passenger.PassengerId, passenger.DefaultCompanion)
	if insertErr != nil {
		return insertErr
	}
	return nil
}

func (p *PostgresDatabase) RemovePassengerForUser(ctx context.Context, user model.UserId, passenger model.UserId) error {
	_, deleteErr := p.db.Exec(ctx, removePassengerForUserSql, user, passenger)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func (p *PostgresDatabase) SetDefaultStatusForPassengerOfUser(ctx context.Context, user model.UserId, passenger model.UserId, defaultStatus bool) error {
	_, updateErr := p.db.Exec(ctx, setDefaultPassengerStatusSql, user, passenger, defaultStatus)
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func (p *PostgresDatabase) AddPassengerToFlight(ctx context.Context, flightId model.FlightId, user model.UserId, passenger model.UserId) error {
	const flightDetailsSql string = "SELECT user_id FROM flights WHERE id = $1"
	// Verify flight user allowed
	vQ := p.db.QueryRow(ctx, flightDetailsSql, flightId)
	var fetchedId model.UserId
	vQScanErr := vQ.Scan(&fetchedId)
	if vQScanErr != nil {
		return vQScanErr
	}

	if fetchedId != user {
		return fmt.Errorf("user and flight owner must match")
	}

	// Add passenger
	_, insertErr := p.db.Exec(ctx, addPassengerToFlightSql, flightId, passenger)
	if insertErr != nil {
		return insertErr
	}
	return nil
}

func (p *PostgresDatabase) RemovePassengerFromFlight(ctx context.Context, flight model.FlightId, user model.UserId, passenger model.UserId) error {
	const flightDetailsSql string = "SELECT user_id FROM flights WHERE id = $1"
	// Verify flight user allowed
	vQ := p.db.QueryRow(ctx, flightDetailsSql, flight)
	var fetchedId model.UserId
	vQScanErr := vQ.Scan(&fetchedId)
	if vQScanErr != nil {
		return vQScanErr
	}

	if fetchedId != user {
		return fmt.Errorf("user and flight owner must match")
	}

	// Remove passenger
	_, deleteErr := p.db.Exec(ctx, removePassengerFromFlightSql, flight, passenger)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func (p *PostgresDatabase) GetFlightsAsPassenger(ctx context.Context, passenger model.UserId) ([]model.Flight, error) {
	flightList, fErr := p.db.Query(ctx, getFlightsAsPassengerSql, passenger)
	if fErr != nil {
		return nil, fErr
	}
	defer flightList.Close()

	var flights []model.Flight
	for flightList.Next() {
		var flight model.Flight
		sErr := flightList.Scan(&flight.Id, &flight.Origin, &flight.Destination, &flight.TailNumber, &flight.Date, &flight.DateAdded, &flight.PassengerCount)
		if sErr != nil {
			return nil, sErr
		}
		flights = append(flights, flight)
	}

	return flights, nil
}

func (p *PostgresDatabase) GetPassengersForFlightUser(ctx context.Context, flight model.FlightId, user model.UserId) ([]model.UserId, error) {
	const flightDetailsSql string = "SELECT user_id FROM flights WHERE id = $1"
	// Verify flight user allowed
	vQ := p.db.QueryRow(ctx, flightDetailsSql, flight)
	var fetchedId model.UserId
	vQScanErr := vQ.Scan(&fetchedId)
	if vQScanErr != nil {
		return nil, vQScanErr
	}

	// Get passengers
	passengerListQ, plqErr := p.db.Query(ctx, getPassengersForFlightUser, flight)
	if plqErr != nil {
		return nil, plqErr
	}

	foundUserAsPassenger := false
	passengerList := make([]model.UserId, 0)
	for passengerListQ.Next() {
		var newPassenger model.UserId
		sErr := passengerListQ.Scan(&newPassenger)
		if sErr != nil {
			return nil, sErr
		}
		passengerList = append(passengerList, newPassenger)
		if newPassenger == user {
			foundUserAsPassenger = true
		}
	}

	if (fetchedId != user) && (!foundUserAsPassenger) {
		return nil, fmt.Errorf("user and flight owner must match, or user must be a passenger on flight")
	}

	return passengerList, nil
}
