package ui

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fatih/color"
)

type fprintFunc func(w io.Writer, msg string)

type printConfig struct {
	prefix     string
	colorAttrs []color.Attribute
	indent     int
	allColor   bool
}

type printType int

const (
	printTypeOutput printType = iota
	printTypeSection
	printTypeSubsection
	printTypeWarn
	printTypeError
	printTypeItemSuccess
	printTypeItemSkipped
	printTypeItemFailure
)

const (
	indentSizeItem = 4
)

var (
	configByPrintType = map[printType]printConfig{
		printTypeSection: {
			prefix:     "➜",
			colorAttrs: []color.Attribute{color.FgYellow},
			allColor:   true,
		},
		printTypeSubsection: {
			prefix:     "▸",
			colorAttrs: []color.Attribute{color.FgBlue},
			allColor:   true,
		},
		printTypeWarn: {
			prefix:     "⚠",
			colorAttrs: []color.Attribute{color.FgHiYellow},
			allColor:   true,
		},
		printTypeError: {
			prefix:     "☓",
			colorAttrs: []color.Attribute{color.FgHiRed},
			allColor:   true,
		},
		printTypeItemSuccess: {
			prefix:     "✔",
			colorAttrs: []color.Attribute{color.Bold, color.FgGreen},
			indent:     indentSizeItem,
		},
		printTypeItemSkipped: {
			prefix:     "╌",
			colorAttrs: []color.Attribute{color.Bold, color.FgBlue},
			indent:     indentSizeItem,
		},
		printTypeItemFailure: {
			prefix:     "✗",
			colorAttrs: []color.Attribute{color.Bold, color.FgRed},
			indent:     indentSizeItem,
		},
	}
	fprintlnFuncByPrintType = map[printType]fprintFunc{}
)

func init() {
	for pt, cfg := range configByPrintType {
		cfg := cfg
		fmtStr := "%s"
		if cfg.indent > 0 {
			fmtStr = "%" + strconv.FormatInt(int64(cfg.indent), 10) + "s"
		}
		if cfg.allColor {
			colored := color.New(cfg.colorAttrs...).FprintfFunc()
			fprintlnFuncByPrintType[pt] = func(w io.Writer, msg string) {
				colored(w, "  "+fmtStr+"  %s\n", cfg.prefix, msg)
			}
		} else {
			prefix := color.New(cfg.colorAttrs...).Sprintf(fmtStr, cfg.prefix)
			fprintlnFuncByPrintType[pt] = func(w io.Writer, msg string) {
				fmt.Fprintf(w, "  %s  %s\n", prefix, msg)
			}
		}
	}
}

func (pt printType) Fprintln(w io.Writer, msg string) {
	fprintlnFuncByPrintType[pt](w, msg)
}

// UI is an interface for intaracting with the terminal.
type UI interface {
	Output(msg string)
	Section(msg string)
	Subsection(msg string)
	Warn(msg string)
	Error(msg string)
	ItemSuccess(msg string)
	ItemSkipped(msg string)
	ItemFailure(msg string)
}

// New creates a new UI instance.
func New(out io.Writer) UI {
	return &uiImpl{
		out: out,
	}
}

type uiImpl struct {
	out       io.Writer
	inSection bool
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
