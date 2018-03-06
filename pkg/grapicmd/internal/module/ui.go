package module

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
	Confirm(msg string) (bool, error)
}
