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
	err = git.Clone(config.Repository, pathToRepository, config.GitUsername, config.GitPassword)
	defer deleteRepository(pathToRepository)
	if err != nil {
		return err
	}

	logger.Log("Repository cloned!")

	// TODO: fix PlainOpen repository
	// err = git.Checkout(pathToRepository, buildRequest.Branch)
	// if err != nil {
	// 	fmt.Print("Checkout error " + pathToRepository)
	// 	fmt.Println(err)
	// 	return err
	// }

	pathToTemplate := path.Join(pathToRepository, constants.TemplateRelativePath, buildRequest.TemplateName)
	err = packer.Build(pathToTemplate, getEnvVariables(config), buildRequest.PackerOptions, &logger)

	return err
}

func deleteRepository(pathToRepository string) error {
	err := os.RemoveAll(pathToRepository)
	fmt.Println(err)
	return err
}

func getEnvVariables(config types.Config) string {
	// TODO: this shold return [] of env varibles
	return fmt.Sprintf(constants.PackerEnvTemplate, config.Username, config.Password, config.StoragePath)
}
