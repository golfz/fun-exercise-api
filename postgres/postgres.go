package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	ErrDBEnvNotSet = errors.New("environment variables not set correctly for database connection")
)

type Postgres struct {
	Db *sql.DB
}

func New() (*Postgres, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal(ErrDBEnvNotSet)
		return nil, ErrDBEnvNotSet
	}
	databaseSource := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", databaseSource)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return &Postgres{Db: db}, nil
}
