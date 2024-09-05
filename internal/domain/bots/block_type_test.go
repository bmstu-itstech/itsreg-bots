package bots_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

func TestNewBlockTypeFromString(t *testing.T) {
	t.Run("should create message block type", func(t *testing.T) {
		bt, err := bots.NewBlockTypeFromString("message")
		require.NoError(t, err)
		require.Equal(t, bt, bots.MessageBlock)
	})

	t.Run("should create question block type", func(t *testing.T) {
		bt, err := bots.NewBlockTypeFromString("question")
		require.NoError(t, err)
		require.Equal(t, bt, bots.QuestionBlock)
	})

	t.Run("should create selection block type", func(t *testing.T) {
		bt, err := bots.NewBlockTypeFromString("selection")
		require.NoError(t, err)
		require.Equal(t, bt, bots.SelectionBlock)
	})
}
