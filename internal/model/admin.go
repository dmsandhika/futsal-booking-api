package model

type Admin struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex,not null,type:varchar(100)"`
	Email    string `gorm:"uniqueIndex,not null,type:varchar(100)"`
	Password string `gorm:"not null"`
}
