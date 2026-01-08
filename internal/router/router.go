package router

import (
	"futsal-booking/internal/handler"

	"futsal-booking/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCourtRoutes(r *gin.Engine,
	authHandler *handler.AuthHandler,
	courtHandler *handler.CourtHandler,
	bookingHandler *handler.BookingHandler,
	closeDateHandler *handler.CloseDateHandler,
) {
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
	bookings := r.Group("/bookings")
	{
		bookings.GET("/", bookingHandler.GetBookings)
		bookings.POST("/", bookingHandler.CreateBooking)
		bookings.PUT("/:id/payment-status", bookingHandler.UpdatePaymentStatus)
		bookings.PUT("/:id/cancel", bookingHandler.CancelBooking)
	}
	close_dates := r.Group("/close-dates")
	{
		close_dates.GET("/", closeDateHandler.GetAllCloseDates)
		close_dates.POST("/", closeDateHandler.CreateCloseDate)
		close_dates.DELETE("/", closeDateHandler.DeleteCloseDate)
	}
}
