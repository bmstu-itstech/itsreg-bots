package bots_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/bots"
)

func TestNewTable(t *testing.T) {
	t.Run("should create table", func(t *testing.T) {
		// +----+------------+-----------+
		// | UserID | First name | Last name |
		// +----+------------+-----------+
		// | 10 | Ivan       | Ivanov    |
		// | 11 | John       |           |
		// +----+------------+-----------+
		blocks := []bots.Block{
			bots.MustNewQuestionBlock(1, 2, "First name", "What is your first name?"),
			bots.MustNewQuestionBlock(2, 0, "Last name", "What is your last name?"),
		}
		bot := bots.MustNewBot(uuid.NewString(), blocks, 1, "Test bot", "xxx-yyy")

		answers := []*bots.Answer{
			bots.MustNewAnswer(10, 1, "Ivan"),
			bots.MustNewAnswer(11, 1, "John"),
			bots.MustNewAnswer(10, 2, "Ivanov"),
		}

		table, err := bots.NewTable(bot, answers)
		require.NoError(t, err)
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
