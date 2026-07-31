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
	GetTimeActive(ctx context.Context, id int64) (domain.TimeInterval, error)
}

type BookingRepository interface {
	GetTimelotByBarber(
		ctx context.Context,
		barberID int64,
	) ([]domain.TimeInterval, error)
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

type BookingService struct {
	serviceRepo       ServiceRepository
	barberRepo        BarberRepository
	bookingRepository BookingRepository
}

func NewService(
	serviceRepo ServiceRepository,
	barberRepo BarberRepository,
	bookingRepository BookingRepository,
) *BookingService {
	return &BookingService{
		serviceRepo:       serviceRepo,
		barberRepo:        barberRepo,
		bookingRepository: bookingRepository,
	}
}

func (s *BookingService) ListActive(ctx context.Context) ([]domain.Service, error) {
	return s.serviceRepo.ListActive(ctx)
}

func (s *BookingService) ListBarberByService(ctx context.Context, serviceID int64) ([]domain.Barber, error) {
	return s.barberRepo.ListBarberByService(ctx, serviceID)
}

func (s *BookingService) ListTimeBarberActive(
	ctx context.Context,
	serviceID int64,
	barberID int64,
) ([]domain.TimeInterval, error) {
	service, err := s.serviceRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	barberTime, err := s.barberRepo.GetTimeActive(ctx, barberID)
	if err != nil {
		return nil, err
	}

	bookedTimes, err := s.bookingRepository.GetTimelotByBarber(ctx, barberID)
	if err != nil {
		return nil, err
	}

	freeTimes := calculateFreeSlots(
		service.DurationMinutes,
		barberTime,
		bookedTimes,
	)

	return freeTimes, nil
}

func (s *BookingService) PostBooking(
	ctx context.Context,
	telegramUserID int64,
	username string,
	firstName string,
	serviceID int64,
	barberID int64,
	time domain.TimeInterval,
) (domain.Booking, error) {
	return s.bookingRepository.PostBooking(
		ctx,
		telegramUserID,
		username,
		firstName,
		serviceID,
		barberID,
		time,
	)
}

func calculateFreeSlots(
	duration int,
	barberTime domain.TimeInterval,
	bookedTimes []domain.TimeInterval,
) []domain.TimeInterval {
	slotDuration := time.Duration(duration) * time.Minute

	freeSlots := make([]domain.TimeInterval, 0)

	for slotStart := barberTime.StartsAt; !slotStart.Add(slotDuration).After(barberTime.EndsAt); slotStart = slotStart.Add(time.Hour) {

		slotEnd := slotStart.Add(slotDuration)

		isBusy := false

		for _, booked := range bookedTimes {
			hasOverlap :=
				slotStart.Before(booked.EndsAt) &&
					slotEnd.After(booked.StartsAt)

			if hasOverlap {
				isBusy = true
				break
			}
		}

		if !isBusy {
			freeSlots = append(freeSlots, domain.TimeInterval{
				StartsAt: slotStart,
				EndsAt:   slotEnd,
			})
		}
	}

	return freeSlots
}

func (s *BookingService) ListMyBooking(
	ctx context.Context,
	telegramUserID int64,
) ([]domain.Booking, error) {
	return s.bookingRepository.ListMyBooking(ctx, telegramUserID)
}
