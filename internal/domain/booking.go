package domain

import "time"

type Booking struct {
	ID int64

	CustomerID       int64
	CustomerName     string
	CustomerUsername string
	TelegramUserId   int64
	BarberID         int64
	ServiceID        int64
	BarberName       string
	ServiceName      string

	StartsAt time.Time
	EndsAt   time.Time

	PriceMinorUnits int64
	CreatedAt       time.Time
}
