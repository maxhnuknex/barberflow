package booking

import (
	"testing"
	"time"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

func TestBookingPreviewText(t *testing.T) {
	service := domain.Service{
		Name:            "Мужская стрижка",
		PriceMinorUnits: 250000,
	}
	barber := domain.Barber{
		Name: "Алексей",
	}
	slot := domain.TimeInterval{
		StartsAt: time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC),
		EndsAt:   time.Date(2026, 8, 4, 16, 0, 0, 0, time.UTC),
	}

	got := bookingPreviewText(service, barber, slot)
	want := "✅ Проверьте запись\n\n✂️ Услуга: Мужская стрижка\n💈 Мастер: Алексей\n📅 Дата: 4 августа 2026\n🕒 Время: 15:00–16:00\n💳 Стоимость: 2 500 ₽\n\nВсё верно?"

	if got != want {
		t.Fatalf("unexpected preview text:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestCreatedBookingText(t *testing.T) {
	booking := domain.Booking{
		ServiceName:     "Мужская стрижка",
		BarberName:      "Алексей",
		StartsAt:        time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC),
		EndsAt:          time.Date(2026, 8, 4, 16, 0, 0, 0, time.UTC),
		PriceMinorUnits: 250000,
	}

	got := createdBookingText(booking)
	want := "🎉 Запись подтверждена\n\n✂️ Услуга: Мужская стрижка\n💈 Мастер: Алексей\n📅 Дата: 4 августа 2026\n🕒 Время: 15:00–16:00\n💳 Стоимость: 2 500 ₽\n\nЗапись сохранена в разделе «Мои записи»."

	if got != want {
		t.Fatalf("unexpected created text:\nwant: %q\ngot:  %q", want, got)
	}
}
