package mocks

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

const runnerTopic = "runner"

type mockRunnerPublisher struct {
	pub message.Publisher
}

func NewMockRunnerPublisher() (bots.RunnerPublisher, <-chan *message.Message) {
	goCh := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})

	messages, err := goCh.Subscribe(context.Background(), messagesTopic)
	if err != nil {
		panic(err)
	}

	return &mockRunnerPublisher{
		pub: goCh,
	}, messages
}

func (s *mockRunnerPublisher) PublishStart(ctx context.Context, botUUID string) error {
	return s.publish(ctx, botUUID, "start")
}

func (s *mockRunnerPublisher) PublishStop(ctx context.Context, botUUID string) error {
	return s.publish(ctx, botUUID, "stop")
}

func (s *mockRunnerPublisher) publish(_ context.Context, botUUID string, command string) error {
	b, err := json.Marshal(runnerMessage{
		BotUUID: botUUID,
		Command: command,
	})
	if err != nil {
		return err
	}

	wmMsg := message.NewMessage(watermill.NewUUID(), b)
	return s.pub.Publish(runnerTopic, wmMsg)
}

type runnerMessage struct {
	BotUUID string `json:"bot_uuid"`
	Command string `json:"command"`
}
