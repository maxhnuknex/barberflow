package ai

import (
	"context"
)

type Orchestrator struct {
	llm LLMClient
}

func NewOrchestrator(llm LLMClient) *Orchestrator {
	return &Orchestrator{
		llm: llm,
	}
}

func (o *Orchestrator) HandleMessage(
	ctx context.Context,
	text string,
) (BookingIntent, error) {
	return o.llm.ParseBookingIntent(ctx, text)
}
