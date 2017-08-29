package git

import (
	"fmt"

	"github.com/tomanikolov/packer-daemon/constants"
	"github.com/tomanikolov/packer-daemon/logger"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Clone ...
func Clone(url, path, username, password string, l *logger.Logger) (*git.Worktree, error) {
	httpAuth := githttp.NewBasicAuth(username, password)
	logStreamerStdout := logger.NewLogstreamer(l, constants.Stdout, false)
	r, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: logStreamerStdout,
		Auth:     httpAuth,
	})

	if err != nil {
		return nil, err
	}

	return r.Worktree()
}

// Checkout ...
func Checkout(w *git.Worktree, branch string) error {
	err := w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Create: false,
		Force:  true,
	})

	return err
}
