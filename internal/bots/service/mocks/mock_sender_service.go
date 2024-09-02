package mocks

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/google/uuid"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"
)

type mockSenderService struct {
	pub message.Publisher
}

func NewMockSenderService() (interfaces.SenderService, message.Subscriber) {
	goCh := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})
	return &mockSenderService{
		pub: goCh,
	}, goCh
}

func (s *mockSenderService) Send(
	ctx context.Context,
	botUUID string,
	userID int64,
	msg bots.Message,
) error {
	b, err := marshall(botUUID, userID, msg)
	if err != nil {
		return err
	}
	wmMsg := message.NewMessage(uuid.NewString(), b)
	return s.pub.Publish("messages", wmMsg)
}

func marshall(botUUID string, userID int64, msg bots.Message) ([]byte, error) {
	type dto struct {
		BotUUID string   `json:"bot_uuid"`
		UserID  int64    `json:"user_id"`
		Text    string   `json:"text"`
		Buttons []string `json:"buttons"`
	}
	d := dto{
		BotUUID: botUUID,
		UserID:  userID,
		Text:    msg.Text,
		Buttons: msg.Buttons,
	}

	return json.Marshal(d)
}
