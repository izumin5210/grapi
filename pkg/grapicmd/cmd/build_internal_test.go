package cmd

import (
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
