package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/DavidHuie/gomigrate"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	migrator, _ := gomigrate.NewMigrator(db, gomigrate.Postgres{}, "./migrations")
	err = migrator.Migrate()
	if err != nil {
		log.Fatal(err)
	}
}
