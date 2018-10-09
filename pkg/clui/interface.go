package clui

type UI interface {
	Section(msg string)
	Subsection(msg string)
	ItemSuccess(msg string)
	ItemSkipped(msg string)
	ItemFailure(msg string)
	Confirm(msg string) (bool, error)
}
