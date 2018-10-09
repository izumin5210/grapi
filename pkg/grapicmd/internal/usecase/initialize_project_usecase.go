package usecase

import (
	"context"

	"github.com/izumin5210/gex"
	"github.com/pkg/errors"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

// InitializeProjectUsecase is an interface to create a new grapi project.
type InitializeProjectUsecase interface {
	Perform(rootDir, pkgName string, headUsed bool) error
	GenerateProject(rootDir, pkgName string, headUsed bool) error
	InstallDeps(rootDir string) error
}

// NewInitializeProjectUsecase creates a new InitializeProjectUsecase instance.
func NewInitializeProjectUsecase(ui module.UI, generator module.ProjectGenerator, commandFactory module.CommandFactory, gexCfg *gex.Config, version string) InitializeProjectUsecase {
	return &initializeProjectUsecase{
		ui:             ui,
		generator:      generator,
		commandFactory: commandFactory,
		gexCfg:         gexCfg,
		version:        version,
	}
}

type initializeProjectUsecase struct {
	ui             module.UI
	generator      module.ProjectGenerator
	commandFactory module.CommandFactory
	gexCfg         *gex.Config
	version        string
}

func (u *initializeProjectUsecase) Perform(rootDir, pkgName string, headUsed bool) error {
	u.ui.Section("Initialize project")

	var err error
	err = u.GenerateProject(rootDir, pkgName, headUsed)
	if err != nil {
		return errors.Wrap(err, "failed to initialize project")
	}

	u.ui.Subsection("Install dependencies")
	err = u.InstallDeps(rootDir)
	if err != nil {
		return errors.Wrap(err, "failed to execute `dep ensure`")
	}

	return nil
}

func (u *initializeProjectUsecase) GenerateProject(rootDir, pkgName string, headUsed bool) error {
	return errors.WithStack(u.generator.GenerateProject(rootDir, pkgName, module.ProjectGenerationConfig{UseHEAD: true}))
}

func (u *initializeProjectUsecase) InstallDeps(rootDir string) error {
	u.gexCfg.WorkingDir = rootDir
	repo, err := u.gexCfg.Create()
	if err == nil {
		err = repo.Add(
			context.TODO(),
			"github.com/izumin5210/grapi/cmd/grapi",
			// TODO: make configurable
			"github.com/golang/protobuf/protoc-gen-go",
			"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway",
			"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger",
		)
	}
	return errors.WithStack(err)
}
