package scheduler

import "gorm.io/gorm"

func runWithLock(db *gorm.DB, lockID int64, job func()) {
	var locked bool

	err := db.Raw(
		"SELECT pg_try_advisory_lock(?)",
		lockID,
	).Scan(&locked).Error

	if err != nil || !locked {
		return
	}

	defer db.Exec("SELECT pg_advisory_unlock(?)", lockID)

	job()
}
