package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	var err error
	dsn := "host=localhost user=postgres password=postgres dbname=statuspage port=5432 sslmode=disable"

	// Überprüfe auf Umgebungsvariable für die Datenbankverbindung
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-Migrate der Modelle
	err = db.AutoMigrate(&Service{}, &Incident{}, &IncidentUpdate{}, &Component{}, &Tag{}, &ComponentTag{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
