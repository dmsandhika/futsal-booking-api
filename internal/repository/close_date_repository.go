package repository

import (
	"time"

	"futsal-booking/internal/model"

	"gorm.io/gorm"
)

type CloseDateRepository struct {
	DB *gorm.DB
}

func (r *CloseDateRepository) CreateCloseDate(date time.Time, reason string) (*model.CloseDate, error) {
	closeDate := &model.CloseDate{
		Date:      date,
		Reason:    reason,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := r.DB.Create(closeDate).Error; err != nil {
		return nil, err
	}
	return closeDate, nil
}

func (r *CloseDateRepository) GetAllCloseDates() ([]model.CloseDate, error) {
	var closeDates []model.CloseDate
	if err := r.DB.Find(&closeDates).Error; err != nil {
		return nil, err
	}
	return closeDates, nil
}

func (r *CloseDateRepository) DeleteCloseDate(id int64) error {
	if err := r.DB.Delete(&model.CloseDate{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *CloseDateRepository) IsDateClosed(date time.Time) (bool, error) {
	var count int64
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	if err := r.DB.Model(&model.CloseDate{}).Where("DATE(date) = ?", dateOnly.Format("2006-01-02")).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}