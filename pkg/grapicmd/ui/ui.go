package ui

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

type Status interface {
	fmt.Stringer
	Level() Level
}

type UI interface {
	PrintWithStatus(msg string, status Status)
	PrintSuccess(msg, status string)
	PrintInfo(msg, status string)
	PrintWarn(msg, status string)
	PrintFail(msg, status string)
}

func New(out io.Writer) UI {
	return &uiImpl{
		out: out,
	}
}

type uiImpl struct {
	out io.Writer
}

func (u *uiImpl) PrintWithStatus(msg string, st Status) {
	u.printWithLevel(msg, st.String(), st.Level())
}

func (u *uiImpl) PrintSuccess(msg, status string) {
	u.printWithLevel(msg, status, LevelSuccess)
}

func (u *uiImpl) PrintInfo(msg, status string) {
	u.printWithLevel(msg, status, LevelInfo)
}

func (u *uiImpl) PrintWarn(msg, status string) {
	u.printWithLevel(msg, status, LevelWarn)
}

func (u *uiImpl) PrintFail(msg, status string) {
	u.printWithLevel(msg, status, LevelFail)
}

func (u *uiImpl) printWithLevel(msg, status string, lv Level) {
	colored := color.New(lv.colorAttrs()...).SprintfFunc()
	fmt.Fprintf(u.out, "%s  %s\n", colored("%12s", status), msg)
}
