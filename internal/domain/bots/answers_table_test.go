package bots_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

func TestNewTable(t *testing.T) {
	t.Run("should create table", func(t *testing.T) {
		// +--------+------------+-----------+
		// | UserID | First name | Last name |
		// +--------+------------+-----------+
		// | 10     | Ivan       | Ivanov    |
		// | 11     | John       |           |
		// +--------+------------+-----------+
		entries := []bots.EntryPoint{
			bots.MustNewEntryPoint("start", 1),
		}
		blocks := []bots.Block{
			bots.MustNewQuestionBlock(1, 2, "First name", "What is your first name?"),
			bots.MustNewQuestionBlock(2, 0, "Last name", "What is your last name?"),
		}
		botUUID := uuid.NewString()
		bot := bots.MustNewBot(botUUID, entries, blocks, "Test bot", "xxx-yyy")

		ivan := bots.MustNewParticipant(botUUID, 10)
		ivan.SwitchTo(1)
		require.NoError(t, ivan.AddAnswer("Ivan"))
		ivan.SwitchTo(2)
		require.NoError(t, ivan.AddAnswer("Ivanov"))

		john := bots.MustNewParticipant(botUUID, 11)
		john.SwitchTo(1)
		require.NoError(t, john.AddAnswer("John"))

		participants := []*bots.Participant{ivan, john}

		table := bots.NewAnswersTable(bot, participants)
		require.NotNil(t, table)

		require.Equal(t, [][]string{
			{"10", "Ivan", "Ivanov"},
			{"11", "John", ""},
		}, table.Body)
		require.Equal(t, []string{
			"UserID", "First name", "Last name",
		}, table.Head)
	})
}
