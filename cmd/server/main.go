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
	db.AutoMigrate(&model.Court{}, &model.Admin{}, &model.Booking{})

	courtRepo := &repository.CourtRepository{DB: db}
	courtHandler := &handler.CourtHandler{Repo: courtRepo}
	adminRepo := &repository.AdminRepository{DB: db}
	authHandler := &handler.AuthHandler{Repo: adminRepo}
	bookingRepo := &repository.BookingRepository{DB: db}
	bookingHandler := &handler.BookingHandler{
		Repo:      bookingRepo,
		CourtRepo: courtRepo,
	}

	r := gin.Default()
	r.Static("/uploads", "./uploads")
	router.SetupCourtRoutes(r, authHandler, courtHandler, bookingHandler)
	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}
