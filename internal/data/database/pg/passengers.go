package pg

import (
	"context"

	"github.com/hschimke/planeTracker/internal/data/model"
)

const (
	getUsersPassengersSql     string = "SELECT passenger_id, default_passenger FROM passengers WHERE user_id = $1"
	addPassengerForUserSql    string = "INSERT INTO passengers (user_id, passenger_id, default_passenger) values ($1, $2, $3)"
	removePassengerForUserSql string = "DELETE FROM passengers WHERE user_id = $1 AND passenger_id = $2"
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

func (p *PostgresDatabase) SetDefaultStatusForPassengerOfUeser(ctx context.Context, user model.UserId, passenger model.UserId, defaultStatus bool) error
func (p *PostgresDatabase) AddPassengerToFlight(ctx context.Context, flight model.FlightId, passenger model.UserId) error
func (p *PostgresDatabase) RemovePassengerFromFlight(ctx context.Context, flight model.FlightId, passenger model.UserId) error
func (p *PostgresDatabase) GetFlightsAsPassenger(ctx context.Context, passenger model.UserId) ([]model.Flight, error)
