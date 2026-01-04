package repository

import (
	"futsal-booking/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRepository struct {
	DB *gorm.DB
}

func (r *BookingRepository) GetBookings(courtID *uuid.UUID, bookingDate *time.Time, limit, offset int) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var total int64
	db := r.DB.Model(&model.Booking{}).Preload("Court")
	if courtID != nil {
		db = db.Where("court_id = ?", *courtID)
	}
	if bookingDate != nil {
		db = db.Where("booking_date = ?", *bookingDate)
	}
	db.Count(&total)
	result := db.Limit(limit).Offset(offset).Find(&bookings)
	return bookings, total, result.Error
}

func (r *BookingRepository) CreateBooking(booking *model.Booking) error {
	result := r.DB.Create(booking)
	return result.Error
}

func (r *BookingRepository) GetBookingsByCourtID(courtID uuid.UUID) ([]model.Booking, error) {
	var bookings []model.Booking
	result := r.DB.Where("court_id = ?", courtID).Find(&bookings)
	return bookings, result.Error
}

func (r *BookingRepository) GetBookingByID(id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	result := r.DB.First(&booking, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &booking, nil
}

func (r *BookingRepository) GetBookingsByUserName(name string) ([]model.Booking, error) {
	var bookings []model.Booking
	result := r.DB.Where("customer_name = ?", name).Find(&bookings)
	return bookings, result.Error
}

func (r *BookingRepository) UpdateBooking(booking *model.Booking) error {
	result := r.DB.Model(&model.Booking{}).Where("id = ?", booking.ID).Updates(booking)
	return result.Error
}

func (r *BookingRepository) DeleteBooking(id uuid.UUID) error {
	result := r.DB.Delete(&model.Booking{}, "id = ?", id)
	return result.Error
}
