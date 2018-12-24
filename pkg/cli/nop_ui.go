package cli

var NopUI = &nopUIImpl{}

type nopUIImpl struct{}

func (*nopUIImpl) Section(msg string)                    {}
func (*nopUIImpl) Subsection(msg string)                 {}
func (*nopUIImpl) ItemSuccess(msg string)                {}
func (*nopUIImpl) ItemSkipped(msg string)                {}
func (*nopUIImpl) ItemFailure(msg string, errs ...error) {}
func (*nopUIImpl) Confirm(msg string) (bool, error)      { return true, nil }
