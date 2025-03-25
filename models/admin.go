package models

type Admin struct {
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Contact string `gorm:"not null" json:"contact"`
	Role    string `gorm:"not null" json:"role"`
	Address string `gorm:"not null" json:"address"`
}
