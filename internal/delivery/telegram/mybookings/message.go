package mybookings

import (
	"fmt"
	"strings"

	tgui "github.com/maxhnucknex/barberflow/internal/delivery/telegram/ui"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

func myBookingsListText(bookings []domain.Booking) string {
	if len(bookings) == 0 {
		return "📭 Нет активных записей\n\nВы можете создать новую запись за пару минут."
	}

	return "📋 Мои записи\n\nВыберите запись, чтобы посмотреть подробности."
}

func bookingDetailText(booking domain.Booking) string {
	var builder strings.Builder
	builder.WriteString("📋 Ваша запись")

	fmt.Fprintf(
		&builder,
		"\n\n✂️ Услуга: %s\n💈 Мастер: %s\n📅 Дата: %s\n🕒 Время: %s\n💳 Стоимость: %s",
		booking.ServiceName,
		booking.BarberName,
		tgui.FullDate(booking.StartsAt),
		tgui.TimeInterval(booking.StartsAt, booking.EndsAt),
		tgui.Price(booking.PriceMinorUnits),
	)
	return builder.String()
}
