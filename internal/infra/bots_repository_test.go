package infra_test

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/infra"
)

func TestPgBotsRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	url := os.Getenv("DATABASE_URI")
	db := sqlx.MustConnect("postgres", url)
	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})

	repos := infra.NewPgBotsRepository(db)
	testBotsRepository(t, repos)
}

func testBotsRepository(t *testing.T, repos bots.Repository) {
	t.Parallel()

	t.Run("should save bot", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		bot := createBot(gofakeit.UUID())

		err := repos.UpdateOrCreate(ctx, bot)
		require.NoError(t, err)

		got, err := repos.Bot(ctx, bot.UUID)
		require.NoError(t, err)
		require.NotNil(t, got)
		requireBot(t, *bot, *got)
	})

	t.Run("should replace bot", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		botUUID := gofakeit.UUID()
		bot := createBot(botUUID)

		err := repos.UpdateOrCreate(ctx, bot)
		require.NoError(t, err)

		newBot := createBot(botUUID)
		err = repos.UpdateOrCreate(ctx, newBot)
		require.NoError(t, err)

		got, err := repos.Bot(ctx, bot.UUID)
		require.NoError(t, err)
		require.NotNil(t, got)
		requireBot(t, *bot, *got)
	})

	t.Run("should return error if bot not found", func(t *testing.T) {
		t.Parallel()

		fakeUUID := gofakeit.UUID()
		_, err := repos.Bot(context.Background(), fakeUUID)
		require.ErrorAs(t, err, &bots.BotNotFoundError{})
		require.EqualError(t, err, fmt.Sprintf("bot not found: %s", fakeUUID))
	})

	t.Run("should delete bot", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		bot := createBot(gofakeit.UUID())

		err := repos.UpdateOrCreate(ctx, bot)
		require.NoError(t, err)

		err = repos.Delete(ctx, bot.UUID)
		require.NoError(t, err)

		_, err = repos.Bot(ctx, bot.UUID)
		require.ErrorAs(t, err, &bots.BotNotFoundError{})
	})

	t.Run("should return owner's bots", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		ownerUUID := gofakeit.UUID()

		bot1 := createBot(ownerUUID)
		err := repos.UpdateOrCreate(context.Background(), bot1)
		require.NoError(t, err)

		bot2 := createBot(ownerUUID)
		err = repos.UpdateOrCreate(context.Background(), bot2)
		require.NoError(t, err)

		bs, err := repos.UserBots(ctx, ownerUUID)
		expected := []*bots.Bot{bot1, bot2}
		requireBotsSlices(t, expected, bs)
	})

	t.Run("should return empty list if owner has no bots", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		ownerUUID := gofakeit.UUID()

		bs, err := repos.UserBots(ctx, ownerUUID)
		require.NoError(t, err)
		require.Empty(t, bs)
	})

	t.Run("should update status", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		bot := createBot(gofakeit.UUID())
		err := repos.UpdateOrCreate(ctx, bot)
		require.NoError(t, err)

		err = repos.UpdateStatus(ctx, bot.UUID, bots.Started)
		require.NoError(t, err)

		got, err := repos.Bot(ctx, bot.UUID)
		require.NoError(t, err)
		require.Equal(t, bots.Started, got.Status)
	})
}

func createBot(ownerUUID string) *bots.Bot {
	return bots.MustNewBot(
		gofakeit.UUID(),
		ownerUUID,
		[]bots.EntryPoint{
			bots.MustNewEntryPoint("start", 1),
			bots.MustNewEntryPoint("mailing-1", 1),
		},
		[]bots.Mailing{
			bots.MustNewMailing("example", "mailing-1", 0),
		},
		[]bots.Block{
			bots.MustNewSelectionBlock(1, 2, []bots.Option{
				bots.MustNewOption("To 3", 3),
				bots.MustNewOption("To 4", 4),
			},
				"Selection",
				"Choose option",
			),
			bots.MustNewMessageBlock(2, 1, "Error", "Choose one option!"),
			bots.MustNewQuestionBlock(3, 0, "Question 3", "Some question"),
			bots.MustNewQuestionBlock(4, 0, "Question 4", "Some question"),
		},
		gofakeit.Name(),
		gofakeit.UUID(),
	)
}

func equalOptions(a, b []bots.Option) bool {
	if len(a) != len(b) {
		return false
	}

	a = slices.Clone(a)
	b = slices.Clone(b)

	optionComparator := func(a, b bots.Option) int {
		return strings.Compare(a.Text, b.Text)
	}

	slices.SortFunc(a, optionComparator)
	slices.SortFunc(b, optionComparator)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalBlocks(a, b bots.Block) bool {
	return a.Type == b.Type &&
		a.State == b.State &&
		a.Title == b.Title &&
		a.Text == b.Text &&
		equalOptions(a.Options, b.Options)
}

func equalBlocksSlices(a, b []bots.Block) bool {
	if len(a) != len(b) {
		return false
	}

	a = slices.Clone(a)
	b = slices.Clone(b)

	blockComparator := func(a, b bots.Block) int {
		return a.State - b.State
	}

	slices.SortFunc(a, blockComparator)
	slices.SortFunc(b, blockComparator)

	for i := range a {
		if !equalBlocks(a[i], b[i]) {
			return false
		}
	}
	return true
}

func equalEntries(a, b bots.EntryPoint) bool {
	return a == b
}

func equalEntriesSlices(a, b []bots.EntryPoint) bool {
	if len(a) != len(b) {
		return false
	}

	a = slices.Clone(a)
	b = slices.Clone(b)

	entryComparator := func(a, b bots.EntryPoint) int {
		return strings.Compare(a.Key, b.Key)
	}

	slices.SortFunc(a, entryComparator)
	slices.SortFunc(b, entryComparator)

	for i := range a {
		if !equalEntries(a[i], b[i]) {
			return false
		}
	}
	return true
}

func equalMailings(a, b bots.Mailing) bool {
	return a == b
}

func equalMailingsSlices(a, b []bots.Mailing) bool {
	if len(a) != len(b) {
		return false
	}

	a = slices.Clone(a)
	b = slices.Clone(b)

	mailingComparator := func(a, b bots.Mailing) int {
		return strings.Compare(a.EntryKey, b.EntryKey)
	}

	slices.SortFunc(a, mailingComparator)
	slices.SortFunc(b, mailingComparator)

	for i := range a {
		if !equalMailings(a[i], b[i]) {
			return false
		}
	}
	return true
}

func equalBots(a, b bots.Bot) bool {
	return a.UUID == b.UUID &&
		a.OwnerUUID == b.OwnerUUID &&
		a.Name == b.Name &&
		a.Token == b.Token &&
		a.CreatedAt.Sub(b.CreatedAt).Abs() < time.Microsecond &&
		a.UpdatedAt.Sub(b.UpdatedAt).Abs() < time.Microsecond &&
		equalEntriesSlices(a.Entries(), b.Entries()) &&
		equalMailingsSlices(a.Mailings(), b.Mailings()) &&
		equalBlocksSlices(a.Blocks(), b.Blocks())
}

func requireBot(t *testing.T, expected bots.Bot, actual bots.Bot) {
	require.Truef(t,
		equalBots(expected, actual),
		"expected %v, got %v", expected, actual,
	)
}

func equalBotsSlices(a, b []*bots.Bot) bool {
	if len(a) != len(b) {
		return false
	}

	a = slices.Clone(a)
	b = slices.Clone(b)

	botsComparator := func(a, b *bots.Bot) int {
		return strings.Compare(a.UUID, b.UUID)
	}

	slices.SortFunc(a, botsComparator)
	slices.SortFunc(b, botsComparator)

	for i := range a {
		if !equalBots(*a[i], *b[i]) {
			return false
		}
	}
	return true
}

func requireBotsSlices(t *testing.T, expected, actual []*bots.Bot) {
	require.Truef(t,
		equalBotsSlices(expected, actual),
		"expected %v, got %v", expected, actual,
	)
}
