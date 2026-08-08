package mybookings

import (
	"testing"
	"time"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

func TestMyBookingsListText(t *testing.T) {
	bookings := []domain.Booking{
		{
			ID:          7,
			ServiceName: "Мужская стрижка",
			BarberName:  "Алексей",
			StartsAt:    time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC),
			EndsAt:      time.Date(2026, 8, 1, 11, 0, 0, 0, time.UTC),
		},
	}

	got := myBookingsListText(bookings)
	want := "📋 Мои записи\n\nВыберите запись, чтобы посмотреть подробности."

	if got != want {
		t.Fatalf("unexpected bookings text:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestMyBookingsListTextEmpty(t *testing.T) {
	got := myBookingsListText(nil)
	want := "📭 Нет активных записей\n\nВы можете создать новую запись за пару минут."

	if got != want {
		t.Fatalf("unexpected empty bookings text:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestBookingDetailText(t *testing.T) {
	booking := domain.Booking{
		ServiceName:     "Мужская стрижка",
		BarberName:      "Алексей",
		StartsAt:        time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC),
		EndsAt:          time.Date(2026, 8, 1, 11, 0, 0, 0, time.UTC),
		PriceMinorUnits: 250000,
	}

	got := bookingDetailText(booking)
	want := "📋 Ваша запись\n\n✂️ Услуга: Мужская стрижка\n💈 Мастер: Алексей\n📅 Дата: 1 августа 2026\n🕒 Время: 10:00–11:00\n💳 Стоимость: 2 500 ₽"

	if got != want {
		t.Fatalf("unexpected booking detail text:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestKeyboardListMyBookingEmpty(t *testing.T) {
	keyboard := keyboardListMyBooking(nil)
	button := keyboard.InlineKeyboard[0][0]

	if button.Text != "✂️ Записаться" {
		t.Fatalf("unexpected empty button text: %q", button.Text)
	}
	if button.CallbackData != "booking:start" {
		t.Fatalf("unexpected empty callback data: %q", button.CallbackData)
	}
}

func TestKeyboardListMyBooking(t *testing.T) {
	bookings := []domain.Booking{
		{
			ID:          7,
			ServiceName: "Мужская стрижка",
			BarberName:  "Алексей",
			StartsAt:    time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC),
		},
	}

	keyboard := keyboardListMyBooking(bookings)
	button := keyboard.InlineKeyboard[0][0]

	if button.Text != "📅 01.08 • 10:00 • Алексей" {
		t.Fatalf("unexpected button text: %q", button.Text)
	}
	if button.CallbackData != "bookings:booking:7" {
		t.Fatalf("unexpected callback data: %q", button.CallbackData)
	}
}

func TestKeyboardBookingDetailBackButton(t *testing.T) {
	keyboard := keyboardBookingDetail(7)
	button := keyboard.InlineKeyboard[1][0]

	if button.Text != "↩️ К моим записям" {
		t.Fatalf("unexpected back text: %q", button.Text)
	}
	if button.CallbackData != "bookings:list" {
		t.Fatalf("unexpected back callback data: %q", button.CallbackData)
	}
}
