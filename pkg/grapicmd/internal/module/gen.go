//go:generate mockgen -package=moduletesting -source=ui.go -destination=testing/ui_mock.go
//go:generate mockgen -package=moduletesting -source=generator.go -destination=testing/generator_mock.go
//go:generate mockgen -package=moduletesting -source=script.go -destination=testing/script_mock.go

package module
