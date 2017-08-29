package builder

import (
	"fmt"
	"os"
	"path"

	"github.com/tomanikolov/packer-daemon/constants"
	"github.com/tomanikolov/packer-daemon/git"
	"github.com/tomanikolov/packer-daemon/logger"
	"github.com/tomanikolov/packer-daemon/packer"
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
	workTree, err := git.Clone(config.Repository, pathToRepository, config.GitUsername, config.GitPassword, &logger)
	defer deleteRepository(pathToRepository)
	if err != nil {
		return err
	}

	logger.Log(fmt.Sprintf("Checkout bracnch: %s", buildRequest.Branch))
	err = git.Checkout(workTree, buildRequest.Branch)
	if err != nil {
		return err
	}

	pathToTemplate := path.Join(pathToRepository, constants.TemplateRelativePath, buildRequest.TemplateName)
	logger.Log(fmt.Sprintf("Running pcaker build whit options : %s", buildRequest.PackerOptions))
	err = packer.Build(pathToTemplate, getEnvVariables(config), buildRequest.PackerOptions, &logger)

	return err
}

func deleteRepository(pathToRepository string) error {
	err := os.RemoveAll(pathToRepository)
	fmt.Println(err)
	return err
}

func getEnvVariables(config types.Config) []string {
	return []string{
		fmt.Sprintf(constants.UserNameEnv, config.Username),
		fmt.Sprintf(constants.PasswordEnv, config.Password),
		fmt.Sprintf(constants.StoragePathEnv, config.StoragePath),
	}
}
