//go:generate go-assets-builder -p template -s="/init" -o init.go -v Init init
//go:generate go-assets-builder -p template -s="/service" -o service.go -v Service service
//go:generate go-assets-builder -p template -s="/command" -o command.go -v Command command

package template
