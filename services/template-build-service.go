package services

import (
	"fmt"
	"os"
	"path"

	"github.com/tomanikolov/packer-daemon/constants"
	"github.com/tomanikolov/packer-daemon/logger"
	"github.com/tomanikolov/packer-daemon/types"
	"github.com/tomanikolov/packer-daemon/utils"
)

// TemplateBuildService ...
type TemplateBuildService struct {
	buildRequest types.BuildRequest
	config       types.Config
	logger       logger.Logger
}

// NewTemplateBuildService ...
func NewTemplateBuildService(buildRequest types.BuildRequest, config types.Config, logger logger.Logger) TemplateBuildService {
	return TemplateBuildService{
		buildRequest: buildRequest,
		config:       config,
		logger:       logger,
	}
}

// Start ...
func (service *TemplateBuildService) Start() error {
	userDir, err := utils.GetUserDir()
	if err != nil {
		return err
	}

	pathToRepository := path.Join(userDir, constants.RepositoryName)
	service.logger.Log(fmt.Sprintf("Clonning repository %s", service.config.Repository))
	git := NewGitSErvice(service.config.Repository, service.config.GitUsername, service.config.GitPassword, &service.logger)
	err = git.Clone(pathToRepository)
	defer service.deleteRepository(pathToRepository)
	if err != nil {
		return err
	}

	service.logger.Log(fmt.Sprintf("Checkout bracnch: %s", service.buildRequest.Branch))
	err = git.Checkout(service.buildRequest.Branch)
	if err != nil {
		return err
	}

	pathToPackerDir := path.Join(pathToRepository, constants.TemplateRelativePath)
	pathToTemplate := path.Join(pathToPackerDir, service.buildRequest.TemplateName)
	service.logger.Log(fmt.Sprintf("Running pcaker build whit options : %s", service.buildRequest.PackerOptions))
	packerService := NewPackerService(service.logger)
	err = packerService.Build(pathToPackerDir, pathToTemplate, service.getEnvVariables(), service.buildRequest.PackerOptions)

	return err
}

func (service *TemplateBuildService) deleteRepository(pathToRepository string) error {
	err := os.RemoveAll(pathToRepository)
	return err
}

func (service *TemplateBuildService) getEnvVariables() []string {
	return []string{
		fmt.Sprintf(constants.UserNameEnv, service.config.Username),
		fmt.Sprintf(constants.PasswordEnv, service.config.Password),
		fmt.Sprintf(constants.StoragePathEnv, service.config.StoragePath),
		fmt.Sprintf(constants.GovcPasswordEnv, service.config.GovcPassword),
		fmt.Sprintf(constants.GovcUsernameEnv, service.config.GovcUsername),
		fmt.Sprintf(constants.GovcURLEnv, service.config.GovcURL),
		fmt.Sprintf(constants.GovcInsecureEnv, service.config.GovcInsecure),
		fmt.Sprintf(constants.GovcDataCenterEnv, service.config.GovcDataCenter),
		fmt.Sprintf(constants.GovcDataStoreEnv, service.config.GovcDataStore),
	}
}
