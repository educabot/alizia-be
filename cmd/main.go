package main

import (
	"log"

	"github.com/educabot/alizia-be/config"
	"github.com/educabot/team-ai-toolkit/dbconn"
)

func main() {
	cfg := config.Load()

	db, err := dbconn.NewPostgresConnector(dbconn.PostgresConfig{
		URL: cfg.DatabaseURL,
	}).Connect()
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	app := NewApp(cfg, db)
	defer app.Close()
	app.Run()
}
