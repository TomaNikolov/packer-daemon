package git

import (
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Clone ...
func Clone(url, path, username, password string) error {
	httpAuth := githttp.NewBasicAuth(username, password)
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth:     httpAuth,
	})

	return err
}

// Checkout ...
func Checkout(path string, branch string) error {
	r, err := git.PlainOpen(path + "/.")
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branch),
	})

	return err
}
