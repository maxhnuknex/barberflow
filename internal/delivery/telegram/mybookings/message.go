package mybookings

import (
	"fmt"
	"strings"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

func myBookingsText(bookings []domain.Booking) string {
	if len(bookings) == 0 {
		return "У вас пока нет активных записей."
	}

	var builder strings.Builder
	builder.WriteString("Ваши записи:")

	for i, booking := range bookings {
		fmt.Fprintf(
			&builder,
			"\n\n%d. Услуга: %s\nМастер: %s\nВремя: %s-%s",
			i+1,
			booking.ServiceName,
			booking.BarberName,
			booking.StartsAt.Format("02.01.2006 15:04"),
			booking.EndsAt.Format("15:04"),
		)
	}

	return builder.String()
}
