package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (repo *BookingRepository) Cancel(ctx context.Context, bookingID int64) error {
	const query = `
		DELETE FROM bookings
		WHERE id = $1
	`

	if _, err := repo.db.Exec(ctx, query, bookingID); err != nil {
		return fmt.Errorf("cancel booking: %w", err)
	}

	return nil
}

func (repo *BookingRepository) GetByID(
	ctx context.Context,
	bookingID int64,
) (domain.Booking, error) {
	const query = `
		SELECT
			b.id,
			b.customer_id,
			c.first_name,
			COALESCE(c.username, ''),
			b.service_id,
			s.name,
			b.barber_id,
			br.name,
			b.starts_at,
			b.ends_at,
			s.price_minor_units,
			b.created_at
		FROM bookings b
		JOIN customers c
			ON c.id = b.customer_id
		JOIN services s
			ON s.id = b.service_id
		JOIN barbers br
			ON br.id = b.barber_id
		WHERE b.id = $1
	`

	var booking domain.Booking
	if err := repo.db.QueryRow(ctx, query, bookingID).Scan(
		&booking.ID,
		&booking.CustomerID,
		&booking.CustomerName,
		&booking.CustomerUsername,
		&booking.ServiceID,
		&booking.ServiceName,
		&booking.BarberID,
		&booking.BarberName,
		&booking.StartsAt,
		&booking.EndsAt,
		&booking.PriceMinorUnits,
		&booking.CreatedAt,
	); err != nil {
		return domain.Booking{}, fmt.Errorf("get booking by id: %w", err)
	}

	return booking, nil
}

func (repo *BookingRepository) GetBookingsForReminder(ctx context.Context) ([]domain.Booking, error) {
	const query = `
		SELECT
			b.id,
			b.customer_id,
			c.first_name,
			COALESCE(c.username, ''),
			c.telegram_user_id,
			b.service_id,
			s.name,
			b.barber_id,
			br.name,
			b.starts_at,
			b.ends_at,
			s.price_minor_units,
			b.created_at
		FROM bookings b
		JOIN customers c
			ON c.id = b.customer_id
		JOIN services s
			ON s.id = b.service_id
		JOIN barbers br
			ON br.id = b.barber_id
		WHERE b.starts_at > NOW()
			AND b.starts_at <= NOW() + INTERVAL '2 hours'
			AND b.reminder_sent_at IS NULL
		ORDER BY b.starts_at
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list bookings for reminder: %w", err)
	}
	defer rows.Close()

	bookings := make([]domain.Booking, 0)

	for rows.Next() {
		var booking domain.Booking

		if err := rows.Scan(
			&booking.ID,
			&booking.CustomerID,
			&booking.CustomerName,
			&booking.CustomerUsername,
			&booking.TelegramUserId,
			&booking.ServiceID,
			&booking.ServiceName,
			&booking.BarberID,
			&booking.BarberName,
			&booking.StartsAt,
			&booking.EndsAt,
			&booking.PriceMinorUnits,
			&booking.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan booking for reminder: %w", err)
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate bookings for reminder: %w", err)
	}

	return bookings, nil
}

func (repo *BookingRepository) GetTimelotByBarber(
	ctx context.Context,
	barberID int64,
	selectedDate time.Time,
) ([]domain.TimeInterval, error) {
	const query = `
		SELECT starts_at, ends_at
		FROM bookings
		WHERE barber_id = $1
			AND starts_at >= $2::date
			AND starts_at < $2::date + INTERVAL '1 day'
	`
	rows, err := repo.db.Query(ctx, query, barberID, selectedDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	times := make([]domain.TimeInterval, 0)
	for rows.Next() {
		var time domain.TimeInterval

		if err := rows.Scan(
			&time.StartsAt,
			&time.EndsAt,
		); err != nil {
			return nil, err
		}

		times = append(times, time)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return times, nil
}

func (repo *BookingRepository) ListByDate(
	ctx context.Context,
	startsAt time.Time,
	endsAt time.Time,
) ([]domain.Booking, error) {
	const query = `
		SELECT
			b.id,
			b.customer_id,
			c.first_name,
			COALESCE(c.username, ''),
			b.service_id,
			s.name,
			b.barber_id,
			br.name,
			b.starts_at,
			b.ends_at,
			s.price_minor_units,
			b.created_at
		FROM bookings b
		JOIN customers c
			ON c.id = b.customer_id
		JOIN services s
			ON s.id = b.service_id
		JOIN barbers br
			ON br.id = b.barber_id
		WHERE b.starts_at >= $1
			AND b.starts_at < $2
		ORDER BY b.starts_at
	`

	rows, err := repo.db.Query(ctx, query, startsAt, endsAt)
	if err != nil {
		return nil, fmt.Errorf("list bookings by date: %w", err)
	}
	defer rows.Close()

	bookings := make([]domain.Booking, 0)

	for rows.Next() {
		var booking domain.Booking

		if err := rows.Scan(
			&booking.ID,
			&booking.CustomerID,
			&booking.CustomerName,
			&booking.CustomerUsername,
			&booking.ServiceID,
			&booking.ServiceName,
			&booking.BarberID,
			&booking.BarberName,
			&booking.StartsAt,
			&booking.EndsAt,
			&booking.PriceMinorUnits,
			&booking.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan booking by date: %w", err)
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate bookings by date: %w", err)
	}

	return bookings, nil
}

func (repo *BookingRepository) PostBooking(
	ctx context.Context,
	telegramUserID int64,
	username string,
	firstName string,
	serviceID int64,
	barberID int64,
	time domain.TimeInterval,
) (domain.Booking, error) {
	const query = `
		WITH customer AS (
			INSERT INTO customers (telegram_user_id, username, first_name)
			VALUES ($1, $2, $3)
			ON CONFLICT (telegram_user_id) DO UPDATE
				SET username = EXCLUDED.username,
					first_name = EXCLUDED.first_name
			RETURNING id
		),
		created_booking AS (
			INSERT INTO bookings (customer_id, service_id, barber_id, starts_at, ends_at)
			SELECT customer.id, $4, $5, $6, $7
			FROM customer
			RETURNING id, customer_id, service_id, barber_id, starts_at, ends_at, created_at
		)
		SELECT
			cb.id,
			cb.customer_id,
			cb.service_id,
			s.name,
			cb.barber_id,
			b.name,
			cb.starts_at,
			cb.ends_at,
			s.price_minor_units,
			cb.created_at
		FROM created_booking cb
		JOIN services s ON s.id = cb.service_id
		JOIN barbers b ON b.id = cb.barber_id
	`

	var booking domain.Booking

	if err := repo.db.QueryRow(
		ctx,
		query,
		telegramUserID,
		username,
		firstName,
		serviceID,
		barberID,
		time.StartsAt,
		time.EndsAt,
	).Scan(
		&booking.ID,
		&booking.CustomerID,
		&booking.ServiceID,
		&booking.ServiceName,
		&booking.BarberID,
		&booking.BarberName,
		&booking.StartsAt,
		&booking.EndsAt,
		&booking.PriceMinorUnits,
		&booking.CreatedAt,
	); err != nil {
		return domain.Booking{}, err
	}

	return booking, nil
}

func (repo *BookingRepository) ListMyBooking(
	ctx context.Context,
	telegramUserID int64,
) ([]domain.Booking, error) {
	const query = `
		SELECT
			b.id,
			b.customer_id,
			c.first_name,
			COALESCE(c.username, ''),
			b.service_id,
			s.name,
			b.barber_id,
			br.name,
			b.starts_at,
			b.ends_at,
			s.price_minor_units,
			b.created_at
		FROM bookings b
		JOIN customers c
			ON c.id = b.customer_id
		JOIN services s
			ON s.id = b.service_id
		JOIN barbers br
			ON br.id = b.barber_id
		WHERE c.telegram_user_id = $1
		  AND b.starts_at >= NOW()
		ORDER BY b.starts_at
	`

	rows, err := repo.db.Query(ctx, query, telegramUserID)
	if err != nil {
		return nil, fmt.Errorf("list customer bookings: %w", err)
	}
	defer rows.Close()

	bookings := make([]domain.Booking, 0)

	for rows.Next() {
		var booking domain.Booking

		if err := rows.Scan(
			&booking.ID,
			&booking.CustomerID,
			&booking.CustomerName,
			&booking.CustomerUsername,
			&booking.ServiceID,
			&booking.ServiceName,
			&booking.BarberID,
			&booking.BarberName,
			&booking.StartsAt,
			&booking.EndsAt,
			&booking.PriceMinorUnits,
			&booking.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan customer booking: %w", err)
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate customer bookings: %w", err)
	}

	return bookings, nil
}

func (repo *BookingRepository) MarkReminderSent(ctx context.Context, bookingID int64) error {
	const query = `
		UPDATE bookings
		SET reminder_sent_at = NOW()
		WHERE id = $1
	`

	if _, err := repo.db.Exec(ctx, query, bookingID); err != nil {
		return fmt.Errorf("mark reminder sent: %w", err)
	}

	return nil
}
