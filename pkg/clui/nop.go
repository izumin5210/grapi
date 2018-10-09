package clui

var Nop = &nopImpl{}

type nopImpl struct{}

func (*nopImpl) Section(msg string)               {}
func (*nopImpl) Subsection(msg string)            {}
func (*nopImpl) ItemSuccess(msg string)           {}
func (*nopImpl) ItemSkipped(msg string)           {}
func (*nopImpl) ItemFailure(msg string)           {}
func (*nopImpl) Confirm(msg string) (bool, error) { return true, nil }
