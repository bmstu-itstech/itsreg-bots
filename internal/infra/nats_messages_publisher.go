package infra

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/logs/sl"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	nc "github.com/nats-io/nats.go"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/logs"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type natsMessagesPublisher struct {
	pub *nats.Publisher
}

func NewNATSMessagesPublisher() (bots.MessagesPublisher, <-chan *message.Message, func() error) {
	marshaller := &nats.GobMarshaler{}
	logger := sl.NewWatermillLoggerAdapter(logs.DefaultLogger())
	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(10 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}

	jsConfig := nats.JetStreamConfig{Disabled: true}

	uri := os.Getenv("NATS_URI")
	if uri == "" {
		panic("NATS_URI environment variable not set")
	}

	sub, err := nats.NewSubscriber(
		nats.SubscriberConfig{
			URL:            uri,
			CloseTimeout:   10 * time.Second,
			AckWaitTimeout: 10 * time.Second,
			NatsOptions:    options,
			Unmarshaler:    marshaller,
			JetStream:      jsConfig,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	messages, err := sub.Subscribe(context.Background(), messagesTopic)
	if err != nil {
		panic(err)
	}

	pub, err := nats.NewPublisher(
		nats.PublisherConfig{
			URL:         uri,
			NatsOptions: options,
			Marshaler:   marshaller,
			JetStream:   jsConfig,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	return natsMessagesPublisher{
			pub: pub,
		}, messages, func() error {
			err = sub.Close()
			err = errors.Join(err, pub.Close())
			return err
		}
}

func (m natsMessagesPublisher) Publish(_ context.Context, botUUID string, userID int64, msg bots.Message) error {
	dto := mapBotMessageToNATS(botUUID, userID, msg)

	b, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	wmMsg := message.NewMessage(watermill.NewUUID(), b)
	return m.pub.Publish(messagesTopic, wmMsg)
}

type natsBotMessage struct {
	BotUUID string   `json:"bot_uuid"`
	UserID  int64    `json:"user_id"`
	Text    string   `json:"text"`
	Buttons []string `json:"buttons"`
}

func mapBotMessageToNATS(botUUID string, userID int64, msg bots.Message) ampqBotMessage {
	return ampqBotMessage{
		BotUUID: botUUID,
		UserID:  userID,
		Text:    msg.Text,
		Buttons: msg.Buttons,
	}
}
