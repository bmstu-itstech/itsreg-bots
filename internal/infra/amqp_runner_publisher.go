package infra

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

const (
	amqpRunnerURI = "amqp://guest:guest@localhost:5672/"
	runnerTopic   = "runner"
)

type amqpRunnerPublisher struct {
	pub *amqp.Publisher
}

func NewAmqpRunnerPublisher() (bots.RunnerPublisher, <-chan *message.Message, func() error) {
	amqpConfig := amqp.NewDurableQueueConfig(amqpRunnerURI)

	wmLogger := watermill.NewStdLogger(false, false)

	publisher, err := amqp.NewPublisher(amqpConfig, wmLogger)
	if err != nil {
		panic(err)
	}

	sub, err := amqp.NewSubscriber(amqpConfig, wmLogger)
	if err != nil {
		panic(err)
	}

	messages, err := sub.Subscribe(context.Background(), runnerTopic)
	if err != nil {
		panic(err)
	}

	return &amqpRunnerPublisher{
		pub: publisher,
	}, messages, nil
}

func (s *amqpRunnerPublisher) PublishStart(ctx context.Context, botUUID string) error {
	return s.publish(ctx, botUUID, "start")
}

func (s *amqpRunnerPublisher) PublishStop(ctx context.Context, botUUID string) error {
	return s.publish(ctx, botUUID, "stop")
}

func (s *amqpRunnerPublisher) publish(ctx context.Context, botUUID string, command string) error {
	dto := amqpRunnerMessage{
		BotUUID: botUUID,
		Command: command,
	}

	b, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	wmMsg := message.NewMessage(watermill.NewUUID(), b)
	return s.pub.Publish(runnerTopic, wmMsg)
}

type amqpRunnerMessage struct {
	BotUUID string `json:"bot_uuid"`
	Command string `json:"command"`
}
