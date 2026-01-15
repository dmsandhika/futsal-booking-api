package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"futsal-booking/internal/config"
	"futsal-booking/internal/handler"
	"futsal-booking/internal/repository"
	"futsal-booking/internal/router"
	"futsal-booking/internal/scheduler"
)

func main() {
	godotenv.Load()

	db := config.InitDB()

	courtRepo := &repository.CourtRepository{DB: db}
	courtHandler := &handler.CourtHandler{Repo: courtRepo}
	adminRepo := &repository.AdminRepository{DB: db}
	authHandler := &handler.AuthHandler{Repo: adminRepo}
	bookingRepo := &repository.BookingRepository{DB: db}
	closeDateRepo := &repository.CloseDateRepository{DB: db}
	closeDateHandler := &handler.CloseDateHandler{Repo: closeDateRepo}
	bookingHandler := &handler.BookingHandler{
		Repo:      bookingRepo,
		CourtRepo: courtRepo,
		CloseDateRepo: closeDateRepo,
	}

	scheduler.Start(db)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	r.Static("/uploads", "./uploads")
	router.SetupCourtRoutes(r, authHandler, courtHandler, bookingHandler, closeDateHandler)
	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}
