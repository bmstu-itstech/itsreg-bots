package infra

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

const (
	messagesTopic = "messages"
)

type amqpMessagesPublisher struct {
	pub *amqp.Publisher
}

func NewAmqpMessagesPublisher() (bots.MessagesPublisher, <-chan *message.Message, func() error) {
	uri := os.Getenv("AMQP_URI")
	amqpConfig := amqp.NewDurableQueueConfig(uri)

	wmLogger := watermill.NewStdLogger(false, false)

	publisher, err := amqp.NewPublisher(amqpConfig, wmLogger)
	if err != nil {
		panic(err)
	}

	sub, err := amqp.NewSubscriber(amqpConfig, wmLogger)
	if err != nil {
		panic(err)
	}

	messages, err := sub.Subscribe(context.Background(), messagesTopic)
	if err != nil {
		panic(err)
	}

	return &amqpMessagesPublisher{
		pub: publisher,
	}, messages, nil
}

func (s *amqpMessagesPublisher) Publish(_ context.Context, botUUID string, userID int64, msg bots.Message) error {
	dto := mapBotMessageToAMPQ(botUUID, userID, msg)

	b, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	wmMsg := message.NewMessage(watermill.NewUUID(), b)
	return s.pub.Publish(messagesTopic, wmMsg)
}

type ampqBotMessage struct {
	BotUUID string   `json:"bot_uuid"`
	UserID  int64    `json:"user_id"`
	Text    string   `json:"text"`
	Buttons []string `json:"buttons"`
}

func mapBotMessageToAMPQ(botUUID string, userID int64, msg bots.Message) ampqBotMessage {
	return ampqBotMessage{
		BotUUID: botUUID,
		UserID:  userID,
		Text:    msg.Text,
		Buttons: msg.Buttons,
	}
}
