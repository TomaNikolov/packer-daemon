package packer

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tomanikolov/packer-daemon/logger"
)

// Build ...
func Build(templatePath, envVariables, options string, l *logger.Logger) error {
	commandArgs := []string{"build"}
	if len(options) > 0 {
		commandArgs = append(commandArgs, fmt.Sprintf("%s %s", "--var", options))
	}

	logStreamer := logger.NewLogstreamer(l, "stdout", false)

	commandArgs = append(commandArgs, templatePath)

	cmd := exec.Command("packer", commandArgs...)
	env := os.Environ()
	env = append(env, envVariables)
	cmd.Env = env
	fmt.Println(env)
	cmd.Stdout = logStreamer
	//cmd.Stdout = os.Stdout
	err := cmd.Run()

	return err
}
