package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hschimke/planeTracker/internal/data/model"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

// SQL commands for CRUD operations
const (
	getUserFlightsSql string = "SELECT id, origin, destination, tail, flight_date, added, count(passenger_id) as cidc FROM flights LEFT JOIN flight_passengers ON flights.id = flight_passengers.flight_id WHERE user_id = $1 GROUP BY id ORDER BY flight_date DESC, added DESC"
	addFlightSql      string = "INSERT INTO flights(id, origin, destination, tail, flight_date, user_id, added) VALUES($1,$2,$3,$4,$5,$6,$7)"
	updateFlightSql   string = "UPDATE flights SET user_id = $2, origin = $3, destination = $4, tail = $5, flight_date = $6 WHERE id = $1"
	deleteFlightSql   string = "DELETE FROM flights WHERE id = $1"
)

// SQL commands to fetch plane details
const (
	getPlaneFlightsForUserSql string = "SELECT id, origin, destination, tail, flight_date, added FROM flights WHERE user_id = $1 AND tail = $2 ORDER BY flight_date DESC, added DESC"
	getPlaneRoutesSql         string = "SELECT origin, destination, count(*) AS total FROM flights WHERE user_id = $1 AND tail = $2 GROUP BY origin, destination"
)

// SQL commands to create tables and indexes (will be executed when NewPostgresDatabase() is called
const (
	createFlightTableSql      string = "CREATE TABLE IF NOT EXISTS flights (id VARCHAR(50), origin VARCHAR(10), destination VARCHAR(10), tail VARCHAR(10), flight_date DATE, added TIMESTAMP WITH TIME ZONE, user_id TEXT, PRIMARY KEY (id))"
	createFlightTableIndexSql string = "CREATE INDEX IF NOT EXISTS flights_table_index ON flights (id, origin, destination, tail, flight_date, user_id, added)"

	createPassengersTableSql string = "CREATE TABLE IF NOT EXISTS passengers (user_id TEXT, passenger_id TEXT, default_passenger BOOLEAN, PRIMARY KEY (user_id, passenger_id))"

	createFlightPassengerTableSql string = "CREATE TABLE IF NOT EXISTS flight_passengers (flight_id VARCHAR(50), passenger_id TEXT)"
	createFlightPassengerIndexSql string = "CREATE INDEX IF NOT EXISTS flight_passengers_index ON flight_passengers (flight_id, passenger_id)"
)

// SQL commands to handle passengers for adding/removing flights
const (
	getDefaultPassengersForUserSql         string = "SELECT passenger_id FROM passengers WHERE user_id = $1 AND default_passenger = true"
	removeAllPassengersForDeletedFlightSql string = "DELETE FROM flight_passengers WHERE flight_id = $1"
)

type PostgresDatabase struct {
	db *pgxpool.Pool
}

func (p *PostgresDatabase) GetFlightsForUser(ctx context.Context, user model.UserId) ([]model.Flight, error) {
	query, queryErr := p.db.Query(ctx, getUserFlightsSql, user)
	if queryErr != nil {
		return nil, queryErr
	}
	defer query.Close()

	var flights []model.Flight
	for query.Next() {
		flight := model.Flight{
			FlightUser: user,
		}
		scanErr := query.Scan(&flight.Id, &flight.Origin, &flight.Destination, &flight.TailNumber, &flight.Date, &flight.DateAdded, &flight.PassengerCount)
		if scanErr != nil {
			return nil, scanErr
		}
		flights = append(flights, flight)
	}
	return flights, nil
}

func (p *PostgresDatabase) AddFlight(ctx context.Context, flight model.Flight, includeDefaultPassengers bool) (model.FlightId, error) {
	if flight.Id == "" {
		flight.Id = model.FlightId(uuid.New().String())
	}

	flight.Normalize()

	// Add default passengers for this user

	_, queryErr := p.db.Exec(ctx, addFlightSql, flight.Id, flight.Origin, flight.Destination, flight.TailNumber, flight.Date, flight.FlightUser, time.Now())
	if queryErr != nil {
		return "", queryErr
	}

	if includeDefaultPassengers {
		defaultUsersQuery, duqErr := p.db.Query(ctx, getDefaultPassengersForUserSql, flight.FlightUser)
		if duqErr != nil {
			return "", duqErr
		}
		defer defaultUsersQuery.Close()

		for defaultUsersQuery.Next() {
			var passenger_id model.UserId
			sErr := defaultUsersQuery.Scan(&passenger_id)
			if sErr != nil {
				return "", sErr
			}
			pAddErr := p.AddPassengerToFlight(ctx, flight.Id, flight.FlightUser, passenger_id)
			if pAddErr != nil {
				return "", pAddErr
			}
		}
	}

	return flight.Id, nil
}

func (p *PostgresDatabase) DeleteFlight(ctx context.Context, flight model.Flight) error {
	_, execErr := p.db.Exec(ctx, deleteFlightSql, flight.Id)
	if execErr != nil {
		return execErr
	}
	_, passengerRemoveErr := p.db.Exec(ctx, removeAllPassengersForDeletedFlightSql, flight.Id)
	if passengerRemoveErr != nil {
		return passengerRemoveErr
	}
	return nil
}

func (p *PostgresDatabase) UpdateFlight(ctx context.Context, flight model.Flight) error {
	_, execErr := p.db.Exec(ctx, updateFlightSql, flight.Id, flight.FlightUser, flight.Origin, flight.Destination, flight.TailNumber, flight.Date)
	return execErr
}

func (p *PostgresDatabase) GetTailDetails(ctx context.Context, tail model.PlaneTail, user model.UserId) (model.PlaneDetail, error) {
	queryFlights, queryErr := p.db.Query(ctx, getPlaneFlightsForUserSql, user, tail)
	if queryErr != nil {
		return model.PlaneDetail{}, queryErr
	}
	defer queryFlights.Close()

	var flights []model.Flight
	for queryFlights.Next() {
		flight := model.Flight{
			FlightUser: user,
		}
		scanErr := queryFlights.Scan(&flight.Id, &flight.Origin, &flight.Destination, &flight.TailNumber, &flight.Date, &flight.DateAdded)
		if scanErr != nil {
			return model.PlaneDetail{}, scanErr
		}
		flights = append(flights, flight)
	}

	queryRoutes, routesErr := p.db.Query(ctx, getPlaneRoutesSql, user, tail)
	if routesErr != nil {
		return model.PlaneDetail{}, routesErr
	}
	defer queryRoutes.Close()

	type statRow struct {
		Origin      model.AirportCode
		Destination model.AirportCode
		Count       uint64
	}

	returnData := model.PlaneDetail{
		Tail:    tail,
		User:    user,
		Flights: flights,
	}

	for queryRoutes.Next() {
		var row statRow
		scanErr := queryRoutes.Scan(&row.Origin, &row.Destination, &row.Count)
		if scanErr != nil {
			return model.PlaneDetail{}, scanErr
		}
		returnData.Seen += row.Count
		returnData.Routes = append(returnData.Routes, row)
	}

	return returnData, nil
}

func NewPostgresDatabase(connectionString string) *PostgresDatabase {
	pool, poolErr := pgxpool.Connect(context.Background(), connectionString)
	if poolErr != nil {
		fmt.Println(poolErr.Error())
		panic(poolErr.Error())
	}
	if setupErr := setupDatabase(pool); setupErr != nil {
		panic(setupErr.Error())
	}
	return &PostgresDatabase{
		db: pool,
	}
}

func setupDatabase(pool *pgxpool.Pool) error {
	dbSetupBatch := &pgx.Batch{}
	dbSetupBatch.Queue(createFlightTableSql)
	dbSetupBatch.Queue(createFlightTableIndexSql)
	dbSetupBatch.Queue(createPassengersTableSql)
	dbSetupBatch.Queue(createFlightPassengerTableSql)
	dbSetupBatch.Queue(createFlightPassengerIndexSql)

	results := pool.SendBatch(context.Background(), dbSetupBatch)
	return results.Close()
}
