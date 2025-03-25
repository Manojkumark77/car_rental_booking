package models

type Review struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID uint   `gorm:"not null" json:"customer_id"`
	VehicleID  uint   `gorm:"not null" json:"vehicle_id"`
	Rating     int    `gorm:"not null;check:rating BETWEEN 1 AND 5" json:"rating"`
	Comments   string `json:"comments"`

	Customer Customer `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;" json:"customer"`
	Vehicle  Vehicle  `gorm:"foreignKey:VehicleID;constraint:OnDelete:CASCADE;" json:"vehicle"`
}
