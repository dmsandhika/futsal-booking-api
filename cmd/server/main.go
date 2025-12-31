package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"futsal-booking/internal/config"
	"futsal-booking/internal/model"
)

func main() {
	godotenv.Load()

	db := config.InitDB()
	db.AutoMigrate(
		model.Court{},
	)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Fatal(r.Run(":" + os.Getenv("APP_PORT")))
}