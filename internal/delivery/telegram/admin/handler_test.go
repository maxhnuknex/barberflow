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
	want := "📋 Запись #42\n\n🙍 Клиент: Иван (@ivan)\n✂️ Услуга: Мужская стрижка\n💈 Мастер: Алексей\n📅 Дата: 3 августа 2026\n🕒 Время: 10:00–11:00\n💳 Стоимость: 2 500 ₽"

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

	if button.Text != "🕒 10:00 • Алексей • Мужская стрижка" {
		t.Fatalf("unexpected button text: %q", button.Text)
	}
	if button.CallbackData != "admin:booking:42" {
		t.Fatalf("unexpected callback data: %q", button.CallbackData)
	}
}

func TestKeyboardServicesList(t *testing.T) {
	services := []domain.Service{
		{
			ID:              7,
			Name:            "Мужская стрижка",
			DurationMinutes: 60,
			PriceMinorUnits: 200000,
			Active:          true,
		},
	}

	keyboard := keyboardServicesList(services)
	button := keyboard.InlineKeyboard[0][0]

	if button.Text != "✅ Мужская стрижка • 60 мин • 2 000 ₽" {
		t.Fatalf("unexpected service button text: %q", button.Text)
	}
	if button.CallbackData != "admin:service:7" {
		t.Fatalf("unexpected service callback data: %q", button.CallbackData)
	}
}

func TestKeyboardBarbersList(t *testing.T) {
	barbers := []domain.Barber{
		{
			ID:     5,
			Name:   "Алексей",
			Active: true,
		},
	}

	keyboard := keyboardBarbersList(barbers)
	button := keyboard.InlineKeyboard[0][0]

	if button.Text != "✅ Алексей" {
		t.Fatalf("unexpected barber button text: %q", button.Text)
	}
	if button.CallbackData != "admin:barber:5" {
		t.Fatalf("unexpected barber callback data: %q", button.CallbackData)
	}
}

func TestServiceDetailText(t *testing.T) {
	service := domain.Service{
		ID:              7,
		Name:            "Мужская стрижка",
		DurationMinutes: 60,
		PriceMinorUnits: 250000,
		Active:          true,
	}

	got := serviceDetailText(service)
	want := "✂️ Услуга #7\n\nНазвание: Мужская стрижка\n⏱ Длительность: 60 мин\n💳 Стоимость: 2 500 ₽\nСтатус: ✅ Активна"

	if got != want {
		t.Fatalf("unexpected service detail:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestBarberDetailText(t *testing.T) {
	barber := domain.Barber{
		ID:     5,
		Name:   "Алексей",
		Active: false,
	}

	got := barberDetailText(barber)
	want := "💈 Мастер #5\n\nИмя: Алексей\nСтатус: ⛔ Отключён"

	if got != want {
		t.Fatalf("unexpected barber detail:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestStatusText(t *testing.T) {
	if serviceStatusText(true) != "✅ Активна" {
		t.Fatalf("unexpected active service status: %q", serviceStatusText(true))
	}
	if serviceStatusText(false) != "⛔ Отключена" {
		t.Fatalf("unexpected inactive service status: %q", serviceStatusText(false))
	}
	if barberStatusText(true) != "✅ Работает" {
		t.Fatalf("unexpected active barber status: %q", barberStatusText(true))
	}
	if barberStatusText(false) != "⛔ Отключён" {
		t.Fatalf("unexpected inactive barber status: %q", barberStatusText(false))
	}
}

func TestKeyboardBookingsMenuHasNoFindByID(t *testing.T) {
	keyboard := keyboardBookingsMenu()

	for _, row := range keyboard.InlineKeyboard {
		for _, button := range row {
			if button.Text == "Найти по ID" || button.CallbackData == "admin:bookings:find" {
				t.Fatalf("find by ID button should not be visible")
			}
		}
	}
}
