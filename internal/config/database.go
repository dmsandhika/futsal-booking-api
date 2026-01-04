package config

import (
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func InitDB() *gorm.DB {
    dsn := os.Getenv("DATABASE_URL")
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect database: " + err.Error())
    }
    return db
}