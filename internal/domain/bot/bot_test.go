package bot_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/domain/bot"
	"github.com/zhikh23/itsreg-bots/internal/domain/module"
	"testing"
)

func TestBot_New(t *testing.T) {
	type args struct {
		ownerId int32
		name    string
		tgToken string
		gsToken string
		limit   int32
		start   *module.Module
	}

	type posCase struct {
		test string
		res  *bot.Bot
		args
	}

	type negCase struct {
		test string
		errs []error
		args
	}

	posCases := []posCase{
		{
			test: "successfully create a bot",
			args: args{
				ownerId: 1,
				name:    "example",
				tgToken: "XXXXXXXXXX:XXXXXXXXXXXXXXXXXXXXXXXXXX-XXXXXXXX",
				gsToken: "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY",
				limit:   1,
				start:   &module.Module{},
			},
			res: &bot.Bot{
				OwnerId: 1,
				Name:    "example",
				TgToken: "XXXXXXXXXX:XXXXXXXXXXXXXXXXXXXXXXXXXX-XXXXXXXX",
				GsToken: "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY",
				Limit:   1,
				Start:   &module.Module{},
			},
		},
	}

	for _, pc := range posCases {
		t.Run(pc.test, func(t *testing.T) {
			b, errs := bot.New(
				pc.ownerId,
				pc.name,
				pc.tgToken,
				pc.gsToken,
				pc.limit,
				pc.start)
			require.Nil(t, errs)
			require.NotNil(t, b)
			require.Equal(t, pc.res.OwnerId, b.OwnerId)
			require.Equal(t, pc.res.Name, b.Name)
			require.Equal(t, pc.res.TgToken, b.TgToken)
			require.Equal(t, pc.res.GsToken, b.GsToken)
			require.Equal(t, pc.res.Limit, b.Limit)
			require.Equal(t, pc.res.Start, b.Start)
		})
	}

	negCases := []negCase{
		{
			test: "error to create bot with invalid owner id",
			errs: []error{
				bot.ErrInvalidOwnerId,
			},
			args: args{
				ownerId: 0,
				name:    "example",
				tgToken: "XXXXXXXXXX:XXXXXXXXXXXXXXXXXXXXXXXXXX-XXXXXXXX",
				gsToken: "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY",
				limit:   1,
				start:   &module.Module{},
			},
		},
		{
			test: "error to create bot with empty name",
			errs: []error{
				bot.ErrEmptyName,
			},
			args: args{
				ownerId: 1,
				name:    "",
				tgToken: "XXXXXXXXXX:XXXXXXXXXXXXXXXXXXXXXXXXXX-XXXXXXXX",
				gsToken: "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY",
				limit:   1,
				start:   &module.Module{},
			},
		},
		{
			test: "error to create bot with empty tg token",
			errs: []error{
				bot.ErrEmptyTgToken,
			},
			args: args{
				ownerId: 1,
				name:    "example",
				tgToken: "",
				gsToken: "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY",
				limit:   1,
				start:   &module.Module{},
			},
		},
		{
			test: "error to create bot with empty gs token",
			errs: []error{
				bot.ErrEmptyGsToken,
			},
			args: args{
				ownerId: 1,
				name:    "example",
				tgToken: "XXXXXXXXXX:XXXXXXXXXXXXXXXXXXXXXXXXXX-XXXXXXXX",
				gsToken: "",
				limit:   1,
				start:   &module.Module{},
			},
		},
		{
			test: "error to create bot with invalid team limit token",
			errs: []error{
				bot.ErrInvalidLimit,
			},
			args: args{
				ownerId: 1,
				name:    "example",
				tgToken: "XXXXXXXXXX:XXXXXXXXXXXXXXXXXXXXXXXXXX-XXXXXXXX",
				gsToken: "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY",
				limit:   0,
				start:   &module.Module{},
			},
		},
		{
			test: "error to create bot with nil start module",
			errs: []error{
				bot.ErrInvalidStartModule,
			},
			args: args{
				ownerId: 1,
				name:    "example",
				tgToken: "XXXXXXXXXX:XXXXXXXXXXXXXXXXXXXXXXXXXX-XXXXXXXX",
				gsToken: "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY",
				limit:   1,
				start:   nil,
			},
		},
		{
			test: "error to create bot with several errors",
			errs: []error{
				bot.ErrEmptyName,
				bot.ErrInvalidStartModule,
			},
			args: args{
				ownerId: 1,
				name:    "",
				tgToken: "XXXXXXXXXX:XXXXXXXXXXXXXXXXXXXXXXXXXX-XXXXXXXX",
				gsToken: "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY",
				limit:   1,
				start:   nil,
			},
		},
	}

	for _, nc := range negCases {
		t.Run(nc.test, func(t *testing.T) {
			_, errs := bot.New(
				nc.ownerId,
				nc.name,
				nc.tgToken,
				nc.gsToken,
				nc.limit,
				nc.start)
			require.NotNil(t, errs)
			require.ElementsMatch(t, nc.errs, errs, "errors are not equal")
		})
	}
}
