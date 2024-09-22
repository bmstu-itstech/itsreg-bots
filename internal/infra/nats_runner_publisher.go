package infra

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	nc "github.com/nats-io/nats.go"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/logs"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/logs/sl"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type natsRunnerPublisher struct {
	pub *nats.Publisher
}

func NewNATSRunnerPublisher() (bots.RunnerPublisher, <-chan *message.Message, func() error) {
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

	messages, err := sub.Subscribe(context.Background(), runnerTopic)
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

	return natsRunnerPublisher{
			pub: pub,
		}, messages, func() error {
			err = sub.Close()
			err = errors.Join(err, pub.Close())
			return err
		}
}

func (p natsRunnerPublisher) PublishStart(ctx context.Context, botUUID string) error {
	return p.publish(ctx, botUUID, "start")
}

func (p natsRunnerPublisher) PublishStop(ctx context.Context, botUUID string) error {
	return p.publish(ctx, botUUID, "stop")
}

func (p natsRunnerPublisher) publish(_ context.Context, botUUID string, command string) error {
	dto := natsRunnerMessage{
		BotUUID: botUUID,
		Command: command,
	}

	b, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	wmMsg := message.NewMessage(watermill.NewUUID(), b)
	return p.pub.Publish(runnerTopic, wmMsg)
}

type natsRunnerMessage struct {
	BotUUID string `json:"bot_uuid"`
	Command string `json:"command"`
}
