package booking

import (
	"context"
	"time"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

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

func (s *BookingService) GetServiceByID(ctx context.Context, id int64) (domain.Service, error) {
	return s.serviceRepo.GetByID(ctx, id)
}

func (s *BookingService) GetBarberByID(ctx context.Context, id int64) (domain.Barber, error) {
	return s.barberRepo.GetByID(ctx, id)
}

func (s *BookingService) Cancel(ctx context.Context, bookingID int64) error {
	return s.bookingRepository.Cancel(ctx, bookingID)
}

func (s *BookingService) GetByID(ctx context.Context, bookingID int64) (domain.Booking, error) {
	return s.bookingRepository.GetByID(ctx, bookingID)
}

func (s *BookingService) ListByDate(ctx context.Context, date time.Time) ([]domain.Booking, error) {
	loc := moscowLocation()
	selectedDate := date.In(loc)
	startsAt := time.Date(
		selectedDate.Year(),
		selectedDate.Month(),
		selectedDate.Day(),
		0,
		0,
		0,
		0,
		loc,
	)
	endsAt := startsAt.AddDate(0, 0, 1)

	return s.bookingRepository.ListByDate(ctx, startsAt, endsAt)
}

func (s *BookingService) ListBarberByService(ctx context.Context, serviceID int64) ([]domain.Barber, error) {
	return s.barberRepo.ListBarberByService(ctx, serviceID)
}

func (s *BookingService) NextSevenDays() []time.Time {
	loc := moscowLocation()
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	days := make([]time.Time, 0, 7)
	for i := 0; i < 7; i++ {
		days = append(days, today.AddDate(0, 0, i))
	}

	return days
}

func (s *BookingService) ListAvailableSlots(
	ctx context.Context,
	serviceID int64,
	barberID int64,
	selectedDate time.Time,
) ([]domain.TimeInterval, error) {
	service, err := s.serviceRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	barberTime, err := s.barberRepo.GetTimeActive(ctx, barberID, selectedDate)
	if err != nil {
		return nil, err
	}

	bookedTimes, err := s.bookingRepository.GetTimelotByBarber(ctx, barberID, selectedDate)
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

func (s *BookingService) ListTimeBarberActive(
	ctx context.Context,
	serviceID int64,
	barberID int64,
) ([]domain.TimeInterval, error) {
	return s.ListAvailableSlots(ctx, serviceID, barberID, time.Now())
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

func moscowLocation() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.FixedZone("Europe/Moscow", 3*60*60)
	}

	return loc
}
