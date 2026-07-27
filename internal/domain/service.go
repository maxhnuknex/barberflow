package domain

import "time"

type Service struct {
	ID              int64
	Name            string
	DurationMinutes int
	PriceMinorUnits int64
	Active          bool
	CreatedAt       time.Time
}
