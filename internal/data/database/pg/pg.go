package pg

import (
	"context"
	"github.com/hschimke/planeTracker/internal/data/model"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	getUserFlightsSql string = ""
	addFlightSql      string = ""
	updateFlightSql   string = ""
	deleteFlightSql   string = ""
	getUserSql        string = ""
	updateUserSql     string = ""
	deleteUserSql     string = ""
)

const (
	createUserTableSql      string = "CREATE TABLE IF NOT EXISTS users (id UUID, email TEXT, banned BOOLEAN)"
	createUserTableIndexSql string = "CREATE INDEX IF NOT EXISTS users_table_index ON users (id, email, banned)"

	createFlightTableSql      string = "CREATE TABLE IF NOT EXISTS flights (origin VARCHAR(10), destination VARCHAR(10), tail VARCHAR(10), flight_date DATE, user UUID)"
	createFlightTableIndexSql string = "CREATE INDEX IF NOT EXISTS flights_table_index ON flights (origin, destination, tail, flight_date, user)"
)

type PostgresDatabase struct {
	db *pgxpool.Pool
}

func (p *PostgresDatabase) GetFlightsForUser(user model.User) []model.Flight {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) AddFlight(flight model.Flight) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) DeleteFlight(flight model.Flight) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) UpdateFlight(flight model.Flight) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) FillUser(user *model.User) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) GetUser(id model.UserId) model.User {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) DeleteUser(user model.User) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) UpdateUser(user model.User) error {
	//TODO implement me
	panic("implement me")
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
	dbSetupBatch.Queue(createUserTableSql)
	dbSetupBatch.Queue(createUserTableIndexSql)
	results := pool.SendBatch(context.Background(), dbSetupBatch)
	return results.Close()
}
