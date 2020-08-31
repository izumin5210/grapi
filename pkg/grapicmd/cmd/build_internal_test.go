package cmd

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/testing"
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
	newTestCases := []struct {
		args              []string
		argsSplited       []string
		opt               []string
		scriptLoaderNames []string
		// 単体のbuildの結果
		buildResult map[string]struct {
			err error
		}
		// 全体のbuildの結果
		errWhole bool
	}{
		{
			// "aaa"だけ指定した場合，buildには成功する
			// "bbb"と"ccc"は呼ばれない
			args:              []string{"aaa"},
			argsSplited:       []string{"aaa"},
			opt:               []string{},
			scriptLoaderNames: []string{"aaa", "bbb", "bbb"},
			buildResult: map[string]struct {
				err error
			}{
				"aaa": {
					err: nil,
				},
			},
			errWhole: false,
		},
		{
			// なんらかの理由で"bbb"のbuildができずに失敗（args指定なし）
			// "ccc"はbuildされない
			args:              []string{},
			argsSplited:       []string{},
			opt:               []string{},
			scriptLoaderNames: []string{"aaa", "bbb", "ccc"},
			buildResult: map[string]struct {
				err error
			}{
				"aaa": {
					err: nil,
				},
				"bbb": {
					err: errBuildFailed,
				},
			},
			errWhole: true,
		},
		{
			// なんらかの理由で"bbb"でエラーが発生した場合（args指定あり）
			// "ccc"は実行されない
			args:              []string{"aaa", "bbb", "ddd"},
			argsSplited:       []string{"aaa", "bbb", "ddd"},
			opt:               []string{},
			scriptLoaderNames: []string{"aaa", "bbb", "ddd"},
			buildResult: map[string]struct {
				err error
			}{
				"aaa": {
					err: nil,
				},
				"bbb": {
					err: errBuildFailed,
				},
			},
			errWhole: true,
		},
		{
			// 与えたオプションによってエラーが発生する場合
			// "bbb"以降は実行されない
			args:              []string{"--", "-b", "-c"},
			argsSplited:       []string{},
			opt:               []string{"-b", "-c"},
			scriptLoaderNames: []string{"aaa", "bbb", "bbb"},
			buildResult: map[string]struct {
				err error
			}{
				"aaa": {
					err: errBuildFailed,
				},
			},
			errWhole: true,
		},
		{
			// build対象は与えず与えたオプションが正当でありすべて成功する場合
			args:              []string{"--", "-a"},
			argsSplited:       []string{},
			opt:               []string{"-a"},
			scriptLoaderNames: []string{"aaa", "bbb", "ccc"},
			buildResult: map[string]struct {
				err error
			}{
				"aaa": {
					err: nil,
				},
				"bbb": {
					err: nil,
				},
				"ccc": {
					err: nil,
				},
			},
			errWhole: false,
		},
		{
			// build対象とオプションの両者を与えず，すべて成功する場合
			args:              []string{},
			argsSplited:       []string{},
			opt:               []string{},
			scriptLoaderNames: []string{"aaa", "bbb", "ccc"},
			buildResult: map[string]struct {
				err error
			}{
				"aaa": {
					err: nil,
				},
				"bbb": {
					err: nil,
				},
				"ccc": {
					err: nil,
				},
			},
			errWhole: false,
		},
	}

	for _, c := range newTestCases {
		loader := moduletesting.NewMockScriptLoader(ctrl)

		loader.EXPECT().Names().Return(c.scriptLoaderNames)
		for _, arg := range c.scriptLoaderNames {
			script := moduletesting.NewMockScript(ctrl)
			loader.EXPECT().Get(arg).Return(script, true).AnyTimes()
			script.EXPECT().Name().Return(arg).AnyTimes()

			if buildResult, ok := c.buildResult[arg]; len(c.argsSplited) == 0 || ok {
				script.EXPECT().Build(gomock.Any(), c.opt).Return(buildResult.err).Times(1)
				// if once build error occur, not build subsequent build args
				if buildResult.err != nil {
					break
				}
			}
		}

		cmd := newBuildCommandMocked(true, loader, cli.NopUI)
		cmd.SetArgs(c.args)
		err := cmd.Execute()

		if c.errWhole != (err != nil) {
			t.Errorf("wantErr: %v, gotErr: %v", c.errWhole, err != nil)
		}
	}
}
