package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	DB *pgxpool.Pool
}

func InitDB(dbUrl string) PostgresDB {
	db, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		panic("database does not work: " + err.Error())
	}
	return PostgresDB{DB: db}
}
