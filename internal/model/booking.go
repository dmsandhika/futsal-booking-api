package model

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	CourtID   uuid.UUID `gorm:"type:char(36);index"`
	BookingDate time.Time `gorm:"type:date"`
	TimeSlot  string    `gorm:"type:varchar(20)"`
	CustomerName string `gorm:"type:varchar(100)"`
	CustomerPhone string `gorm:"type:varchar(20)"`
	CustomerEmail string `gorm:"type:varchar(100);null"`
	TotalPrice   int
  PaymentType PaymentType `gorm:"type:varchar(50)"`
	AmountPaid  int
	RemainingAmount int `gorm:"default:0"`
	PaymentStatus PaymentStatus `gorm:"type:varchar(50); default:'unpaid'"`
	PaymentDeadline time.Time
	Status   Status `gorm:"type:varchar(50); default:'pending'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PaymentType string

const (
    PaymentTypeDp   PaymentType = "dp"
    PaymentTypeFull PaymentType = "full"
)

type PaymentStatus string

const (
		PaymentStatusUnpaid   PaymentStatus = "unpaid"
		PaymentStatusPaid PaymentStatus = "paid"
		PaymentStatusExpired    PaymentStatus = "expired"
)

type Status string

const (
		StatusPending  Status = "pending"
		StatusConfirmed Status = "confirmed"
		StatusCancelled    Status = "cancelled"
)