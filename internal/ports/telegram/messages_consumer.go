package telegram

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
)

type botMessageHandler func(ctx context.Context, msg botMessage) error

type messagesConsumer struct {
	ch  <-chan *message.Message
	h   botMessageHandler
	log *slog.Logger
}

func newMessagesConsumer(
	ch <-chan *message.Message,
	h botMessageHandler,
	log *slog.Logger,
) *messagesConsumer {
	return &messagesConsumer{
		ch:  ch,
		h:   h,
		log: log,
	}
}

func (c *messagesConsumer) Process() {
	for msg := range c.ch {
		botMsg, err := unmarshalBotMessage(msg)
		if err != nil {
			c.log.Error("failed to unmarshal bot message", "error", err.Error())
			msg.Nack()
			continue
		}

		err = c.h(msg.Context(), botMsg)
		if err != nil {
			c.log.Error("failed to handle bot message", "error", err.Error())
			msg.Nack()
			continue
		}
		msg.Ack()
	}
}

type botMessage struct {
	BotUUID string   `json:"bot_uuid"`
	UserID  int64    `json:"user_id"`
	Text    string   `json:"text"`
	Buttons []string `json:"buttons"`
}

func unmarshalBotMessage(msg *message.Message) (botMessage, error) {
	var res botMessage
	err := json.Unmarshal(msg.Payload, &res)
	return res, err
}
