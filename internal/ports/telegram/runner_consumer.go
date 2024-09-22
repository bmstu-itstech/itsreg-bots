package telegram

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
)

type runnerMessagesHandler func(ctx context.Context, msg runnerMessage) error

type runnerConsumer struct {
	ch  <-chan *message.Message
	h   runnerMessagesHandler
	log *slog.Logger
}

func newRunnerConsumer(
	ch <-chan *message.Message,
	h runnerMessagesHandler,
	log *slog.Logger,
) *runnerConsumer {
	return &runnerConsumer{
		ch:  ch,
		h:   h,
		log: log,
	}
}

func (c *runnerConsumer) Process() {
	for msg := range c.ch {
		runnerMsg, err := unmarshalRunnerMessage(msg)
		if err != nil {
			c.log.Error("failed to unmarshal bot message", "error", err.Error())
			msg.Nack()
			continue
		}

		err = c.h(msg.Context(), runnerMsg)
		if err != nil {
			c.log.Error("failed to handle bot message", "error", err.Error())
			msg.Nack()
			continue
		}

		msg.Ack()
	}
}

type runnerMessage struct {
	BotUUID string `json:"bot_uuid"`
	Command string `json:"command"`
}

func unmarshalRunnerMessage(msg *message.Message) (runnerMessage, error) {
	var res runnerMessage
	err := json.Unmarshal(msg.Payload, &res)
	return res, err
}
