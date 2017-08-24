package builder

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/tomanikolov/packer-daemon/constants"
	"github.com/tomanikolov/packer-daemon/git"
	"github.com/tomanikolov/packer-daemon/packer"
	"github.com/tomanikolov/packer-daemon/types"
	"github.com/tomanikolov/packer-daemon/utils"
)

// Start ...
func Start(buildMessage string, config types.Config) error {
	buildRequest, err := getBuildRequest(buildMessage)
	if err != nil {
		return err
	}

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

	// TODO: fix PlainOpen repository
	// err = git.Checkout(pathToRepository, buildRequest.Branch)
	// if err != nil {
	// 	fmt.Print("Checkout error " + pathToRepository)
	// 	fmt.Println(err)
	// 	return err
	// }

	pathToTemplate := path.Join(pathToRepository, constants.TemplateRelativePath, buildRequest.TemplateName)
	err = packer.Build(pathToTemplate, getEnvVariables(config), buildRequest.PackerOptions)

	return err
}

func getBuildRequest(buildMessage string) (types.BuildRequest, error) {
	buildRequest := types.BuildRequest{}
	err := json.Unmarshal([]byte(buildMessage), &buildRequest)
	return buildRequest, err
}

func deleteRepository(pathToRepository string) error {
	fmt.Println("remove " + pathToRepository)
	err := os.RemoveAll(pathToRepository)
	fmt.Println(err)
	return err
}

func getEnvVariables(config types.Config) string {
	return fmt.Sprintf(constants.PackerEnvTemplate, config.Username, config.Password, config.StoragePath)
}
