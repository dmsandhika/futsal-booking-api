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
	Repo          *repository.BookingRepository
	CourtRepo     *repository.CourtRepository
	CloseDateRepo *repository.CloseDateRepository
}

type BookingRequest struct {
	CourtID       string `json:"court_id"`
	BookingDate   string `json:"booking_date"`
	TimeSlot      string `json:"time_slot"`
	CustomerName  string `json:"customer_name"`
	CustomerContact string `json:"customer_contact"`
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
	timeSlot := c.Query("time_slot")
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

	bookings, total, err := h.Repo.GetBookings(&courtID, &bookingDate, timeSlot, limit, offset)
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

	bookingDate := parseDate(req.BookingDate)
	
	isClosed, err := h.CloseDateRepo.IsDateClosed(bookingDate)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if isClosed {
		c.JSON(400, gin.H{"error": "Tanggal pemesanan adalah tanggal libur, tidak bisa membuat booking"})
		return
	}
	
	existing, _, err := h.Repo.GetBookings(&courtUUID, &bookingDate, req.TimeSlot, 100, 0)
	for _, b := range existing {
		if b.TimeSlot == req.TimeSlot && b.Status != model.StatusCancelled {
			c.JSON(400, gin.H{"error": "Booking untuk lapangan, tanggal, dan jam tersebut sudah ada"})
			return
		}
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
		CustomerContact:   req.CustomerContact,
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

func (h*BookingHandler) UpdatePaymentStatus(c *gin.Context) {
	bookingId := c.Param("id")
	bookingUUID, err := uuid.Parse(bookingId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid booking ID"})
		return
	}
	booking, err := h.Repo.GetBookingByID(bookingUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if booking == nil {
		c.JSON(404, gin.H{"error": "Booking not found"})
		return
	}
	var req struct {
		PaymentStatus string `json:"payment_status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if(req.PaymentStatus == string(model.PaymentStatusPaid)){
		booking.RemainingAmount = 0
		booking.Status = model.StatusConfirmed
	}
	booking.PaymentStatus = model.PaymentStatus(req.PaymentStatus)
	booking.UpdatedAt = time.Now()
	if err := h.Repo.UpdateBooking(booking); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Booking payment status updated successfully", "data": booking})
}

func (h *BookingHandler) CancelBooking (c *gin.Context) {
	bookingId := c.Param("id")
	bookingUUID, err := uuid.Parse(bookingId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid booking ID"})
		return
	}
	booking, err := h.Repo.GetBookingByID(bookingUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if booking == nil {
		c.JSON(404, gin.H{"error": "Booking not found"})
		return
	}
	booking.Status = model.StatusCancelled
	booking.UpdatedAt = time.Now()
	if err := h.Repo.UpdateBooking(booking); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Booking cancelled successfully", "data": booking})
}


