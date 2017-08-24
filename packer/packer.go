package packer

import (
	"fmt"
	"os"
	"os/exec"
)

// Build ...
func Build(templatePath, envVariables, options string) error {
	commandArgs := []string{"build"}
	if len(options) > 0 {
		commandArgs = append(commandArgs, fmt.Sprintf("%s %s", "--var", options))
	}

	commandArgs = append(commandArgs, templatePath)

	cmd := exec.Command("packer", commandArgs...)
	env := os.Environ()
	env = append(env, envVariables)
	cmd.Env = env
	fmt.Println(env)
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	return err
}
