package mybookings

import (
	"fmt"

	"github.com/go-telegram/bot/models"
	tgui "github.com/maxhnucknex/barberflow/internal/delivery/telegram/ui"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

func keyboardListMyBooking(
	bookings []domain.Booking,
) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(bookings)+1)

	if len(bookings) == 0 {
		rows = append(rows, []models.InlineKeyboardButton{
			{
				Text:         "✂️ Записаться",
				CallbackData: "booking:start",
			},
		})
		rows = append(rows, userMainMenuButtonRow())

		return &models.InlineKeyboardMarkup{
			InlineKeyboard: rows,
		}
	}

	for _, booking := range bookings {
		button := models.InlineKeyboardButton{
			Text: fmt.Sprintf(
				"📅 %s • %s • %s",
				tgui.ShortDate(booking.StartsAt),
				booking.StartsAt.Format("15:04"),
				booking.BarberName,
			),
			CallbackData: fmt.Sprintf("bookings:booking:%d", booking.ID),
		}

		rows = append(rows, []models.InlineKeyboardButton{
			button,
		})
	}

	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text:         "✂️ Новая запись",
			CallbackData: "booking:start",
		},
	})
	rows = append(rows, userMainMenuButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func keyboardBookingDetail(bookingID int64) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "❌ Отменить запись",
					CallbackData: fmt.Sprintf("bookings:cancel:%d", bookingID),
				},
			},
			{
				{
					Text:         "↩️ К моим записям",
					CallbackData: "bookings:list",
				},
			},
			userMainMenuButtonRow(),
		},
	}
}

func cancelBookingKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "✂️ Записаться снова",
					CallbackData: "booking:start",
				},
			},
			{
				{
					Text:         "📋 Мои записи",
					CallbackData: "bookings:list",
				},
			},
			userMainMenuButtonRow(),
		},
	}
}

func mainMenuKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			userMainMenuButtonRow(),
		},
	}
}

func userMainMenuButtonRow() []models.InlineKeyboardButton {
	return []models.InlineKeyboardButton{
		{
			Text:         "🏠 Главное меню",
			CallbackData: "/start",
		},
	}
}
