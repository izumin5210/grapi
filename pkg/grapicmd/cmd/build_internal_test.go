package cmd

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	moduletesting "github.com/izumin5210/grapi/pkg/grapicmd/internal/module/testing"
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

var errBuildFailed = errors.New("error occur")

func Test_newBuildCommandMocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type errBuild struct {
		err error
	}

	testCases := []struct {
		args      []string
		optExpect []string
		// 単体のbuildの結果
		buildResult map[string]errBuild
		// 全体のbuildの結果
		wantErr bool
	}{
		{
			// "aaa"だけ指定した場合，buildには成功する
			// "bbb"と"ccc"は呼ばれない
			args:      []string{"aaa"},
			optExpect: []string{},
			buildResult: map[string]errBuild{
				"aaa": {},
			},
		},
		{
			// なんらかの理由で"bbb"のbuildができずに失敗（args指定なし）
			// "ccc"はbuildされない
			args:      []string{},
			optExpect: []string{},
			buildResult: map[string]errBuild{
				"aaa": {},
				"bbb": {
					err: errBuildFailed,
				},
			},
			wantErr: true,
		},
		{
			// なんらかの理由で"bbb"でエラーが発生した場合（args指定あり）
			// "ccc"は実行されない
			args:      []string{"aaa", "bbb", "ccc"},
			optExpect: []string{},
			buildResult: map[string]errBuild{
				"aaa": {},
				"bbb": {
					err: errBuildFailed,
				},
			},
			wantErr: true,
		},
		{
			// 与えたオプションによってエラーが発生する場合
			// "bbb"以降は実行されない
			args:      []string{"--", "-b", "-c"},
			optExpect: []string{"-b", "-c"},
			buildResult: map[string]errBuild{
				"aaa": {
					err: errBuildFailed,
				},
			},
			wantErr: true,
		},
		{
			// build対象は与えず与えたオプションが正当でありすべて成功する場合
			args:      []string{"--", "-a"},
			optExpect: []string{"-a"},
			buildResult: map[string]errBuild{
				"aaa": {},
				"bbb": {},
				"ccc": {},
			},
		},
		{
			// build対象とオプションの両者を与えず，すべて成功する場合
			args:      []string{},
			optExpect: []string{},
			buildResult: map[string]errBuild{
				"aaa": {},
				"bbb": {},
				"ccc": {},
			},
		},
	}

	ctx := &grapicmd.Ctx{}
	err := ctx.Init()
	if err != nil {
		t.Fatal(err)
	}

	commandNames := []string{"aaa", "bbb", "ccc"}
	for _, c := range testCases {
		loader := moduletesting.NewMockScriptLoader(ctrl)

		loader.EXPECT().Load(gomock.Any()).Return(nil)
		loader.EXPECT().Names().Return(commandNames)
		for _, arg := range commandNames {
			script := moduletesting.NewMockScript(ctrl)
			loader.EXPECT().Get(arg).Return(script, true).AnyTimes()
			script.EXPECT().Name().Return(arg).AnyTimes()
			if buildResult, ok := c.buildResult[arg]; ok {
				script.EXPECT().Build(gomock.Any(), c.optExpect).Return(buildResult.err)
			}
		}

		cmd := newBuildCommandMocked(true, ctx, loader, cli.NopUI)
		cmd.SetArgs(c.args)
		err := cmd.Execute()

		if c.wantErr != (err != nil) {
			t.Errorf("wantErr: %v, gotErr: %v", c.wantErr, err != nil)
		}
	}
}
