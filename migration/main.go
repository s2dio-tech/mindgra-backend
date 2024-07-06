package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/neo4j"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	host, _ := os.LookupEnv("DB_HOST")
	port, _ := os.LookupEnv("DB_PORT")
	username, _ := os.LookupEnv("DB_USERNAME")
	password, _ := os.LookupEnv("DB_PASSWORD")

	m, err := migrate.New(
		"file://./migrations",
		"neo4j://"+username+":"+password+"@"+host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
