package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASS") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=require" +
		" search_path=" + os.Getenv("DB_SCHEMA")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	
	log.Println("Database connected successfully")
	
	// Ensure schema exists and set search path
	db.Exec("CREATE SCHEMA IF NOT EXISTS " + os.Getenv("DB_SCHEMA"))
	db.Exec("SET search_path TO " + os.Getenv("DB_SCHEMA"))
	
	return db
}
