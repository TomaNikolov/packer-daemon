package services

import (
	"fmt"

	"github.com/tomanikolov/packer-daemon/constants"
	"github.com/tomanikolov/packer-daemon/logger"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// GitService ...
type GitService struct {
	httpAuth   *githttp.BasicAuth
	logger     *logger.Logger
	repository *git.Repository
	remoteURL  string
}

// NewGitSErvice ...
func NewGitSErvice(remoteURL, username, password string, l *logger.Logger) GitService {
	return GitService{
		httpAuth:  githttp.NewBasicAuth(username, password),
		logger:    l,
		remoteURL: remoteURL,
	}
}

// Clone ...
func (service *GitService) Clone(path string) error {
	logStreamerStdout := logger.NewLogstreamer(service.logger, constants.Stdout, false)
	r, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      service.remoteURL,
		Progress: logStreamerStdout,
		Auth:     service.httpAuth,
	})

	service.repository = r

	return err
}

// Checkout ...
func (service *GitService) Checkout(branch string) error {
	w, err := service.repository.Worktree()

	if err != nil {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", branch)),
		Create: false,
		Force:  true,
	})
}

// Fetch ...
func (service *GitService) Fetch() error {
	logStreamerStdout := logger.NewLogstreamer(service.logger, constants.Stdout, false)
	err := service.repository.Fetch(&git.FetchOptions{
		Auth:       service.httpAuth,
		Progress:   logStreamerStdout,
		RemoteName: "origin",
	})

	return err
}
