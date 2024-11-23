package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const port = "8081"

var tryCounts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// Connect to DB
	dbConnection := connectToDB()
	if dbConnection == nil {
		log.Panic("Can't connect to Postgres!")
	}

	app := Config{
		DB:     dbConnection,
		Models: data.New(dbConnection),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			tryCounts++
		} else {
			log.Println("Connected to postgres db")
			return connection
		}

		if tryCounts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds....")
		time.Sleep(2 * time.Second)
	}
}
