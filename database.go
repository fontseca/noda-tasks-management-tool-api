package database

import (
	"database/sql"
	"log"
	"noda/config"
	"sync"

	_ "github.com/lib/pq"
)

var (
	once sync.Once
	db   *sql.DB
)

func Connect() {
	if db == nil {
		once.Do(func() {
			var err error
			dbconfig := config.GetDatabaseConfig()
			db, err = sql.Open("postgres", dbconfig.Conn())
			if err != nil {
				log.Fatal(err)
			}
			if err := db.Ping(); err != nil {
				log.Fatal(err)
			}
			dbconfig.LogSuccess()
		})
		return
	}
	log.Fatal("already connected to database")
}

func Get() *sql.DB {
	if db == nil {
		log.Fatal("connection to database not established yet")
	}
	return db
}
