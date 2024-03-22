package infrastructure

import (
	"database/sql"
	"fmt"
)

type SqlPostgres interface {
	GetConnection() *sql.DB
}

type sqlPostgresImpl struct {
	master *sql.DB
}

func NewSqlPostgres() SqlPostgres {
	return &sqlPostgresImpl{
		master: connectDb(),
	}
}

func (sql *sqlPostgresImpl) GetConnection() *sql.DB {
	return sql.master 
}

func connectDb() *sql.DB {
	host := "127.0.0.1"
	port := "5432"
	user := "postgres"
	password := "postgres"
	dbname := "postgres"

	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	conn, err := sql.Open("postgres", connection)

	if err != nil {
		return nil
	}

	return conn
}
