package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Database struct {
	Name string
	Host string
	Port string
	User struct {
		Name     string
		Password string
	}
}

type Server struct {
	Port              string
	Host              string
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
}

func (db *Database) Conn() string {
	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		db.User.Name, db.User.Password, db.Host, db.Name)
}

func (db *Database) LogSuccess() {
	log.Printf("Connection established to database \033[0;32m`%s'\033[0m as user \033[0;34m`%s'\033[0m...\n", db.Name, db.User.Name)
}

func GetDatabaseConfig() *Database {
	return &Database{
		Name: getEnv("DB_NAME"),
		Host: getEnv("DB_HOST", "localhost"),
		Port: getEnv("DB_PORT", "5432"),
		User: struct {
			Name     string
			Password string
		}{
			Name:     getEnv("DB_USER"),
			Password: getEnv("DB_USER_PASSWORD"),
		},
	}
}

func GetDatabaseConfigWithValues(
	dbname,
	host,
	port,
	user,
	password string,
) *Database {
	if port == "" {
		port = "5432"
	}

	if host == "" {
		host = "localhost"
	}

	return &Database{
		Name: dbname,
		Host: host,
		Port: port,
		User: struct {
			Name     string
			Password string
		}{
			Name:     user,
			Password: password,
		},
	}
}

func GetServerConfig() *Server {
	return &Server{
		Host:              getEnv("SERVER_HOST", "localhost"),
		Port:              getEnv("SERVER_PORT", "2846"),
		WriteTimeout:      5 * time.Second,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

func getEnv(env string, fallback ...string) string {
	if len(fallback) > 1 {
		log.Fatalf("configuration failed with error: invalid operation getEnvOrDefault(%s, %v) expects 1 or 2 arguments; found %d",
			env, strings.Join(fallback, ", "), len(fallback))
	}
	if val := os.Getenv(env); val != "" {
		return val
	}
	if len(fallback) == 0 {
		log.Fatalf("config failed: not found environment variable `%s'", env)
	}
	return fallback[0]
}
