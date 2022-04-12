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
	getUserFlightsSql string = "SELECT id, origin, destination, tail, flight_date, added FROM flights WHERE user_id = $1 ORDER BY flight_date DESC, added DESC"
	addFlightSql      string = "INSERT INTO flights(id, origin, destination, tail, flight_date, user_id, added) VALUES($1,$2,$3,$4,$5,$6,$7)"
	updateFlightSql   string = "UPDATE flights SET user_id = $2, origin = $3, destination = $4, tail = $5, flight_date = $6 WHERE id = $1"
	deleteFlightSql   string = "DELETE FROM flights WHERE id = $1"
)

// SQL commands to create tables and indexes (will be executed when NewPostgresDatabase() is called
const (
	createFlightTableSql      string = "CREATE TABLE IF NOT EXISTS flights (id VARCHAR(50), origin VARCHAR(10), destination VARCHAR(10), tail VARCHAR(10), flight_date DATE, added TIMESTAMP WITH TIME ZONE, user_id TEXT, PRIMARY KEY (id))"
	createFlightTableIndexSql string = "CREATE INDEX IF NOT EXISTS flights_table_index ON flights (id, origin, destination, tail, flight_date, user_id, added)"
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
		scanErr := query.Scan(&flight.Id, &flight.Origin, &flight.Destination, &flight.TailNumber, &flight.Date, &flight.DateAdded)
		if scanErr != nil {
			return nil, scanErr
		}
		flights = append(flights, flight)
	}
	return flights, nil
}

func (p *PostgresDatabase) AddFlight(ctx context.Context, flight model.Flight) (model.FlightId, error) {
	if flight.Id == "" {
		flight.Id = model.FlightId(uuid.New().String())
	}

	_, queryErr := p.db.Exec(ctx, addFlightSql, flight.Id, flight.Origin, flight.Destination, flight.TailNumber, flight.Date, flight.FlightUser, time.Now())
	return flight.Id, queryErr
}

func (p *PostgresDatabase) DeleteFlight(ctx context.Context, flight model.Flight) error {
	_, execErr := p.db.Exec(ctx, deleteFlightSql, flight.Id)
	return execErr
}

func (p *PostgresDatabase) UpdateFlight(ctx context.Context, flight model.Flight) error {
	_, execErr := p.db.Exec(ctx, updateFlightSql, flight.Id, flight.FlightUser, flight.Origin, flight.Destination, flight.TailNumber, flight.Date)
	return execErr
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

	results := pool.SendBatch(context.Background(), dbSetupBatch)
	return results.Close()
}
