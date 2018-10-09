package clui

import (
	"fmt"
	"io"

	"github.com/izumin5210/clicontrib/pkg/clog"
	"github.com/tcnksm/go-input"
)

// New creates a new UI instance.
func New(out io.Writer, in io.Reader) UI {
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
