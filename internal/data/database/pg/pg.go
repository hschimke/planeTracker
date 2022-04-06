package pg

import (
	"context"

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
