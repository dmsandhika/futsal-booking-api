package repository

import (
	"futsal-booking/internal/model"
	"gorm.io/gorm"
)

type AdminRepository struct {
	DB *gorm.DB
}

func (repo *AdminRepository) CreateAdmin(admin *model.Admin) error {
	return repo.DB.Create(admin).Error
}

func (repo *AdminRepository) GetAdminByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	if err := repo.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}
func (repo *AdminRepository) GetAdminByEmail(email string) (*model.Admin, error) {
	var admin model.Admin
	if err := repo.DB.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}
func (repo *AdminRepository) GetAdminByID(id uint) (*model.Admin, error) {
	var admin model.Admin
	if err := repo.DB.First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}
func (repo *AdminRepository) UpdateAdmin(admin *model.Admin) error {
	return repo.DB.Save(admin).Error
}

func (repo *AdminRepository) DeleteAdmin(id uint) error {
	return repo.DB.Delete(&model.Admin{}, id).Error
}
