package router

import (
	"futsal-booking/internal/handler"

	"github.com/gin-gonic/gin"
	"futsal-booking/internal/middleware"
)

func SetupCourtRoutes(r *gin.Engine, 
	authHandler *handler.AuthHandler,
	courtHandler *handler.CourtHandler,){
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.GET("/me", middleware.JWTAuth(), authHandler.GetMe)
	}
	courts := r.Group("/courts")
	{
		courts.GET("/", courtHandler.GetAllCourts)
		courts.POST("/", courtHandler.CreateCourt)
		courts.GET("/:id", courtHandler.GetCourtByID)
		courts.PUT("/:id", courtHandler.UpdateCourt)
		courts.DELETE("/:id", courtHandler.DeleteCourt)
	}
}