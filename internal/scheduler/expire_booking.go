package scheduler

import (
	"log"
	"time"

	"gorm.io/gorm"
)

func ExpireBooking(db *gorm.DB) {
	now := time.Now()

	result := db.Exec(`
		UPDATE bookings
		SET status = 'expired'
		WHERE status = 'pending'
		AND payment_deadline < ?
	`, now)

	if result.Error != nil {
		log.Println("ExpireBooking error:", result.Error)
		return
	}

	log.Println("ExpireBooking success, rows:", result.RowsAffected)
}
