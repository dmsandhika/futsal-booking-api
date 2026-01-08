package model

import "time"

type CloseDate struct {
	ID        int64     `db:"id"`
	Date      time.Time `db:"date"`
	Reason    string    `db:"reason"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}