package admin

import (
	"fmt"
	"time"

	"github.com/go-telegram/bot/models"
	tgui "github.com/maxhnucknex/barberflow/internal/delivery/telegram/ui"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

func keyboardAdminStart() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "📋 Записи",
					CallbackData: "admin:bookings",
				},
			},
			{
				{
					Text:         "✂️ Услуги",
					CallbackData: "admin:services",
				},
			},
			{
				{
					Text:         "💈 Мастера",
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
					Text:         "📅 Записи на сегодня",
					CallbackData: "admin:bookings:today",
				},
			},
			{
				{
					Text:         "🗓 Выбрать дату",
					CallbackData: "admin:bookings:date",
				},
			},
			adminMainMenuButtonRow(),
		},
	}
}

func keyboardBookingsDate(days []time.Time) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(days)+2)
	today := time.Now()
	if len(days) > 0 {
		today = days[0]
	}

	for _, day := range days {
		button := models.InlineKeyboardButton{
			Text: tgui.DateButton(day, today),
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
	rows = append(rows, adminMainMenuButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func keyboardBookingsList(bookings []domain.Booking) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(bookings)+2)

	for _, booking := range bookings {
		button := models.InlineKeyboardButton{
			Text: fmt.Sprintf(
				"🕒 %s • %s • %s",
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

	rows = append(rows, periodBackButtonRow())
	rows = append(rows, adminMainMenuButtonRow())

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
					CallbackData: fmt.Sprintf("admin:booking:cancel:%d", bookingID),
				},
			},
			backButtonRow(),
			adminMainMenuButtonRow(),
		},
	}
}

func keyboardServicesList(services []domain.Service) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(services)+1)

	for _, service := range services {
		button := models.InlineKeyboardButton{
			Text: fmt.Sprintf(
				"%s %s • %d мин • %s",
				serviceButtonStatus(service.Active),
				service.Name,
				service.DurationMinutes,
				tgui.Price(service.PriceMinorUnits),
			),
			CallbackData: fmt.Sprintf("admin:service:%d", service.ID),
		}

		rows = append(rows, []models.InlineKeyboardButton{
			button,
		})
	}

	rows = append(rows, adminMainMenuButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func keyboardServiceDetail(service domain.Service) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, 2)

	if service.Active {
		rows = append(rows, []models.InlineKeyboardButton{
			{
				Text:         "⛔ Отключить услугу",
				CallbackData: fmt.Sprintf("admin:service:disable:%d", service.ID),
			},
		})
	}

	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text:         "↩️ К услугам",
			CallbackData: "admin:services",
		},
	})
	rows = append(rows, adminMainMenuButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func keyboardBackToServices() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "↩️ К услугам",
					CallbackData: "admin:services",
				},
			},
			adminMainMenuButtonRow(),
		},
	}
}

func keyboardBarbersList(barbers []domain.Barber) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(barbers)+1)

	for _, barber := range barbers {
		button := models.InlineKeyboardButton{
			Text: fmt.Sprintf(
				"%s %s",
				barberButtonStatus(barber.Active),
				barber.Name,
			),
			CallbackData: fmt.Sprintf("admin:barber:%d", barber.ID),
		}

		rows = append(rows, []models.InlineKeyboardButton{
			button,
		})
	}

	rows = append(rows, adminMainMenuButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func keyboardBarberDetail(barber domain.Barber) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, 2)

	if barber.Active {
		rows = append(rows, []models.InlineKeyboardButton{
			{
				Text:         "⛔ Отключить мастера",
				CallbackData: fmt.Sprintf("admin:barber:disable:%d", barber.ID),
			},
		})
	}

	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text:         "↩️ К мастерам",
			CallbackData: "admin:barbers",
		},
	})
	rows = append(rows, adminMainMenuButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func keyboardBackToBarbers() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "↩️ К мастерам",
					CallbackData: "admin:barbers",
				},
			},
			adminMainMenuButtonRow(),
		},
	}
}

func backButtonRow() []models.InlineKeyboardButton {
	return []models.InlineKeyboardButton{
		{
			Text:         "↩️ К записям",
			CallbackData: "admin:bookings",
		},
	}
}

func periodBackButtonRow() []models.InlineKeyboardButton {
	return []models.InlineKeyboardButton{
		{
			Text:         "↩️ К выбору периода",
			CallbackData: "admin:bookings",
		},
	}
}

func adminMainMenuButtonRow() []models.InlineKeyboardButton {
	return []models.InlineKeyboardButton{
		{
			Text:         "🛠 Админ-панель",
			CallbackData: "/admin",
		},
	}
}

func keyboardBookingCanceled() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "📋 К записям",
					CallbackData: "admin:bookings",
				},
			},
			adminMainMenuButtonRow(),
		},
	}
}

func serviceButtonStatus(active bool) string {
	if active {
		return "✅"
	}

	return "⛔"
}

func barberButtonStatus(active bool) string {
	if active {
		return "✅"
	}

	return "⛔"
}
