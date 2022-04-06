package pg

import (
	"context"
	"github.com/hschimke/planeTracker/internal/data/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	getUserFlightsSql string = "SELECT * FROM flights WHERE user_id = $1"
	addFlightSql      string = "INSERT INTO flights("
)

type PostgresDatabase struct {
	db *pgxpool.Pool
}

func NewPostgresDatabase(connectionString string) *PostgresDatabase {
	pool, poolErr := pgxpool.Connect(context.Background(), connectionString)
	if poolErr != nil {
		panic(poolErr.Error())
	}
	return &PostgresDatabase{
		db: pool,
	}
}

func setupDatabase(pool *pgxpool.Pool) {
	//TODO setup database (tables index etc)
}

func (p *PostgresDatabase) GetFlightsForUser(user model.User) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) AddFlight(flight model.Flight) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) DeleteFlight(flight model.Flight) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) UpdateFlight(flight model.Flight) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) GetUser(user model.User) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) DeleteUser(user model.User) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDatabase) UpdateUser(user model.User) {
	//TODO implement me
	panic("implement me")
}
