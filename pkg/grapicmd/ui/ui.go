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
	colored := color.New(st.Level().colorAttrs()...).SprintfFunc()
	fmt.Fprintf(u.out, "%s  %s\n", colored("%12s", st.String()), msg)
}
