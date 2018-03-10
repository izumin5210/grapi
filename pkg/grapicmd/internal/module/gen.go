//go:generate mockgen -package=moduletesting -source=ui.go -destination=testing/ui_mock.go
//go:generate mockgen -package=moduletesting -imports=.=github.com/izumin5210/grapi/pkg/grapicmd/internal/module -source=generator.go -destination=testing/generator_mock.go

package module
