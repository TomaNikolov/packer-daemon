package packer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tomanikolov/packer-daemon/constants"
	"github.com/tomanikolov/packer-daemon/logger"
)

// Build ...
func Build(templatePath string, envVariables []string, options string, l *logger.Logger) error {
	commandArgs := []string{"build"}
	if len(options) > 0 {
		commandArgs = append(commandArgs, fmt.Sprintf("%s %s", "--var", options))
	}

	logStreamerStdout := logger.NewLogstreamer(l, constants.Stdout, false)
	logStreamerStderr := logger.NewLogstreamer(l, constants.Stderr, false)

	commandArgs = append(commandArgs, templatePath)

	cmd := exec.Command("packer", commandArgs...)
	env := os.Environ()
	env = append(env, envVariables...)
	cmd.Env = env
	l.Log(strings.Join(env, ", "))
	cmd.Stdout = logStreamerStdout
	cmd.Stderr = logStreamerStderr
	err := cmd.Run()

	return err
}
