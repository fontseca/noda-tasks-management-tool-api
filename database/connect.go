package database

import (
	"database/sql"
	"log"
	"noda/config"

	_ "github.com/lib/pq"
)

func ConnectAndGet() *sql.DB {
	dbconfig := config.GetDatabaseConfig()
	db, err := sql.Open("postgres", dbconfig.Conn())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	dbconfig.LogSuccess()
	return db
}
