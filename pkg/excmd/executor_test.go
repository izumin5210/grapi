package excmd_test

import (
	"context"
	"testing"

	clibtesting "github.com/izumin5210/clig/pkg/clib/testing"

	"github.com/izumin5210/grapi/pkg/excmd"
)

func TestExecutor_WithConnectedIO(t *testing.T) {
	cases := []struct {
		test   string
		cmd    string
		opts   []excmd.Option
		out    string
		stdout string
	}{
		{
			test: "not connected",
			cmd:  "/bin/sh",
			opts: []excmd.Option{excmd.WithArgs("-c", "echo foo")},
			out:  "foo\n",
		},
		{
			test:   "connected",
			cmd:    "/bin/sh",
			opts:   []excmd.Option{excmd.WithArgs("-c", "read i && echo $i-$i"), excmd.WithIOConnected()},
			out:    "foo-foo\n",
			stdout: "foo-foo\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.test, func(t *testing.T) {
			io := clibtesting.NewFakeIO()
			io.InBuf.WriteString("foo\n")

			execer := excmd.NewExecutor(io)

			out, err := execer.Exec(context.TODO(), tc.cmd, tc.opts...)
			if err != nil {
				t.Errorf("returned %v, want nil", err)
			}

			if got, want := string(out), tc.out; got != want {
				t.Errorf("returned %q, want %q", got, want)
			}

			if got, want := io.OutBuf.String(), tc.stdout; got != want {
				t.Errorf("printed %q, want %q", got, want)
			}
		})
	}
}
