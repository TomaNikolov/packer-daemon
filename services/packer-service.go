package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tomanikolov/packer-daemon/constants"
	"github.com/tomanikolov/packer-daemon/logger"
)

// PackerService ...
type PackerService struct {
	logger logger.Logger
}

// NewPackerService ...
func NewPackerService(logger logger.Logger) PackerService {
	return PackerService{
		logger: logger,
	}
}

// Build ...
func (service *PackerService) Build(cwd string, templatePath string, envVariables []string, options string) error {
	commandArgs := []string{"build"}

	if len(options) > 0 {
		commandArgs = append(commandArgs, fmt.Sprintf("%s %s", "--var", options))
	}

	service.logger.Log(strings.Join(envVariables, ", "))
	commandArgs = append(commandArgs, templatePath)
	logStreamerStdout := logger.NewLogstreamer(&service.logger, constants.Stdout, false)
	logStreamerStderr := logger.NewLogstreamer(&service.logger, constants.Stderr, false)
	env := os.Environ()
	env = append(env, envVariables...)

	cmd := exec.Command("packer", commandArgs...)
	cmd.Env = env
	cmd.Stdout = logStreamerStdout
	cmd.Stderr = logStreamerStderr
	cmd.Dir = cwd

	return cmd.Run()
}
