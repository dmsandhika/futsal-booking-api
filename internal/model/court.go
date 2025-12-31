package model

import (
	"time"
)

type Court struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string
	Location      string `gorm:"type:varchar(255)"`
	PricePerHour  int
	Status        string
	TimeSlots     string `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}