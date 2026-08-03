package admin

import (
	"testing"
	"time"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

func TestBookingDetailText(t *testing.T) {
	booking := domain.Booking{
		ID:               42,
		CustomerName:     "Иван",
		CustomerUsername: "ivan",
		BarberName:       "Алексей",
		ServiceName:      "Мужская стрижка",
		StartsAt:         time.Date(2026, 8, 3, 10, 0, 0, 0, time.UTC),
		EndsAt:           time.Date(2026, 8, 3, 11, 0, 0, 0, time.UTC),
		PriceMinorUnits:  250000,
	}

	got := bookingDetailText(booking)
	want := "ID: 42\nКлиент: Иван (@ivan)\nМастер: Алексей\nУслуга: Мужская стрижка\nДата: 03.08.2026\nВремя: 10:00-11:00\nЦена: 2500 ₽"

	if got != want {
		t.Fatalf("unexpected booking detail text:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestKeyboardBookingsList(t *testing.T) {
	bookings := []domain.Booking{
		{
			ID:          42,
			BarberName:  "Алексей",
			ServiceName: "Мужская стрижка",
			StartsAt:    time.Date(2026, 8, 3, 10, 0, 0, 0, time.UTC),
		},
	}

	keyboard := keyboardBookingsList(bookings)
	button := keyboard.InlineKeyboard[0][0]

	if button.Text != "10:00 • Алексей • Мужская стрижка" {
		t.Fatalf("unexpected button text: %q", button.Text)
	}
	if button.CallbackData != "admin:booking:42" {
		t.Fatalf("unexpected callback data: %q", button.CallbackData)
	}
}
