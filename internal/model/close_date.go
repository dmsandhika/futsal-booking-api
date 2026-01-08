package model

import "time"

type CloseDate struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Date      time.Time `gorm:"column:date" json:"date"`
	Reason    string    `gorm:"column:reason" json:"reason"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}