package infra_test

import (
	"context"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand/v2"
	"strconv"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/infra"
)

func TestNATSMessagesPublisher(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pub, messages, closeFn := infra.NewNATSMessagesPublisher()
	t.Cleanup(func() {
		err := closeFn()
		require.NoError(t, err)
	})

	testMessagesPublisher(t, pub, messages)
}

func testMessagesPublisher(t *testing.T, pub bots.MessagesPublisher, messages <-chan *message.Message) {
	t.Run("should publish message", func(t *testing.T) {
		ctx := context.Background()

		botUUID := gofakeit.UUID()
		userID := rand.Int64()
		botMessage := bots.MustNewPlainMessage(gofakeit.Sentence(5))
		err := pub.Publish(ctx, botUUID, userID, botMessage)
		require.NoError(t, err)

		require.Eventually(t,
			func() bool {
				jsonMsg := <-messages
				if jsonMsg == nil {
					return false
				}
				rec, err := unmarshalJSONMessage(jsonMsg)
				if err != nil {
					t.Log(err.Error())
					return false
				}
				jsonMsg.Ack()
				return botUUID == rec.BotUUID && userID == rec.UserID && botMessage.Text == rec.Text
			},
			time.Second, time.Second/10,
		)
	})

	t.Run("should guarantee message order", func(t *testing.T) {
		ctx := context.Background()

		botUUID := gofakeit.UUID()
		userID := rand.Int64()

		const messagesNum = 50

		for i := range messagesNum {
			msg := bots.MustNewPlainMessage(strconv.Itoa(i))
			err := pub.Publish(ctx, botUUID, userID, msg)
			require.NoError(t, err)
		}

		isOrdered := func(s []string) bool {
			for i := range s {
				if strconv.Itoa(i) != s[i] {
					return false
				}
			}
			return true
		}

		require.Eventually(t, func() bool {
			received := make([]string, 0, messagesNum)
			for jsonMsg := range messages {
				rec, err := unmarshalJSONMessage(jsonMsg)
				require.NoError(t, err)
				jsonMsg.Ack()

				received = append(received, rec.Text)

				if len(received) == messagesNum {
					return isOrdered(received)
				}
			}
			return false
		},
			time.Second, time.Second/10,
		)
	})
}

type jsonBotMessage struct {
	BotUUID string   `json:"bot_uuid"`
	UserID  int64    `json:"user_id"`
	Text    string   `json:"text"`
	Buttons []string `json:"buttons"`
}

func unmarshalJSONMessage(msg *message.Message) (jsonBotMessage, error) {
	var res jsonBotMessage
	err := json.Unmarshal(msg.Payload, &res)
	return res, err
}
