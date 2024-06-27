package module_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-tgservice/internal/domain/module"
	"testing"
)

func TestCreateModule(t *testing.T) {
	type args struct {
		title    string
		text     string
		isSilent bool
		typ      module.Type
		next     *module.Module
		buttons  []module.Button
	}

	type posCase struct {
		test string
		res  *module.Module
		args
	}

	type negCase struct {
		test string
		errs []error
		args
	}

	posCases := []posCase{
		{
			test: "created module successfully",
			args: args{
				title:    "title",
				text:     "text",
				isSilent: true,
				typ:      module.String,
				next:     nil,
				buttons:  make([]module.Button, 0),
			},
			res: &module.Module{
				Title:    "title",
				Text:     "text",
				IsSilent: true,
				Type:     module.String,
				Next:     nil,
				Buttons:  make([]module.Button, 0),
			},
		},
	}

	for _, pc := range posCases {
		mod, errs := module.New(
			pc.title,
			pc.text,
			pc.isSilent,
			pc.typ,
			pc.next,
			pc.buttons)
		require.Nil(t, errs)
		require.NotNil(t, mod)
		require.Equal(t, pc.res.Title, mod.Title)
		require.Equal(t, pc.res.Text, mod.Text)
		require.Equal(t, pc.res.IsSilent, mod.IsSilent)
		require.Equal(t, pc.res.Type, mod.Type)
		require.Equal(t, pc.res.Next, mod.Next)
		require.Equal(t, pc.res.Buttons, mod.Buttons)
	}

	negCases := []negCase{
		{
			test: "error to create module with empty title",
			args: args{
				title:    "",
				text:     "text",
				isSilent: true,
				typ:      module.String,
				next:     nil,
				buttons:  make([]module.Button, 0),
			},
			errs: []error{
				module.ErrEmptyModuleTitle,
			},
		},
		{
			test: "error to create module with empty text",
			args: args{
				title:    "title",
				text:     "",
				isSilent: true,
				typ:      module.String,
				next:     nil,
				buttons:  make([]module.Button, 0),
			},
			errs: []error{
				module.ErrEmptyModuleText,
			},
		},
		{
			test: "error to create module with invalid type",
			args: args{
				title:    "title",
				text:     "text",
				isSilent: true,
				typ:      0,
				next:     nil,
				buttons:  make([]module.Button, 0),
			},
			errs: []error{
				module.ErrInvalidModuleType,
			},
		},
	}

	for _, nc := range negCases {
		t.Run(nc.test, func(t *testing.T) {
			_, errs := module.New(
				nc.title,
				nc.text,
				nc.isSilent,
				nc.typ,
				nc.next,
				nc.buttons)
			require.NotNil(t, errs)
			require.ElementsMatch(t, nc.errs, errs, "errors are not equal")
		})
	}
}
