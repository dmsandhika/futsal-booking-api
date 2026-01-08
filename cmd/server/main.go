package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"futsal-booking/internal/config"
	"futsal-booking/internal/handler"
	"futsal-booking/internal/model"
	"futsal-booking/internal/repository"
	"futsal-booking/internal/router"
)

func main() {
	godotenv.Load()

	db := config.InitDB()
	
	// Migrate models dengan option untuk skip table yang sudah ada
	db = db.Session(&gorm.Session{SkipDefaultTransaction: true})
	
	if err := db.AutoMigrate(&model.Court{}); err != nil {
		log.Println("Failed to migrate Court:", err)
	}
	if err := db.AutoMigrate(&model.Admin{}); err != nil {
		log.Println("Failed to migrate Admin:", err)
	}
	if err := db.AutoMigrate(&model.Booking{}); err != nil {
		log.Println("Failed to migrate Booking:", err)
	}
	if err := db.AutoMigrate(&model.CloseDate{}); err != nil {
		log.Println("Failed to migrate CloseDate:", err)
	}
	
	log.Println("All migrations completed")

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

	r := gin.Default()
	r.Static("/uploads", "./uploads")
	router.SetupCourtRoutes(r, authHandler, courtHandler, bookingHandler, closeDateHandler)
	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}
