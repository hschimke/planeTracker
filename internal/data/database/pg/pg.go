package pg

import (
	"context"

	"github.com/google/uuid"
	"github.com/hschimke/planeTracker/internal/data/model"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

// SQL commands for CRUD operations
const (
	getUserFlightsSql string = "SELECT origin,destination,tail,flight_date FROM flights WHERE user = $1"
	addFlightSql      string = "INSERT INTO flights(id, origin,destination,tail,flight_date,user) VALUES($1,$2,$3,$4,$5,$6)"
	updateFlightSql   string = "UPDATE flights SET user = $2, origin = $3, destination = $4, tail = $5, flight_date = $6 WHERE id = $1"
	deleteFlightSql   string = "DELETE FROM flights WHERE id = $1"
)

// SQL commands to create tables and indexes (will be executed when NewPostgresDatabase() is called
const (
	createFlightTableSql      string = "CREATE TABLE IF NOT EXISTS flights (id UUID, origin VARCHAR(10), destination VARCHAR(10), tail VARCHAR(10), flight_date DATE, user TEXT, PRIMARY KEY (id))"
	createFlightTableIndexSql string = "CREATE INDEX IF NOT EXISTS flights_table_index ON flights (id, origin, destination, tail, flight_date, user)"
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
		scanErr := query.Scan(&flight.Origin, &flight.Destination, &flight.TailNumber, &flight.Date)
		if scanErr != nil {
			return nil, scanErr
		}
		flights = append(flights, flight)
	}
	return flights, nil
}

func (p *PostgresDatabase) AddFlight(ctx context.Context, flight model.Flight) error {
	if flight.Id == "" {
		flight.Id = model.FlightId(uuid.New().String())
	}

	_, queryErr := p.db.Exec(ctx, addFlightSql, flight.Id, flight.Origin, flight.Destination, flight.TailNumber, flight.Date, flight.FlightUser)
	return queryErr
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
