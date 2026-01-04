package handler

import (
	"futsal-booking/internal/repository"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandler struct {
	Repo *repository.BookingRepository
}


func (h *BookingHandler) GetBookings(c *gin.Context) {
	courtIDParam := c.Query("court_id")
	bookingDateParam := c.Query("booking_date")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	courtID, err := uuid.Parse(courtIDParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid court ID"})
		return
	}
	bookingDate, err := time.Parse("2006-01-02", bookingDateParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid booking_date format. Use YYYY-MM-DD"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	bookings, total, err := h.Repo.GetBookings(&courtID, &bookingDate, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var result []gin.H
	for _, booking := range bookings {
		result = append(result, gin.H{
			"id": booking.ID,
			"court_id": booking.CourtID,
			"customer_name": booking.CustomerName,
			"booking_date": booking.BookingDate,
		})
	}
	totalPages := (int(total) + limit - 1) / limit
	c.JSON(200, gin.H{
		"message": "Bookings retrieved successfully",
		"data": result,
		"pagination": gin.H{"page": page, "limit": limit, "total": total, "total_pages": totalPages},
	})
}

