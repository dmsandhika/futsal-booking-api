package repository

import (
	"futsal-booking/internal/model"	
	"gorm.io/gorm"
)

type CourtRepository struct {
	DB *gorm.DB
}

func (r *CourtRepository) GetAllCourts() ([]model.Court, error) {
	var courts []model.Court
	result := r.DB.Find(&courts)
	return courts, result.Error
}

func (r *CourtRepository) GetAllCourtsPaginated(page, limit int) ([]model.Court, int64, error) {
	var courts []model.Court
	var total int64

	offset := (page - 1) * limit
	result := r.DB.Model(&model.Court{}).Count(&total).Offset(offset).Limit(limit).Find(&courts)
	return courts, total, result.Error
}

func (r *CourtRepository) CreateCourt(court *model.Court) error {
	result := r.DB.Create(court)
	return result.Error
}

func (r *CourtRepository) GetCourtByID(id uint) (*model.Court, error) {
	var court model.Court
	result := r.DB.First(&court, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &court, nil
}

func (r *CourtRepository) UpdateCourt(court *model.Court) error {
	result := r.DB.Model(&model.Court{}).Where("id = ?", court.ID).Updates(court)
	return result.Error
}

func (r *CourtRepository) DeleteCourt(id uint) error {
	result := r.DB.Delete(&model.Court{}, id)
	return result.Error
}