package handler

import (
	"futsal-booking/internal/model"
	"futsal-booking/internal/repository"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandler struct {
	Repo      *repository.BookingRepository
	CourtRepo *repository.CourtRepository
}

type BookingRequest struct {
	CourtID       string `json:"court_id"`
	BookingDate   string `json:"booking_date"`
	TimeSlot      string `json:"time_slot"`
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	TotalPrice    int    `json:"total_price"`
	PaymentType   string `json:"payment_type"`
	AmountPaid    int    `json:"amount_paid"`
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
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
		var courtName string
		if booking.Court != nil {
			courtName = booking.Court.Name
		}
		result = append(result, gin.H{
			"id":               booking.ID,
			"court_id":         booking.CourtID,
			"court_name":       courtName,
			"customer_name":    booking.CustomerName,
			"booking_date":     booking.BookingDate,
			"time_slot":        booking.TimeSlot,
			"total_price":      booking.TotalPrice,
			"payment_type":     booking.PaymentType,
			"amount_paid":      booking.AmountPaid,
			"remaining_amount": booking.RemainingAmount,
			"payment_status":   booking.PaymentStatus,
			"payment_deadline": booking.PaymentDeadline,
			"status":           booking.Status,
			"created_at":       booking.CreatedAt,
			"updated_at":       booking.UpdatedAt,
		})
	}
	totalPages := (int(total) + limit - 1) / limit
	c.JSON(200, gin.H{
		"message":    "Bookings retrieved successfully",
		"data":       result,
		"pagination": gin.H{"page": page, "limit": limit, "total": total, "total_pages": totalPages},
	})
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	courtUUID, err := uuid.Parse(req.CourtID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid court_id format"})
		return
	}
	court, err := h.CourtRepo.GetCourtByID(courtUUID)
	if err != nil || court == nil {
		c.JSON(400, gin.H{"error": "court_id not found"})
		return
	}

	booking := model.Booking{
		ID:              uuid.New(),
		CourtID:         courtUUID,
		BookingDate:     parseDate(req.BookingDate),
		TimeSlot:        req.TimeSlot,
		CustomerName:    req.CustomerName,
		CustomerPhone:   req.CustomerPhone,
		TotalPrice:      req.TotalPrice,
		PaymentType:     model.PaymentType(req.PaymentType),
		AmountPaid:      req.AmountPaid,
		RemainingAmount: req.TotalPrice - req.AmountPaid,
		PaymentStatus:   model.PaymentStatusUnpaid,
		PaymentDeadline: time.Now().Add(24 * time.Hour),
		Status:          model.StatusPending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := h.Repo.CreateBooking(&booking); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Booking created successfully", "data": booking})
}
