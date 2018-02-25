package generate

import "github.com/izumin5210/grapi/pkg/grapicmd/ui"

type status int

const (
	statusCreate status = iota
	statusExist
	statusIdentical
	statusConflicted
	statusForce
	statusSkipped
)

var (
	nameByStatus = map[status]string{
		statusCreate:     "create",
		statusExist:      "exist",
		statusIdentical:  "identical",
		statusConflicted: "conflicted",
		statusForce:      "force",
		statusSkipped:    "skipped",
	}
	levelByStatus = map[status]ui.Level{
		statusCreate:     ui.LevelSuccess,
		statusExist:      ui.LevelInfo,
		statusIdentical:  ui.LevelInfo,
		statusConflicted: ui.LevelFail,
		statusForce:      ui.LevelWarn,
		statusSkipped:    ui.LevelWarn,
	}
	creatableStatusSet = map[status]struct{}{
		statusCreate: struct{}{},
		statusForce:  struct{}{},
	}
)

func (s status) String() string {
	return nameByStatus[s]
}

func (s status) Level() ui.Level {
	return levelByStatus[s]
}

func (s status) ShouldCreate() bool {
	_, ok := creatableStatusSet[s]
	return ok
}
