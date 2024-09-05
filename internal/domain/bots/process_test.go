package bots_test

import (
	"math/rand/v2"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

func TestBot_Process(t *testing.T) {
	/*
		TEST BOT SCHEME
		            4
		          /   \
		          \   /
		1 --> 2 --> 3 --> 5 --> 6 --> 0
		             \
		              \ --> 7 --> 0
	*/
	greetingBlock := bots.MustNewMessageBlock(1, 2, "Greeting", "Hello, user!")
	usernameBlock := bots.MustNewQuestionBlock(2, 3, "User's name", "What's your name?")
	selectionBlock := bots.MustNewSelectionBlock(3, 4, []bots.Option{
		bots.MustNewOption("To 5", 5),
		bots.MustNewOption("To 7", 7),
	}, "Next", "Choose next step")
	errorBlock := bots.MustNewMessageBlock(4, 3, "Error", "Choose one option!")
	redirectBlock := bots.MustNewMessageBlock(5, 6, "Redirect", "Redirecting to block 6...")
	endMessageBlock := bots.MustNewMessageBlock(6, 0, "End", "End message")
	endQuestionBlock := bots.MustNewQuestionBlock(7, 0, "End", "End question")

	blocks := []bots.Block{
		greetingBlock,
		usernameBlock,
		selectionBlock,
		errorBlock,
		redirectBlock,
		endMessageBlock,
		endQuestionBlock,
	}
	entries := []bots.EntryPoint{
		bots.MustNewEntryPoint("start", 1),
	}
	botUUID := uuid.NewString()
	bot := bots.MustNewBot(botUUID, entries, blocks, "Test bot", "random token")

	t.Run("should entry bot", func(t *testing.T) {
		userID := rand.Int64()
		prt := bots.MustNewParticipant(botUUID, userID)

		resp, err := bot.Entry(prt, "start")
		require.NoError(t, err)
		requireMessages(t, []bots.Message{
			bots.MustNewPlainMessage(greetingBlock.Text),
			bots.MustNewPlainMessage(usernameBlock.Text),
		}, resp)
		requireAnswers(t, []bots.Answer{}, prt.Answers())
	})

	t.Run("should process answer for question", func(t *testing.T) {
		userID := rand.Int64()
		prt := bots.MustNewParticipant(botUUID, userID)
		prt.SwitchTo(2)

		name := "Ivan"
		resp, err := bot.Process(prt, name)
		require.NoError(t, err)
		requireMessages(t, []bots.Message{
			bots.MustNewMessageWithButtons(selectionBlock.Text, selectionBlock.Options),
		}, resp)
		requireAnswers(t, []bots.Answer{
			bots.MustNewAnswer(2, name),
		}, prt.Answers())
	})

	t.Run("should recursively return to block", func(t *testing.T) {
		userID := rand.Int64()
		prt := bots.MustNewParticipant(botUUID, userID)
		prt.SwitchTo(3)

		invalidAnswer := "invalid option"
		resp, err := bot.Process(prt, invalidAnswer)
		require.NoError(t, err)
		requireMessages(t, []bots.Message{
			bots.MustNewPlainMessage(errorBlock.Text),
			bots.MustNewMessageWithButtons(selectionBlock.Text, selectionBlock.Options),
		}, resp)
		requireAnswers(t, []bots.Answer{
			bots.MustNewAnswer(3, invalidAnswer),
		}, prt.Answers())
	})

	t.Run("should overwrite answer", func(t *testing.T) {
		userID := rand.Int64()
		prt := bots.MustNewParticipant(botUUID, userID)
		prt.SwitchTo(3)

		invalidAnswer := "invalid option"
		_, err := bot.Process(prt, invalidAnswer)
		require.NoError(t, err)

		validAnswer := "To 7"
		_, err = bot.Process(prt, validAnswer)
		requireAnswers(t, []bots.Answer{
			bots.MustNewAnswer(3, validAnswer),
		}, prt.Answers())
	})

	t.Run("should end bot", func(t *testing.T) {
		userID := rand.Int64()
		prt := bots.MustNewParticipant(botUUID, userID)
		prt.SwitchTo(3)

		resp, err := bot.Process(prt, "To 5")
		require.NoError(t, err)
		requireMessages(t, []bots.Message{
			bots.MustNewPlainMessage(redirectBlock.Text),
			bots.MustNewPlainMessage(endMessageBlock.Text),
		}, resp)
		requireAnswers(t, []bots.Answer{
			bots.MustNewAnswer(3, "To 5"),
		}, prt.Answers())
	})

	t.Run("should end bot after question", func(t *testing.T) {
		userID := rand.Int64()
		prt := bots.MustNewParticipant(botUUID, userID)
		prt.SwitchTo(7)

		resp, err := bot.Process(prt, "something")
		require.NoError(t, err)
		requireMessages(t, []bots.Message{}, resp)
		requireAnswers(t, []bots.Answer{
			bots.MustNewAnswer(7, "something"),
		}, prt.Answers())
	})

	t.Run("should clear answers if user re-entry bot", func(t *testing.T) {
		userID := rand.Int64()
		prt := bots.MustNewParticipant(botUUID, userID)

		_, err := bot.Entry(prt, "start")
		require.NoError(t, err)

		_, err = bot.Process(prt, "Ivan")
		require.NoError(t, err)
		requireAnswers(t, []bots.Answer{
			bots.MustNewAnswer(2, "Ivan"),
		}, prt.Answers())

		_, err = bot.Entry(prt, "start")
		require.NoError(t, err)
		requireAnswers(t, []bots.Answer{}, prt.Answers())
	})
}

func requireMessages(t *testing.T, expected []bots.Message, actual []bots.Message) {
	require.Lenf(t, actual, len(expected), "expected %d messages, got %d", len(expected), len(actual))
	for i, msg := range actual {
		require.Truef(t, expected[i].Equal(msg), "expected message %v, got %v", expected[i], msg)
	}
}

func requireAnswers(t *testing.T, expected []bots.Answer, actual []bots.Answer) {
	require.Lenf(t, actual, len(expected), "expected %d answers, got %d", len(expected), len(actual))
	for i, ans := range actual {
		require.Equalf(t, expected[i], ans, "expected answer %v, got %v", expected[i], ans)
	}
}
