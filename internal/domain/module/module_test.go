package module_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-tgservice/internal/domain/module"
	"testing"
)

func TestModule_New(t *testing.T) {
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
			test: "created common module successfully",
			args: args{
				title:    "title",
				text:     "text",
				isSilent: false,
				typ:      module.String,
				next:     &module.Module{},
				buttons:  nil,
			},
			res: &module.Module{
				Title:    "title",
				Text:     "text",
				IsSilent: false,
				Type:     module.String,
				Next:     &module.Module{},
				Buttons:  make([]module.Button, 0),
			},
		},
		{
			test: "created module with buttons successfully",
			args: args{
				title:    "title",
				text:     "text",
				isSilent: false,
				typ:      module.String,
				next:     &module.Module{},
				buttons:  []module.Button{{}},
			},
			res: &module.Module{
				Title:    "title",
				Text:     "text",
				IsSilent: false,
				Type:     module.String,
				Next:     &module.Module{},
				Buttons:  []module.Button{{}},
			},
		},
		{
			test: "created silent module successfully",
			args: args{
				title:    "title",
				text:     "text",
				isSilent: true,
				typ:      module.String,
				next:     &module.Module{},
				buttons:  nil,
			},
			res: &module.Module{
				Title:    "title",
				Text:     "text",
				IsSilent: true,
				Type:     module.String,
				Next:     &module.Module{},
				Buttons:  make([]module.Button, 0),
			},
		},
	}

	for _, pc := range posCases {
		t.Run(pc.test, func(t *testing.T) {
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
		})
	}

	negCases := []negCase{
		{
			test: "error to create module with empty title",
			args: args{
				title:    "",
				text:     "text",
				isSilent: true,
				typ:      module.String,
				next:     &module.Module{},
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
				next:     &module.Module{},
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
				next:     &module.Module{},
				buttons:  make([]module.Button, 0),
			},
			errs: []error{
				module.ErrInvalidModuleType,
			},
		},
		{
			test: "error to create silent module with buttons",
			args: args{
				title:    "title",
				text:     "text",
				isSilent: true,
				typ:      module.Number,
				next:     &module.Module{},
				buttons:  []module.Button{{}},
			},
			errs: []error{
				module.ErrInvalidModuleSilent,
			},
		},
		{
			test: "error to create last module with buttons",
			args: args{
				title:    "title",
				text:     "text",
				isSilent: false,
				typ:      module.Number,
				next:     nil,
				buttons:  []module.Button{{}},
			},
			errs: []error{
				module.ErrInvalidLastModule,
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

func TestModule_IsLast(t *testing.T) {
	t.Run("module is last", func(t *testing.T) {
		mod := &module.Module{
			Title:    "title",
			Text:     "text",
			IsSilent: true,
			Type:     module.String,
			Next:     nil,
			Buttons:  make([]module.Button, 0),
		}
		require.True(t, mod.IsLast())
	})

	t.Run("module is not last", func(t *testing.T) {
		mod := &module.Module{
			Title:    "title",
			Text:     "text",
			IsSilent: true,
			Type:     module.String,
			Next:     &module.Module{},
			Buttons:  make([]module.Button, 0),
		}
		require.False(t, mod.IsLast())
	})
}

func TestModule_HasButtons(t *testing.T) {
	t.Run("module has buttons", func(t *testing.T) {
		mod := &module.Module{
			Title:    "title",
			Text:     "text",
			IsSilent: true,
			Type:     module.String,
			Next:     &module.Module{},
			Buttons:  []module.Button{{}},
		}
		require.True(t, mod.HasButtons())
	})

	t.Run("module has no buttons", func(t *testing.T) {
		mod := &module.Module{
			Title:    "title",
			Text:     "text",
			IsSilent: true,
			Type:     module.String,
			Next:     &module.Module{},
			Buttons:  make([]module.Button, 0),
		}
		require.False(t, mod.HasButtons())
	})
}
