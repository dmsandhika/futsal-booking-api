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
		" options=-csearch_path=" + os.Getenv("DB_SCHEMA")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	
	log.Println("Database connected successfully")
	
	return db

	//migrate -path migrations -database "postgres://postgres:PASS@$localhost:5432/futsal_db?sslmode=require&search_path=public" up
}
