package router

import (
	"futsal-booking/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupCourtRoutes(r *gin.Engine, 
	courtHandler *handler.CourtHandler){
	courts := r.Group("/courts")
	{
		courts.GET("/", courtHandler.GetAllCourts)
		courts.POST("/", courtHandler.CreateCourt)
		courts.GET("/:id", courtHandler.GetCourtByID)
		courts.PUT("/:id", courtHandler.UpdateCourt)
		courts.DELETE("/:id", courtHandler.DeleteCourt)
	}
}