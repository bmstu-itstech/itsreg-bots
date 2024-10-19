package bots_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

func TestBot_AddMailing(t *testing.T) {
	bot := bots.MustNewBot(
		"1234",
		"1234",
		[]bots.EntryPoint{
			bots.MustNewEntryPoint("start", 1),
		},
		[]bots.Mailing{},
		[]bots.Block{
			bots.MustNewMessageBlock(1, 0, "Title", "Test text"),
		},
		"Test bot",
		"12345678:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	)

	require.Empty(t, bot.Mailings())
	err := bot.AddMailing(
		"Mailing 0", 0,
		bots.MustNewEntryPoint("mailing_00", 2),
		[]bots.Block{
			bots.MustNewMessageBlock(2, 0, "Mailing 01", "Test text"),
		},
	)
	require.NoError(t, err)
}
