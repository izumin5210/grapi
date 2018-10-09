package cmd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/grapi/pkg/clui"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/testing"
)

func Test_userDefinedCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		args      []string
		buildArgs []interface{}
		runArgs   []interface{}
		verbose   bool
	}{
		{
			args:      []string{},
			buildArgs: []interface{}{},
			runArgs:   []interface{}{},
			verbose:   false,
		},
		{
			args:      []string{"-v"},
			buildArgs: []interface{}{},
			runArgs:   []interface{}{},
			verbose:   true,
		},
		{
			args:      []string{"--", "-v"},
			buildArgs: []interface{}{"-v"},
			runArgs:   []interface{}{},
			verbose:   false,
		},
		{
			args:      []string{"-v", "--", "-v"},
			buildArgs: []interface{}{"-v"},
			runArgs:   []interface{}{},
			verbose:   true,
		},
		{
			args:      []string{"--", "--", "-v"},
			buildArgs: []interface{}{},
			runArgs:   []interface{}{"-v"},
			verbose:   false,
		},
		{
			args:      []string{"-v", "--", "--", "-v"},
			buildArgs: []interface{}{},
			runArgs:   []interface{}{"-v"},
			verbose:   true,
		},
		{
			args:      []string{"-v", "--", "-a", "--", "-b"},
			buildArgs: []interface{}{"-a"},
			runArgs:   []interface{}{"-b"},
			verbose:   true,
		},
		{
			args:      []string{"-v", "--", "-foo", "bar", "-baz", "--", "qux", "quux"},
			buildArgs: []interface{}{"-foo", "bar", "-baz"},
			runArgs:   []interface{}{"qux", "quux"},
			verbose:   true,
		},
	}

	name := "testcommand"

	for _, c := range cases {
		loader := moduletesting.NewMockScriptLoader(ctrl)
		script := moduletesting.NewMockScript(ctrl)

		var verbose bool
		cmd := newUserDefinedCommand(clui.Nop, loader, name)
		cmd.SetArgs(c.args)
		cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "")

		loader.EXPECT().Get(name).Return(script, true)

		script.EXPECT().Build(c.buildArgs...)
		script.EXPECT().Run(c.runArgs...)
		script.EXPECT().Name().Return(name)

		err := cmd.Execute()

		if err != nil {
			t.Errorf("Execute() returned an error %v", err)
		}

		if got, want := verbose, c.verbose; got != want {
			t.Errorf("verbose option is %t, want %t", got, want)
		}
	}
}
