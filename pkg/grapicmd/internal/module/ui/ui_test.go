package ui

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

type errReader struct {
}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("failed to read")
}

func Test_UI_Confirm(t *testing.T) {
	type TestContext struct {
		in, out, err *bytes.Buffer
		ui           module.UI
	}

	createTestContext := func() *TestContext {
		in := new(bytes.Buffer)
		out := new(bytes.Buffer)
		return &TestContext{
			in:  in,
			out: out,
			ui:  New(out, in),
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
			ctx.in.WriteString(c.input)

			ok, err := ctx.ui.Confirm(c.test)

			if got, want := ctx.out.String(), c.test; !strings.HasPrefix(got, want) {
				t.Errorf("Confirm() wrote %q, want %q", got, want)
			}

			wantErrMsg := "input must be Y or n\n"
			if got, want := strings.Count(ctx.out.String(), wantErrMsg), c.errMsgCnt; got != want {
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
		ui := New(new(bytes.Buffer), &errReader{})

		ok, err := ui.Confirm("test")

		if err == nil {
			t.Error("Confirm() should return an error")
		}

		if ok {
			t.Error("Confirm() should return false")
		}
	})
}
