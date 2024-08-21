package bots_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

func TestBot_Process(t *testing.T) {
	startState := 1
	blocks := []bots.Block{
		bots.MustNewQuestionBlock(startState, 2, "Question 1", "Some text"),
		bots.MustNewSelectionBlock(2, 2, []bots.Option{
			bots.MustNewOption("3", 3),
			bots.MustNewOption("4", 4),
			bots.MustNewOption("5", 5),
			bots.MustNewOption("0", 0),
		}, "Selection", "Some text"),
		bots.MustNewMessageBlock(3, 4, "Message 3", "Some text"),
		bots.MustNewMessageBlock(4, 5, "Message 4", "Some text"),
		bots.MustNewQuestionBlock(5, 6, "Question 5", "Some text"),
		bots.MustNewMessageBlock(6, 0, "Message 6", "Some text"),
	}
	bot := bots.MustNewBot(uuid.NewString(), blocks, startState, "Test bot", "xxx-yyy")

	t.Run("should process question block", func(t *testing.T) {
		prt := bots.MustNewParticipant(1, startState)

		processed, answer, err := bot.Process(prt, "random answer")
		require.NoError(t, err)
		require.NotNil(t, processed)
		require.NotNil(t, answer)

		require.Len(t, processed, 1)
		require.Equal(t, 2, processed[0].State)

		require.Equal(t, prt.ID, answer.UserID)
		require.Equal(t, "random answer", answer.Text)
		require.Equal(t, 1, answer.State)

		require.Equal(t, 2, prt.State)
	})

	t.Run("should process option of selection block", func(t *testing.T) {
		prt := bots.MustNewParticipant(1, 2)

		processed, answer, err := bot.Process(prt, "5")
		require.NoError(t, err)
		require.NotNil(t, processed)
		require.NotNil(t, answer)

		require.Len(t, processed, 1)
		require.Equal(t, 5, processed[0].State)

		require.Equal(t, prt.ID, answer.UserID)
		require.Equal(t, "5", answer.Text)
		require.Equal(t, 2, answer.State)

		require.Equal(t, 5, prt.State)
	})

	t.Run("should process unconditional branch of selection block", func(t *testing.T) {
		prt := bots.MustNewParticipant(1, 2)

		processed, answer, err := bot.Process(prt, "random answer")
		require.NoError(t, err)
		require.NotNil(t, processed)
		require.NotNil(t, answer)

		require.Len(t, processed, 1)
		require.Equal(t, 2, processed[0].State)

		require.Equal(t, prt.ID, answer.UserID)
		require.Equal(t, "random answer", answer.Text)
		require.Equal(t, 2, answer.State)

		require.Equal(t, 2, prt.State)
	})

	t.Run("should process message block", func(t *testing.T) {
		prt := bots.MustNewParticipant(1, 2)

		processed, answer, err := bot.Process(prt, "4")
		require.NoError(t, err)
		require.NotNil(t, processed)
		require.NotNil(t, answer)

		require.Equal(t, prt.ID, answer.UserID)
		require.Equal(t, "4", answer.Text)
		require.Equal(t, 2, answer.State)

		require.Len(t, processed, 2)
		require.Equal(t, 4, processed[0].State)
		require.Equal(t, 5, processed[1].State)

		require.Equal(t, 5, prt.State)
	})

	t.Run("should process some message blocks", func(t *testing.T) {
		prt := bots.MustNewParticipant(1, 2)

		processed, answer, err := bot.Process(prt, "3")
		require.NoError(t, err)
		require.NotNil(t, processed)
		require.NotNil(t, answer)

		require.Equal(t, prt.ID, answer.UserID)
		require.Equal(t, "3", answer.Text)
		require.Equal(t, 2, answer.State)

		require.Len(t, processed, 3)
		require.Equal(t, 3, processed[0].State)
		require.Equal(t, 4, processed[1].State)
		require.Equal(t, 5, processed[2].State)

		require.Equal(t, 5, prt.State)
	})

	t.Run("should finish bots script", func(t *testing.T) {
		prt := bots.MustNewParticipant(1, 2)

		processed, answer, err := bot.Process(prt, "0")
		require.Len(t, processed, 0)
		require.NoError(t, err)
		require.NotNil(t, answer)

		require.Equal(t, prt.ID, answer.UserID)
		require.Equal(t, "0", answer.Text)
		require.Equal(t, 2, answer.State)

		require.Equal(t, 0, prt.State)
	})

	t.Run("should process message block and finish", func(t *testing.T) {
		prt := bots.MustNewParticipant(1, 5)

		processed, answer, err := bot.Process(prt, "answer")
		require.NotNil(t, processed)
		require.NoError(t, err)
		require.NotNil(t, answer)

		require.Len(t, processed, 1)
		require.Equal(t, 6, processed[0].State)

		require.Equal(t, prt.ID, answer.UserID)
		require.Equal(t, "answer", answer.Text)
		require.Equal(t, 5, answer.State)

		require.Equal(t, 0, prt.State)
	})

	t.Run("should ignore if participant already finished", func(t *testing.T) {
		prt := bots.MustNewParticipant(1, bots.FinishState)

		processed, answer, err := bot.Process(prt, "answer")
		require.Len(t, processed, 0)
		require.NoError(t, err)
		require.Nil(t, answer)
	})
}
