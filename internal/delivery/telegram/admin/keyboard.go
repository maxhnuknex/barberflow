package admin

import (
	"fmt"
	"time"

	"github.com/go-telegram/bot/models"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

func keyboardAdminStart() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Записи",
					CallbackData: "admin:bookings",
				},
			},
			{
				{
					Text:         "Услуги",
					CallbackData: "admin:services",
				},
			},
			{
				{
					Text:         "Мастера",
					CallbackData: "admin:barbers",
				},
			},
		},
	}

}

func keyboardBookingsMenu() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Записи на сегодня",
					CallbackData: "admin:bookings:today",
				},
			},
			{
				{
					Text:         "Записи на дату",
					CallbackData: "admin:bookings:date",
				},
			},
			{
				{
					Text:         "Найти по ID",
					CallbackData: "admin:bookings:find",
				},
			},
		},
	}
}

func keyboardBookingsDate(days []time.Time) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(days)+1)

	for _, day := range days {
		button := models.InlineKeyboardButton{
			Text: day.Format("02.01.2006"),
			CallbackData: fmt.Sprintf(
				"admin:bookings:date:%s",
				day.Format("2006-01-02"),
			),
		}

		rows = append(rows, []models.InlineKeyboardButton{
			button,
		})
	}

	rows = append(rows, backButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func keyboardBookingsList(bookings []domain.Booking) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(bookings)+1)

	for _, booking := range bookings {
		button := models.InlineKeyboardButton{
			Text: fmt.Sprintf(
				"%s • %s • %s",
				booking.StartsAt.Format("15:04"),
				booking.BarberName,
				booking.ServiceName,
			),
			CallbackData: fmt.Sprintf("admin:booking:%d", booking.ID),
		}

		rows = append(rows, []models.InlineKeyboardButton{
			button,
		})
	}

	rows = append(rows, backButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func keyboardBookingDetail(bookingID int64) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Отменить запись",
					CallbackData: fmt.Sprintf("admin:booking:cancel:%d", bookingID),
				},
			},
			backButtonRow(),
		},
	}
}

func backButtonRow() []models.InlineKeyboardButton {
	return []models.InlineKeyboardButton{
		{
			Text:         "Назад",
			CallbackData: "admin:bookings",
		},
	}
}
