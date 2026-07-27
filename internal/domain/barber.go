package domain

import "time"

type Barber struct {
	ID        int64
	Name      string
	Active    bool
	CreatedAt time.Time
}
