package booking

import (
	"context"
	"time"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

type ServiceRepository interface {
	GetByID(ctx context.Context, id int64) (domain.Service, error)
	ListActive(ctx context.Context) ([]domain.Service, error)
}

type BarberRepository interface {
	GetByID(ctx context.Context, id int64) (domain.Barber, error)
	ListBarberByService(ctx context.Context, id int64) ([]domain.Barber, error)
	GetTimeActive(ctx context.Context, id int64, selectedDate time.Time) (domain.TimeInterval, error)
}

type BookingRepository interface {
	Cancel(ctx context.Context, bookingID int64) error
	GetByID(ctx context.Context, bookingID int64) (domain.Booking, error)
	GetTimelotByBarber(
		ctx context.Context,
		barberID int64,
		selectedDate time.Time,
	) ([]domain.TimeInterval, error)
	ListByDate(
		ctx context.Context,
		startsAt time.Time,
		endsAt time.Time,
	) ([]domain.Booking, error)
	PostBooking(
		ctx context.Context,
		telegramUserID int64,
		username string,
		firstName string,
		serviceID int64,
		barberID int64,
		time domain.TimeInterval,
	) (domain.Booking, error)
	ListMyBooking(
		ctx context.Context,
		telegramUserID int64,
	) ([]domain.Booking, error)
}
