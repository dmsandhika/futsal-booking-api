package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"futsal-booking/internal/config"
	"futsal-booking/internal/handler"
	"futsal-booking/internal/model"
	"futsal-booking/internal/repository"
	"futsal-booking/internal/router"
)

func main() {
	godotenv.Load()

	db := config.InitDB()
	db.AutoMigrate(&model.Court{})

	
	courtRepo := &repository.CourtRepository{DB: db}
	courtHandler := &handler.CourtHandler{Repo: courtRepo}

	r := gin.Default()
	r.Static("/uploads", "./uploads")
	router.SetupCourtRoutes(r, courtHandler)

	log.Fatal(r.Run(":" + os.Getenv("APP_PORT")))
}