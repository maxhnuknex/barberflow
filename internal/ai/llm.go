package ai

import (
	"context"
)

type LLMClient interface {
	ParseBookingIntent(
		ctx context.Context,
		text string,
	) (BookingIntent, error)
}
