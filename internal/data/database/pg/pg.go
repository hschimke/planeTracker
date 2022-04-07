package pg

import (
	"context"
	"github.com/hschimke/planeTracker/internal/data/model"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	getUserFlightsSql string = "SELECT origin,destination,tail,flight_date FROM flights WHERE user = $1"
	addFlightSql      string = "INSERT INTO flights(origin,destination,tail,flight_date,user) VALUES($1,$2,$3,$4,$5)"
	updateFlightSql   string = "UPDATE flights SET user = $6, origin = $7, destination = $8, tail = $9, flight_date = $10 WHERE user = $1 AND origin = $2 AND destination = $3 AND tail = $4 AND flight_date = $5"
	deleteFlightSql   string = "DELETE FROM flights WHERE user = $1 AND origin = $2 AND destination = $3 AND tail = $4 AND flight_date = $5"
)

const (
	createFlightTableSql      string = "CREATE TABLE IF NOT EXISTS flights (origin VARCHAR(10), destination VARCHAR(10), tail VARCHAR(10), flight_date DATE, user TEXT)"
	createFlightTableIndexSql string = "CREATE INDEX IF NOT EXISTS flights_table_index ON flights (origin, destination, tail, flight_date, user)"
)

type PostgresDatabase struct {
	db *pgxpool.Pool
}

func (p *PostgresDatabase) GetFlightsForUser(user model.UserId) ([]model.Flight, error) {
	query, queryErr := p.db.Query(context.Background(), getUserFlightsSql, user)
	if queryErr != nil {
		return nil, queryErr
	}
	defer query.Close()

	var flights []model.Flight
	for query.Next() {
		flight := model.Flight{
			FlightUser: user,
		}
		scanErr := query.Scan(&flight.Origin, &flight.Destination, &flight.TailNumber, &flight.Date)
		if scanErr != nil {
			return nil, scanErr
		}
		flights = append(flights, flight)
	}
	return flights, nil
}

func (p *PostgresDatabase) AddFlight(flight model.Flight) error {
	_, queryErr := p.db.Exec(context.Background(), addFlightSql, flight.Origin, flight.Destination, flight.TailNumber, flight.Date, flight.FlightUser)
	return queryErr
}

func (p *PostgresDatabase) DeleteFlight(flight model.Flight) error {
	_, execErr := p.db.Exec(context.Background(), deleteFlightSql, flight.FlightUser, flight.Origin, flight.Destination, flight.TailNumber, flight.Date)
	return execErr
}

func (p *PostgresDatabase) UpdateFlight(updated_flight model.Flight, existing_flight model.Flight) error {
	//"UPDATE flights SET user = $6, origin = $7, destination = $8, tail = $9, flight_date = $10 WHERE user = $1 AND origin = $2 AND destination = $3 AND tail = $4 AND flight_date = $5"
	_, execErr := p.db.Exec(context.Background(), updateFlightSql, existing_flight.FlightUser, existing_flight.Origin, existing_flight.Destination, existing_flight.TailNumber, existing_flight.Date, updated_flight.FlightUser, updated_flight.Origin, updated_flight.Destination, updated_flight.TailNumber, updated_flight.Date)
	return execErr
}

func NewPostgresDatabase(connectionString string) *PostgresDatabase {
	pool, poolErr := pgxpool.Connect(context.Background(), connectionString)
	if poolErr != nil {
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

	results := pool.SendBatch(context.Background(), dbSetupBatch)
	return results.Close()
}
