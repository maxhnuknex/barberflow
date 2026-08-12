package ai

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type stubClient struct {
	input string
}

func (c *stubClient) Complete(_ context.Context, text string) (string, error) {
	c.input = text
	return "Здравствуйте! Чем могу помочь?", nil
}

type telegramHTTPClient struct {
	messageText string
}

func (c *telegramHTTPClient) Do(request *http.Request) (*http.Response, error) {
	if err := request.ParseMultipartForm(1 << 20); err != nil {
		return nil, err
	}
	c.messageText = request.FormValue("text")

	return &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(
			`{"ok":true,"result":{"message_id":1,"date":0,"chat":{"id":42,"type":"private"}}}`,
		)),
		Header: make(http.Header),
	}, nil
}

func TestHandlerRepliesWithCompletion(t *testing.T) {
	t.Parallel()

	llm := &stubClient{}
	telegramClient := &telegramHTTPClient{}
	b, err := bot.New(
		"test-token",
		bot.WithSkipGetMe(),
		bot.WithNotAsyncHandlers(),
		bot.WithHTTPClient(time.Second, telegramClient),
	)
	if err != nil {
		t.Fatalf("create bot: %v", err)
	}

	RegisterHandler(b, NewHandler(llm, slog.New(slog.NewTextHandler(io.Discard, nil))))
	b.ProcessUpdate(context.Background(), &models.Update{
		Message: &models.Message{
			Text: "Привет, работаешь?",
			Chat: models.Chat{ID: 42},
		},
	})

	if llm.input != "Привет, работаешь?" {
		t.Fatalf("Complete received %q, want original message", llm.input)
	}
	if telegramClient.messageText != "Здравствуйте! Чем могу помочь?" {
		t.Fatalf("sent message %q, want LLM response", telegramClient.messageText)
	}
}
