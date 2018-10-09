package clui

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
