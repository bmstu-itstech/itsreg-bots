package module_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-tgservice/internal/domain/module"
	"testing"
)

func TestModule_NewButton(t *testing.T) {
	type args struct {
		text  string
		value string
		next  *module.Module
	}

	type posCase struct {
		test string
		res  *module.Button
		args
	}

	type negCase struct {
		test string
		errs []error
		args
	}

	posCases := []posCase{
		{
			test: "create button successfully",
			args: args{
				text:  "text",
				value: "value",
				next:  &module.Module{},
			},
			res: &module.Button{
				Text:  "text",
				Value: "value",
				Next:  &module.Module{},
			},
		},
	}

	for _, pc := range posCases {
		t.Run(pc.test, func(t *testing.T) {
			btn, errs := module.NewButton(
				pc.text,
				pc.value,
				pc.next)
			require.Nil(t, errs)
			require.NotNil(t, btn)
			require.Equal(t, pc.res.Text, btn.Text)
			require.Equal(t, pc.res.Value, btn.Value)
			require.Equal(t, pc.res.Next, btn.Next)
		})
	}

	negCases := []negCase{
		{
			test: "error to create button with empty text",
			args: args{
				text:  "",
				value: "value",
				next:  &module.Module{},
			},
			errs: []error{
				module.ErrEmptyButtonText,
			},
		},
		{
			test: "error to create button with empty value",
			args: args{
				text:  "text",
				value: "",
				next:  &module.Module{},
			},
			errs: []error{
				module.ErrEmptyButtonValue,
			},
		},
		{
			test: "error to create button with nil next module",
			args: args{
				text:  "text",
				value: "value",
				next:  nil,
			},
			errs: []error{
				module.ErrInvalidButtonNextModule,
			},
		},
	}

	for _, nc := range negCases {
		t.Run(nc.test, func(t *testing.T) {
			_, errs := module.NewButton(
				nc.text,
				nc.value,
				nc.next)
			require.NotNil(t, errs)
			require.ElementsMatch(t, nc.errs, errs, "errors are not equal")
		})
	}
}
