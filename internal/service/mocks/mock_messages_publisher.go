package mocks

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

const messagesTopic = "messages"

type mockSenderService struct {
	pub message.Publisher
}

func NewMockMessagesPublisher() (bots.MessagesPublisher, <-chan *message.Message) {
	goCh := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})

	messages, err := goCh.Subscribe(context.Background(), messagesTopic)
	if err != nil {
		panic(err)
	}

	return &mockSenderService{
		pub: goCh,
	}, messages
}

func (s *mockSenderService) Publish(
	_ context.Context,
	botUUID string,
	userID int64,
	msg bots.Message,
) error {
	b, err := json.Marshal(botMessage{
		BotUUID: botUUID,
		UserID:  userID,
		Text:    msg.Text,
		Buttons: msg.Buttons,
	})
	if err != nil {
		return err
	}

	wmMsg := message.NewMessage(watermill.NewUUID(), b)
	return s.pub.Publish(messagesTopic, wmMsg)
}

type botMessage struct {
	BotUUID string   `json:"bot_uuid"`
	UserID  int64    `json:"user_id"`
	Text    string   `json:"text"`
	Buttons []string `json:"buttons"`
}
