package models

type Vehicle struct {
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Model        string  `gorm:"not null" json:"model"`
	Year         int     `gorm:"not null" json:"year"`
	RentalRate   float64 `gorm:"type:decimal(10,2);not null" json:"rental_rate"`
	Availability string  `gorm:"type:enum('Available','Rented','Maintenance');default:'Available';not null" json:"availability"`
	Type         string  `gorm:"type:enum('Suzuki','Toyota','Honda','Tata');not null" json:"type"`
	Mileage      int     `gorm:"not null" json:"mileage"`
}
