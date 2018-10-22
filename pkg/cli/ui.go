package cli

import (
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/fatih/color"
	"github.com/izumin5210/clicontrib/pkg/clog"
	"github.com/tcnksm/go-input"
)

type UI interface {
	Section(msg string)
	Subsection(msg string)
	ItemSuccess(msg string)
	ItemSkipped(msg string)
	ItemFailure(msg string)
	Confirm(msg string) (bool, error)
}

var (
	ui   UI
	uiMu sync.Mutex
)

func UIInstance(io *IO) UI {
	uiMu.Lock()
	defer uiMu.Unlock()
	if ui == nil {
		ui = NewUI(io)
	}
	return ui
}

// NewUI creates a new UI instance.
func NewUI(io *IO) UI {
	return &uiImpl{
		out: io.Out,
		inputUI: &input.UI{
			Reader: io.In,
			Writer: io.Out,
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

type fprintFunc func(w io.Writer, msg string)

type printConfig struct {
	prefix     string
	colorAttrs []color.Attribute
	indent     int
	allColor   bool
}

type printType int

const (
	printTypeUnknown printType = iota
	printTypeSection
	printTypeSubsection
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
