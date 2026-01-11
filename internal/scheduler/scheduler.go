package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func Start(db *gorm.DB) {
	c := cron.New(cron.WithSeconds())

	// expire booking tiap 1 menit
	c.AddFunc("0 * * * * *", func() {
		runWithLock(db, 1001, func() {
			ExpireBooking(db)
		})
	})

	c.Start()
	log.Println("Scheduler started")
}
