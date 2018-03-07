package ui

import (
	"fmt"
	"io"

	"github.com/izumin5210/clicontrib/clog"
	"github.com/tcnksm/go-input"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

// New creates a new UI instance.
func New(out io.Writer, in io.Reader) module.UI {
	return &uiImpl{
		out: out,
		inputUI: &input.UI{
			Reader: in,
			Writer: out,
		},
	}
}

type uiImpl struct {
	out       io.Writer
	inSection bool
	inputUI   *input.UI
}

func (u *uiImpl) Output(msg string) {
	u.inSection = true
	fmt.Fprintf(u.out, "    %s\n", msg)
}

func (u *uiImpl) Section(msg string) {
	if u.inSection {
		fmt.Fprintln(u.out)
		u.inSection = false
	}
	printTypeSection.Fprintln(u.out, msg)
}

func (u *uiImpl) Subsection(msg string) {
	if u.inSection {
		fmt.Fprintln(u.out)
		u.inSection = false
	}
	printTypeSubsection.Fprintln(u.out, msg)
}

func (u *uiImpl) Warn(msg string) {
	u.inSection = true
	printTypeWarn.Fprintln(u.out, msg)
}

func (u *uiImpl) Error(msg string) {
	u.inSection = true
	printTypeError.Fprintln(u.out, msg)
}

func (u *uiImpl) ItemSuccess(msg string) {
	u.inSection = true
	printTypeItemSuccess.Fprintln(u.out, msg)
}

func (u *uiImpl) ItemSkipped(msg string) {
	u.inSection = true
	printTypeItemSkipped.Fprintln(u.out, msg)
}

func (u *uiImpl) ItemFailure(msg string) {
	u.inSection = true
	printTypeItemFailure.Fprintln(u.out, msg)
}

func (u *uiImpl) Confirm(msg string) (bool, error) {
	ans, err := u.inputUI.Ask(fmt.Sprintf("%s [Y/n]", msg), &input.Options{
		HideOrder: true,
		Loop:      true,
		ValidateFunc: func(ans string) error {
			clog.Debug("receive user input", "query", msg, "input", ans)
			if ans != "Y" && ans != "n" {
				return fmt.Errorf("input must be Y or n")
			}
			return nil
		},
	})
	if err != nil {
		return false, err
	}
	return ans == "Y", nil
}
