package ai

type BookingIntent struct {
	Intent   string  `json:"intent"`
	Service  *string `json:"service"`
	Barber   *string `json:"barber"`
	Date     *string `json:"date"`
	TimeFrom *string `json:"time_from"`
}
