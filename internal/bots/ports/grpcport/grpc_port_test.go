package grpcport_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/ports/grpcport"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/service"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/server"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/tests"

	botspb "github.com/bmstu-itstech/itsreg-bots/api/grpc/gen/bots"
)

var messageSubscriber message.Subscriber

func TestGrpcPort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Parallel()

	client := newBotsGRPCClient(t)
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second, cancel)

	ch, err := messageSubscriber.Subscribe(ctx, "messages")
	require.NoError(t, err)

	botUUID := gofakeit.UUID()

	_, err = client.CreateBot(ctx, &botspb.CreateBotRequest{
		BotUUID: botUUID,
		Name:    gofakeit.Name(),
		Token:   gofakeit.UUID(),
		Entries: []*botspb.EntryPoint{
			{Key: "start", State: 1},
		},
		Blocks: []*botspb.Block{
			{
				Type:      botspb.BlockType_BlockMessage,
				State:     1,
				NextState: 2,
				Title:     "Greeting",
				Text:      "Hello, user!",
			},
			{
				Type:      botspb.BlockType_BlockQuestion,
				State:     2,
				NextState: 3,
				Title:     "User's name",
				Text:      "What's your name?",
			},
			{
				Type:      botspb.BlockType_BlockSelection,
				State:     3,
				NextState: 4,
				Options: []*botspb.Option{
					{Text: "Red", Next: 5},
					{Text: "Blue", Next: 6},
				},
				Title: "Pill's color",
				Text:  "Choose a pill: red or blue",
			},
			{
				Type:      botspb.BlockType_BlockMessage,
				State:     4,
				NextState: 3,
				Title:     "Error",
				Text:      "Oops, choose one: Red or Blue!",
			},
			{
				Type:  botspb.BlockType_BlockMessage,
				State: 5,
				Title: "Red pill",
				Text:  "Okay, you are choose red pill...",
			},
			{
				Type:  botspb.BlockType_BlockMessage,
				State: 6,
				Title: "Blue pill",
				Text:  "Okay, you are choose blue pill...",
			},
		},
	})
	require.NoError(t, err)

	t.Run("should process participants", func(t *testing.T) {
		t.Parallel()

		ivan := newParticipant(client, botUUID)
		require.NoError(t, ivan.Entry(ctx, "start"))
		require.NoError(t, ivan.Process(ctx, "Ivan"))
		require.NoError(t, ivan.Process(ctx, "Red"))

		john := newParticipant(client, botUUID)
		require.NoError(t, john.Entry(ctx, "start"))
		require.NoError(t, john.Process(ctx, "John"))
		require.NoError(t, john.Process(ctx, "Green"))
		require.NoError(t, john.Process(ctx, "Blue"))

		time.Sleep(100 * time.Millisecond)

		rec, err := receiveMessagesFromBot(ch)
		require.NoError(t, err)

		require.Contains(t, rec, newPlainBotMessage(botUUID, ivan.UserID, "Hello, user!"))
		require.Contains(t, rec, newPlainBotMessage(botUUID, ivan.UserID, "What's your name?"))
		require.Contains(t, rec,
			newBotMessageWithButtons(botUUID, ivan.UserID, "Choose a pill: red or blue", []string{"Red", "Blue"}))
		require.Contains(t, rec, newPlainBotMessage(botUUID, ivan.UserID, "Okay, you are choose red pill..."))

		require.Contains(t, rec, newPlainBotMessage(botUUID, john.UserID, "Hello, user!"))
		require.Contains(t, rec, newPlainBotMessage(botUUID, john.UserID, "What's your name?"))
		require.Contains(t, rec,
			newBotMessageWithButtons(botUUID, ivan.UserID, "Choose a pill: red or blue", []string{"Red", "Blue"}))
		require.Contains(t, rec, newPlainBotMessage(botUUID, john.UserID, "Oops, choose one: Red or Blue!"))
		require.Contains(t, rec, newPlainBotMessage(botUUID, john.UserID, "Okay, you are choose blue pill..."))
	})

	t.Run("should return error if bot not found", func(t *testing.T) {
		t.Parallel()

		fakeBotUUID := gofakeit.UUID()
		_, err := client.GetBot(ctx, &botspb.GetBotRequest{BotUUID: fakeBotUUID})
		require.Equal(t,
			status.Error(codes.NotFound, fmt.Sprintf("bot not found: %s", fakeBotUUID)), err,
			"gRPC error does not match",
		)
	})
}

func TestMain(m *testing.M) {
	if !startService() {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func startService() bool {
	app, sub := service.NewComponentTestApplication()
	messageSubscriber = sub

	port := os.Getenv("PORT")
	addr := fmt.Sprintf("localhost:%s", port)
	go server.RunGRPCServerOnAddr(addr, func(server *grpc.Server) {
		grpcport.RegisterGRPCServer(server, app)
	})

	ok := tests.WaitForPort(addr)
	if !ok {
		log.Println("Timed out waiting for auth gRPC to come up")
	}

	return ok
}

func newBotsGRPCClient(t *testing.T) botspb.BotsServiceClient {
	port := os.Getenv("PORT")
	addr := fmt.Sprintf("localhost:%s", port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	return botspb.NewBotsServiceClient(conn)
}

type participant struct {
	client  botspb.BotsServiceClient
	botUUID string
	UserID  int64
}

func newParticipant(client botspb.BotsServiceClient, botUUID string) *participant {
	return &participant{
		client:  client,
		botUUID: botUUID,
		UserID:  gofakeit.Int64(),
	}
}

func (p *participant) Entry(ctx context.Context, key string) error {
	_, err := p.client.Entry(ctx, &botspb.EntryRequest{
		BotUUID: p.botUUID,
		UserID:  p.UserID,
		Key:     key,
	})
	return err
}

func (p *participant) Process(ctx context.Context, text string) error {
	_, err := p.client.Process(ctx, &botspb.ProcessRequest{
		BotUUID: p.botUUID,
		UserID:  p.UserID,
		Text:    text,
	})
	return err
}

type botMessage struct {
	BotUUID string   `json:"bot_uuid"`
	UserID  int64    `json:"user_id"`
	Text    string   `json:"text"`
	Buttons []string `json:"buttons"`
}

func newPlainBotMessage(botUUID string, userID int64, text string) botMessage {
	return botMessage{
		BotUUID: botUUID,
		UserID:  userID,
		Text:    text,
		Buttons: make([]string, 0),
	}
}

func newBotMessageWithButtons(botUUID string, userID int64, text string, buttons []string) botMessage {
	return botMessage{
		BotUUID: botUUID,
		UserID:  userID,
		Text:    text,
		Buttons: buttons,
	}
}

func (m botMessage) Equal(o botMessage) bool {
	return m.Text == o.Text && slices.Equal(m.Buttons, o.Buttons)
}

func receiveMessagesFromBot(messages <-chan *message.Message) ([]botMessage, error) {
	botMessages := make([]botMessage, 0)

	var rec botMessage
	for msg := range messages {
		err := json.Unmarshal(msg.Payload, &rec)
		if err != nil {
			return nil, err
		}
		botMessages = append(botMessages, rec)

		msg.Ack()
	}

	return botMessages, nil
}
