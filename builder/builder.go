package builder

import (
	"fmt"
	"os"
	"path"

	"github.com/tomanikolov/packer-daemon/constants"
	"github.com/tomanikolov/packer-daemon/logger"
	"github.com/tomanikolov/packer-daemon/packer"
	"github.com/tomanikolov/packer-daemon/services"
	"github.com/tomanikolov/packer-daemon/types"
	"github.com/tomanikolov/packer-daemon/utils"
)

// Start ...
func Start(buildRequest types.BuildRequest, config types.Config, logger logger.Logger) error {
	userDir, err := utils.GetUserDir()
	if err != nil {
		return err
	}

	pathToRepository := path.Join(userDir, constants.RepositoryName)
	logger.Log(fmt.Sprintf("Clonning repository %s", config.Repository))
	git := services.NewGitSErvice(config.Repository, config.GitUsername, config.GitPassword, &logger)
	err = git.Clone(pathToRepository)
	defer deleteRepository(pathToRepository)
	if err != nil {
		return err
	}

	logger.Log(fmt.Sprintf("Checkout bracnch: %s", buildRequest.Branch))
	err = git.Checkout(buildRequest.Branch)
	if err != nil {
		return err
	}

	pathToPackerDir := path.Join(pathToRepository, constants.TemplateRelativePath)
	pathToTemplate := path.Join(pathToPackerDir, buildRequest.TemplateName)
	logger.Log(fmt.Sprintf("Running pcaker build whit options : %s", buildRequest.PackerOptions))
	err = packer.Build(pathToPackerDir, pathToTemplate, getEnvVariables(config), buildRequest.PackerOptions, &logger)

	return err
}

func deleteRepository(pathToRepository string) error {
	err := os.RemoveAll(pathToRepository)
	return err
}

func getEnvVariables(config types.Config) []string {
	return []string{
		fmt.Sprintf(constants.UserNameEnv, config.Username),
		fmt.Sprintf(constants.PasswordEnv, config.Password),
		fmt.Sprintf(constants.StoragePathEnv, config.StoragePath),
		fmt.Sprintf(constants.GovcPasswordEnv, config.GovcPassword),
		fmt.Sprintf(constants.GovcUsernameEnv, config.GovcUsername),
		fmt.Sprintf(constants.GovcURLEnv, config.GovcURL),
		fmt.Sprintf(constants.GovcInsecureEnv, config.GovcInsecure),
	}
}
