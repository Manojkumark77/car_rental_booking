package models

type Customer struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string `gorm:"not null" json:"name"`
	Contact       string `gorm:"not null" json:"contact"`
	Address       string `gorm:"not null" json:"address"`
	LicenseNumber string `gorm:"unique;not null" json:"license_number"`
}
