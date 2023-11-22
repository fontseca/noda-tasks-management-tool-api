package server

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var logConnState bool = false

type DatabaseConfig struct {
	Name string
	Host string
	Port string
	User struct {
		Name     string
		Password string
	}
}

type serverConfig struct {
	Port              string
	Host              string
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
}

func (db *DatabaseConfig) Conn() string {
	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		db.User.Name, db.User.Password, db.Host, db.Name)
}

func (db *DatabaseConfig) LogSuccess() {
	fmt.Printf("Connection established to database \033[0;32m`%s'\033[0m as user \033[0;34m`%s'\033[0m ...\n", db.Name, db.User.Name)
}

func getDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
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
) *DatabaseConfig {
	if port == "" {
		port = "5432"
	}

	if host == "" {
		host = "localhost"
	}

	return &DatabaseConfig{
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

func getServerConfig() *serverConfig {
	return &serverConfig{
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
