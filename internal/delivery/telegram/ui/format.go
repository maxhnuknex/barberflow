package ui

import (
	"strconv"
	"time"
)

var months = [...]string{
	"",
	"января",
	"февраля",
	"марта",
	"апреля",
	"мая",
	"июня",
	"июля",
	"августа",
	"сентября",
	"октября",
	"ноября",
	"декабря",
}

var weekdays = map[time.Weekday]string{
	time.Monday:    "Пн",
	time.Tuesday:   "Вт",
	time.Wednesday: "Ср",
	time.Thursday:  "Чт",
	time.Friday:    "Пт",
	time.Saturday:  "Сб",
	time.Sunday:    "Вс",
}

func Price(minorUnits int64) string {
	rubles := minorUnits / 100
	value := strconv.FormatInt(rubles, 10)
	if len(value) <= 3 {
		return value + " ₽"
	}

	result := make([]byte, 0, len(value)+len(value)/3)
	firstGroup := len(value) % 3
	if firstGroup == 0 {
		firstGroup = 3
	}

	result = append(result, value[:firstGroup]...)
	for i := firstGroup; i < len(value); i += 3 {
		result = append(result, ' ')
		result = append(result, value[i:i+3]...)
	}

	return string(result) + " ₽"
}

func FullDate(date time.Time) string {
	month := months[int(date.Month())]
	return strconv.Itoa(date.Day()) + " " + month + " " + strconv.Itoa(date.Year())
}

func DayMonth(date time.Time) string {
	month := months[int(date.Month())]
	return strconv.Itoa(date.Day()) + " " + month
}

func ShortDate(date time.Time) string {
	return date.Format("02.01")
}

func TimeInterval(startsAt time.Time, endsAt time.Time) string {
	return startsAt.Format("15:04") + "–" + endsAt.Format("15:04")
}

func DateButton(date time.Time, today time.Time) string {
	current := dayStart(date)
	base := dayStart(today)

	switch {
	case current.Equal(base):
		return "Сегодня • " + ShortDate(date)
	case current.Equal(base.AddDate(0, 0, 1)):
		return "Завтра • " + ShortDate(date)
	default:
		return weekdays[date.Weekday()] + " • " + ShortDate(date)
	}
}

func dayStart(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}
