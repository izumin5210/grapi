package ui

import "github.com/fatih/color"

// Level represents is an output priority.
type Level int

// Enum values of Level.
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
