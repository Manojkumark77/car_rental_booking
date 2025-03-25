package models

import (
	"time"
)

type Booking struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID uint      `gorm:"not null;index" json:"customer_id"`
	VehicleID  uint      `gorm:"not null;index" json:"vehicle_id"`
	StartDate  time.Time `gorm:"not null" json:"start_date"`
	EndDate    time.Time `gorm:"not null" json:"end_date"`
	Status     string    `gorm:"type:enum('Pending', 'Confirmed', 'Cancelled', 'Completed');default:'Pending'" json:"status"`

	Customer Customer `gorm:"foreignKey:CustomerID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"customer"`
	Vehicle  Vehicle  `gorm:"foreignKey:VehicleID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"vehicle"`
}
