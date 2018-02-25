package ui

import "github.com/fatih/color"

type Level int

const (
	LevelSuccess Level = iota
	LevelInfo
	LevelWarn
	LevelFail
)

var (
	colorByLevel = map[Level][]color.Attribute{
		LevelSuccess: {color.Bold, color.FgGreen},
		LevelInfo:    {color.Bold, color.FgBlue},
		LevelWarn:    {color.Bold, color.FgYellow},
		LevelFail:    {color.Bold, color.FgRed},
	}
)

func (l Level) colorAttrs() []color.Attribute {
	return colorByLevel[l]
}
