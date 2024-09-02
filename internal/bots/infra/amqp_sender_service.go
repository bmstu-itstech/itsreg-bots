package infra

import (
	"context"
	"encoding/json"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
)

const (
	amqpURI = "amqp://guest:guest@localhost:5672/"
)

type amqpSenderService struct {
	pub *amqp.Publisher
}

func NewAmqpSenderService() (interfaces.SenderService, func() error) {
	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)

	publisher, err := amqp.NewPublisher(amqpConfig, watermill.NewStdLogger(false, false))
	if err != nil {
		panic(err)
	}

	return &amqpSenderService{
		pub: publisher,
	}, publisher.Close
}

func (s *amqpSenderService) Send(ctx context.Context, botUUID string, userID int64, msg bots.Message) error {
	dto := mapMessageToDTO(botUUID, userID, msg)

	b, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	wmMsg := message.NewMessage(watermill.NewUUID(), b)
	return s.pub.Publish("messages", wmMsg)
}

type amqpMessage struct {
	BotUUID string   `json:"bot_uuid"`
	UserID  int64    `json:"user_id"`
	Text    string   `json:"text"`
	Buttons []string `json:"buttons"`
}

func mapMessageToDTO(botUUID string, userID int64, msg bots.Message) amqpMessage {
	return amqpMessage{
		BotUUID: botUUID,
		UserID:  userID,
		Text:    msg.Text,
		Buttons: msg.Buttons,
	}
}
