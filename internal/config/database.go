package config

import (
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
	       " sslmode=disable TimeZone=Asia/Jakarta"

       db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
       if err != nil {
	       panic("Failed to connect database")
       }
       return db
}