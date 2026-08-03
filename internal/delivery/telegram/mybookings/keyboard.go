package mybookings

import (
	"fmt"

	"github.com/go-telegram/bot/models"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

func keyboardListMyBooking(
	bookings []domain.Booking,
) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(bookings)+1)

	for _, booking := range bookings {
		button := models.InlineKeyboardButton{
			Text: fmt.Sprintf(
				"%s %s",
				booking.ServiceName,
				booking.StartsAt.Format("02.01 15:04"),
			),
			CallbackData: "bookings:list",
		}

		rows = append(rows, []models.InlineKeyboardButton{
			button,
		})
	}

	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text:         "Главное меню",
			CallbackData: "/start",
		},
	})

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}
