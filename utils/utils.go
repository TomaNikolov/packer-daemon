package utils

import (
	"log"
	"os/user"
)

// GetUserDir ...
func GetUserDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir, err
}
