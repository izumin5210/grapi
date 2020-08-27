package cmd

import (
	"errors"
	//"github.com/spf13/pflag"
	"github.com/golang/mock/gomock"
	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/testing"

	"reflect"
	"testing"
)

func TestSplitOptions(t *testing.T) {
	tests := []struct {
		input     []string
		outputArg []string
		outputOpt []string
	}{
		{
			[]string{"foo", "bar"},
			[]string{"foo", "bar"},
			[]string{},
		},
		{
			[]string{"foo", "-b"},
			[]string{"foo"},
			[]string{"-b"},
		},
		{
			[]string{"foo", "-b", "-h"},
			[]string{"foo"},
			[]string{"-b", "-h"},
		},
		{
			[]string{"foo", "-b", "-h"},
			[]string{"foo"},
			[]string{"-b", "-h"},
		},
		{
			[]string{"-b", "-h"},
			[]string{},
			[]string{"-b", "-h"},
		},
		{
			[]string{"foo", "-b", "-h", "ooo"},
			[]string{"foo"},
			[]string{"-b", "-h", "ooo"},
		},
		{
			[]string{},
			[]string{},
			[]string{},
		},
	}

	for i, test := range tests {
		gotArg, gotOpt := splitOptions(test.input)
		if !reflect.DeepEqual(test.outputArg, gotArg) {
			t.Errorf("(%v) Expected: %v gotArg: %v", i, test.outputArg, gotArg)
		}
		if !reflect.DeepEqual(test.outputOpt, gotOpt) {
			t.Errorf("(%v) Expected: %v gotOpt: %v", i, test.outputOpt, gotOpt)
		}
	}
}

func Test_newBuildCommandMocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		args              []string
		argsStriped       []string
		opt               []string
		scriptLoaderNames []string
		hasNameSet        []bool
		wantErr           bool
		err               error
	}{
		{
			args:              []string{"aaa", "bbb", "ccc", "--", "-a", "-b", "-c"},
			argsStriped:       []string{"aaa", "bbb", "ccc"},
			opt:               []string{"-a", "-b", "-c"},
			scriptLoaderNames: []string{"aaa", "bbb", "ccc"},
			hasNameSet:        []bool{true, true, true},
			wantErr:           true,
			err:               errors.New("error occur"),
		},
		{
			args:              []string{"aaa", "bbb", "ccc"},
			argsStriped:       []string{"aaa", "bbb", "ccc"},
			opt:               []string{},
			scriptLoaderNames: []string{"aaa", "bbb", "ccc"},
			hasNameSet:        []bool{true, true, true},
			wantErr:           false,
			err:               nil,
		},
		{
			args:              []string{"aaa"},
			argsStriped:       []string{"aaa"},
			opt:               []string{},
			scriptLoaderNames: []string{"aaa", "bbb", "ccc"},
			hasNameSet:        []bool{true, false, false},
			wantErr:           false,
			err:               nil,
		},
		{
			args:              []string{"aaa", "--", "-b", "c"},
			argsStriped:       []string{"aaa"},
			opt:               []string{"-b", "c"},
			scriptLoaderNames: []string{"aaa"},
			hasNameSet:        []bool{true},
			wantErr:           true,
			err:               errors.New("error occur"),
		},
		{
			args:              []string{"--", "-a", "-b"},
			argsStriped:       []string{},
			opt:               []string{"-a", "-b"},
			scriptLoaderNames: []string{"aaa", "bbb", "ccc"},
			hasNameSet:        []bool{false, false, false},
			wantErr:           true,
			err:               errors.New("error occur"),
		},
		{
			args:              []string{"--", "-a", "b"},
			argsStriped:       []string{},
			opt:               []string{"-a", "b"},
			scriptLoaderNames: []string{"aaa", "bbb", "ccc"},
			hasNameSet:        []bool{false, false, false},
			wantErr:           true,
			err:               errors.New("error occur"),
		},
		{
			args:              []string{},
			argsStriped:       []string{},
			opt:               []string{},
			scriptLoaderNames: []string{"aaa", "bbb", "ccc"},
			hasNameSet:        []bool{false, false, false},
			wantErr:           false,
			err:               nil,
		},
	}

	for _, c := range cases {
		loader := moduletesting.NewMockScriptLoader(ctrl)

		loader.EXPECT().Names().Return(c.scriptLoaderNames)

		for i, arg := range c.scriptLoaderNames {
			script := moduletesting.NewMockScript(ctrl)
			loader.EXPECT().Get(arg).Return(script, true)
			script.EXPECT().Name().Return(arg).AnyTimes()

			if len(c.argsStriped) == 0 || c.hasNameSet[i] {
				script.EXPECT().Build(gomock.Any(), c.opt).Return(c.err)
			}
			if c.err != nil {
				break
			}
		}

		cmd := newBuildCommandMocked(loader, cli.NopUI)
		cmd.SetArgs(c.args)
		err := cmd.Execute()

		if c.wantErr != (err != nil) {
			t.Errorf("wantErr: %v, gotErr: %v", c.wantErr, err != nil)
		}
	}
}
