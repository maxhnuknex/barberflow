package booking

import (
	"testing"
	"time"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

func TestMyBookingsText(t *testing.T) {
	bookings := []domain.Booking{
		{
			ServiceName: "Мужская стрижка",
			BarberName:  "Алексей",
			StartsAt:    time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC),
			EndsAt:      time.Date(2026, 8, 1, 11, 0, 0, 0, time.UTC),
		},
	}

	got := myBookingsText(bookings)
	want := "Ваши записи:\n\n1. Услуга: Мужская стрижка\nМастер: Алексей\nВремя: 01.08.2026 10:00-11:00"

	if got != want {
		t.Fatalf("unexpected bookings text:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestMyBookingsTextEmpty(t *testing.T) {
	got := myBookingsText(nil)
	want := "У вас пока нет активных записей."

	if got != want {
		t.Fatalf("unexpected empty bookings text:\nwant: %q\ngot:  %q", want, got)
	}
}
