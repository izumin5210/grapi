package cli_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/izumin5210/clig/pkg/clib"

	"github.com/izumin5210/grapi/pkg/cli"
)

func TestUI(t *testing.T) {
	defer func(b bool) { color.NoColor = b }(color.NoColor)
	color.NoColor = true

	want := `  ➜  section 1
  ▸  subsection 1.1
     ✔  created
     ╌  skipped
     ✔  ok

  ▸  subsection 1.2
     ✗  failure
        foobar
        baz

  ➜  section 2
     ✗  fail!!!
`

	io := clib.NewBufferedIO()
	ui := cli.NewUI(io)

	ui.Section("section 1")
	ui.Subsection("subsection 1.1")
	ui.ItemSuccess("created")
	ui.ItemSkipped("skipped")
	ui.ItemSuccess("ok")
	ui.Subsection("subsection 1.2")
	ui.ItemFailure("failure", errors.New("foobar"), errors.New("baz"))
	ui.Section("section 2")
	ui.ItemFailure("fail!!!")

	if got := io.Out.(*bytes.Buffer).String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

type errReader struct {
}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("failed to read")
}

func TestUI_Confirm(t *testing.T) {
	type TestContext struct {
		io *clib.IO
		ui cli.UI
	}

	createTestContext := func() *TestContext {
		io := clib.NewBufferedIO()
		return &TestContext{
			io: io,
			ui: cli.NewUI(io),
		}
	}

	cases := []struct {
		test      string
		errMsgCnt int
		input     string
		output    bool
	}{
		{
			test:      "with inputs 'n'",
			errMsgCnt: 0,
			input:     "n\n",
			output:    false,
		},
		{
			test:      "with inputs 'Y'",
			errMsgCnt: 0,
			input:     "Y\n",
			output:    true,
		},
		{
			test:      "with inputs 2 invalid chars and Y",
			errMsgCnt: 2,
			input:     "y\nN\nY\n",
			output:    true,
		},
		{
			test:      "with inputs 1 invalid chars and n",
			errMsgCnt: 2,
			input:     "N\ny\nn\n",
			output:    false,
		},
	}

	for _, c := range cases {
		t.Run(c.test, func(t *testing.T) {
			ctx := createTestContext()
			ctx.io.In.(*bytes.Buffer).WriteString(c.input)

			ok, err := ctx.ui.Confirm(c.test)

			if got, want := ctx.io.Out.(*bytes.Buffer).String(), c.test; !strings.HasPrefix(got, want) {
				t.Errorf("Confirm() wrote %q, want %q", got, want)
			}

			wantErrMsg := "input must be Y or n\n"
			if got, want := strings.Count(ctx.io.Out.(*bytes.Buffer).String(), wantErrMsg), c.errMsgCnt; got != want {
				t.Errorf("Confirm() wrote %q as error %d times, want %d times", wantErrMsg, got, want)
			}

			if err != nil {
				t.Errorf("Confirm() should not return errors, but returned %v", err)
			}

			if got, want := ok, c.output; got != want {
				t.Errorf("Confirm() returned %t, want %t", got, want)
			}
		})
	}

	t.Run("when failed to read", func(t *testing.T) {
		ui := cli.NewUI(clib.NewBufferedIO())

		ok, err := ui.Confirm("test")

		if err == nil {
			t.Error("Confirm() should return an error")
		}

		if ok {
			t.Error("Confirm() should return false")
		}
	})
}
