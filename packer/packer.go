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
func Build(cwd string, templatePath string, envVariables []string, options string, l *logger.Logger) error {
	commandArgs := []string{"build"}

	if len(options) > 0 {
		commandArgs = append(commandArgs, fmt.Sprintf("%s %s", "--var", options))
	}

	l.Log(strings.Join(envVariables, ", "))
	commandArgs = append(commandArgs, templatePath)
	logStreamerStdout := logger.NewLogstreamer(l, constants.Stdout, false)
	logStreamerStderr := logger.NewLogstreamer(l, constants.Stderr, false)
	env := os.Environ()
	env = append(env, envVariables...)

	cmd := exec.Command("packer", commandArgs...)
	cmd.Env = env
	cmd.Stdout = logStreamerStdout
	cmd.Stderr = logStreamerStderr
	cmd.Dir = cwd

	return cmd.Run()
}
