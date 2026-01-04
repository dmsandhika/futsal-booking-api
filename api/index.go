package api

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"futsal-booking/internal/config"
	"futsal-booking/internal/handler"
	"futsal-booking/internal/model"
	"futsal-booking/internal/repository"
	"futsal-booking/internal/router"
)

var ginEngine *gin.Engine

func init() {
	godotenv.Load()
	db := config.InitDB()
	db.AutoMigrate(&model.Court{}, &model.Admin{})

	courtRepo := &repository.CourtRepository{DB: db}
	courtHandler := &handler.CourtHandler{Repo: courtRepo}
	adminRepo := &repository.AdminRepository{DB: db}
	authHandler := &handler.AuthHandler{Repo: adminRepo}

	gin.SetMode(gin.ReleaseMode)
	ginEngine = gin.New()
	ginEngine.Static("/uploads", "./uploads")
	router.SetupCourtRoutes(ginEngine, authHandler, courtHandler)
}

// Vercel entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	ginEngine.ServeHTTP(w, r)
}

// For Vercel Go API
func Main(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	Handler(w, r)
}
