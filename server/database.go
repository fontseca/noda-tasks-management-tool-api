package server

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	once sync.Once
	db   *sql.DB
)

func connectToDatabase() {
	if db == nil {
		once.Do(func() {
			var err error
			dbconfig := getDatabaseConfig()
			db, err = sql.Open("postgres", dbconfig.Conn())
			if err != nil {
				log.Fatal(err)
			}
			if err := db.Ping(); err != nil {
				log.Fatal(err)
			}
		})
		return
	}
	log.Fatal("already connected to getDatabase")
}

func getDatabase() *sql.DB {
	if db == nil {
		log.Fatal("connection to getDatabase not established yet")
	}
	return db
}
