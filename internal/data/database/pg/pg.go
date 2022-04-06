package pg

import (
	"context"
	"github.com/hschimke/planeTracker/internal/data/model"

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
	createUserTableSql      string = ""
	createUserTableIndexSql string = ""

	createFlightTableSql      string = ""
	createFlightTableIndexSql string = ""
)

type PostgresDatabase struct {
	db *pgxpool.Pool
}

func (p PostgresDatabase) GetFlightsForUser(user model.User) []model.Flight {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDatabase) AddFlight(flight model.Flight) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDatabase) DeleteFlight(flight model.Flight) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDatabase) UpdateFlight(flight model.Flight) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDatabase) FillUser(user *model.User) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDatabase) GetUser(id model.UserId) model.User {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDatabase) DeleteUser(user model.User) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDatabase) UpdateUser(user model.User) error {
	//TODO implement me
	panic("implement me")
}

func NewPostgresDatabase(connectionString string) *PostgresDatabase {
	pool, poolErr := pgxpool.Connect(context.Background(), connectionString)
	if poolErr != nil {
		panic(poolErr.Error())
	}
	setupDatabase(pool)
	return &PostgresDatabase{
		db: pool,
	}
}

func setupDatabase(pool *pgxpool.Pool) {
	//TODO setup database (tables index etc)
	panic("implement me")
}
