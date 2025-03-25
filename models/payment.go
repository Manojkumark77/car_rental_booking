package models

type Payment struct {
	ID            uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID     uint    `gorm:"not null" json:"booking_id"`
	Amount        float64 `gorm:"not null" json:"amount"`
	PaymentMethod string  `gorm:"type:enum('Credit Card', 'Debit Card', 'Cash');not null" json:"payment_method"`
	Status        string  `gorm:"type:enum('Pending', 'Completed', 'Failed');default:'Pending'" json:"status"`

	Booking Booking `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"booking"`
}
